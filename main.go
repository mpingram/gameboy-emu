package main

import (
	"os"

	"github.com/mpingram/gameboy-emu/mmu"
)

func main() {
	mmu := mmu.New(mmu.MMUOptions{})

	memdump, err := os.Create("dumps/memdump.bin")
	if err != nil {
		panic(err)
	}
	// dump memory to file
	mmu.Dump(memdump)
}
