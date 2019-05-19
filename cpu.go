package cpu

import "github.com/mpingram/gameboy-emu/mmu"

// New initializes and returns an instance of CPU.
func New(mmu *mmu.MMU) *CPU {
	cpu := &CPU{mem: mmu.CPUInterface}
	return cpu
}

type CPU struct {
	Registers

	TClock <-chan int
	MClock <-chan int

	mem mmu.MemoryReadWriter

	halted  bool
	stopped bool

	ime    bool // Interrupt master enable
	setIME bool // set IME next instruction (used for Ei() command)

}

/**
 * Memory addresses
 */

// ADDR_IE is the address of the IE (Interrupts Enabled) byte in memory,
// IE byte:
// FIXME documentation
const ADDR_IE uint16 = 0x0

// ADDR_IF is the address of the IF (Interrupt Flags) byte in memory
// IF byte:
// FIXME documentation
const ADDR_IF uint16 = 0x0
