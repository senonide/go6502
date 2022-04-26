package device6502

import "log"

type Memory interface {
	Read(address uint16) byte
	Write(address uint16, value byte)
}

type cpuMemory struct {
	device *Device
}

func NewCPUMemory(device *Device) Memory {
	return &cpuMemory{device}
}

func (mem *cpuMemory) Read(address uint16) byte {
	switch {
	case address < 0x2000:
		return mem.device.RAM[address%0x0800]
	case address < 0x4000:
		return mem.device.PPU.readRegister(0x2000 + address%8)
	case address == 0x4014:
		return mem.device.PPU.readRegister(address)
	case address == 0x4015:
		return mem.device.APU.readRegister(address)
	case address == 0x4016:
		return mem.device.Controller1.Read()
	case address == 0x4017:
		return mem.device.Controller2.Read()
	case address < 0x6000:
	case address >= 0x6000:
		return mem.device.Mapper.Read(address)
	default:
		log.Fatalf("unhandled cpu memory read at address: 0x%04X", address)
	}
	return 0
}

func (mem *cpuMemory) Write(address uint16, value byte) {
	switch {
	case address < 0x2000:
		mem.device.RAM[address%0x0800] = value
	case address < 0x4000:
		mem.device.PPU.writeRegister(0x2000+address%8, value)
	case address < 0x4014:
		mem.device.APU.writeRegister(address, value)
	case address == 0x4014:
		mem.device.PPU.writeRegister(address, value)
	case address == 0x4015:
		mem.device.APU.writeRegister(address, value)
	case address == 0x4016:
		mem.device.Controller1.Write(value)
		mem.device.Controller2.Write(value)
	case address == 0x4017:
		mem.device.APU.writeRegister(address, value)
	case address < 0x6000:
	case address >= 0x6000:
		mem.device.Mapper.Write(address, value)
	default:
		log.Fatalf("unhandled cpu memory write at address: 0x%04X", address)
	}
}

type ppuMemory struct {
	device *Device
}

func NewPPUMemory(device *Device) Memory {
	return &ppuMemory{device}
}

func (mem *ppuMemory) Read(address uint16) byte {
	address = address % 0x4000
	switch {
	case address < 0x2000:
		return mem.device.Mapper.Read(address)
	case address < 0x3F00:
		mode := mem.device.Cartridge.Mirror
		return mem.device.PPU.nameTableData[MirrorAddress(mode, address)%2048]
	case address < 0x4000:
		return mem.device.PPU.readPalette(address % 32)
	default:
		log.Fatalf("unhandled ppu memory read at address: 0x%04X", address)
	}
	return 0
}

func (mem *ppuMemory) Write(address uint16, value byte) {
	address = address % 0x4000
	switch {
	case address < 0x2000:
		mem.device.Mapper.Write(address, value)
	case address < 0x3F00:
		mode := mem.device.Cartridge.Mirror
		mem.device.PPU.nameTableData[MirrorAddress(mode, address)%2048] = value
	case address < 0x4000:
		mem.device.PPU.writePalette(address%32, value)
	default:
		log.Fatalf("unhandled ppu memory write at address: 0x%04X", address)
	}
}

const (
	MirrorHorizontal = 0
	MirrorVertical   = 1
	MirrorSingle0    = 2
	MirrorSingle1    = 3
	MirrorFour       = 4
)

var MirrorLookup = [...][4]uint16{
	{0, 0, 1, 1},
	{0, 1, 0, 1},
	{0, 0, 0, 0},
	{1, 1, 1, 1},
	{0, 1, 2, 3},
}

func MirrorAddress(mode byte, address uint16) uint16 {
	address = (address - 0x2000) % 0x1000
	table := address / 0x0400
	offset := address % 0x0400
	return 0x2000 + MirrorLookup[mode][table]*0x0400 + offset
}
