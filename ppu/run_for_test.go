package ppu

import (
	"testing"

	"github.com/mpingram/gameboy-emu/mmu"
)

func testSetup() (*PPU, *mmu.MMU) {
	m := mmu.New(mmu.MMUOptions{})
	p := New(m.PPUInterface)
	return p, m
}

func Test_RunForAdvancesCycles(t *testing.T) {
	tests := []struct {
		name           string
		cyclesIn       int
		cyclesToRunFor int
		cyclesOut      int
	}{
		{"0 + 1 -> 1", 0, 1, 1},
		{"1 + 1 -> 2", 1, 1, 2},
		{"1 + 80 -> 81", 1, 80, 81},
		{"wraparound: 455 + 1 -> 0", 455, 1, 0},
		{"wraparound: 0 + 456 -> 0", 0, 456, 0},
		{"wraparound: 1 + 456 -> 1", 1, 456, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, _ := testSetup()
			p.cycles = tt.cyclesIn
			p.RunFor(tt.cyclesToRunFor)
			if p.cycles != tt.cyclesOut {
				t.Errorf("Got cycles %v; expected %v", p.cycles, tt.cyclesOut)
			}
		})
	}
}

func Test_ModeSwitching(t *testing.T) {
	// This test isn't meant to test detailed timings, but instead to confirm that the timings
	// are roughly correct. This is mostly relevant when checking transitions from PixelDrawing -> HBlank,
	// as PixelDrawing can take anywhere from ~172 -> ~289 cycles.
	tests := []struct {
		name           string
		cyclesIn       int
		cyclesToRunFor int
		ly             byte
		modeIn         Mode
		modeOut        Mode
	}{
		{"79 + 1 -> 80 switches from OAMSearch -> PixelDrawing", 79, 1, 0, OAMSearch, PixelDrawing},
		{"80 + 290 -> 370 switches from PixelDrawing -> HBlank", 80, 290, 0, PixelDrawing, HBlank},
		{"455 + 1 -> 0 switches from HBlank -> OAMSearch if ly < 143", 455, 1, 142, HBlank, OAMSearch},
		{"455 + 1 -> 0 switches from HBlank -> VBlank if 143 => ly < 153", 455, 1, 143, HBlank, VBlank},
		{"455 + 1 -> 0 switches from VBlank -> OAMSearch if ly=153", 455, 1, 153, VBlank, OAMSearch},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, _ := testSetup()
			p.cycles = tt.cyclesIn
			p.setLY(tt.ly)
			p.setMode(tt.modeIn)
			p.RunFor(tt.cyclesToRunFor)
			if p.readLCDStat().Mode != tt.modeOut {
				t.Errorf("Got mode %v; expected %v", p.readLCDStat().Mode, tt.modeOut)
			}
		})
	}
}
