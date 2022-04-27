package gamepad

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/se-nonide/go6502/pkg/controller"
)

func readKey(window *glfw.Window, key glfw.Key) bool {
	return window.GetKey(key) == glfw.Press
}

func ReadKeys(window *glfw.Window, turbo bool) [8]bool {
	var result [8]bool
	result[controller.ButtonA] = readKey(window, glfw.KeyZ) || (turbo && readKey(window, glfw.KeyA))
	result[controller.ButtonB] = readKey(window, glfw.KeyX) || (turbo && readKey(window, glfw.KeyS))
	result[controller.ButtonSelect] = readKey(window, glfw.KeyRightShift)
	result[controller.ButtonStart] = readKey(window, glfw.KeyEnter)
	result[controller.ButtonUp] = readKey(window, glfw.KeyUp)
	result[controller.ButtonDown] = readKey(window, glfw.KeyDown)
	result[controller.ButtonLeft] = readKey(window, glfw.KeyLeft)
	result[controller.ButtonRight] = readKey(window, glfw.KeyRight)
	return result
}

func ReadJoystick(joy glfw.Joystick, turbo bool) [8]bool {
	var result [8]bool
	if !glfw.Joystick1.Present() {
		return result
	}
	//axes := glfw.Joystick1.GetAxes()
	buttons := glfw.Joystick1.GetButtons()
	if len(buttons) < 8 {
		return result
	}
	result[controller.ButtonA] = buttons[1] == 1 || (turbo && buttons[2] == 1)
	result[controller.ButtonB] = buttons[0] == 1 || (turbo && buttons[3] == 1)
	result[controller.ButtonSelect] = buttons[6] == 1
	result[controller.ButtonStart] = buttons[7] == 1
	result[controller.ButtonUp] = buttons[11] == 1
	result[controller.ButtonDown] = buttons[13] == 1
	result[controller.ButtonLeft] = buttons[14] == 1
	result[controller.ButtonRight] = buttons[12] == 1
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
