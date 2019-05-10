package cpu

import (
	"fmt"
)

const negativeFlag = 0x10
const halfCarryFlag = 0x20
const carryFlag = 0x40
const zeroFlag = 0x80

// Registers represents the Sharp LR35902's registers.
type Registers struct {
	A                uint8 // accumulator
	B, C, D, E, H, L uint8 // gen purpose
	F                uint8 // flag register

	SP uint16 // stack pointer

	PC uint16 // program counter
}

// Each 8-bit register can be combined with its complement
// to function as one 16-bit register.
// To emulate this, we use getters and setters that
// change the state of the underlying registers.
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

func (r *Registers) setAF(val uint16) {
	r.A = uint8((val & 0xFF00) >> 8)
	r.F = uint8(val & 0x00FF)
}
func (r *Registers) setBC(val uint16) {
	r.B = uint8((val & 0xFF00) >> 8)
	r.C = uint8(val & 0x00FF)
}
func (r *Registers) setDE(val uint16) {
	r.D = uint8((val & 0xFF00) >> 8)
	r.E = uint8(val & 0x00FF)
}
func (r *Registers) setHL(val uint16) {
	r.H = uint8((val & 0xFF00) >> 8)
	r.L = uint8(val & 0x00FF)
}

// Flag setters/getters for convenience.
func (r *Registers) getFlagC() byte {
	if r.F&carryFlag != 0 {
		return 1
	} else {
		return 0
	}
}

func (r *Registers) setFlagC(bit byte) {
	switch bit {
	case 0:
		r.F ^= carryFlag
	case 1:
		r.F |= carryFlag
	default:
		panic(fmt.Errorf("Attempted to set bit to value %v. Value must be 0 or 1", bit))
	}
}

func (r *Registers) getFlagH() byte {
	if r.F&halfCarryFlag != 0 {
		return 1
	} else {
		return 0
	}
}
func (r *Registers) setFlagH(bit byte) {
	switch bit {
	case 0:
		r.F ^= halfCarryFlag
	case 1:
		r.F |= halfCarryFlag
	default:
		panic(fmt.Errorf("Attempted to set bit to value %v. Value must be 0 or 1", bit))
	}
}

func (r *Registers) getFlagN() byte {
	if r.F&negativeFlag != 0 {
		return 1
	} else {
		return 0
	}
}
func (r *Registers) setFlagN(bit byte) {
	switch bit {
	case 0:
		r.F ^= negativeFlag
	case 1:
		r.F |= negativeFlag
	default:
		panic(fmt.Errorf("Attempted to set bit to value %v. Value must be 0 or 1", bit))
	}
}

func (r *Registers) getFlagZ() byte {
	if r.F&zeroFlag != 0 {
		return 1
	} else {
		return 0
	}
}
func (r *Registers) setFlagZ(bit byte) {
	switch bit {
	case 0:
		r.F ^= zeroFlag
	case 1:
		r.F |= zeroFlag
	default:
		panic(fmt.Errorf("Attempted to set bit to value %v. Value must be 0 or 1", bit))
	}
}
