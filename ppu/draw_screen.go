package ppu

func (p *PPU) drawScreen() Screen {

	screen := make(Screen, int(screenHeight)*int(screenWidth)*3)

	for y := byte(0); y < screenHeight; y++ {

		// re-read these registers before every scanline
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

	// I. OAM Search
	// =====================
	p.setMode(OAMSearch)
	// get the first 10 sprites that are on this y-coordinate
	sprites := p.getOAMEntries(y, lcdc)

	// II. Pixel Drawing mode
	// =====================
	p.setMode(PixelDrawing)

	drawingWindow := false

	// 1. Fill the pixel fifo with pixels from a) the tile that intersects with SCX-8 and b) the tile to its right.
	// Note that we often have extra pixels from tile A in the pixel fifo when SCX-8 cuts into the middle of a tile.
	pixelFifo := pixelFifo{}
	// scY is an absolute coordinate (based on 0,0, just like the background tile map),
	// while y is a screen-space coordinate of the current scanline (based on scY.)
	// So, scY + y is the absolute coordinate of the current scanline. Because
	// absolute coordinates have one tile every 8px, (scY + y) % 8 gets us the row
	// of the tile that the current scanline is on.
	row := (scY + y) % 8
	bgTileRow1 := p.getBackgroundTileRow((scX - 8), scY+y, row, lcdc)
	bgTileRow2 := p.getBackgroundTileRow((scX-8)+8, scY+y, row, lcdc)
	pixelFifo.addTile(bgTileRow1)
	pixelFifo.addTile(bgTileRow2)

	// 2. Discard pixels from pixel fifo until it lines up with SCX-8. (Dequeue (SCX-8) % 8 times)
	for i := byte(0); i < scX%8; i++ {
		pixelFifo.dequeue()
	}

	// 3. For x = 0 up to 159+8: (screen width is 160px, and we started 8px to the left of the screen edge)
	for x := byte(0); x < screenWidth+8; x++ {

		// 3a. Check if the window starts at this pixel. (Because the window's x coordinate
		// is relative to the screen, and we started 8px to the left of the screen, we
		// don't need to worry about the edge case of wX (window.X) being less than x.)
		if lcdc.WindowEnable {
			// If window starts at this xpos and we're past the window ypos, clear fifo and start filling it with window tiles instead.
			if wY <= y && wX == x {
				drawingWindow = true
				pixelFifo.clear()
				// Because window tiles are aligned to the screen, the row we're looking
				// for is just y (current screen coordinate) mod 8.
				row := y % 8
				windowTileRowA := p.getWindowTileRow(x, y, row, lcdc)
				windowTileRowB := p.getWindowTileRow(x, y, row, lcdc)
				pixelFifo.addTile(windowTileRowA)
				pixelFifo.addTile(windowTileRowB)
			}
		}

		// 3b. Overlay all sprites that start on this pixel onto the pixel fifo,
		if lcdc.SpriteEnable {
			// For every sprite that starts at this xpos, overlay that sprite on top of
			// the pixel fifo.
			for sprite := sprites[0]; sprite.x == x; {
				// sprite.y is the screen coordinate of the top of the sprite, and we
				// know that sprite.y <= y (because the sprite is on this row). So we can
				// get the row of the sprite by y - sprite.y. E.g., if sprite.y=20 and y=26,
				// the row we want is 26 - 20 == 6.
				row := y - sprite.y
				pixels := p.getSpriteRow(sprite, row)
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
		r, g, b := p.colorize(px)
		// Draw the pixel to the screen.
		scanline = append(scanline, r, g, b)

		// 3d. If the pixel fifo contains only one tile, add the next tile to it.
		if pixelFifo.size() <= 8 {
			// Get next tile -- if we're drawing background, get the next background tile,
			// otherwise get the next window tile.
			var tile []pixel
			if drawingWindow {
				tile = p.getWindowTileRow(x+8, y, y%screenHeight, lcdc)
			} else {
				tile = p.getBackgroundTileRow((scX-8)+8, scY+y, (scY+y)%8, lcdc)
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
