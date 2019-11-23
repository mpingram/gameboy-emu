package cpu

import (
	"testing"
)

func TestCPU_Ld_r1_r2(t *testing.T) {
	type args struct {
		r1 Reg8
		r2 Reg8
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"A <- A, A=0x01", Registers{A: 0x01}, args{r1: RegA, r2: RegA}},
		{"A <- B, A=0x01, B=0x02", Registers{A: 0x01, B: 0x02}, args{r1: RegA, r2: RegB}},
		{"C <- D, C=0x00, D=0xDD", Registers{C: 0x00, D: 0xDD}, args{r1: RegC, r2: RegD}},
		{"D <- H, D=0xDD, H=0x00", Registers{D: 0xDD, H: 0x00}, args{r1: RegD, r2: RegH}},
		{"L <- H, L=0x08, H=0x00", Registers{L: 0x08, H: 0x00}, args{r1: RegL, r2: RegH}},

		{"F <- B, F=0xF0, B=0xFF", Registers{F: 0xF0, B: 0xFF}, args{r1: RegF, r2: RegB}},
		{"A <- F, A=0x01, F=0xFF", Registers{A: 0x01, F: 0xFF}, args{r1: RegA, r2: RegF}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.Ld_r1_r2(tt.args.r1, tt.args.r2)

			// expect target register to be value of src register
			getr1, _ := c.getReg8(tt.args.r1)
			getr2, _ := c.getReg8(tt.args.r2)
			if getr1() != getr2() {
				t.Errorf("Expected r1 to be %02x, got %02x", getr1(), getr2())
			}
		})
	}
}

func TestCPU_Ld_r_d8(t *testing.T) {
	type args struct {
		r1 Reg8
		d8 byte
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"A <- 0xBF, A=0x01", Registers{A: 0x01}, args{r1: RegA, d8: 0xBF}},
		{"C <- 0x00, C=0x00", Registers{C: 0x00}, args{r1: RegC, d8: 0x00}},
		{"L <- 0xFF, L=0x08", Registers{L: 0x08}, args{r1: RegL, d8: 0xFF}},
		{"F <- 0xFF, F=0x70", Registers{F: 0x70}, args{r1: RegF, d8: 0xFF}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.Ld_r_d8(tt.args.r1, tt.args.d8)

			// expect target register to be value of src register
			getr1, _ := c.getReg8(tt.args.r1)
			if getr1() != tt.args.d8 {
				t.Errorf("Expected r1 to be %02x, got %02x", tt.args.d8, getr1())
			}
		})
	}
}

func TestCPU_Ld_r_valHL(t *testing.T) {
	type args struct {
		r1 Reg8
	}
	tests := []struct {
		name  string
		regs  Registers
		valHL byte
		args  args
	}{
		{"A <-(HL), A=0x01, H=0x12, L=0x34, (HL)=0x10", Registers{A: 0x01, H: 0x12, L: 0x34}, 0x10, args{r1: RegA}},
		{"C <-(HL), C=0x00, H=0x17, L=0x38, (HL)=0x21", Registers{C: 0x00, H: 0x17, L: 0x38}, 0x21, args{r1: RegC}},
		{"L <-(HL), L=0x08, H=0x00, (HL)=0x32", Registers{L: 0x08, H: 0x00}, 0x32, args{r1: RegL}},
		{"H <-(HL), H=0x0A, L=0x00, (HL)=0x43", Registers{H: 0x0A, L: 0x00}, 0x43, args{r1: RegH}},
		{"F <-(HL), F=0x70, H=0xD0, L=0x0B, (HL)=0x54", Registers{F: 0x70, H: 0xD0, L: 0x0B}, 0x54, args{r1: RegF}},
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

			c.Ld_r_valHL(tt.args.r1)

			// expect target register to be value of HL
			getr1, _ := c.getReg8(tt.args.r1)
			if getr1() != tt.valHL {
				t.Errorf("Expected r1 to be %02x, got %02x", tt.valHL, getr1())
			}
		})
	}
}

func TestCPU_Ld_valHL_r(t *testing.T) {
	type args struct {
		r Reg8
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"(HL) <- A, A=0x01, H=0x12, L=0x34, (HL)=0x10", Registers{A: 0x01, H: 0x12, L: 0x34}, args{r: RegA}},
		{"(HL) <- C, C=0x00, H=0x17, L=0x38, (HL)=0x21", Registers{C: 0x00, H: 0x17, L: 0x38}, args{r: RegC}},
		{"(HL) <- L, L=0x08, H=0x00, (HL)=0x32", Registers{L: 0x08, H: 0x00}, args{r: RegL}},
		{"(HL) <- H, H=0x0A, L=0x00, (HL)=0x43", Registers{H: 0x0A, L: 0x00}, args{r: RegH}},
		{"(HL) <- F, F=0x70, H=0xD0, L=0x0B, (HL)=0x54", Registers{F: 0x70, H: 0xD0, L: 0x0B}, args{r: RegF}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.Ld_valHL_r(tt.args.r)

			// expect memory at HL to be value of src register
			getr, _ := c.getReg8(tt.args.r)
			valHL, err := c.mem.Rb(c.getHL())
			if err != nil {
				t.Error(err)
			}
			if getr() != valHL {
				t.Errorf("Expected r1 to be %02x, got %02x", valHL, getr())
			}
		})
	}
}

func TestCPU_Ld_valHL_d8(t *testing.T) {
	type args struct {
		d8 byte
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"(HL) <- 0x00, H=0x12, L=0x34", Registers{H: 0x12, L: 0x34}, args{d8: 0x00}},
		{"(HL) <- 0xB0, H=0x17, L=0x38", Registers{H: 0x17, L: 0x38}, args{d8: 0xB0}},
		{"(HL) <- 0x0F, H=0x00, L=0x08", Registers{H: 0x00, L: 0x08}, args{d8: 0x0F}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.Ld_valHL_d8(tt.args.d8)

			// Expect (HL) to be d8
			valHL, err := c.mem.Rb(c.getHL())
			if err != nil {
				t.Error(err)
			}
			if valHL != tt.args.d8 {
				t.Errorf("Expected (HL) to be %02x, got %02x", tt.args.d8, valHL)
			}
		})
	}
}

func TestCPU_Ld_A_valBC(t *testing.T) {
	tests := []struct {
		name  string
		regs  Registers
		valBC byte
	}{
		{"A <-(BC), B=0x12, C=0x34, (BC)=0x00", Registers{B: 0x12, C: 0x34}, 0x00},
		{"A <-(BC), B=0x17, C=0x38, (BC)=0x80", Registers{B: 0x17, C: 0x38}, 0x80},
		{"A <-(BC), B=0x00, C=0x08, (BC)=0x01", Registers{B: 0x00, C: 0x08}, 0x01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs
			// initialize memory
			err := c.mem.Wb(c.getBC(), tt.valBC)
			if err != nil {
				t.Error(err)
			}

			c.Ld_A_valBC()

			// Expect A to be (BC)
			if c.A != tt.valBC {
				t.Errorf("Expected A to be %02x, got %02x", tt.valBC, c.A)
			}
		})
	}
}

func TestCPU_Ld_A_valDE(t *testing.T) {
	tests := []struct {
		name  string
		regs  Registers
		valDE byte
	}{
		{"A <-(DE), D=0x12, E=0x34, (DE)=0x00", Registers{D: 0x12, E: 0x34}, 0x00},
		{"A <-(DE), D=0x17, E=0x38, (DE)=0x80", Registers{D: 0x17, E: 0x38}, 0x80},
		{"A <-(DE), D=0x00, E=0x08, (DE)=0x01", Registers{D: 0x00, E: 0x08}, 0x01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs
			// initialize memory
			err := c.mem.Wb(c.getDE(), tt.valDE)
			if err != nil {
				t.Error(err)
			}

			c.Ld_A_valDE()

			// Expect A to be (DE)
			if c.A != tt.valDE {
				t.Errorf("Expected A to be %02x, got %02x", tt.valDE, c.A)
			}
		})
	}
}

func TestCPU_Ld_A_valA16(t *testing.T) {
	type args struct {
		a16 uint16
	}
	tests := []struct {
		name   string
		regs   Registers
		args   args
		valA16 byte
	}{
		{"A <-(a16), a16=0x1234, (a16)=0x00", Registers{}, args{a16: 0x1234}, 0x00},
		{"A <-(a16), a16=0x1738, (a16)=0x80", Registers{}, args{a16: 0x1738}, 0x80},
		{"A <-(a16), a16=0x0008, (a16)=0x01", Registers{}, args{a16: 0x0008}, 0x01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs
			// initialize memory
			err := c.mem.Wb(tt.args.a16, tt.valA16)
			if err != nil {
				t.Error(err)
			}

			c.Ld_A_valA16(tt.args.a16)

			// Expect A to be (a16)
			if c.A != tt.valA16 {
				t.Errorf("Expected A to be %02x, got %02x", tt.valA16, c.A)
			}
		})
	}
}

func TestCPU_Ld_valBC_A(t *testing.T) {
	tests := []struct {
		name string
		regs Registers
	}{
		{"(BC)<- A, B=0x12, C=0x34, A=0x00", Registers{B: 0x12, C: 0x34, A: 0x00}},
		{"(BC)<- A, B=0x17, C=0x38, A=0x80", Registers{B: 0x17, C: 0x38, A: 0x80}},
		{"(BC)<- A, B=0x00, C=0x08, A=0x01", Registers{B: 0x00, C: 0x08, A: 0x01}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.Ld_valBC_A()

			// Expect (BC) to be A
			valBC, err := c.mem.Rb(c.getBC())
			if err != nil {
				t.Error(err)
			}
			if valBC != tt.regs.A {
				t.Errorf("Expected (BC) to be %02x, got %02x", tt.regs.A, valBC)
			}
		})
	}
}

func TestCPU_Ld_valDE_A(t *testing.T) {
	tests := []struct {
		name string
		regs Registers
	}{
		{"(DE)<- A, D=0x12, E=0x34, A=0x00", Registers{D: 0x12, E: 0x34, A: 0x00}},
		{"(DE)<- A, D=0x17, E=0x38, A=0x80", Registers{D: 0x17, E: 0x38, A: 0x80}},
		{"(DE)<- A, D=0x00, E=0x08, A=0x01", Registers{D: 0x00, E: 0x08, A: 0x01}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.Ld_valDE_A()

			// Expect (DE) to be A
			valDE, err := c.mem.Rb(c.getDE())
			if err != nil {
				t.Error(err)
			}
			if valDE != tt.regs.A {
				t.Errorf("Expected (BC) to be %02x, got %02x", tt.regs.A, valDE)
			}
		})
	}
}

func TestCPU_Ld_valA16_A(t *testing.T) {
	type args struct {
		a16 uint16
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"(a16)<- A, a16=0x1234, A=0x00", Registers{A: 0x00}, args{a16: 0x1234}},
		{"(a16)<- A, a16=0x1738, A=0x80", Registers{A: 0x80}, args{a16: 0x1738}},
		{"(a16)<- A, a16=0x0801, A=0x01", Registers{A: 0x01}, args{a16: 0x0801}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.Ld_valA16_A(tt.args.a16)

			// Expect (a16) to be A
			valA16, err := c.mem.Rb(tt.args.a16)
			if err != nil {
				t.Error(err)
			}
			if valA16 != tt.regs.A {
				t.Errorf("Expected (BC) to be %02x, got %02x", tt.regs.A, valA16)
			}
		})
	}
}

func TestCPU_Ld_A_FF00_plus_a8(t *testing.T) {
	type args struct {
		a8 byte
	}
	tests := []struct {
		name          string
		regs          Registers
		args          args
		valFF00plusA8 byte
	}{
		{"A <-($FF00+a8), a8=0x00", Registers{}, args{a8: 0x00}, 0x12},
		{"A <-($FF00+a8), a8=0x80", Registers{}, args{a8: 0x80}, 0x00},
		{"A <-($FF00+a8), a8=0xFF", Registers{}, args{a8: 0xFF}, 0xF0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs
			// initialize memory
			addr := 0xFF00 + uint16(tt.args.a8)
			err := c.mem.Wb(addr, tt.valFF00plusA8)
			if err != nil {
				t.Error(err)
			}

			c.Ld_A_FF00_plus_a8(tt.args.a8)

			// Expect A to be ($FF00 + a8)
			if c.A != tt.valFF00plusA8 {
				t.Errorf("Expected A to be %02x, got %02x", tt.valFF00plusA8, c.A)
			}
		})
	}
}

func TestCPU_Ld_FF00_plus_a8_A(t *testing.T) {
	type args struct {
		a8 byte
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"($FF00+a8) <- A, a8=0x00, A=0x12", Registers{A: 0x12}, args{0x00}},
		{"($FF00+a8) <- A, a8=0x80, A=0x00", Registers{A: 0x00}, args{0x80}},
		{"($FF00+a8) <- A, a8=0xFF, A=0xF0", Registers{A: 0xF0}, args{0xFF}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.Ld_FF00_plus_a8_A(tt.args.a8)

			// Expect ($FF00 + a8) to be A
			addr := 0xFF00 + uint16(tt.args.a8)
			valFF00plusA8, err := c.mem.Rb(addr)
			if err != nil {
				t.Error(err)
			}
			if valFF00plusA8 != tt.regs.A {
				t.Errorf("Expected ($FF00+a8) to be %02x, got %02x", tt.regs.A, valFF00plusA8)
			}
		})
	}
}

func TestCPU_Ld_A_FF00_plus_C(t *testing.T) {
	tests := []struct {
		name         string
		regs         Registers
		valFF00plusC byte
	}{
		{"A <-($FF00+C), C=0x00", Registers{C: 0x00}, 0x12},
		{"A <-($FF00+C), C=0x80", Registers{C: 0x80}, 0x00},
		{"A <-($FF00+C), C=0xFF", Registers{C: 0xFF}, 0xF0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs
			// initialize memory
			addr := 0xFF00 + uint16(tt.regs.C)
			err := c.mem.Wb(addr, tt.valFF00plusC)
			if err != nil {
				t.Error(err)
			}

			c.Ld_A_FF00_plus_C()

			// Expect A to be ($FF00 + a8)
			if c.A != tt.valFF00plusC {
				t.Errorf("Expected A to be %02x, got %02x", tt.valFF00plusC, c.A)
			}
		})
	}
}

func TestCPU_Ld_FF00_plus_C_A(t *testing.T) {
	tests := []struct {
		name string
		regs Registers
	}{
		{"($FF00+C) <- A, C=0x00, A=0x12", Registers{A: 0x12, C: 0x00}},
		{"($FF00+C) <- A, C=0x80, A=0x00", Registers{A: 0x00, C: 0x80}},
		{"($FF00+C) <- A, C=0xFF, A=0xF0", Registers{A: 0xF0, C: 0xFF}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.Ld_FF00_plus_C_A()

			// Expect ($FF00 + C) to be A
			addr := 0xFF00 + uint16(tt.regs.C)
			valFF00plusC, err := c.mem.Rb(addr)
			if err != nil {
				t.Error(err)
			}
			if valFF00plusC != tt.regs.A {
				t.Errorf("Expected ($FF00+C) to be %02x, got %02x", tt.regs.A, valFF00plusC)
			}
		})
	}
}

func TestCPU_Ld_valHLinc_A(t *testing.T) {
	tests := []struct {
		name string
		regs Registers
	}{
		{"(HL)<- A, HL=HL+1, H=0x12, L=0x01, A=0x00", Registers{H: 0x12, L: 0x01, A: 0x00}},
		{"(HL)<- A, HL=HL+1, H=0x01, L=0xFF, A=0xF0", Registers{H: 0x00, L: 0xFF, A: 0xF0}},
		{"(HL)<- A, HL=HL+1, H=0xFE, L=0xFF, A=0x12", Registers{H: 0xFF, L: 0xFF, A: 0x12}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.Ld_valHLinc_A()

			// Expect (oldHL) to be A
			oldHL := uint16(tt.regs.H)<<8 | uint16(tt.regs.L)
			valHL, err := c.mem.Rb(oldHL)
			if err != nil {
				t.Error(err)
			}
			if valHL != tt.regs.A {
				t.Errorf("Expected (HL) to be %02x, got %02x", tt.regs.A, valHL)
			}

			// Expect HL to be oldHL + 1
			if c.getHL() != oldHL+1 {
				t.Errorf("Expected HL to be %04x, got %04x", oldHL+1, c.getHL())
			}
		})
	}
}

func TestCPU_Ld_A_valHLinc(t *testing.T) {
	tests := []struct {
		name  string
		regs  Registers
		valHL byte
	}{
		{"A <-(HL), HL=HL+1, H=0x12, L=0x01, (HL)=0x00", Registers{H: 0x12, L: 0x01}, 0x00},
		{"A <-(HL), HL=HL+1, H=0x01, L=0xFF, (HL)=0xF0", Registers{H: 0x00, L: 0xFF}, 0xF0},
		{"A <-(HL), HL=HL+1, H=0xFE, L=0xFF, (HL)=0x12", Registers{H: 0xFF, L: 0xFF}, 0x12},
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

			c.Ld_A_valHLinc()

			// Expect A to be (oldHL)
			oldHL := uint16(tt.regs.H)<<8 | uint16(tt.regs.L)
			if c.A != tt.valHL {
				t.Errorf("Expected A to be %02x, got %02x", tt.valHL, c.A)
			}
			// Expect HL to be incremented by 1
			if c.getHL() != oldHL+1 {
				t.Errorf("Expected HL to be %04x, got %04x", oldHL+1, c.getHL())
			}
		})
	}
}

func TestCPU_Ld_valHLdec_A(t *testing.T) {
	tests := []struct {
		name string
		regs Registers
	}{
		{"(HL)<- A, HL=HL-1, H=0x12, L=0x00, A=0x00", Registers{H: 0x12, L: 0x00, A: 0x00}},
		{"(HL)<- A, HL=HL-1, H=0x01, L=0xFF, A=0xF0", Registers{H: 0x00, L: 0xFF, A: 0xF0}},
		{"(HL)<- A, HL=HL-1, H=0xFE, L=0xFF, A=0x12", Registers{H: 0xFE, L: 0xFF, A: 0x12}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testSetup()

			c.Registers = tt.regs

			c.Ld_valHLdec_A()

			// Expect (oldHL) to be A
			oldHL := uint16(tt.regs.H)<<8 | uint16(tt.regs.L)
			valHL, err := c.mem.Rb(oldHL)
			if err != nil {
				t.Error(err)
			}
			if valHL != tt.regs.A {
				t.Errorf("Expected (HL) to be %02x, got %02x", tt.regs.A, valHL)
			}

			// Expect HL to be oldHL - 1
			if c.getHL() != oldHL-1 {
				t.Errorf("Expected HL to be %04x, got %04x", oldHL-1, c.getHL())
			}
		})
	}
}

func TestCPU_Ld_A_valHLdec(t *testing.T) {
	tests := []struct {
		name  string
		regs  Registers
		valHL byte
	}{
		{"A <-(HL), HL=HL-1, H=0x12, L=0x00, (HL)=0x00", Registers{H: 0x12, L: 0x00}, 0x00},
		{"A <-(HL), HL=HL-1, H=0x01, L=0xFF, (HL)=0xF0", Registers{H: 0x00, L: 0xFF}, 0xF0},
		{"A <-(HL), HL=HL-1, H=0xFE, L=0xFF, (HL)=0x12", Registers{H: 0xFF, L: 0xFF}, 0x12},
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

			c.Ld_A_valHLdec()

			// Expect A to be (oldHL)
			oldHL := uint16(tt.regs.H)<<8 | uint16(tt.regs.L)
			if c.A != tt.valHL {
				t.Errorf("Expected A to be %02x, got %02x", tt.valHL, c.A)
			}
			// Expect HL to be oldHL - 1
			if c.getHL() != oldHL-1 {
				t.Errorf("Expected HL to be %04x, got %04x", oldHL-1, c.getHL())
			}
		})
	}
}
