package cpu

import (
	"testing"
)

type flags struct {
	z, n, h, c byte
}

func checkFlags(cpu *CPU, f flags, t *testing.T) {
	if f.z != (cpu.F&0x80)>>7 {
		t.Errorf("Got flags %04b, wanted Z to be %v", cpu.F>>4, f.z)
	}
	if f.n != (cpu.F&0x40)>>6 {
		t.Errorf("Got flags %04b, wanted N to be %v", cpu.F>>4, f.n)
	}
	if f.h != (cpu.F&0x20)>>5 {
		t.Errorf("Got flags %04b, wanted H to be %v", cpu.F>>4, f.h)
	}
	if f.c != (cpu.F&0x10)>>4 {
		t.Errorf("Got flags %04b, wanted C to be %v", cpu.F>>4, f.c)
	}
}

func TestCPU_Add_r(t *testing.T) {
	type args struct {
		r Reg8
	}
	tests := []struct {
		name  string
		regs  Registers
		args  args
		flags flags
	}{
		{"A+=B, A=0x01, B=0x02", Registers{A: 0x01, B: 0x02}, args{r: RegB}, flags{0, 0, 0, 0}},
		{"A+=B, A=0x00, B=0x00", Registers{A: 0x00, B: 0x00}, args{r: RegB}, flags{1, 0, 0, 0}},
		{"A+=B, A=0x0F, B=0x01", Registers{A: 0x0F, B: 0x01}, args{r: RegB}, flags{0, 0, 1, 0}},
		{"A+=B, A=0x01, B=0xFF", Registers{A: 0x01, B: 0xFF}, args{r: RegB}, flags{1, 0, 1, 1}},
		{"A+=B, A=0xF0, B=0x10", Registers{A: 0xF0, B: 0x10}, args{r: RegB}, flags{1, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			c.Registers = tt.regs

			c.Add_r(tt.args.r)

			// Expect A=A+r
			getr, _ := c.getReg8(tt.args.r)
			expected := tt.regs.A + getr()
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Add_d8(t *testing.T) {
	type args struct {
		d8 byte
	}
	tests := []struct {
		name  string
		regs  Registers
		args  args
		flags flags
	}{
		{"A+=0x02, A=0x01", Registers{A: 0x01}, args{d8: 0x02}, flags{0, 0, 0, 0}},
		{"A+=0x00, A=0x00", Registers{A: 0x00}, args{d8: 0x00}, flags{1, 0, 0, 0}},
		{"A+=0x01, A=0x0F", Registers{A: 0x0F}, args{d8: 0x01}, flags{0, 0, 1, 0}},
		{"A+=0xFF, A=0x01", Registers{A: 0x01}, args{d8: 0xFF}, flags{1, 0, 1, 1}},
		{"A+=0x10, A=0xF0", Registers{A: 0xF0}, args{d8: 0x10}, flags{1, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			c.Registers = tt.regs

			c.Add_d8(tt.args.d8)

			// Expect A=A+r
			expected := tt.regs.A + tt.args.d8
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Add_valHL(t *testing.T) {
	tests := []struct {
		name  string
		regs  Registers
		valHL byte
		flags flags
	}{
		{"A+=(HL), (HL)=0x02, A=0x01, H=0x0B, L=0x0A", Registers{A: 0x01, H: 0x0B, L: 0x0A}, 0x02, flags{0, 0, 0, 0}},
		{"A+=(HL), (HL)=0x00, A=0x00, H=0x0B, L=0x0A", Registers{A: 0x00, H: 0x0B, L: 0x0A}, 0x00, flags{1, 0, 0, 0}},
		{"A+=(HL), (HL)=0x01, A=0x0F, H=0x0B, L=0x0A", Registers{A: 0x0F, H: 0x0B, L: 0x0A}, 0x01, flags{0, 0, 1, 0}},
		{"A+=(HL), (HL)=0xFF, A=0x01, H=0x0B, L=0x0A", Registers{A: 0x01, H: 0x0B, L: 0x0A}, 0xFF, flags{1, 0, 1, 1}},
		{"A+=(HL), (HL)=0x10, A=0xF0, H=0x0B, L=0x0A", Registers{A: 0xF0, H: 0x0B, L: 0x0A}, 0x10, flags{1, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			c.Registers = tt.regs
			// initialize memory
			err := c.mem.Wb(c.getHL(), tt.valHL)
			if err != nil {
				t.Error(err)
			}

			c.Add_valHL()

			// Expect A=A+r
			expected := tt.regs.A + tt.valHL
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Adc_r(t *testing.T) {
	type args struct {
		r Reg8
	}
	tests := []struct {
		name  string
		regs  Registers
		args  args
		flags flags
	}{
		{"A+=B, A=0x01, B=0x02, carry=1", Registers{A: 0x01, B: 0x02, F: 0x10}, args{r: RegB}, flags{0, 0, 0, 0}},
		{"A+=B, A=0x00, B=0x00, carry=0", Registers{A: 0x00, B: 0x00, F: 0x00}, args{r: RegB}, flags{1, 0, 0, 0}},
		{"A+=B, A=0x0F, B=0x01, carry=1", Registers{A: 0x0F, B: 0x01, F: 0x10}, args{r: RegB}, flags{0, 0, 1, 0}},
		{"A+=B, A=0x00, B=0xFF, carry=1", Registers{A: 0x01, B: 0xFF, F: 0x10}, args{r: RegB}, flags{1, 0, 1, 1}},
		{"A+=B, A=0xF0, B=0x0F, carry=1", Registers{A: 0xF0, B: 0x10, F: 0x10}, args{r: RegB}, flags{1, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			c.Registers = tt.regs

			c.Adc_r(tt.args.r)

			// Expect A=A+r
			getr, _ := c.getReg8(tt.args.r)
			expected := tt.regs.A + getr()
			// if carry, add 1
			if (c.F & 0x10) != 0 {
				expected++
			}
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Adc_d8(t *testing.T) {
	type args struct {
		d8 byte
	}
	tests := []struct {
		name  string
		regs  Registers
		args  args
		flags flags
	}{
		{"A+=d8, A=0x01, d8=0x02, carry=1", Registers{A: 0x01, F: 0x10}, args{d8: 0x02}, flags{0, 0, 0, 0}},
		{"A+=d8, A=0x00, d8=0x00, carry=0", Registers{A: 0x00, F: 0x00}, args{d8: 0x00}, flags{1, 0, 0, 0}},
		{"A+=d8, A=0x0F, d8=0x01, carry=1", Registers{A: 0x0F, F: 0x10}, args{d8: 0x01}, flags{0, 0, 1, 0}},
		{"A+=d8, A=0x00, d8=0xFF, carry=1", Registers{A: 0x01, F: 0x10}, args{d8: 0xFF}, flags{1, 0, 1, 1}},
		{"A+=d8, A=0xF0, d8=0x0F, carry=1", Registers{A: 0xF0, F: 0x10}, args{d8: 0x0F}, flags{1, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			c.Registers = tt.regs

			c.Adc_d8(tt.args.d8)

			// Expect A=A+r
			expected := tt.regs.A + tt.args.d8
			// if carry, add 1
			if (c.F & 0x10) != 0 {
				expected++
			}
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Adc_valHL(t *testing.T) {
	tests := []struct {
		name  string
		regs  Registers
		valHL byte
		flags flags
	}{
		{"A+=(HL), A=0x01, d8=0x02, carry=1", Registers{A: 0x01, F: 0x10}, 0x02, flags{0, 0, 0, 0}},
		{"A+=(HL), A=0x00, d8=0x00, carry=0", Registers{A: 0x00, F: 0x00}, 0x00, flags{1, 0, 0, 0}},
		{"A+=(HL), A=0x0F, d8=0x01, carry=1", Registers{A: 0x0F, F: 0x10}, 0x01, flags{0, 0, 1, 0}},
		{"A+=(HL), A=0x00, d8=0xFF, carry=1", Registers{A: 0x01, F: 0x10}, 0xFF, flags{1, 0, 1, 1}},
		{"A+=(HL), A=0xF0, d8=0x0F, carry=1", Registers{A: 0xF0, F: 0x10}, 0x0F, flags{1, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			c.Registers = tt.regs
			// initialize memory
			err := c.mem.Wb(c.getHL(), tt.valHL)
			if err != nil {
				t.Error(err)
			}

			c.Adc_valHL()

			// Expect A=A+r
			expected := tt.regs.A + tt.valHL
			// if carry, add 1
			if (c.F & 0x10) != 0 {
				expected++
			}
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Sub_r(t *testing.T) {
	type args struct {
		r Reg8
	}
	tests := []struct {
		name  string
		regs  Registers
		args  args
		flags flags
	}{
		{"A-=B, A=0x01, B=0x02", Registers{A: 0x01, B: 0x02}, args{r: RegB}, flags{0, 1, 0, 1}},
		{"A-=B, A=0x00, B=0x00", Registers{A: 0x00, B: 0x00}, args{r: RegB}, flags{1, 1, 0, 0}},
		{"A-=B, A=0x10, B=0x01", Registers{A: 0x10, B: 0x01}, args{r: RegB}, flags{0, 1, 1, 0}},
		{"A-=B, A=0x01, B=0xFF", Registers{A: 0x01, B: 0xFF}, args{r: RegB}, flags{0, 1, 1, 1}},
		{"A-=B, A=0xF0, B=0x10", Registers{A: 0xF0, B: 0x10}, args{r: RegB}, flags{0, 1, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			c.Registers = tt.regs

			c.Sub_r(tt.args.r)

			// Expect A=A-r
			getr, _ := c.getReg8(tt.args.r)
			expected := tt.regs.A - getr()
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Sub_d8(t *testing.T) {
	type args struct {
		d8 byte
	}
	tests := []struct {
		name  string
		regs  Registers
		args  args
		flags flags
	}{
		{"A-=d8, A=0x01, d8=0x02", Registers{A: 0x01}, args{d8: 0x02}, flags{0, 1, 0, 1}},
		{"A-=d8, A=0x00, d8=0x00", Registers{A: 0x00}, args{d8: 0x00}, flags{1, 1, 0, 0}},
		{"A-=d8, A=0x10, d8=0x01", Registers{A: 0x10}, args{d8: 0x01}, flags{0, 1, 1, 0}},
		{"A-=d8, A=0x01, d8=0xFF", Registers{A: 0x01}, args{d8: 0xFF}, flags{0, 1, 1, 1}},
		{"A-=d8, A=0xF0, d8=0x10", Registers{A: 0xF0}, args{d8: 0x10}, flags{0, 1, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			c.Registers = tt.regs

			c.Sub_d8(tt.args.d8)

			// Expect A=A-d8
			expected := tt.regs.A - tt.args.d8
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Sub_valHL(t *testing.T) {
	tests := []struct {
		name  string
		regs  Registers
		valHL byte
		flags flags
	}{
		{"A-=(HL), A=0x01, (HL)=0x02", Registers{A: 0x01}, 0x02, flags{0, 1, 0, 1}},
		{"A-=(HL), A=0x00, (HL)=0x00", Registers{A: 0x00}, 0x00, flags{1, 1, 0, 0}},
		{"A-=(HL), A=0x10, (HL)=0x01", Registers{A: 0x10}, 0x01, flags{0, 1, 1, 0}},
		{"A-=(HL), A=0x01, (HL)=0xFF", Registers{A: 0x01}, 0xFF, flags{0, 1, 1, 1}},
		{"A-=(HL), A=0xF0, (HL)=0x10", Registers{A: 0xF0}, 0x10, flags{0, 1, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			c.Registers = tt.regs
			// initialize memory
			err := c.mem.Wb(c.getHL(), tt.valHL)
			if err != nil {
				t.Error(err)
			}

			c.Sub_valHL()

			// Expect A=A-(HL)
			expected := tt.regs.A - tt.valHL
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Sbc_r(t *testing.T) {
	type args struct {
		r Reg8
	}
	tests := []struct {
		name  string
		regs  Registers
		args  args
		flags flags
	}{
		{"A-=B, A=0x01, B=0x02, carry=0", Registers{A: 0x01, B: 0x02, F: 0x00}, args{r: RegB}, flags{0, 1, 0, 1}},
		{"A-=B, A=0x00, B=0x00, carry=1", Registers{A: 0x00, B: 0x00, F: 0x10}, args{r: RegB}, flags{0, 1, 1, 1}},
		{"A-=B, A=0x10, B=0x01, carry=0", Registers{A: 0x10, B: 0x01, F: 0x00}, args{r: RegB}, flags{0, 1, 1, 0}},
		{"A-=B, A=0x01, B=0xFF, carry=1", Registers{A: 0x01, B: 0xFF, F: 0x10}, args{r: RegB}, flags{1, 1, 1, 1}},
		{"A-=B, A=0xF0, B=0x0F, carry=1", Registers{A: 0xF0, B: 0x10, F: 0x10}, args{r: RegB}, flags{0, 1, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			c.Registers = tt.regs

			c.Sbc_r(tt.args.r)

			// Expect A=A-r-carry
			getr, _ := c.getReg8(tt.args.r)
			expected := tt.regs.A - getr()
			// if carry, subtract one more from expected
			if c.F&0x10 != 0 {
				expected--
			}
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Sbc_d8(t *testing.T) {
}

func TestCPU_Sbc_valHL(t *testing.T) {
}

func TestCPU_And_r(t *testing.T) {
}

func TestCPU_And_d8(t *testing.T) {
}

func TestCPU_And_valHL(t *testing.T) {
}

func TestCPU_Xor_r(t *testing.T) {
}

func TestCPU_Xor_d8(t *testing.T) {
}

func TestCPU_Xor_valHL(t *testing.T) {
}

func TestCPU_Or_r(t *testing.T) {
}

func TestCPU_Or_d8(t *testing.T) {
}

func TestCPU_Or_valHL(t *testing.T) {
}

func TestCPU_Cp_r(t *testing.T) {
}

func TestCPU_Cp_d8(t *testing.T) {
}

func TestCPU_Cp_valHL(t *testing.T) {
}

func TestCPU_Inc_r(t *testing.T) {

}

func TestCPU_Inc_d8(t *testing.T) {
}

func TestCPU_Inc_valHL(t *testing.T) {
}

func TestCPU_Dec_r(t *testing.T) {
}

func TestCPU_Dec_d8(t *testing.T) {
}

func TestCPU_Dec_valHL(t *testing.T) {
}

func TestCPU_Daa(t *testing.T) {
}

func TestCPU_Cpl(t *testing.T) {
}
