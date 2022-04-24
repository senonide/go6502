package main

import (
	"log"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
)

const width = 256
const height = 240

func init() {

	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(width*4, height*4, "Go 6502", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	window.MakeContextCurrent()

	for !window.ShouldClose() {
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
