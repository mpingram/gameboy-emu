package cpu

type UnprefixedOpcode struct {
	val          OpcodeValue
	prefixCB     bool
	mnemonic     string
	length       uint8
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
