package main

import (
	"os"

	"github.com/mpingram/gameboy-emu/cpu"
	"github.com/mpingram/gameboy-emu/mmu"
)

func main() {
	mmu := mmu.New(mmu.MMUOptions{})
	cpu := cpu.New(mmu)
	cpu.SetBreakpoint(0x100) // this should be just after boot rom

	// this function call takes over the main thread.
	// Should terminate once breakpoint is hit
	cpu.Run()

	memdump, err := os.Create("dumps/memdump.bin")
	defer memdump.Close()
	if err != nil {
		panic(err)
	}
	// dump memory to file
	mmu.Dump(memdump)
}
