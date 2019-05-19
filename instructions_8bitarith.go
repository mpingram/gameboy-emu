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
	if addHalfCarriesByte(c.A, d8, false) {
		c.setFlagH(true)
	} else {
		c.setFlagH(false)
	}
	// c
	if addOverflowsByte(c.A, d8, false) {
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
	if addHalfCarriesByte(c.A, d8, false) {
		c.setFlagH(true)
	} else {
		c.setFlagH(false)
	}
	// c
	if addOverflowsByte(c.A, d8, false) {
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
	if addHalfCarriesByte(c.A, b, false) {
		c.setFlagH(true)
	} else {
		c.setFlagH(false)
	}
	// c = 1 if carry else 0
	if addOverflowsByte(c.A, b, false) {
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
	get, _ := c.getReg8(r)
	valR := get()
	carry := c.getFlagC()
	sum := c.A + valR
	if carry {
		sum++
	}

	// flags
	// Z
	if sum == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// N
	c.setFlagN(false)
	// H
	c.setFlagH(addHalfCarriesByte(c.A, valR, carry))
	// C
	c.setFlagC(addOverflowsByte(c.A, valR, carry))

	c.A = sum
}

// Adc_d8 sets A = A + d8 + cy
//
// flags affected (znhc): z0hc
func (c *CPU) Adc_d8(d8 byte) {
	carry := c.getFlagC()
	sum := c.A + d8
	if carry {
		sum++
	}

	// flags
	// Z
	if sum == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// N
	c.setFlagN(false)
	// H
	c.setFlagH(addHalfCarriesByte(c.A, d8, carry))
	// C
	c.setFlagC(addOverflowsByte(c.A, d8, carry))

	c.A = sum
}

// Adc_valHL sets A = A + (HL) + cy
//
// flags affected (znhc): z0hc
func (c *CPU) Adc_valHL() {
	valHL, err := c.mem.Rb(c.getHL())
	if err != nil {
		panic(err)
	}
	carry := c.getFlagC()
	sum := c.A + valHL
	if carry {
		sum++
	}

	// flags
	// Z
	if sum == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// N
	c.setFlagN(false)
	// H
	c.setFlagH(addHalfCarriesByte(c.A, valHL, carry))
	// C
	c.setFlagC(addOverflowsByte(c.A, valHL, carry))

	c.A = sum
}

// Sub_r sets A=A-r
//
// Flags affected(zhnc): z1hc
func (c *CPU) Sub_r(r Reg8) {
	get, _ := c.getReg8(r)
	valR := get()
	sub := c.A - valR

	// flags
	// Z
	if sub == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// N
	c.setFlagN(true)
	// H
	c.setFlagH(subHalfCarriesByte(c.A, valR, false))
	// C
	c.setFlagC(subUnderflowsByte(c.A, valR, false))

	c.A = sub
}

// Sub_d8 sets A=A-d8
//
// Flags affected(zhnc): z1hc
func (c *CPU) Sub_d8(d8 byte) {
	sub := c.A - d8

	// flags
	// Z
	if sub == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// N
	c.setFlagN(true)
	// H
	c.setFlagH(subHalfCarriesByte(c.A, d8, false))
	// C
	c.setFlagC(subUnderflowsByte(c.A, d8, false))

	c.A = sub
}

// Sub_valHL sets A=A-(HL)
//
// Flags affected(zhnc): z1hc
func (c *CPU) Sub_valHL() {
	valHL, err := c.mem.Rb(c.getHL())
	if err != nil {
		panic(err)
	}
	sub := c.A - valHL

	// flags
	// Z
	if sub == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// N
	c.setFlagN(true)
	// H
	c.setFlagH(subHalfCarriesByte(c.A, valHL, false))
	// C
	c.setFlagC(subUnderflowsByte(c.A, valHL, false))

	c.A = sub
}

// Sbc_r sets A=A-r-carry
//
// Flags affected(zhnc): z1hc
func (c *CPU) Sbc_r(r Reg8) {
	carry := c.getFlagC()
	get, _ := c.getReg8(r)
	valR := get()
	sub := c.A - valR
	if carry {
		sub--
	}

	// flags
	// Z
	if sub == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// N
	c.setFlagN(true)
	// H
	c.setFlagH(subHalfCarriesByte(c.A, valR, carry))
	// C
	c.setFlagC(subUnderflowsByte(c.A, valR, carry))

	c.A = sub
}

// Sbc_d8 sets A=A-d8-carry
//
// Flags affected(zhnc): z1hc
func (c *CPU) Sbc_d8(d8 byte) {
	carry := c.getFlagC()
	sub := c.A - d8
	if carry {
		sub--
	}

	// flags
	// Z
	if sub == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// N
	c.setFlagN(true)
	// H
	c.setFlagH(subHalfCarriesByte(c.A, d8, carry))
	// C
	c.setFlagC(subUnderflowsByte(c.A, d8, carry))

	c.A = sub
}

// Sbc_valHL sets A=A-(HL)-carry
//
// Flags affected(zhnc): z1hc
func (c *CPU) Sbc_valHL() {
	carry := c.getFlagC()
	valHL, err := c.mem.Rb(c.getHL())
	if err != nil {
		panic(err)
	}
	sub := c.A - valHL
	if carry {
		sub--
	}

	// flags
	// Z
	if sub == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// N
	c.setFlagN(true)
	// H
	c.setFlagH(subHalfCarriesByte(c.A, valHL, carry))
	// C
	c.setFlagC(subUnderflowsByte(c.A, valHL, carry))

	c.A = sub
}

// And_r sets A=A&r
//
// flags affected(znhc): z010
func (c *CPU) And_r(r Reg8) {
	get, _ := c.getReg8(r)
	d8 := get()
	res := c.A & d8

	// FLAGS
	// z
	if res == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// n
	c.setFlagN(false)
	// h
	c.setFlagH(true)
	// c
	c.setFlagC(false)

	c.A = res
}

// And_d8 sets A=A&d8
//
// flags affected(znhc): z010
func (c *CPU) And_d8(d8 byte) {
	res := c.A & d8

	// FLAGS
	// z
	if res == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// n = 0
	c.setFlagN(false)
	// h = 1
	c.setFlagH(true)
	// c = 0
	c.setFlagC(false)

	c.A = res
}

// And_valHL sets A=A&(HL)
//
// flags affected(znhc): z010
func (c *CPU) And_valHL() {
	valHL, err := c.mem.Rb(c.getHL())
	if err != nil {
		panic(err)
	}
	res := c.A & valHL

	// FLAGS
	// z
	if res == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// n
	c.setFlagN(false)
	// h
	c.setFlagH(true)
	// c
	c.setFlagC(false)

	c.A = res
}

// Xor_r sets A=A^r
//
// flags affected(znhc): z000
func (c *CPU) Xor_r(r Reg8) {
	get, _ := c.getReg8(r)
	d8 := get()
	res := c.A ^ d8

	// FLAGS
	// z
	if res == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// n
	c.setFlagN(false)
	// h
	c.setFlagH(false)
	// c
	c.setFlagC(false)

	c.A = res
}

// Xor_d8 sets A=A^d8
//
// flags affected(znhc): z000
func (c *CPU) Xor_d8(d8 byte) {
	res := c.A ^ d8

	// FLAGS
	// z
	if res == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// n = 0
	c.setFlagN(false)
	// h = 0
	c.setFlagH(false)
	// c = 0
	c.setFlagC(false)

	c.A = res
}

// Xor_valHL sets A=A^(HL)
//
// Flags affected(zhnc): z000
func (c *CPU) Xor_valHL() {
	valHL, err := c.mem.Rb(c.getHL())
	if err != nil {
		panic(err)
	}
	res := c.A ^ valHL

	// FLAGS
	// z
	if res == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// n = 0
	c.setFlagN(false)
	// h = 0
	c.setFlagH(false)
	// c = 0
	c.setFlagC(false)

	c.A = res
}

// Or_r sets A=A|r
//
// flags affected(znhc): z000
func (c *CPU) Or_r(r Reg8) {
	get, _ := c.getReg8(r)
	d8 := get()
	res := c.A | d8

	// FLAGS
	// z
	if res == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// n = 0
	c.setFlagN(false)
	// h = 0
	c.setFlagH(false)
	// c = 0
	c.setFlagC(false)

	c.A = res
}

// Or_d8 sets A=A|d8
//
// flags affected(znhc): z000
func (c *CPU) Or_d8(d8 byte) {
	res := c.A | d8

	// FLAGS
	// z
	if res == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// n = 0
	c.setFlagN(false)
	// h = 0
	c.setFlagH(false)
	// c = 0
	c.setFlagC(false)

	c.A = res
}

// Or_valHL sets A=A|(HL)
//
// Flags affected(zhnc): z000
func (c *CPU) Or_valHL() {
	valHL, err := c.mem.Rb(c.getHL())
	if err != nil {
		panic(err)
	}
	res := c.A | valHL

	// FLAGS
	// z
	if res == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// n = 0
	c.setFlagN(false)
	// h = 0
	c.setFlagH(false)
	// c = 0
	c.setFlagC(false)

	c.A = res
}

// Cp_r sets flags to result of A-r
//
// Flags affected(zhnc): z1hc
func (c *CPU) Cp_r(r Reg8) {
	get, _ := c.getReg8(r)
	valR := get()
	sub := c.A - valR

	// flags
	// Z
	if sub == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// N = 1
	c.setFlagN(true)
	// H
	c.setFlagH(subHalfCarriesByte(c.A, valR, false))
	// C
	c.setFlagC(subUnderflowsByte(c.A, valR, false))
}

// Cp_d8 sets flags to result of A-d8
//
// Flags affected(zhnc): z1hc
func (c *CPU) Cp_d8(d8 byte) {
	sub := c.A - d8

	// flags
	// Z
	if sub == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// N
	c.setFlagN(true)
	// H
	c.setFlagH(subHalfCarriesByte(c.A, d8, false))
	// C
	c.setFlagC(subUnderflowsByte(c.A, d8, false))
}

// Cp_valHL sets flags to result of A-(HL)
//
// Flags affected(zhnc): z1hc
func (c *CPU) Cp_valHL() {
	valHL, err := c.mem.Rb(c.getHL())
	if err != nil {
		panic(err)
	}
	sub := c.A - valHL

	// flags
	// Z
	if sub == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// N
	c.setFlagN(true)
	// H
	c.setFlagH(subHalfCarriesByte(c.A, valHL, false))
	// C
	c.setFlagC(subUnderflowsByte(c.A, valHL, false))
}

// Inc_r increments register r
//
// flags affected(zhnc): z0h-
func (c *CPU) Inc_r(r Reg8) {
	get, set := c.getReg8(r)
	valR := get()
	res := valR + 1

	// FLAGS
	// z
	if res == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// n = 0
	c.setFlagN(false)
	// h
	c.setFlagH(addHalfCarriesByte(valR, 1, false))
	// c
	// do nothing

	set(res)
}

// Inc_valHL increments byte at address HL
//
// flags affected(zhnc): z0h-
func (c *CPU) Inc_valHL() {
	valHL, err := c.mem.Rb(c.getHL())
	if err != nil {
		panic(err)
	}
	res := valHL + 1

	// FLAGS
	// z
	if res == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// n = 0
	c.setFlagN(false)
	// h
	c.setFlagH(addHalfCarriesByte(valHL, 1, false))
	// c
	// do nothing

	err = c.mem.Wb(c.getHL(), res)
	if err != nil {
		panic(err)
	}
}

// Dec_r decrements register r
//
// flags affected(zhnc): z1h-
func (c *CPU) Dec_r(r Reg8) {
	get, set := c.getReg8(r)
	valR := get()
	res := valR - 1

	// FLAGS
	// z
	if res == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// n = 1
	c.setFlagN(true)
	// h
	c.setFlagH(subHalfCarriesByte(valR, 1, false))
	// c
	// do nothing

	set(res)
}

// Dec_valHL decrements byte at address HL
//
// flags affected(zhnc): z1h-
func (c *CPU) Dec_valHL() {
	valHL, err := c.mem.Rb(c.getHL())
	if err != nil {
		panic(err)
	}
	res := valHL - 1

	// FLAGS
	// z
	if res == 0 {
		c.setFlagZ(true)
	} else {
		c.setFlagZ(false)
	}
	// n = 1
	c.setFlagN(true)
	// h
	c.setFlagH(subHalfCarriesByte(valHL, 1, false))
	// c
	// do nothing

	err = c.mem.Wb(c.getHL(), res)
	if err != nil {
		panic(err)
	}
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

// Cpl takes the complement of A
// (sets A = A ^ 0xFF)
//
// flags affected(zhnc): -11-
func (c *CPU) Cpl() {
	c.A ^= 0xFF

	// Z
	// do nothing
	// H
	c.setFlagH(true)
	// N
	c.setFlagN(true)
	// C
	// do nothing
}
