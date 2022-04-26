package device6502

import (
	"encoding/gob"
	"fmt"
)

const CPUFrequency = 1789773

// Interruption types
const (
	_ = iota
	interruptNone
	interruptNMI
	interruptIRQ
)

// Addressing modes
const (
	_ = iota
	modeAbsolute
	modeAbsoluteX
	modeAbsoluteY
	modeAccumulator
	modeImmediate
	modeImplied
	modeIndexedIndirect
	modeIndirect
	modeIndirectIndexed
	modeRelative
	modeZeroPage
	modeZeroPageX
	modeZeroPageY
)

// Addressing mode for each instruction
var instructionModes = [256]byte{
	6, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	1, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	6, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	6, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 8, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 13, 13, 6, 3, 6, 3, 2, 2, 3, 3,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 13, 13, 6, 3, 6, 3, 2, 2, 3, 3,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
}

// Size of each instruction in bytes
var instructionSizes = [256]byte{
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	3, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	1, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	1, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 0, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 0, 3, 0, 0,
	2, 2, 2, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
}

// Cycles used by each instruction
var instructionCycles = [256]byte{
	7, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 3, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 5, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	2, 6, 2, 6, 4, 4, 4, 4, 2, 5, 2, 5, 5, 5, 5, 5,
	2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	2, 5, 2, 5, 4, 4, 4, 4, 2, 4, 2, 4, 4, 4, 4, 4,
	2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
}

// Cycles used by each instruction when a page is crossed
var instructionPageCycles = [256]byte{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
}

// Name of each instruction
var instructionNames = [256]string{
	"BRK", "ORA", "KIL", "SLO", "NOP", "ORA", "ASL", "SLO",
	"PHP", "ORA", "ASL", "ANC", "NOP", "ORA", "ASL", "SLO",
	"BPL", "ORA", "KIL", "SLO", "NOP", "ORA", "ASL", "SLO",
	"CLC", "ORA", "NOP", "SLO", "NOP", "ORA", "ASL", "SLO",
	"JSR", "AND", "KIL", "RLA", "BIT", "AND", "ROL", "RLA",
	"PLP", "AND", "ROL", "ANC", "BIT", "AND", "ROL", "RLA",
	"BMI", "AND", "KIL", "RLA", "NOP", "AND", "ROL", "RLA",
	"SEC", "AND", "NOP", "RLA", "NOP", "AND", "ROL", "RLA",
	"RTI", "EOR", "KIL", "SRE", "NOP", "EOR", "LSR", "SRE",
	"PHA", "EOR", "LSR", "ALR", "JMP", "EOR", "LSR", "SRE",
	"BVC", "EOR", "KIL", "SRE", "NOP", "EOR", "LSR", "SRE",
	"CLI", "EOR", "NOP", "SRE", "NOP", "EOR", "LSR", "SRE",
	"RTS", "ADC", "KIL", "RRA", "NOP", "ADC", "ROR", "RRA",
	"PLA", "ADC", "ROR", "ARR", "JMP", "ADC", "ROR", "RRA",
	"BVS", "ADC", "KIL", "RRA", "NOP", "ADC", "ROR", "RRA",
	"SEI", "ADC", "NOP", "RRA", "NOP", "ADC", "ROR", "RRA",
	"NOP", "STA", "NOP", "SAX", "STY", "STA", "STX", "SAX",
	"DEY", "NOP", "TXA", "XAA", "STY", "STA", "STX", "SAX",
	"BCC", "STA", "KIL", "AHX", "STY", "STA", "STX", "SAX",
	"TYA", "STA", "TXS", "TAS", "SHY", "STA", "SHX", "AHX",
	"LDY", "LDA", "LDX", "LAX", "LDY", "LDA", "LDX", "LAX",
	"TAY", "LDA", "TAX", "LAX", "LDY", "LDA", "LDX", "LAX",
	"BCS", "LDA", "KIL", "LAX", "LDY", "LDA", "LDX", "LAX",
	"CLV", "LDA", "TSX", "LAS", "LDY", "LDA", "LDX", "LAX",
	"CPY", "CMP", "NOP", "DCP", "CPY", "CMP", "DEC", "DCP",
	"INY", "CMP", "DEX", "AXS", "CPY", "CMP", "DEC", "DCP",
	"BNE", "CMP", "KIL", "DCP", "NOP", "CMP", "DEC", "DCP",
	"CLD", "CMP", "NOP", "DCP", "NOP", "CMP", "DEC", "DCP",
	"CPX", "SBC", "NOP", "ISC", "CPX", "SBC", "INC", "ISC",
	"INX", "SBC", "NOP", "SBC", "CPX", "SBC", "INC", "ISC",
	"BEQ", "SBC", "KIL", "ISC", "NOP", "SBC", "INC", "ISC",
	"SED", "SBC", "NOP", "ISC", "NOP", "SBC", "INC", "ISC",
}

type CPU struct {
	Memory           // memory interface
	Cycles    uint64 // number of cycles
	PC        uint16 // program counter
	SP        byte   // stack pointer
	A         byte   // accumulator
	X         byte   // x register
	Y         byte   // y register
	C         byte   // carry flag
	Z         byte   // zero flag
	I         byte   // interrupt disable flag
	D         byte   // decimal mode flag
	B         byte   // break command flag
	U         byte   // unused flag
	V         byte   // overflow flag
	N         byte   // negative flag
	interrupt byte   // interrupt type to perform
	stall     int    // number of cycles to stall
	table     [256]func(*stepInfo)
}

func NewCPU(device *Device) *CPU {
	cpu := CPU{Memory: NewCPUMemory(device)}
	cpu.createTable()
	cpu.Reset()
	return &cpu
}

func (c *CPU) createTable() {
	c.table = [256]func(*stepInfo){
		c.brk, c.ora, c.kil, c.slo, c.nop, c.ora, c.asl, c.slo,
		c.php, c.ora, c.asl, c.anc, c.nop, c.ora, c.asl, c.slo,
		c.bpl, c.ora, c.kil, c.slo, c.nop, c.ora, c.asl, c.slo,
		c.clc, c.ora, c.nop, c.slo, c.nop, c.ora, c.asl, c.slo,
		c.jsr, c.and, c.kil, c.rla, c.bit, c.and, c.rol, c.rla,
		c.plp, c.and, c.rol, c.anc, c.bit, c.and, c.rol, c.rla,
		c.bmi, c.and, c.kil, c.rla, c.nop, c.and, c.rol, c.rla,
		c.sec, c.and, c.nop, c.rla, c.nop, c.and, c.rol, c.rla,
		c.rti, c.eor, c.kil, c.sre, c.nop, c.eor, c.lsr, c.sre,
		c.pha, c.eor, c.lsr, c.alr, c.jmp, c.eor, c.lsr, c.sre,
		c.bvc, c.eor, c.kil, c.sre, c.nop, c.eor, c.lsr, c.sre,
		c.cli, c.eor, c.nop, c.sre, c.nop, c.eor, c.lsr, c.sre,
		c.rts, c.adc, c.kil, c.rra, c.nop, c.adc, c.ror, c.rra,
		c.pla, c.adc, c.ror, c.arr, c.jmp, c.adc, c.ror, c.rra,
		c.bvs, c.adc, c.kil, c.rra, c.nop, c.adc, c.ror, c.rra,
		c.sei, c.adc, c.nop, c.rra, c.nop, c.adc, c.ror, c.rra,
		c.nop, c.sta, c.nop, c.sax, c.sty, c.sta, c.stx, c.sax,
		c.dey, c.nop, c.txa, c.xaa, c.sty, c.sta, c.stx, c.sax,
		c.bcc, c.sta, c.kil, c.ahx, c.sty, c.sta, c.stx, c.sax,
		c.tya, c.sta, c.txs, c.tas, c.shy, c.sta, c.shx, c.ahx,
		c.ldy, c.lda, c.ldx, c.lax, c.ldy, c.lda, c.ldx, c.lax,
		c.tay, c.lda, c.tax, c.lax, c.ldy, c.lda, c.ldx, c.lax,
		c.bcs, c.lda, c.kil, c.lax, c.ldy, c.lda, c.ldx, c.lax,
		c.clv, c.lda, c.tsx, c.las, c.ldy, c.lda, c.ldx, c.lax,
		c.cpy, c.cmp, c.nop, c.dcp, c.cpy, c.cmp, c.dec, c.dcp,
		c.iny, c.cmp, c.dex, c.axs, c.cpy, c.cmp, c.dec, c.dcp,
		c.bne, c.cmp, c.kil, c.dcp, c.nop, c.cmp, c.dec, c.dcp,
		c.cld, c.cmp, c.nop, c.dcp, c.nop, c.cmp, c.dec, c.dcp,
		c.cpx, c.sbc, c.nop, c.isc, c.cpx, c.sbc, c.inc, c.isc,
		c.inx, c.sbc, c.nop, c.sbc, c.cpx, c.sbc, c.inc, c.isc,
		c.beq, c.sbc, c.kil, c.isc, c.nop, c.sbc, c.inc, c.isc,
		c.sed, c.sbc, c.nop, c.isc, c.nop, c.sbc, c.inc, c.isc,
	}
}

func (cpu *CPU) Save(encoder *gob.Encoder) error {
	encoder.Encode(cpu.Cycles)
	encoder.Encode(cpu.PC)
	encoder.Encode(cpu.SP)
	encoder.Encode(cpu.A)
	encoder.Encode(cpu.X)
	encoder.Encode(cpu.Y)
	encoder.Encode(cpu.C)
	encoder.Encode(cpu.Z)
	encoder.Encode(cpu.I)
	encoder.Encode(cpu.D)
	encoder.Encode(cpu.B)
	encoder.Encode(cpu.U)
	encoder.Encode(cpu.V)
	encoder.Encode(cpu.N)
	encoder.Encode(cpu.interrupt)
	encoder.Encode(cpu.stall)
	return nil
}

func (cpu *CPU) Load(decoder *gob.Decoder) error {
	decoder.Decode(&cpu.Cycles)
	decoder.Decode(&cpu.PC)
	decoder.Decode(&cpu.SP)
	decoder.Decode(&cpu.A)
	decoder.Decode(&cpu.X)
	decoder.Decode(&cpu.Y)
	decoder.Decode(&cpu.C)
	decoder.Decode(&cpu.Z)
	decoder.Decode(&cpu.I)
	decoder.Decode(&cpu.D)
	decoder.Decode(&cpu.B)
	decoder.Decode(&cpu.U)
	decoder.Decode(&cpu.V)
	decoder.Decode(&cpu.N)
	decoder.Decode(&cpu.interrupt)
	decoder.Decode(&cpu.stall)
	return nil
}

func (cpu *CPU) Reset() {
	cpu.PC = cpu.Read16(0xFFFC)
	cpu.SP = 0xFD
	cpu.SetFlags(0x24)
}

func (cpu *CPU) PrintInstruction() {
	opcode := cpu.Read(cpu.PC)
	bytes := instructionSizes[opcode]
	name := instructionNames[opcode]
	w0 := fmt.Sprintf("%02X", cpu.Read(cpu.PC+0))
	w1 := fmt.Sprintf("%02X", cpu.Read(cpu.PC+1))
	w2 := fmt.Sprintf("%02X", cpu.Read(cpu.PC+2))
	if bytes < 2 {
		w1 = "  "
	}
	if bytes < 3 {
		w2 = "  "
	}
	fmt.Printf(
		"%4X  %s %s %s  %s %28s"+
			"A:%02X X:%02X Y:%02X P:%02X SP:%02X CYC:%3d\n",
		cpu.PC, w0, w1, w2, name, "",
		cpu.A, cpu.X, cpu.Y, cpu.Flags(), cpu.SP, (cpu.Cycles*3)%341)
}

func pagesDiffer(a, b uint16) bool {
	return a&0xFF00 != b&0xFF00
}

func (cpu *CPU) addBranchCycles(info *stepInfo) {
	cpu.Cycles++
	if pagesDiffer(info.pc, info.address) {
		cpu.Cycles++
	}
}

func (cpu *CPU) compare(a, b byte) {
	cpu.setZN(a - b)
	if a >= b {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
}

func (cpu *CPU) Read16(address uint16) uint16 {
	lo := uint16(cpu.Read(address))
	hi := uint16(cpu.Read(address + 1))
	return hi<<8 | lo
}

func (cpu *CPU) read16bug(address uint16) uint16 {
	a := address
	b := (a & 0xFF00) | uint16(byte(a)+1)
	lo := cpu.Read(a)
	hi := cpu.Read(b)
	return uint16(hi)<<8 | uint16(lo)
}

func (cpu *CPU) push(value byte) {
	cpu.Write(0x100|uint16(cpu.SP), value)
	cpu.SP--
}

func (cpu *CPU) pull() byte {
	cpu.SP++
	return cpu.Read(0x100 | uint16(cpu.SP))
}

func (cpu *CPU) push16(value uint16) {
	hi := byte(value >> 8)
	lo := byte(value & 0xFF)
	cpu.push(hi)
	cpu.push(lo)
}

func (cpu *CPU) pull16() uint16 {
	lo := uint16(cpu.pull())
	hi := uint16(cpu.pull())
	return hi<<8 | lo
}

func (cpu *CPU) Flags() byte {
	var flags byte
	flags |= cpu.C << 0
	flags |= cpu.Z << 1
	flags |= cpu.I << 2
	flags |= cpu.D << 3
	flags |= cpu.B << 4
	flags |= cpu.U << 5
	flags |= cpu.V << 6
	flags |= cpu.N << 7
	return flags
}

func (cpu *CPU) SetFlags(flags byte) {
	cpu.C = (flags >> 0) & 1
	cpu.Z = (flags >> 1) & 1
	cpu.I = (flags >> 2) & 1
	cpu.D = (flags >> 3) & 1
	cpu.B = (flags >> 4) & 1
	cpu.U = (flags >> 5) & 1
	cpu.V = (flags >> 6) & 1
	cpu.N = (flags >> 7) & 1
}

func (cpu *CPU) setZ(value byte) {
	if value == 0 {
		cpu.Z = 1
	} else {
		cpu.Z = 0
	}
}

func (cpu *CPU) setN(value byte) {
	if value&0x80 != 0 {
		cpu.N = 1
	} else {
		cpu.N = 0
	}
}

func (cpu *CPU) setZN(value byte) {
	cpu.setZ(value)
	cpu.setN(value)
}

func (cpu *CPU) triggerNMI() {
	cpu.interrupt = interruptNMI
}

func (cpu *CPU) triggerIRQ() {
	if cpu.I == 0 {
		cpu.interrupt = interruptIRQ
	}
}

type stepInfo struct {
	address uint16
	pc      uint16
	mode    byte
}

func (cpu *CPU) Step() int {
	if cpu.stall > 0 {
		cpu.stall--
		return 1
	}

	cycles := cpu.Cycles

	switch cpu.interrupt {
	case interruptNMI:
		cpu.nmi()
	case interruptIRQ:
		cpu.irq()
	}
	cpu.interrupt = interruptNone

	opcode := cpu.Read(cpu.PC)
	mode := instructionModes[opcode]

	var address uint16
	var pageCrossed bool
	switch mode {
	case modeAbsolute:
		address = cpu.Read16(cpu.PC + 1)
	case modeAbsoluteX:
		address = cpu.Read16(cpu.PC+1) + uint16(cpu.X)
		pageCrossed = pagesDiffer(address-uint16(cpu.X), address)
	case modeAbsoluteY:
		address = cpu.Read16(cpu.PC+1) + uint16(cpu.Y)
		pageCrossed = pagesDiffer(address-uint16(cpu.Y), address)
	case modeAccumulator:
		address = 0
	case modeImmediate:
		address = cpu.PC + 1
	case modeImplied:
		address = 0
	case modeIndexedIndirect:
		address = cpu.read16bug(uint16(cpu.Read(cpu.PC+1) + cpu.X))
	case modeIndirect:
		address = cpu.read16bug(cpu.Read16(cpu.PC + 1))
	case modeIndirectIndexed:
		address = cpu.read16bug(uint16(cpu.Read(cpu.PC+1))) + uint16(cpu.Y)
		pageCrossed = pagesDiffer(address-uint16(cpu.Y), address)
	case modeRelative:
		offset := uint16(cpu.Read(cpu.PC + 1))
		if offset < 0x80 {
			address = cpu.PC + 2 + offset
		} else {
			address = cpu.PC + 2 + offset - 0x100
		}
	case modeZeroPage:
		address = uint16(cpu.Read(cpu.PC + 1))
	case modeZeroPageX:
		address = uint16(cpu.Read(cpu.PC+1)+cpu.X) & 0xff
	case modeZeroPageY:
		address = uint16(cpu.Read(cpu.PC+1)+cpu.Y) & 0xff
	}

	cpu.PC += uint16(instructionSizes[opcode])
	cpu.Cycles += uint64(instructionCycles[opcode])
	if pageCrossed {
		cpu.Cycles += uint64(instructionPageCycles[opcode])
	}
	info := &stepInfo{address, cpu.PC, mode}
	cpu.table[opcode](info)

	return int(cpu.Cycles - cycles)
}

func (cpu *CPU) nmi() {
	cpu.push16(cpu.PC)
	cpu.php(nil)
	cpu.PC = cpu.Read16(0xFFFA)
	cpu.I = 1
	cpu.Cycles += 7
}

func (cpu *CPU) irq() {
	cpu.push16(cpu.PC)
	cpu.php(nil)
	cpu.PC = cpu.Read16(0xFFFE)
	cpu.I = 1
	cpu.Cycles += 7
}

func (cpu *CPU) adc(info *stepInfo) {
	a := cpu.A
	b := cpu.Read(info.address)
	c := cpu.C
	cpu.A = a + b + c
	cpu.setZN(cpu.A)
	if int(a)+int(b)+int(c) > 0xFF {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
	if (a^b)&0x80 == 0 && (a^cpu.A)&0x80 != 0 {
		cpu.V = 1
	} else {
		cpu.V = 0
	}
}

func (cpu *CPU) and(info *stepInfo) {
	cpu.A = cpu.A & cpu.Read(info.address)
	cpu.setZN(cpu.A)
}

func (cpu *CPU) asl(info *stepInfo) {
	if info.mode == modeAccumulator {
		cpu.C = (cpu.A >> 7) & 1
		cpu.A <<= 1
		cpu.setZN(cpu.A)
	} else {
		value := cpu.Read(info.address)
		cpu.C = (value >> 7) & 1
		value <<= 1
		cpu.Write(info.address, value)
		cpu.setZN(value)
	}
}

func (cpu *CPU) bcc(info *stepInfo) {
	if cpu.C == 0 {
		cpu.PC = info.address
		cpu.addBranchCycles(info)
	}
}

func (cpu *CPU) bcs(info *stepInfo) {
	if cpu.C != 0 {
		cpu.PC = info.address
		cpu.addBranchCycles(info)
	}
}

func (cpu *CPU) beq(info *stepInfo) {
	if cpu.Z != 0 {
		cpu.PC = info.address
		cpu.addBranchCycles(info)
	}
}

func (cpu *CPU) bit(info *stepInfo) {
	value := cpu.Read(info.address)
	cpu.V = (value >> 6) & 1
	cpu.setZ(value & cpu.A)
	cpu.setN(value)
}

func (cpu *CPU) bmi(info *stepInfo) {
	if cpu.N != 0 {
		cpu.PC = info.address
		cpu.addBranchCycles(info)
	}
}

func (cpu *CPU) bne(info *stepInfo) {
	if cpu.Z == 0 {
		cpu.PC = info.address
		cpu.addBranchCycles(info)
	}
}

func (cpu *CPU) bpl(info *stepInfo) {
	if cpu.N == 0 {
		cpu.PC = info.address
		cpu.addBranchCycles(info)
	}
}

func (cpu *CPU) brk(info *stepInfo) {
	cpu.push16(cpu.PC)
	cpu.php(info)
	cpu.sei(info)
	cpu.PC = cpu.Read16(0xFFFE)
}

func (cpu *CPU) bvc(info *stepInfo) {
	if cpu.V == 0 {
		cpu.PC = info.address
		cpu.addBranchCycles(info)
	}
}

func (cpu *CPU) bvs(info *stepInfo) {
	if cpu.V != 0 {
		cpu.PC = info.address
		cpu.addBranchCycles(info)
	}
}

func (cpu *CPU) clc(info *stepInfo) {
	cpu.C = 0
}

func (cpu *CPU) cld(info *stepInfo) {
	cpu.D = 0
}

func (cpu *CPU) cli(info *stepInfo) {
	cpu.I = 0
}

func (cpu *CPU) clv(info *stepInfo) {
	cpu.V = 0
}

func (cpu *CPU) cmp(info *stepInfo) {
	value := cpu.Read(info.address)
	cpu.compare(cpu.A, value)
}

func (cpu *CPU) cpx(info *stepInfo) {
	value := cpu.Read(info.address)
	cpu.compare(cpu.X, value)
}

func (cpu *CPU) cpy(info *stepInfo) {
	value := cpu.Read(info.address)
	cpu.compare(cpu.Y, value)
}

func (cpu *CPU) dec(info *stepInfo) {
	value := cpu.Read(info.address) - 1
	cpu.Write(info.address, value)
	cpu.setZN(value)
}

func (cpu *CPU) dex(info *stepInfo) {
	cpu.X--
	cpu.setZN(cpu.X)
}

func (cpu *CPU) dey(info *stepInfo) {
	cpu.Y--
	cpu.setZN(cpu.Y)
}

func (cpu *CPU) eor(info *stepInfo) {
	cpu.A = cpu.A ^ cpu.Read(info.address)
	cpu.setZN(cpu.A)
}

func (cpu *CPU) inc(info *stepInfo) {
	value := cpu.Read(info.address) + 1
	cpu.Write(info.address, value)
	cpu.setZN(value)
}

func (cpu *CPU) inx(info *stepInfo) {
	cpu.X++
	cpu.setZN(cpu.X)
}

func (cpu *CPU) iny(info *stepInfo) {
	cpu.Y++
	cpu.setZN(cpu.Y)
}

func (cpu *CPU) jmp(info *stepInfo) {
	cpu.PC = info.address
}

func (cpu *CPU) jsr(info *stepInfo) {
	cpu.push16(cpu.PC - 1)
	cpu.PC = info.address
}

func (cpu *CPU) lda(info *stepInfo) {
	cpu.A = cpu.Read(info.address)
	cpu.setZN(cpu.A)
}

func (cpu *CPU) ldx(info *stepInfo) {
	cpu.X = cpu.Read(info.address)
	cpu.setZN(cpu.X)
}

func (cpu *CPU) ldy(info *stepInfo) {
	cpu.Y = cpu.Read(info.address)
	cpu.setZN(cpu.Y)
}

func (cpu *CPU) lsr(info *stepInfo) {
	if info.mode == modeAccumulator {
		cpu.C = cpu.A & 1
		cpu.A >>= 1
		cpu.setZN(cpu.A)
	} else {
		value := cpu.Read(info.address)
		cpu.C = value & 1
		value >>= 1
		cpu.Write(info.address, value)
		cpu.setZN(value)
	}
}

func (cpu *CPU) nop(info *stepInfo) {
}

func (cpu *CPU) ora(info *stepInfo) {
	cpu.A = cpu.A | cpu.Read(info.address)
	cpu.setZN(cpu.A)
}

func (cpu *CPU) pha(info *stepInfo) {
	cpu.push(cpu.A)
}

func (cpu *CPU) php(info *stepInfo) {
	cpu.push(cpu.Flags() | 0x10)
}

func (cpu *CPU) pla(info *stepInfo) {
	cpu.A = cpu.pull()
	cpu.setZN(cpu.A)
}

func (cpu *CPU) plp(info *stepInfo) {
	cpu.SetFlags(cpu.pull()&0xEF | 0x20)
}

func (cpu *CPU) rol(info *stepInfo) {
	if info.mode == modeAccumulator {
		c := cpu.C
		cpu.C = (cpu.A >> 7) & 1
		cpu.A = (cpu.A << 1) | c
		cpu.setZN(cpu.A)
	} else {
		c := cpu.C
		value := cpu.Read(info.address)
		cpu.C = (value >> 7) & 1
		value = (value << 1) | c
		cpu.Write(info.address, value)
		cpu.setZN(value)
	}
}

func (cpu *CPU) ror(info *stepInfo) {
	if info.mode == modeAccumulator {
		c := cpu.C
		cpu.C = cpu.A & 1
		cpu.A = (cpu.A >> 1) | (c << 7)
		cpu.setZN(cpu.A)
	} else {
		c := cpu.C
		value := cpu.Read(info.address)
		cpu.C = value & 1
		value = (value >> 1) | (c << 7)
		cpu.Write(info.address, value)
		cpu.setZN(value)
	}
}

func (cpu *CPU) rti(info *stepInfo) {
	cpu.SetFlags(cpu.pull()&0xEF | 0x20)
	cpu.PC = cpu.pull16()
}

func (cpu *CPU) rts(info *stepInfo) {
	cpu.PC = cpu.pull16() + 1
}

func (cpu *CPU) sbc(info *stepInfo) {
	a := cpu.A
	b := cpu.Read(info.address)
	c := cpu.C
	cpu.A = a - b - (1 - c)
	cpu.setZN(cpu.A)
	if int(a)-int(b)-int(1-c) >= 0 {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
	if (a^b)&0x80 != 0 && (a^cpu.A)&0x80 != 0 {
		cpu.V = 1
	} else {
		cpu.V = 0
	}
}

func (cpu *CPU) sec(info *stepInfo) {
	cpu.C = 1
}

func (cpu *CPU) sed(info *stepInfo) {
	cpu.D = 1
}

func (cpu *CPU) sei(info *stepInfo) {
	cpu.I = 1
}

func (cpu *CPU) sta(info *stepInfo) {
	cpu.Write(info.address, cpu.A)
}

func (cpu *CPU) stx(info *stepInfo) {
	cpu.Write(info.address, cpu.X)
}

func (cpu *CPU) sty(info *stepInfo) {
	cpu.Write(info.address, cpu.Y)
}

func (cpu *CPU) tax(info *stepInfo) {
	cpu.X = cpu.A
	cpu.setZN(cpu.X)
}

func (cpu *CPU) tay(info *stepInfo) {
	cpu.Y = cpu.A
	cpu.setZN(cpu.Y)
}

func (cpu *CPU) tsx(info *stepInfo) {
	cpu.X = cpu.SP
	cpu.setZN(cpu.X)
}

func (cpu *CPU) txa(info *stepInfo) {
	cpu.A = cpu.X
	cpu.setZN(cpu.A)
}

func (cpu *CPU) txs(info *stepInfo) {
	cpu.SP = cpu.X
}

func (cpu *CPU) tya(info *stepInfo) {
	cpu.A = cpu.Y
	cpu.setZN(cpu.A)
}

func (cpu *CPU) ahx(info *stepInfo) {
}

func (cpu *CPU) alr(info *stepInfo) {
}

func (cpu *CPU) anc(info *stepInfo) {
}

func (cpu *CPU) arr(info *stepInfo) {
}

func (cpu *CPU) axs(info *stepInfo) {
}

func (cpu *CPU) dcp(info *stepInfo) {
}

func (cpu *CPU) isc(info *stepInfo) {
}

func (cpu *CPU) kil(info *stepInfo) {
}

func (cpu *CPU) las(info *stepInfo) {
}

func (cpu *CPU) lax(info *stepInfo) {
}

func (cpu *CPU) rla(info *stepInfo) {
}

func (cpu *CPU) rra(info *stepInfo) {
}

func (cpu *CPU) sax(info *stepInfo) {
}

func (cpu *CPU) shx(info *stepInfo) {
}

func (cpu *CPU) shy(info *stepInfo) {
}

func (cpu *CPU) slo(info *stepInfo) {
}

func (cpu *CPU) sre(info *stepInfo) {
}

func (cpu *CPU) tas(info *stepInfo) {
}

func (cpu *CPU) xaa(info *stepInfo) {
}
