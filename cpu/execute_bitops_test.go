package cpu

import (
	"testing"
)

func TestCPU_Rlc_A(t *testing.T) {
	tests := []struct {
		name     string
		regAIn   byte
		regAOut  byte
		flagsIn  flags
		flagsOut flags
	}{
		{"No carry - carry bit was 0", 0b0000_0100, 0b0000_1000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 1", 0b0000_0100, 0b0000_1000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 0", 0b0000_1000, 0b0001_0000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 1", 0b0000_1000, 0b0001_0000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Carry - carry bit was 1", 0b1000_0000, 0b0000_0001, flags{1, 1, 1, 1}, flags{0, 0, 0, 1}},
		{"Carry - carry bit was 0", 0b1000_0000, 0b0000_0001, flags{1, 1, 1, 0}, flags{0, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.A = tt.regAIn
			setFlags(c, tt.flagsIn)

			c.Rlc_A()

			if c.A != tt.regAOut {
				t.Errorf("Expected A to be 0b%08b, got 0b%08b", tt.regAOut, c.A)
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Rl_A(t *testing.T) {
	tests := []struct {
		name     string
		regAIn   byte
		regAOut  byte
		flagsIn  flags
		flagsOut flags
	}{
		{"No carry - carry bit was 0", 0b0000_0100, 0b0000_1000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 1", 0b0000_0100, 0b0000_1001, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 0", 0b0000_1000, 0b0001_0000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 1", 0b0000_1000, 0b0001_0001, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Carry - carry bit was 1", 0b1000_0000, 0b0000_0001, flags{1, 1, 1, 1}, flags{0, 0, 0, 1}},
		{"Carry - carry bit was 0", 0b1000_0000, 0b0000_0000, flags{1, 1, 1, 0}, flags{0, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.A = tt.regAIn
			setFlags(c, tt.flagsIn)

			c.Rl_A()

			if c.A != tt.regAOut {
				t.Errorf("Expected A to be 0b%08b, got 0b%08b", tt.regAOut, c.A)
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Rrc_A(t *testing.T) {
	tests := []struct {
		name     string
		regAIn   byte
		regAOut  byte
		flagsIn  flags
		flagsOut flags
	}{
		{"No carry - carry bit was 0", 0b0000_0100, 0b0000_0010, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 1", 0b0000_0100, 0b0000_0010, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 0", 0b0001_0000, 0b0000_1000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 1", 0b0001_0000, 0b0000_1000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Carry - carry bit was 1", 0b0000_0001, 0b1000_0000, flags{1, 1, 1, 1}, flags{0, 0, 0, 1}},
		{"Carry - carry bit was 0", 0b0000_0001, 0b1000_0000, flags{1, 1, 1, 0}, flags{0, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.A = tt.regAIn
			setFlags(c, tt.flagsIn)

			c.Rrc_A()

			if c.A != tt.regAOut {
				t.Errorf("Expected A to be 0b%08b, got 0b%08b", tt.regAOut, c.A)
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Rr_A(t *testing.T) {
	tests := []struct {
		name     string
		regAIn   byte
		regAOut  byte
		flagsIn  flags
		flagsOut flags
	}{
		{"No carry - carry bit was 0", 0b0000_0100, 0b0000_0010, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 1", 0b0000_0100, 0b1000_0010, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 0", 0b0001_0000, 0b0000_1000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 1", 0b0001_0000, 0b1000_1000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Carry - carry bit was 0", 0b0000_0001, 0b0000_0000, flags{1, 1, 1, 0}, flags{0, 0, 0, 1}},
		{"Carry - carry bit was 1", 0b0000_0001, 0b1000_0000, flags{1, 1, 1, 1}, flags{0, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.A = tt.regAIn
			setFlags(c, tt.flagsIn)

			c.Rr_A()

			if c.A != tt.regAOut {
				t.Errorf("Expected A to be 0b%08b, got 0b%08b", tt.regAOut, c.A)
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Rlc_r(t *testing.T) {
	type args struct {
		r Reg8
	}
	tests := []struct {
		name string
		args
		regIn    byte
		regOut   byte
		flagsIn  flags
		flagsOut flags
	}{
		{"All 0s", args{RegA}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 1}, flags{1, 0, 0, 0}},
		{"No carry - carry bit was 0", args{RegA}, 0b0000_0100, 0b0000_1000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 1", args{RegB}, 0b0000_0100, 0b0000_1000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 0", args{RegC}, 0b0000_1000, 0b0001_0000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 1", args{RegD}, 0b0000_1000, 0b0001_0000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Carry - carry bit was 1", args{RegH}, 0b1000_0000, 0b0000_0001, flags{1, 1, 1, 1}, flags{0, 0, 0, 1}},
		{"Carry - carry bit was 0", args{RegL}, 0b1000_0000, 0b0000_0001, flags{1, 1, 1, 0}, flags{0, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			getr, setr := c.getReg8(tt.args.r)
			setr(tt.regIn)
			setFlags(c, tt.flagsIn)

			c.Rlc_r(tt.args.r)

			if getr() != tt.regOut {
				t.Errorf("Expected register %v to be 0b%08b, got 0b%08b", tt.args.r, tt.regOut, getr())
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Rlc_valHL(t *testing.T) {
	tests := []struct {
		name string
		Registers
		valHLIn  byte
		valHLOut byte
		flagsIn  flags
		flagsOut flags
	}{
		{"All 0s", Registers{H: 0x10, L: 0x10}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 1}, flags{1, 0, 0, 0}},
		{"No carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_0100, 0b0000_1000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_0100, 0b0000_1000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_1000, 0b0001_0000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_1000, 0b0001_0000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b1000_0000, 0b0000_0001, flags{1, 1, 1, 1}, flags{0, 0, 0, 1}},
		{"Carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b1000_0000, 0b0000_0001, flags{1, 1, 1, 0}, flags{0, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.mem.Wb(c.getHL(), tt.valHLIn)
			setFlags(c, tt.flagsIn)

			c.Rlc_valHL()

			actualValHL := c.mem.Rb(c.getHL())
			if actualValHL != tt.valHLOut {
				t.Errorf("Expected (HL) to be 0b%08b, got 0b%08b", tt.valHLOut, actualValHL)
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Rl_r(t *testing.T) {
	type args struct {
		r Reg8
	}
	tests := []struct {
		name string
		args
		regIn    byte
		regOut   byte
		flagsIn  flags
		flagsOut flags
	}{
		{"All 0s - carry bit was 0", args{RegA}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 0}},
		{"All 0s - carry bit was 1", args{RegA}, 0b0000_0000, 0b0000_0001, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 0", args{RegA}, 0b0000_0100, 0b0000_1000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 1", args{RegB}, 0b0000_0100, 0b0000_1001, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 0", args{RegC}, 0b0000_1000, 0b0001_0000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 1", args{RegD}, 0b0000_1000, 0b0001_0001, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Carry - carry bit was 0", args{RegL}, 0b1000_0000, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 1}},
		{"Carry - carry bit was 1", args{RegH}, 0b1000_0000, 0b0000_0001, flags{1, 1, 1, 1}, flags{0, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			getr, setr := c.getReg8(tt.args.r)
			setr(tt.regIn)
			setFlags(c, tt.flagsIn)

			c.Rl_r(tt.args.r)

			if getr() != tt.regOut {
				t.Errorf("Expected register %v to be 0b%08b, got 0b%08b", tt.args.r, tt.regOut, getr())
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Rl_valHL(t *testing.T) {
	tests := []struct {
		name string
		Registers
		valHLIn  byte
		valHLOut byte
		flagsIn  flags
		flagsOut flags
	}{
		{"All 0s - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 0}},
		{"All 0s - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_0000, 0b0000_0001, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_0100, 0b0000_1000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_0100, 0b0000_1001, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_1000, 0b0001_0000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_1000, 0b0001_0001, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b1000_0000, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 1}},
		{"Carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b1000_0000, 0b0000_0001, flags{1, 1, 1, 1}, flags{0, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.mem.Wb(c.getHL(), tt.valHLIn)
			setFlags(c, tt.flagsIn)

			c.Rl_valHL()

			actualValHL := c.mem.Rb(c.getHL())
			if actualValHL != tt.valHLOut {
				t.Errorf("Expected (HL) to be 0b%08b, got 0b%08b", tt.valHLOut, actualValHL)
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Rrc_r(t *testing.T) {
	type args struct {
		r Reg8
	}
	tests := []struct {
		name string
		args
		regIn    byte
		regOut   byte
		flagsIn  flags
		flagsOut flags
	}{
		{"All 0s", args{RegA}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 1}, flags{1, 0, 0, 0}},
		{"No carry - carry bit was 0", args{RegA}, 0b0000_0100, 0b0000_0010, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 1", args{RegB}, 0b0000_0100, 0b0000_0010, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 0", args{RegC}, 0b0000_1000, 0b0000_0100, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 1", args{RegD}, 0b0000_1000, 0b0000_0100, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Carry - carry bit was 1", args{RegH}, 0b0000_0001, 0b1000_0000, flags{1, 1, 1, 1}, flags{0, 0, 0, 1}},
		{"Carry - carry bit was 0", args{RegL}, 0b0000_0001, 0b1000_0000, flags{1, 1, 1, 0}, flags{0, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			getr, setr := c.getReg8(tt.args.r)
			setr(tt.regIn)
			setFlags(c, tt.flagsIn)

			c.Rrc_r(tt.args.r)

			if getr() != tt.regOut {
				t.Errorf("Expected register %v to be 0b%08b, got 0b%08b", tt.args.r, tt.regOut, getr())
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Rrc_valHL(t *testing.T) {
	tests := []struct {
		name string
		Registers
		valHLIn  byte
		valHLOut byte
		flagsIn  flags
		flagsOut flags
	}{
		{"All 0s", Registers{H: 0x10, L: 0x10}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 1}, flags{1, 0, 0, 0}},
		{"No carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_0100, 0b0000_0010, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_0100, 0b0000_0010, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_1000, 0b0000_0100, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_1000, 0b0000_0100, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_0001, 0b1000_0000, flags{1, 1, 1, 1}, flags{0, 0, 0, 1}},
		{"Carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_0001, 0b1000_0000, flags{1, 1, 1, 0}, flags{0, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.mem.Wb(c.getHL(), tt.valHLIn)
			setFlags(c, tt.flagsIn)

			c.Rrc_valHL()

			actualValHL := c.mem.Rb(c.getHL())
			if actualValHL != tt.valHLOut {
				t.Errorf("Expected (HL) to be 0b%08b, got 0b%08b", tt.valHLOut, actualValHL)
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Rr_r(t *testing.T) {
	type args struct {
		r Reg8
	}
	tests := []struct {
		name string
		args
		regIn    byte
		regOut   byte
		flagsIn  flags
		flagsOut flags
	}{
		{"All 0s - carry bit was 0", args{RegA}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 0}},
		{"All 0s - carry bit was 1", args{RegA}, 0b0000_0000, 0b1000_0000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 0", args{RegA}, 0b0000_0100, 0b0000_0010, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 1", args{RegB}, 0b0000_0100, 0b1000_0010, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 0", args{RegC}, 0b0000_1000, 0b0000_0100, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 1", args{RegD}, 0b0000_1000, 0b1000_0100, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Carry - carry bit was 0", args{RegL}, 0b0000_0001, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 1}},
		{"Carry - carry bit was 1", args{RegH}, 0b0000_0001, 0b1000_0000, flags{1, 1, 1, 1}, flags{0, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			getr, setr := c.getReg8(tt.args.r)
			setr(tt.regIn)
			setFlags(c, tt.flagsIn)

			c.Rr_r(tt.args.r)

			if getr() != tt.regOut {
				t.Errorf("Expected register %v to be 0b%08b, got 0b%08b", tt.args.r, tt.regOut, getr())
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Rr_valHL(t *testing.T) {
	tests := []struct {
		name string
		Registers
		valHLIn  byte
		valHLOut byte
		flagsIn  flags
		flagsOut flags
	}{
		{"All 0s - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 0}},
		{"All 0s - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_0000, 0b1000_0000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_0100, 0b0000_0010, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_0100, 0b1000_0010, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_1000, 0b0000_0100, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_1000, 0b1000_0100, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_0001, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 1}},
		{"Carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_0001, 0b1000_0000, flags{1, 1, 1, 1}, flags{0, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.mem.Wb(c.getHL(), tt.valHLIn)
			setFlags(c, tt.flagsIn)

			c.Rr_valHL()

			actualValHL := c.mem.Rb(c.getHL())
			if actualValHL != tt.valHLOut {
				t.Errorf("Expected (HL) to be 0b%08b, got 0b%08b", tt.valHLOut, actualValHL)
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Sla_r(t *testing.T) {
	type args struct {
		r Reg8
	}
	tests := []struct {
		name string
		args
		regIn    byte
		regOut   byte
		flagsIn  flags
		flagsOut flags
	}{
		{"All 0s - carry bit was 0", args{RegA}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 0}},
		{"All 0s - carry bit was 1", args{RegA}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 1}, flags{1, 0, 0, 0}},
		{"No carry - carry bit was 0", args{RegA}, 0b0000_0100, 0b0000_1000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 1", args{RegB}, 0b0000_0100, 0b0000_1000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 0", args{RegC}, 0b0000_1000, 0b0001_0000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 1", args{RegD}, 0b0000_1000, 0b0001_0000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Carry - carry bit was 0", args{RegL}, 0b1000_0000, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 1}},
		{"Carry - carry bit was 1", args{RegH}, 0b1000_0000, 0b0000_0000, flags{1, 1, 1, 1}, flags{1, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			getr, setr := c.getReg8(tt.args.r)
			setr(tt.regIn)
			setFlags(c, tt.flagsIn)

			c.Sla_r(tt.args.r)

			if getr() != tt.regOut {
				t.Errorf("Expected register %v to be 0b%08b, got 0b%08b", tt.args.r, tt.regOut, getr())
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Sla_valHL(t *testing.T) {
	tests := []struct {
		name string
		Registers
		valHLIn  byte
		valHLOut byte
		flagsIn  flags
		flagsOut flags
	}{
		{"All 0s - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 0}},
		{"All 0s - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 1}, flags{1, 0, 0, 0}},
		{"No carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_0100, 0b0000_1000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_0100, 0b0000_1000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_1000, 0b0001_0000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_1000, 0b0001_0000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b1000_0000, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 1}},
		{"Carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b1000_0000, 0b0000_0000, flags{1, 1, 1, 1}, flags{1, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.mem.Wb(c.getHL(), tt.valHLIn)
			setFlags(c, tt.flagsIn)

			c.Sla_valHL()

			actualValHL := c.mem.Rb(c.getHL())
			if actualValHL != tt.valHLOut {
				t.Errorf("Expected (HL) to be 0b%08b, got 0b%08b", tt.valHLOut, actualValHL)
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Swap_r(t *testing.T) {
	type args struct {
		r Reg8
	}
	tests := []struct {
		name string
		args
		regIn    byte
		regOut   byte
		flagsIn  flags
		flagsOut flags
	}{
		{"All 0s", args{RegA}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 0}},
		{"0xF 0x0 -> 0x0 0xF", args{RegA}, 0b1111_0000, 0b0000_1111, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"0x8 0xF -> 0xF 0x8", args{RegA}, 0b1000_1111, 0b1111_1000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			getr, setr := c.getReg8(tt.args.r)
			setr(tt.regIn)
			setFlags(c, tt.flagsIn)

			c.Swap_r(tt.args.r)

			if getr() != tt.regOut {
				t.Errorf("Expected register %v to be 0b%08b, got 0b%08b", tt.args.r, tt.regOut, getr())
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Swap_valHL(t *testing.T) {
	tests := []struct {
		name string
		Registers
		valHLIn  byte
		valHLOut byte
		flagsIn  flags
		flagsOut flags
	}{
		{"All 0s", Registers{H: 0x10, L: 0x10}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 0}},
		{"0xF 0x0 -> 0x0 0xF", Registers{H: 0x10, L: 0x10}, 0b1111_0000, 0b0000_1111, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"0x8 0xF -> 0xF 0x8", Registers{H: 0x10, L: 0x10}, 0b1000_1111, 0b1111_1000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.mem.Wb(c.getHL(), tt.valHLIn)
			setFlags(c, tt.flagsIn)

			c.Swap_valHL()

			actualValHL := c.mem.Rb(c.getHL())
			if actualValHL != tt.valHLOut {
				t.Errorf("Expected (HL) to be 0b%08b, got 0b%08b", tt.valHLOut, actualValHL)
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Sra_r(t *testing.T) {
	type args struct {
		r Reg8
	}
	tests := []struct {
		name string
		args
		regIn    byte
		regOut   byte
		flagsIn  flags
		flagsOut flags
	}{
		{"All 0s - carry bit was 0", args{RegA}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 0}},
		{"All 0s - carry bit was 1", args{RegA}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 1}, flags{1, 0, 0, 0}},
		{"No carry - carry bit was 0", args{RegA}, 0b0000_0100, 0b0000_0010, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 1", args{RegB}, 0b0000_0100, 0b0000_0010, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 0", args{RegC}, 0b0000_1000, 0b0000_0100, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 1", args{RegD}, 0b0000_1000, 0b0000_0100, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Carry - carry bit was 0", args{RegL}, 0b0000_0001, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 1}},
		{"Carry - carry bit was 1", args{RegH}, 0b0000_0001, 0b0000_0000, flags{1, 1, 1, 1}, flags{1, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			getr, setr := c.getReg8(tt.args.r)
			setr(tt.regIn)
			setFlags(c, tt.flagsIn)

			c.Sra_r(tt.args.r)

			if getr() != tt.regOut {
				t.Errorf("Expected register %v to be 0b%08b, got 0b%08b", tt.args.r, tt.regOut, getr())
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Sra_valHL(t *testing.T) {
	tests := []struct {
		name string
		Registers
		valHLIn  byte
		valHLOut byte
		flagsIn  flags
		flagsOut flags
	}{
		{"All 0s - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 0}},
		{"All 0s - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 1}, flags{1, 0, 0, 0}},
		{"No carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_0100, 0b0000_0010, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_0100, 0b0000_0010, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_1000, 0b0000_0100, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_1000, 0b0000_0100, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_0001, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 1}},
		{"Carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_0001, 0b0000_0000, flags{1, 1, 1, 1}, flags{1, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.mem.Wb(c.getHL(), tt.valHLIn)
			setFlags(c, tt.flagsIn)

			c.Sra_valHL()

			actualValHL := c.mem.Rb(c.getHL())
			if actualValHL != tt.valHLOut {
				t.Errorf("Expected (HL) to be 0b%08b, got 0b%08b", tt.valHLOut, actualValHL)
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Srl_r(t *testing.T) {
	type args struct {
		r Reg8
	}
	tests := []struct {
		name string
		args
		regIn    byte
		regOut   byte
		flagsIn  flags
		flagsOut flags
	}{
		{"All 0s - carry bit was 0", args{RegA}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 0}},
		{"All 0s - carry bit was 1", args{RegA}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 1}, flags{1, 0, 0, 0}},
		{"No carry - carry bit was 0", args{RegA}, 0b0000_0100, 0b0000_0010, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 1", args{RegB}, 0b0000_0100, 0b0000_0010, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 0", args{RegC}, 0b0000_1000, 0b0000_0100, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 1", args{RegD}, 0b0000_1000, 0b0000_0100, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Carry - carry bit was 0", args{RegL}, 0b0000_0001, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 1}},
		{"Carry - carry bit was 1", args{RegH}, 0b0000_0001, 0b0000_0000, flags{1, 1, 1, 1}, flags{1, 0, 0, 1}},
		{"Copy bit 7 - carry bit was 0", args{RegL}, 0b1000_0000, 0b1100_0000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Copy bit 7 - carry bit was 1", args{RegH}, 0b1000_0000, 0b1100_0000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			getr, setr := c.getReg8(tt.args.r)
			setr(tt.regIn)
			setFlags(c, tt.flagsIn)

			c.Srl_r(tt.args.r)

			if getr() != tt.regOut {
				t.Errorf("Expected register %v to be 0b%08b, got 0b%08b", tt.args.r, tt.regOut, getr())
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Srl_valHL(t *testing.T) {
	tests := []struct {
		name string
		Registers
		valHLIn  byte
		valHLOut byte
		flagsIn  flags
		flagsOut flags
	}{
		{"All 0s - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 0}},
		{"All 0s - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_0000, 0b0000_0000, flags{1, 1, 1, 1}, flags{1, 0, 0, 0}},
		{"No carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_0100, 0b0000_0010, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"No carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_0100, 0b0000_0010, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_1000, 0b0000_0100, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Half carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_1000, 0b0000_0100, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Carry - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b0000_0001, 0b0000_0000, flags{1, 1, 1, 0}, flags{1, 0, 0, 1}},
		{"Carry - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b0000_0001, 0b0000_0000, flags{1, 1, 1, 1}, flags{1, 0, 0, 1}},
		{"Copy bit 7 - carry bit was 0", Registers{H: 0x10, L: 0x10}, 0b1000_0000, 0b1100_0000, flags{1, 1, 1, 0}, flags{0, 0, 0, 0}},
		{"Copy bit 7 - carry bit was 1", Registers{H: 0x10, L: 0x10}, 0b1000_0000, 0b1100_0000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.mem.Wb(c.getHL(), tt.valHLIn)
			setFlags(c, tt.flagsIn)

			c.Srl_valHL()

			actualValHL := c.mem.Rb(c.getHL())
			if actualValHL != tt.valHLOut {
				t.Errorf("Expected (HL) to be 0b%08b, got 0b%08b", tt.valHLOut, actualValHL)
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Bit_n_r(t *testing.T) {
	type args struct {
		n uint8
		r Reg8
	}
}

func TestCPU_Bit_n_valHL(t *testing.T) {
	type args struct {
		n uint8
	}
}

func TestCPU_Set_n_r(t *testing.T) {
	type args struct {
		n uint8
		r Reg8
	}
}

func TestCPU_Set_n_valHL(t *testing.T) {
	type args struct {
		n uint8
	}
}

func TestCPU_Res_n_r(t *testing.T) {
	type args struct {
		n uint8
		r Reg8
	}
}

func TestCPU_Res_n_valHL(t *testing.T) {
	type args struct {
		n uint8
	}
}
