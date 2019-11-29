package cpu

import (
	"fmt"
	"time"
)

type MemoryReadWriter interface {
	MemoryReader
	MemoryWriter
}

type MemoryReader interface {
	Rb(addr uint16) byte
	Rw(addr uint16) uint16
}

type MemoryWriter interface {
	Wb(addr uint16, b byte)
	Ww(addr uint16, bb uint16)
}

// New initializes and returns an instance of CPU.
func New(memoryInterface MemoryReadWriter) *CPU {
	cpu := &CPU{mem: memoryInterface}
	// DEBUG this is an exceptionally slow clock
	cpu.Clock = time.NewTicker(time.Microsecond).C
	cpu.readyForStart = true
	return cpu
}

type CPU struct {
	Registers

	Clock <-chan time.Time
	//TClock <-chan int
	//MClock <-chan int

	mem MemoryReadWriter

	halted  bool // set by call to HALT: when halted, CPU is still `running`
	stopped bool // set by call to STOP

	readyForStart bool // meta-status set after first call to Start()

	ime    bool // Interrupt master enable
	setIME bool // set IME next instruction (used for Ei() command)

	breakpoint uint16
}

func (c *CPU) SetBreakpoint(pc uint16) {
	c.breakpoint = pc
}

func (c *CPU) Step() (Instruction, int) {
	// Decode and execute one instruction
	instr := Decode(c.PC, c.mem)

	origPC := c.PC

	c.Execute(instr)

	// Increment the PC by the instruction size, IF the instruction
	// didn't already jump the pc.
	didJump := c.PC != origPC
	cycles := instr.opc.cycles
	if !didJump {
		c.PC += instr.opc.length
		// If an instruction has a nonzero instr.opc.cyclesNoop,
		// this represents the number of cycles taken when a branching
		// instruction did not jump. Therefore, we use it instead.
		if instr.opc.cyclesIfNoop > 0 {
			cycles = instr.opc.cyclesIfNoop
		}
	}
	return instr, cycles
}

// Run begins operation of the CPU. It only has
// effect if the CPU is currently in `stopped` state.
func (c *CPU) Run() {
	if !c.readyForStart {
		fmt.Printf("cpu.Start called before it was ready!")
		return
	}

	for {
		if c.PC == c.breakpoint {
			break
		}

		// Decode and execute one instruction
		instr := Decode(c.PC, c.mem)

		// DEBUG: print instruction mnemonic
		fmt.Printf("($%04x)\t%s\n", c.PC, instr.String())
		c.Execute(instr)

		orig_pc := c.PC
		// Increment the PC by the instruction size, IF the instruction
		// didn't already jump the pc.
		didNotJump := c.PC == orig_pc
		cycles := instr.opc.cycles
		if didNotJump {
			c.PC += instr.opc.length
			// If an instruction has a nonzero instr.opc.cyclesNoop,
			// this represents the number of cycles taken when a branching
			// instruction did not jump. Therefore, we use it instead.
			if instr.opc.cyclesIfNoop > 0 {
				cycles = instr.opc.cyclesIfNoop
			}
		}

		// Wait for the correct number of clock ticks.
		ticks := cycles * 4 // Each cpu cycle == 4 clock ticks.
		for i := 0; i < ticks; i++ {
			<-c.Clock
		}
	}
}

/**
 * Memory addresses
 */

// ADDR_IE is the address of the IE (Interrupts Enabled) byte in memory,
// IE byte:
// FIXME documentation
const ADDR_IE uint16 = 0x0

// ADDR_IF is the address of the IF (Interrupt Flags) byte in memory
// IF byte:
// FIXME documentation
const ADDR_IF uint16 = 0x0

const zeroFlag = 0x80      // bit 7 of F
const negativeFlag = 0x40  // bit 6 of F
const halfCarryFlag = 0x20 // bit 5 of F
const carryFlag = 0x10     // bit 4 of

// Registers represents the Sharp LR35902's registers.
type Registers struct {
	A                uint8 // accumulator
	B, C, D, E, H, L uint8 // gen purpose

	F uint8 // flag register

	SP uint16 // stack pointer

	PC uint16 // program counter
}

// 8-bit register getters/setters
func (r *Registers) getA() byte {
	return r.A
}
func (r *Registers) getB() byte {
	return r.B
}
func (r *Registers) getC() byte {
	return r.C
}
func (r *Registers) getD() byte {
	return r.D
}
func (r *Registers) getE() byte {
	return r.E
}
func (r *Registers) getF() byte {
	return r.F
}
func (r *Registers) getH() byte {
	return r.H
}
func (r *Registers) getL() byte {
	return r.L
}

func (r *Registers) setA(b byte) {
	r.A = b
}
func (r *Registers) setB(b byte) {
	r.B = b
}
func (r *Registers) setC(b byte) {
	r.C = b
}
func (r *Registers) setD(b byte) {
	r.D = b
}
func (r *Registers) setE(b byte) {
	r.E = b
}
func (r *Registers) setF(b byte) {
	r.F = b
}
func (r *Registers) setH(b byte) {
	r.H = b
}
func (r *Registers) setL(b byte) {
	r.L = b
}

// 16bit register getters/setters
//
// NOTE The '16bit registers' are the 8bit registers
// combined to function as one 16-bit register.
func (r *Registers) getAF() uint16 {
	return uint16(r.A)<<8 | uint16(r.F)
}
func (r *Registers) getBC() uint16 {
	return uint16(r.B)<<8 | uint16(r.C)
}
func (r *Registers) getDE() uint16 {
	return uint16(r.D)<<8 | uint16(r.E)
}
func (r *Registers) getHL() uint16 {
	return uint16(r.H)<<8 | uint16(r.L)
}

func (r *Registers) setAF(bb uint16) {
	r.A = uint8((bb & 0xFF00) >> 8)
	r.F = uint8(bb & 0x00FF)
}
func (r *Registers) setBC(bb uint16) {
	r.B = uint8((bb & 0xFF00) >> 8)
	r.C = uint8(bb & 0x00FF)
}
func (r *Registers) setDE(bb uint16) {
	r.D = uint8((bb & 0xFF00) >> 8)
	r.E = uint8(bb & 0x00FF)
}
func (r *Registers) setHL(bb uint16) {
	r.H = uint8((bb & 0xFF00) >> 8)
	r.L = uint8(bb & 0x00FF)
}

type Reg8 int

const (
	RegA Reg8 = 0
	RegB Reg8 = 1
	RegC Reg8 = 2
	RegD Reg8 = 3
	RegE Reg8 = 4
	RegL Reg8 = 5
	RegH Reg8 = 6
	RegF Reg8 = 7
)

type Reg16 int

const (
	RegAF Reg16 = 7
	RegBC Reg16 = 8
	RegDE Reg16 = 9
	RegHL Reg16 = 10
)

// getReg8 returns a pair of getter/setter functions to a 8bit cpu register.
func (c *CPU) getReg8(r Reg8) (func() byte, func(byte)) {
	switch r {
	case RegA:
		return c.getA, c.setA
	case RegB:
		return c.getB, c.setB
	case RegC:
		return c.getC, c.setC
	case RegD:
		return c.getD, c.setD
	case RegE:
		return c.getE, c.setE
	case RegF:
		return c.getF, c.setF
	case RegH:
		return c.getH, c.setH
	case RegL:
		return c.getL, c.setL
	default:
		panic(fmt.Errorf("Incorrect register %v passed to getReg8", r))
	}
}

// reg16 gets a pair of getter/setter functions to a 16bit register.
func (c *CPU) getReg16(rr Reg16) (func() uint16, func(uint16)) {
	switch rr {
	case RegAF:
		return c.getAF, c.setAF
	case RegBC:
		return c.getBC, c.setBC
	case RegDE:
		return c.getDE, c.setDE
	case RegHL:
		return c.getHL, c.setHL
	default:
		panic(fmt.Errorf("Incorrect register %v passed to getReg16", rr))
	}
}

// Flags register getters/setters

// getFlagC returns true if Carry flag is
// set else false.
func (r *Registers) getFlagC() bool {
	return r.F&carryFlag != 0
}

// setFlagC sets or resets the carry flag.
func (r *Registers) setFlagC(bit bool) {
	if bit {
		// 1110 | 0001 => 1111
		// 1111 | 0001 => 1111
		r.F |= carryFlag
	} else {
		// 1010 & 1110 => 1010
		// 1011 & 1110 => 1010
		r.F = r.F &^ carryFlag
	}
}

// getFlagH returns true if the half carry
// flag is set else false.
func (r *Registers) getFlagH() bool {
	return r.F&halfCarryFlag != 0
}

// setFlagH sets or resets the half carry flag.
func (r *Registers) setFlagH(bit bool) {
	if bit {
		r.F |= halfCarryFlag
	} else {
		r.F = r.F &^ halfCarryFlag
	}
}

// getFlagN returns true if the negative
// operation flag is set else false.
func (r *Registers) getFlagN() bool {
	return r.F&negativeFlag != 0
}

// setFlagN sets or resets the negative operation flag.
func (r *Registers) setFlagN(bit bool) {
	if bit {
		r.F |= negativeFlag
	} else {
		r.F = r.F &^ negativeFlag
	}
}

// getFlagZ returns true if the zero flag
// is set else false.
func (r *Registers) getFlagZ() bool {
	return r.F&zeroFlag != 0
}

// setFlagZ sets or resets the zero flag.
func (r *Registers) setFlagZ(bit bool) {
	if bit {
		r.F |= zeroFlag
	} else {
		r.F = r.F &^ zeroFlag
	}
}
