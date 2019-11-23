package mmu

import (
	"io"
	"os"
)

const bootRomFileLocation = "./roms/boot/DMG_ROM.bin"

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

type MMUOptions struct {
	GameRom io.Reader
}

func New(opt MMUOptions) *MMU {
	mmu := &MMU{}
	mmu.init()
	if opt.GameRom != nil {
		mmu.loadGameRom(opt.GameRom)
	}
	return mmu
}

type MMU struct {
	mem          []byte
	CPUInterface MemoryReadWriter
	PPUInterface MemoryReadWriter
}

func (m *MMU) init() {
	// create backing memory array
	m.mem = make([]byte, 0x10000)
	// zero out all memory (FIXME for true emulation, this should actually randomize all(?) values)
	for i := 0; i < 0x10000; i++ {
		m.mem[i] = 0x00
	}

	// load boot rom into memory at $0000-$0100
	// DEBUG -- for now, read boot rom from hardcoded file location
	bootRom, err := os.Open(bootRomFileLocation)
	if err != nil {
		panic(err)
	}
	bootRomMemory := m.mem[0x000:0x0100]
	_, err = bootRom.Read(bootRomMemory)
	if err != nil {
		panic(err)
	}

	// TODO memory bank switching

	// Set up CPU interface and PPU interface
	m.CPUInterface = &cpuMemoryInterface{mmu: m}
	m.PPUInterface = &ppuMemoryInterface{mmu: m}
}

func (m *MMU) loadGameRom(io.Reader) {
	// do nothing
}

func (m *MMU) Dump(out io.Writer) {
	_, err := out.Write(m.mem)
	if err != nil {
		panic(err)
	}
}
