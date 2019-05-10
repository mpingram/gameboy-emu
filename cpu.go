package cpu

type CPU struct {
	Registers

	TClock <-chan int
	MClock <-chan int

	mem memoryReadWriter

	halted  bool
	stopped bool
}

// Memory addresses
const ADDR_INTERRUPTS uint16 = 0x0

// Interrupt flags
const IME_BIT = 0x0
