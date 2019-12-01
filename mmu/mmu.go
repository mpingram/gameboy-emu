package mmu

import (
	"fmt"
	"io"
)

type MMUOptions struct {
	GameRom io.Reader
	BootRom io.Reader
}

func New(opt MMUOptions) *MMU {
	mmu := &MMU{}
	mmu.init()
	if opt.BootRom != nil {
		mmu.loadBootRom(opt.BootRom)
	}
	if opt.GameRom != nil {
		mmu.gameRom = opt.GameRom
	}
	return mmu
}

type MMU struct {
	mem          []byte
	CPUInterface *cpuMemoryInterface
	PPUInterface *ppuMemoryInterface
	gameRom      io.Reader
}

func (m *MMU) init() {
	// create backing memory array
	m.mem = make([]byte, 0x10000)
	// zero out all memory (FIXME for true emulation, this should actually randomize all(?) values)
	for i := 0; i < 0x10000; i++ {
		m.mem[i] = 0x00
	}

	// TODO memory bank switching

	// Set up CPU interface and PPU interface
	m.CPUInterface = &cpuMemoryInterface{mmu: m}
	m.PPUInterface = &ppuMemoryInterface{mmu: m}
}

// loadGameRom is called after the boot rom executes. The
// loading is triggered by a write to $FF50
func (m *MMU) loadGameRom() {
	if m.gameRom != nil {
		// TODO implement bank switching
		gameRomMemory := m.mem[0x0000:0x07FF]
		_, err := m.gameRom.Read(gameRomMemory)
		if err != nil {
			panic(err)
		}
		m.gameRom = nil
	}
}

func (m *MMU) loadBootRom(rom io.Reader) {
	bootRomMemory := m.mem[0x0000:0x0100]
	_, err := rom.Read(bootRomMemory)
	if err != nil {
		panic(err)
	}
}

func (m *MMU) Dump(out io.Writer) {
	_, err := out.Write(m.mem)
	if err != nil {
		panic(err)
	}
}

func (m *MMU) rb(addr uint16) byte {
	return m.mem[addr]
}

func (m *MMU) wb(addr uint16, b byte) {
	// Handle memory mapped registers
	switch addr {
	case 0xff50: // writing 0x1 to $ff50 unmaps the boot ROM from memory.
		fmt.Print("Wrote to $ff50")
		if b == 0x1 {
			fmt.Print("Wrote 0x1 to $FF50")
			m.loadGameRom()
		}
	default:
		m.mem[addr] = b
	}
}

func (m *MMU) rw(addr uint16) uint16 {
	hi := m.mem[addr]
	lo := m.mem[addr+1]
	return uint16(hi)<<8 | uint16(lo)
}

func (m *MMU) ww(addr uint16, w uint16) {
	hi := byte(w >> 8)
	lo := byte(w)
	m.mem[addr] = hi
	m.mem[addr+1] = lo
}
