package cpu

import (
	"fmt"
)

/**
 * IV. 16bit load commands
 */

// Ld_rr_d16 loads word d16 into 16bit register rr,
// where rr is BC, DE, or HL.
func (c *CPU) Ld_rr_d16(rr Reg16, d16 uint16) {
	switch rr {
	case RegBC:
		c.setBC(d16)
	case RegDE:
		c.setDE(d16)
	case RegHL:
		c.setHL(d16)
	default:
		panic(fmt.Errorf("Ld_rr_d16 called with invalid register %v", rr))
	}
}

// Ld_SP_d16 loads word d16 into the SP(stack pointer) register.
func (c *CPU) Ld_SP_d16(d16 uint16) {
	c.SP = d16
}

// Ld_SP_HL loads HL into the SP(stack pointer) register.
func (c *CPU) Ld_SP_HL() {
	c.SP = c.getHL()
}

// Push_rr pushes the value of 16-bit register onto the stack.
// (I.e., it increments SP, then loads rr into the two bytes at address of SP)
// rr can be BC, DE, HL, AF.
func (c *CPU) Push_rr(rr Reg16) {
	c.SP -= 2 // stack grows downwards
	get, _ := c.getReg16(rr)
	err := c.mem.Ww(c.SP, get())
	if err != nil {
		panic(err)
	}
}

// Pop_rr pops a value off the stack and places it in 16bit register rr.
// (I.e., it loads word at address SP into rr, then decrements SP.)
func (c *CPU) Pop_rr(rr Reg16) {
	w, err := c.mem.Rw(c.SP)
	if err != nil {
		panic(err)
	}
	_, set := c.getReg16(rr)
	set(w)
	c.SP += 2 // stack grows downwards
}
