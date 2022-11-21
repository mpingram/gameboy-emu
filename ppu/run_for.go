package ppu

import "fmt"

// RunFor runs the PPU for a certain number of 4.14xxx MHz cycles.
// Cycles should always be divisible by 2.
func (p *PPU) RunFor(cycles int) {
	for i := 0; i < cycles; i++ {
		p.step()
	}
}

// Step executes 1 cycle's worth of work on the PPU.
func (p *PPU) step() {
	lcdstat := p.mem.Rb(LCDStatAddr)
	lastCycle := 455
	if lcdstat&LCDMode != VBlank {
		// In non-VBlank mode, iterate through OAMSearch -> Pixel Drawing -> HBLank modes,
		// drawing a scanline every 456 clocks.
		switch p.cycles {
		case 79:
			p.setMode(PixelDrawing)
			scanline := p.drawScanline(p.getLY(), p.getScrollX(), p.getScrollY())
			p.screen = append(p.screen, scanline...)
		case 159:
			p.setMode(HBlank)
		case lastCycle:
			if p.getLY() < 143 { // 143 is last onscreen scanline
				p.setMode(OAMSearch)
				p.setLY(p.getLY() + 1)
			} else if p.getLY() == 143 {
				p.setMode(VBlank)
				// send screen to output channel, but don't block
				select {
				case p.VideoOut <- p.screen:
				default:
				}
				p.screen = make([]byte, 0)
				p.setLY(p.getLY() + 1)
			} else {
				panic(fmt.Sprintf("LY is %v (>143), but mode is %v (should be VBlank)", p.getLY(), lcdstat&LCDMode))
			}
		}
	} else if lcdstat&LCDMode == VBlank {
		// In VBlank mode, increment LY every 456 clocks until LY == 153,
		// at which point re-enter OAMSearch mode and resume drawing scanlines.
		if p.cycles == lastCycle {
			if p.getLY() < 153 {
				p.setLY(p.getLY() + 1)
			} else if p.getLY() == 153 {
				p.setLY(0)
				p.setMode(OAMSearch)
			} else {
				panic(fmt.Sprintf("LY is %v, should be less than 153", p.getLY()))
			}
		}
	}
	p.cycles = (p.cycles + 1) % 456
}
