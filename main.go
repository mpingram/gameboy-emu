package main

import (
	"bufio"
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

	cpuClock := time.NewTicker(time.Nanosecond)
	defer cpuClock.Stop()
	paused := false
	// cpu goroutine
	go func() {
		breakpoint := uint16(0x0000)
		var instr cpu.Instruction
		for {
			<-cpuClock.C
			if paused {
				fmt.Print("> ")
				command := waitForInput()
				if command == "d\n" || command == "dump\n" {
					fmt.Println(dumpCPUState(c))
				} else {
					pc := c.PC
					instr, _ = c.Step()
					fmt.Printf("($%04x)\t%s\n", pc, instr.String())
					fmt.Printf("c.PC is now: %04x\n", c.PC)
				}
			} else {
				pc := c.PC
				instr, _ = c.Step()
				if pc == breakpoint {
					paused = true
					fmt.Println("HALTED after executing")
					fmt.Printf("($%04x)\t%s\n", pc, instr.String())
				}
			}
		}
	}()
	// ppu goroutine
	ppuClock := time.NewTicker(time.Nanosecond)
	defer ppuClock.Stop()
	go func() {
		for {
			<-ppuClock.C // FIXME we'll need to find a different way to
			// coordinate the cpu and ppu timings.
			// screen := placeholderScreen()
			videoChannel <- p.DrawScreen()
			time.Sleep(time.Microsecond)
		}
	}()

	frontend.ConnectVideo(videoChannel)

	memdump, err := os.Create("dumps/memdump.bin")
	defer memdump.Close()
	if err != nil {
		panic(err)
	}
	// dump memory to file
	m.Dump(memdump)
}

func waitForInput() string {
	// block until the user types anything
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return text
}

func dumpCPUState(c *cpu.CPU) string {
	return fmt.Sprintf(`
	===== CPU =====
	PC: %0x \t SP: %0x
	A: %0x \t F: %0x
	B: %0x \t C: %0x
	D: %0x \t E: %0x
	H: %0x \t L: %0x
	F: %08b
	`, c.PC, c.SP, c.A, c.F, c.B, c.C, c.D, c.E, c.H, c.L, c.F)
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
