package main

type Grid struct {
	width  int
	height int
	cells  []uint8
}

func NewGrid(width, height int, fill uint8) Grid {
	g := Grid{
		width:  width,
		height: height,
		cells:  make([]uint8, width*height),
	}

	for i := range g.cells {
		g.cells[i] = fill
	}

	return g
}

// Get the value of a cell.
func (g *Grid) At(x, y int) uint8 {
	return g.cells[y*g.width+x]
}

// Set the value of a cell.
func (g *Grid) Set(x, y int, value uint8) {
	g.cells[y*g.width+x] = value
}

// Toggle an individual Braille dot.
// Dot index ranges from 0 to 7.
func (g *Grid) ToggleDot(x, y, dot int) {
	g.cells[y*g.width+x] ^= 1 << dot
}

func (g *Grid) Width() int {
	return g.width
}

func (g *Grid) Height() int {
	return g.height
}
