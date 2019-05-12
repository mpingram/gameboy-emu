package cpu

import (
	"testing"
)

func TestCcf(t *testing.T) {

	tests := []struct {
		// in, out = top nibble of F register
		name string
		in   byte
		out  byte
	}{
		// Expect carry flag to be flipped
		{"0001 -> 0000", 0x1, 0x0},
		{"0000 -> 0001", 0x0, 0x1},

		// Expect N and H flags to be reset
		{"0111 -> 0000", 0x7, 0x0},
		{"0110 -> 0001", 0x6, 0x1},

		// Expect Z flag to be unchanged
		{"1111 -> 1000", 0xF, 0x8},
		{"1110 -> 1001", 0xE, 0x9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CPU{}
			c.setF(tt.in << 4)
			c.Ccf()
			res := c.getF() >> 4
			if res != tt.out {
				t.Errorf("got %04b, want %04b", res, tt.out)
			}
		})
	}
}

func TestScf(t *testing.T) {
	tests := []struct {
		// in, out = top nibble of F register
		name string
		in   byte
		out  byte
	}{
		// Expect carry flag to be set
		{"0001 -> 0001", 0x1, 0x1},
		{"0000 -> 0001", 0x0, 0x1},

		// Expect N and H flags to be reset
		{"0111 -> 0001", 0x7, 0x1},
		{"0110 -> 0001", 0x6, 0x1},

		// Expect Z flag to be unchanged
		{"1111 -> 1001", 0xF, 0x9},
		{"1110 -> 1001", 0xE, 0x9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CPU{}
			c.setF(tt.in << 4)
			c.Scf()
			res := c.getF() >> 4
			if res != tt.out {
				t.Errorf("got %04b, want %04b", res, tt.out)
			}
		})
	}
}

func TestNop(t *testing.T) {
	// not sure what needs testing here
}

func TestHalt(t *testing.T) {
	// NOTE there are bugs/subtleties in the behavior of halt
	// See https://github.com/AntonioND/giibiiadvance/blob/master/docs/TCAGBD.pdf
	//
	// From TCAGBD:
	// `HALT mode is exited when a flag in memory register IF is set and the
	// corresponding flag in IE is also set, regardless of the value of IME. The only difference
	// is that IME=1 will make the CPU jump to the interrupt vector(and clear the IF flag), while IME=0
	// will only make the CPU continue executing instructions, but the jump won't be performed (and the
	// IF flag won't be cleared.)
	//
	// ...
	//
	// The HALT instruction has three different behaviors depending on IME, IE, and IF. It behaves the same
	// way in all Gameboy models.
	//
	// - IME = 1
	//  HALT executed normally. CPU stops executing instructions until (IE & IF & 1F) != 0. When a flag in
	//	IF is set and the corresponding flag in IE is also set, the CPU jumps to the interrupt vector. The
	//	return address pushed to the stack is the next instruction to the halt, not the halt itself. The IF
	//  flag corresponding to the vector the CPU has jumped in is cleared.
	// - IME = 0
	// 	- (IE & IF & 0x1F) = 0
	//			HALT mode is entered. It works like the IME = 1 case, but when a IF flag is set and the corresponding
	//			IE flag is also set, the CPU doesn't jump to the interrupt vector, it just continues executing instructions.
	//      The IF flags aren't cleared.
	//	- (IE & IF & 0x1F) != 0
	//			HALT mode is not entered. HALT bug occurs. The CPU fails to increase PC when executing the next instruction. The
	//			IF flags aren't cleared. This results in weird behavior. For example:
	//			> halt
	//			> ld a, $14 // $3E $14 is executed as $3E $3E $14 (ld $3E; inc D)
	// `

	// Expect CPU to be enter halt state
	c := CPU{}
	c.halted = false
	c.Halt()
	if c.halted != true {
		t.Error("Expected CPU to halt")
	}

	// FIXME figure out how to test halt bug / interrupts. Maybe I should test interrupts separately
}

func TestDi(t *testing.T) {

	tests := []struct {
		// in, out = IME bit
		name string
		in   bool
		out  bool
	}{
		// Expect IME bit to be reset in this clock cycle
		{"1 -> 0", true, false},
		{"0 -> 0", false, false},
	}

	for _, tt := range tests {
		c := CPU{}
		c.ime = tt.in
		c.Di()
		if c.ime != tt.out {
			t.Errorf("Expected IME bit to be %v, got %v", tt.out, c.ime)
		}
	}
}

func TestEi(t *testing.T) {
	// TODO test Ei behavior of enable interrupts on next machine cycle
	// TODO test Ei behavior interaction with Di
	tests := []struct {
		name      string
		imeIn     bool
		imeOut    bool
		setIMEin  bool
		setIMEout bool
	}{
		// Expect IME bit to be unchanged in this clock cycle
		// and setIME internal flag to be set
		{"IME: (1 -> 1), setIME: 1 -> 1", true, true, true, true},
		{"IME: (0 -> 0), setIME: 0 -> 1", false, false, true, true},
	}

	for _, tt := range tests {
		c := CPU{}
		c.ime = tt.imeIn
		c.setIME = tt.setIMEin
		c.Ei()
		if c.ime != tt.imeOut {
			t.Errorf("Expected IME bit to be %v, got %v", tt.imeOut, c.ime)
		}
		if c.setIME != tt.setIMEout {
			t.Errorf("Expected IME bit to be %v, got %v", tt.setIMEout, c.setIME)
		}
	}
}
