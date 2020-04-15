package mmu

type cpuMemoryInterface struct {
	mmu *MMU
}

func (cmi *cpuMemoryInterface) Rb(addr uint16) byte {
	return cmi.mmu.rb(addr)
}

func (cmi *cpuMemoryInterface) Wb(addr uint16, b byte) {
	cmi.mmu.wb(addr, b)
}

func (cmi *cpuMemoryInterface) Rw(addr uint16) uint16 {
	return cmi.mmu.rw(addr)
}

func (cmi *cpuMemoryInterface) Ww(addr uint16, w uint16) {
	cmi.mmu.ww(addr, w)
}
