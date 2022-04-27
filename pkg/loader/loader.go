package loader

import (
	"encoding/binary"
	"errors"
	"io"
	"os"

	"github.com/se-nonide/go6502/pkg/cartridge"
)

const iNESFileMagic = 0x1a53454e

type iNESFileHeader struct {
	Magic    uint32  // iNES magic number
	NumPRG   byte    // number of PRG-ROM banks (16KB each)
	NumCHR   byte    // number of CHR-ROM banks (8KB each)
	Control1 byte    // control bits
	Control2 byte    // control bits
	NumRAM   byte    // PRG-RAM size (x 8KB)
	_        [7]byte // unused padding
}

func LoadNESFile(path string) (*cartridge.Cartridge, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	header := iNESFileHeader{}
	if err := binary.Read(file, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	if header.Magic != iNESFileMagic {
		return nil, errors.New("invalid .device6502 file")
	}

	mapper1 := header.Control1 >> 4
	mapper2 := header.Control2 >> 4
	mapper := mapper1 | mapper2<<4

	mirror1 := header.Control1 & 1
	mirror2 := (header.Control1 >> 3) & 1
	mirror := mirror1 | mirror2<<1

	battery := (header.Control1 >> 1) & 1

	if header.Control1&4 == 4 {
		trainer := make([]byte, 512)
		if _, err := io.ReadFull(file, trainer); err != nil {
			return nil, err
		}
	}

	prg := make([]byte, int(header.NumPRG)*16384)
	if _, err := io.ReadFull(file, prg); err != nil {
		return nil, err
	}

	chr := make([]byte, int(header.NumCHR)*8192)
	if _, err := io.ReadFull(file, chr); err != nil {
		return nil, err
	}

	if header.NumCHR == 0 {
		chr = make([]byte, 8192)
	}

	return cartridge.NewCartridge(prg, chr, mapper, mirror, battery), nil
}
