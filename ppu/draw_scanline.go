package ppu

func (p *PPU) drawScanline(ly, scX, scY byte) []Color {
	// NOTE this implementation currently completely ignores the Window and sprites.
	var scanline []Color
	y := scY + ly // y is the global y-coordinate of the current scanline.

	// Initialize the pixel fifo with pixels from the tile that intersects with scX.
	pixelFifo := pixelFifo{}
	bgTile := p.getBackgroundTileRow(scX-(scX%8), y)
	pixelFifo.addTile(bgTile)

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
			palette := p.getBGPalette()
			color := palette[px.colorNumber]
			scanline = append(scanline, color)
		}
	}
	return scanline
}
