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
	gameRomFileLocation := os.Args[2]
	gameRom, err := os.Open(gameRomFileLocation)
	if err != nil {
		panic(err)
	}

	videoChannel := make(chan []byte, 1)
	m := mmu.New(mmu.MMUOptions{BootRom: bootRom, GameRom: gameRom})
	p := ppu.New(m.PPUInterface, videoChannel)
	c := cpu.New(m.CPUInterface)

	cpuClock := time.NewTicker(time.Nanosecond)
	defer cpuClock.Stop()
	paused := false
	// cpu goroutine
	go func() {
		breakpoint := uint16(0x50)
		var instr cpu.Instruction
		for {
			<-cpuClock.C
			if paused {
				fmt.Print("> ")
				command := waitForInput()
				if command == "p\n" || command == "print\n" {
					fmt.Println(printCPUState(c))
				} else if command == "m\n" || command == "memdump\n" {
					memdump, err := os.Create("dumps/memdump.bin")
					defer memdump.Close()
					if err != nil {
						panic(err)
					}
					// dump memory to file
					m.Dump(memdump)
				} else if command == "e\n" || command == "exit\n" {
					break
				} else {
					pc := c.PC
					instr, _ = c.Step()
					fmt.Printf("($%04x)\t%s\n", pc, instr.String())
					fmt.Printf("c.PC is now: %04x\n", c.PC)
				}
			} else {
				pc := c.PC
				instr, cycles := c.Step()
				p.RunFor(cycles)
				if pc == breakpoint {
					paused = true
					fmt.Println("HALTED after executing")
					fmt.Printf("($%04x)\t%s\n", pc, instr.String())
				}
			}
		}
	}()

	frontend.ConnectVideo(videoChannel)

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

func printCPUState(c *cpu.CPU) string {
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
