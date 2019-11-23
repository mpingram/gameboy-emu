package cpu

import (
	"testing"
)

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
			c, _ := testSetup()
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
			c, _ := testSetup()
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
			c, _ := testSetup()
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
		{"A+=B, A=0x00, B=0xFF, carry=1", Registers{A: 0x00, B: 0xFF, F: 0x10}, args{r: RegB}, flags{1, 0, 1, 1}},
		{"A+=B, A=0xF0, B=0x0F, carry=1", Registers{A: 0xF0, B: 0x0F, F: 0x10}, args{r: RegB}, flags{1, 0, 1, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.Registers = tt.regs

			prevCarry := c.F&0x10 != 0

			c.Adc_r(tt.args.r)

			// Expect A=A+r
			getr, _ := c.getReg8(tt.args.r)
			expected := tt.regs.A + getr()
			// if carry, add 1
			if prevCarry {
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
		{"A+=d8, A=0x00, d8=0xFF, carry=1", Registers{A: 0x00, F: 0x10}, args{d8: 0xFF}, flags{1, 0, 1, 1}},
		{"A+=d8, A=0xF0, d8=0x0F, carry=1", Registers{A: 0xF0, F: 0x10}, args{d8: 0x0F}, flags{1, 0, 1, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.Registers = tt.regs

			carry := c.F&0x10 != 0
			c.Adc_d8(tt.args.d8)

			// Expect A=A+r
			expected := tt.regs.A + tt.args.d8
			if carry {
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
		{"A+=(HL), A=0x00, d8=0xFF, carry=1", Registers{A: 0x00, F: 0x10}, 0xFF, flags{1, 0, 1, 1}},
		{"A+=(HL), A=0xF0, d8=0x0F, carry=1", Registers{A: 0xF0, F: 0x10}, 0x0F, flags{1, 0, 1, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.Registers = tt.regs
			// initialize memory
			err := c.mem.Wb(c.getHL(), tt.valHL)
			if err != nil {
				t.Error(err)
			}
			carry := c.F&0x10 != 0

			c.Adc_valHL()

			// Expect A=A+r
			expected := tt.regs.A + tt.valHL
			// if carry, add 1
			if carry {
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
		{"A-=B, A=0x01, B=0x02", Registers{A: 0x01, B: 0x02}, args{r: RegB}, flags{0, 1, 1, 1}},
		{"A-=B, A=0x00, B=0x00", Registers{A: 0x00, B: 0x00}, args{r: RegB}, flags{1, 1, 0, 0}},
		{"A-=B, A=0x10, B=0x01", Registers{A: 0x10, B: 0x01}, args{r: RegB}, flags{0, 1, 1, 0}},
		{"A-=B, A=0x01, B=0xFF", Registers{A: 0x01, B: 0xFF}, args{r: RegB}, flags{0, 1, 1, 1}},
		{"A-=B, A=0xF0, B=0x10", Registers{A: 0xF0, B: 0x10}, args{r: RegB}, flags{0, 1, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
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
		{"A-=d8, A=0x01, d8=0x02", Registers{A: 0x01}, args{d8: 0x02}, flags{0, 1, 1, 1}},
		{"A-=d8, A=0x00, d8=0x00", Registers{A: 0x00}, args{d8: 0x00}, flags{1, 1, 0, 0}},
		{"A-=d8, A=0x10, d8=0x01", Registers{A: 0x10}, args{d8: 0x01}, flags{0, 1, 1, 0}},
		{"A-=d8, A=0x01, d8=0xFF", Registers{A: 0x01}, args{d8: 0xFF}, flags{0, 1, 1, 1}},
		{"A-=d8, A=0xF0, d8=0x10", Registers{A: 0xF0}, args{d8: 0x10}, flags{0, 1, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
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
		{"A-=(HL), A=0x01, (HL)=0x02", Registers{A: 0x01}, 0x02, flags{0, 1, 1, 1}},
		{"A-=(HL), A=0x00, (HL)=0x00", Registers{A: 0x00}, 0x00, flags{1, 1, 0, 0}},
		{"A-=(HL), A=0x10, (HL)=0x01", Registers{A: 0x10}, 0x01, flags{0, 1, 1, 0}},
		{"A-=(HL), A=0x01, (HL)=0xFF", Registers{A: 0x01}, 0xFF, flags{0, 1, 1, 1}},
		{"A-=(HL), A=0xF0, (HL)=0x10", Registers{A: 0xF0}, 0x10, flags{0, 1, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
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
		{"A-=B, A=0x01, B=0x02, carry=0", Registers{A: 0x01, B: 0x02, F: 0x00}, args{r: RegB}, flags{0, 1, 1, 1}},
		{"A-=B, A=0x00, B=0x00, carry=1", Registers{A: 0x00, B: 0x00, F: 0x10}, args{r: RegB}, flags{0, 1, 1, 1}},
		{"A-=B, A=0x10, B=0x01, carry=0", Registers{A: 0x10, B: 0x01, F: 0x00}, args{r: RegB}, flags{0, 1, 1, 0}},
		{"A-=B, A=0x00, B=0xFF, carry=1", Registers{A: 0x00, B: 0xFF, F: 0x10}, args{r: RegB}, flags{1, 1, 1, 1}},
		{"A-=B, A=0xF0, B=0x0F, carry=1", Registers{A: 0xF0, B: 0x0F, F: 0x10}, args{r: RegB}, flags{0, 1, 1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.Registers = tt.regs
			carry := c.F&0x10 != 0

			c.Sbc_r(tt.args.r)

			// Expect A=A-r-carry
			getr, _ := c.getReg8(tt.args.r)
			expected := tt.regs.A - getr()
			// if carry, subtract one more from expected
			if carry {
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
	type args struct {
		d8 byte
	}
	tests := []struct {
		name  string
		regs  Registers
		args  args
		flags flags
	}{
		{"A-=d8, A=0x01, d8=0x02, carry=0", Registers{A: 0x01, F: 0x00}, args{d8: 0x02}, flags{0, 1, 1, 1}},
		{"A-=d8, A=0x00, d8=0x00, carry=1", Registers{A: 0x00, F: 0x10}, args{d8: 0x00}, flags{0, 1, 1, 1}},
		{"A-=d8, A=0x10, d8=0x01, carry=0", Registers{A: 0x10, F: 0x00}, args{d8: 0x01}, flags{0, 1, 1, 0}},
		{"A-=d8, A=0x00, d8=0xFF, carry=1", Registers{A: 0x00, F: 0x10}, args{d8: 0xFF}, flags{1, 1, 1, 1}},
		{"A-=d8, A=0xF0, d8=0x0F, carry=1", Registers{A: 0xF0, F: 0x10}, args{d8: 0x0F}, flags{0, 1, 1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.Registers = tt.regs
			carry := c.F&0x10 != 0

			c.Sbc_d8(tt.args.d8)

			// Expect A=A-d8-carry
			expected := tt.regs.A - tt.args.d8
			if carry {
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

func TestCPU_Sbc_valHL(t *testing.T) {
	tests := []struct {
		name  string
		regs  Registers
		valHL byte
		flags flags
	}{
		{"A-=(HL), A=0x01, (HL)=0x02, carry=0", Registers{A: 0x01, F: 0x00}, 0x02, flags{0, 1, 1, 1}},
		{"A-=(HL), A=0x00, (HL)=0x00, carry=1", Registers{A: 0x00, F: 0x10}, 0x00, flags{0, 1, 1, 1}},
		{"A-=(HL), A=0x10, (HL)=0x01, carry=0", Registers{A: 0x10, F: 0x00}, 0x01, flags{0, 1, 1, 0}},
		{"A-=(HL), A=0x00, (HL)=0xFF, carry=1", Registers{A: 0x00, F: 0x10}, 0xFF, flags{1, 1, 1, 1}},
		{"A-=(HL), A=0xF0, (HL)=0x0F, carry=1", Registers{A: 0xF0, F: 0x10}, 0x0F, flags{0, 1, 1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()
			c.Registers = tt.regs
			// initialize memory
			err := c.mem.Wb(c.getHL(), tt.valHL)
			if err != nil {
				t.Error(err)
			}
			carry := c.F&0x10 != 0

			c.Sbc_valHL()

			// Expect A=A-(HL)-carry
			expected := tt.regs.A - tt.valHL
			if carry {
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

func TestCPU_And_r(t *testing.T) {
	type args struct {
		r Reg8
	}
	tests := []struct {
		name  string
		regs  Registers
		args  args
		flags flags
	}{
		{"A&=A, A=0xF0", Registers{A: 0xF0}, args{r: RegA}, flags{0, 0, 1, 0}},
		{"A&=B, A=0x01, B=0x03", Registers{A: 0x01, B: 0x03}, args{r: RegB}, flags{0, 0, 1, 0}},
		{"A&=C, A=0x00, C=0x00", Registers{A: 0x00, C: 0x00}, args{r: RegC}, flags{1, 0, 1, 0}},
		{"A&=F, A=0x10, F=0x01", Registers{A: 0x10, F: 0x01}, args{r: RegF}, flags{1, 0, 1, 0}},
		{"A&=H, A=0x01, H=0xFF", Registers{A: 0x01, H: 0xFF}, args{r: RegH}, flags{0, 0, 1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.And_r(tt.args.r)

			// Expect A=A&r
			getr, _ := c.getReg8(tt.args.r)
			expected := tt.regs.A & getr()
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_And_d8(t *testing.T) {
	type args struct {
		d8 byte
	}
	tests := []struct {
		name  string
		regs  Registers
		args  args
		flags flags
	}{
		{"A&=d8, A=0xF0, d8=0xF0", Registers{A: 0xF0}, args{d8: 0xF0}, flags{0, 0, 1, 0}},
		{"A&=d8, A=0x01, d8=0x03", Registers{A: 0x01}, args{d8: 0x03}, flags{0, 0, 1, 0}},
		{"A&=d8, A=0x00, d8=0x00", Registers{A: 0x00}, args{d8: 0x00}, flags{1, 0, 1, 0}},
		{"A&=d8, A=0x10, d8=0x01", Registers{A: 0x10}, args{d8: 0x01}, flags{1, 0, 1, 0}},
		{"A&=d8, A=0x01, d8=0xFF", Registers{A: 0x01}, args{d8: 0xFF}, flags{0, 0, 1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.And_d8(tt.args.d8)

			// Expect A=A&r
			expected := tt.regs.A & tt.args.d8
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_And_valHL(t *testing.T) {
	tests := []struct {
		name  string
		regs  Registers
		valHL byte
		flags flags
	}{
		{"A&=(HL), A=0xF0, (HL)=0xF0", Registers{A: 0xF0}, 0xF0, flags{0, 0, 1, 0}},
		{"A&=(HL), A=0x01, (HL)=0x03", Registers{A: 0x01}, 0x03, flags{0, 0, 1, 0}},
		{"A&=(HL), A=0x00, (HL)=0x00", Registers{A: 0x00}, 0x00, flags{1, 0, 1, 0}},
		{"A&=(HL), A=0x10, (HL)=0x01", Registers{A: 0x10}, 0x01, flags{1, 0, 1, 0}},
		{"A&=(HL), A=0x01, (HL)=0xFF", Registers{A: 0x01}, 0xFF, flags{0, 0, 1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs
			// initialize memory
			err := c.mem.Wb(c.getHL(), tt.valHL)
			if err != nil {
				t.Error(err)
			}

			c.And_valHL()

			// Expect A=A&(HL)
			expected := tt.regs.A & tt.valHL
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Xor_r(t *testing.T) {
	type args struct {
		r Reg8
	}
	tests := []struct {
		name  string
		regs  Registers
		args  args
		flags flags
	}{
		{"A^=A, A=0xF0", Registers{A: 0xF0}, args{r: RegA}, flags{1, 0, 0, 0}},
		{"A^=B, A=0x01, B=0x03", Registers{A: 0x01, B: 0x03}, args{r: RegB}, flags{0, 0, 0, 0}},
		{"A^=C, A=0x00, C=0x00", Registers{A: 0x00, C: 0x00}, args{r: RegC}, flags{1, 0, 0, 0}},
		{"A^=F, A=0x10, F=0x01", Registers{A: 0x10, F: 0x01}, args{r: RegF}, flags{0, 0, 0, 0}},
		{"A^=H, A=0x01, H=0xFF", Registers{A: 0x01, H: 0xFF}, args{r: RegH}, flags{0, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs
			get, _ := c.getReg8(tt.args.r)
			valR := get()

			c.Xor_r(tt.args.r)

			// Expect A=A^r
			expected := tt.regs.A ^ valR
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Xor_d8(t *testing.T) {
	type args struct {
		d8 byte
	}
	tests := []struct {
		name  string
		regs  Registers
		args  args
		flags flags
	}{
		{"A^=d8, A=0xF0, d8=0xF0", Registers{A: 0xF0}, args{d8: 0xF0}, flags{1, 0, 0, 0}},
		{"A^=d8, A=0x01, d8=0x03", Registers{A: 0x01}, args{d8: 0x03}, flags{0, 0, 0, 0}},
		{"A^=d8, A=0x00, d8=0x00", Registers{A: 0x00}, args{d8: 0x00}, flags{1, 0, 0, 0}},
		{"A^=d8, A=0x10, d8=0x01", Registers{A: 0x10}, args{d8: 0x01}, flags{0, 0, 0, 0}},
		{"A^=d8, A=0x01, d8=0xFF", Registers{A: 0x01}, args{d8: 0xFF}, flags{0, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.Xor_d8(tt.args.d8)

			// Expect A=A^r
			expected := tt.regs.A ^ tt.args.d8
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Xor_valHL(t *testing.T) {
	tests := []struct {
		name  string
		regs  Registers
		valHL byte
		flags flags
	}{
		{"A^=(HL), A=0xF0, (HL)=0xF0", Registers{A: 0xF0}, 0xF0, flags{1, 0, 0, 0}},
		{"A^=(HL), A=0x01, (HL)=0x03", Registers{A: 0x01}, 0x03, flags{0, 0, 0, 0}},
		{"A^=(HL), A=0x00, (HL)=0x00", Registers{A: 0x00}, 0x00, flags{1, 0, 0, 0}},
		{"A^=(HL), A=0x10, (HL)=0x01", Registers{A: 0x10}, 0x01, flags{0, 0, 0, 0}},
		{"A^=(HL), A=0x01, (HL)=0xFF", Registers{A: 0x01}, 0xFF, flags{0, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs
			// initialize memory
			err := c.mem.Wb(c.getHL(), tt.valHL)
			if err != nil {
				t.Error(err)
			}

			c.Xor_valHL()

			// Expect A=A^(HL)
			expected := tt.regs.A ^ tt.valHL
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Or_r(t *testing.T) {
	type args struct {
		r Reg8
	}
	tests := []struct {
		name  string
		regs  Registers
		args  args
		flags flags
	}{
		{"A|=A, A=0xF0", Registers{A: 0xF0}, args{r: RegA}, flags{0, 0, 0, 0}},
		{"A|=B, A=0x01, B=0x03", Registers{A: 0x01, B: 0x03}, args{r: RegB}, flags{0, 0, 0, 0}},
		{"A|=C, A=0x00, C=0x00", Registers{A: 0x00, C: 0x00}, args{r: RegC}, flags{1, 0, 0, 0}},
		{"A|=F, A=0x10, F=0x01", Registers{A: 0x10, F: 0x01}, args{r: RegF}, flags{0, 0, 0, 0}},
		{"A|=H, A=0x01, H=0xFF", Registers{A: 0x01, H: 0xFF}, args{r: RegH}, flags{0, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.Or_r(tt.args.r)

			// Expect A=A|r
			getr, _ := c.getReg8(tt.args.r)
			expected := tt.regs.A | getr()
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Or_d8(t *testing.T) {
	type args struct {
		d8 byte
	}
	tests := []struct {
		name  string
		regs  Registers
		args  args
		flags flags
	}{
		{"A|=d8, A=0xF0, d8=0xF0", Registers{A: 0xF0}, args{d8: 0xF0}, flags{0, 0, 0, 0}},
		{"A|=d8, A=0x01, d8=0x03", Registers{A: 0x01}, args{d8: 0x03}, flags{0, 0, 0, 0}},
		{"A|=d8, A=0x00, d8=0x00", Registers{A: 0x00}, args{d8: 0x00}, flags{1, 0, 0, 0}},
		{"A|=d8, A=0x10, d8=0x01", Registers{A: 0x10}, args{d8: 0x01}, flags{0, 0, 0, 0}},
		{"A|=d8, A=0x01, d8=0xFF", Registers{A: 0x01}, args{d8: 0xFF}, flags{0, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.Or_d8(tt.args.d8)

			// Expect A=A|r
			expected := tt.regs.A | tt.args.d8
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Or_valHL(t *testing.T) {
	tests := []struct {
		name  string
		regs  Registers
		valHL byte
		flags flags
	}{
		{"A|=(HL), A=0xF0, (HL)=0xF0", Registers{A: 0xF0}, 0xF0, flags{0, 0, 0, 0}},
		{"A|=(HL), A=0x01, (HL)=0x03", Registers{A: 0x01}, 0x03, flags{0, 0, 0, 0}},
		{"A|=(HL), A=0x00, (HL)=0x00", Registers{A: 0x00}, 0x00, flags{1, 0, 0, 0}},
		{"A|=(HL), A=0x10, (HL)=0x01", Registers{A: 0x10}, 0x01, flags{0, 0, 0, 0}},
		{"A|=(HL), A=0x01, (HL)=0xFF", Registers{A: 0x01}, 0xFF, flags{0, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs
			// initialize memory
			err := c.mem.Wb(c.getHL(), tt.valHL)
			if err != nil {
				t.Error(err)
			}

			c.Or_valHL()

			// Expect A=A|(HL)
			expected := tt.regs.A | tt.valHL
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Cp_r(t *testing.T) {
	type args struct {
		r Reg8
	}
	tests := []struct {
		name  string
		regs  Registers
		args  args
		flags flags
	}{
		{"A=0x01, B=0x02", Registers{A: 0x01, B: 0x02}, args{r: RegB}, flags{0, 1, 1, 1}},
		{"A=0x00, B=0x00", Registers{A: 0x00, B: 0x00}, args{r: RegB}, flags{1, 1, 0, 0}},
		{"A=0x10, B=0x01", Registers{A: 0x10, B: 0x01}, args{r: RegB}, flags{0, 1, 1, 0}},
		{"A=0x01, B=0xFF", Registers{A: 0x01, B: 0xFF}, args{r: RegB}, flags{0, 1, 1, 1}},
		{"A=0xF0, B=0x10", Registers{A: 0xF0, B: 0x10}, args{r: RegB}, flags{0, 1, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.Cp_r(tt.args.r)

			// Expect A to be unchanged
			expected := tt.regs.A
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}

}

func TestCPU_Cp_d8(t *testing.T) {
	type args struct {
		d8 byte
	}
	tests := []struct {
		name  string
		regs  Registers
		args  args
		flags flags
	}{
		{"A=0x01, d8=0x02", Registers{A: 0x01}, args{d8: 0x02}, flags{0, 1, 1, 1}},
		{"A=0x00, d8=0x00", Registers{A: 0x00}, args{d8: 0x00}, flags{1, 1, 0, 0}},
		{"A=0x10, d8=0x01", Registers{A: 0x10}, args{d8: 0x01}, flags{0, 1, 1, 0}},
		{"A=0x01, d8=0xFF", Registers{A: 0x01}, args{d8: 0xFF}, flags{0, 1, 1, 1}},
		{"A=0xF0, d8=0x10", Registers{A: 0xF0}, args{d8: 0x10}, flags{0, 1, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.Cp_d8(tt.args.d8)

			// Expect A to be unchanged
			expected := tt.regs.A
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Cp_valHL(t *testing.T) {
	tests := []struct {
		name  string
		regs  Registers
		valHL byte
		flags flags
	}{
		{"A=0x01, (HL)=0x02", Registers{A: 0x01}, 0x02, flags{0, 1, 1, 1}},
		{"A=0x00, (HL)=0x00", Registers{A: 0x00}, 0x00, flags{1, 1, 0, 0}},
		{"A=0x10, (HL)=0x01", Registers{A: 0x10}, 0x01, flags{0, 1, 1, 0}},
		{"A=0x01, (HL)=0xFF", Registers{A: 0x01}, 0xFF, flags{0, 1, 1, 1}},
		{"A=0xF0, (HL)=0x10", Registers{A: 0xF0}, 0x10, flags{0, 1, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs
			// initialize memory
			err := c.mem.Wb(c.getHL(), tt.valHL)
			if err != nil {
				t.Error(err)
			}

			c.Cp_valHL()

			// Expect A to be unchanged
			expected := tt.regs.A
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Inc_r(t *testing.T) {
	type args struct {
		r Reg8
	}
	tests := []struct {
		name  string
		regs  Registers
		args  args
		flags flags
	}{
		{"A++, A=0x01", Registers{A: 0x01, F: 0x10}, args{r: RegA}, flags{0, 0, 0, 1}},
		{"B++, B=0x00", Registers{B: 0x00, F: 0x00}, args{r: RegB}, flags{0, 0, 0, 0}},
		{"C++, C=0x0F", Registers{C: 0x0F, F: 0x10}, args{r: RegC}, flags{0, 0, 1, 1}},
		{"D++, D=0xFF", Registers{D: 0xFF, F: 0x00}, args{r: RegD}, flags{1, 0, 1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs
			getR, _ := c.getReg8(tt.args.r)
			oldR := getR()

			c.Inc_r(tt.args.r)

			// Expect r = oldr + 1
			expected := oldR + 1
			if getR() != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, getR())
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Inc_valHL(t *testing.T) {
	tests := []struct {
		name  string
		regs  Registers
		valHL byte
		flags flags
	}{
		{"(HL)++, (HL)=0x01", Registers{F: 0x00}, 0x01, flags{0, 0, 0, 0}},
		{"(HL)++, (HL)=0x00", Registers{F: 0x10}, 0x00, flags{0, 0, 0, 1}},
		{"(HL)++, (HL)=0x0F", Registers{F: 0x00}, 0x0F, flags{0, 0, 1, 0}},
		{"(HL)++, (HL)=0xFF", Registers{F: 0x10}, 0xFF, flags{1, 0, 1, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs
			err := c.mem.Wb(c.getHL(), tt.valHL)
			if err != nil {
				t.Error(err)
			}

			c.Inc_valHL()

			// Expect (HL) = (oldHL) + 1
			expected := tt.valHL + 1
			valHL, err := c.mem.Rb(c.getHL())
			if err != nil {
				t.Error(err)
			}
			if valHL != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, valHL)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Dec_r(t *testing.T) {
	type args struct {
		r Reg8
	}
	tests := []struct {
		name  string
		regs  Registers
		args  args
		flags flags
	}{
		{"A--, A=0x01", Registers{A: 0x01, F: 0x00}, args{r: RegA}, flags{1, 1, 0, 0}},
		{"B--, B=0x00", Registers{B: 0x00, F: 0x10}, args{r: RegB}, flags{0, 1, 1, 1}},
		{"C--, C=0x10", Registers{C: 0x10, F: 0x00}, args{r: RegC}, flags{0, 1, 1, 0}},
		{"D--, D=0xFF", Registers{D: 0xFF, F: 0x10}, args{r: RegD}, flags{0, 1, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs
			getR, _ := c.getReg8(tt.args.r)
			oldR := getR()

			c.Dec_r(tt.args.r)

			// Expect r = oldr - 1
			expected := oldR - 1
			if getR() != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, getR())
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Dec_valHL(t *testing.T) {
	tests := []struct {
		name  string
		regs  Registers
		valHL byte
		flags flags
	}{
		{"(HL)--, (HL)=0x01", Registers{F: 0x00}, 0x01, flags{1, 1, 0, 0}},
		{"(HL)--, (HL)=0x00", Registers{F: 0x10}, 0x00, flags{0, 1, 1, 1}},
		{"(HL)--, (HL)=0x10", Registers{F: 0x00}, 0x10, flags{0, 1, 1, 0}},
		{"(HL)--, (HL)=0xFF", Registers{F: 0x10}, 0xFF, flags{0, 1, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs
			err := c.mem.Wb(c.getHL(), tt.valHL)
			if err != nil {
				t.Error(err)
			}

			c.Dec_valHL()

			// Expect (HL) = (oldHL) - 1
			expected := tt.valHL - 1
			valHL, err := c.mem.Rb(c.getHL())
			if err != nil {
				t.Error(err)
			}
			if valHL != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, valHL)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Daa(t *testing.T) {
	tests := []struct {
		name  string
		regs  Registers
		Aout  byte
		flags flags
	}{
		// DAA:
		// if last operation was addition,
		// 	add 0x60 to A if carry or if A > 0x99
		//  add 0x06 to A if half carry or if (A & 0x0F) > 0x09
		// if last operation was subtraction,
		// 	subtr 0x60 from A if carry
		//  subtr 0x06 from A if half carry
		//
		// Z = 1 if A is 0 else 0
		// H = 0

		{"addition, no h, no c, top not oob, bottom not oob", Registers{A: 0x34, F: 0x00}, 0x34, flags{0, 0, 0, 0}},
		{"addition, no h, no c, top not oob, bottom YES oob", Registers{A: 0x3F, F: 0x00}, 0x45, flags{0, 0, 0, 0}},
		{"addition, no h, no c, top YES oob, bottom not oob", Registers{A: 0xA4, F: 0x00}, 0x04, flags{0, 0, 0, 1}},
		{"addition, no h, no c, top YES oob, bottom YES oob", Registers{A: 0xAA, F: 0x00}, 0x10, flags{0, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.Daa()

			// Expect A=Aout
			expected := tt.Aout
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

func TestCPU_Cpl(t *testing.T) {
	tests := []struct {
		name  string
		regs  Registers
		flags flags
	}{
		{"A=0xF0", Registers{A: 0xF0, F: 0xF0}, flags{1, 1, 1, 1}},
		{"A=0x00", Registers{A: 0x00, F: 0x80}, flags{1, 1, 1, 0}},
		{"A=0xAB", Registers{A: 0xAB, F: 0x20}, flags{0, 1, 1, 0}},
		{"A=0x0E", Registers{A: 0x0E, F: 0x00}, flags{0, 1, 1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.Cpl()

			// Expect A=A^0xFF
			expected := tt.regs.A ^ 0xFF
			if c.A != expected {
				t.Errorf("Expected A to be %02x, got %02x", expected, c.A)
			}

			// Expect flags to be set
			checkFlags(c, tt.flags, t)
		})
	}
}

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
