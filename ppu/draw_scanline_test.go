package ppu

import (
	"fmt"
	"testing"

	"github.com/mpingram/gameboy-emu/mmu"
)

func setupTest() (*PPU, *mmu.MMU) {
	m := mmu.New(mmu.MMUOptions{})
	// Set up memory
	// Put two tiles (16bytes) in tile memory
	for i := 0; i < 16; i++ {
		// The tile is uniformly color 3
		m.Mem[mmu.AddrVRAMBlock0+i] = 0xFF
	}
	for i := 0; i < 16; i++ {
		m.Mem[mmu.AddrVRAMBlock0+16+i] = 0x11 + byte(i)
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
	m.Mem[mmu.AddrTileMap0] = 0x10 // 16 byte offset; 1 tile
	// make a diagonal pattern, starting from the top left and going down to the right
	for i := 0; i < 18; i++ {
		m.Mem[mmu.AddrTileMap0+(32*i)+i] = 0x10
	}
	// The BG tile map for all other spaces should already be set to 0
	p := New(m)
	return p, m
}

func Test_drawScanline_renderBG(t *testing.T) {
	p, m := setupTest()
	// Put two tiles (16bytes) in tile memory
	for i := 0; i < 16; i++ {
		// The tile is uniformly color 3
		m.Mem[mmu.AddrVRAMBlock0+i] = 0xFF
	}
	for i := 0; i < 16; i++ {
		m.Mem[mmu.AddrVRAMBlock0+16+i] = 0x11 + byte(i)
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
	m.Mem[mmu.AddrTileMap0] = 0x10 // 16 byte offset; 1 tile
	// make a diagonal pattern, starting from the top left and going down to the right
	for i := 0; i < 18; i++ {
		m.Mem[mmu.AddrTileMap0+(32*i)+i] = 0x10
	}
	// The BG tile map for all other spaces should already be set to 0

	go func() {
		for i := 0; ; i++ {
			m.Mem[mmu.AddrSCY] = byte(i)
			m.Mem[mmu.AddrSCX] = 50
			// run the ppu for 1 screen
			p.RunFor(456 * 154)
		}
	}()
}

// Test tile A:
// _  _  1  3  3  1  _  _
// _  1  3  3  3  3  1  _
// _  3  3  1  1  3  3  _
// -  3  3  _  _  3  3  -
// -  3  1  _  _  1  3  -
// -  3  3  3  3  3  3  -
// _  3  _  _  _  _  3  -
// 3  3  _  _  _  _  3  _
var testTileA = []byte{
	0, 0, 1, 3, 3, 1, 0, 0, // row 0
	0, 1, 3, 3, 3, 3, 1, 0, // row 1
	0, 3, 3, 1, 1, 3, 3, 0, // row 2
	0, 3, 3, 0, 0, 3, 3, 0, // row 3
	0, 3, 1, 0, 0, 1, 3, 0, // row 4
	0, 3, 3, 3, 3, 3, 3, 0, // row 5
	0, 3, 0, 0, 0, 0, 3, 0, // row 6
	3, 3, 0, 0, 0, 0, 3, 0, // row 7
}

// Test tile Y:
// 3  3  _  _  _  3  3  _
// 1  3  _  _  _  3  1  _
// 1  3  1  _  1  3  1  _
// _  1  3  3  3  1  _  _
// _  _  2  3  2  _  _  _
// _  _  2  3  2  _  _  _
// _  _  2  3  2  _  _  _
// _  2  3  3  3  2  _  _
var testTileY = []byte{
	3, 3, 0, 0, 0, 3, 3, 0, // row 0
	1, 3, 0, 0, 0, 3, 1, 0, // row 1
	1, 3, 1, 0, 1, 3, 1, 0, // row 2
	0, 1, 3, 3, 3, 1, 0, 0, // row 3
	0, 0, 2, 3, 2, 0, 0, 0, // row 4
	0, 0, 2, 3, 2, 0, 0, 0, // row 5
	0, 0, 2, 3, 2, 0, 0, 0, // row 6
	0, 2, 3, 3, 3, 2, 0, 0, // row 7
}

// Test tile Checkerboard:
// _  1  _  2  _  1 _  1
// 1  _  2  _  2  _  1  _
// _  1  _  2  _  1  _  1
// 1  _  2  _  2  _  1  _
// _  1  _  2  _  3  _  3
// 1  _  2  _  3  _  3  _
// _  1  _  2  _  3  _  3
// 1  _  2  _  3  _  3  _
var testTileCheckered = []byte{
	0, 2, 0, 2, 0, 2, 0, 2, // row 0
	2, 0, 2, 0, 2, 0, 2, 0, // row 0
	0, 2, 0, 2, 0, 2, 0, 2, // row 0
	2, 0, 2, 0, 2, 0, 2, 0, // row 0
	0, 2, 0, 2, 1, 3, 1, 3, // row 0
	2, 0, 2, 0, 3, 1, 3, 1, // row 0
	0, 2, 0, 2, 1, 3, 1, 3, // row 0
	2, 0, 2, 0, 3, 1, 3, 1, // row 0
}

// packTile formats a background tile for the Gameboy VRAM.
// The input format is an 8x8 array of pallete numbers (0-3).
// The GB represents tiles in a 'packed' 16-byte format,
// aka "2 Bits per pixel" / 2bpp.
// Each row of the sprite corresponds to two bytes:
// the first byte contains the low bits of the 8 palette numbers,
// on the row, and the second byte contains the high bits.
func packTile(tile []byte) []byte {
	packed := make([]byte, 0, 16)
	// for each row of the tile from top (0) to bottom (7)
	for r := 0; r < 8; r++ {
		row := tile[r*8 : (r*8)+8]
		// add the low bit of each palette number to the low byte,
		// and the high bit to the high byte.
		var lo, hi byte
		for i, px := range row {
			hi |= ((px & 0b10) >> 1) << (7 - i)
			lo |= (px & 0b01) << (7 - i)
		}
		packed = append(packed, lo, hi)
	}
	return packed
}

func Test_packTile(t *testing.T) {
	tile := packTile(testTileCheckered)
	tileStr := fmt.Sprintf("% x", tile)
	// confirmed using renderer at https://www.huderlem.com/demos/gameboy2bpp.html
	expectedTileStr := "00 55 00 aa 00 55 00 aa 0f 55 0f aa 0f 55 0f aa"
	if tileStr != expectedTileStr {
		t.Errorf("Expected tile to be %v, got %v", tileStr, expectedTileStr)
	}
}
