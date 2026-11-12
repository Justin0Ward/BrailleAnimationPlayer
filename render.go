package main

import (
	"strings"
)

func RenderGrid(g *Grid) string {
	var b strings.Builder

	for y := 0; y < g.Height(); y++ {
		for x := 0; x < g.Width(); x++ {

			// Convert the cell's bitmask to Unicode Braille.
			char := rune(0x2800 + int(g.At(x, y)))

			b.WriteRune(char)
		}

		b.WriteByte('\n')
	}

	return b.String()
}
