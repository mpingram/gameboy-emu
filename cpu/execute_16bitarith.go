package cpu

import "fmt"

/**
 * 16bit arithmetic
 */

// Add_HL_rr adds 16bit register r to HL.
// Flags affected (znhc): -0hc
func (c *CPU) Add_HL_rr(rr Reg16) {
	getrr, _ := c.getReg16(rr)
	hl := c.getHL() + getrr()
	c.setHL(hl)

	// n
	c.setFlagN(false)
	// h // FIXME understand half carry for uint16
	// c // FIXME understand carry for uint16
}
func (c *CPU) Add_HL_SP() {
	fmt.Println("Instruction not implemented: ADD HL SP")
}

// Inc_rr increments 16bit register rr by 1.
func (c *CPU) Inc_rr(rr Reg16) {
	getrr, setrr := c.getReg16(rr)
	setrr(getrr() + 1)
}

func (c *CPU) Inc_SP() {
	fmt.Println("Instruction not implemented: INC SP")
}
func (c *CPU) Dec_rr(rr Reg16) {
	fmt.Println("Instruction not implemented: DEC rr")
}
func (c *CPU) Dec_SP() {
	fmt.Println("Instruction not implemented: DEC sp")
}
func (c *CPU) Add_SP_r8(r8 int8) {
	fmt.Println("Instruction not implemented: ADD SP r8")
}
func (c *CPU) Ld_HL_SPplusr8(r8 int8) {
	fmt.Println("Instruction not implemented: LD HL SP+r8")
}
