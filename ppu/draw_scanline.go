package ppu

import (
	"fmt"
)

const BGPaletteAddr = 0xFF47

type pixelFifo struct {
	fifo []byte
}

func (pf *pixelFifo) overlay(sprite []byte, bg byte) {
	// overlay the sprite's pixels on top of the leftmost 8 pixels in the fifo
	for i, px := range sprite {
		// Sprite color #0 is transparent -- don't overlay it.
		// TODO implement OBJ-to-BG priority (https://gbdev.gg8.se/wiki/articles/Video_Display#FF48_-_OBP0_-_Object_Palette_0_Data_.28R.2FW.29_-_Non_CGB_Mode_Only)
		if px != 0 {
			pf.fifo[i] = px
		}
	}
}

// FIXME this is one of the more expensive operations left in ppu.
// I think the best thing to do is to figure out how we can
// make the fifo dirt simple.
func (pf *pixelFifo) dequeue() (byte, error) {
	if pf.size() > 0 {
		px := pf.fifo[0]
		pf.fifo = pf.fifo[1:]
		return px, nil
	}
	return 0, fmt.Errorf("Pixel fifo is empty")
}

func (pf *pixelFifo) addTile(tile []byte) {
	pf.fifo = append(pf.fifo, tile...)
}

func (pf *pixelFifo) clear() {
	pf.fifo = make([]byte, 0)
}

func (pf *pixelFifo) size() int {
	return len(pf.fifo)
}

// NOTE this implementation currently completely ignores the Window and sprites.
func (p *PPU) drawScanline(ly, scX, scY byte) []byte {

	// Fetch the current background palette.
	palette := p.mem.Rb(BGPaletteAddr)
	col0 := palette & 0b0000_0011
	col1 := (palette & 0b0000_1100) >> 2
	col2 := (palette & 0b0011_0000) >> 4
	col3 := (palette & 0b1100_0000) >> 6

	// scanline := make([]byte, 0, 160)
	y := scY + ly // y is the global y-coordinate of the current scanline.

	// Initialize the pixel fifo with pixels from the tile that intersects with scX.
	// FIXME optimization: don't allocate a new pixelFifo each time
	pixelFifo := pixelFifo{}
	bgTile := p.getBackgroundTileRow(scX-(scX%8), y)
	pixelFifo.addTile(bgTile)

	// food for thought... what happens when x is 255?
	i := 0
	for x := scX - (scX % 8); x < scX+screenWidth; x++ {
		// Every 8 pixels (including x=0), fill the pixel fifo with the next tile.
		if x%8 == 0 {
			tile := p.getBackgroundTileRow(x+8, y)
			pixelFifo.addTile(tile)
		}
		// Dequeue a pixel from the pixel fifo.
		px, err := pixelFifo.dequeue()
		if err != nil {
			panic(err)
		}
		// If x is onscreen, draw the pixel to the current scanline.
		if x >= scX {
			var color byte
			switch px {
			case 0:
				color = col0
			case 1:
				color = col1
			case 2:
				color = col2
			case 3:
				color = col3
			default:
				panic(fmt.Errorf("Got bad color number: %d", px))
			}

			p.scanline[i] = color
			i += 1
		}
	}
	return p.scanline
}
