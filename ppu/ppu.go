package ppu

import (
	"github.com/mpingram/gameboy-emu/mmu"
)

type PPU struct {
	mem mmu.MemoryReadWriter
}

type PPUMode int

const (
	OAMSearch    PPUMode = 2
	PixelDrawing         = 3
	HBlank               = 0
	VBlank               = 1
)

func (p *PPU) enterMode(mode PPUMode) {
	// FIXME implement - set bits 1 and 0 of LCDSTAT register
}

type LCDControl struct {
	LCDEnable               bool // bit 7
	WindowTileMapSelect     bool // bit 6 (0=$9800-$9BFF, 1=$9C00-$9FFF)
	WindowEnable            bool // bit 5
	TileAddressingMode      bool // bit 4 (0=8800-97FF, 1=8000-8FFF) // See http://gbdev.gg8.se/wiki/articles/Video_Display#LCDC.4_-_BG_.26_Window_Tile_Data_Select
	BGTileMapSelect         bool // bit 3 (0=9800-9BFF, 1=9C00-9FFF)
	SpriteSize              bool // bit 2 (0=8x8, 1=8x16)
	SpriteEnable            bool // bit 1
	WindowDisplayORPriority bool // bit 0 -- See http://gbdev.gg8.se/wiki/articles/Video_Display#LCDC.0_-_BG.2FWindow_Display.2FPriority
}

type LCDStat struct {
	LYCoincidenceInterruptEnable bool    // bit 6
	OAMInterruptEnable           bool    // bit 5
	VBlankInterruptEnable        bool    // bit 4
	HBlankInterruptEnalbe        bool    // bit 3
	LYCoincidenceStatus          bool    // bit 2 (0: LYC<>LY, 1: LYC=LY)
	Mode                         PPUMode // bits 1,0
}

func (p *PPU) getLCDStat() LCDStat {
	return LCDStat{}
}

func (p *PPU) getLCDControl() LCDControl {
	return LCDControl{}
}

// spriteAttrib represents the Sprite (aka Object) Attributes
// stored in OAM ram.
type spriteAttrib struct {
	y    byte // byte 0
	x    byte // byte 1
	addr byte // byte 2 -- represents memory location of sprite tile ($8000 + tileNo)
	bool      // byte 4
}

func (p *PPU) getOAMEntries(y byte, lcdc LCDControl) []spriteAttrib {
	return make([]spriteAttrib, 10)
}

func (p *PPU) fetchSpritePixels(sa spriteAttrib, row byte) []pixel {
	return make([]pixel, 8)
}

func (p *PPU) fetchTilePixels(addr byte, row byte) []pixel {
	return make([]pixel, 8)
}

type pallette int

const (
	bg   pallette = 0
	obj1          = 1
	obj2          = 2
)

type color int

const (
	col0 color = 0
	col1       = 1
	col2       = 2
	col3       = 3
)

type pixel struct {
	color   color
	palette pallette
}

func (p *PPU) getWindowCoords() (byte, byte) {
	return 0, 0
}

func (p *PPU) getScrollOffsets() (byte, byte) {
	return 0, 0
}

const screenHeight byte = 144
const screenWidth byte = 160

func shiftTileLeft(pixels []pixel, shift byte) []pixel {
	return pixels
}

func shiftTileRight(pixels []pixel, shift byte) []pixel {
	return pixels
}

func (p *PPU) colorize(px pixel) (r, g, b byte) {
	return 0, 0, 0
}

func (p *PPU) getWindowTile(scX, scY, viewportX, viewportY byte, lcdc LCDControl) []pixel {
	// parse lcdc to see where to look up window tile map
	return make([]pixel, 8)
}

func (p *PPU) getBackgroundTile(scX, scY, viewportX, viewportY byte, lcdc LCDControl) []pixel {
	return make([]pixel, 8)
}
