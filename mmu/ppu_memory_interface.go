package mmu

type ppuMemoryInterface struct {
	mmu *MMU
}

func (pmi *ppuMemoryInterface) Rb(addr uint16) (byte, error) {
	m := pmi.mmu
	return m.mem[addr], nil
}
func (pmi *ppuMemoryInterface) Wb(addr uint16, b byte) error {
	m := pmi.mmu
	m.mem[addr] = b
	return nil
}

func (pmi *ppuMemoryInterface) Rw(addr uint16) (uint16, error) {
	m := pmi.mmu
	hi := m.mem[addr]
	lo := m.mem[addr+1]
	return uint16(hi)<<8 | uint16(lo), nil
}

func (pmi *ppuMemoryInterface) Ww(addr uint16, w uint16) error {
	m := pmi.mmu
	hi := byte(w >> 8)
	lo := byte(w)
	m.mem[addr] = hi
	m.mem[addr+1] = lo
	return nil
}
