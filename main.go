package main

import (
	"os"

	"github.com/mpingram/gameboy-emu/cpu"
	"github.com/mpingram/gameboy-emu/mmu"
	"github.com/mpingram/gameboy-emu/ppu"
)

func main() {
	mmu := mmu.New(mmu.MMUOptions{})
	ppu := ppu.New(mmu.PPUInterface)
	cpu := cpu.New(mmu.CPUInterface)
	//cpu.SetBreakpoint(0x0100) // this should be just after boot rom
	// this function call takes over the main thread.
	// Should terminate once breakpoint is hit
	// cpu.Run()

	breakpoint := uint16(0x0100)
	for {
		if cpu.PC == breakpoint {
			break
		}
		cpu.Step()
		ppu.DrawScreen()
	}

	memdump, err := os.Create("dumps/memdump.bin")
	defer memdump.Close()
	if err != nil {
		panic(err)
	}
	// dump memory to file
	mmu.Dump(memdump)
}
