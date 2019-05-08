package cpu

const (
	NOP OpcodeValue = 0x00
)

const (
	RLC_B UnprefixedOpcodeValue = 0x00
)

var unprefixedOpcodes = map[uint8]Opcode{
	0x00: Opcode{0x00, false, "NOP", 1, 4, Flags{}},
}

var prefixedOpcodes = map[uint8]Opcode{
	0x00: Opcode{0x00, true, "RLC B", 2, 8, Flags{}},
}
