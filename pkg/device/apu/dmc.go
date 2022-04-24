package apu

import "github.com/se-nonide/go6502/pkg/device/cpu"

type DMC struct {
	cpu            *cpu.CPU
	enabled        bool
	value          byte
	sampleAddress  uint16
	sampleLength   uint16
	currentAddress uint16
	currentLength  uint16
	shiftRegister  byte
	bitCount       byte
	tickPeriod     byte
	tickValue      byte
	loop           bool
	irq            bool
}
