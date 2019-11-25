package ppu

import "fmt"

type MemoryReadWriter interface {
	MemoryReader
	MemoryWriter
}

type MemoryReader interface {
	Rb(addr uint16) byte
	Rw(addr uint16) uint16
}

type MemoryWriter interface {
	Wb(addr uint16, b byte)
	Ww(addr uint16, bb uint16)
}

type PPU struct {
	mem MemoryReadWriter
}

func New(mem MemoryReadWriter) *PPU {
	ppu := &PPU{mem}
	return ppu
}

const screenHeight = 144
const screenWidth = 160

// Mode represents the 'drawing mode' of the Gameboy,
// which is stored in bits 0 and 1 of the LCDStat memory register.
// The drawing mode cycles between OAMSearch, PixelDrawing, and HBlank
// for each scanline, and enters VBlank once a full screen has been drawn.
type Mode int

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
	// map of the background from. (0=9800-9BFF, 1=9C00-9FFF)
	BGTileMapSelect bool // bit 3 (0=9800-9BFF, 1=9C00-9FFF)
	// SpriteSize toggles whether all sprites are one tile in size or two tiles (0=8x8px, 1=8x16px)
	SpriteSize bool // bit 2
	// SpriteEnable toggles whether all sprites are rendered or not.
	SpriteEnable bool // bit 1
	// WindowDisplayORPriority
	// See https://gbdev.gg8.se/wiki/articles/LCDC#LCDC.0_-_BG.2FWindow_Display.2FPriority
	WindowDisplayORPriority bool // bit 0
}

func (p *PPU) readLCDControl() LCDControl {
	var lcdcAddr uint16 = 0xff40
	b := p.mem.Rb(lcdcAddr)
	return LCDControl{
		LCDEnable:               b&0b1000_0000 > 0,
		WindowTileMapSelect:     b&0b0100_0000 > 0,
		WindowEnable:            b&0b0010_0000 > 0,
		TileAddressingMode:      b&0b0001_0000 > 0,
		BGTileMapSelect:         b&0b0000_1000 > 0,
		SpriteSize:              b&0b0000_0100 > 0,
		SpriteEnable:            b&0b0000_0010 > 0,
		WindowDisplayORPriority: b&0b0000_0001 > 0,
	}
}

// LCDStat represents a memory register located at [FIXME address] which is
// used to enable some interrupts related to drawing the LCD screen
// FIXME should PPU care about this? Or MMU?
// Why should the PPU care about it? The only bits it needs to set are Mode (verify this?)
type LCDStat struct {
	LYCoincidenceInterruptEnable bool // bit 6
	OAMInterruptEnable           bool // bit 5
	VBlankInterruptEnable        bool // bit 4
	HBlankInterruptEnalbe        bool // bit 3
	LYCoincidenceStatus          bool // bit 2 (0: LYC<>LY, 1: LYC=LY)
	Mode                         Mode // bits 1,0
}

func (p *PPU) readLCDStat() LCDStat {
	lcdStatAddr := uint16(0xFF41)
	b := p.mem.Rb(lcdStatAddr)
	return LCDStat{
		LYCoincidenceInterruptEnable: b&0b0100_0000 > 0,
		OAMInterruptEnable:           b&0b0010_0000 > 0,
		VBlankInterruptEnable:        b&0b0001_0000 > 0,
		HBlankInterruptEnalbe:        b&0b0000_1000 > 0,
		LYCoincidenceStatus:          b&0b0000_0100 > 0,
		Mode:                         Mode(b & 0b0000_0010),
	}
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
	// FIXME implement
	return make([]oamEntry, 10)
}

// getSpriteRow returns the 8 pixels, from left to right, of a certain row of the sprite.
// Sprite rows are 0-indexed and run from top to bottom.
// Sprites can be either 8 or 16 pixels tall, so the bottom row of a sprite can either be
// row 7 or row 15.
func (p *PPU) getSpriteRow(spriteData oamEntry, row byte) []pixel {
	// FIXME implement
	return make([]pixel, 8)
}

// getWindowTileRow returns the 8 pixels of a row of a window tile located at
// screen-based coordinate screenX, screenY. Note that the xy coordinates are based on
// _the top left of the screen_.
// Reference: https://gbdev.gg8.se/wiki/articles/Video_Display#VRAM_Tile_Data
func (p *PPU) getWindowTileRow(screenX, screenY byte, lcdc LCDControl) []pixel {
	// FIXME implement
	// calculate location in window tile map based on scX / scY
	// use lcdc to check where window tile map is stored
	return make([]pixel, 0)
}

// getBackgroundTileRow returns the 8 pixels of a row of a background tile located at
// coordinate x, y.
func (p *PPU) getBackgroundTileRow(x, y byte, lcdc LCDControl) []pixel {
	// check lcdc to see where bg tile map is stored
	var tileMapLocation uint16
	if lcdc.BGTileMapSelect == false {
		tileMapLocation = 0x9800
	} else {
		tileMapLocation = 0x9C00
	}
	// calculate byte offset in bg tile map based on x,y
	y = y % screenHeight
	x = x % screenWidth
	offset := (y/8)*32 + (x/8)
	// Explanation:
	// Tile memory is laid out like this:
	// $9BFF/$9FFF +-------------------+
	//             |                   |
	//             |                   | 
	// 32 tile ptrs|                   |
	//         |   |32|33|34|...       |
	//         y   |0 |1 |2 |...       |
	// $9800/$9C00 +--32 tile ptrs-----+ 
	// 					x ----->
	// Where each tile represents a 8x8px area. So the formula for 
	// getting the tile byte offset is floor(y/8)*32 + floor(x/8).
	// To accommodate wraparound, we set y=y%144 and x=x%160.

	// Read the correct byte of the tile map to get the address of the tile data
	// (Remember that the address of the tile data is an offset, not a full uint16 addresss.)
	tileAddrOffset := p.mem.Rb(tileMapLocation+uint16(offset))	
	// The row of the tile that intersects with this y-coordinate. Rows go from 0-7,
	// where 7 is the bottom row.
	row := y % 8
	tileData := p.getTileRowData(tileAddrOffset, row, lcdc)
	pixels := make([]pixel, 8)
	for _, colorNumber := range tileData {
		px := pixel{colorNumber, *bgPalette}
		pixels = append(pixels, px)
	}
	return pixels
}

// getTileRowData returns the color numbers, from left to right, of a certain row of a tile located
// at a location in video memory determined by `addrOffset`.
// The way the tile's memory address is determined from `addrOffset` depends on the
// TileAddressingMode bit of the LCDControl register.
// If the bit is 0, the address is determined using the '$8800' method:
// `addrOffset` is treated as a signed byte and the memory address is $8800 +/- addrOffset.
// If the bit is 1, the address is determined using the '$8000' method (the same method sprites use):
// `addrOffset` is treated as an unsigned byte and the memory address is $8000 + addrOffset.
// Tile rows are 0-indexed and run from top to bottom, so the bottom row of a tile is row 7.
func (p *PPU) getTileRowData(addrOffset byte, row byte, lcdc LCDControl) []colorNumber {
	if row > 7 {
		panic(fmt.Sprintf("Got tile row > 7: %v", row))
	}
	var tileAddr uint16
	if lcdc.TileAddressingMode == false {
		// convert addrOffset to a signed byte
		signedAddrOffset := int8(addrOffset)
		// NOTE this is potentially buggy!
		// promote the signed int8 to int in order to add it to 0x8800,
		// then convert the result back to uint16.
		tileAddr = uint16(0x8800 + int(signedAddrOffset))
	} else {
		tileAddr = 0x8000 + uint16(addrOffset)
	}

	// Each 2 bytes of the tile is a row of the tile, and they are stored from top
	// to bottom. The bytes that represents row n of the tile are at (tileAddr + n*2)
	b1 := p.mem.Rb(tileAddr + uint16(row*2))
	b2 := p.mem.Rb(tileAddr + uint16(row*2) + 1)
	tileData := make([]colorNumber, 8)
	for i := 7; i >= 0; i-- {
		// WARNING Possibly buggy
		// https://gbdev.gg8.se/wiki/articles/Video_Display#VRAM_Tile_Data
		// b1 contains the low bit of each px, from left (bit 7) to right (bit 0)
		// b2 contains the high bit of each px, as above.
		mask := byte(1) << i
		lo := (b1 & mask) >> i
		hi := (b2 & mask) >> i
		color := (hi << 1) | lo
		tileData = append(tileData, colorNumber(color))
	}

	return tileData
}

// palette represents the color palette used to color a tile.
// Each pixel of a tile has color, numbered from 1-4. The palette is
// a map from a color number to a RGB color, allowing each tile to be
// colored with up to four different colors. If a tile is a sprite, color 4
// is always colored as transparent.
// In the original Gameboy, there are only 4 colors total to choose from.
type palette map[colorNumber]color

var bgPalette, obj0Palette, obj1Palette *palette

func (p *PPU) getBGPalette() palette {
	var bgpAddr uint16 = 0xFF47
	b := p.mem.Rb(bgpAddr)
	pal := map[colorNumber]color{
		col3: color(b & 0b1100_0000 >> 6),
		col2: color(b & 0b0011_0000 >> 4),
		col1: color(b & 0b0000_1100 >> 2),
		col0: color(b & 0b0000_0011),
	}
	return palette(pal)
}

func (p *PPU) getObj0Palette() palette {
	var obp0Addr uint16 = 0xFF48
	b := p.mem.Rb(obp0Addr)
	pal := map[colorNumber]color{
		col3: color(b & 0b1100_0000 >> 6),
		col2: color(b & 0b0011_0000 >> 4),
		col1: color(b & 0b0000_1100 >> 2),
		col0: color(b & 0b0000_0011),
	}
	return palette(pal)
}

func (p *PPU) getObj1Palette() palette {
	var obp1Addr uint16 = 0xFF48
	b := p.mem.Rb(obp1Addr)
	pal := map[colorNumber]color{
		col3: color(b & 0b1100_0000 >> 6),
		col2: color(b & 0b0011_0000 >> 4),
		col1: color(b & 0b0000_1100 >> 2),
		col0: color(b & 0b0000_0011),
	}
	return palette(pal)
}

type color byte

const (
	white     color = 0
	lightGray       = 1
	darkGray        = 2
	black           = 3
)

type colorNumber byte

const (
	col0 colorNumber = 0
	col1             = 1
	col2             = 2
	col3             = 3
)

type pixel struct {
	color   colorNumber
	palette palette
}

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

// func shiftTileLeft(pixels []pixel, shift byte) []pixel {
// 	return pixels
// }

// func shiftTileRight(pixels []pixel, shift byte) []pixel {
// 	return pixels
// }

// func (p *PPU) colorize(px pixel) [3]byte {
// 	return px.palette[px.color]
// }
