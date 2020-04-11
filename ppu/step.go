package ppu

import "fmt"

// RunFor runs the PPU for a certain number of 4.14xxx MHz cycles.
// Cycles should always be divisible by 2.
func (p *PPU) RunFor(cycles int) {
	if cycles%2 != 0 {
		panic(fmt.Sprintf("ppu.RunFor(%v) called with cycles not divisible by 2", cycles))
	}
	for i := 0; i < cycles; i++ {
		p.Step()
	}
}

// Step executes 1 cycle's worth of work on the PPU.
func (p *PPU) Step() {
	p.cycles = (p.cycles + 1) % 456

	lcdstat := p.readLCDStat()

	if lcdstat.Mode != VBlank {
		// In non-VBlank mode, iterate through OAMSearch -> Pixel Drawing -> HBLank modes,
		// drawing a scanline every 456 clocks.
		switch p.cycles {
		case 0:
			p.setMode(OAMSearch)
		case 80:
			p.setMode(PixelDrawing)
			scanline := p.drawScanline(p.readLCDControl(), p.getLY(), p.getScrollX(), p.getScrollY())
			p.screen = append(p.screen, scanline...)
		case 160:
			p.setMode(HBlank)
		case 456:
			if p.getLY() < 143 { // 143 is last onscreen scanline
				p.setMode(OAMSearch)
				p.setLY(p.getLY() + 1)
			} else if p.getLY() == 143 {
				p.setMode(VBlank)
				// send screen to output channel
				p.videoOut <- p.screen
				p.screen = make([]byte, 0)
				p.setLY(p.getLY() + 1)
			} else {
				panic(fmt.Sprintf("LY is %v (>143), but mode is %v (should be VBlank)", p.getLY(), lcdstat.Mode))
			}
		}
		return

	} else if lcdstat.Mode == VBlank {
		// In VBlank mode, increment LY every 456 clocks until LY == 153,
		// at which point re-enter OAMSearch mode and resume drawing scanlines.
		if p.cycles == 456 {
			if p.getLY() < 153 {
				p.setLY(p.getLY() + 1)
			} else if p.getLY() == 153 {
				p.setLY(0)
				p.setMode(OAMSearch)
			} else {
				panic(fmt.Sprintf("LY is %v, should be less than 153", p.getLY(), lcdstat.Mode))
			}
		}
		return
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
		// Every 8 pixels (including x=0), fill the pixel fifo with the next tile.
		if x%8 == 0 {
			tile := p.getBackgroundTileRow(x+8, y, lcdc)
			pixelFifo.addTile(tile)
		}
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
	}
	return scanline
}
