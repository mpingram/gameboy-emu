package main

import (
	"os"
	"fmt"
	"github.com/mpingram/gameboy-emu/cpu"
	"github.com/mpingram/gameboy-emu/mmu"
	"github.com/mpingram/gameboy-emu/ppu"
)

func main() {

	const bootRomFileLocation = "./roms/boot/DMG_ROM.bin"
	bootRom, err := os.Open(bootRomFileLocation)
	if err != nil {
		panic(err)
	}
	m := mmu.New(mmu.MMUOptions{BootRom: bootRom})
	p := ppu.New(m.PPUInterface)
	c := cpu.New(m.CPUInterface)
	//cpu.SetBreakpoint(0x0100) // this should be just after boot rom
	// this function call takes over the main thread.
	// Should terminate once breakpoint is hit
	// cpu.Run()

	breakpoint := uint16(0x08e)
	var screen ppu.Screen
	for {
		if c.PC == breakpoint {
			break
		}
		c.Step()
		screen = p.DrawScreen()
	}

	fmt.Printf("%v", screen)

	memdump, err := os.Create("dumps/memdump.bin")
	defer memdump.Close()
	if err != nil {
		panic(err)
	}
	// dump memory to file
	m.Dump(memdump)
}
