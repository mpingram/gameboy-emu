package cpu

/**
 * 8bit load commands
 */

// Ld_r1_r2 -- Load register 2 into register 1.
func (c *CPU) Ld_r1_r2(r1 Reg8, r2 Reg8) {
	// set r1 to r2
	_, setr1 := c.getReg8(r1)
	getr2, _ := c.getReg8(r2)
	setr1(getr2())
}

// Ld_r_d8 -- Load byte d8 into r
func (c *CPU) Ld_r_d8(r Reg8, d8 byte) {
	_, set := c.getReg8(r)
	set(d8)
}

// Ld_r_valHL loads byte at address HL into r
func (c *CPU) Ld_r_valHL(r Reg8) {
	b := c.mem.Rb(c.getHL())
	_, set := c.getReg8(r)
	set(b)
}

// Ld_r_valHL loads r into byte at address HL
func (c *CPU) Ld_valHL_r(r Reg8) {
	get, _ := c.getReg8(r)
	b := get()
	c.mem.Wb(c.getHL(), b)

}

// Ld_valHL_d8 loads byte d8 into memory at address HL.
func (c *CPU) Ld_valHL_d8(d8 byte) {
	c.mem.Wb(c.getHL(), d8)
}

// Ld_A_valBC loads byte at address BC into A.
func (c *CPU) Ld_A_valBC() {
	b := c.mem.Rb(c.getBC())
	c.A = b
}

// Ld_A_valDE loads byte at address DE into A.
func (c *CPU) Ld_A_valDE() {
	b := c.mem.Rb(c.getDE())
	c.A = b
}

// Ld_A_valA16 loads byte at address a16 into A.
func (c *CPU) Ld_A_valA16(a16 uint16) {
	b := c.mem.Rb(a16)
	c.A = b
}

// Ld_valBC_A loads A into byte at address BC.
func (c *CPU) Ld_valBC_A() {
	c.mem.Wb(c.getBC(), c.A)
}

// Ld_valDE_A loads A into byte at address DE.
func (c *CPU) Ld_valDE_A() {
	c.mem.Wb(c.getDE(), c.A)
}

// Ld_valA16_A loads A into byte at address a16.
func (c *CPU) Ld_valA16_A(a16 uint16) {
	c.mem.Wb(a16, c.A)
}

// Ld_A_FF00_plus_a8 loads byte at address $FF00 + a8 into A.
func (c *CPU) Ld_A_FF00_plus_a8(a8 byte) {
	b := c.mem.Rb(0xFF00 + uint16(a8))
	c.A = b
}

// Ld_FF00_plus_a8_A loads A into byte at address $FF00+a8.
func (c *CPU) Ld_FF00_plus_a8_A(a8 byte) {
	c.mem.Wb(0xFF00+uint16(a8), c.A)
}

// Ld_A_FF00_plus_C loads byte at address $FF00+C into A.
func (c *CPU) Ld_A_FF00_plus_C() {
	b := c.mem.Rb(0xFF00 + uint16(c.C))
	c.A = b
}

// Ld_FF00_plus_C_A loads A into byte at address $FF00+C.
func (c *CPU) Ld_FF00_plus_C_A() {
	c.mem.Wb(0xFF00+uint16(c.C), c.A)
}

// Ld_valHLinc_A loads A into byte at address HL and then increments HL.
func (c *CPU) Ld_valHLinc_A() {
	c.mem.Wb(c.getHL(), c.A)
	c.setHL(c.getHL() + 1)
}

// Ld_A_valHLinc loads byte at address HL into A, then increments HL.
func (c *CPU) Ld_A_valHLinc() {
	b := c.mem.Rb(c.getHL())
	c.A = b
	c.setHL(c.getHL() + 1)
}

// Ld_valHLdec_A loads A into byte at address HL, then decrements HL.
func (c *CPU) Ld_valHLdec_A() {
	c.mem.Wb(c.getHL(), c.A)
	c.setHL(c.getHL() - 1)
}

// Ld_A_valHLdec loads byte at address HL into A, then decrements HL.
func (c *CPU) Ld_A_valHLdec() {
	b := c.mem.Rb(c.getHL())
	c.A = b
	c.setHL(c.getHL() - 1)
}
