package device

import (
	"image"

	"github.com/se-nonide/go6502/pkg/cartridge"
	"github.com/se-nonide/go6502/pkg/device/apu"
	"github.com/se-nonide/go6502/pkg/device/cpu"
	"github.com/se-nonide/go6502/pkg/device/mappers"
	"github.com/se-nonide/go6502/pkg/device/ppu"
	"github.com/se-nonide/go6502/pkg/gamepad"
)

type NintendoEntertainmentSystem struct {
	CPU       *cpu.CPU             // CentralProcessingUnit
	RAM       []byte               // RandomAccessMemory
	PPU       *ppu.PPU             // PictureProcessingUnit
	APU       *apu.APU             // AudioProccesingUnit
	Cartridge *cartridge.Cartridge // Cartridge
	Gamepad1  *gamepad.Gamepad     // Gamepad1
	Gamepad2  *gamepad.Gamepad     // Gamepad2
	Mapper    mappers.Mapper       // Mapper
}

func NewNintendoEntertainmentSystem() (*NintendoEntertainmentSystem, error) {
	/*cartridge, err := LoadNESFile(path)
	if err != nil {
		return nil, err
	}*/
	ram := make([]byte, 2048)
	/*controller1 := NewController()
	controller2 := NewController()*/
	nes := NintendoEntertainmentSystem{
		nil, ram, nil, nil, nil, nil, nil, nil}
	/*mapper, err := NewMapper(&console)
	if err != nil {
		return nil, err
	}
	console.Mapper = mapper
	console.CPU = NewCPU(&console)
	console.APU = NewAPU(&console)*/
	nes.PPU = ppu.NewPPU()
	return &nes, nil
}

func (nes *NintendoEntertainmentSystem) Reset() {

}

func (nes *NintendoEntertainmentSystem) Buffer() *image.RGBA {
	return nes.PPU.Front
}
