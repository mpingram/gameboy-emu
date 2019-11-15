package cpu

import (
	"github.com/mpingram/gameboy-emu/mmu"
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
	cyclesIfNoop int
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
	return i.opc.mnemonic
}

func Decode(addr uint16, mem mmu.MemoryReader) (Instruction, error) {

	// FIXME: EDGE CASE: 'HALT' opcode may be 1 or 2 bytes long. Officially,
	// it's supposed to be 0x01 0x00 (which looks like HALT, NOP), so some
	// programs just give 0x01.
	// Need to figure out how to handle that information.

	var opc Opcode

	// =======================
	// Get the first byte of the opcode and look it up in our opcode table.
	// =======================
	firstByte, err := mem.Rb(addr)
	if err != nil {
		return Instruction{}, err
	}
	// 0xCB is a special 'prefix' opcode. If it is present, the next
	// byte represents an opcode in the CB prefix table.
	cbPrefix := byte(0xCB)
	isPrefixed := firstByte == cbPrefix

	if isPrefixed {
		secondByte, err := mem.Rb(addr + 1)
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
			argByte, err := mem.Rb(addr + 1)
			if err != nil {
				return Instruction{}, err
			}
			data[0] = argByte

		case 3:
			// if length is 3, next two bytes in memory are argument
			data = make([]byte, 2)
			argByte1, err := mem.Rb(addr + 1)
			argByte2, err := mem.Rb(addr + 2)
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
