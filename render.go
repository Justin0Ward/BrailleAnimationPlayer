package main

import (
	"strconv"
	"strings"
)

const (
	testImage = `
0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
0 0 0 0 0 0 1 1 1 0 0 0 0 0 0 0
0 0 0 0 0 1 0 0 0 1 0 0 0 0 0 0
0 0 0 0 0 1 0 0 0 1 0 0 0 0 0 0
0 0 0 0 0 1 0 0 0 1 0 0 0 0 0 0
0 0 0 0 0 0 1 1 1 0 0 0 0 0 0 0
0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0
0 0 0 1 1 1 1 1 1 1 1 1 1 0 0 0
0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0
0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0
0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0
0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0
0 0 0 0 0 0 1 0 1 0 0 0 0 0 0 0
0 0 0 0 0 1 0 0 0 1 0 0 0 0 0 0
0 0 0 0 1 0 0 0 0 0 1 0 0 0 0 0
0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0
0 0 1 0 0 0 0 0 0 0 0 0 1 0 0 0
0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
`
)

type Image struct {
	width  int
	height int
	pixels []uint8
}

func ParseImage(data string, width, height int) Image {
	fields := strings.Fields(data)

	img := Image{
		width:  width,
		height: height,
		pixels: make([]uint8, len(fields)),
	}

	for i, field := range fields {
		n, err := strconv.Atoi(field)
		if err != nil {
			panic(err)
		}

		img.pixels[i] = uint8(n)
	}

	return img
}

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
