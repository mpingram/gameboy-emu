package ppu

func (p *PPU) drawScreen() Screen {

	screen := make(Screen, int(screenHeight)*int(screenWidth)*3)

	for y := byte(0); y < screenHeight; y++ {

		lcdc := p.readLCDControl()
		stat := p.readLCDStat()
		wX := p.getWindowX()
		wY := p.getWindowY()
		scX := p.getScrollX()
		scY := p.getScrollY()

		scanline := p.drawScanline(y, lcdc, stat, wX, wY, scX, scY)
		screen = append(screen, scanline...)
	}

	return screen
}

func (p *PPU) drawScanline(
	y byte, // the y-coordinate of this scanline
	lcdc LCDControl,
	stat LCDStat,
	wX, // the x-coordinate of the window
	wY, // the y-coordinate of the window
	scX, // the x scroll amount (ie the x-coordinate of the screen top left)
	scY byte, // the y scroll amount (ie the y-coordinate of the screen top left)
) scanline {

	scanline := make([]byte, int(screenWidth)*3)

	// are we currently drawing the window (instead of the background)?
	drawingWindow := false

	// I. OAM Search
	// =====================
	p.setMode(OAMSearch)
	// get the first 10 sprites that are on this y-coordinate
	sprites := p.getOAMEntries(y, lcdc)

	// II. Pixel Drawing mode
	// =====================
	p.setMode(PixelDrawing)
	// Start at SCX-8. (This gives us room to draw sprites offscreen.)
	// Fill pixel fifo with the leftmost two tiles on this scanline
	pixelFifo := pixelFifo{}
	bgTileRow1 := p.getBackgroundTileRow(scX, scY, 0, y, lcdc)
	bgTileRow2 := p.getBackgroundTileRow(scX+8, scY, 0, y, lcdc)
	pixelFifo.addTile(bgTileRow1)
	pixelFifo.addTile(bgTileRow2)
	// dequeue pixel fifo XScroll % 8 times. This aligns our tiles to the edge of the screen.
	for i := byte(0); i < scX%8; i++ {
		pixelFifo.dequeue()
	}

	// For each pixel on this scanline:
	for x := byte(0); x < screenWidth; x++ {

		// Draw sprites
		if lcdc.SpriteEnable {
			// If a sprite starts at this xpos, layer that sprite on top of the bg tile.
			// Do this for every sprite that starts at this xpos.
			//
			// NOTE sprite X,Y coordinates are different from screen coordinates.
			// Sprite X coordinate starts 8px to left of screen, and sprite Y coordinate starts
			// 16px off the top of the screen. This is so sprites can be drawn partially off screen.
			//
			// Draw all partially-offscreen sprites.
			if isFirstTile := x == 0; isFirstTile {
				// Draw all the partially-offscreen sprites on the left of the screen.
				// We need to overlay these sprites with the tiles in the pixel fifo, but we need to shift the
				// sprite left by (8-sprite.x) pixels.
				for sprite := sprites[0]; sprite.x < 8; {
					// overlay this sprite, but shift it left by (8-x) pixels.
					pixels := p.getSpriteRow(sprite, y)
					shifted := shiftTileLeft(pixels, 8-sprite.x)
					pixelFifo.overlay(shifted)
					// pop this element off the sprites array
					sprites = sprites[1:]
				}
			} else if isLastTile := x == screenWidth-1-8; isLastTile {
				// Draw all the partially-offscreen sprites on the right of the screen.
				for sprite := sprites[0]; sprite.x > screenWidth; {
					// overlay this sprite, but shift it right by (???) pixels.
					pixels := p.getSpriteRow(sprite, y)
					shifted := shiftTileRight(pixels, 8-sprite.x)
					pixelFifo.overlay(shifted)
					// pop this element off the sprites array
					sprites = sprites[1:]
				}
			}

			// Draw the fully onscreen sprite(s) that start at this x-coordinate.
			for sprite := sprites[0]; sprite.x-8 == x; {
				pixels := p.getSpriteRow(sprite, y)
				pixelFifo.overlay(pixels)
				// pop this element off the sprites array
				sprites = sprites[1:]
			}
		}

		// Check for window start
		if lcdc.WindowEnable {
			// If window starts at this xpos and we're past the window ypos, clear fifo and start filling it with window tiles instead.
			if wY >= y && wX == x {
				drawingWindow = true
				pixelFifo.clear()
				windowTile := p.getWindowTileRow(scX, scY, x, y, lcdc)
				pixelFifo.addTile(windowTile)
			}
		}

		// pop a pixel off the pixel fifo, colorize it, and send it to the screen
		px, err := pixelFifo.dequeue()
		if err != nil {
			panic(err)
		}
		r, g, b := p.colorize(px)
		scanline = append(scanline, r, g, b)

		// Refill the pixel fifo if necessary
		if pixelFifo.size() <= 8 {
			// get next tile -- if we're drawing background, get the next background tile,
			// otherwise get the next window tile.
			var tile []pixel
			if drawingWindow {
				tile = p.getBackgroundTileRow(x, y, scX, scY, lcdc)
			} else {
				tile = p.getBackgroundTileRow(x, y, scX, scY, lcdc)

			}
			pixelFifo.addTile(tile)
		}
	}

	// III. H-Blank mode (Horizontal blank) -- wait out the rest of the cycle.
	// ======================
	p.setMode(HBlank)

	return scanline
}

// Screen is a byte array representing the colorized pixels
// of a gameboy screen. Its format is
//
// 	1 pixel
//  |-----|
// [R, G, B, R, G, B, R, G, B]
// Where R,G,B are one byte representing the reg, green, blue
// component of each pixel.
type Screen []byte // byte array of length 144 * 160 * 3, consisting

// scanLine is a byte array representing the colorized pixels
// of one row of the screen. Its format is the same as the format of Screen,
// but it is 160 * 3 bytes long.
type scanline []byte
