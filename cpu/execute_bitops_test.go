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
		{"No carry", 0b0000_0100, 0b0000_1000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry", 0b0000_1000, 0b0001_0000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Carry - carry bit was 1", 0b1000_0100, 0b0000_0001, flags{1, 1, 1, 1}, flags{0, 0, 0, 1}},
		{"Carry - carry bit was 0", 0b1000_0100, 0b0000_0001, flags{1, 1, 1, 0}, flags{0, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.A = tt.regAIn
			setFlags(c, tt.flagsIn)

			c.Rlc_A()

			if c.A != tt.regAOut {
				t.Errorf("Expected A to be %v, got %v", tt.regAOut, c.A)
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
		{"No carry", 0b0000_0100, 0b0000_1000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
		{"Half carry", 0b0000_1000, 0b0001_0000, flags{1, 1, 1, 1}, flags{0, 0, 0, 0}},
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
				t.Errorf("Expected A to be %v, got %v", tt.regAOut, c.A)
			}
			checkFlags(c, tt.flagsOut, t)
		})
	}
}

func TestCPU_Rrc_A(t *testing.T) {
}

func TestCPU_Rr_A(t *testing.T) {
}

func TestCPU_Rlc_r(t *testing.T) {
}

func TestCPU_Rlc_valHL(t *testing.T) {
}

func TestCPU_Rl_r(t *testing.T) {
}

func TestCPU_Rl_valHL(t *testing.T) {
}

func TestCPU_Rrc_r(t *testing.T) {
}

func TestCPU_Rrc_valHL(t *testing.T) {
}

func TestCPU_Rr_r(t *testing.T) {
}

func TestCPU_Rr_valHL(t *testing.T) {
}

func TestCPU_Sla_r(t *testing.T) {
	type args struct {
		r Reg8
	}
}

func TestCPU_Sla_valHL(t *testing.T) {
}

func TestCPU_Swap_r(t *testing.T) {
	type args struct {
		r Reg8
	}
}

func TestCPU_Swap_valHL(t *testing.T) {
}

func TestCPU_Sra_r(t *testing.T) {
}

func TestCPU_Sra_valHL(t *testing.T) {
}

func TestCPU_Srl_r(t *testing.T) {
}

func TestCPU_Srl_valHL(t *testing.T) {
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
