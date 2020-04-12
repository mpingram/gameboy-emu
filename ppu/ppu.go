package ppu

import "fmt"

type MemoryReadWriter interface {
	MemoryReader
	MemoryWriter
}

type MemoryReader interface {
	// Rb reads one byte.
	Rb(addr uint16) byte
	// Rw reads one word.
	Rw(addr uint16) uint16
}

type MemoryWriter interface {
	// Wb writes one byte.
	Wb(addr uint16, b byte)
	// Ww writes one word.
	Ww(addr uint16, bb uint16)
}

type PPU struct {
	mem      MemoryReadWriter
	cycles   int
	screen   []Pixel
	videoOut chan []Pixel
}

func New(mem MemoryReadWriter, videoOut chan []Pixel) *PPU {
	ppu := &PPU{mem, 0, []Pixel{}, videoOut}
	ppu.setMode(OAMSearch)
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
	lcdStatAddr := uint16(0xFF41)
	// set the bottom two bits of the lcdstat byte equal to the binary
	// representation of the mode
	b := p.mem.Rb(lcdStatAddr)
	// clear the bottom two bits
	b &= 0b1111_1100
	// set the bottom two bits
	b |= byte(mode)
	p.mem.Wb(lcdStatAddr, b)
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
		Mode:                         Mode(b & 0b0000_0011),
	}
}

// getBackgroundTileRow returns the 8 pixels of a row of a background tile located at
// coordinate x, y.
func (p *PPU) getBackgroundTileRow(x, y byte) []pixelData {
	// check lcdc to see where bg tile map is stored
	lcdc := p.readLCDControl()
	var tileMapLocation uint16
	if lcdc.BGTileMapSelect == false {
		tileMapLocation = 0x9800
	} else {
		tileMapLocation = 0x9C00
	}
	// calculate byte offset in bg tile map based on x,y
	offset := (y/8)*32 + (x / 8)
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
	tileAddrOffset := p.mem.Rb(tileMapLocation + uint16(offset))
	var tileAddr uint16
	if lcdc.TileAddressingMode == true {
		// convert addrOffset to a signed byte
		signedAddrOffset := int8(tileAddrOffset)
		// NOTE this is potentially buggy!
		// promote the signed int8 to int in order to add it to 0x8800,
		// then convert the result back to uint16.
		tileAddr = uint16(0x8800 + int(signedAddrOffset))
	} else {
		tileAddr = 0x8000 + uint16(tileAddrOffset)
	}
	// The row of the tile that intersects with this y-coordinate. Rows go from 0-7,
	// where 7 is the bottom row.
	row := y % 8
	tileData := p.getTileRowData(tileAddr, row)
	pixels := make([]pixelData, 8)
	for _, paletteIndex := range tileData {
		px := pixelData{paletteIndex, bg}
		pixels = append(pixels, px)
	}
	return pixels
}

// getTileRowData returns the color numbers, from left to right, of a certain row of a tile located
// at a location in video memory determined by `addrOffset`.
// Tile rows are 0-indexed and run from top to bottom, so the bottom row of a tile is row 7.
func (p *PPU) getTileRowData(tileAddr uint16, row byte) []byte {
	if row > 7 {
		panic(fmt.Sprintf("Got tile row > 7: %v", row))
	}

	// Each 2 bytes of the tile is a row of the tile, and they are stored from top
	// to bottom. The bytes that represents row n of the tile are at (tileAddr + n*2)
	b1 := p.mem.Rb(tileAddr + uint16(row*2))
	b2 := p.mem.Rb(tileAddr + uint16(row*2) + 1)
	tileData := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		// WARNING Possibly buggy
		// https://gbdev.gg8.se/wiki/articles/Video_Display#VRAM_Tile_Data
		// b1 contains the low bit of each px, from left (bit 7) to right (bit 0)
		// b2 contains the high bit of each px, as above.
		mask := byte(1) << i
		lo := (b1 & mask) >> i
		hi := (b2 & mask) >> i
		color := (hi << 1) | lo
		tileData = append(tileData, color)
	}

	return tileData
}

// palette represents the color palette used to color a tile.
// Each pixel of a tile has color, numbered from 1-4. The palette is
// a map from a color number to a RGB color, allowing each tile to be
// colored with up to four different colors. If a tile is a sprite, color 4
// is always colored as transparent.
// In the original Gameboy, there are only 4 colors total to choose from.
type palette map[byte]Pixel

func (p *PPU) getBGPalette() palette {
	var bgpAddr uint16 = 0xFF47
	b := p.mem.Rb(bgpAddr)
	pal := map[byte]Pixel{
		3: Pixel(b & 0b1100_0000 >> 6),
		2: Pixel(b & 0b0011_0000 >> 4),
		1: Pixel(b & 0b0000_1100 >> 2),
		0: Pixel(b & 0b0000_0011),
	}
	return palette(pal)
}

type Pixel byte

const (
	White     Pixel = 0
	LightGray       = 1
	DarkGray        = 2
	Black           = 3
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
	// fmt.Printf("LY: (%08x) %d\n", p.mem.Rb(lyAddr), p.mem.Rb(lyAddr))
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
// 	1 pixel
//  |-----|
// [R, G, B, R, G, B, R, G, B]
// Where R,G,B are one byte representing the red, green, blue
// component of each pixel.
type Screen []byte

type pixelFifo struct {
	fifo []pixelData
}

func (pf *pixelFifo) overlay(sprite []pixelData) {
	// overlay the sprite's pixels on top of the leftmost 8 pixels in the fifo
	for i, px := range sprite {
		// Sprite color #0 is transparent -- don't overlay it.
		// TODO implement OBJ-to-BG priority (https://gbdev.gg8.se/wiki/articles/Video_Display#FF48_-_OBP0_-_Object_Palette_0_Data_.28R.2FW.29_-_Non_CGB_Mode_Only)
		if px.colorNumber != 0 {
			pf.fifo[i] = px
		}
	}
}

func (pf *pixelFifo) dequeue() (pixelData, error) {
	if len(pf.fifo) > 0 {
		px := pf.fifo[0]
		pf.fifo = pf.fifo[1:]
		return px, nil
	}
	return pixelData{}, fmt.Errorf("Pixel fifo is empty")
}

func (pf *pixelFifo) addTile(tile []pixelData) {
	pf.fifo = append(pf.fifo, tile...)
}

func (pf *pixelFifo) clear() {
	pf.fifo = make([]pixelData, 0)
}

func (pf *pixelFifo) size() int {
	return len(pf.fifo)
}
