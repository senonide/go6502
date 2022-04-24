package gamepad

const (
	A = iota
	B
	Select
	Start
	Up
	Down
	Left
	Right
)

type Gamepad struct {
	buttons [8]bool
	index   byte
	strobe  byte
}

func NewGamepad() *Gamepad {
	return &Gamepad{}
}

func (gp *Gamepad) SetButtons(buttons [8]bool) {
	gp.buttons = buttons
}

func (gp *Gamepad) Read() byte {
	value := byte(0)
	if gp.index < 8 && gp.buttons[gp.index] {
		value = 1
	}
	gp.index++
	if gp.strobe&1 == 1 {
		gp.index = 0
	}
	return value
}

func (gp *Gamepad) Write(value byte) {
	gp.strobe = value
	if gp.strobe&1 == 1 {
		gp.index = 0
	}
}
