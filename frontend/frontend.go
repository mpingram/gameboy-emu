package frontend

import (
	"bytes"
	"encoding/binary"
	"syscall/js"
)

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

	// We need to convert these native Golang floats to
	// js typed arrays in order for webgl to consume them.
	// The way we do this is by converting the floats to bytes
	// and then calling js.CopyBytesToJS to get our typedArray.
	verticesNative := []float32{
		-0.5, 0.5, 0, // top left
		-0.5, -0.5, 0, // bottom left
		0.5, 0.5, 0, // top right
		0.5, -0.5, 0, // bottom right
	}
	vertices := toFloat32Array(verticesNative)

	indicesNative := []uint16{
		0, 1, 3, // first triangle
		0, 3, 2, // second triangle
	}
	indices := toUint16Array(indicesNative)

	// Create buffer
	vertexBuffer := gl.Call("createBuffer", glTypes.arrayBuffer)

	// Bind to buffer
	gl.Call("bindBuffer", glTypes.arrayBuffer, vertexBuffer)

	// Pass data to buffer
	gl.Call("bufferData", glTypes.arrayBuffer, vertices, glTypes.staticDraw)

	// Unbind buffer
	gl.Call("bindBuffer", glTypes.arrayBuffer, nil)

	//// INDEX BUFFER ////

	// Create buffer
	indexBuffer := gl.Call("createBuffer", glTypes.elementArrayBuffer)

	// Bind to buffer
	gl.Call("bindBuffer", glTypes.elementArrayBuffer, indexBuffer)

	// Pass data to buffer
	gl.Call("bufferData", glTypes.elementArrayBuffer, indices, glTypes.staticDraw)

	// Unbind buffer
	gl.Call("bindBuffer", glTypes.elementArrayBuffer, nil)

	//// Shaders ////

	// Vertex shader source code
	vertCode := `
	attribute vec3 coordinates;
		
	void main(void) {
		gl_Position = vec4(coordinates, 1.0);
	}`

	// Create a vertex shader object
	vertShader := gl.Call("createShader", glTypes.vertexShader)

	// Attach vertex shader source code
	gl.Call("shaderSource", vertShader, vertCode)

	// Compile the vertex shader
	gl.Call("compileShader", vertShader)

	//fragment shader source code
	fragCode := `
	void main(void) {
		gl_FragColor = vec4(0.0, 0.0, 0.0, 0.1);
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

	// Bind vertex buffer object
	gl.Call("bindBuffer", glTypes.arrayBuffer, vertexBuffer)

	// Bind index buffer object
	gl.Call("bindBuffer", glTypes.elementArrayBuffer, indexBuffer)

	// Get the attribute location
	coord := gl.Call("getAttribLocation", shaderProgram, "coordinates")

	// Point an attribute to the currently bound VBO
	gl.Call("vertexAttribPointer", coord, 3, glTypes.float, false, 0, 0)

	// Enable the attribute
	gl.Call("enableVertexAttribArray", coord)

	//// Drawing the triangle ////

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
