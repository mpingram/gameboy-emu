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

	ly := p.getLY()
	lcdstat := p.readLCDStat()
	lcdc := p.readLCDControl()

	switch lcdstat.Mode {
	case OAMSearch: // 80 dots total
		if p.dots == 2 {
			// Find the first 10 sprites on this line (8 dots each)
			// get the first 10 sprites that are on this y-coordinate
			p.sprites = p.getOAMEntries(ly, lcdc)
			//fmt.Println("Got OAM sprites")
		} else if p.dots == 80 {
			p.setMode(PixelDrawing)
			//fmt.Println("Entered PixelDrawing mode")
		}

	case PixelDrawing: // takes 160 + about 10 dots per sprite total
		if p.dots == 82 {
			p.screen = append(p.screen, p.drawScanline()...)
			//fmt.Println("Drew scanline!")
		} else if p.dots == 160+10*len(p.sprites) {
			p.setMode(HBlank)
			//fmt.Println("Entered HBlank mode")
		}

	case HBlank: // ends after 456 dots
		// do nothing
		if p.dots == 456 {
			if ly == 143 { // 143 is last onscreen scanline
				p.setMode(VBlank)
				//fmt.Println("Entered VBlank mode")
			} else {
				p.setMode(OAMSearch)
				p.setLY(ly + 1)
				//fmt.Println("Entered OAMSearch mode")
			}
			// reset the dot clock
			p.dots = 0
		}

	case VBlank: // takes 456 dots
		if p.dots == 456 {
			fmt.Printf("Finished line %v in VBlank!\n", ly)
			// If it's the last scanline, render the screen
			if ly == 153 { // 153 is the last scanline (?)
				// send screen to output channel
				p.videoOut <- p.screen
				fmt.Println("Sent screen to video out!")
				p.screen = make([]byte, 0)

				// enter OAMSearch mode
				p.setMode(OAMSearch)
				//fmt.Println("Entered OAMSearch mode")
				p.setLY(0)
				p.dots = 0
			} else {
				// otherwise just move to the next scanline; remain
				// in VBlank mode.
				p.setLY(ly + 1)
			}
		}
	}
}

func (p *PPU) drawScanline() []byte {

	var scanline []byte
	ly := p.getLY()

	// 1. Start 8x to the left of the screen edge (scX-8). Fill the pixel fifo with pixels
	// from a) the tile that intersects with scX-8 and b) the tile to its right.
	// Note that the pixel fifo probably contains extra pixels from tile A, unless tile A's
	// left edge happens to line up perfectly with scX-8.
	pixelFifo := pixelFifo{}
	// scY is an absolute coordinate (based on 0,0, just like the background tile map),
	// while y is a screen-space coordinate of the current scanline (based on scY.)
	// So, scY + y is the absolute coordinate of the current scanline. Because
	// absolute coordinates have one tile every 8px, (scY + y) % 8 gets us the row
	// of the tile that the current scanline is on.
	bgTileRow1 := p.getBackgroundTileRow((p.scX - 8), p.scY+ly, p.lcdc)
	bgTileRow2 := p.getBackgroundTileRow(p.scX, p.scY+ly, p.lcdc)
	pixelFifo.addTile(bgTileRow1)
	pixelFifo.addTile(bgTileRow2)

	// 2. Discard pixels from pixel fifo until it lines up with SCX-8. (Dequeue (SCX-8) % 8 times)
	for i := byte(0); i < (p.scX-8)%8; i++ {
		pixelFifo.dequeue()
	}

	drawingWindow := false
	// 3. For x = -8 up to 159: (x=-8 corresponds to scX-8)
	for x := -8; x < screenWidth; x++ {

		// 3a. Check if the window starts at this pixel. (Because the window's x coordinate
		// is relative to the screen, and we started 8px to the left of the screen, we
		// don't need to worry about the edge case of wX (window.X) being less than x.)
		if p.lcdc.WindowEnable {
			// If window starts at this xpos and we're past the window ypos, clear fifo and start filling it with window tiles instead.
			// wX and x are cast to ints to avoid underflow issues.
			if p.wY <= ly && int(p.wX) == x {
				drawingWindow = true
				pixelFifo.clear()
				// Because window tiles are aligned to the screen, the row we're looking
				// for is just y (current screen coordinate) mod 8.
				// We know that x is at least 0 and x is at most 159.
				// (x is at least zero because of the int(wX) == x check earlier,
				// and at most 159 due to the for loop.)
				// Therefore the byte(x) conversion is safe (i.e, x will always fit into a byte.)
				windowTileRowA := p.getWindowTileRow(byte(x), ly, p.lcdc)
				windowTileRowB := p.getWindowTileRow(byte(x)+8, ly, p.lcdc)
				pixelFifo.addTile(windowTileRowA)
				pixelFifo.addTile(windowTileRowB)
			}
		}

		// 3b. Overlay all sprites that start on this pixel onto the pixel fifo,
		if p.lcdc.SpriteEnable {
			// For every sprite that starts at this xpos, overlay that sprite on top of
			// the pixel fifo. sprite.x is shifted 8px to the left of the screen
			// coordinate system.
			for sprite := p.sprites[0]; int(sprite.x)-8 == x; {
				// sprite.y is the screen coordinate of the top of the sprite, shifted up
				// by 16px. We know that (sprite.y + 16) <= y (because the sprite is on this row
				// according to fetchOAMEntries.). So we can get the row of the sprite
				// by y - sprite.y. E.g., if sprite.y=20 and y=26, the row we want is 26 - 20 == 6.
				row := ly - (sprite.y + 16)
				pixels := p.getSpriteRow(sprite, row, p.lcdc)
				pixelFifo.overlay(pixels)
				// pop this element off the sprites array
				p.sprites = p.sprites[1:]
			}
		}

		// 3c. Dequeue a pixel from the pixel fifo, colorize it, and draw it to the screen.
		px, err := pixelFifo.dequeue()
		if err != nil {
			panic(err)
		}
		// colorize the pixel -- look up its color number in the provided palette.
		var palette palette
		if px.paletteNumber == bg {
			palette = p.getBGPalette()
		} else if px.paletteNumber == obj0 {
			palette = p.getObj0Palette()
		} else if px.paletteNumber == obj1 {
			palette = p.getObj1Palette()
		} else {
			panic(fmt.Sprintf("Bad pallete number: %v", px.paletteNumber))
		}
		color := palette[px.color]

		// Draw the pixel to the scanline
		rgb := toRGB(color)
		scanline = append(scanline, rgb...)

		// 3d. If the pixel fifo contains only one tile, add the next tile to it.
		if pixelFifo.size() <= 8 {
			// Get next tile -- if we're drawing background, get the next background tile,
			// otherwise get the next window tile.
			var tile []pixel
			if drawingWindow {
				// The byte(x+8) conversion here is safe:
				// Let n represent the number of pixels in the pixel fifo at the start
				// of the loop (x=-8). The pixel fifo is filled with 16 pixels and
				// dequeued up to 7 times, so n >= 9. When n=9, because the pixel fifo
				// is dequeued once every loop and x is incremented once every loop,
				// the pixel fifo will contain 8 pixels and trigger this code path when x=-7.
				// In that case, x+8 > 0, meaning byte(x+8) is safe.
				tile = p.getWindowTileRow(byte(x+8), ly, p.lcdc)
			} else {
				// Because background tiles are on the global coordinate system,
				// the row we're looking for is (scY + y) % 8.
				tile = p.getBackgroundTileRow(p.scX+byte(x+8), p.scY+ly, p.lcdc)
			}
			pixelFifo.addTile(tile)
		}
	}

	return scanline
}
