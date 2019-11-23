package frontend

import (
	"bytes"
	"encoding/binary"
	"syscall/js"
)

// https://www.tutorialspoint.com/webgl/webgl_quick_guide.htm

var (
	width   int
	height  int
	gl      js.Value
	glTypes GLTypes
)

type GLTypes struct {
	staticDraw         js.Value
	arrayBuffer        js.Value
	elementArrayBuffer js.Value
	vertexShader       js.Value
	fragmentShader     js.Value
	float              js.Value
	depthTest          js.Value
	colorBufferBit     js.Value
	triangles          js.Value
	unsignedShort      js.Value
}

func (types *GLTypes) New() {
	types.staticDraw = gl.Get("STATIC_DRAW")
	types.arrayBuffer = gl.Get("ARRAY_BUFFER")
	types.elementArrayBuffer = gl.Get("ELEMENT_ARRAY_BUFFER")
	types.vertexShader = gl.Get("VERTEX_SHADER")
	types.fragmentShader = gl.Get("FRAGMENT_SHADER")
	types.float = gl.Get("FLOAT")
	types.depthTest = gl.Get("DEPTH_TEST")
	types.colorBufferBit = gl.Get("COLOR_BUFFER_BIT")
	types.triangles = gl.Get("TRIANGLES")
	types.unsignedShort = gl.Get("UNSIGNED_SHORT")
}

func alert(val interface{}) {
	js.Global().Call("alert", val)
}

func log(val interface{}) {
	js.Global().Get("console").Call("log", val)
}

func toUint8Array(val interface{}) js.Value {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, val)
	if err != nil {
		panic(err)
	}
	// create a js UInt8Array
	typedArray := js.Global().Get("Uint8Array").New(len(buf.Bytes()))
	js.CopyBytesToJS(typedArray, buf.Bytes())
	return typedArray
}

func toFloat32Array(val interface{}) js.Value {
	uint8Array := toUint8Array(val)
	// represent bytes in uint8Array (TypedArray.buffer) as float32Array instead
	// https://stackoverflow.com/questions/46196990/convert-float-values-to-uint8-array-in-javascript?rq=1
	float32Array := js.Global().Get("Float32Array").New(uint8Array.Get("buffer"))
	return float32Array
}

func toUint16Array(val interface{}) js.Value {
	uint8Array := toUint8Array(val)
	// represent bytes in uint8Array (TypedArray.buffer) as uint16Array instead
	uint16Array := js.Global().Get("Uint16Array").New(uint8Array.Get("buffer"))
	return uint16Array
}

func DrawTriangle() {

	// Init Canvas stuff
	doc := js.Global().Get("document")
	canvasEl := doc.Call("getElementById", "gocanvas")
	width = doc.Get("body").Get("clientWidth").Int()
	height = doc.Get("body").Get("clientHeight").Int()
	canvasEl.Set("width", width)
	canvasEl.Set("height", height)

	gl = canvasEl.Call("getContext", "webgl")
	if gl == js.Undefined() {
		gl = canvasEl.Call("getContext", "experimental-webgl")
	}
	// once again
	if gl == js.Undefined() {
		js.Global().Call("alert", "browser might not support webgl")
		return
	}

	glTypes.New()

	//// VERTEX BUFFER ////

	// our vertices are organized
	// with three spatial coordinates (xyz)
	// and two texture coordinates (st)
	verticesNative := []float32{
		// x		y		z		s		t
		-0.5, 0.5, 0, 0, 1, // top left
		-0.5, -0.5, 0, 0, 0, // bottom left
		0.5, 0.5, 0, 1, 1, // top right
		0.5, -0.5, 0, 1, 0, // bottom right
	}
	vertices := toFloat32Array(verticesNative)

	indicesNative := []uint16{
		0, 1, 3, // first triangle
		0, 3, 2, // second triangle
	}
	indices := toUint16Array(indicesNative)

	// Create vertex buffer
	vertexBuffer := gl.Call("createBuffer", glTypes.arrayBuffer)
	// Bind to vertex buffer
	gl.Call("bindBuffer", glTypes.arrayBuffer, vertexBuffer)
	// Pass vertices to vertex buffer
	gl.Call("bufferData", glTypes.arrayBuffer, vertices, glTypes.staticDraw)
	// Unbind vertex buffer
	gl.Call("bindBuffer", glTypes.arrayBuffer, nil)

	// Create index buffer
	indexBuffer := gl.Call("createBuffer", glTypes.elementArrayBuffer)
	// Bind to index buffer
	gl.Call("bindBuffer", glTypes.elementArrayBuffer, indexBuffer)
	// Pass data to index buffer
	gl.Call("bufferData", glTypes.elementArrayBuffer, indices, glTypes.staticDraw)
	// Unbind index buffer
	gl.Call("bindBuffer", glTypes.elementArrayBuffer, nil)

	//// Shaders ////

	// Vertex shader source code
	vertCode := `
	attribute vec3 a_position;
	attribute vec2 a_texcoord;

	varying vec2 v_texcoord;
		
	void main() {
		gl_Position = vec4(a_position, 1.0);
		v_texcoord = a_texcoord;
	}`

	// Create a vertex shader object
	vertShader := gl.Call("createShader", glTypes.vertexShader)
	// Attach vertex shader source code
	gl.Call("shaderSource", vertShader, vertCode)
	// Compile the vertex shader
	gl.Call("compileShader", vertShader)

	//fragment shader source code
	fragCode := `
	varying highp vec2 v_texcoord;
	uniform sampler2D u_texture;

	void main() {
		gl_FragColor = texture2D(u_texture, v_texcoord);
	}`

	// Create fragment shader object
	fragShader := gl.Call("createShader", glTypes.fragmentShader)

	// Attach fragment shader source code
	gl.Call("shaderSource", fragShader, fragCode)
	// Compile the fragmentt shader
	gl.Call("compileShader", fragShader)
	// Create a shader program object to store
	// the combined shader program
	shaderProgram := gl.Call("createProgram")

	// Attach a vertex shader
	gl.Call("attachShader", shaderProgram, vertShader)
	// Attach a fragment shader
	gl.Call("attachShader", shaderProgram, fragShader)
	// Link both the programs
	gl.Call("linkProgram", shaderProgram)
	// Use the combined shader program object
	gl.Call("useProgram", shaderProgram)

	//// Associating shaders to buffer objects ////
	// Our vertices have two attributes in them: the spatial coordinates (xyz)
	// and the texture coordinates (st). Both of these attributes are packed into
	// the same vertex buffer, so we need to tell openGL how to find these
	// two attributes using `vertexAttribPointer`

	// Bind vertex buffer object
	gl.Call("bindBuffer", glTypes.arrayBuffer, vertexBuffer)
	// Bind index buffer object
	gl.Call("bindBuffer", glTypes.elementArrayBuffer, indexBuffer)
	// Get the `position` attribute location
	position := gl.Call("getAttribLocation", shaderProgram, "a_position")
	// Configure the `coordinates` attribute pointer
	gl.Call(
		"vertexAttribPointer",
		position,      // attribute location
		3,             // number of elements in attribute
		glTypes.float, // type
		false,         // not normalized (clamp to [-1,1], not [0,1])
		5*4,           // stride: number of bytes between different vertex data elements
		0,             // offset of first element (in bytes)
	)
	// Enable the `position` attribute
	gl.Call("enableVertexAttribArray", position)

	// Get the `texcoord` attribute
	texcoord := gl.Call("getAttribLocation", shaderProgram, "a_texcoord")
	// Configure the `texcoord` attribute pointer
	gl.Call(
		"vertexAttribPointer",
		texcoord,      // attribute location
		2,             // number of elements in attribute
		glTypes.float, // type
		false,         // not normalized (clamp to [-1,1], not [0,1])
		5*4,           // stride: number of bytes between different vertex data elements
		3*4,           // offset of first element (in bytes): it is 3 floats from the start
	)
	// Enable the `position` attribute
	gl.Call("enableVertexAttribArray", texcoord)

	//// Create a standby texture ////
	texture := gl.Call("createTexture")
	standbyTexturePixel := toUint8Array([]uint8{99, 200, 77})
	gl.Call("bindTexture", gl.Get("TEXTURE_2D"), texture)
	// NOTE -- WebGL 1.0 requires special treatment of textures that
	// aren't powers of 2 (need to look into this more.) In any case,
	// the GB screen texture (160x144) isn't a power of 2 in either dimension,
	// so we'll need to do the special treatment thing no matter what.
	gl.Call(
		"texImage2D",
		gl.Get("TEXTURE_2D"),    // target
		0,                       // level
		gl.Get("RGB"),           // internal format
		1,                       // width
		1,                       // height
		0,                       // border
		gl.Get("RGB"),           // src format -- in WebGL 1, must be same as internal format
		gl.Get("UNSIGNED_BYTE"), // 1 byte per color channel (RGB), 3 bytes per pixel
		standbyTexturePixel,     // pixel data
	)

	//// Draw the screen ////
	// Clear the canvas
	gl.Call("clearColor", 0.5, 0.5, 0.5, 0.9)
	gl.Call("clear", glTypes.colorBufferBit)
	// Enable the depth test
	gl.Call("enable", glTypes.depthTest)
	// Set the view port
	gl.Call("viewport", 0, 0, width, height)
	// Draw the triangle
	gl.Call("drawElements", glTypes.triangles, indices.Get("length"), glTypes.unsignedShort, 0)
}
