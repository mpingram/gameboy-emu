package ppu

import "fmt"

// RunFor runs the PPU for a certain number of 4.14xxx MHz cycles (dots).
// Cycles should always be divisible by 2.
func (p *PPU) RunFor(dots int) {
	if dots%2 != 0 {
		panic(fmt.Sprintf("ppu.RunFor(%v) called with dots not divisible by 2", dots))
	}
	for i := 0; i < dots; i += 2 {
		// Step does 2 dots worth of work; call it every 2 dots
		p.Step()
	}
}

// Step executes 2 dots' worth of work on the PPU.
func (p *PPU) Step() {
	p.dots += 2

	lcdstat := p.readLCDStat()

	switch lcdstat.Mode {
	case OAMSearch: // 80 dots total
		if p.dots == 2 {
			// Find the first 10 sprites on this line (8 dots each)
			// get the first 10 sprites that are on this y-coordinate
			p.sprites = p.getOAMEntries(p.getLY(), p.readLCDControl())
		} else if p.dots == 80 {
			p.setMode(PixelDrawing)
		}

	case PixelDrawing: // takes 160 + about 10 dots per sprite total
		if p.dots == 82 {
			scanline := p.drawScanline(p.readLCDControl(), p.getLY(), p.getScrollX(), p.getScrollY())
			p.screen = append(p.screen, scanline...)
		} else if p.dots == 160+10*len(p.sprites) {
			p.setMode(HBlank)
		}

	case HBlank: // ends after 456 dots
		// do nothing
		if p.dots == 456 {
			if p.getLY() == 143 { // 143 is last onscreen scanline
				p.setMode(VBlank)
				p.setLY(p.getLY() + 1)
			} else {
				p.setMode(OAMSearch)
				p.setLY(p.getLY() + 1)
			}
			// reset the dot clock
			p.dots = 0
		}

	case VBlank: // takes 456 dots
		if p.dots == 456 {
			// If it's the first scanline, render the screen
			if p.getLY() == 144 {
				// send screen to output channel
				p.videoOut <- p.screen
				p.screen = make([]byte, 0)
				p.setLY(p.getLY() + 1)
			} else if p.getLY() == 153 { // 153 is the last scanline (?)
				// If it's the last scanline, enter OAM search
				// enter OAMSearch mode
				p.setMode(OAMSearch)
				p.setLY(0)
			} else {
				// otherwise just move to the next scanline; remain
				// in VBlank mode.
				p.setLY(p.getLY() + 1)
			}
			p.dots = 0
		}
	}
}

func (p *PPU) drawScanline(lcdc LCDControl, ly, scX, scY byte) []byte {
	// NOTE this implementation currently completely ignores the Window and sprites.
	var scanline []byte
	y := scY + ly // y is the global y-coordinate of the current scanline.

	// Initialize the pixel fifo with pixels from the tile that intersects with scX.
	pixelFifo := pixelFifo{}
	bgTile := p.getBackgroundTileRow(scX-(scX%8), y, lcdc)
	pixelFifo.addTile(bgTile)

	for x := scX - scX%8; x < scX+screenWidth; x++ {
		// Dequeue a pixel from the pixel fifo.
		px, err := pixelFifo.dequeue()
		if err != nil {
			panic(err)
		}
		// If x is onscreen, draw the pixel to the current scanline.
		if x >= scX {
			palette := p.getBGPalette()
			color := palette[px.color]
			rgb := toRGB(color)
			scanline = append(scanline, rgb...)
		}
		// Every 8 pixels (including x=0), fill the pixel fifo with the next tile.
		if x%8 == 0 {
			tile := p.getBackgroundTileRow(x+8, y, lcdc)
			pixelFifo.addTile(tile)
		}
	}
	return scanline
}
