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

type MMU struct {
	Mem          []byte
	CPUInterface *cpuMemoryInterface
	PPUInterface *ppuMemoryInterface
	gameRom      []byte
	bootRom      []byte
	mapBootRom   bool
}

type MMUOptions struct {
	GameRom io.Reader
	BootRom io.Reader
}

func New(opt MMUOptions) *MMU {
	m := &MMU{}
	m.CPUInterface = &cpuMemoryInterface{mmu: m}
	m.PPUInterface = &ppuMemoryInterface{mmu: m}

	m.Mem = make([]byte, 0x10000)
	if opt.BootRom != nil {
		// The boot ROM is 0x100 bytes long and is mapped to 0x0-0x100 at boot. When the boot
		// ROM finishes, it writes to the register 0xFF50, which unmaps the boot rom from memory.
		// At this point 0x0000-0x0100 becomes mapped to the start of game rom bank 0.
		m.bootRom = make([]byte, 0x0100)
		_, err := opt.BootRom.Read(m.bootRom)
		if err != nil {
			panic(fmt.Errorf("ERR reading boot ROM: %v", err))
		}
		m.mapBootRom = true
	}
	if opt.GameRom != nil {
		// The game ROM can contain up to 125 switchable 16kb rom banks in addition to bank 0.
		// The banks are addressed 0x01-0x7F, although bank numbers 0x20,0x40,0x60 cannot be used.
		m.gameRom = make([]byte, 0x80*0x4000) // Space for 0x80=128 16kb banks.
		_, err := opt.GameRom.Read(m.gameRom)
		if err != nil {
			panic(fmt.Errorf("ERR reading game ROM: %v", err))
		}
		// load game ROM banks 0 and 1 into 0x0000-0x7FFF
		for i := 0x0; i < 0x8000; i++ {
			m.Mem[i] = m.gameRom[i]
		}
	}
	return m
}

func (m *MMU) Dump(out io.Writer) {
	_, err := out.Write(m.Mem)
	if err != nil {
		panic(err)
	}
}

func (m *MMU) rb(addr uint16) byte {
	switch {
	case addr < 0x0100:
		if m.mapBootRom {
			return m.bootRom[addr]
		}
		return m.Mem[addr]
	default:
		return m.Mem[addr]
	}
}

func (m *MMU) wb(addr uint16, b byte) {
	// TODO Handle memory mapped registers
	switch addr {
	case 0xff50: // writing 0x1 to $ff50 unmaps the boot ROM from memory.
		if b == 0x1 {
			m.mapBootRom = false
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
