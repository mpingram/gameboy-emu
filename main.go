package main

import (
	"github.com/mpingram/gameboy-emu/frontend"
)

func main() {
	renderer := frontend.NewWebGLRenderer()
	fakeScreen := make([]byte, 0)

	var red, green, blue byte
	for row := 0; row < 144; row++ {
		red++
		for col := 0; col < 160; col++ {
			green++
			fakeScreen = append(fakeScreen, red, green, blue)
		}
	}

	renderer.Render(fakeScreen)
}
