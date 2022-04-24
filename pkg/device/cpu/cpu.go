package cpu

import "github.com/se-nonide/go6502/pkg/device/memory"

const Frequency = 1789773

type CPU struct {
	Memory    memory.Memory // memory interface
	Cycles    uint64        // number of cycles
	PC        uint16        // program counter
	SP        byte          // stack pointer
	A         byte          // accumulator
	X         byte          // x register
	Y         byte          // y register
	C         byte          // carry flag
	Z         byte          // zero flag
	I         byte          // interrupt disable flag
	D         byte          // decimal mode flag
	B         byte          // break command flag
	U         byte          // unused flag
	V         byte          // overflow flag
	N         byte          // negative flag
	interrupt byte          // interrupt type to perform
	stall     int           // number of cycles to stall
	table     [256]func(*stepData)
}

type stepData struct {
	address uint16
	pc      uint16
	mode    byte
}
