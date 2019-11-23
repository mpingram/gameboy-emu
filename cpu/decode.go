package cpu

import (
	"fmt"
	"strings"
)

type Instruction struct {
	opc  Opcode
	data []byte
}

type Opcode struct {
	val          OpcodeValue
	prefixCB     bool
	mnemonic     string
	length       uint16
	cycles       int
	cyclesIfNoop int // Some conditional instructions, eg JR NZ, take different amounts of time depending on if the condition is executed or not
	flags        Flags
}

type Flags struct {
	Z, H, N, C FlagState
}
type FlagState uint8

const (
	NoChange  FlagState = 0
	CanChange FlagState = 1
	IsSet     FlagState = 2
	IsReset   FlagState = 3
)

type OpcodeValue uint8

func (i *Instruction) String() string {
	mnemonic := i.opc.mnemonic
	if strings.Contains(mnemonic, "d8") {
		return fmt.Sprintf("%s: 0x%02x", mnemonic, d8(i.data))
	} else if strings.Contains(mnemonic, "d16") {
		return fmt.Sprintf("%s: 0x%04x", mnemonic, d16(i.data))
	} else if strings.Contains(mnemonic, "a8") {
		return fmt.Sprintf("%s: $%02X", mnemonic, a8(i.data))
	} else if strings.Contains(mnemonic, "a16") {
		return fmt.Sprintf("%s: $%04X", mnemonic, a16(i.data))
	} else if strings.Contains(mnemonic, "r8") {
		return fmt.Sprintf("%s: 0x%02x", mnemonic, r8(i.data))
	} else {
		return fmt.Sprintf("%s", mnemonic)
	}
}

func Decode(addr uint16, mem MemoryReader) Instruction {

	// FIXME: EDGE CASE: 'HALT' opcode may be 1 or 2 bytes long. Officially,
	// it's supposed to be 0x01 0x00 (which looks like HALT, NOP), so some
	// programs just give 0x01.
	// Need to figure out how to handle that information.

	var opc Opcode

	// =======================
	// Get the first byte of the opcode and look it up in our opcode table.
	// =======================
	firstByte := mem.Rb(addr)
	// 0xCB is a special 'prefix' opcode. If it is present, the next
	// byte represents an opcode in the CB prefix table.
	cbPrefix := byte(0xCB)
	isPrefixed := firstByte == cbPrefix

	if isPrefixed {
		secondByte := mem.Rb(addr + 1)
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
			argByte := mem.Rb(addr + 1)
			data[0] = argByte

		case 3:
			// if length is 3, next two bytes in memory are argument
			data = make([]byte, 2)
			argByte1 := mem.Rb(addr + 1)
			argByte2 := mem.Rb(addr + 2)
			data[0] = argByte1
			data[1] = argByte2
		}
	}

	// =======================
	// Wrap our opcode in an Instruction, including the extra data.
	// =======================
	return Instruction{opc, data}

}
