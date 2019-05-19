package mmu

import "math/rand"

// NOTE -- MMU needs to know which areas of memory can't be read by CPU, vs which can
// be read by PPU, *at different times depending on PPU state*
// (ie, OAM RAM is inaccessible to CPU during OAM search / Pixel transfer modes, and VRAM is
// inaccessible to CPU during Pixel transfer mode).
//
// The current design does not take that into account -- we could dynamically set which
// areas of memory are writable / readable, but MMU doesn't know about CPU/PPU and would
// have no way of knowing which reads are coming from which source. (ie, ppu could say 'lock this
// area of memory', but how would we know that future reads come from PPU? Wait... that's what a
// lock is. Holy- ok that might be an area we can research)
//
// Need to consider a change in design from a 'dumb' mmu? One that's aware of the state
// of the PPU (through LCDStatus register, also located in memory...) ... or, can have
// CPU self-police its reads/writes to memory during OAM search / pixel transfer modes.
// Or, can implement locks.
//
// Could have mmu.AcquireLock(addr, size) (mmu.Lock, err)
//

type MemoryReadWriter interface {
	MemoryReader
	MemoryWriter
}

type MemoryReader interface {
	Rb(addr uint16) (byte, error)
	Rw(addr uint16) (uint16, error)
}

type MemoryWriter interface {
	Wb(addr uint16, b byte) error
	Ww(addr uint16, bb uint16) error
}

func New() *MMU {
	mmu := &MMU{}
	mmu.init()
	return mmu
}

type MMU struct {
	mem          []byte
	CPUInterface MemoryReadWriter
	PPUInterface MemoryReadWriter
}

func (m *MMU) init() {
	// initialize backing memory array
	m.mem = make([]byte, 0x10000)
	// initialize memory to randomized values
	rand.Read(m.mem)

	// TODO add boot ROM, other initialization things
	// TODO memory bank switching

	// Set up CPU interface and PPU interface
	m.CPUInterface = &cpuMemoryInterface{mmu: m}
	m.PPUInterface = &ppuMemoryInterface{mmu: m}
}
