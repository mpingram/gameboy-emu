package cpu

/**
 * VI. Rotates/ shifts
 */

// Rlc_A rotates A one bit to the left with bit 7 being moved to bit 0 and also
// stored into the carry.
func (c *CPU) Rlc_A() {
	notImpl()
}

// Rl_A rotates A one bit to the left with the carry's value put into bit 0
// and bit 7 put into the carry.
func (c *CPU) Rl_A() {
	notImpl()
}

// Rrc_A rotates A one bit to the right with bit 0 being moved to bit 7 and also
// stored into the carry.
func (c *CPU) Rrc_A() {
	notImpl()
}

// Rr_A rotates A one bit to the right with the carry's value put into bit 7 and
// bit 0 being moved to the carry.
func (c *CPU) Rr_A() {
	notImpl()
}

// Rlc_r rotates register r one bit to the left with bit 7 being moved to bit 0 and also
// stored into the carry.
func (c *CPU) Rlc_r(r Reg8) {
	notImpl()
}

// Rlc_valHL rotates the byte at $HL one bit to the left with bit 7 being moved to bit 0
// and also stored into the carry.
func (c *CPU) Rlc_valHL() {
	notImpl()
}

// Rl_r rotates register r one bit to the left with the carry's value put into bit 0
// and bit 7 put into the carry.
func (c *CPU) Rl_r(r Reg8) {
	notImpl()
}

// Rl_valHL rotates the byte at $HL one bit to the left with the carry's value put into bit 0
// and bit 7 put into the carry.
func (c *CPU) Rl_valHL() {
	notImpl()
}

// Rrc_r rotates register r one bit to the right with bit 0 being moved to bit 7 and also
// stored into the carry.
func (c *CPU) Rrc_r(r Reg8) {
	notImpl()
}

// Rrc_valHL rotates the byte at $HL one bit to the right with bit 0 being moved to bit 7 and also
// stored into the carry.
func (c *CPU) Rrc_valHL() {
	notImpl()
}

// Rr_r rotates register r one bit to the right with the carry's value put into bit 7 and
// bit 0 being moved to the carry.
func (c *CPU) Rr_r(r Reg8) {
	notImpl()
}

// Rr_valHL rotates the byte at $HL one bit to the right with the carry's value put into bit 7 and
// bit 0 being moved to the carry.
func (c *CPU) Rr_valHL() {
	notImpl()
}

// Sla_r shifts register r to the left with bit 7 moved to the carry flag and bit 0 reset (zeroed).
func (c *CPU) Sla_r(r Reg8) {
	notImpl()
}

// Sla_valHL shifts the byte at $HL to the left with bit 7 moved to the carry flag and bit 0 reset (zeroed).
func (c *CPU) Sla_valHL() {
	notImpl()
}

// Swap_r swaps the low and high nibble of register r
func (c *CPU) Swap_r(r Reg8) {
	notImpl()
}

// Swap_valHL swaps the low and high nibble of the byte at $HL
func (c *CPU) Swap_valHL() {
	notImpl()
}

// Sra_r shifts register r to the right with bit 0 moved to the carry flag
// and bit 7 retaining its original value.
func (c *CPU) Sra_r(r Reg8) {
	notImpl()
}

// Sra_valHL shifts the byte at $HL to the right with bit 0 moved to the carry flag
// and bit 7 retaining its original value.
func (c *CPU) Sra_valHL() {
	notImpl()
}

// Srl_r shifts register r to the right with bit 0 moved to the carry flag and bit 7 zeroed.
func (c *CPU) Srl_r(r Reg8) {
	notImpl()
}

// Srl_valHL shifts the byte at $HL to the right with bit 0 moved to the carry flag and bit 7 zeroed.
func (c *CPU) Srl_valHL() {
	notImpl()
}

/**
 * Bit operations
 */

// Bit_n_r sets Z=1 if bit n of register r is 0.
// Flags set (znhc): z01-
func (c *CPU) Bit_n_r(n uint8, r Reg8) {
	getr, _ := c.getReg8(r)
	b := getr()
	mask := byte(0x1) << n
	if b^mask == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	c.setFlagN(false)
	c.setFlagH(true)
}

// Bit_n_valHL sets z=1 if bit n of the byte at $HL is 0.
func (c *CPU) Bit_n_valHL(n uint8) {
	hl := c.getHL()
	b := c.mem.Rb(hl)
	mask := byte(0x1) << n
	c.mem.Wb(hl, b&mask)
}

// Set_n_r sets bit n of register r.
func (c *CPU) Set_n_r(n uint8, r Reg8) {
	getr, setr := c.getReg8(r)
	b := getr()
	mask := byte(0x1) << n
	setr(b | mask)
}

// Set_n_valHL sets bit n of the byte at $HL.
func (c *CPU) Set_n_valHL(n uint8) {
	// FIXME investigate -- does n wrap if n>7?
	hl := c.getHL()
	b := c.mem.Rb(hl)
	mask := byte(0x1) << n
	c.mem.Wb(hl, b|mask)
}

// Res_n_r unsets bit n of register r.
func (c *CPU) Res_n_r(n uint8, r Reg8) {
	getr, setr := c.getReg8(r)
	b := getr()
	mask := ^(byte(0x1) << n)
	setr(b & mask)
}

// Res_n_valHL unsets bit n of the byte at $HL.
func (c *CPU) Res_n_valHL(n uint8) {
	hl := c.getHL()
	b := c.mem.Rb(hl)
	mask := ^(byte(0x1) << n)
	c.mem.Wb(hl, b&mask)
}
