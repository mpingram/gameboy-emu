package cpu

import (
	"github.com/mpingram/gameboy-emu/mmu"
	"testing"
)

func TestCPU_Jp(t *testing.T) {

	type args struct {
		a16 uint16
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"Jp to should move PC to a16", Registers{}, args{0x9BFF}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			c.Jp(tt.args.a16)
			// Expect PC to be a16
			if c.PC != tt.args.a16 {
				t.Errorf("Expected PC %04x, got %04x", tt.args.a16, c.PC)
			}
			// Expect SP to not have moved
			if c.SP != tt.regs.SP {
				t.Errorf("Expected SP to be unchanged from %04x; got %04x", tt.regs.SP, c.SP)
			}
		})
	}
}

func TestCPU_Jp_HL(t *testing.T) {
	tests := []struct {
		name string
		regs Registers
	}{
		{"HL = 0xF00F", Registers{H: 0xF0, L: 0x0F}},
		{"HL = 0xB00B", Registers{H: 0xB0, L: 0x0B}},
		{"HL = 0x1234", Registers{H: 0x12, L: 0x34}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs
			c.Jp_HL()
			// Expect PC to be HL
			if c.PC != c.getHL() {
				t.Errorf("Want PC: %04x, got %04x", c.getHL(), c.PC)
			}
		})
	}
}

func TestCPU_JpNZ(t *testing.T) {
	type args struct {
		a16 uint16
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"Z = 0, a16 = 0xD00B", Registers{F: 0x00}, args{a16: 0xD00B}},
		{"Z = 1, a16 = 0xD00B", Registers{F: 0xF0}, args{a16: 0xD00B}},
		{"Z = 1, a16 = 0x0000", Registers{F: 0x00}, args{a16: 0x0000}},
		{"Z = 0, a16 = 0x0000", Registers{F: 0xF0}, args{a16: 0x0000}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs
			c.JpNZ(tt.args.a16)
			if c.getFlagZ() == false {
				if c.PC != tt.args.a16 {
					t.Errorf("Want PC: %04x, got %04x", tt.args.a16, c.PC)
				}
			}
		})
	}
}

func TestCPU_JpZ(t *testing.T) {
	type args struct {
		a16 uint16
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"Z = 0, a16 = 0xD00B", Registers{F: 0x00}, args{a16: 0xD00B}},
		{"Z = 1, a16 = 0xD00B", Registers{F: 0xF0}, args{a16: 0xD00B}},
		{"Z = 1, a16 = 0x0000", Registers{F: 0x00}, args{a16: 0x0000}},
		{"Z = 0, a16 = 0x0000", Registers{F: 0xF0}, args{a16: 0x0000}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			c.JpZ(tt.args.a16)
			if c.getFlagZ() == true {
				if c.PC != tt.args.a16 {
					t.Errorf("Want PC: %04x, got %04x", tt.args.a16, c.PC)
				}
			}
		})
	}
}

func TestCPU_JpNC(t *testing.T) {
	type args struct {
		a16 uint16
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"Z = 0, a16 = 0xD00B", Registers{F: 0x00}, args{a16: 0xD00B}},
		{"Z = 1, a16 = 0xD00B", Registers{F: 0xF0}, args{a16: 0xD00B}},
		{"Z = 1, a16 = 0x0000", Registers{F: 0xF0}, args{a16: 0x0000}},
		{"Z = 0, a16 = 0x0000", Registers{F: 0x00}, args{a16: 0x0000}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			c.JpNC(tt.args.a16)
			if c.getFlagC() == false {
				if c.PC != tt.args.a16 {
					t.Errorf("Want PC: %04x, got %04x", tt.args.a16, c.PC)
				}
			}
		})
	}
}

func TestCPU_JpC(t *testing.T) {
	type args struct {
		a16 uint16
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"C = 0, a16 = 0xD00B", Registers{F: 0x00}, args{a16: 0xD00B}},
		{"C = 1, a16 = 0xD00B", Registers{F: 0xF0}, args{a16: 0xD00B}},
		{"C = 1, a16 = 0x0000", Registers{F: 0xF0}, args{a16: 0x0000}},
		{"C = 0, a16 = 0x0000", Registers{F: 0x00}, args{a16: 0x0000}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			c.JpC(tt.args.a16)
			if c.getFlagC() == true {
				if c.PC != tt.args.a16 {
					t.Errorf("Want PC: %04x, got %04x", tt.args.a16, c.PC)
				}
			} else {
				if c.PC != tt.regs.PC {
					t.Errorf("Expected PC to be unchanged from %04x, got %04x", tt.regs.PC, c.PC)
				}
			}
		})
	}
}

func TestCPU_Jr(t *testing.T) {
	type args struct {
		r8 int8
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"r8=0, PC=0", Registers{PC: 0}, args{0}},
		{"r8=0x7F, PC=0", Registers{PC: 0}, args{0x7F}},
		{"r8=-0x80, PC=0x80", Registers{PC: 128}, args{-0x80}},
		{"r8=0x7F, PC=0xFFFE", Registers{PC: 0xFFFE}, args{0x7F}},
		{"r8=-0x80, PC=1", Registers{PC: 1}, args{-0x80}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			c.Jr(tt.args.r8)
			// Expect PC to have changed by r8
			var expected uint16
			if tt.args.r8 < 0 {
				expected = tt.regs.PC - uint16(-1*tt.args.r8)
			} else {
				expected = tt.regs.PC + uint16(tt.args.r8)
			}
			/* FIXME what is the expected wraparound behavior?
			if expected < 0 {
				expected = 0
			} else if expected > 0xFFFF {
				expected = 0xFFFF
			}
			*/
			if c.PC != expected {
				t.Errorf("Want PC: %04x, got %04x", expected, c.PC)
			}
		})
	}
}

func TestCPU_JrNZ(t *testing.T) {
	type args struct {
		r8 int8
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"r8=0, PC=0, Z=0", Registers{PC: 0}, args{0}},
		{"r8=0, PC=0, Z=1", Registers{PC: 0, F: 0x80}, args{0}},
		{"r8=0x7F, PC=0, Z=0", Registers{PC: 0}, args{0x7F}},
		{"r8=0x7F, PC=0, Z=1", Registers{PC: 0, F: 0x80}, args{0x7F}},
		{"r8=-0x80, PC=0x80, Z=0", Registers{PC: 128}, args{-0x80}},
		{"r8=-0x80, PC=0x80, Z=1", Registers{PC: 128, F: 0x80}, args{-0x80}},
		{"r8=0x7F, PC=0xFFFE, Z=0", Registers{PC: 0xFFFE}, args{0x7F}},
		{"r8=0x7F, PC=0xFFFE, Z=1", Registers{PC: 0xFFFE, F: 0x80}, args{0x7F}},
		{"r8=-0x80, PC=1, Z=0", Registers{PC: 1}, args{-0x80}},
		{"r8=-0x80, PC=1, Z=1", Registers{PC: 1, F: 0x80}, args{-0x80}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			c.JrNZ(tt.args.r8)

			// Expect PC to have changed by r8 IF Flags & Zbit = 0
			var expected uint16
			if tt.regs.F&0x80 != 0 {
				// If Z bit is 1, expect PC to not change
				expected = tt.regs.PC
			} else {
				// Otherwise expect jump
				if tt.args.r8 < 0 {
					expected = tt.regs.PC - uint16(-1*tt.args.r8)
				} else {
					expected = tt.regs.PC + uint16(tt.args.r8)
				}
			}
			if c.PC != expected {
				t.Errorf("Want PC: %04x, got %04x", expected, c.PC)
			}
		})
	}
}

func TestCPU_JrZ(t *testing.T) {
	type args struct {
		r8 int8
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"r8=0, PC=0, Z=0", Registers{PC: 0}, args{0}},
		{"r8=0, PC=0, Z=1", Registers{PC: 0, F: 0x80}, args{0}},
		{"r8=0x7F, PC=0, Z=0", Registers{PC: 0}, args{0x7F}},
		{"r8=0x7F, PC=0, Z=1", Registers{PC: 0, F: 0x80}, args{0x7F}},
		{"r8=-0x80, PC=0x80, Z=0", Registers{PC: 128}, args{-0x80}},
		{"r8=-0x80, PC=0x80, Z=1", Registers{PC: 128, F: 0x80}, args{-0x80}},
		{"r8=0x7F, PC=0xFFFE, Z=0", Registers{PC: 0xFFFE}, args{0x7F}},
		{"r8=0x7F, PC=0xFFFE, Z=1", Registers{PC: 0xFFFE, F: 0x80}, args{0x7F}},
		{"r8=-0x80, PC=1, Z=0", Registers{PC: 1}, args{-0x80}},
		{"r8=-0x80, PC=1, Z=1", Registers{PC: 1, F: 0x80}, args{-0x80}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			c.JrZ(tt.args.r8)

			// Expect PC to have changed by r8 IF Flags & Zbit = 1
			var expected uint16
			if tt.regs.F&0x80 == 0 {
				// If Z is 0, expect PC to not change
				expected = tt.regs.PC
			} else {
				// Otherwise expect jump
				if tt.args.r8 < 0 {
					expected = tt.regs.PC - uint16(-1*tt.args.r8)
				} else {
					expected = tt.regs.PC + uint16(tt.args.r8)
				}
			}
			if c.PC != expected {
				t.Errorf("Want PC: %04x, got %04x", expected, c.PC)
			}
		})
	}
}

func TestCPU_JrNC(t *testing.T) {
	type args struct {
		r8 int8
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"r8=0, PC=0, C=0", Registers{PC: 0}, args{0}},
		{"r8=0, PC=0, C=1", Registers{PC: 0, C: 0x10}, args{0}},
		{"r8=0x7F, PC=0, C=0", Registers{PC: 0}, args{0x7F}},
		{"r8=0x7F, PC=0, C=1", Registers{PC: 0, C: 0x10}, args{0x7F}},
		{"r8=-0x80, PC=0x80, C=0", Registers{PC: 128}, args{-0x80}},
		{"r8=-0x80, PC=0x80, C=1", Registers{PC: 128, C: 0x10}, args{-0x80}},
		{"r8=0x7F, PC=0xFFFE, C=0", Registers{PC: 0xFFFE}, args{0x7F}},
		{"r8=0x7F, PC=0xFFFE, C=1", Registers{PC: 0xFFFE, C: 0x10}, args{0x7F}},
		{"r8=-0x80, PC=1, C=0", Registers{PC: 1}, args{-0x80}},
		{"r8=-0x80, PC=1, C=1", Registers{PC: 1, C: 0x10}, args{-0x80}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			c.JrNC(tt.args.r8)

			// Expect PC to have changed by r8 IF Flags & Cbit = 0
			var expected uint16
			if tt.regs.F&0x10 != 0 {
				// If C is 1, expect PC to not change
				expected = tt.regs.PC
			} else {
				// Otherwise expect jump
				if tt.args.r8 < 0 {
					expected = tt.regs.PC - uint16(-1*tt.args.r8)
				} else {
					expected = tt.regs.PC + uint16(tt.args.r8)
				}
			}
			if c.PC != expected {
				t.Errorf("Want PC: %04x, got %04x", expected, c.PC)
			}
		})
	}
}

func TestCPU_JrC(t *testing.T) {
	type args struct {
		r8 int8
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"r8=0, PC=0, C=0", Registers{PC: 0}, args{0}},
		{"r8=0, PC=0, C=1", Registers{PC: 0, C: 0x10}, args{0}},
		{"r8=0x7F, PC=0, C=0", Registers{PC: 0}, args{0x7F}},
		{"r8=0x7F, PC=0, C=1", Registers{PC: 0, C: 0x10}, args{0x7F}},
		{"r8=-0x80, PC=0x80, C=0", Registers{PC: 128}, args{-0x80}},
		{"r8=-0x80, PC=0x80, C=1", Registers{PC: 128, C: 0x10}, args{-0x80}},
		{"r8=0x7F, PC=0xFFFE, C=0", Registers{PC: 0xFFFE}, args{0x7F}},
		{"r8=0x7F, PC=0xFFFE, C=1", Registers{PC: 0xFFFE, C: 0x10}, args{0x7F}},
		{"r8=-0x80, PC=1, C=0", Registers{PC: 1}, args{-0x80}},
		{"r8=-0x80, PC=1, C=1", Registers{PC: 1, C: 0x10}, args{-0x80}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			c.JrC(tt.args.r8)

			// Expect PC to have changed by r8 IF Flags & Cbit = 1
			var expected uint16
			if tt.regs.F&0x10 == 0 {
				// If C is 0, expect PC to not change
				expected = tt.regs.PC
			} else {
				// Otherwise expect jump
				if tt.args.r8 < 0 {
					expected = tt.regs.PC - uint16(-1*tt.args.r8)
				} else {
					expected = tt.regs.PC + uint16(tt.args.r8)
				}
			}
			if c.PC != expected {
				t.Errorf("Want PC: %04x, got %04x", expected, c.PC)
			}
		})
	}
}

func TestCPU_Call(t *testing.T) {
	type args struct {
		a16 uint16
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"a16=0xF00F, SP=0xCFFF, PC=0x0000", Registers{PC: 0x0000, SP: 0xCFFF}, args{a16: 0xF00F}},
		{"a16=0x0000, SP=0xCFFF, PC=0x1234", Registers{PC: 0x1234, SP: 0xCFFF}, args{a16: 0xF00F}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			c.Call(tt.args.a16)

			// Expect PC to be a16
			if c.PC != tt.args.a16 {
				t.Errorf("Wanted PC to be %04x, got %04x", tt.args.a16, c.PC)
			}
			// Expect SP to be SP-2
			if c.SP != tt.regs.SP-2 {
				t.Errorf("Wanted SP to be %04x, got %04x", tt.regs.SP-2, c.SP)
			}
			// Expect word at SP to contain previous PC
			valSP, err := c.mem.Rw(c.SP)
			if err != nil {
				t.Error(err)
			}
			if valSP != tt.regs.PC {
				t.Errorf("Wanted (SP) to be %04x, got %04x", tt.regs.PC, valSP)
			}
		})
	}
}

func TestCPU_CallNZ(t *testing.T) {
	type args struct {
		a16 uint16
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"a16=0xF00F, SP=0xCFFF, PC=0x0000, Z=0", Registers{PC: 0x0000, SP: 0xCFFF}, args{a16: 0xF00F}},
		{"a16=0xF00F, SP=0xCFFF, PC=0x0000, Z=1", Registers{PC: 0x0000, SP: 0xCFFF, F: 0xF0}, args{a16: 0xF00F}},
		{"a16=0x0000, SP=0xCFFF, PC=0xB00F, Z=0", Registers{PC: 0xB00F, SP: 0xCFFF, F: 0x70}, args{a16: 0xF00F}},
		{"a16=0x0000, SP=0xCFFF, PC=0xB00F, Z=1", Registers{PC: 0xB00F, SP: 0xCFFF, F: 0x80}, args{a16: 0xF00F}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			c.CallNZ(tt.args.a16)

			// If Z is 1, expect no change
			if tt.regs.F&0x80 != 0 {
				// Expect PC to be unchanged
				if c.PC != tt.regs.PC {
					t.Errorf("Wanted PC to be %04x, got %04x", tt.regs.PC, c.PC)
				}
				// Expect SP to be unchanged
				if c.SP != tt.regs.SP {
					t.Errorf("Wanted SP to be %04x, got %04x", tt.regs.SP, c.SP)
				}

				// Otherwise expect call
			} else {
				// Expect PC to be a16
				if c.PC != tt.args.a16 {
					t.Errorf("Wanted PC to be %04x, got %04x", tt.args.a16, c.PC)
				}
				// Expect SP to be SP-2
				if c.SP != tt.regs.SP-2 {
					t.Errorf("Wanted SP to be %04x, got %04x", tt.regs.SP-2, c.SP)
				}
				// Expect word at SP to contain previous PC
				valSP, err := c.mem.Rw(c.SP)
				if err != nil {
					t.Error(err)
				}
				if valSP != tt.regs.PC {
					t.Errorf("Wanted (SP) to be %04x, got %04x", tt.regs.PC, valSP)
				}
			}
		})
	}
}

func TestCPU_CallZ(t *testing.T) {
	type args struct {
		a16 uint16
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"a16=0xF00F, SP=0xCFFF, PC=0x0000, Z=0", Registers{PC: 0x0000, SP: 0xCFFF}, args{a16: 0xF00F}},
		{"a16=0xF00F, SP=0xCFFF, PC=0x0000, Z=1", Registers{PC: 0x0000, SP: 0xCFFF, F: 0xF0}, args{a16: 0xF00F}},
		{"a16=0x0000, SP=0xCFFF, PC=0xB00F, Z=0", Registers{PC: 0xB00F, SP: 0xCFFF, F: 0x70}, args{a16: 0xF00F}},
		{"a16=0x0000, SP=0xCFFF, PC=0xB00F, Z=1", Registers{PC: 0xB00F, SP: 0xCFFF, F: 0x80}, args{a16: 0xF00F}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			c.CallZ(tt.args.a16)

			// If Z is 0, expect no change
			if tt.regs.F&0x80 == 0 {
				// Expect PC to be unchanged
				if c.PC != tt.regs.PC {
					t.Errorf("Wanted PC to be %04x, got %04x", tt.regs.PC, c.PC)
				}
				// Expect SP to be unchanged
				if c.SP != tt.regs.SP {
					t.Errorf("Wanted SP to be %04x, got %04x", tt.regs.SP, c.SP)
				}

				// Otherwise expect call
			} else {
				// Expect PC to be a16
				if c.PC != tt.args.a16 {
					t.Errorf("Wanted PC to be %04x, got %04x", tt.args.a16, c.PC)
				}
				// Expect SP to be SP-2
				if c.SP != tt.regs.SP-2 {
					t.Errorf("Wanted SP to be %04x, got %04x", tt.regs.SP-2, c.SP)
				}
				// Expect word at SP to contain previous PC
				valSP, err := c.mem.Rw(c.SP)
				if err != nil {
					t.Error(err)
				}
				if valSP != tt.regs.PC {
					t.Errorf("Wanted (SP) to be %04x, got %04x", tt.regs.PC, valSP)
				}
			}
		})
	}
}

func TestCPU_CallNC(t *testing.T) {
	type args struct {
		a16 uint16
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"a16=0xF00F, SP=0xCFFF, PC=0x0000, C=0", Registers{PC: 0x0000, SP: 0xCFFF}, args{a16: 0xF00F}},
		{"a16=0xF00F, SP=0xCFFF, PC=0x0000, C=1", Registers{PC: 0x0000, SP: 0xCFFF, F: 0xF0}, args{a16: 0xF00F}},
		{"a16=0x0000, SP=0xCFFF, PC=0xB00F, C=0", Registers{PC: 0xB00F, SP: 0xCFFF}, args{a16: 0xF00F}},
		{"a16=0x0000, SP=0xCFFF, PC=0xB00F, C=1", Registers{PC: 0xB00F, SP: 0xCFFF, F: 0x10}, args{a16: 0xF00F}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			c.CallNC(tt.args.a16)

			// If C is 1, expect no change
			if tt.regs.F&0x10 != 0 {
				// Expect PC to be unchanged
				if c.PC != tt.regs.PC {
					t.Errorf("Wanted PC to be %04x, got %04x", tt.regs.PC, c.PC)
				}
				// Expect SP to be unchanged
				if c.SP != tt.regs.SP {
					t.Errorf("Wanted SP to be %04x, got %04x", tt.regs.SP, c.SP)
				}

				// Otherwise expect call
			} else {
				// Expect PC to be a16
				if c.PC != tt.args.a16 {
					t.Errorf("Wanted PC to be %04x, got %04x", tt.args.a16, c.PC)
				}
				// Expect SP to be SP-2
				if c.SP != tt.regs.SP-2 {
					t.Errorf("Wanted SP to be %04x, got %04x", tt.regs.SP-2, c.SP)
				}
				// Expect word at SP to contain previous PC
				valSP, err := c.mem.Rw(c.SP)
				if err != nil {
					t.Error(err)
				}
				if valSP != tt.regs.PC {
					t.Errorf("Wanted (SP) to be %04x, got %04x", tt.regs.PC, valSP)
				}
			}
		})
	}
}

func TestCPU_CallC(t *testing.T) {
	type args struct {
		a16 uint16
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		{"a16=0xF00F, SP=0xCFFF, PC=0x0000, C=0", Registers{PC: 0x0000, SP: 0xCFFF}, args{a16: 0xF00F}},
		{"a16=0xF00F, SP=0xCFFF, PC=0x0000, C=1", Registers{PC: 0x0000, SP: 0xCFFF, F: 0xF0}, args{a16: 0xF00F}},
		{"a16=0x0000, SP=0xCFFF, PC=0xB00F, C=0", Registers{PC: 0xB00F, SP: 0xCFFF}, args{a16: 0xF00F}},
		{"a16=0x0000, SP=0xCFFF, PC=0xB00F, C=1", Registers{PC: 0xB00F, SP: 0xCFFF, F: 0x10}, args{a16: 0xF00F}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			c.CallC(tt.args.a16)

			// If C is 0, expect no change
			if tt.regs.F&0x10 == 0 {
				// Expect PC to be unchanged
				if c.PC != tt.regs.PC {
					t.Errorf("Wanted PC to be %04x, got %04x", tt.regs.PC, c.PC)
				}
				// Expect SP to be unchanged
				if c.SP != tt.regs.SP {
					t.Errorf("Wanted SP to be %04x, got %04x", tt.regs.SP, c.SP)
				}

				// Otherwise expect call
			} else {
				// Expect PC to be a16
				if c.PC != tt.args.a16 {
					t.Errorf("Wanted PC to be %04x, got %04x", tt.args.a16, c.PC)
				}
				// Expect SP to be SP-2
				if c.SP != tt.regs.SP-2 {
					t.Errorf("Wanted SP to be %04x, got %04x", tt.regs.SP-2, c.SP)
				}
				// Expect word at SP to contain previous PC
				valSP, err := c.mem.Rw(c.SP)
				if err != nil {
					t.Error(err)
				}
				if valSP != tt.regs.PC {
					t.Errorf("Wanted (SP) to be %04x, got %04x", tt.regs.PC, valSP)
				}
			}
		})
	}
}

func TestCPU_Ret(t *testing.T) {
	tests := []struct {
		name  string
		valSP uint16
		regs  Registers
	}{
		{"SP=0xCFFF, (SP,SP+1)=0xF00F, PC=0x0000", 0xF00F, Registers{PC: 0x0000, SP: 0xCFFF}},
		{"SP=0xCFFD, (SP,SP+1)=0x000B, PC=0x1234", 0x000B, Registers{PC: 0x1234, SP: 0xCFFF}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			// initialize memory
			err := c.mem.Ww(c.SP, tt.valSP)
			if err != nil {
				t.Error(err)
			}

			c.Ret()

			// Expect new PC to be word at old SP
			expectedPC, err := c.mem.Rw(tt.regs.SP)
			if err != nil {
				t.Error(err)
			}
			if c.PC != expectedPC {
				t.Errorf("Wanted PC to be %04x, got %04x", expectedPC, c.PC)
			}
			// Expect SP to be SP+2
			if c.SP != tt.regs.SP+2 {
				t.Errorf("Wanted SP to be %04x, got %04x", tt.regs.SP+2, c.SP)
			}
		})
	}
}

func TestCPU_RetNZ(t *testing.T) {
	tests := []struct {
		name  string
		valSP uint16
		regs  Registers
	}{
		{"SP=0xCFFF, (SP,SP+1)=0xF00F, PC=0x0000, Z=0", 0xF00F, Registers{PC: 0x0000, SP: 0xCFFF}},
		{"SP=0xCFFF, (SP,SP+1)=0xF00F, PC=0x0000, Z=1", 0xF00F, Registers{PC: 0x0000, SP: 0xCFFF, F: 0xF0}},
		{"SP=0xCFFD, (SP,SP+1)=0x000B, PC=0x1234, Z=0", 0x000B, Registers{PC: 0x1234, SP: 0xCFFF, F: 0x70}},
		{"SP=0xCFFD, (SP,SP+1)=0x000B, PC=0x1234, Z=1", 0x000B, Registers{PC: 0x1234, SP: 0xCFFF, F: 0x80}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			// initialize memory
			err := c.mem.Ww(c.SP, tt.valSP)
			if err != nil {
				t.Error(err)
			}

			c.RetNZ()

			// if Z is 1, expect no change
			if c.F&0x80 != 0 {
				// expect PC to be unchanged
				if c.PC != tt.regs.PC {
					t.Errorf("Wanted PC to be %04x, got %04x", tt.regs.PC, c.PC)
				}
				// expect SP to be unchanged
				if c.SP != tt.regs.SP {
					t.Errorf("Wanted SP to be %04x, got %04x", tt.regs.SP, c.SP)
				}

				// otherwise expect a return
			} else {
				// Expect new PC to be word at old SP
				expectedPC, err := c.mem.Rw(tt.regs.SP)
				if err != nil {
					t.Error(err)
				}
				if c.PC != expectedPC {
					t.Errorf("Wanted PC to be %04x, got %04x", expectedPC, c.PC)
				}
				// Expect SP to be SP+2
				if c.SP != tt.regs.SP+2 {
					t.Errorf("Wanted SP to be %04x, got %04x", tt.regs.SP+2, c.SP)
				}
			}
		})
	}
}

func TestCPU_RetZ(t *testing.T) {
	tests := []struct {
		name  string
		valSP uint16
		regs  Registers
	}{
		{"SP=0xCFFF, (SP,SP+1)=0xF00F, PC=0x0000, Z=0", 0xF00F, Registers{PC: 0x0000, SP: 0xCFFF}},
		{"SP=0xCFFF, (SP,SP+1)=0xF00F, PC=0x0000, Z=1", 0xF00F, Registers{PC: 0x0000, SP: 0xCFFF, F: 0xF0}},
		{"SP=0xCFFD, (SP,SP+1)=0x000B, PC=0x1234, Z=0", 0x000B, Registers{PC: 0x1234, SP: 0xCFFF, F: 0x70}},
		{"SP=0xCFFD, (SP,SP+1)=0x000B, PC=0x1234, Z=1", 0x000B, Registers{PC: 0x1234, SP: 0xCFFF, F: 0x80}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			// initialize memory
			err := c.mem.Ww(c.SP, tt.valSP)
			if err != nil {
				t.Error(err)
			}

			c.RetZ()

			// if Z is 0, expect no change
			if c.F&0x80 == 0 {
				// expect PC to be unchanged
				if c.PC != tt.regs.PC {
					t.Errorf("Wanted PC to be %04x, got %04x", tt.regs.PC, c.PC)
				}
				// expect SP to be unchanged
				if c.SP != tt.regs.SP {
					t.Errorf("Wanted SP to be %04x, got %04x", tt.regs.SP, c.SP)
				}

				// otherwise expect a return
			} else {
				// Expect new PC to be word at old SP
				expectedPC, err := c.mem.Rw(tt.regs.SP)
				if err != nil {
					t.Error(err)
				}
				if c.PC != expectedPC {
					t.Errorf("Wanted PC to be %04x, got %04x", expectedPC, c.PC)
				}
				// Expect SP to be SP+2
				if c.SP != tt.regs.SP+2 {
					t.Errorf("Wanted SP to be %04x, got %04x", tt.regs.SP+2, c.SP)
				}
			}
		})
	}
}

func TestCPU_RetNC(t *testing.T) {
	tests := []struct {
		name  string
		valSP uint16
		regs  Registers
	}{
		{"SP=0xCFFF, (SP,SP+1)=0xF00F, PC=0x0000, C=0", 0xF00F, Registers{PC: 0x0000, SP: 0xCFFF}},
		{"SP=0xCFFF, (SP,SP+1)=0xF00F, PC=0x0000, C=1", 0xF00F, Registers{PC: 0x0000, SP: 0xCFFF, F: 0xF0}},
		{"SP=0xCFFD, (SP,SP+1)=0x000B, PC=0x1234, C=0", 0x000B, Registers{PC: 0x1234, SP: 0xCFFF, F: 0x20}},
		{"SP=0xCFFD, (SP,SP+1)=0x000B, PC=0x1234, C=1", 0x000B, Registers{PC: 0x1234, SP: 0xCFFF, F: 0x10}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			// initialize memory
			err := c.mem.Ww(c.SP, tt.valSP)
			if err != nil {
				t.Error(err)
			}

			c.RetNC()

			// if C is 1, expect no change
			if c.F&0x10 != 0 {
				// expect PC to be unchanged
				if c.PC != tt.regs.PC {
					t.Errorf("Wanted PC to be %04x, got %04x", tt.regs.PC, c.PC)
				}
				// expect SP to be unchanged
				if c.SP != tt.regs.SP {
					t.Errorf("Wanted SP to be %04x, got %04x", tt.regs.SP, c.SP)
				}

				// otherwise expect a return
			} else {
				// Expect new PC to be word at old SP
				expectedPC, err := c.mem.Rw(tt.regs.SP)
				if err != nil {
					t.Error(err)
				}
				if c.PC != expectedPC {
					t.Errorf("Wanted PC to be %04x, got %04x", expectedPC, c.PC)
				}
				// Expect SP to be SP+2
				if c.SP != tt.regs.SP+2 {
					t.Errorf("Wanted SP to be %04x, got %04x", tt.regs.SP+2, c.SP)
				}
			}
		})
	}
}

func TestCPU_RetC(t *testing.T) {
	tests := []struct {
		name  string
		valSP uint16
		regs  Registers
	}{
		{"SP=0xCFFF, (SP,SP+1)=0xF00F, PC=0x0000, C=0", 0xF00F, Registers{PC: 0x0000, SP: 0xCFFF}},
		{"SP=0xCFFF, (SP,SP+1)=0xF00F, PC=0x0000, C=1", 0xF00F, Registers{PC: 0x0000, SP: 0xCFFF, F: 0xF0}},
		{"SP=0xCFFD, (SP,SP+1)=0x000B, PC=0x1234, C=0", 0x000B, Registers{PC: 0x1234, SP: 0xCFFF, F: 0x20}},
		{"SP=0xCFFD, (SP,SP+1)=0x000B, PC=0x1234, C=1", 0x000B, Registers{PC: 0x1234, SP: 0xCFFF, F: 0x10}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			// initialize memory
			err := c.mem.Ww(c.SP, tt.valSP)
			if err != nil {
				t.Error(err)
			}

			c.RetC()

			// if C is 0, expect no change
			if c.F&0x10 == 0 {
				// expect PC to be unchanged
				if c.PC != tt.regs.PC {
					t.Errorf("Wanted PC to be %04x, got %04x", tt.regs.PC, c.PC)
				}
				// expect SP to be unchanged
				if c.SP != tt.regs.SP {
					t.Errorf("Wanted SP to be %04x, got %04x", tt.regs.SP, c.SP)
				}

				// otherwise expect a return
			} else {
				// Expect new PC to be word at old SP
				expectedPC, err := c.mem.Rw(tt.regs.SP)
				if err != nil {
					t.Error(err)
				}
				if c.PC != expectedPC {
					t.Errorf("Wanted PC to be %04x, got %04x", expectedPC, c.PC)
				}
				// Expect SP to be SP+2
				if c.SP != tt.regs.SP+2 {
					t.Errorf("Wanted SP to be %04x, got %04x", tt.regs.SP+2, c.SP)
				}
			}
		})
	}
}

func TestCPU_Reti(t *testing.T) {
	tests := []struct {
		name  string
		valSP uint16
		ime   bool
		regs  Registers
	}{
		{"SP=0xCFFF, (SP,SP+1)=0xF00F, IME=1, PC=0x0000", 0xF00F, true, Registers{PC: 0x0000, SP: 0xCFFF}},
		{"SP=0xCFFF, (SP,SP+1)=0xF00F, IME=0, PC=0x0000", 0xF00F, false, Registers{PC: 0x0000, SP: 0xCFFF}},
		{"SP=0xCFFD, (SP,SP+1)=0x000B, IME=1, PC=0x1234", 0x000B, true, Registers{PC: 0x1234, SP: 0xCFFF}},
		{"SP=0xCFFD, (SP,SP+1)=0x000B, IME=0, PC=0x1234", 0x000B, false, Registers{PC: 0x1234, SP: 0xCFFF}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs
			c.ime = tt.ime

			// initialize memory
			err := c.mem.Ww(c.SP, tt.valSP)
			if err != nil {
				t.Error(err)
			}

			c.Reti()

			// Expect interrupts to be enabled
			if c.ime != true {
				t.Error("Wanted IME to be true, was false")
			}

			// Expect new PC to be word at old SP
			expectedPC, err := c.mem.Rw(tt.regs.SP)
			if err != nil {
				t.Error(err)
			}
			if c.PC != expectedPC {
				t.Errorf("Wanted PC to be %04x, got %04x", expectedPC, c.PC)
			}
			// Expect SP to be SP+2
			if c.SP != tt.regs.SP+2 {
				t.Errorf("Wanted SP to be %04x, got %04x", tt.regs.SP+2, c.SP)
			}
		})
	}
}

func TestCPU_Rst(t *testing.T) {
	type args struct {
		n byte
	}
	tests := []struct {
		name string
		regs Registers
		args args
	}{
		// valid n
		{"PC=0xF00D, SP=0xCFFF, n=00", Registers{PC: 0xF00D, SP: 0xCFFF}, args{n: 0x00}},
		{"PC=0xF00D, SP=0xCFFF, n=08", Registers{PC: 0xF00D, SP: 0xCFFF}, args{n: 0x08}},
		{"PC=0xF00D, SP=0xCFFF, n=10", Registers{PC: 0xF00D, SP: 0xCFFF}, args{n: 0x10}},
		{"PC=0xF00D, SP=0xCFFF, n=18", Registers{PC: 0xF00D, SP: 0xCFFF}, args{n: 0x18}},
		{"PC=0xF00D, SP=0xCFFF, n=20", Registers{PC: 0xF00D, SP: 0xCFFF}, args{n: 0x20}},
		{"PC=0xF00D, SP=0xCFFF, n=28", Registers{PC: 0xF00D, SP: 0xCFFF}, args{n: 0x28}},
		{"PC=0xF00D, SP=0xCFFF, n=30", Registers{PC: 0xF00D, SP: 0xCFFF}, args{n: 0x30}},
		{"PC=0xF00D, SP=0xCFFF, n=38", Registers{PC: 0xF00D, SP: 0xCFFF}, args{n: 0x38}},
		// invalid n
		{"PC=0xF00D, SP=0xCFFF, n=FF", Registers{PC: 0xF00D, SP: 0xCFFF}, args{n: 0xFF}},
		{"PC=0xF00D, SP=0xCFFF, n=01", Registers{PC: 0xF00D, SP: 0xCFFF}, args{n: 0x01}},
		{"PC=0xF00D, SP=0xCFFF, n=83", Registers{PC: 0xF00D, SP: 0xCFFF}, args{n: 0x83}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := mmu.New()
			c := New(mmu)
			c.Registers = tt.regs

			// Expect call to $0000 + n IF n is 00, 08, 10, 18, 20, 28, 30, 38
			switch tt.args.n {
			case 0x00, 0x08, 0x10, 0x18, 0x20, 0x28, 0x30, 0x38:

				c.Rst(RstTarget(tt.args.n))

				// Expect PC to be $0000 + n
				if c.PC != uint16(tt.args.n) {
					t.Errorf("Expected PC to be %04x, got %04x", uint16(tt.args.n), c.PC)
				}
				// Expect word at SP to be old PC
				valSP, err := c.mem.Rw(c.SP)
				if err != nil {
					t.Error(err)
				}
				if valSP != tt.regs.PC {
					t.Errorf("Expected (SP) to be %04x, got %04x", tt.regs.PC, valSP)
				}
				// Expect SP to be old SP-2
				if c.SP != tt.regs.SP-2 {
					t.Errorf("Expected SP to be %04x, got %04x", tt.regs.SP-2, c.SP)
				}

			default:
				// Expect panic if n is any other value
				func() {
					defer func() {
						r := recover()
						if r == nil {
							t.Errorf("Expected reset to panic when passed invalid value %02x", tt.args.n)
						}
					}()
					c.Rst(RstTarget(tt.args.n))
				}()
			}
		})
	}
}
