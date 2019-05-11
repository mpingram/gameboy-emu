package cpu

import (
	"fmt"
)

/**
 * I .CPU control instructions
 */

// Ccf complements(flips) the carry flag and
// resets the N and H flags.
//
// Flags(zhnc):  -00c
func (c *CPU) Ccf() {
	// xor c
	c.setFlagC(c.getFlagC() ^ 1)
	// zero n
	c.setFlagN(0)
	// zero h
	c.setFlagH(0)
}

// Scf sets the carry flag to 1 and resets
// the N and H flags.
//
// Flags(zhnc): -001
func (c *CPU) Scf() {
	// set c to 1
	c.setFlagC(1)
	// zero n
	c.setFlagN(0)
	// zero h
	c.setFlagH(0)
}

// Nop does no operation.
func (c *CPU) Nop() {
	// nothin
}

// Halt halts CPU until an interrupt occurs.
func (c *CPU) Halt() {
	c.halted = true
}

// Stop enters low-power standby mode.
func (c *CPU) Stop() {
	c.stopped = true
}

// Di disables all interrupts (sets IME bit to 0)
func (c *CPU) Di() {
	// Set IME = 0
	interrupts, err := c.mem.rb(ADDR_INTERRUPTS)
	if err != nil {
		panic(err)
	}
	interrupts ^= IME_BIT
	err = c.mem.wb(ADDR_INTERRUPTS, interrupts)
	if err != nil {
		panic(err)
	}
}

// Ei enables all interrupts (sets IME bit to 1)
func (c *CPU) Ei() {
	// Set IME = 0
	interrupts, err := c.mem.rb(ADDR_INTERRUPTS)
	if err != nil {
		panic(err)
	}
	interrupts |= IME_BIT
	err = c.mem.wb(ADDR_INTERRUPTS, interrupts)
	if err != nil {
		panic(err)
	}
}

/**
 * II. Jump commands
 */

// Jp jumps to address a16 (PC=a16)
func (c *CPU) Jp(a16 uint16) {
	c.PC = a16
}

// Jp_HL jumps to address in HL (PC=HL)
func (c *CPU) Jp_HL() {
	valHL, err := c.mem.rw(c.getHL())
	if err != nil {
		panic(err)
	}
	c.PC = valHL
}

// JpNZ conditionally jumps to a16 if Zero flag == 0.
func (c *CPU) JpNZ(a16 uint16) {
	// Jump if Z is 0
	if c.getFlagZ() == 0 {
		c.PC = a16
	}
}

// JpZ conditionally jumps to a16 if Zero flag == 1
func (c *CPU) JpZ(a16 uint16) {
	if c.getFlagZ() == 1 {
		c.PC = a16
	}
}

// JpNC conditionally jumps to a16 if Carry flag == 0
func (c *CPU) JpNC(a16 uint16) {
	if c.getFlagC() == 0 {
		c.PC = a16
	}
}

// JpC conditionally jumps to a16 if Carry flag == 1
func (c *CPU) JpC(a16 uint16) {
	if c.getFlagC() == 1 {
		c.PC = a16
	}
}

// Jr does a relative jump to PC +/- r8. (r8 is a signed byte).
func (c *CPU) Jr(r8 int8) {
	if isNegative := r8 < 0; isNegative {
		// if our signed byte r8 is negative,
		// make it positive(multiply it by -1), convert it to a uint16,
		// and subtract it from PC.
		r8 *= -1
		c.PC -= uint16(r8)
	} else {
		// our signed byte r8 is positive,
		// so we can just convert it to a uint16 and add it.
		c.PC += uint16(r8)
	}
}

// JrNZ does a Conditional relative jump to PC+/-r8 if Zero flag == 0
func (c *CPU) JrNZ(r8 int8) {
	if c.getFlagZ() == 0 {
		c.Jr(r8)
	}
}

// JrZ does a conditional relative jump to PC+/-r8 if Zero flag == 1
func (c *CPU) JrZ(r8 int8) {
	if c.getFlagZ() == 1 {
		c.Jr(r8)
	}
}

// JrNC does a conditional relative jump to PC+/-r8 if Carry flag == 0
func (c *CPU) JrNC(r8 int8) {
	if c.getFlagC() == 0 {
		c.Jr(r8)
	}
}

// JrC does a conditional relative jump to PC+/-r8 if Carry flag == 1
func (c *CPU) JrC(r8 int8) {
	if c.getFlagC() == 1 {
		c.Jr(r8)
	}
}

// Call calls a subroutine at address a16 (push PC onto stack, jump to a16)
func (c *CPU) Call(a16 uint16) {
	c.SP -= 2 // stack grows downward in memory
	err := c.mem.ww(c.SP, c.PC)
	if err != nil {
		panic(err)
	}
	c.PC = a16
}

// CallNZ conditionally call subroutine at address a16 if Zero flag == 0
func (c *CPU) CallNZ(a16 uint16) {
	// TODO in future will need to figure out a way
	// to handle the timing differences if condition
	// succeeds or fails. Return bool from instruction?
	if c.getFlagZ() == 0 {
		c.Call(a16)
	}
}

// CallZ conditionally calls a subroutine at address a16 if Zero flag == 1
func (c *CPU) CallZ(a16 uint16) {
	if c.getFlagZ() == 1 {
		c.Call(a16)
	}
}

// CallNC conditionally calls a subroutine at address a16 if Carry flag == 0
func (c *CPU) CallNC(a16 uint16) {
	if c.getFlagC() == 0 {
		c.Call(a16)
	}
}

// CallC conditionally calls a subroutine at address a16 if Carry flag == 1
func (c *CPU) CallC(a16 uint16) {
	if c.getFlagC() == 1 {
		c.Call(a16)
	}
}

// Ret returns from a subroutine. (Pop stack and jump to that address)
func (c *CPU) Ret() {
	// jump to address at top of stack
	c.PC, err = c.mem.rw(c.SP)
	if err != nil {
		panic(err)
	}
	// pop stack (stack grows downwards)
	c.SP += 2
}

// RetNZ conditionally returns from subroutine if Zero flag == 1
func (c *CPU) RetNZ() {
	if c.getFlagZ() == 1 {
		c.Ret()
	}
}

// RetZ conditionally returns from subroutine if Zero flag == 0
func (c *CPU) RetZ() {
	if c.getFlagZ() == 0 {
		c.Ret()
	}
}

// RetNC conditionally returns from subroutine if Carry flag == 1
func (c *CPU) RetNC() {
	if c.getFlagC() == 1 {
		c.Ret()
	}
}

// RetC conditionally returns from subroutine if Carry flag == 0
func (c *CPU) RetC() {
	if c.getFlagC() == 0 {
		c.Ret()
	}
}

// Reti returns and enables all interrupts (IME bit=1)
func (c *CPU) Reti() {
	c.Ei()
	c.Ret()
}

// Rst (Restart) calls to special memory addresses.
// (Push PC onto stack and jump to $0000 + n, where
// n in {0x00, 0x08, 0x10, 0x18, 0x20, 0x28, 0x30, 0x38}
func (c *CPU) Rst(n byte) {
	// Manually call to address -- c.Call() should
	// restrict access to this rea of memory.
	c.SP -= 2 // stack grows downward
	err := c.mem.ww(c.SP, c.PC)
	if err != nil {
		panic(err)
	}
	c.PC = n
}

/**
 * 8bit load commands
 */

// Ld_r1_r2 -- Load register 2 into register 1.
func (c *CPU) Ld_r1_r2(r1 Reg8, r2 Reg8) {
	// set r1 to r2
	_, setr1 := getReg8(r1)
	getr2, _ := getReg8(r2)
	setr1(getr2())
}

// Ld_r_d8 -- Load byte d8 into r
func (c *CPU) Ld_r_d8(r Reg8, d8 byte) {
	_, set := getReg8(r)
	set(d8)
}

// Ld_r_valHL loads byte at address HL into r
func (c *CPU) Ld_r_valHL(r Reg8) {
	b, err := mem.rb(c.getHL())
	if err != nil {
		panic(err)
	}
	_, set := getReg8(r)
	set(b)
}

// Ld_r_valHL loads r into byte at address HL
func (c *CPU) Ld_valHL_r(r Reg8) {
	get, _ := getReg8(r)
	b := get()
	err := c.mem.wb(c.HL, b)
	if err != nil {
		panic(err)
	}

}

// Ld_valHL_d8 loads byte d8 into memory at address HL.
func (c *CPU) Ld_valHL_d8(d8 byte) {
	err := c.mem.wb(c.HL, d8)
	if err != nil {
		panic(err)
	}
}

// Ld_A_valBC loads byte at address BC into A.
func (c *CPU) Ld_A_valBC() {
	b, err := c.mem.rb(c.BC)
	if err != nil {
		panic(err)
	}
	c.A = b
}

// Ld_A_valDE loads byte at address DE into A.
func (c *CPU) Ld_A_valDE() {
	b, err := c.mem.rb(c.BC)
	if err != nil {
		panic(err)
	}
	c.A = b
}

// Ld_A_valA16 loads byte at address a16 into A.
func (c *CPU) Ld_A_valA16(a16 uint16) {
	b, err := c.mem.rb(a16)
	if err != nil {
		panic(err)
	}
	c.A = b
}

// Ld_valBC_A loads A into byte at address BC.
func (c *CPU) Ld_valBC_A() {
	err := c.mem.wb(c.BC, c.A)
	if err != nil {
		panic(err)
	}
}

// Ld_valDE_A loads A into byte at address DE.
func (c *CPU) Ld_valDE_A() {
	err := c.mem.wb(c.DE, c.A)
	if err != nil {
		panic(err)
	}
}

// Ld_valA16_A loads A into byte at address a16.
func (c *CPU) Ld_valA16_A(a16 uint16) {
	err := c.mem.wb(a16, c.A)
	if err != nil {
		panic(err)
	}
}

// Ld_A_FF00_plus_a8 loads byte at address $FF00 + a8 into A.
func (c *CPU) Ld_A_FF00_plus_a8(a8 byte) {
	b, err := c.mem.rb(0xFF00 + a8)
	if err != nil {
		panic(err)
	}
	c.A = b
}

// Ld_FF00_plus_a8_A loads A into byte at address $FF00+a8.
func (c *CPU) Ld_FF00_plus_a8_A(a8 byte) {
	err := c.mem.wb(0xFF00+a8, c.A)
	if err != nil {
		panic(err)
	}
}

// Ld_A_FF00_plus_C loads byte at address $FF00+C into A.
func (c *CPU) Ld_A_FF00_plus_C() {
	b, err := c.mem.rb(0xFF00 + c.C)
	if err != nil {
		panic(err)
	}
	c.A = b
}

// Ld_FF00_plus_C_A loads A into byte at address $FF00+C.
func (c *CPU) Ld_FF00_plus_C_A() {
	err := c.mem.wb(0xFF00+c.C, c.A)
	if err != nil {
		panic(err)
	}
}

// Ld_valHLinc_A loads A into byte at address HL and then increments HL.
func (c *CPU) Ld_valHLinc_A() {
	err := c.mem.wb(c.HL, c.A)
	if err != nil {
		panic(err)
	}
	c.HL++
}

// Ld_A_valHLinc loads byte at address HL into A, then increments HL.
func (c *CPU) Ld_A_valHLinc() {
	b, err := c.mem.rb(c.HL)
	if err != nil {
		panic(err)
	}
	c.A = b
	c.HL++
}

// Ld_valHLdec_A loads A into byte at address HL, then decrements HL.
func (c *CPU) Ld_valHLdec_A() {
	err := c.mem.wb(c.HL, c.A)
	if err != nil {
		panic(err)
	}
	c.HL--
}

// Ld_A_valHLdec loads byte at address HL into A, then decrements HL.
func (c *CPU) Ld_A_valHLdec() {
	b, err := c.mem.rb(c.HL)
	if err != nil {
		panic(err)
	}
	c.A = b
	c.HL--
}

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
	c.SP = c.HL
}

// Push_rr pushes the value of 16-bit register onto the stack.
// (I.e., it increments SP, then loads rr into the two bytes at address of SP)
// rr can be BC, DE, HL, AF.
func (c *CPU) Push_rr(rr Reg16) {
	c.SP -= 2 // stack grows downwards
	get, _ := getReg16(rr)
	err := c.mem.ww(c.SP, get())
	if err != nil {
		panic(err)
	}
}

// Pop_rr pops a value off the stack and places it in 16bit register rr.
// (I.e., it loads word at address SP into rr, then decrements SP.)
func (c *CPU) Pop_rr(rr Reg16) {
	w, err := c.mem.rw(c.SP)
	_, set := getReg16(rr)
	set(w)
	c.SP += 2 // stack grows downwards
}

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
		panic(fmt.Errorf("Incorrect register %v passed to getReg8"))
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
		panic(fmt.Errorf("Incorrect register %v passed to getReg16"))
	}
}
