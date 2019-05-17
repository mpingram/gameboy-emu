package cpu

import (
	"fmt"
)

type RstTarget byte

const (
	RstTarget0x00 RstTarget = 0x00
	RstTarget0x08 RstTarget = 0x08
	RstTarget0x10 RstTarget = 0x10
	RstTarget0x18 RstTarget = 0x18
	RstTarget0x20 RstTarget = 0x20
	RstTarget0x28 RstTarget = 0x28
	RstTarget0x30 RstTarget = 0x30
	RstTarget0x38 RstTarget = 0x38
)

func addHalfCarriesByte(a, b, carry byte) bool {
	// add the lower nibbles of a,b,c, and check
	// if the sum carries over to the higher nibble.
	return ((a&0x0f)+(b&0x0f)+(carry))&0x10 == 0x10
}

func addOverflowsByte(a, b, carry byte) bool {
	if uint16(a)+uint16(b)+uint16(carry) != uint16(a+b+carry) {
		return true
	}
	return false
}

func subHalfCarriesByte(a, b, carry byte) bool {
	// if the lower nibble of a is less than
	// the lower nibble of b, there will be a
	// carry from bit 4, the first bit of the upper nibble.
	if a&0x0f < (b&0x0f + carry) {
		return true
	}
	return false
}

func subUnderflowsByte(a, b, carry byte) bool {
	if uint16(a) < (uint16(b) - uint16(carry)) {
		return true
	}
	return false
}

func notImpl() {
	fmt.Printf("Instruction not implemented.")
}
