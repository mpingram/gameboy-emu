package openglfrontend

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/mpingram/gameboy-emu/ppu"
)

const scale = 3
const width = 160
const height = 144

// gbColorsToRGB converts the screen data from the gameboy, which are single-byte
// color identifiers, into RGB pixels that can be fed to openGL.
func gbColorsToRGB(colors []ppu.Color) []byte {
	pixels := make([]byte, 0)
	for _, c := range colors {
		switch c {
		case ppu.White:
			pixels = append(pixels, 255, 255, 255)
		case ppu.LightGray:
			pixels = append(pixels, 151, 150, 149)
		case ppu.DarkGray:
			pixels = append(pixels, 76, 75, 74)
		case ppu.Black:
			pixels = append(pixels, 0, 0, 0)
		default:
			panic(fmt.Sprintf("toRGB: Got bad color: %v", c))
		}
	}
	return pixels
}

var (
	window        *glfw.Window
	shaderProgram uint32
	screenTexture uint32
	vao           uint32
	vertices      []float32
	eboIndices    []uint32
)

func ConnectVideo(screens <-chan []ppu.Color) {
	// ensure that this runs on main thread
	runtime.LockOSThread()

	// Initialize GLFW
	err := glfw.Init()
	if err != nil {
		fmt.Println("Error initializing GLFW")
		panic(err)
	}
	defer glfw.Terminate()
	w, err := glfw.CreateWindow(160*scale, 144*scale, "Gameboy", nil, nil)
	window = w
	if err != nil {
		fmt.Println("Error creating GLFW window")
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window.MakeContextCurrent()

	// Initialize OpenGL
	vertexShaderSource := []byte(`
	attribute vec3 a_position;
	attribute vec2 a_texcoord;

	varying vec2 v_texcoord;

	void main() {
		gl_Position = vec4(a_position, 1.0);
		v_texcoord = a_texcoord;
	}
	`)

	fragShaderSource := []byte(`
	varying vec2 v_texcoord;
	uniform sampler2D u_texture;

	void main() {
		gl_FragColor = texture2D(u_texture, v_texcoord);
	}
	`)

	// opengl initialization
	// ====================================
	err = gl.Init()
	if err != nil {
		fmt.Println("Error initializing openGL")
		panic(err)
	}
	// DEBUG
	// gl.PolygonMode(gl.FRONT, gl.LINE)
	// END DEBUG
	// version := gl.GoStr(gl.GetString(gl.VERSION))
	// log.Printf("OpenGL version: %s\n", version)

	gl.Viewport(0, 0, 160, 144)

	// link vertex and fragment shaders into shader program
	// and use it for rendering
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		fmt.Println("Error compiling vertex shader")
		panic(err)
	}
	fragShader, err := compileShader(fragShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		fmt.Println("Error compiling fragment shader")
		panic(err)
	}

	shaderProgram = gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragShader)
	gl.LinkProgram(shaderProgram)
	// check for linking errors
	var status int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		logLength := int32(512)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))
		panic(log)
	}
	// free our shaders once we've linked them into a shader program
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragShader)

	checkGLErr()

	// VERTICES
	// --------------------
	// create our vertices. These will be a simple rect (2 triangles)
	// covering the 2 screens -- a canvas for us to render our screen 'texture'
	// to.
	vertices = []float32{
		// x 	y		z  		s   t
		// -------------------------
		-0.5, 0.5, 0.0, 0.0, 1.0, // top-left
		-0.5, -0.5, 0.0, 0.0, 0.0, // bottom-left
		0.5, 0.5, 0.0, 1.0, 1.0, // top-right
		0.5, -0.5, 0.0, 1.0, 0.0, // bottom-right
	}
	eboIndices = []uint32{
		0, 1, 3, // first triangle
		0, 3, 2, // second triangle
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	checkGLErr()
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	checkGLErr()

	// load our vertices into our vertex buffer object
	// gl.BindBuffer call sets vbo as the active vertex buffer; now things we configure
	// in gl.VertexAttribPointer will be stored to the bound VAO for later use.
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,  // load into the current array buffer
		4*len(vertices),  // total number of bytes in the array to be loaded (each float32 is 4 bytes wide)
		gl.Ptr(vertices), // openGL pointer to the array of vertices
		gl.STATIC_DRAW,   // hint to openGL that we won't be changing these vertices often at all
	)
	checkGLErr()

	// load the indices of the vertices we want to draw into the element buffer object
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(
		gl.ELEMENT_ARRAY_BUFFER, // load into current element buffer object
		4*len(eboIndices),       // total number of bytes to be loaded (each uint32 is 4 bytes wide)
		gl.Ptr(eboIndices),      // openGL pointer to array of indices
		gl.STATIC_DRAW,          // hint to openGL that we won't be changing these indices often at all
	)
	checkGLErr()

	// enable vertex attribute coordinates
	positionLoc := gl.GetAttribLocation(shaderProgram, gl.Str("a_position\x00"))
	// tell openGL about the shape of our vertex buffer
	gl.VertexAttribPointer(
		uint32(positionLoc), // configure the vertex attribute with id 0 (location)
		3,                   // each vertex attribute is made of three components (in this case, xyz coordinates)
		gl.FLOAT,            // each component is a 32bit float
		false,               // there are no delimiters between each attribute in the array (array is tightly packed)
		5*4,                 // the span of bytes of one vertex attribute is 5 float32s (3 for location attrib, 2 for texel coordinate attrib). Each float32 is 4 bytes.
		nil,                 // the offset of the first vertex attribute in the array is zero. For some reason, this requires a void pointer cast, represented in go-gl as nil.
	)
	gl.EnableVertexAttribArray(uint32(positionLoc))

	checkGLErr()

	textureLoc := gl.GetAttribLocation(shaderProgram, gl.Str("a_texcoord\x00"))
	gl.VertexAttribPointer(
		uint32(textureLoc), // configure the vertex attribute with id 1 (texture coordinates)
		2,                  // each vertex attribute is made of two components (in this case, st texture coordinates)
		gl.FLOAT,           // each component is a 32bit float
		false,              // there are no delimiters between each ser of components in the array (array is tightly packed)
		5*4,                // the span of bytes of one vertex attribute is 5 float32s (3 for location attrib, 2 for texel coordinate attrib). Each float32 is 4 bytes.
		gl.PtrOffset(3*4),  // the offset of the first vertex attribute in the array is 12.
	)
	gl.EnableVertexAttribArray(uint32(textureLoc))

	checkGLErr()
	// ----------------------------

	// create our screen texture
	gl.GenTextures(1, &screenTexture)
	gl.BindTexture(gl.TEXTURE_2D, screenTexture)
	// clamp texture to border: do not render texels outside of texture coordinate area
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)
	// use nearest neighbor texture filtering when zooming up; it's pixel graphics, let's keep it blocky
	// when zooming down, use bilinear filtering
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	checkGLErr()

	// create texture data from initial Chip8 screen
	var emptyScreen []byte
	for row := 0; row < 144; row++ {
		for col := 0; col < 160; col++ {
			emptyScreen = append(emptyScreen, 0, 0, 100)
		}
	}

	gl.TexImage2D(
		gl.TEXTURE_2D, // target the 2D texture
		0,             // mipmap level 0
		gl.RGB,        // internal texture format; a subtype of the texture format parameter above. gl.R8 is A single 'red' channel represented by one byte.
		160,
		144,
		0,                   // fun fact, this parameter apparently is 'border', and it must always be 0 or else!
		gl.RGB,              // texture format; gl.RED is a single red channel.
		gl.UNSIGNED_BYTE,    // type; gl.R8 requires an unsigned byte. Some other internal texture formats can have different types, that's why this is here.
		gl.Ptr(emptyScreen), // last but not least, the pixel data of the texture
	)

	// set our texture uniform in our shader to our texture (NOTE why 0 and not texture id?)
	gl.UseProgram(shaderProgram)
	texUniform := gl.GetUniformLocation(shaderProgram, gl.Str("texture1\000"))
	gl.Uniform1i(texUniform, 0)
	checkGLErr()
	// =====================================

	// Main loop: render screens. Does not return -- any calling code must
	// use a separate goroutine to do work.
	for screen := range screens {
		// convert screen to RGB pixels
		pixels := gbColorsToRGB(screen)
		gl.ClearColor(0.9, 0.9, 0.7, 1.0) // gross pale yellow
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(shaderProgram)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, screenTexture)

		// replace the current texture with new texture
		gl.TexSubImage2D(
			gl.TEXTURE_2D,
			0,                // mipmap level 0
			0,                // x offset
			0,                // y offset
			160,              // width
			144,              // height
			gl.RGB,           // format
			gl.UNSIGNED_BYTE, // type,
			gl.Ptr(pixels),   // data
		)
		checkGLErr()

		numVerticesToDraw := int32(6)
		gl.DrawElements(gl.TRIANGLES, numVerticesToDraw, gl.UNSIGNED_INT, gl.PtrOffset(0))

		window.SwapBuffers()
		glfw.PollEvents()
		// Break out of loop on esc keypress
		if k := window.GetKey(glfw.KeyEscape); k == glfw.Press {
			break
		}
	}
}

func checkGLErr() {
	// 'handle' errors
	for err := gl.GetError(); err != gl.NO_ERROR; err = gl.GetError() {
		fmt.Printf("Encountered openGL error %v\n", err)
	}
}

func compileShader(sourceBytes []byte, shaderType uint32) (uint32, error) {
	sourceStr := string(sourceBytes) + string('\000')
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(sourceStr)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile shader:\n%v\n%v", sourceStr, log)
	}

	return shader, nil
}
