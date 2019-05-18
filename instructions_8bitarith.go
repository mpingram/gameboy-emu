package cpu

/**
 * 8bit arithmetic
 */

// Add_r adds value of 8bit register r to A. (A=A+r)
//
// Flags affected (znhc): z0hc
func (c *CPU) Add_r(r Reg8) {
	get, _ := c.getReg8(r)
	d8 := get()
	sum := c.A + d8

	// FLAGS
	// z
	if sum == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// n
	c.setFlagN(false)
	// h
	if addHalfCarriesByte(c.A, d8, 0) {
		c.setFlagH(true)
	} else {
		c.setFlagH(false)
	}
	// c
	if addOverflowsByte(c.A, d8, 0) {
		c.setFlagC(true)
	} else {
		c.setFlagC(false)
	}

	c.A = sum
}

// Add_d8 adds d8 to A (A=A+d8)
//
// flags affected (znhc): z0hc
func (c *CPU) Add_d8(d8 byte) {
	sum := c.A + d8

	// FLAGS
	// z
	if sum == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// n
	c.setFlagN(false)
	// h
	if addHalfCarriesByte(c.A, d8, 0) {
		c.setFlagH(true)
	} else {
		c.setFlagH(false)
	}
	// c
	if addOverflowsByte(c.A, d8, 0) {
		c.setFlagC(true)
	} else {
		c.setFlagC(false)
	}

	c.A = sum
}

// Add_valHL adds the byte at address HL to A.
//
// flags affected (znhc): z0hc
func (c *CPU) Add_valHL() {
	b, err := c.mem.Rb(c.getHL())
	if err != nil {
		panic(err)
	}
	sum := c.A + b

	// FLAGS
	// z = 1 if sum is 0 else 0
	if sum == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// n = 0
	c.setFlagN(false)
	// h = 1 if half carry else 0
	if addHalfCarriesByte(c.A, b, 0) {
		c.setFlagH(true)
	} else {
		c.setFlagH(false)
	}
	// c = 1 if carry else 0
	if addOverflowsByte(c.A, b, 0) {
		c.setFlagC(true)
	} else {
		c.setFlagC(false)
	}

	c.A = sum
}

// Adc_r sets A = A + r + cy
//
// flags affected (znhc): z0hc
func (c *CPU) Adc_r(r Reg8) {
	notImpl()
}

func (c *CPU) Adc_d8(d8 byte) {
	notImpl()
}

func (c *CPU) Adc_valHL() {
	notImpl()
}

func (c *CPU) Sub_r(r Reg8) {
	notImpl()
}

func (c *CPU) Sub_d8(d8 byte) {
	notImpl()
}

func (c *CPU) Sub_valHL() {
	notImpl()
}

func (c *CPU) Sbc_r(r Reg8) {
	notImpl()
}

func (c *CPU) Sbc_d8(d8 byte) {
	notImpl()
}

func (c *CPU) Sbc_valHL() {
	notImpl()
}

func (c *CPU) And_r(r Reg8) {
	notImpl()
}

func (c *CPU) And_d8(d8 byte) {
	notImpl()
}

func (c *CPU) And_valHL() {
	notImpl()
}

func (c *CPU) Xor_r(r Reg8) {
	notImpl()
}

func (c *CPU) Xor_d8(d8 byte) {
	notImpl()
}

func (c *CPU) Xor_valHL() {
	notImpl()
}

func (c *CPU) Or_r(r Reg8) {
	notImpl()
}

func (c *CPU) Or_d8(d8 byte) {
	notImpl()
}

func (c *CPU) Or_valHL() {
	notImpl()
}

func (c *CPU) Cp_r(r Reg8) {
	notImpl()
}

func (c *CPU) Cp_d8(d8 byte) {
	notImpl()
}

func (c *CPU) Cp_valHL() {
	notImpl()
}

func (c *CPU) Inc_r(r Reg8) {
	notImpl()
}

func (c *CPU) Inc_d8(d8 byte) {
	notImpl()
}

func (c *CPU) Inc_valHL() {
	notImpl()
}

func (c *CPU) Dec_r(r Reg8) {
	notImpl()
}

func (c *CPU) Dec_d8(d8 byte) {
	notImpl()
}

func (c *CPU) Dec_valHL() {
	notImpl()
}

func (c *CPU) Daa() {
	notImpl()
}

func (c *CPU) Cpl() {
	notImpl()
}
