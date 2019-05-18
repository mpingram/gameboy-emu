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

	// flags
	// Z
	if c.A == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// N
	c.setFlagN(false)

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

func (c *CPU) Inc_valHL() {
	notImpl()
}

func (c *CPU) Dec_r(r Reg8) {
	notImpl()
}

func (c *CPU) Dec_valHL() {
	notImpl()
}

func (c *CPU) Daa() {
	/**
	* From https://www.reddit.com/r/EmuDev/comments/4ycoix/a_guide_to_the_gameboys_halfcarry_flag/
	*
	* --------------
	* u = 0;
	* if (FH || (!FN && (RA & 0xf) > 9)) { // If half carry, or if last op was addition and lower nyb of A is not oob (ie, less than 9)
	*   u = 6;
	* }
	* if (FC || (!FN && RA > 0x99)) {
	*   u |= 0x60;
	*   FC = 1;
	* }
	* RA += FN ? -u : u;
	* FZ_EQ0(RA);
	* FH = 0;
	* -----------------
	*
	* FZ, FN, FH, FC = flags
	* RA = A
	 */

	// If last op had half carry,
	// or if last op was addition and lower nyb of A needs to be adjusted (ie, is greater than 9)
	var u byte
	if c.getFlagH() || (!c.getFlagN() && (c.A&0x0F) > 0x09) {
		u = 0x06
	}
	// If last op had carry,
	// or if upper nyb of A needs to be adjusted (ie, is greater than 99)
	if c.getFlagC() || (!c.getFlagN() && c.A > 0x99) {
		u |= 0x60
		c.setFlagC(true)
	}
	// Adjust A by subtracting u if last operation was a subtraction,
	// else adjust A by adding u
	if c.getFlagN() {
		c.A -= u
	} else {
		c.A += u
	}

	// Set zero flag if A is 0
	if c.A == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// Set half carry to 0
	c.setFlagH(false)
}

func (c *CPU) Cpl() {
	notImpl()
}
