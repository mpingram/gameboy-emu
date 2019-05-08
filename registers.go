package cpu

const subFlag = 0x10
const halfCarryFlag = 0x20
const carryFlag = 0x40
const zeroFlag = 0x80

// Registers represents the Sharp LR35902's registers.
type registers struct {
	A                uint8 // accumulator
	B, C, D, E, H, L uint8 // gen purpose
	F                uint8 // flag register

	SP uint16 // stack pointer

	PC uint16 // program counter
}

// Each 8-bit register can be combined with its complement
// to function as one 16-bit register.
// To emulate this, we use getters and setters that
// change the state of the underlying registers.
func (r *registers) getAF() uint16 {
	return uint16(r.A)<<8 | uint16(r.F)
}
func (r *registers) getBC() uint16 {
	return uint16(r.B)<<8 | uint16(r.C)
}
func (r *registers) getDE() uint16 {
	return uint16(r.D)<<8 | uint16(r.E)
}
func (r *registers) getHL() uint16 {
	return uint16(r.H)<<8 | uint16(r.L)
}

func (r *registers) setAF(val uint16) {
	r.A = uint8((val & 0xFF00) >> 8)
	r.F = uint8(val & 0x00FF)
}
func (r *registers) setBC(val uint16) {
	r.B = uint8((val & 0xFF00) >> 8)
	r.C = uint8(val & 0x00FF)
}
func (r *registers) setDE(val uint16) {
	r.D = uint8((val & 0xFF00) >> 8)
	r.E = uint8(val & 0x00FF)
}
func (r *registers) setHL(val uint16) {
	r.H = uint8((val & 0xFF00) >> 8)
	r.L = uint8(val & 0x00FF)
}
