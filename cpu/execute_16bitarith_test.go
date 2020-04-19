package cpu

import (
	"testing"
	"time"
)

func TestCPU_Add_HL_rr(t *testing.T) {
	type fields struct {
		Registers     Registers
		Clock         <-chan time.Time
		mem           MemoryReadWriter
		halted        bool
		stopped       bool
		readyForStart bool
		ime           bool
		setIME        bool
		breakpoint    uint16
	}
	type args struct {
		rr Reg16
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
				Registers:     tt.fields.Registers,
				Clock:         tt.fields.Clock,
				mem:           tt.fields.mem,
				halted:        tt.fields.halted,
				stopped:       tt.fields.stopped,
				readyForStart: tt.fields.readyForStart,
				ime:           tt.fields.ime,
				setIME:        tt.fields.setIME,
				breakpoint:    tt.fields.breakpoint,
			}
			c.Add_HL_rr(tt.args.rr)
		})
	}
}

func TestCPU_Add_HL_SP(t *testing.T) {
	type fields struct {
		Registers     Registers
		Clock         <-chan time.Time
		mem           MemoryReadWriter
		halted        bool
		stopped       bool
		readyForStart bool
		ime           bool
		setIME        bool
		breakpoint    uint16
	}
	tests := []struct {
		name      string
		fieldsIn  fields
		fieldsOut fields
	}{
		// https://stackoverflow.com/questions/57958631/game-boy-half-carry-flag-and-16-bit-instructions-especially-opcode-0xe8
		{
			"SP=0xFFFF, r8=-1 -> SP=0xFFFE, H=1, C=1",
			fields{Registers: Registers{SP: 0xFFFF, F: 0b0000_1110}},
			fields{Registers: Registers{SP: 0xFFFE, F: 0b0000_1011}},
		},
		{
			"SP=0x0000, r8=-1 -> SP=0xFFFF, H=1, C=1",
			fields{Registers: Registers{SP: 0xFFFF, F: 0b0000_1110}},
			fields{Registers: Registers{SP: 0xFFFF, F: 0b0000_1011}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers:     tt.fieldsIn.Registers,
				Clock:         tt.fieldsIn.Clock,
				mem:           tt.fieldsIn.mem,
				halted:        tt.fieldsIn.halted,
				stopped:       tt.fieldsIn.stopped,
				readyForStart: tt.fieldsIn.readyForStart,
				ime:           tt.fieldsIn.ime,
				setIME:        tt.fieldsIn.setIME,
				breakpoint:    tt.fieldsIn.breakpoint,
			}
			c.Add_HL_SP()
		})
	}
}

func TestCPU_Inc_rr(t *testing.T) {
	type fields struct {
		Registers     Registers
		Clock         <-chan time.Time
		mem           MemoryReadWriter
		halted        bool
		stopped       bool
		readyForStart bool
		ime           bool
		setIME        bool
		breakpoint    uint16
	}
	type args struct {
		rr Reg16
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
				Registers:     tt.fields.Registers,
				Clock:         tt.fields.Clock,
				mem:           tt.fields.mem,
				halted:        tt.fields.halted,
				stopped:       tt.fields.stopped,
				readyForStart: tt.fields.readyForStart,
				ime:           tt.fields.ime,
				setIME:        tt.fields.setIME,
				breakpoint:    tt.fields.breakpoint,
			}
			c.Inc_rr(tt.args.rr)
		})
	}
}

func TestCPU_Inc_SP(t *testing.T) {
	type fields struct {
		Registers     Registers
		Clock         <-chan time.Time
		mem           MemoryReadWriter
		halted        bool
		stopped       bool
		readyForStart bool
		ime           bool
		setIME        bool
		breakpoint    uint16
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
				Registers:     tt.fields.Registers,
				Clock:         tt.fields.Clock,
				mem:           tt.fields.mem,
				halted:        tt.fields.halted,
				stopped:       tt.fields.stopped,
				readyForStart: tt.fields.readyForStart,
				ime:           tt.fields.ime,
				setIME:        tt.fields.setIME,
				breakpoint:    tt.fields.breakpoint,
			}
			c.Inc_SP()
		})
	}
}

func TestCPU_Dec_rr(t *testing.T) {
	type fields struct {
		Registers     Registers
		Clock         <-chan time.Time
		mem           MemoryReadWriter
		halted        bool
		stopped       bool
		readyForStart bool
		ime           bool
		setIME        bool
		breakpoint    uint16
	}
	type args struct {
		rr Reg16
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
				Registers:     tt.fields.Registers,
				Clock:         tt.fields.Clock,
				mem:           tt.fields.mem,
				halted:        tt.fields.halted,
				stopped:       tt.fields.stopped,
				readyForStart: tt.fields.readyForStart,
				ime:           tt.fields.ime,
				setIME:        tt.fields.setIME,
				breakpoint:    tt.fields.breakpoint,
			}
			c.Dec_rr(tt.args.rr)
		})
	}
}

func TestCPU_Dec_SP(t *testing.T) {
	type fields struct {
		Registers     Registers
		Clock         <-chan time.Time
		mem           MemoryReadWriter
		halted        bool
		stopped       bool
		readyForStart bool
		ime           bool
		setIME        bool
		breakpoint    uint16
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
				Registers:     tt.fields.Registers,
				Clock:         tt.fields.Clock,
				mem:           tt.fields.mem,
				halted:        tt.fields.halted,
				stopped:       tt.fields.stopped,
				readyForStart: tt.fields.readyForStart,
				ime:           tt.fields.ime,
				setIME:        tt.fields.setIME,
				breakpoint:    tt.fields.breakpoint,
			}
			c.Dec_SP()
		})
	}
}

func TestCPU_Add_SP_r8(t *testing.T) {
	type fields struct {
		Registers     Registers
		Clock         <-chan time.Time
		mem           MemoryReadWriter
		halted        bool
		stopped       bool
		readyForStart bool
		ime           bool
		setIME        bool
		breakpoint    uint16
	}
	type args struct {
		r8 int8
	}
	tests := []struct {
		name      string
		args      args
		fieldsIn  fields
		fieldsOut fields
	}{
		{
			"SP=0x0000, r8=0x20 -> SP=0x0020, H=0, C=0",
			args{r8: 0x20},
			fields{Registers: Registers{SP: 0x0000, F: 0b0000_1101}},
			fields{Registers: Registers{SP: 0x0020, F: 0b0000_1000}},
		},
		{
			"Test half carry: SP=0x0008, r8=0x01 -> SP=0x0009, H=0, C=0",
			args{r8: 0x20},
			fields{Registers: Registers{SP: 0x0000, F: 0b0000_1101}},
			fields{Registers: Registers{SP: 0x0020, F: 0b0000_1000}},
		},
		// https://stackoverflow.com/questions/57958631/game-boy-half-carry-flag-and-16-bit-instructions-especially-opcode-0xe8
		{
			"Test carry/half carry edge case: SP=0xFFFF, r8=-1 -> SP=0xFFFE, H=1, C=1",
			args{r8: -1},
			fields{Registers: Registers{SP: 0xFFFF, F: 0b0000_1110}},
			fields{Registers: Registers{SP: 0xFFFE, F: 0b0000_1011}},
		},
		{
			"Test overflow case: SP=0x0000, r8=-1 -> SP=0xFFFF, H=0, C=0",
			args{r8: -1},
			fields{Registers: Registers{SP: 0xFFFF, F: 0b0000_0101}},
			fields{Registers: Registers{SP: 0xFFFF, F: 0b0000_0000}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CPU{
				Registers:     tt.fields.Registers,
				Clock:         tt.fields.Clock,
				mem:           tt.fields.mem,
				halted:        tt.fields.halted,
				stopped:       tt.fields.stopped,
				readyForStart: tt.fields.readyForStart,
				ime:           tt.fields.ime,
				setIME:        tt.fields.setIME,
				breakpoint:    tt.fields.breakpoint,
			}
			c.Add_SP_r8(tt.args.r8)
		})
	}
}

func TestCPU_Ld_HL_SPplusr8(t *testing.T) {
	type fields struct {
		Registers     Registers
		Clock         <-chan time.Time
		mem           MemoryReadWriter
		halted        bool
		stopped       bool
		readyForStart bool
		ime           bool
		setIME        bool
		breakpoint    uint16
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
				Registers:     tt.fields.Registers,
				Clock:         tt.fields.Clock,
				mem:           tt.fields.mem,
				halted:        tt.fields.halted,
				stopped:       tt.fields.stopped,
				readyForStart: tt.fields.readyForStart,
				ime:           tt.fields.ime,
				setIME:        tt.fields.setIME,
				breakpoint:    tt.fields.breakpoint,
			}
			c.Ld_HL_SPplusr8(tt.args.r8)
		})
	}
}
