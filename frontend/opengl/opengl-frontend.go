package openglfrontend

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v3.2-compatibility/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	window        *glfw.Window
	shaderProgram uint32
	screenTexture uint32
	vao           uint32
	vertices      []float32
	eboIndices    []uint32
)

func init() {
	// ensure that this runs on main thread
	runtime.LockOSThread()

	// Initialize GLFW
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	window, err = glfw.CreateWindow(640, 320, "Gameboy", nil, nil)
	if err != nil {
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
	#version 300 es

	layout (location=0) in vec3 aPos;
	layout (location=1) in vec2 aTexCoord;

	out vec2 TexCoord;

	void main()
	{
		gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0f);
		TexCoord = aTexCoord;
	}
	`)

	fragShaderSource := []byte(`
	#version 300 es

	precision mediump float;
	
	out vec4 FragColor;
	in vec2 TexCoord;
	
	uniform sampler2D texture1;
	
	void main()
	{
		FragColor = texture(texture1, TexCoord);
	}
	`)

	// opengl initialization
	// ====================================
	err = gl.Init()
	if err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Printf("OpenGL version: %s\n", version)

	gl.Viewport(0, 0, 640, 320)

	// link vertex and fragment shaders into shader program
	// and use it for rendering
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragShader, err := compileShader(fragShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
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

	// VERTICES
	// --------------------
	// create our vertices. These will be a simple rect (2 triangles)
	// covering the 2 screens -- a canvas for us to render our screen 'texture'
	// to.
	vertices = []float32{
		// x 	y		z  		s   t
		// -------------------------
		-1.0, 1.0, 0.0, 0.0, 1.0, // top-left
		-1.0, -1.0, 0.0, 0.0, 0.0, // bottom-left
		1.0, 1.0, 0.0, 1.0, 1.0, // top-right
		1.0, -1.0, 0.0, 1.0, 0.0, // bottom-right
	}
	eboIndices = []uint32{
		0, 1, 3, // first triangle
		0, 3, 2, // second triangle
	}

	// initialize a vertex array object -- this is a convenience object that
	// stores information about how the currently bound vertex buffer object is
	// configured. The information that the vertex array object stores is the
	// info about the shape of the VBO that we set in gl.VertexAttribPointer below.
	gl.GenVertexArrays(1, &vao)
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	var ebo uint32
	gl.GenBuffers(1, &ebo)

	// set our new vertex array object as the active VAO
	gl.BindVertexArray(vao)

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

	// load the indices of the vertices we want to draw into the element buffer object
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(
		gl.ELEMENT_ARRAY_BUFFER, // load into current element buffer object
		4*len(eboIndices),       // total number of bytes to be loaded (each uint32 is 4 bytes wide)
		gl.Ptr(eboIndices),      // openGL pointer to array of indices
		gl.STATIC_DRAW,          // hint to openGL that we won't be changing these indices often at all
	)

	// tell openGL about the shape of our vertex buffer
	gl.VertexAttribPointer(
		0,        // configure the vertex attribute with id 0 (location)
		3,        // each vertex attribute is made of three components (in this case, xyz coordinates)
		gl.FLOAT, // each component is a 32bit float
		false,    // there are no delimiters between each ser of components in the array (array is tightly packed)
		5*4,      // the span of bytes of one vertex attribute is 5 float32s (3 for location attrib, 2 for texel coordinate attrib). Each float32 is 4 bytes.
		nil,      // the offset of the first vertex attribute in the array is zero. For some reason, this requires a void pointer cast, represented in go-gl as nil.
	)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(
		1,                 // configure the vertex attribute with id 1 (texture coordinates)
		2,                 // each vertex attribute is made of two components (in this case, st texture coordinates)
		gl.FLOAT,          // each component is a 32bit float
		false,             // there are no delimiters between each ser of components in the array (array is tightly packed)
		5*4,               // the span of bytes of one vertex attribute is 5 float32s (3 for location attrib, 2 for texel coordinate attrib). Each float32 is 4 bytes.
		gl.PtrOffset(3*4), // the offset of the first vertex attribute in the array is 12.
	)
	gl.EnableVertexAttribArray(1)
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
	// =====================================

	// 'handle' errors
	if err := gl.GetError(); err != gl.NO_ERROR {
		panic(err)
	}
}

func ConnectVideo(screens <-chan []byte) {
	for screen := range screens {
		render(screen, screenTexture, shaderProgram, vao, window)
	}
}

// Render takes a `screen`, which is a slice of bytes which represent
// 160x144 RGB pixels.
func render(screen []byte, screenTexture, shaderProgram, vao uint32, window *glfw.Window) {

	gl.ClearColor(0.1, 0.2, 0.1, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

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
		gl.Ptr(screen),   // data
	)

	// use our screen shader program
	gl.UseProgram(shaderProgram)
	// draw the vertices in our vertex array object as triangles
	// containing our screen-sized rectangle
	gl.BindVertexArray(vao)
	numVerticesToDraw := int32(6)
	gl.DrawElements(gl.TRIANGLES, numVerticesToDraw, gl.UNSIGNED_INT, gl.PtrOffset(0))

	// render screen
	window.SwapBuffers()

	// 'handle' errors
	if err := gl.GetError(); err != gl.NO_ERROR {
		panic(err)
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
