package main

import (
	"testing"

	frontend "github.com/mpingram/gameboy-emu/frontend/opengl"

	"github.com/mpingram/gameboy-emu/mmu"
	"github.com/mpingram/gameboy-emu/ppu"
)

var (
	p *ppu.PPU
	m *mmu.MMU
)

func setupTest() (*ppu.PPU, *mmu.MMU) {
	m := mmu.New(mmu.MMUOptions{})
	p := ppu.New(m.PPUInterface)
	return p, m
}

func Test_renderScreen(t *testing.T) {
	p, m := setupTest()
	// Put two tiles (16bytes) in tile memory
	for i := 0; i < 16; i++ {
		// The tile is uniformly color 3
		m.Mem[mmu.AddrVRAMBlock0+i] = 0xFF
	}
	for i := 0; i < 16; i++ {
		m.Mem[mmu.AddrVRAMBlock0+16+i] = 0x01
	}
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

	// change the tile mapping at 0,0 to point to tile 1 instead of tile 0.
	m.Mem[mmu.AddrTileMap0] = 0b0000_0001
	// The BG tile map for all other spaces should already be set to 0, which
	// should now indicate this tile we've made

	// run the ppu for 60 screens
	go p.RunFor(456 * 154 * 60)
	frontend.ConnectVideo(p.VideoOut)
	// render the screen the ppu creates
	// screen := <-p.VideoOut
	// frontend.Render(screen)
}
