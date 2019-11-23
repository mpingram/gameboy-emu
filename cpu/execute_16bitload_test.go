package cpu

import (
	"testing"
)

func TestCPU_Ld_rr_d16(t *testing.T) {
	type args struct {
		rr  Reg16
		d16 uint16
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"BC <- 0x0001", Registers{}, args{rr: RegBC, d16: 0x0001}},
		{"DE <- 0xFFFF", Registers{}, args{rr: RegDE, d16: 0xFFFF}},
		{"HL <- 0x0D0B", Registers{}, args{rr: RegHL, d16: 0x0D0B}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.Registers = tt.regs

			c.Ld_rr_d16(tt.args.rr, tt.args.d16)

			// Expect register to be d16
			getrr, _ := c.getReg16(tt.args.rr)
			if getrr() != tt.args.d16 {
				t.Errorf("Expected rr to be %04x, got %04x", tt.args.d16, getrr())
			}
		})
	}
}

func TestCPU_Ld_SP_d16(t *testing.T) {
	type args struct {
		d16 uint16
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"SP <- 0x0001", Registers{}, args{d16: 0x0001}},
		{"SP <- 0xFFFF", Registers{}, args{d16: 0xFFFF}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.Registers = tt.regs

			c.Ld_SP_d16(tt.args.d16)

			// Expect SP to be d16
			if c.SP != tt.args.d16 {
				t.Errorf("Expected SP to be %04x, got %04x", tt.args.d16, c.SP)
			}
		})
	}
}

func TestCPU_Ld_SP_HL(t *testing.T) {
	tests := []struct {
		name string
		regs Registers
	}{
		{"SP <- HL, H=0x00, L=0x01", Registers{H: 0x00, L: 0x01}},
		{"SP <- HL, H=0xFF, L=0xFF", Registers{H: 0xFF, L: 0xFF}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.Registers = tt.regs

			c.Ld_SP_HL()

			// Expect SP to be HL
			if c.SP != c.getHL() {
				t.Errorf("Expected SP to be %04x, got %04x", c.getHL(), c.SP)
			}
		})
	}
}

func TestCPU_Push_rr(t *testing.T) {
	type args struct {
		rr Reg16
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"SP-=2, (SP)=BC, B=0x00, C=0xF0, SP=0xFFFF", Registers{B: 0x00, C: 0xF0, SP: 0xFFFF}, args{rr: RegBC}},
		{"SP-=2, (SP)=DE, D=0x01, C=0xDE, SP=0x4450", Registers{D: 0x01, C: 0xDE, SP: 0x4450}, args{rr: RegDE}},
		{"SP-=2, (SP)=HL, H=0xFF, L=0xFE, SP=0x0000", Registers{H: 0xFF, L: 0xFE, SP: 0x0000}, args{rr: RegHL}},
		{"SP-=2, (SP)=AF, A=0x04, F=0xF0, SP=0x5000", Registers{A: 0x04, F: 0xF0, SP: 0x5000}, args{rr: RegAF}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.Registers = tt.regs

			c.Push_rr(tt.args.rr)

			// Expect SP to be oldSP - 2
			if c.SP != tt.regs.SP-2 {
				t.Errorf("Expected SP to be %04x, got %04x", tt.regs.SP-2, c.SP)
			}
			// Expect (SP) to be rr
			getrr, _ := c.getReg16(tt.args.rr)
			valSP, err := c.mem.Rw(c.SP)
			if err != nil {
				t.Error(err)
			}
			if valSP != getrr() {
				t.Errorf("Expected (SP) to be %04x, got %04x", getrr(), valSP)
			}
		})
	}
}

func TestCPU_Pop_rr(t *testing.T) {
	type args struct {
		rr Reg16
	}
	tests := []struct {
		name  string
		regs  Registers
		args  args
		valSP uint16
	}{
		{"BC=(SP), SP+=2, (SP)=0x00F0, SP=0xFFFF", Registers{SP: 0xFFFF}, args{rr: RegBC}, 0x00F0},
		{"DE=(SP), SP+=2, (SP)=0x01DE, SP=0x4450", Registers{SP: 0x4450}, args{rr: RegDE}, 0x01DE},
		{"HL=(SP), SP+=2, (SP)=0xFFFE, SP=0x0000", Registers{SP: 0x0000}, args{rr: RegHL}, 0xFFFE},
		{"AF=(SP), SP+=2, (SP)=0x04F0, SP=0x5000", Registers{SP: 0x5000}, args{rr: RegAF}, 0x04F0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.Registers = tt.regs
			// initialize memory
			err := c.mem.Ww(c.SP, tt.valSP)
			if err != nil {
				t.Error(err)
			}

			c.Pop_rr(tt.args.rr)

			// Expect SP to be oldSP + 2
			if c.SP != tt.regs.SP+2 {
				t.Errorf("Expected SP to be %04x, got %04x", tt.regs.SP-2, c.SP)
			}
			// Expect rr to be value of old SP
			getrr, _ := c.getReg16(tt.args.rr)
			if getrr() != tt.valSP {
				t.Errorf("Expected rr to be %04x, got %04x", tt.valSP, getrr())
			}
		})
	}
}
