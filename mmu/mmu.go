package mmu

import (
	"fmt"
	"io"
)

const (
	AddrCartRomBank00         = 0x0000
	AddrCartRomSwitchableBank = 0x4000
	AddrVRAM                  = 0x8000

	AddrVRAMBlock0 = 0x8000
	AddrVRAMBlock1 = 0x8800
	AddrVRAMBlock2 = 0x9000
	AddrTileMap0   = 0x9800
	AddrTileMap1   = 0x9C00

	AddrCartRAM               = 0xA000
	AddrWorkRAMBank0          = 0xC000
	AddrWorkRAMSwitchableBank = 0xD000
	AddrEchoRAM               = 0xE000
	AddrOamRAM                = 0xFE00
	AddrIORegs                = 0xFF00

	AddrLCDC    = 0xFF40
	AddrLCDStat = 0xFF41
	AddrSCY     = 0xFF42
	AddrSCX     = 0xFF43
	AddrLY      = 0xFF44
	AddrLYC     = 0xFF45
	AddrDMA     = 0xFF46
	AddrBGP     = 0xFF47
	AddrOBP0    = 0xFF48
	AddrOBP1    = 0xFF49
	AddrWY      = 0xFF4A
	AddrWX      = 0xFF4B

	AddrHighRAM            = 0xFF80
	AddrInterruptEnableReg = 0xFFFF
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
	Mem          []byte
	CPUInterface *cpuMemoryInterface
	PPUInterface *ppuMemoryInterface
	gameRom      io.Reader
}

func (m *MMU) init() {
	// create backing memory array
	m.Mem = make([]byte, 0x10000)
	// zero out all memory (FIXME for true emulation, this should actually randomize all(?) values)
	for i := 0; i < 0x10000; i++ {
		m.Mem[i] = 0x00
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
		gameRomMemory := m.Mem[0x0000:0x07FF]
		_, err := m.gameRom.Read(gameRomMemory)
		if err != nil {
			panic(err)
		}
		m.gameRom = nil
	}
}

func (m *MMU) loadBootRom(rom io.Reader) {
	bootRomMemory := m.Mem[0x0000:0x0100]
	_, err := rom.Read(bootRomMemory)
	if err != nil {
		panic(err)
	}
}

func (m *MMU) Dump(out io.Writer) {
	_, err := out.Write(m.Mem)
	if err != nil {
		panic(err)
	}
}

func (m *MMU) rb(addr uint16) byte {
	return m.Mem[addr]
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
		m.Mem[addr] = b
	}
}

func (m *MMU) rw(addr uint16) uint16 {
	hi := m.Mem[addr]
	lo := m.Mem[addr+1]
	return uint16(hi)<<8 | uint16(lo)
}

func (m *MMU) ww(addr uint16, w uint16) {
	hi := byte(w >> 8)
	lo := byte(w)
	m.Mem[addr] = hi
	m.Mem[addr+1] = lo
}
