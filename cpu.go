package cpu

type MMU struct {
	mem []byte
}

func (m *MMU) init() {
	m.mem = make([]byte, 0xFFFF)
}

func (m *MMU) rb(addr uint16) (byte, error) {
	if len(m.mem) == 0 {
		m.init()
	}
	return m.mem[addr], nil
}
func (m *MMU) wb(addr uint16, b byte) error {
	if len(m.mem) == 0 {
		m.init()
	}
	m.mem[addr] = b
	return nil
}

func (m *MMU) rw(addr uint16) (uint16, error) {
	if len(m.mem) == 0 {
		m.init()
	}
	hi := m.mem[addr]
	lo := m.mem[addr+1]
	return uint16(hi)<<8 | uint16(lo), nil
}

func (m *MMU) ww(addr uint16, w uint16) error {
	if len(m.mem) == 0 {
		m.init()
	}
	m.mem[addr] = byte(w)
	m.mem[addr+1] = byte(w >> 8)
	return nil
}

type CPU struct {
	Registers

	TClock <-chan int
	MClock <-chan int

	mem MMU

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
