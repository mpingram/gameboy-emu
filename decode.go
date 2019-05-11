package cpu

type memoryReadWriter interface {
	memoryReader
	memoryWriter
}

type memoryReader interface {
	rb(addr uint16) (byte, error)
	rw(addr uint16) (uint16, error)
}

type memoryWriter interface {
	wb(addr uint16, b byte) error
	ww(addr uint16, bb uint16) error
}

type Instruction struct {
	opc  Opcode
	data []byte
}

func (i *Instruction) String() string {
	return i.opc.mnemonic
}

func Decode(addr uint16, mem memoryReader) (Instruction, error) {

	// FIXME: EDGE CASE: 'HALT' opcode may be 1 or 2 bytes long. Officially,
	// it's supposed to be 0x01 0x00 (which looks like HALT, NOP), so some
	// programs just give 0x01.
	// Need to figure out how to handle that information.

	var opc Opcode

	// =======================
	// Get the first byte of the opcode and look it up in our opcode table.
	// =======================
	firstByte, err := mem.rb(addr)
	if err != nil {
		return Instruction{}, err
	}
	// 0xCB is a special 'prefix' opcode. If it is present, the next
	// byte represents an opcode in the CB prefix table.
	cbPrefix := byte(0xCB)
	isPrefixed := firstByte == cbPrefix

	if isPrefixed {
		secondByte, err := mem.rb(addr + 1)
		if err != nil {
			return Instruction{}, err
		}
		val := OpcodeValue(secondByte)
		// look up opcode in opcode value table for cb-prefixed opcodes
		opc = prefixedOpcodes[val]
	} else {
		val := OpcodeValue(firstByte)
		// look up opcode in opcode value table for non-prefixed opcodes
		opc = unprefixedOpcodes[val]
	}

	// =======================
	// Get the extra data that goes along with this opcode, if any.
	// If there is extra data, it may be 1 or 2 bytes long.
	// =======================
	var data []byte
	if isPrefixed {
		// all prefixed opcodes happen to have no extra data.
		data = make([]byte, 0)
	} else {
		switch opc.length {
		case 1:
			// opcodes of byte length 1 take no argruments
			data = make([]byte, 0)
		case 2:
			// if length is 2, next byte in memory is argument
			data = make([]byte, 1)
			argByte, err := mem.rb(addr + 1)
			if err != nil {
				return Instruction{}, err
			}
			data[0] = argByte

		case 3:
			// if length is 3, next two bytes in memory are argument
			data = make([]byte, 2)
			argByte1, err := mem.rb(addr + 1)
			argByte2, err := mem.rb(addr + 2)
			if err != nil {
				return Instruction{}, err
			}
			data[0] = argByte1
			data[1] = argByte2
		}
	}

	// =======================
	// Wrap our opcode in an Instruction, including the extra data.
	// =======================
	return Instruction{opc, data}, nil

}
