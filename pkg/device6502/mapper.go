package device6502

import (
	"encoding/gob"
	"fmt"
	"log"
)

type Mapper interface {
	Read(address uint16) byte
	Write(address uint16, value byte)
	Step()
	Save(encoder *gob.Encoder) error
	Load(decoder *gob.Decoder) error
}

func NewMapper(device *Device) (Mapper, error) {
	cartridge := device.Cartridge
	switch cartridge.Mapper {
	case 0:
		log.Print("Mapper 0 -> 2")
		return NewMapper2(cartridge), nil
	case 1:
		log.Print("Mapper 1")
		return NewMapper1(cartridge), nil
	case 2:
		log.Print("Mapper 2")
		return NewMapper2(cartridge), nil
	case 3:
		log.Print("Mapper 3")
		return NewMapper3(cartridge), nil
	case 4:
		log.Print("Mapper 4")
		return NewMapper4(device, cartridge), nil
	case 7:
		log.Print("Mapper 7")
		return NewMapper7(cartridge), nil
	case 40:
		log.Print("Mapper 40")
		return NewMapper40(device, cartridge), nil
	case 225:
		log.Print("Mapper 225")
		return NewMapper225(cartridge), nil
	}
	err := fmt.Errorf("unsupported mapper: %d", cartridge.Mapper)
	return nil, err
}
