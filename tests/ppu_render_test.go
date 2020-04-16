package tests

import (
	"testing"

	frontend "github.com/mpingram/gameboy-emu/frontend/opengl"

	"github.com/mpingram/gameboy-emu/mmu"
	"github.com/mpingram/gameboy-emu/ppu"
)

func setupTest() (*ppu.PPU, *mmu.MMU) {
	m := mmu.New(mmu.MMUOptions{})
	p := ppu.New(m.PPUInterface)
	return p, m
}

func renderScreen(t *testing.T) {

	p, m := setupTest()
	tileOffset := 0x2
	// Put two tiles (16bytes) in tile memory
	for i := 0; i < 16; i++ {
		addr := 0x8000 + (tileOffset * 0x10)
		// addr := 0x8010
		// m.Mem[addr+i] = PackedTestTileA[i]
		m.Mem[addr+i] = PackedTestTileA[i]
	}
	// for i := 0; i < 16; i++ {
	// m.Mem[0x08190+i] = PackedTestTileY[i]
	// }
	// Set a bg palette in which color 0 is white,
	// color 1 is lightgray, color 2 is darkgray, and color 3 is black.
	// https://gbdev.gg8.se/wiki/articles/Video_Display#FF47_-_BGP_-_BG_Palette_Data_.28R.2FW.29_-_Non_CGB_Mode_Only
	// 	Bit 7-6 - Shade for Color Number 3
	// 	Bit 5-4 - Shade for Color Number 2
	// 	Bit 3-2 - Shade for Color Number 1
	// 	Bit 1-0 - Shade for Color Number 0
	//    The four possible gray shades are:
	// 	0  White
	// 	1  Light gray
	// 	2  Dark gray
	// 	3  Black
	m.Mem[mmu.AddrBGP] = 0b11_10_01_00
	// Set Tile addressing method to use the 8000-8FFF range
	m.Mem[mmu.AddrLCDC] |= 0b0001_0000

	// // change the tile mapping at 0,0 to point to tile 1 instead of tile 0.
	// m.Mem[mmu.AddrTileMap0] = 0x10 // 16 byte offset; 1 tile
	// Write 'AYYYY' in the middle of the screen
	y := 10
	x := 10
	m.Mem[mmu.AddrTileMap0+(32*y)+x] = byte(tileOffset)
	// The BG tile map for all other spaces should already be set to 0

	go func() {
		for i := byte(0); ; i-- {
			if i%4 == 0 {
				m.Mem[mmu.AddrSCY] = byte(i / 4)
			}
			// m.Mem[mmu.AddrSCX] = 0
			// run the ppu for 1 screen
			p.RunFor(456 * 154)
		}
	}()
	frontend.ConnectVideo(p.VideoOut)
}
