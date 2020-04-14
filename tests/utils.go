package tests

// Test tile A:
// _  _  3  3  3  _  _  _
// _  3  3  3  3  3  _  _
// _  3  _  _  _  3  _  _
// -  3  _  _  _  3  _  -
// -  3  3  3  3  3  _  -
// -  3  _  _  _  3  _  -
// _  3  _  _  _  3  _  _
// _  _  _  _  _  _  _  _
var TestTileA = []byte{
	0, 0, 0, 3, 3, 0, 0, 0, // row 0
	0, 0, 3, 3, 3, 3, 0, 0, // row 1
	0, 3, 3, 0, 0, 3, 3, 0, // row 2
	0, 3, 3, 0, 0, 3, 3, 0, // row 3
	0, 3, 0, 0, 0, 0, 3, 0, // row 4
	0, 3, 3, 3, 3, 3, 3, 0, // row 5
	0, 3, 0, 0, 0, 0, 3, 0, // row 6
	3, 3, 0, 0, 0, 0, 3, 0, // row 7
}
var PackedTestTileA = PackTile(TestTileA)

// Test tile Y:
// _  3  _  _  _  _  3  _
// _  _  3  _  _  3  _  _
// _  _  _  3  3  _  _  _
// _  _  _  2  3  _  _  _
// _  _  _  3  2  _  _  _
// _  _  _  3  1  _  _  _
// _  _  _  3  1  _  _  _
// _  _  _  _  _  _  _  _
var TestTileY = []byte{
	0, 3, 0, 0, 0, 0, 3, 0, // row 0
	0, 0, 3, 0, 0, 3, 0, 0, // row 1
	0, 0, 0, 3, 3, 0, 0, 0, // row 2
	0, 0, 0, 2, 3, 0, 0, 0, // row 3
	0, 0, 0, 3, 2, 0, 0, 0, // row 4
	0, 0, 0, 3, 1, 0, 0, 0, // row 5
	0, 0, 0, 3, 1, 0, 0, 0, // row 6
	0, 0, 0, 0, 0, 0, 0, 0, // row 7
}
var PackedTestTileY = PackTile(TestTileY)

// Test tile Checkerboard:
// _  1  _  2  _  1 _  1
// 1  _  2  _  2  _  1  _
// _  1  _  2  _  1  _  1
// 1  _  2  _  2  _  1  _
// _  1  _  2  _  3  _  3
// 1  _  2  _  3  _  3  _
// _  1  _  2  _  3  _  3
// 1  _  2  _  3  _  3  _
var TestTileCheckered = []byte{
	0, 2, 0, 2, 0, 2, 0, 2, // row 0
	2, 0, 2, 0, 2, 0, 2, 0, // row 0
	0, 2, 0, 2, 0, 2, 0, 2, // row 0
	2, 0, 2, 0, 2, 0, 2, 0, // row 0
	0, 2, 0, 2, 1, 3, 1, 3, // row 0
	2, 0, 2, 0, 3, 1, 3, 1, // row 0
	0, 2, 0, 2, 1, 3, 1, 3, // row 0
	2, 0, 2, 0, 3, 1, 3, 1, // row 0
}

var PackedTestTileCheckered = PackTile(TestTileCheckered)

// PackTile formats a background tile for the Gameboy VRAM.
// The input format is an 8x8 array of pallete numbers (0-3).
// The GB represents tiles in a 'packed' 16-byte format,
// aka "2 Bits per pixel" / 2bpp.
// Each row of the sprite corresponds to two bytes:
// the first byte contains the low bits of the 8 palette numbers,
// on the row, and the second byte contains the high bits.
func PackTile(tile []byte) []byte {
	packed := make([]byte, 0, 16)
	// for each row of the tile from top (0) to bottom (7)
	for r := 0; r < 8; r++ {
		row := tile[r*8 : (r*8)+8]
		// add the low bit of each palette number to the low byte,
		// and the high bit to the high byte.
		var lo, hi byte
		for i, px := range row {
			hi |= ((px & 0b10) >> 1) << (7 - i)
			lo |= (px & 0b01) << (7 - i)
		}
		packed = append(packed, lo, hi)
	}
	return packed
}
