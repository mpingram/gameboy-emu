package mmu

import "math/rand"

func New() *MMU {
	mmu := &MMU{}
	mmu.init()
	return mmu
}

type MMU struct {
	mem []byte
}

func (m *MMU) init() {
	m.mem = make([]byte, 0xFFFF)
	// initialize to randomized values
	rand.Read(m.mem)
}

func (m *MMU) Rb(addr uint16) (byte, error) {
	if len(m.mem) == 0 {
		m.init()
	}
	return m.mem[addr], nil
}
func (m *MMU) Wb(addr uint16, b byte) error {
	if len(m.mem) == 0 {
		m.init()
	}
	m.mem[addr] = b
	return nil
}

func (m *MMU) Rw(addr uint16) (uint16, error) {
	if len(m.mem) == 0 {
		m.init()
	}
	hi := m.mem[addr]
	lo := m.mem[addr+1]
	return uint16(hi)<<8 | uint16(lo), nil
}

func (m *MMU) Ww(addr uint16, w uint16) error {
	if len(m.mem) == 0 {
		m.init()
	}
	hi := byte(w >> 8)
	lo := byte(w)
	m.mem[addr] = hi
	m.mem[addr+1] = lo
	return nil
}
