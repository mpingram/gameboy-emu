package cpu

import (
	"fmt"
)

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
