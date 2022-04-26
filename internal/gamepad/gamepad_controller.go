package gamepad

import (
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/se-nonide/go6502/pkg/device6502"
)

func readKey(window *glfw.Window, key glfw.Key) bool {
	return window.GetKey(key) == glfw.Press
}

func ReadKeys(window *glfw.Window, turbo bool) [8]bool {
	var result [8]bool
	result[device6502.ButtonA] = readKey(window, glfw.KeyZ) || (turbo && readKey(window, glfw.KeyA))
	result[device6502.ButtonB] = readKey(window, glfw.KeyX) || (turbo && readKey(window, glfw.KeyS))
	result[device6502.ButtonSelect] = readKey(window, glfw.KeyRightShift)
	result[device6502.ButtonStart] = readKey(window, glfw.KeyEnter)
	result[device6502.ButtonUp] = readKey(window, glfw.KeyUp)
	result[device6502.ButtonDown] = readKey(window, glfw.KeyDown)
	result[device6502.ButtonLeft] = readKey(window, glfw.KeyLeft)
	result[device6502.ButtonRight] = readKey(window, glfw.KeyRight)
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
		result[device6502.ButtonA] = buttons[14] == 1 || (turbo && buttons[2] == 1)
		result[device6502.ButtonB] = buttons[13] == 1 || (turbo && buttons[3] == 1)
		result[device6502.ButtonSelect] = buttons[0] == 1
		result[device6502.ButtonStart] = buttons[3] == 1
		result[device6502.ButtonUp] = buttons[4] == 1 || axes[1] < -0.5
		result[device6502.ButtonDown] = buttons[6] == 1 || axes[1] > 0.5
		result[device6502.ButtonLeft] = buttons[7] == 1 || axes[0] < -0.5
		result[device6502.ButtonRight] = buttons[5] == 1 || axes[0] > 0.5
		return result
	}
	if len(buttons) < 8 {
		return result
	}
	result[device6502.ButtonA] = buttons[0] == 1 || (turbo && buttons[2] == 1)
	result[device6502.ButtonB] = buttons[1] == 1 || (turbo && buttons[3] == 1)
	result[device6502.ButtonSelect] = buttons[6] == 1
	result[device6502.ButtonStart] = buttons[7] == 1
	result[device6502.ButtonUp] = axes[1] < -0.5
	result[device6502.ButtonDown] = axes[1] > 0.5
	result[device6502.ButtonLeft] = axes[0] < -0.5
	result[device6502.ButtonRight] = axes[0] > 0.5
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
