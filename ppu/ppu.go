package ppu

import (
	"fmt"
	"github.com/mpingram/gameboy-emu/mmu"
)

type Screen []byte

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

type sprite struct {
	row [8]pixel
	x   byte
}

func (p *PPU) getSprites(yCoord byte, stat LCDStat) []sprite {
	return make([]sprite, 10)
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

type pixelFifo struct {
	fifo []pixel
}

func (pf *pixelFifo) overlay(sprite [8]pixel) {
	// TODO Implement
}

func (pf *pixelFifo) dequeue() (pixel, error) {
	if len(pf.fifo) > 0 {
		px := pf.fifo[0]
		pf.fifo = pf.fifo[1:]
		return px, nil
	}
	return pixel{}, fmt.Errorf("Pixel fifo is empty")
}

func (pf *pixelFifo) addTile(tile []pixel) {
	pf.fifo = append(pf.fifo, tile...)
}

func (pf *pixelFifo) clear() {
	pf.fifo = make([]pixel, 0)
}

func (pf *pixelFifo) size() int {
	return len(pf.fifo)
}

func (p *PPU) makePixelFifo(y byte, lcdc LCDControl) pixelFifo {
	return pixelFifo{}
}

func (p *PPU) getWindowCoords() (byte, byte) {
	return 0, 0
}

func (p *PPU) getScrollOffsets() (byte, byte) {
	return 0, 0
}

const screenHeight byte = 144
const screenWidth byte = 160

func (p *PPU) Update() {

	screen := make([]byte, int(screenWidth)*int(screenHeight)*3)

	// 144 times:
	// 		Read LCDStat
	// 		OAM Search (Read ten sprites on this line) (Read OAM, Write to LCDStat)
	// 		Pixel Drawing (Fill pixel fifo with bg tiles and sprites) (Read OAM(?) and VRAM, Write to LCDStat)
	// 		H-Blank (write to LCDStat)
	// V-Blank (write to LCDStat)

	// OAM Search
	// enter mode [2]
	// sprites = 10 sprites visible on this line

	// Pixel Drawing
	// fill pixel fifo with first 8 pixels
	// for x := 0; x < 160; x++ {
	//   for all sprites that start on this pixel(-8?), fetch the sprite and overlay it with bg tile (+10 clocks) (Problem -- sprites at pos 0 are off screen, they're on screen starting at pox 8!)
	//   if window starts on this pixel, throw away pixels in fifo, fetch window tile (+10 clocks?)
	//   (if end of line, clear fetcher and enter hblank)
	// }

	// For each scanline (144 times)
	for y := byte(0); y < screenHeight; y++ {

		// are we currently drawing the window?
		drawingWindow := false

		// Read relevant memory registers
		lcdStat := p.getLCDStat()
		lcdControl := p.getLCDControl()
		windowX, windowY := p.getWindowCoords()
		scX, scY := p.getScrollOffsets()

		// I. OAM Search -- find the sprites on this line
		p.enterMode(OAMSearch)
		// get the first 10 sprites that are on this y-coordinate
		sprites := p.getSprites(y, lcdStat)

		// II. Pixel Drawing mode
		p.enterMode(PixelDrawing)
		// Fill pixel fifo with the leftmost two tiles on this scanline
		pixelFifo := pixelFifo{}
		bgTile1 := p.getBackgroundTile(scX, scY, 0, y, lcdControl)
		bgTile2 := p.getBackgroundTile(scX+8, scY, 0, y, lcdControl)
		pixelFifo.addTile(bgTile1)
		pixelFifo.addTile(bgTile2)
		// dequeue pixel fifo XScroll % 8 times. This aligns our tiles to the edge of the screen.
		for i := byte(0); i < scX%8; i++ {
			pixelFifo.dequeue()
		}

		// For each pixel on this scanline:
		for x := byte(0); x < screenWidth; x++ {

			if lcdControl.SpriteEnable {
				// If a sprite starts at this xpos, layer that sprite on top of the bg tile.
				// Do this for every sprite that starts at this xpos.
				// FIXME -- need to do this before dequeing pixel fifo scX % 8 times. Or no, actually no. Huh.
				// How to deal with fact that sprite at coordinate 0,0 means 'off the screen'? Are sprite coordinates
				// even relative to screen space? Or are they relative to coordinate space? (Answer: screen space)
				for sprite := sprites[0]; sprite.x == x; {
					pixelFifo.overlay(sprite.row)
					// pop this element off the sprites array
					sprites = sprites[1:]
				}
			}

			if lcdControl.WindowEnable {
				// If window starts at this xpos and we're past the window ypos, clear fifo and start filling it with window tiles instead.
				if windowY >= y && windowX == x {
					drawingWindow = true
					pixelFifo.clear()
					windowTile := p.getWindowTile(x, y, scX, scY, lcdControl)
					pixelFifo.addTile(windowTile)
				}
			}

			// pop a pixel off the pixel fifo, colorize it, and send it to the screen
			px, err := pixelFifo.dequeue()
			if err != nil {
				panic(err)
			}
			r, g, b := p.colorize(px)
			screen = append(screen, r, g, b)

			// if pixel fifo has 8 pixels, fill it with the next tile.
			if pixelFifo.size() <= 8 {
				// get next tile -- if we're drawing background, get the next background tile,
				// otherwise get the next window tile.
				var tile []pixel
				if drawingWindow {
					tile = p.getWindowTile(x, y, scX, scY, lcdControl)
				} else {
					tile = p.getBackgroundTile(x, y, scX, scY, lcdControl)

				}
				pixelFifo.addTile(tile)
			}
		}
		// III. H-Blank mode (Horizontal blank) -- wait out the rest of the cycle.
		p.enterMode(HBlank)

	}
	// IV. V-Blank mode
	p.enterMode(VBlank)

}

func (p *PPU) colorize(px pixel) (r, g, b byte) {
	return 0, 0, 0
}

func (p *PPU) getWindowTile(scX, scY, xOffset, yOffset byte, lcdc LCDControl) []pixel {
	return make([]pixel, 8)
}

func (p *PPU) getBackgroundTile(scX, scY, xOffset, yOffset byte, lcdc LCDControl) []pixel {
	return make([]pixel, 8)
}
