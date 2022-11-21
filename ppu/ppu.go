package ppu

import (
	"fmt"

	"github.com/mpingram/gameboy-emu/mmu"
)

const screenHeight = 144
const screenWidth = 160

type PPU struct {
	mem      mmu.MemoryReadWriter
	cycles   int
	screen   []byte
	VideoOut chan []byte

	scanline []byte
	tileData []byte
}

func New(mem mmu.MemoryReadWriter) *PPU {
	videoOut := make(chan []byte, 1) // videoOut channel is buffered by one screen
	scanline := make([]byte, screenWidth)
	tileData := make([]byte, 8)
	ppu := &PPU{mem, 0, []byte{}, videoOut, scanline, tileData}
	ppu.setMode(OAMSearch)
	return ppu
}

// LCDControl represents a memory register located at 0xFF4
// which is used to configure the behavior of the PPU while the Gameboy is running.
// See https://gbdev.gg8.se/wiki/articles/LCDC.
// LCDEnable controls whether or not the screen and PPU are turned on. (0=off, 1=on)
// WindowTileMapSelect switches the region of memory that the PPU reads
//   the Window's tile map from. (0=$9800-$9BFF, 1=$9C00-$9FFF).
// WindowEnable controls whether or not the Window is displayed. (0=off, 1=on)
// TileAddressingMode changes the way that the PPU determines the memory addresses of Background and Window tiles
//   from the tile address offsets provided by the Background and Window tile maps.
//   If 0, tiles are addressed by $8000 + byte(offset); if 1, tiles are addressed by $8800 +/- int8(offset).
//   For a more detailed explanation, see the comments on the `getTileRow` function below.
//   See https://gbdev.gg8.se/wiki/articles/Video_Display#VRAM_Tile_Data
// BGTileMapSelect switches the region of memory the PPU reads the tile
//   map of the background from. (0=9800-9BFF, 1=9C00-9FFF)
// SpriteSize toggles whether all sprites are one tile in size or two tiles (0=8x8px, 1=8x16px)
// SpriteEnable toggles whether all sprites are rendered or not.
// WindowDisplayORPriority
// See https://gbdev.gg8.se/wiki/articles/LCDC#LCDC.0_-_BG.2FWindow_Display.2FPriority

const LCDCAddr = 0xFF40

const LCDEnable = 0b1000_0000
const WindowTileMapSelect = 0b0100_0000
const WindowEnable = 0b0010_0000
const TileAddressingMode = 0b0001_0000
const BGTileMapSelect = 0b0000_1000
const SpriteSize = 0b0000_0100
const SpriteEnable = 0b0000_0010
const WindowDisplayORPriority = 0b0000_0001 // bit 0

// LCDStat is a memory register which is
// used to enable some interrupts related to drawing the LCD screen.
const LCDStatAddr = 0xFF41

const LYCoincidenceInterruptEnable = 0b0100_0000 // bit 6
const OAMInterruptEnable = 0b0010_0000           // bit 5
const VBlankInterruptEnable = 0b0001_0000        // bit 4
const HBlankInterruptEnalbe = 0b0000_1000        // bit 3
const LYCoincidenceStatus = 0b0000_0100          // bit 2 (0: LYC<>LY, 1: LYC=LY)
// Mode represents the 'drawing mode' of the Gameboy,
// which is stored in bits 0 and 1 of the LCDStat memory register.
// The drawing mode cycles between OAMSearch, PixelDrawing, and HBlank
// for each scanline, and enters VBlank once a full screen has been drawn.
const LCDMode = 0b0000_0011 // bits 0,1
type Mode byte

const (
	// OAMSearch (mode 2) is the mode in which the PPU searches through the Object Attribute Memory
	// for active sprites on this scanline.
	OAMSearch Mode = 2
	// PixelDrawing (mode 3) represents the mode in which the PPU draws a scanline to the LCD screen,
	PixelDrawing = 3
	// HBlank (Horizontal Blank) (mode 0) represents the mode in which the PPU waits [FIXME: how many?] clocks after drawing a scanline.
	// Then, after drawing the full screen (30 scanlines), the PPU enters a 4th mode:
	HBlank = 0
	// VBlank (Vertical Blank) (mode 1) represents the mode in which the PPU waits [FIXME: some number] clocks after drawing the screen.
	VBlank = 1
)

// setMode sets the LCDStat's Mode register by writing to the MMU.
// This has a side effect in the MMU: In some modes, some regions of memory
// are locked to either the CPU or the PPU.
func (p *PPU) setMode(mode Mode) {
	// set the bottom two bits of the lcdstat byte equal to the binary
	// representation of the mode
	b := p.mem.Rb(LCDStatAddr)
	// clear the bottom two bits
	b &= 0b1111_1100
	// set the bottom two bits
	b |= byte(mode)
	p.mem.Wb(LCDStatAddr, b)
}

// getBackgroundTileRow returns the 8 pixels of a row of a background tile located at
// coordinate x, y. Note that it does NOT know the palette color for these pixels.
func (p *PPU) getBackgroundTileRow(x, y byte) []byte {
	// func (p *PPU) getBackgroundTileRow(x, y byte) []pixelData {
	// check lcdc to see where bg tile map is stored
	lcdc := p.mem.Rb(LCDCAddr)
	var tileMapLocation uint16
	if lcdc&BGTileMapSelect == 0 {
		tileMapLocation = 0x9800
	} else {
		tileMapLocation = 0x9C00
	}
	// calculate byte offset in bg tile map based on x,y.
	offset := (uint16(y)/8)*32 + (uint16(x) / 8)
	// Explanation:
	// Tile memory is laid out like this:
	// $9BFF/$9FFF +-------------------+
	//             |992|993|...   |1023|
	//             |                   |
	// 32 tile ptrs|                   |
	//         |   |32|33|34|...    |63|
	//         y   |0 |1 |2 |...    |31|
	// $9800/$9C00 +--32 tile ptrs-----+
	// 					x ----->
	// Where each tile represents a 8x8px area. So the formula for
	// getting the tile byte offset is floor(y/8)*32 + floor(x/8).
	// Wraparound for x and y is handled by byte overflow.

	// Read the correct byte of the tile map to get the address of the tile data
	// (Remember that the address of the tile data is an offset, not a full uint16 addresss.)
	tileAddrOffset := p.mem.Rb(tileMapLocation + offset)
	var tileAddr uint16
	// LCDC tileAddressingMode 0 = $8800 1 = $8000
	if lcdc&TileAddressingMode == 0 {
		// convert addrOffset to a signed byte
		signedAddrOffset := int8(tileAddrOffset)
		// NOTE this is potentially buggy!
		// promote the signed int8 to int in order to add it to 0x8800,
		// then convert the result back to uint16.
		tileAddr = uint16(0x8800 + int(signedAddrOffset)*0x10)
	} else {
		tileAddr = 0x8000 + uint16(tileAddrOffset)*0x10
	}
	// The row of the tile that intersects with this y-coordinate. Rows go from 0-7,
	// where 7 is the bottom row.
	row := y % 8
	return p.getTileRowData(tileAddr, row)

	// tileData := p.getTileRowData(tileAddr, row)
	// pixels := make([]pixelData, 0)
	// for _, paletteIndex := range tileData {
	// 	px := pixelData{paletteIndex, bg}
	// 	pixels = append(pixels, px)
	// }
	// return pixels
}

// getTileRowData returns the color numbers, from left to right, of a certain row of a tile located
// at a 16-byte region in video memory.
// Tile rows are 0-indexed and run from top to bottom, so the bottom row of a tile is row 7.
// FIXME almost certainly a bug here!
// FIXME introduced a bug here :>
func (p *PPU) getTileRowData(tileAddr uint16, row byte) []byte {
	if row > 7 {
		panic(fmt.Sprintf("Got tile row > 7: %v", row))
	}
	// Each 2 bytes of the tile is a row of the tile, and they are stored from top
	// to bottom. The bytes that represents row n of the tile are at (tileAddr + n*2)
	b1 := p.mem.Rb(tileAddr + uint16(row*2))
	b2 := p.mem.Rb(tileAddr + uint16(row*2) + 1)
	i := 0
	for j := 7; j >= 0; j-- {
		// WARNING Possibly buggy
		// https://gbdev.gg8.se/wiki/articles/Video_Display#VRAM_Tile_Data
		// b1 contains the low bit of each px, from left (bit 7) to right (bit 0)
		// b2 contains the high bit of each px, as above.
		mask := byte(1) << j
		lo := (b1 & mask) >> j
		hi := (b2 & mask) >> j
		color := (hi << 1) | lo
		p.tileData[i] = color
		i += 1
	}
	return p.tileData
}

const (
	White     byte = 0
	LightGray      = 1
	DarkGray       = 2
	Black          = 3
)

type pixelData struct {
	colorNumber   byte
	paletteNumber paletteNumber
}

type paletteNumber byte

const (
	bg   paletteNumber = 2
	obj0 paletteNumber = 0
	obj1               = 1
)

// getYScroll gets the y-coordinate of the top-left of the LCD screen.
// Reads the SCY ($FF42) memory register.
func (p *PPU) getScrollY() byte {
	var scyAddr uint16 = 0xFF42
	return p.mem.Rb(scyAddr)
}

// getXScroll gets the y-coordinate of the top-left of the LCD screen.
// Reads the SCX ($FF43) memory register.
func (p *PPU) getScrollX() byte {
	var scxAddr uint16 = 0xff43
	return p.mem.Rb(scxAddr)
}

// getLY gets the current scanline. The max 'scanline'
// is 153 and not 143 (the LCD screen only has 144 scanlines) because scanlines
// 144-153 represent the time spent in VBlank mode.
// Reads the LY ($FF44) memory register.
func (p *PPU) getLY() byte {
	var lyAddr uint16 = 0xff44
	return p.mem.Rb(lyAddr)
}

// setLY sets the current scanline. The max 'scanline'
// is 153 and not 143 (the LCD screen only has 144 scanlines) because scanlines
// 144-153 represent the time spent in VBlank mode.
// Reads the LY ($FF44) memory register.
func (p *PPU) setLY(ly byte) {
	var lyAddr uint16 = 0xff44
	p.mem.Wb(lyAddr, ly)
}

// getLYCompare gets the value of a register that is used to trigger
// an interrupt on a specific scanline.
// When a new scanline is started, the MMU compares the value of the LYC register
// with the current scanline (the LY register). If the values match and the
// LYCoincidenceInterruptEnabled bit is set in the LCDStat register, then the
// LYCoincidenceInterrupt is triggered and program execution jumps to that routine.
// Reads the LYC (LYCompare) ($FF45) memory register.
func (p *PPU) getLYCompare() byte {
	var lycAddr uint16 = 0xff45
	return p.mem.Rb(lycAddr)
}

// Gets the X coordinate of the Window top left, minus 7.
// Reads WindowX-7($FF4B) memory register.
func (p *PPU) getWindowX() byte {
	var wxAddr uint16 = 0xff4a
	return p.mem.Rb(wxAddr)
}

// Gets the Y coordinate of the Window top left.
// Reads the WindowY($FF4A) memory registers.
func (p *PPU) getWindowY() byte {
	var wyAddr uint16 = 0xff4a
	return p.mem.Rb(wyAddr)
}

// Screen is a byte array representing the colorized pixels
// of a gameboy screen. Its format is
//
//	 1 pixel
//	|-----|
//
// [R, G, B, R, G, B, R, G, B]
// Where R,G,B are one byte representing the red, green, blue
// component of each pixel.
type Screen []byte
