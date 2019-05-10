package cpu

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
