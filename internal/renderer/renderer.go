package renderer

import (
	"log"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/se-nonide/go6502/internal/gamepad"
	"github.com/se-nonide/go6502/internal/graphics"
	"github.com/se-nonide/go6502/pkg/device6502"
)

const width = 256
const height = 240
const scale = 4
const FPS = 240

type Renderer struct {
	window  *glfw.Window
	nes     *device6502.Device
	texture uint32
}

func NewRenderer(window *glfw.Window, path string) Renderer {
	nes, err := device6502.NewDevice(path)
	if err != nil {
		log.Fatal(err)
	}
	nes.Reset()
	texture := graphics.CreateTexture()
	return Renderer{window: window, nes: nes, texture: texture}
}

func Start(path string) {
	err := glfw.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()
	window, err := glfw.CreateWindow(width*scale, height*scale, "Go 6502", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.SetSizeLimits(width*scale, height*scale, width*scale, height*scale)
	window.MakeContextCurrent()
	err = gl.Init()
	if err != nil {
		log.Fatal(err)
	}
	gl.Enable(gl.TEXTURE_2D)
	gl.ClearColor(0, 0, 0, 1)
	renderer := NewRenderer(window, path)
	renderer.Run()
}

func (r Renderer) Run() {
	var deltaTime float64 = 0
	var timestamp float64 = glfw.GetTime()
	r.window.SetKeyCallback(r.onKey)
	for !r.window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		current := glfw.GetTime()
		deltaTime += (current - timestamp)
		timestamp = current
		r.nes.StepSeconds(deltaTime)
		if deltaTime >= (1.0 / FPS) {
			deltaTime = 0
			r.Render()
		}

		r.window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (r Renderer) Render() {
	updateControllers(r.window, r.nes)
	gl.BindTexture(gl.TEXTURE_2D, r.texture)
	graphics.SetTexture(r.nes.Buffer())
	r.drawBuffer(r.window)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (r Renderer) onKey(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		switch key {
		case glfw.KeyR:
			log.Print("Reset")
			r.nes.Reset()
		}
	}
}

func (r Renderer) drawBuffer(window *glfw.Window) {
	w, h := window.GetFramebufferSize()
	s1 := float32(w) / 256
	s2 := float32(h) / 240
	f := float32(1)
	var x, y float32
	if s1 >= s2 {
		x = f * s2 / s1
		y = f
	} else {
		x = f
		y = f * s1 / s2
	}
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(-x, -y)
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(x, -y)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(x, y)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(-x, y)
	gl.End()
}

func updateControllers(window *glfw.Window, nes *device6502.Device) {
	turbo := nes.PPU.Frame%6 < 3
	k1 := gamepad.ReadKeys(window, turbo)
	j1 := gamepad.ReadJoystick(glfw.Joystick1, turbo)
	j2 := gamepad.ReadJoystick(glfw.Joystick2, turbo)
	nes.SetButtons1(gamepad.CombineButtons(k1, j1))
	nes.SetButtons2(j2)
}
