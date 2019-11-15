package cpu

/**
 * I .CPU control instructions
 */

// Ccf complements(flips) the carry flag and
// resets the N and H flags.
//
// Flags(zhnc):  -00c
func (c *CPU) Ccf() {
	// flip c
	c.setFlagC(!c.getFlagC())
	// zero n
	c.setFlagN(false)
	// zero h
	c.setFlagH(false)
}

// Scf sets the carry flag to 1 and resets
// the N and H flags.
//
// Flags(zhnc): -001
func (c *CPU) Scf() {
	// set c to 1
	c.setFlagC(true)
	// zero n
	c.setFlagN(false)
	// zero h
	c.setFlagH(false)
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
	c.ime = false
}

// Ei enables all interrupts (sets IME bit to 1)
// Ei only takes effect after one machine cycle.
func (c *CPU) Ei() {
	c.setIME = true
}
