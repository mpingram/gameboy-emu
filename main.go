package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mpingram/gameboy-emu/cpu"
	frontend "github.com/mpingram/gameboy-emu/frontend/opengl"
	"github.com/mpingram/gameboy-emu/mmu"
	"github.com/mpingram/gameboy-emu/ppu"
)

func main() {

	bootRomFileLocation := os.Args[1]
	bootRom, err := os.Open(bootRomFileLocation)
	if err != nil {
		panic(err)
	}
	m := mmu.New(mmu.MMUOptions{BootRom: bootRom})
	p := ppu.New(m.PPUInterface)
	c := cpu.New(m.CPUInterface)

	videoChannel := make(chan []byte, 1)
	go func() {
		breakpoint := uint16(0x08e)
		var screen ppu.Screen
		for {
			if c.PC == breakpoint {
				break
			}
			c.Step()
			screen = p.DrawScreen()
			//screen = placeholderScreen()
			videoChannel <- screen
			time.Sleep(100 * time.Millisecond)
		}
	}()
	frontend.ConnectVideo(videoChannel)
	fmt.Println("?")

	memdump, err := os.Create("dumps/memdump.bin")
	defer memdump.Close()
	if err != nil {
		panic(err)
	}
	// dump memory to file
	m.Dump(memdump)
}

func placeholderScreen() []byte {
	screen := make([]byte, 0)
	var row, col byte
	for row = 0; row < 144; row++ {
		for col = 0; col < 160; col++ {
			var r, g, b byte
			r = row
			g = col
			b = byte(time.Now().Second() % 255)
			screen = append(screen, r, g, b)
		}
	}
	return screen
}
