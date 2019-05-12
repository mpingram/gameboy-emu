package cpu

import "testing"

func TestCPU_Jp(t *testing.T) {

	type args struct {
		a16 uint16
	}
	tests := []struct {
		name      string
		Registers Registers
		args      args
	}{
		{"Jp to should move PC to a16", Registers{}, args{0x9BFF}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers: tt.Registers,
			}
			c.Jp(tt.args.a16)
			// Expect PC to be a16
			if c.PC != tt.args.a16 {
				t.Errorf("Expected PC %04x, got %04x", tt.args.a16, c.PC)
			}
			// Expect SP to not have moved
			if c.SP != tt.Registers.SP {
				t.Errorf("Expected SP to be unchanged from %04x; got %04x", tt.Registers.SP, c.SP)
			}
		})
	}
}

func TestCPU_Jp_HL(t *testing.T) {
	tests := []struct {
		name        string
		RegistersIn Registers
	}{
		{"HL = 0xF00F", Registers{H: 0xF0, L: 0x0F}},
		{"HL = 0xB00B", Registers{H: 0xB0, L: 0x0B}},
		{"HL = 0x1234", Registers{H: 0x12, L: 0x34}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers: tt.RegistersIn,
			}
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
			c := &CPU{
				Registers: tt.regs,
			}
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
			c := &CPU{
				Registers: tt.regs,
			}
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
			c := &CPU{
				Registers: tt.regs,
			}
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
			c := &CPU{
				Registers: tt.regs,
			}
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
			c := &CPU{
				Registers: tt.regs,
			}
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
	type fields struct {
		Registers Registers
		TClock    <-chan int
		MClock    <-chan int
		mem       MMU
		halted    bool
		stopped   bool
		ime       bool
		setIME    bool
	}
	type args struct {
		r8 int8
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers: tt.fields.Registers,
				TClock:    tt.fields.TClock,
				MClock:    tt.fields.MClock,
				mem:       tt.fields.mem,
				halted:    tt.fields.halted,
				stopped:   tt.fields.stopped,
				ime:       tt.fields.ime,
				setIME:    tt.fields.setIME,
			}
			c.JrNZ(tt.args.r8)
		})
	}
}

func TestCPU_JrZ(t *testing.T) {
	type fields struct {
		Registers Registers
		TClock    <-chan int
		MClock    <-chan int
		mem       MMU
		halted    bool
		stopped   bool
		ime       bool
		setIME    bool
	}
	type args struct {
		r8 int8
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers: tt.fields.Registers,
				TClock:    tt.fields.TClock,
				MClock:    tt.fields.MClock,
				mem:       tt.fields.mem,
				halted:    tt.fields.halted,
				stopped:   tt.fields.stopped,
				ime:       tt.fields.ime,
				setIME:    tt.fields.setIME,
			}
			c.JrZ(tt.args.r8)
		})
	}
}

func TestCPU_JrNC(t *testing.T) {
	type fields struct {
		Registers Registers
		TClock    <-chan int
		MClock    <-chan int
		mem       MMU
		halted    bool
		stopped   bool
		ime       bool
		setIME    bool
	}
	type args struct {
		r8 int8
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers: tt.fields.Registers,
				TClock:    tt.fields.TClock,
				MClock:    tt.fields.MClock,
				mem:       tt.fields.mem,
				halted:    tt.fields.halted,
				stopped:   tt.fields.stopped,
				ime:       tt.fields.ime,
				setIME:    tt.fields.setIME,
			}
			c.JrNC(tt.args.r8)
		})
	}
}

func TestCPU_JrC(t *testing.T) {
	type fields struct {
		Registers Registers
		TClock    <-chan int
		MClock    <-chan int
		mem       MMU
		halted    bool
		stopped   bool
		ime       bool
		setIME    bool
	}
	type args struct {
		r8 int8
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers: tt.fields.Registers,
				TClock:    tt.fields.TClock,
				MClock:    tt.fields.MClock,
				mem:       tt.fields.mem,
				halted:    tt.fields.halted,
				stopped:   tt.fields.stopped,
				ime:       tt.fields.ime,
				setIME:    tt.fields.setIME,
			}
			c.JrC(tt.args.r8)
		})
	}
}

func TestCPU_Call(t *testing.T) {
	type fields struct {
		Registers Registers
		TClock    <-chan int
		MClock    <-chan int
		mem       MMU
		halted    bool
		stopped   bool
		ime       bool
		setIME    bool
	}
	type args struct {
		a16 uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers: tt.fields.Registers,
				TClock:    tt.fields.TClock,
				MClock:    tt.fields.MClock,
				mem:       tt.fields.mem,
				halted:    tt.fields.halted,
				stopped:   tt.fields.stopped,
				ime:       tt.fields.ime,
				setIME:    tt.fields.setIME,
			}
			c.Call(tt.args.a16)
		})
	}
}

func TestCPU_CallNZ(t *testing.T) {
	type fields struct {
		Registers Registers
		TClock    <-chan int
		MClock    <-chan int
		mem       MMU
		halted    bool
		stopped   bool
		ime       bool
		setIME    bool
	}
	type args struct {
		a16 uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers: tt.fields.Registers,
				TClock:    tt.fields.TClock,
				MClock:    tt.fields.MClock,
				mem:       tt.fields.mem,
				halted:    tt.fields.halted,
				stopped:   tt.fields.stopped,
				ime:       tt.fields.ime,
				setIME:    tt.fields.setIME,
			}
			c.CallNZ(tt.args.a16)
		})
	}
}

func TestCPU_CallZ(t *testing.T) {
	type fields struct {
		Registers Registers
		TClock    <-chan int
		MClock    <-chan int
		mem       MMU
		halted    bool
		stopped   bool
		ime       bool
		setIME    bool
	}
	type args struct {
		a16 uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers: tt.fields.Registers,
				TClock:    tt.fields.TClock,
				MClock:    tt.fields.MClock,
				mem:       tt.fields.mem,
				halted:    tt.fields.halted,
				stopped:   tt.fields.stopped,
				ime:       tt.fields.ime,
				setIME:    tt.fields.setIME,
			}
			c.CallZ(tt.args.a16)
		})
	}
}

func TestCPU_CallNC(t *testing.T) {
	type fields struct {
		Registers Registers
		TClock    <-chan int
		MClock    <-chan int
		mem       MMU
		halted    bool
		stopped   bool
		ime       bool
		setIME    bool
	}
	type args struct {
		a16 uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers: tt.fields.Registers,
				TClock:    tt.fields.TClock,
				MClock:    tt.fields.MClock,
				mem:       tt.fields.mem,
				halted:    tt.fields.halted,
				stopped:   tt.fields.stopped,
				ime:       tt.fields.ime,
				setIME:    tt.fields.setIME,
			}
			c.CallNC(tt.args.a16)
		})
	}
}

func TestCPU_CallC(t *testing.T) {
	type fields struct {
		Registers Registers
		TClock    <-chan int
		MClock    <-chan int
		mem       MMU
		halted    bool
		stopped   bool
		ime       bool
		setIME    bool
	}
	type args struct {
		a16 uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers: tt.fields.Registers,
				TClock:    tt.fields.TClock,
				MClock:    tt.fields.MClock,
				mem:       tt.fields.mem,
				halted:    tt.fields.halted,
				stopped:   tt.fields.stopped,
				ime:       tt.fields.ime,
				setIME:    tt.fields.setIME,
			}
			c.CallC(tt.args.a16)
		})
	}
}

func TestCPU_Ret(t *testing.T) {
	type fields struct {
		Registers Registers
		TClock    <-chan int
		MClock    <-chan int
		mem       MMU
		halted    bool
		stopped   bool
		ime       bool
		setIME    bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers: tt.fields.Registers,
				TClock:    tt.fields.TClock,
				MClock:    tt.fields.MClock,
				mem:       tt.fields.mem,
				halted:    tt.fields.halted,
				stopped:   tt.fields.stopped,
				ime:       tt.fields.ime,
				setIME:    tt.fields.setIME,
			}
			c.Ret()
		})
	}
}

func TestCPU_RetNZ(t *testing.T) {
	type fields struct {
		Registers Registers
		TClock    <-chan int
		MClock    <-chan int
		mem       MMU
		halted    bool
		stopped   bool
		ime       bool
		setIME    bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers: tt.fields.Registers,
				TClock:    tt.fields.TClock,
				MClock:    tt.fields.MClock,
				mem:       tt.fields.mem,
				halted:    tt.fields.halted,
				stopped:   tt.fields.stopped,
				ime:       tt.fields.ime,
				setIME:    tt.fields.setIME,
			}
			c.RetNZ()
		})
	}
}

func TestCPU_RetZ(t *testing.T) {
	type fields struct {
		Registers Registers
		TClock    <-chan int
		MClock    <-chan int
		mem       MMU
		halted    bool
		stopped   bool
		ime       bool
		setIME    bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers: tt.fields.Registers,
				TClock:    tt.fields.TClock,
				MClock:    tt.fields.MClock,
				mem:       tt.fields.mem,
				halted:    tt.fields.halted,
				stopped:   tt.fields.stopped,
				ime:       tt.fields.ime,
				setIME:    tt.fields.setIME,
			}
			c.RetZ()
		})
	}
}

func TestCPU_RetNC(t *testing.T) {
	type fields struct {
		Registers Registers
		TClock    <-chan int
		MClock    <-chan int
		mem       MMU
		halted    bool
		stopped   bool
		ime       bool
		setIME    bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers: tt.fields.Registers,
				TClock:    tt.fields.TClock,
				MClock:    tt.fields.MClock,
				mem:       tt.fields.mem,
				halted:    tt.fields.halted,
				stopped:   tt.fields.stopped,
				ime:       tt.fields.ime,
				setIME:    tt.fields.setIME,
			}
			c.RetNC()
		})
	}
}

func TestCPU_RetC(t *testing.T) {
	type fields struct {
		Registers Registers
		TClock    <-chan int
		MClock    <-chan int
		mem       MMU
		halted    bool
		stopped   bool
		ime       bool
		setIME    bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers: tt.fields.Registers,
				TClock:    tt.fields.TClock,
				MClock:    tt.fields.MClock,
				mem:       tt.fields.mem,
				halted:    tt.fields.halted,
				stopped:   tt.fields.stopped,
				ime:       tt.fields.ime,
				setIME:    tt.fields.setIME,
			}
			c.RetC()
		})
	}
}

func TestCPU_Reti(t *testing.T) {
	type fields struct {
		Registers Registers
		TClock    <-chan int
		MClock    <-chan int
		mem       MMU
		halted    bool
		stopped   bool
		ime       bool
		setIME    bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers: tt.fields.Registers,
				TClock:    tt.fields.TClock,
				MClock:    tt.fields.MClock,
				mem:       tt.fields.mem,
				halted:    tt.fields.halted,
				stopped:   tt.fields.stopped,
				ime:       tt.fields.ime,
				setIME:    tt.fields.setIME,
			}
			c.Reti()
		})
	}
}

func TestCPU_Rst(t *testing.T) {
	type fields struct {
		Registers Registers
		TClock    <-chan int
		MClock    <-chan int
		mem       MMU
		halted    bool
		stopped   bool
		ime       bool
		setIME    bool
	}
	type args struct {
		n byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers: tt.fields.Registers,
				TClock:    tt.fields.TClock,
				MClock:    tt.fields.MClock,
				mem:       tt.fields.mem,
				halted:    tt.fields.halted,
				stopped:   tt.fields.stopped,
				ime:       tt.fields.ime,
				setIME:    tt.fields.setIME,
			}
			c.Rst(tt.args.n)
		})
	}
}
