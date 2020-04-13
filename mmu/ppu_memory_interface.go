package mmu

type ppuMemoryInterface struct {
	mmu *MMU
}

func (pmi *ppuMemoryInterface) Rb(addr uint16) byte {
	m := pmi.mmu
	return m.Mem[addr]
}
func (pmi *ppuMemoryInterface) Wb(addr uint16, b byte) {
	m := pmi.mmu
	m.Mem[addr] = b
}

func (pmi *ppuMemoryInterface) Rw(addr uint16) uint16 {
	m := pmi.mmu
	hi := m.Mem[addr]
	lo := m.Mem[addr+1]
	return uint16(hi)<<8 | uint16(lo)
}

func (pmi *ppuMemoryInterface) Ww(addr uint16, w uint16) {
	m := pmi.mmu
	hi := byte(w >> 8)
	lo := byte(w)
	m.Mem[addr] = hi
	m.Mem[addr+1] = lo
}
