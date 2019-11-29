package main

import (
	"fmt"
	"os"

	"github.com/mpingram/gameboy-emu/cpu"
	frontend "github.com/mpingram/gameboy-emu/frontend/opengl"
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

	videoChannel := make(chan []byte, 1)
	frontend.ConnectVideo(videoChannel)

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
