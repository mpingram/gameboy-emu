package ppu

import "fmt"

func (p *PPU) DrawScreen() Screen {

	screen := make(Screen, 0)

	var maxLY byte = screenHeight + 10 // 10 extra rows indicate VBlank period
	for y := byte(0); y < maxLY; y++ {

		// set the LY register
		p.setLY(y)
		// If y > 143, we are in VBlank mode -- do nothing.
		if y > screenHeight-1 {
			p.setMode(VBlank)
			continue
		}

		// re-read these registers before every scanline (after every V-Blank period)
		lcdc := p.readLCDControl()
		wX := p.getWindowX()
		wY := p.getWindowY()
		scX := p.getScrollX()
		scY := p.getScrollY()

		// reset the flag for whether or not we're drawing the window
		drawingWindow := false

		// I. OAM Search
		// =====================
		p.setMode(OAMSearch)
		// get the first 10 sprites that are on this y-coordinate
		sprites := p.getOAMEntries(y, lcdc)

		// II. Pixel Drawing mode
		// =====================
		p.setMode(PixelDrawing)

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
		bgTileRow1 := p.getBackgroundTileRow((scX - 8), scY+y, lcdc)
		bgTileRow2 := p.getBackgroundTileRow(scX, scY+y, lcdc)
		pixelFifo.addTile(bgTileRow1)
		pixelFifo.addTile(bgTileRow2)

		// 2. Discard pixels from pixel fifo until it lines up with SCX-8. (Dequeue (SCX-8) % 8 times)
		for i := byte(0); i < (scX-8)%8; i++ {
			pixelFifo.dequeue()
		}

		// 3. For x = -8 up to 159: (x=-8 corresponds to scX-8)
		for x := -8; x < screenWidth; x++ {

			// 3a. Check if the window starts at this pixel. (Because the window's x coordinate
			// is relative to the screen, and we started 8px to the left of the screen, we
			// don't need to worry about the edge case of wX (window.X) being less than x.)
			if lcdc.WindowEnable {
				// If window starts at this xpos and we're past the window ypos, clear fifo and start filling it with window tiles instead.
				// wX and x are cast to ints to avoid underflow issues.
				if wY <= y && int(wX) == x {
					drawingWindow = true
					pixelFifo.clear()
					// Because window tiles are aligned to the screen, the row we're looking
					// for is just y (current screen coordinate) mod 8.
					// We know that x is at least 0 and x is at most 159.
					// (x is at least zero because of the int(wX) == x check earlier,
					// and at most 159 due to the for loop.)
					// Therefore the byte(x) conversion is safe (i.e, x will always fit into a byte.)
					windowTileRowA := p.getWindowTileRow(byte(x), y, lcdc)
					windowTileRowB := p.getWindowTileRow(byte(x)+8, y, lcdc)
					pixelFifo.addTile(windowTileRowA)
					pixelFifo.addTile(windowTileRowB)
				}
			}

			// 3b. Overlay all sprites that start on this pixel onto the pixel fifo,
			if lcdc.SpriteEnable {
				// For every sprite that starts at this xpos, overlay that sprite on top of
				// the pixel fifo. sprite.x is shifted 8px to the left of the screen
				// coordinate system.
				for sprite := sprites[0]; int(sprite.x)-8 == x; {
					// sprite.y is the screen coordinate of the top of the sprite, shifted up
					// by 16px. We know that (sprite.y + 16) <= y (because the sprite is on this row
					// according to fetchOAMEntries.). So we can get the row of the sprite
					// by y - sprite.y. E.g., if sprite.y=20 and y=26, the row we want is 26 - 20 == 6.
					row := y - (sprite.y + 16)
					pixels := p.getSpriteRow(sprite, row, lcdc)
					pixelFifo.overlay(pixels)
					// pop this element off the sprites array
					sprites = sprites[1:]
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
			// Draw the pixel to the screen.
			rgb := toRGB(color)
			screen = append(screen, rgb...)

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
					tile = p.getWindowTileRow(byte(x+8), y, lcdc)
				} else {
					// Because background tiles are on the global coordinate system,
					// the row we're looking for is (scY + y) % 8.
					tile = p.getBackgroundTileRow(scX+byte(x+8), scY+y, lcdc)
				}
				pixelFifo.addTile(tile)
			}
		}

		// III. H-Blank mode (Horizontal blank) -- wait out the rest of the cycle.
		// ======================
		p.setMode(HBlank)
	}

	// update mode
	p.setMode(VBlank)
	return screen
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
	fifo []pixel
}

func (pf *pixelFifo) overlay(sprite []pixel) {
	// overlay the sprite's pixels on top of the leftmost 8 pixels in the fifo
	for i, px := range sprite {
		// Sprite color 0 is transparent -- don't overlay it.
		// TODO implement OBJ-to-BG priority (https://gbdev.gg8.se/wiki/articles/Video_Display#FF48_-_OBP0_-_Object_Palette_0_Data_.28R.2FW.29_-_Non_CGB_Mode_Only)
		if px.color != 0 {
			pf.fifo[i] = px
		}
	}
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

func toRGB(c color) []byte {
	switch c {
	case white:
		return []byte{241, 240, 239}
	case lightGray:
		return []byte{151, 150, 149}
	case darkGray:
		return []byte{76, 75, 74}
	case black:
		return []byte{0, 0, 255}
	}
	panic(fmt.Sprintf("toRGB: Got bad color: %v", c))
}
