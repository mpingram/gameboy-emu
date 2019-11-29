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
	//p := ppu.New(m.PPUInterface)
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
			//screen = p.DrawScreen()
			screen = make([]byte, 0)
			for i := 0; i < 144*160; i++ {
				screen = append(screen, byte(time.Now().Second()%255), 0, byte(i%255))
			}
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
