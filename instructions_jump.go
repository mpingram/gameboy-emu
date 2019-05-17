package cpu

import (
	"fmt"
)

/**
 * II. Jump commands
 */

// Jp jumps to address a16 (PC=a16)
func (c *CPU) Jp(a16 uint16) {
	c.PC = a16
}

// Jp_HL jumps to address in HL (PC=HL)
func (c *CPU) Jp_HL() {
	c.PC = c.getHL()
}

// JpNZ conditionally jumps to a16 if Zero flag == 0.
func (c *CPU) JpNZ(a16 uint16) {
	// Jump if Z is 0
	if c.getFlagZ() == false {
		c.PC = a16
	}
}

// JpZ conditionally jumps to a16 if Zero flag == 1
func (c *CPU) JpZ(a16 uint16) {
	if c.getFlagZ() {
		c.PC = a16
	}
}

// JpNC conditionally jumps to a16 if Carry flag == 0
func (c *CPU) JpNC(a16 uint16) {
	if c.getFlagC() == false {
		c.PC = a16
	}
}

// JpC conditionally jumps to a16 if Carry flag == 1
func (c *CPU) JpC(a16 uint16) {
	if c.getFlagC() {
		c.PC = a16
	}
}

// Jr does a relative jump to PC +/- r8. (r8 is a signed byte).
func (c *CPU) Jr(r8 int8) {
	if isNegative := r8 < 0; isNegative {
		// if our signed byte r8 is negative,
		// make it positive(multiply it by -1), convert it to a uint16,
		// and subtract it from PC.
		c.PC -= uint16(r8 * -1)
	} else {
		// our signed byte r8 is positive,
		// so we can just convert it to a uint16 and add it.
		c.PC += uint16(r8)
	}
}

// JrNZ does a Conditional relative jump to PC+/-r8 if Zero flag == 0
func (c *CPU) JrNZ(r8 int8) {
	if c.getFlagZ() == false {
		c.Jr(r8)
	}
}

// JrZ does a conditional relative jump to PC+/-r8 if Zero flag == 1
func (c *CPU) JrZ(r8 int8) {
	if c.getFlagZ() == true {
		c.Jr(r8)
	}
}

// JrNC does a conditional relative jump to PC+/-r8 if Carry flag == 0
func (c *CPU) JrNC(r8 int8) {
	if c.getFlagC() == false {
		c.Jr(r8)
	}
}

// JrC does a conditional relative jump to PC+/-r8 if Carry flag == 1
func (c *CPU) JrC(r8 int8) {
	if c.getFlagC() == true {
		c.Jr(r8)
	}
}

// Call calls a subroutine at address a16 (push PC onto stack, jump to a16)
func (c *CPU) Call(a16 uint16) {
	c.SP -= 2 // stack grows downward in memory
	err := c.mem.Ww(c.SP, c.PC)
	if err != nil {
		panic(err)
	}
	c.PC = a16
}

// CallNZ conditionally call subroutine at address a16 if Zero flag == 0
func (c *CPU) CallNZ(a16 uint16) {
	if c.getFlagZ() == false {
		c.Call(a16)
	}
}

// CallZ conditionally calls a subroutine at address a16 if Zero flag == 1
func (c *CPU) CallZ(a16 uint16) {
	if c.getFlagZ() == true {
		c.Call(a16)
	}
}

// CallNC conditionally calls a subroutine at address a16 if Carry flag == 0
func (c *CPU) CallNC(a16 uint16) {
	if c.getFlagC() == false {
		c.Call(a16)
	}
}

// CallC conditionally calls a subroutine at address a16 if Carry flag == 1
func (c *CPU) CallC(a16 uint16) {
	if c.getFlagC() == true {
		c.Call(a16)
	}
}

// Ret returns from a subroutine. (Pop stack and jump to that address)
func (c *CPU) Ret() {
	// jump to address at top of stack
	b, err := c.mem.Rw(c.SP)
	if err != nil {
		panic(err)
	}
	c.PC = b
	// pop stack (stack grows downwards)
	c.SP += 2
}

// RetNZ conditionally returns from subroutine if Zero flag == 0
func (c *CPU) RetNZ() {
	if c.getFlagZ() == false {
		c.Ret()
	}
}

// RetZ conditionally returns from subroutine if Zero flag == 1
func (c *CPU) RetZ() {
	if c.getFlagZ() == true {
		c.Ret()
	}
}

// RetNC conditionally returns from subroutine if Carry flag == 0
func (c *CPU) RetNC() {
	if c.getFlagC() == false {
		c.Ret()
	}
}

// RetC conditionally returns from subroutine if Carry flag == 1
func (c *CPU) RetC() {
	if c.getFlagC() == true {
		c.Ret()
	}
}

// Reti returns and (immediately) enables all interrupts (IME bit=1)
func (c *CPU) Reti() {
	c.Ret()
	c.ime = true
}

// Rst (Restart) calls to special memory addresses.
// (Push PC onto stack and jump to $0000 + n, where
// n in {0x00, 0x08, 0x10, 0x18, 0x20, 0x28, 0x30, 0x38}
func (c *CPU) Rst(n RstTarget) {
	switch n {
	case 0x00, 0x08, 0x10, 0x18, 0x20, 0x28, 0x30, 0x38:
		// Manually call to address -- c.Call() might
		// restrict access to this rea of memory.
		c.SP -= 2 // stack grows downward
		err := c.mem.Ww(c.SP, c.PC)
		if err != nil {
			panic(err)
		}
		c.PC = uint16(n)

	default:
		// panic if n is not 0x00, 0x08, 0x10, 0x18, 0x20, 0x28, 0x30, 0x38
		panic(fmt.Sprintf("Invalid Rst target %02x", n))
	}

}
