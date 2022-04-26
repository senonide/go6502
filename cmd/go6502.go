package main

import (
	"runtime"

	"github.com/se-nonide/go6502/internal/renderer"
)

func init() {
	//runtime.GOMAXPROCS(2)
	runtime.LockOSThread()
}

func main() {
	renderer.Start()
}
