package mmu

type cpuMemoryInterface struct {
	mmu *MMU
}

func (cmi *cpuMemoryInterface) Rb(addr uint16) (byte, error) {
	m := cmi.mmu
	return m.mem[addr], nil
}
func (cmi *cpuMemoryInterface) Wb(addr uint16, b byte) error {
	m := cmi.mmu
	m.mem[addr] = b
	return nil
}

func (cmi *cpuMemoryInterface) Rw(addr uint16) (uint16, error) {
	m := cmi.mmu
	hi := m.mem[addr]
	lo := m.mem[addr+1]
	return uint16(hi)<<8 | uint16(lo), nil
}

func (cmi *cpuMemoryInterface) Ww(addr uint16, w uint16) error {
	m := cmi.mmu
	hi := byte(w >> 8)
	lo := byte(w)
	m.mem[addr] = hi
	m.mem[addr+1] = lo
	return nil
}