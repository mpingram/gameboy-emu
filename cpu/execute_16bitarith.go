package cpu

/**
 * 16bit arithmetic
 */

// Add_HL_rr adds 16bit register rr to HL.
// Flags affected (znhc): -0hc
// Carry/Half carry is set based on upper byte (carry from bit 15, half carry from bits 11->12)
// https://stackoverflow.com/questions/57958631/game-boy-half-carry-flag-and-16-bit-instructions-especially-opcode-0xe8
func (c *CPU) Add_HL_rr(reg16 Reg16) {
	getrr, _ := c.getReg16(reg16)
	rr := getrr()
	hl := c.getHL()

	// z
	// --
	// n
	c.setFlagN(false)
	// h
	// for this instruction, half carry occurs when there is a carry from bit 11 -> bit 12.
	// In other words, perform the carry and half carry checks on the high byte only.
	r := byte(rr & 0xFF00 >> 8)
	h := byte(hl & 0xFF00 >> 8)
	c.setFlagH(addHalfCarriesByte(r, h, false))
	// c
	c.setFlagC(addOverflowsByte(r, h, false))
	c.setHL(hl + rr)
}

// Add_HL_SP adds the current SP register to HL.
// Flags affected (znhc): -0hc
// Half carry is set based on upper byte (carry from bits 11->12)
// https://stackoverflow.com/questions/57958631/game-boy-half-carry-flag-and-16-bit-instructions-especially-opcode-0xe8
func (c *CPU) Add_HL_SP() {
	hl := c.getHL()

	// z
	// --
	// n
	c.setFlagN(false)
	// h
	// for this instruction, half carry occurs when there is a carry from bit 11 -> bit 12.
	// In other words, perform the carry and half carry checks on the high byte only.
	spHighByte := byte(c.SP & 0xFF00 >> 8)
	h := byte(hl & 0xFF00 >> 8)
	halfCarry := addHalfCarriesByte(spHighByte, h, false)
	c.setFlagH(halfCarry)
	// c
	// FIXME the halfCarry here might make a lot of sense. Who knows.
	c.setFlagC(addOverflowsByte(spHighByte, h, halfCarry))
	c.setHL(hl + c.SP)
}

// Inc_rr increments 16bit register rr.
// Flags affected (znhc): ----
func (c *CPU) Inc_rr(rr Reg16) {
	getrr, setrr := c.getReg16(rr)
	setrr(getrr() + 1)
}

// Inc_SP increments the stack pointer (SP) register.
// Flags affected (znhc): ----
func (c *CPU) Inc_SP() {
	c.SP++
}

// Dec_rr decrements 16bit register rr.
// Flags affected (znhc): ----
func (c *CPU) Dec_rr(rr Reg16) {
	getrr, setrr := c.getReg16(rr)
	setrr(getrr() - 1)
}

// Dec_SP decrements the stack pointer (SP) register.
// Flags affected (znhc): ----
func (c *CPU) Dec_SP() {
	c.SP--
}

// Add_SP_r8 adds signed byte r8 to the stack pointer (SP)
// Flags affected (znhc): 00hc
// Carry/Half carry is based on lower byte (carry from bit 7, half carry from bits 3->4)
// https://stackoverflow.com/questions/57958631/game-boy-half-carry-flag-and-16-bit-instructions-especially-opcode-0xe8
func (c *CPU) Add_SP_r8(r8 int8) {
	// z
	c.setFlagZ(false)
	// n
	c.setFlagN(false)
	if r8 >= 0 {
		// h
		// for this instruction, half carry is set when a carry occurs from bit 3 to bit 4.
		// In other words, perform the carry and half carry checks on the low byte.
		c.setFlagH(addHalfCarriesByte(byte(c.SP), byte(r8), false))
		// c
		c.setFlagC(addOverflowsByte(byte(c.SP), byte(r8), false))
		c.SP = c.SP + uint16(r8)
	} else {
		// h
		c.setFlagH(subHalfCarriesByte(byte(c.SP), byte(r8), false))
		// c
		c.setFlagC(subUnderflowsByte(byte(c.SP), byte(r8), false))
		c.SP = c.SP - uint16(r8)
	}
}

// Ld_HL_SPplusr8 loads the value of memory at SP+r8 (signed byte) into HL.
// Flags affected (znhc): 00hc
// Carry/Half carry is based on lower byte (carry from bit 7, half carry from bits 3->4)
// https://stackoverflow.com/questions/57958631/game-boy-half-carry-flag-and-16-bit-instructions-especially-opcode-0xe8
// https://stackoverflow.com/questions/5159603/gbz80-how-does-ld-hl-spe-affect-h-and-c-flags
func (c *CPU) Ld_HL_SPplusr8(r8 int8) {
	// z
	c.setFlagZ(false)
	// n
	c.setFlagN(false)
	if r8 >= 0 {
		// h
		c.setFlagH(false) // DEBUG no idea what the half carry behavior would be
		// c
		c.setFlagC(false)
		c.setL(c.mem.Rb(c.SP + uint16(r8)))
		c.setH(c.mem.Rb(c.SP + uint16(r8) + 1))
	} else {
		// h
		c.setFlagH(false)
		// c
		c.setFlagC(false)
		c.setL(c.mem.Rb(c.SP - uint16(r8)))
		c.setH(c.mem.Rb(c.SP - uint16(r8) + 1))
	}

}
