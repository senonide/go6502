package main

import (
	"log"
	"os"
	"runtime"

	"github.com/se-nonide/go6502/internal/renderer"
)

func init() {
	//runtime.GOMAXPROCS(2)
	runtime.LockOSThread()
}

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatal("Specify the path for a game to play")
	}
	renderer.Start(args[1])
}
