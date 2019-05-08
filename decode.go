package cpu

type MemoryReadWriter interface {
	rb(addr uint16) uint8
	wb(addr uint16)

	rw(addr uint16) uint16
	ww(addr uint16)
}

type Instruction struct {
	opc      Opcode
	d8, a8   uint8
	d16, a16 uint16
	r8       int8
}

func (i *Instruction) String() {
	return i.opc.mnemonic
}

func Decode(addr uint16, mem *MemoryReadWriter) (Instruction, int, error) {

	var opc Opcode

	// 0xCB is a special 'prefix' opcode. If it is present, the next
	// byte represents an opcode in the CB prefix table.
	opcPrefix := uint8(0xCB)
	if mem.rb(addr) == opcPrefix {
		// we represent this prefixed opcode as 0xCBXX, where XX is the value
		// of the opcode byte after the prefix.
		opc = Opcode(0xCB00 | uint16(mem.rb(addr+1)))
	} else {
		opc = Opcode(mem.rb(addr))
	}

	switch opc {
	case NOP:

	}

	return Instruction{opc, 0, 0, 0, 0, 0, ""}, 0, nil
}
