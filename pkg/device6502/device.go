package device6502

import (
	"encoding/gob"
	"image"
	"image/color"
	"log"
	"os"
	"path"

	"github.com/se-nonide/go6502/pkg/cartridge"
	"github.com/se-nonide/go6502/pkg/controller"
	"github.com/se-nonide/go6502/pkg/loader"
	"github.com/se-nonide/go6502/pkg/pallete"
)

type Device struct {
	CPU         *CPU
	APU         *APU
	PPU         *PPU
	Cartridge   *cartridge.Cartridge
	Controller1 *controller.Controller
	Controller2 *controller.Controller
	Mapper      Mapper
	RAM         []byte
}

func NewDevice(path string) (*Device, error) {
	cartridge, err := loader.LoadNESFile(path)
	if err != nil {
		return nil, err
	}
	ram := make([]byte, 2048)
	controller1 := controller.NewController()
	controller2 := controller.NewController()
	device := Device{
		nil, nil, nil, cartridge, controller1, controller2, nil, ram}
	mapper, err := NewMapper(&device)
	if err != nil {
		return nil, err
	}
	device.Mapper = mapper
	device.CPU = NewCPU(&device)
	device.APU = NewAPU(&device)
	device.PPU = NewPPU(&device)
	log.Printf("Nintendo Entertainment System created")
	return &device, nil
}

func (device *Device) Reset() {
	device.CPU.Reset()
}

func (device *Device) Step() int {
	//log.Print("Step")
	cpuCycles := device.CPU.Step()
	ppuCycles := cpuCycles * 3
	for i := 0; i < ppuCycles; i++ {
		device.PPU.Step()
		device.Mapper.Step()
	}
	for i := 0; i < cpuCycles; i++ {
		device.APU.Step()
	}
	return cpuCycles
}

func (device *Device) StepFrame() int {
	cpuCycles := 0
	frame := device.PPU.Frame
	for frame == device.PPU.Frame {
		cpuCycles += device.Step()
	}
	return cpuCycles
}

func (device *Device) StepSeconds(seconds float64) {
	cycles := int(CPUFrequency * seconds)
	for cycles > 0 {
		cycles -= device.Step()
	}
}

func (device *Device) Buffer() *image.RGBA {
	return device.PPU.front
}

func (device *Device) BackgroundColor() color.RGBA {
	return pallete.Palette[device.PPU.readPalette(0)%64]
}

func (device *Device) SetButtons1(buttons [8]bool) {
	device.Controller1.SetButtons(buttons)
}

func (device *Device) SetButtons2(buttons [8]bool) {
	device.Controller2.SetButtons(buttons)
}

func (device *Device) SetAudioChannel(channel chan float32) {
	device.APU.channel = channel
}

func (device *Device) SetAudioSampleRate(sampleRate float64) {
	if sampleRate != 0 {
		device.APU.sampleRate = CPUFrequency / sampleRate
		device.APU.filterChain = FilterChain{
			HighPassFilter(float32(sampleRate), 90),
			HighPassFilter(float32(sampleRate), 440),
			LowPassFilter(float32(sampleRate), 14000),
		}
	} else {
		device.APU.filterChain = nil
	}
}
func (device *Device) SaveState(filename string) error {
	dir, _ := path.Split(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	return device.Save(encoder)
}

func (device *Device) Save(encoder *gob.Encoder) error {
	encoder.Encode(device.RAM)
	device.CPU.Save(encoder)
	device.APU.Save(encoder)
	device.PPU.Save(encoder)
	device.Cartridge.Save(encoder)
	device.Mapper.Save(encoder)
	return encoder.Encode(true)
}

func (device *Device) LoadState(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	return device.Load(decoder)
}

func (device *Device) Load(decoder *gob.Decoder) error {
	decoder.Decode(&device.RAM)
	device.CPU.Load(decoder)
	device.APU.Load(decoder)
	device.PPU.Load(decoder)
	device.Cartridge.Load(decoder)
	device.Mapper.Load(decoder)
	var dummy bool
	if err := decoder.Decode(&dummy); err != nil {
		return err
	}
	return nil
}
