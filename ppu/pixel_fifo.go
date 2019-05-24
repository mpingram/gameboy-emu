package ppu

import "fmt"

type pixelFifo struct {
	fifo []pixel
}

func (pf *pixelFifo) overlay(sprite []pixel) {
	// TODO Implement
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
