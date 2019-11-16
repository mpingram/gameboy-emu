package ppu

import (
	"github.com/mpingram/gameboy-emu/mmu"
)

type PPU struct {
	TClock <-chan int
	MClock <-chan int

	mem mmu.MemoryReadWriter
}

const screenHeight byte = 144
const screenWidth byte = 160

// Mode represents the 'drawing mode' of the Gameboy,
// which is stored in bits 0 and 1 of the LCDStat memory register.
// The drawing mode cycles between OAMSearch, PixelDrawing, and HBlank
// for each scanline, and enters VBlank once a full screen has been drawn.
type Mode int

const (
	// OAMSearch is the mode in which the PPU searches through the Object Attribute Memory
	// for active sprites on this scanline.
	OAMSearch Mode = 2
	// PixelDrawing represents the mode in which the PPU draws a scanline to the LCD screen,
	PixelDrawing = 3
	// HBlank (Horizontal Blank) represents the mode in which the PPU waits [FIXME: how many?] clocks after drawing a scanline.
	// Then, after drawing the full screen (30 scanlines), the PPU enters a 4th mode:
	HBlank = 0
	// VBlank (Vertical Blank) represents the mode in which the PPU waits [FIXME: some number] clocks after drawing the screen.
	VBlank = 1
)

// setMode sets the LCDStat's Mode register by writing to the MMU.
// This has a side effect in the MMU: In some modes, some regions of memory
// are locked to either the CPU or the PPU.
func (p *PPU) setMode(mode Mode) {
	// FIXME implement - set bits 1 and 0 of LCDSTAT register
}

// LCDControl represents a memory register located at [FIXME address]
// which is used to configure the behavior of the PPU while the Gameboy is running.
// See https://gbdev.gg8.se/wiki/articles/LCDC.
type LCDControl struct {
	// LCDEnable controls whether or not the screen and PPU are turned on. (0=off, 1=on)
	LCDEnable bool // bit 7
	// WindowTileMapSelect switches the region of memory that the PPU reads
	// the Window's tile map from. (0=$9800-$9BFF, 1=$9C00-$9FFF).
	WindowTileMapSelect bool // bit 6 (0=$9800-$9BFF, 1=$9C00-$9FFF)
	// WindowEnable controls whether or not the Window is displayed. (0=off, 1=on)
	WindowEnable bool // bit 5
	// TileAddressingMode changes the way that the PPU determines the memory addresses of Background and Window tiles
	// from the tile address offsets provided by the Background and Window tile maps.
	// If 0, tiles are addressed by $8000 + byte(offset); if 1, tiles are addressed by $8800 +/- int8(offset).
	// For a more detailed explanation, see the comments on the `getTileRow` function below.
	// See https://gbdev.gg8.se/wiki/articles/Video_Display#VRAM_Tile_Data
	TileAddressingMode bool // bit 4 (0=8800-97FF, 1=8000-8FFF)
	// BGTileMapSelect switches the region of memory the PPU reads the tile
	// map of the background from.
	BGTileMapSelect bool // bit 3 (0=9800-9BFF, 1=9C00-9FFF)
	// SpriteSize toggles whether all sprites are one tile in size or two tiles (0=8x8px, 1=8x16px)
	SpriteSize bool // bit 2
	// SpriteEnable toggles whether all sprites are rendered or not.
	SpriteEnable bool // bit 1
	// WindowDisplayORPriority
	// See https://gbdev.gg8.se/wiki/articles/LCDC#LCDC.0_-_BG.2FWindow_Display.2FPriority
	WindowDisplayORPriority bool // bit 0
}

// LCDStat represents a memory register located at [FIXME address] which is
// used to enable some interrupts related to drawing the LCD screen
type LCDStat struct {
	LYCoincidenceInterruptEnable bool // bit 6
	OAMInterruptEnable           bool // bit 5
	VBlankInterruptEnable        bool // bit 4
	HBlankInterruptEnalbe        bool // bit 3
	LYCoincidenceStatus          bool // bit 2 (0: LYC<>LY, 1: LYC=LY)
	Mode                         Mode // bits 1,0
}

func (p *PPU) readLCDStat() LCDStat {
	// FIXME implement -- read LCDStat from mmu
	return LCDStat{}
}

func (p *PPU) readLCDControl() LCDControl {
	// FIXME implement -- read LCDC from mmu
	return LCDControl{}
}

// oamEntry represents one 4-byte entry of sprite data (aka "Object Attributes")
// stored in OAM ram.
type oamEntry struct {
	// y represents the y-coordinate of the top-left of the sprite. (y=16 -> sprite fully visible on y-axis)
	y byte // byte 0
	// x represents the x-coordinate of the top-left of the sprite. (x=8 -> sprite fully visible on x-axis)
	x byte // byte 1
	// tileAddrOffset is used to determine the memory address of the sprite's tile data (tiles are 8x8 blocks of pixels).
	// The memory address is determined by adding tileAddrOffset to $8000.
	// If the SpriteSize bit of LCDControl is set to 1 (in which case all sprites are 1x2 tiles instead of 1 tile),
	// the tile at this address is the top tile of the sprite, and the tile at addr+1 is the bottom tile of the sprite. (FIXME confirm this.)
	tileAddrOffset   byte             // byte 2
	spriteAttributes spriteAttributes // byte 4
}

// spriteAttributes represents a one-byte bitfield
// that stores data and flags for a sprite in byte 4 of an oam Entry.
// See https://gbdev.gg8.se/wiki/articles/Video_Display#VRAM_Tile_Data
type spriteAttributes struct {
	// priority determines whether the sprite is rendered on top of the background tiles
	// or behind the first 3 colors of the background tile (but not the 4th). (0=above BG, 1=behind BG color 1-3)
	priority bool // bit 7
	// yFlip determines if the sprite is flipped vertically (0=normal, 1=flipped)
	yFlip bool // bit 6
	// xFlip determines if the sprite is flipped horizontally (0=normal, 1=flipped)
	xFlip bool // bit 5
	// palletteNumber determines which pallete is used to color the sprite if not in CGB mode.
	// (0=OBP0, 1=OBP1)
	palleteNumber bool // bit 4
	// tileVRAMBank determines which VRAM bank the sprite's tile data is stored in.
	// This option is only available in the Gameboy Color, which is
	tileVRAMBank bool // bit 3
	// cgbPaletteNumber chooses the color palette of the sprite in CGB mode (OBP0-7).
	// The Gameboy color supports 8 swappable palettes (as opposed to the Gameboy's 2 swappable palettes.)
	cgbPaletteNumber int // bit 2,1,0
}

// getOAMEntries reads the first ten sprite data entries that are on the current scanline ('y').
func (p *PPU) getOAMEntries(y byte, lcdc LCDControl) []oamEntry {
	return make([]oamEntry, 10)
}

// getSpriteRow returns the pixels, from left to right, of a certain row of the sprite.
// Sprite rows are 0-indexed and run from top to bottom.
// Sprites can be either 8 or 16 pixels tall, so the bottom row of a sprite can either be
// row 7 or row 15.
func (p *PPU) getSpriteRow(spriteData oamEntry, row int) []pixel {
	return make([]pixel, 8)
}

// getTileRow returns the pixels, from left to right, of a certain row of a tile located
// at a location in video memory determined by `addrOffset`.
// The way the tile's memory address is determined from `addrOffset` depends on the
// TileAddressingMode bit of the LCDControl register.
// If the bit is 0, the address is determined using the '$8800' method:
// `addrOffset` is treated as a signed byte and the memory address is $8800 +/- addrOffset.
// If the bit is 1, the address is determined using the '$8000' method (the same method sprites use):
// `addrOffset` is treated as an unsigned byte and the memory address is $8000 + addrOffset.
// Tile rows are 0-indexed and run from top to bottom, so the bottom row of a tile is row 7.
func (p *PPU) getTileRow(addrOffset byte, row int) []pixel {
	return make([]pixel, 8)
}

// palette represents the color palette used to color a tile.
// Each pixel of a tile has color, numbered from 1-4. The palette is
// a map from a color number to a color, allowing each tile to be
// colored with up to four different colors. If a tile is a sprite, color 4
// is always colored as transparent.
// In the original Gameboy, there are only 4 colors total to choose from.
type palette int

const (
	unspecifiedPalette palette = iota
	bg
	obj0
	obj1
	obj2 // CGB only
	obj3 // CGB only
	obj4 // CGB only
	obj5 // CGB only
	obj6 // CGB only
	obj7 // CGB only
)

type colorNumber int

const (
	unspecifiedColor colorNumber = iota
	col1
	col2
	col3
	col4
)

// FIXME is this the right abstraction? Should palette be included here?
type pixel struct {
	color   colorNumber
	palette palette
}

// getYScroll gets the y-coordinate of the top-left of the LCD screen.
// Reads the SCY ($FF42) memory register.
func (p *PPU) getScrollY() byte {
	return 0
}

// getXScroll gets the y-coordinate of the top-left of the LCD screen.
// Reads the SCX ($FF43) memory register.
func (p *PPU) getScrollX() byte {
	return 0
}

// getScanline gets the current scanline (0 through 153). The max 'scanline'
// is 153 and not 143 (the LCD screen only has 144 scanlines) because scanlines
// 144-153 represent the time spent in VBlank mode.
// Reads the LY ($FF44) memory register.
func (p *PPU) getScanline() byte {
	return 0
}

// getLY sets the current scanline.
// Writes the LY ($FF44) memory register.
func (p *PPU) getLY() byte {
	return 0
}

// getLYCompare gets the value of a register that is used to trigger
// an interrupt on a specific scanline.
// When a new scanline is started, the MMU compares the value of the LYC register
// with the current scanline (the LY register). If the values match and the
// LYCoincidenceInterruptEnabled bit is set in the LCDStat register, then the
// LYCoincidenceInterrupt is triggered and program execution jumps to that routine.
// Reads the LYC (LYCompare) ($FF45) memory register.
func (p *PPU) getLYCompare() byte {
	return 0
}

// Gets the X coordinate of the Window top left, minus 7.
// Reads WindowX-7($FF4B) memory register.
func (p *PPU) getWindowX() byte {
	return 0
}

// Gets the Y coordinate of the Window top left.
// Reads the WindowY($FF4A) memory registers.
func (p *PPU) getWindowY() byte {
	return 0
}

func shiftTileLeft(pixels []pixel, shift byte) []pixel {
	return pixels
}

func shiftTileRight(pixels []pixel, shift byte) []pixel {
	return pixels
}

func (p *PPU) colorize(px pixel) (r, g, b byte) {
	return 0, 0, 0
}

// Reference: https://gbdev.gg8.se/wiki/articles/Video_Display#VRAM_Tile_Data
// Reference: https://www.huderlem.com/demos/gameboy2bpp.html
func (p *PPU) getWindowTileRow(scX, scY, viewportX, viewportY byte, lcdc LCDControl) []pixel {
	// parse lcdc to see where to look up window tile map
	return make([]pixel, 8)
}

func (p *PPU) getBackgroundTileRow(scX, scY, viewportX, viewportY byte, lcdc LCDControl) []pixel {
	return make([]pixel, 8)
}
