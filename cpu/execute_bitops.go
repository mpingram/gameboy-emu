package cpu

/**
 * VI. Rotates/ shifts
 */

// Rlc_A rotates A one bit to the left with bit 7 being moved to bit 0 and also
// stored into the carry.
// Flags affected (znhc): 000c
func (c *CPU) Rlc_A() {
	b := c.A
	mask := byte(0b1000_0000)
	// if carry (if bit 7 of b is 1)
	if b&mask != 0 {
		// shift bits left (discard bit 7)
		b = b << 1
		// set bit new 0 to 1 (wrap bit 7 around)
		b |= 0b1
		c.A = b
		// set the carry flag
		c.setFlagC(true)
	} else {
		// No carry.
		// Shift bits left.
		b = b << 1
		c.A = b
		c.setFlagC(false)
	}
	// Zero all other flags
	c.setFlagZ(false)
	c.setFlagN(false)
	c.setFlagH(false)
}

// Rl_A rotates A one bit to the left with the carry's value put into bit 0
// and bit 7 put into the carry.
// Flags affected (znhc): 000c
func (c *CPU) Rl_A() {
	b := c.A
	mask := byte(0b1000_0000)
	oldCarry := c.getFlagC()
	// if carry will occur (if bit 7 of b is 1)
	if b&mask != 0 {
		// set the carry flag
		c.setFlagC(true)
	} else {
		c.setFlagC(false)
	}
	// shift b left one byte and
	// set bit 0 to value of old carry
	b = b << 1
	if oldCarry {
		b |= 1
	}
	c.A = b
	// Zero all other flags
	c.setFlagZ(false)
	c.setFlagN(false)
	c.setFlagH(false)
}

// Rrc_A rotates A one bit to the right with bit 0 being moved to bit 7 and also
// stored into the carry.
// Flags affected (znhc): 000c
func (c *CPU) Rrc_A() {
	b := c.A
	mask := byte(0b0000_0001)
	// if carry (if bit 0 of b is 1)
	if b&mask != 0 {
		// shift bits right (discard bit 0)
		b = b >> 1
		// set bit new 7 to 1 (wrap bit 0 around)
		b |= 0b1000_0000
		c.A = b
		// set the carry flag
		c.setFlagC(true)
	} else {
		// No carry.
		// Shift bits right.
		b = b >> 1
		c.A = b
		c.setFlagC(false)
	}
	// Zero all other flags
	c.setFlagZ(false)
	c.setFlagN(false)
	c.setFlagH(false)
}

// Rr_A rotates A one bit to the right with the carry's value put into bit 7 and
// bit 0 being moved to the carry.
// Flags affected (znhc): 000c
func (c *CPU) Rr_A() {
	b := c.A
	mask := byte(0b0000_0001)
	oldCarry := c.getFlagC()
	// if carry will occur (if bit 0 of b is 1)
	if b&mask != 0 {
		// set the carry flag
		c.setFlagC(true)
	} else {
		c.setFlagC(false)
	}
	// shift b right one byte and
	// set bit 7 to value of old carry
	b = b >> 1
	if oldCarry {
		b |= 0b1000_0000
	}
	c.A = b
	// Zero all other flags
	c.setFlagZ(false)
	c.setFlagN(false)
	c.setFlagH(false)
}

// Rlc_r rotates register r one bit to the left with bit 7 being moved to bit 0 and also
// stored into the carry.
// Flags affected (znhc): z00c
func (c *CPU) Rlc_r(r Reg8) {
	getr, setr := c.getReg8(r)
	b := getr()
	mask := byte(0b1000_0000)
	// if carry (if bit 7 of b is 1)
	if b&mask != 0 {
		// shift bits left (discard bit 7)
		b = b << 1
		// set bit new 0 to 1 (wrap bit 7 around)
		b |= 0b1
		setr(b)
		// set the carry flag
		c.setFlagC(true)
		// set zero flag
		c.setFlagZ(b == 0)
	} else {
		// No carry.
		// Shift bits left.
		b = b << 1
		setr(b)
		c.setFlagC(false)
		c.setFlagZ(b == 0)
	}
	// Zero all other flags
	c.setFlagN(false)
	c.setFlagH(false)
}

// Rlc_valHL rotates the byte at $HL one bit to the left with bit 7 being moved to bit 0
// and also stored into the carry.
// Flags affected (znhc): z00c
func (c *CPU) Rlc_valHL() {
	b := c.mem.Rb(c.getHL())
	mask := byte(0b1000_0000)
	// if carry (if bit 7 of b is 1)
	if b&mask != 0 {
		// shift bits left (discard bit 7)
		b = b << 1
		// set bit new 0 to 1 (wrap bit 7 around)
		b |= 0b1
		c.mem.Wb(c.getHL(), b)
		// set the carry flag
		c.setFlagC(true)
		c.setFlagZ(b == 0)
	} else {
		// No carry.
		// Shift bits left.
		b = b << 1
		c.mem.Wb(c.getHL(), b)
		c.setFlagC(false)
		c.setFlagZ(b == 0)
	}
	// Zero all other flags
	c.setFlagN(false)
	c.setFlagH(false)
}

// Rl_r rotates register r one bit to the left with the carry's value put into bit 0
// and bit 7 put into the carry.
// Flags affected (znhc): z00c
func (c *CPU) Rl_r(r Reg8) {
	getr, setr := c.getReg8(r)
	b := getr()
	mask := byte(0b1000_0000)
	oldCarry := c.getFlagC()
	// if carry (if bit 7 of b is 1)
	if b&mask != 0 {
		// set the carry flag
		c.setFlagC(true)
	} else {
		c.setFlagC(false)
	}
	// shift b left one byte and
	// set bit 0 to value of old carry
	b = b << 1
	if oldCarry {
		b |= 1
	}
	setr(b)
	c.setFlagZ(b == 0)
	// Zero all other flags
	c.setFlagN(false)
	c.setFlagH(false)
}

// Rl_valHL rotates the byte at $HL one bit to the left with the carry's value put into bit 0
// and bit 7 put into the carry.
// Flags affected (znhc): z00c
func (c *CPU) Rl_valHL() {
	hl := c.getHL()
	b := c.mem.Rb(hl)
	mask := byte(0b1000_0000)
	oldCarry := c.getFlagC()
	// if carry (if bit 7 of b is 1)
	if b&mask != 0 {
		// set the carry flag
		c.setFlagC(true)
	} else {
		c.setFlagC(false)
	}
	// shift b left one byte and
	// set bit 0 to value of old carry
	b = b << 1
	if oldCarry {
		b |= 1
	}
	c.mem.Wb(hl, b)
	c.setFlagZ(b == 0)
	// Zero all other flags
	c.setFlagN(false)
	c.setFlagH(false)
}

// Rrc_r rotates register r one bit to the right with bit 0 being moved to bit 7 and also
// stored into the carry.
// Flags affected (znhc): z00c
func (c *CPU) Rrc_r(r Reg8) {
	getr, setr := c.getReg8(r)
	b := getr()
	mask := byte(0b0000_0001)
	// if carry (if bit 0 of b is 1)
	if b&mask != 0 {
		// shift bits right (discard bit 0)
		b = b >> 1
		// set new bit 7 to 1 (wrap bit 0 around)
		b |= 0b1000_0000
		setr(b)
		// set the carry flag
		c.setFlagC(true)
		// set zero flag
		c.setFlagZ(b == 0)
	} else {
		// No carry.
		// Shift bits right.
		b = b >> 1
		setr(b)
		c.setFlagC(false)
		c.setFlagZ(b == 0)
	}
	// Zero all other flags
	c.setFlagN(false)
	c.setFlagH(false)
}

// Rrc_valHL rotates the byte at $HL one bit to the right with bit 0 being moved to bit 7 and also
// stored into the carry.
// Flags affected (znhc): z00c
func (c *CPU) Rrc_valHL() {
	b := c.mem.Rb(c.getHL())
	mask := byte(0b0000_0001)
	// if carry (if bit 0 of b is 1)
	if b&mask != 0 {
		// shift bits right (discard bit 0)
		b = b >> 1
		// set new bit 7 to 1 (wrap bit 0 around)
		b |= 0b1000_0000
		c.mem.Wb(c.getHL(), b)
		// set the carry flag
		c.setFlagC(true)
		c.setFlagZ(b == 0)
	} else {
		// No carry.
		// Shift bits right.
		b = b >> 1
		c.mem.Wb(c.getHL(), b)
		c.setFlagC(false)
		c.setFlagZ(b == 0)
	}
	// Zero all other flags
	c.setFlagN(false)
	c.setFlagH(false)
}

// Rr_r rotates register r one bit to the right with the carry's value put into bit 7 and
// bit 0 being moved to the carry.
// Flags affected (znhc): z00c
func (c *CPU) Rr_r(r Reg8) {
	getr, setr := c.getReg8(r)
	b := getr()
	mask := byte(0b0000_0001)
	oldCarry := c.getFlagC()
	// if carry (if bit 0 of b is 1)
	if b&mask != 0 {
		// set the carry flag
		c.setFlagC(true)
	} else {
		c.setFlagC(false)
	}
	// shift b right one byte and
	// set bit 7 to value of old carry
	b = b >> 1
	if oldCarry {
		b |= 0b1000_0000
	}
	setr(b)
	c.setFlagZ(b == 0)
	// Zero all other flags
	c.setFlagN(false)
	c.setFlagH(false)
}

// Rr_valHL rotates the byte at $HL one bit to the right with the carry's value put into bit 7 and
// bit 0 being moved to the carry.
// Flags affected (znhc): z00c
func (c *CPU) Rr_valHL() {
	hl := c.getHL()
	b := c.mem.Rb(hl)
	mask := byte(0b0000_0001)
	oldCarry := c.getFlagC()
	// if carry (if bit 0 of b is 1)
	if b&mask != 0 {
		// set the carry flag
		c.setFlagC(true)
	} else {
		c.setFlagC(false)
	}
	// shift b right one bit and
	// set bit 7 to value of old carry
	b = b >> 1
	if oldCarry {
		b |= 0b1000_0000
	}
	c.mem.Wb(hl, b)
	c.setFlagZ(b == 0)
	// Zero all other flags
	c.setFlagN(false)
	c.setFlagH(false)
}

// Sla_r shifts register r to the left with bit 7 moved to the carry flag and bit 0 reset (zeroed).
// Flags affected (znhc): z00c
func (c *CPU) Sla_r(r Reg8) {
	getr, setr := c.getReg8(r)
	b := getr()
	mask := byte(0b1000_0000)
	// if carry (if bit 7 of b is 1)
	if b&mask != 0 {
		// set the carry flag
		c.setFlagC(true)
	} else {
		c.setFlagC(false)
	}
	// shift b left
	b = b << 1
	setr(b)
	c.setFlagZ(b == 0)
	// Zero all other flags
	c.setFlagN(false)
	c.setFlagH(false)
}

// Sla_valHL shifts the byte at $HL to the left with bit 7 moved to the carry flag and bit 0 reset (zeroed).
// Flags affected (znhc): z00c
func (c *CPU) Sla_valHL() {
	hl := c.getHL()
	b := c.mem.Rb(hl)
	mask := byte(0b1000_0000)
	// if carry (if bit 7 of b is 1)
	if b&mask != 0 {
		// set the carry flag
		c.setFlagC(true)
	} else {
		c.setFlagC(false)
	}
	// shift b left
	b = b << 1
	c.mem.Wb(hl, b)
	c.setFlagZ(b == 0)
	// Zero all other flags
	c.setFlagN(false)
	c.setFlagH(false)
}

// Swap_r swaps the low and high nibble of register r
// Flags affected (znhc): z000
func (c *CPU) Swap_r(r Reg8) {
	getr, setr := c.getReg8(r)
	b := getr()
	lo := b & 0x0F
	hi := b & 0xF0
	b = (lo << 4) | (hi >> 4)
	setr(b)
	c.setFlagZ(b == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(false)
}

// Swap_valHL swaps the low and high nibble of the byte at $HL
// Flags affected (znhc): z000
func (c *CPU) Swap_valHL() {
	hl := c.getHL()
	b := c.mem.Rb(hl)
	lo := b & 0x0F
	hi := b & 0xF0
	b = (lo << 4) | (hi >> 4)
	c.mem.Wb(hl, b)
	c.setFlagZ(b == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(false)
}

// Sra_r shifts register r to the right with bit 0 moved to the carry flag
// and bit 7 retaining its original value.
// Flags affected (znhc): z00c
func (c *CPU) Sra_r(r Reg8) {
	getr, setr := c.getReg8(r)
	b := getr()
	mask := byte(0b0000_0001)
	// if carry (if bit 0 of b is 1)
	if b&mask != 0 {
		// set the carry flag
		c.setFlagC(true)
	} else {
		c.setFlagC(false)
	}
	// shift b right
	b = b >> 1
	setr(b)
	c.setFlagZ(b == 0)
	// Zero all other flags
	c.setFlagN(false)
	c.setFlagH(false)
}

// Sra_valHL shifts the byte at $HL to the right with bit 0 moved to the carry flag
// and bit 7 retaining its original value.
// Flags affected (znhc): z00c
func (c *CPU) Sra_valHL() {
	hl := c.getHL()
	b := c.mem.Rb(hl)
	mask := byte(0b0000_0001)
	// if carry (if bit 0 of b is 1)
	if b&mask != 0 {
		// set the carry flag
		c.setFlagC(true)
	} else {
		c.setFlagC(false)
	}
	// shift b right
	b = b >> 1
	c.mem.Wb(hl, b)
	c.setFlagZ(b == 0)
	// Zero all other flags
	c.setFlagN(false)
	c.setFlagH(false)
}

// Srl_r shifts register r to the right with bit 0 moved to the carry flag and bit 7 zeroed.
// Flags affected (znhc): z00c
func (c *CPU) Srl_r(r Reg8) {
	getr, setr := c.getReg8(r)
	b := getr()
	mask := byte(0b0000_0001)
	oldBit7 := b&0b1000_0000 > 0
	// if carry (if bit 0 of b is 1)
	if b&mask != 0 {
		// set the carry flag
		c.setFlagC(true)
	} else {
		c.setFlagC(false)
	}
	// shift b right
	b = b >> 1
	// set bit 7 to old bit 7
	if oldBit7 {
		b |= 0b1000_0000
	}
	setr(b)
	c.setFlagZ(b == 0)
	// Zero all other flags
	c.setFlagN(false)
	c.setFlagH(false)
}

// Srl_valHL shifts the byte at $HL to the right with bit 0 moved to the carry flag and bit 7 zeroed.
// Flags affected (znhc): z00c
func (c *CPU) Srl_valHL() {
	hl := c.getHL()
	b := c.mem.Rb(hl)
	mask := byte(0b0000_0001)
	oldBit7 := b&0b1000_0000 > 0
	// if carry (if bit 0 of b is 1)
	if b&mask != 0 {
		// set the carry flag
		c.setFlagC(true)
	} else {
		c.setFlagC(false)
	}
	// shift b right
	b = b >> 1
	// set bit 7 to old bit 7
	if oldBit7 {
		b |= 0b1000_0000
	}
	c.mem.Wb(hl, b)
	c.setFlagZ(b == 0)
	// Zero all other flags
	c.setFlagN(false)
	c.setFlagH(false)
}

/**
 * Bit operations
 */

// Bit_n_r sets Z=1 if bit n of register r is 0.
// Flags set (znhc): z01-
func (c *CPU) Bit_n_r(n uint8, r Reg8) {
	getr, _ := c.getReg8(r)
	b := getr()
	mask := byte(1) << n
	c.setFlagZ(b&mask == 0)

	c.setFlagN(false)
	c.setFlagH(true)
}

// Bit_n_valHL sets z=1 if bit n of the byte at $HL is 0.
// Flags set (znhc): z01-
func (c *CPU) Bit_n_valHL(n uint8) {
	hl := c.getHL()
	b := c.mem.Rb(hl)
	mask := byte(1) << n
	c.setFlagZ(b&mask == 0)

	c.setFlagN(false)
	c.setFlagH(true)
}

// Set_n_r sets bit n of register r.
// Flags set (znhc): ----
func (c *CPU) Set_n_r(n byte, r Reg8) {
	getr, setr := c.getReg8(r)
	b := getr()
	mask := byte(1) << n
	setr(b | mask)
}

// Set_n_valHL sets bit n of the byte at $HL.
// Flags set (znhc): ----
func (c *CPU) Set_n_valHL(n byte) {
	// FIXME investigate -- does n wrap if n>7?
	hl := c.getHL()
	b := c.mem.Rb(hl)
	mask := byte(1) << n
	c.mem.Wb(hl, b|mask)
}

// Res_n_r unsets bit n of register r.
// Flags set (znhc): ----
func (c *CPU) Res_n_r(n byte, r Reg8) {
	getr, setr := c.getReg8(r)
	b := getr()
	mask := ^(byte(1) << n)
	setr(b & mask)
}

// Res_n_valHL unsets bit n of the byte at $HL.
// Flags set (znhc): ----
func (c *CPU) Res_n_valHL(n byte) {
	hl := c.getHL()
	b := c.mem.Rb(hl)
	mask := ^(byte(1) << n)
	c.mem.Wb(hl, b&mask)
}
