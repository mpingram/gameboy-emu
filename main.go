package main

import (
	"github.com/mpingram/gameboy-emu/frontend"
)

func main() {
	renderer := frontend.NewWebGLRenderer()
	fakeScreen := make([]byte, 0)

	for row := 0; row < 144; row++ {
		for col := 0; col < 160; col++ {
			fakeScreen = append(fakeScreen, 255, 0, 0)
		}
	}

	renderer.Render(fakeScreen)
}
