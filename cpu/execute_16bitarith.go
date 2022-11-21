package cpu

import "fmt"

/**
 * 16bit arithmetic
 */

// Add_HL_rr adds 16bit register rr to HL.
// Flags affected (znhc): -0hc
// Carry/Half carry is set based on upper byte (carry from bit 15, half carry from bits 11->12)
// https://stackoverflow.com/questions/57958631/game-boy-half-carry-flag-and-16-bit-instructions-especially-opcode-0xe8
func (c *CPU) Add_HL_rr(rr Reg16) {
	getrr, _ := c.getReg16(rr)
	sum := c.getHL() + getrr()

	// z
	// --
	// n
	c.setFlagN(false)
	// h
	hlHighByte := c.getHL() & 0xFF00 >> 8
	sumHighByte := sum & 0xFF00 >> 8
	// half carry occurs when bit 11
	// c
	c.setHL(sum)
}

// Add_HL_SP adds the current SP register to HL.
// Flags affected (znhc): -0hc
// Half carry is set based on upper byte (carry from bits 11->12)
// https://stackoverflow.com/questions/57958631/game-boy-half-carry-flag-and-16-bit-instructions-especially-opcode-0xe8
func (c *CPU) Add_HL_SP() {
	fmt.Println("Instruction not implemented: ADD HL SP")
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
	fmt.Println("Instruction not implemented: ADD SP r8")
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
		c.setFlagH((c.SP&0xF)+uint16(r8&0xF) > 0xF)
		// c
		c.setFlagC((c.SP&0xFF)+uint16(r8) > 0xFF)
		c.SP = c.SP + uint16(r8)
	} else {
		// h
		c.setFlagH((c.SP & 0xF) <= (c.SP & 0xF))
		// c
		c.setFlagC((c.SP & 0xFF) <= (c.SP & 0xFF))
		c.SP = c.SP - uint16(r8*-1)
	}

}
