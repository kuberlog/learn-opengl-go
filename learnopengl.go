package main

import (
	"fmt"
	"runtime"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func clearToCyberpunkColor() {
	gl.ClearColor(4.0/255.0, 217.0/255.0, 255.0/255.0, 255.0/255.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func processInput(window *glfw.Window) {
	if(window.GetKey(glfw.KeyEscape) == glfw.Press) {
		window.SetShouldClose(true)
	}
}

func init() {
	runtime.LockOSThread()
}

func initWindow() *glfw.Window {
	fmt.Println("init Window")
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	err = gl.Init()

	if err != nil {
		panic(err)
	}

	window, err := glfw.CreateWindow(800, 600, "LearnOpenGL", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	gl.Viewport(0, 0, 800, 600)

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	return window
}

func bindVerts(vertices [][]float32) []uint32 {
	vbo := make([]uint32, len(vertices))
	for i, vert := range(vertices) {
		gl.GenBuffers(1, &vbo[i])
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo[i])
		gl.BufferData(gl.ARRAY_BUFFER,
			len(vert) * 4,
			gl.Ptr(vert),
			gl.STATIC_DRAW)
	}
	return vbo
}

func triangleShaderProg() uint32 {
	// location = 1 selects the 2nd attribute pointer (index 1) in the VAO
	shaderSource := gl.Str("#version 330 core\n" +
	"layout (location = 1) in vec3 aPos;\n" +
	"void main()\n" +
	"{\n" +
	" gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);\n" +
	"}\x00")

	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vertexShader, 1, &shaderSource, nil)
	gl.CompileShader(vertexShader)

	var success int32
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &success)

	if(success == 0) {
		panic("shader couldn't compile")
	}

	fragShaderSource := gl.Str(`
	#version 330 core
	out vec4 FragColor;
	void main()
	{
		FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
	}` + "\x00")

	fragShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fragShader, 1, &fragShaderSource, nil)
	gl.CompileShader(fragShader)

	shaderProg := gl.CreateProgram()
	gl.AttachShader(shaderProg, vertexShader)
	gl.AttachShader(shaderProg, fragShader)
	gl.LinkProgram(shaderProg)
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragShader)

	return shaderProg

}

func bind_VBO_to_VAO(vbo []uint32) {
	// Vertex attribute configration, stored in the VAO
	// The first attribute pointer (reference location = {0, 1} in shader)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[0])
	gl.VertexAttribPointer(0,3, gl.FLOAT, false, 12, gl.PtrOffset(0))

	// The second attribute pointer (reference location = {0, 1} in shader)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[1])
	gl.VertexAttribPointer(1,3, gl.FLOAT, false, 12, gl.PtrOffset(0))

	// Enable the configurations, stored in the VAO
	// required to use the attribpointers in the shader progs
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
}

func main() {
	window := initWindow()

	defer glfw.Terminate()

	vertices := [][]float32{
		{-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,},

		{-0.3, -0.3, 0.0,
		0.3, -0.3, 0.0,
		0.0, 0.3, 0.0,},
	}

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	vbo := bindVerts(vertices)

	shaderProg := triangleShaderProg()

	bind_VBO_to_VAO(vbo)

	for !window.ShouldClose() {
		// Do OpenGL stuff.
		processInput(window)
		clearToCyberpunkColor()
		gl.UseProgram(shaderProg)
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
