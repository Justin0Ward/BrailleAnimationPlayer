package main

import "strings"

// RenderCell converts one Braille cell (in cell coordinates) into a Unicode rune.
func RenderCell(g *Grid, x, y int) rune {
	return rune(0x2800 + int(g.At(x, y)))
}

// RenderGrid converts the complete grid into rows of Unicode Braille characters.
func RenderGrid(g *Grid) string {
	var b strings.Builder

	for y := 0; y < g.Height(); y++ {
		for x := 0; x < g.Width(); x++ {
			b.WriteRune(RenderCell(g, x, y))
		}

		b.WriteByte('\n')
	}

	return b.String()
}
