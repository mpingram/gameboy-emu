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
	lastCycle := 455

	lcdstat := p.readLCDStat()

	if lcdstat.Mode != VBlank {
		// In non-VBlank mode, iterate through OAMSearch -> Pixel Drawing -> HBLank modes,
		// drawing a scanline every 456 clocks.
		switch p.cycles {
		case 0:
			p.setMode(OAMSearch)
			if p.readLCDStat().Mode != OAMSearch {
				panic(fmt.Sprintf("Mode should be OAMSearch; Got %v", p.readLCDStat().Mode))
			}
		case 80:
			p.setMode(PixelDrawing)
			if p.readLCDStat().Mode != PixelDrawing {
				panic(fmt.Sprintf("Mode should be PixelDrawing; Got %v", p.readLCDStat().Mode))
			}
			scanline := p.drawScanline(p.getLY(), p.getScrollX(), p.getScrollY())
			p.screen = append(p.screen, scanline...)
		case 160:
			p.setMode(HBlank)
			if p.readLCDStat().Mode != HBlank {
				panic(fmt.Sprintf("Mode should be HBlank; Got %v", p.readLCDStat().Mode))
			}
		case lastCycle:
			if p.getLY() < 144 { // 143 is last onscreen scanline
				p.setMode(OAMSearch)
				p.setLY(p.getLY() + 1)
			} else if p.getLY() == 144 {
				p.setMode(VBlank)
				if p.readLCDStat().Mode != VBlank {
					panic(fmt.Sprintf("Mode should be VBlank; Got %v", p.readLCDStat().Mode))
				}
				// send screen to output channel
				p.videoOut <- p.screen
				p.screen = make([]Color, 0)
				p.setLY(p.getLY() + 1)
			} else {
				panic(fmt.Sprintf("LY is %v (>143), but mode is %v (should be VBlank)", p.getLY(), lcdstat.Mode))
			}
		}
		return

	} else if lcdstat.Mode == VBlank {
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
		return
	}
}
