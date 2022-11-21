package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	"io"

	"github.com/mpingram/gameboy-emu/cpu"
	"github.com/mpingram/gameboy-emu/mmu"
	"github.com/mpingram/gameboy-emu/ppu"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Input struct {
	esc bool
}
type GameBoy struct {
	mmu *mmu.MMU
	cpu *cpu.CPU
	ppu *ppu.PPU

	breakpoint        uint16
	breakpointEnabled bool

	screen  *ebiten.Image
	screenW int
	screenH int
}

func NewGameBoy(bootRom, gameRom io.Reader, breakpoint uint16) *GameBoy {

	gb := &GameBoy{}

	gb.breakpoint = breakpoint
	gb.breakpointEnabled = breakpoint != 0

	gb.mmu = mmu.New(mmu.MMUOptions{BootRom: bootRom, GameRom: gameRom})
	gb.ppu = ppu.New(gb.mmu)
	gb.cpu = cpu.New(gb.mmu)

	gb.screen = ebiten.NewImage(160, 144)
	gb.screenW = 160
	gb.screenH = 144

	return gb
}

func (gb *GameBoy) Draw(screen *ebiten.Image) {

	rawImage := <-gb.ppu.VideoOut
	for i, px := range rawImage {
		var col color.Color
		switch px {
		case ppu.White:
			col = color.White
		case ppu.LightGray:
			col = color.White
		case ppu.DarkGray:
			col = color.Black
		case ppu.Black:
			col = color.Black
		}
		x := i % gb.screenW
		y := i / gb.screenH
		gb.screen.Set(x, y, col)
	}

	screen.DrawImage(gb.screen, &ebiten.DrawImageOptions{})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("0x%0x  FPS: %f", gb.cpu.PC, ebiten.CurrentFPS()))
}

func (gb *GameBoy) Update() error {
	MAX_CYCLES_PER_FRAME := 70224

	// execute up to MAX_CYCLES_PER_FRAME cycles
	cycles_this_frame := 0
	for cycles_this_frame < MAX_CYCLES_PER_FRAME {
		_, cycles := gb.cpu.Step()
		gb.ppu.RunFor(cycles)
		cycles_this_frame += cycles
	}

	return nil
}

func (gb *GameBoy) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 160, 144
}

//go:embed roms/boot/DMG_ROM.bin
var bootRom []byte

//go:embed roms/tetris.gb
var gameRom []byte

func main() {
	gb := NewGameBoy(bytes.NewReader(bootRom), bytes.NewReader(gameRom), 0)
	ebiten.RunGame(gb)
}
