package gamepad

import (
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/se-nonide/go6502/pkg/gamepad"
)

func readKey(window *glfw.Window, key glfw.Key) bool {
	return window.GetKey(key) == glfw.Press
}

func ReadKeys(window *glfw.Window, turbo bool) [8]bool {
	var result [8]bool
	result[gamepad.A] = readKey(window, glfw.KeyZ) || (turbo && readKey(window, glfw.KeyA))
	result[gamepad.B] = readKey(window, glfw.KeyX) || (turbo && readKey(window, glfw.KeyS))
	result[gamepad.Select] = readKey(window, glfw.KeyRightShift)
	result[gamepad.Start] = readKey(window, glfw.KeyEnter)
	result[gamepad.Up] = readKey(window, glfw.KeyUp)
	result[gamepad.Down] = readKey(window, glfw.KeyDown)
	result[gamepad.Left] = readKey(window, glfw.KeyLeft)
	result[gamepad.Right] = readKey(window, glfw.KeyRight)
	return result
}

func ReadJoystick(joy glfw.Joystick, turbo bool) [8]bool {
	var result [8]bool
	if !glfw.Joystick1.Present() {
		return result
	}
	joyname := glfw.Joystick1.GetName()
	log.Print(joyname)
	axes := glfw.Joystick1.GetAxes()
	buttons := glfw.Joystick1.GetButtons()
	if joyname == "PLAYSTATION(R)3 Controller" {
		result[gamepad.A] = buttons[14] == 1 || (turbo && buttons[2] == 1)
		result[gamepad.B] = buttons[13] == 1 || (turbo && buttons[3] == 1)
		result[gamepad.Select] = buttons[0] == 1
		result[gamepad.Start] = buttons[3] == 1
		result[gamepad.Up] = buttons[4] == 1 || axes[1] < -0.5
		result[gamepad.Down] = buttons[6] == 1 || axes[1] > 0.5
		result[gamepad.Left] = buttons[7] == 1 || axes[0] < -0.5
		result[gamepad.Right] = buttons[5] == 1 || axes[0] > 0.5
		return result
	}
	if len(buttons) < 8 {
		return result
	}
	result[gamepad.A] = buttons[0] == 1 || (turbo && buttons[2] == 1)
	result[gamepad.B] = buttons[1] == 1 || (turbo && buttons[3] == 1)
	result[gamepad.Select] = buttons[6] == 1
	result[gamepad.Start] = buttons[7] == 1
	result[gamepad.Up] = axes[1] < -0.5
	result[gamepad.Down] = axes[1] > 0.5
	result[gamepad.Left] = axes[0] < -0.5
	result[gamepad.Right] = axes[0] > 0.5
	return result
}

func JoystickReset(joy glfw.Joystick) bool {
	if !glfw.Joystick1.Present() {
		return false
	}
	buttons := glfw.Joystick1.GetButtons()
	if len(buttons) < 6 {
		return false
	}
	return buttons[4] == 1 && buttons[5] == 1
}

func CombineButtons(a, b [8]bool) [8]bool {
	var result [8]bool
	for i := 0; i < 8; i++ {
		result[i] = a[i] || b[i]
	}
	return result
}
