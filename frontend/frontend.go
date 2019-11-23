package frontend

import (
	"bytes"
	"encoding/binary"
	"syscall/js"
)

// https://www.tutorialspoint.com/webgl/webgl_quick_guide.htm

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

type WebGLRenderer struct {
	gl                        js.Value
	width, height             int
	indices, vertices         js.Value
	indexBuffer, vertexBuffer js.Value
	position, texcoord        js.Value
	shaderProgram             js.Value
	texture                   js.Value
}

func NewWebGLRenderer() *WebGLRenderer {
	r := WebGLRenderer{}
	r.initialize()
	return &r
}

func (r *WebGLRenderer) initialize() {
	// Init Canvas stuff
	doc := js.Global().Get("document")
	canvasEl := doc.Call("getElementById", "gocanvas")
	r.width = doc.Get("body").Get("clientWidth").Int()
	r.height = doc.Get("body").Get("clientHeight").Int()
	canvasEl.Set("width", r.width)
	canvasEl.Set("height", r.height)

	r.gl = canvasEl.Call("getContext", "webgl")
	if r.gl == js.Undefined() {
		r.gl = canvasEl.Call("getContext", "experimental-webgl")
	}
	// once again
	if r.gl == js.Undefined() {
		js.Global().Call("alert", "browser might not support webgl")
		return
	}

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
	r.vertices = toFloat32Array(verticesNative)

	indicesNative := []uint16{
		0, 1, 3, // first triangle
		0, 3, 2, // second triangle
	}
	r.indices = toUint16Array(indicesNative)

	// Create vertex buffer
	r.vertexBuffer = r.gl.Call("createBuffer", r.gl.Get("ARRAY_BUFFER"))
	// Bind to vertex buffer
	r.gl.Call("bindBuffer", r.gl.Get("ARRAY_BUFFER"), r.vertexBuffer)
	// Pass vertices to vertex buffer
	r.gl.Call("bufferData", r.gl.Get("ARRAY_BUFFER"), r.vertices, r.gl.Get("STATIC_DRAW"))

	// Create index buffer
	r.indexBuffer = r.gl.Call("createBuffer", r.gl.Get("ELEMENT_ARRAY_BUFFER"))
	// Bind to index buffer
	r.gl.Call("bindBuffer", r.gl.Get("ELEMENT_ARRAY_BUFFER"), r.indexBuffer)
	// Pass data to index buffer
	r.gl.Call("bufferData", r.gl.Get("ELEMENT_ARRAY_BUFFER"), r.indices, r.gl.Get("STATIC_DRAW"))

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
	vertShader := r.gl.Call("createShader", r.gl.Get("VERTEX_SHADER"))
	// Attach vertex shader source code
	r.gl.Call("shaderSource", vertShader, vertCode)
	// Compile the vertex shader
	r.gl.Call("compileShader", vertShader)

	//fragment shader source code
	fragCode := `
	varying highp vec2 v_texcoord;
	uniform sampler2D u_texture;

	void main() {
		gl_FragColor = texture2D(u_texture, v_texcoord);
	}`

	// Create fragment shader object
	fragShader := r.gl.Call("createShader", r.gl.Get("FRAGMENT_SHADER"))

	// Attach fragment shader source code
	r.gl.Call("shaderSource", fragShader, fragCode)
	// Compile the fragmentt shader
	r.gl.Call("compileShader", fragShader)
	// Create a shader program object to store
	// the combined shader program
	r.shaderProgram = r.gl.Call("createProgram")

	// Attach a vertex shader
	r.gl.Call("attachShader", r.shaderProgram, vertShader)
	// Attach a fragment shader
	r.gl.Call("attachShader", r.shaderProgram, fragShader)
	// Link both the programs
	r.gl.Call("linkProgram", r.shaderProgram)
	// Use the combined shader program object
	r.gl.Call("useProgram", r.shaderProgram)

	//// Associating shaders to buffer objects ////
	// Our vertices have two attributes in them: the spatial coordinates (xyz)
	// and the texture coordinates (st). Both of these attributes are packed into
	// the same vertex buffer, so we need to tell openGL how to find these
	// two attributes using `vertexAttribPointer`

	// Bind vertex buffer object
	r.gl.Call("bindBuffer", r.gl.Get("ARRAY_BUFFER"), r.vertexBuffer)
	// Bind index buffer object
	r.gl.Call("bindBuffer", r.gl.Get("ELEMENT_ARRAY_BUFFER"), r.indexBuffer)
	// Get the `position` attribute location
	r.position = r.gl.Call("getAttribLocation", r.shaderProgram, "a_position")
	// Configure the `coordinates` attribute pointer
	r.gl.Call(
		"vertexAttribPointer",
		r.position,        // attribute location
		3,                 // number of elements in attribute
		r.gl.Get("FLOAT"), // type
		false,             // not normalized (clamp to [-1,1], not [0,1])
		5*4,               // stride: number of bytes between different vertex data elements
		0,                 // offset of first element (in bytes)
	)
	// Enable the `position` attribute
	r.gl.Call("enableVertexAttribArray", r.position)

	// Get the `texcoord` attribute
	r.texcoord = r.gl.Call("getAttribLocation", r.shaderProgram, "a_texcoord")
	// Configure the `texcoord` attribute pointer
	r.gl.Call(
		"vertexAttribPointer",
		r.texcoord,        // attribute location
		2,                 // number of elements in attribute
		r.gl.Get("FLOAT"), // type
		false,             // not normalized (clamp to [-1,1], not [0,1])
		5*4,               // stride: number of bytes between different vertex data elements
		3*4,               // offset of first element (in bytes): it is 3 floats from the start
	)
	// Enable the `texcoord` attribute
	r.gl.Call("enableVertexAttribArray", r.texcoord)

	//// Create a screen texture ////
	r.texture = r.gl.Call("createTexture")
	standbyTextureNative := make([]uint8, 0)
	for row := 0; row < 144; row++ {
		for col := 0; col < 160; col++ {
			standbyTextureNative = append(standbyTextureNative, 0, 0, 100)
		}
	}
	standbyTexture := toUint8Array(standbyTextureNative)
	r.gl.Call("bindTexture", r.gl.Get("TEXTURE_2D"), r.texture)
	// NOTE -- WebGL 1.0 requires special treatment of textures that
	// aren't powers of 2 (need to look into this more.) In any case,
	// the GB screen texture (160x144) isn't a power of 2 in either dimension,
	// so we'll need to do the special treatment thing no matter what.
	r.gl.Call(
		"texImage2D",
		r.gl.Get("TEXTURE_2D"),    // target
		0,                         // level
		r.gl.Get("RGB"),           // internal format
		160,                       // width
		144,                       // height
		0,                         // border
		r.gl.Get("RGB"),           // src format -- in WebGL 1, must be same as internal format
		r.gl.Get("UNSIGNED_BYTE"), // 1 byte per color channel (RGB), 3 bytes per pixel
		standbyTexture,            // pixel data
	)

	// turn off mipmapping, set wrapping to clamp to edge
	r.gl.Call("texParameteri", r.gl.Get("TEXTURE_2D"), r.gl.Get("TEXTURE_WRAP_S"), r.gl.Get("CLAMP_TO_EDGE"))
	r.gl.Call("texParameteri", r.gl.Get("TEXTURE_2D"), r.gl.Get("TEXTURE_WRAP_T"), r.gl.Get("CLAMP_TO_EDGE"))
	r.gl.Call("texParameteri", r.gl.Get("TEXTURE_2D"), r.gl.Get("TEXTURE_MIN_FILTER"), r.gl.Get("LINEAR"))

	//// Draw the screen ////
	// Clear the canvas
	r.gl.Call("clearColor", 0.5, 0.5, 0.5, 0.9)
	r.gl.Call("clear", r.gl.Get("COLOR_BUFFER_BIT"))
	// Enable the depth test
	r.gl.Call("enable", r.gl.Get("DEPTH_TEST"))
	// Set the view port
	r.gl.Call("viewport", 0, 0, r.width, r.height)
	// Draw the triangle
	r.gl.Call("drawElements", r.gl.Get("TRIANGLES"), r.indices.Get("length"), r.gl.Get("UNSIGNED_SHORT"), 0)
}

func (r *WebGLRenderer) Render(screen []byte) {

	// Clear the canvas
	r.gl.Call("clearColor", 0.5, 0.5, 0.5, 0.9)
	r.gl.Call("clear", r.gl.Get("COLOR_BUFFER_BIT"))
	r.gl.Call("clear", r.gl.Get("DEPTH_BUFFER_BIT"))

	r.gl.Call("activeTexture", r.gl.Get("TEXTURE0"))
	r.gl.Call(
		"bindTexture",
		r.gl.Get("TEXTURE_2D"),
		r.texture,
	)

	// Enable the `position` attribute
	r.gl.Call("enableVertexAttribArray", r.position)
	// Enable the `texcoord` attribute
	r.gl.Call("enableVertexAttribArray", r.texcoord)

	// Swap screen in as next texture
	r.gl.Call(
		"texSubImage2D",
		r.gl.Get("TEXTURE_2D"),
		0,                         // mipmap level 0
		0,                         // x offset
		0,                         // y offset
		160,                       // width (texels)
		144,                       // height (texels)
		r.gl.Get("RGB"),           // format
		r.gl.Get("UNSIGNED_BYTE"), // type,
		toUint8Array(screen),      // texture data
	)

	//// Draw the screen ////
	// Enable the depth test
	r.gl.Call("enable", r.gl.Get("DEPTH_TEST"))
	// Set the viewport
	r.gl.Call("viewport", 0, 0, r.width, r.height)
	/// use the shader program
	r.gl.Call("useProgram", r.shaderProgram)

	// Bind vertex buffer object
	r.gl.Call("bindBuffer", r.gl.Get("ARRAY_BUFFER"), r.vertexBuffer)
	// Bind index buffer object
	r.gl.Call("bindBuffer", r.gl.Get("ELEMENT_ARRAY_BUFFER"), r.indexBuffer)

	// Draw the triangle
	r.gl.Call("drawElements", r.gl.Get("TRIANGLES"), r.indices.Get("length"), r.gl.Get("UNSIGNED_SHORT"), 0)

}
