package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mpingram/gameboy-emu/cpu"
	"github.com/mpingram/gameboy-emu/mmu"
	"github.com/mpingram/gameboy-emu/ppu"
)

func main() {

	bootRomFileLocation := "./roms/boot/DMG_ROM.bin"
	bootRom, err := os.Open(bootRomFileLocation)
	if err != nil {
		panic(err)
	}
	gameRomFileLocation := os.Args[1]
	gameRom, err := os.Open(gameRomFileLocation)
	if err != nil {
		panic(err)
	}

	var _, debug = os.LookupEnv("DEBUG")

	var breakpointEnabled bool
	var breakpoint int64
	if len(os.Args) > 2 {
		breakpointInput := os.Args[2]
		breakpoint, err = strconv.ParseInt(breakpointInput, 0, 0)
		if err != nil {
			if breakpointInput != "" {
				fmt.Printf("ERR: Failed to parse breakpoint: %v\n", breakpointInput)
				return
			}
			breakpointEnabled = false
		} else {
			breakpointEnabled = true
		}
	}

	m := mmu.New(mmu.MMUOptions{BootRom: bootRom, GameRom: gameRom})
	if err != nil {
		panic(err)
	}
	p := ppu.New(m.PPUInterface)
	c := cpu.New(m.CPUInterface)
	if debug {
		fmt.Println(printCPUState(c))
	}

	cpuClock := time.NewTicker(time.Nanosecond)
	defer cpuClock.Stop()
	_, paused := os.LookupEnv("PAUSED")
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
			} else if command == "c\n" || command == "continue\n" {
				paused = false
			} else if command == "q\n" || command == "quit\n" {
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
			if breakpointEnabled && int64(pc) == breakpoint {
				paused = true
				fmt.Println("HALTED after executing")
				fmt.Printf("($%04x)\t%s\n", pc, instr.String())
			}
		}
	}

	// frontend.ConnectVideo(p.VideoOut)

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
