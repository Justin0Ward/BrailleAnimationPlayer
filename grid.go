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

// Get the value of a Braille cell.
func (g *Grid) At(x, y int) uint8 {
	return g.cells[y*g.width+x]
}

// Set the value of a Braille cell.
func (g *Grid) Set(x, y int, value uint8) {
	g.cells[y*g.width+x] = value
}

// Toggle an individual Braille dot within a cell.
// Dot index ranges from 0 to 7.
func (g *Grid) ToggleDot(x, y, dot int) {
	g.cells[y*g.width+x] ^= 1 << dot
}

func (g *Grid) SetDot(x, y, dot int, value bool) {
	if value {
		g.cells[y*g.width+x] |= 1 << dot
	} else {
		g.cells[y*g.width+x] &^= 1 << dot
	}
}

// DotWidth returns the width of the grid measured in individual Braille dots.
func (g *Grid) DotWidth() int {
	return g.width * 2
}

// DotHeight returns the height of the grid measured in individual Braille dots.
func (g *Grid) DotHeight() int {
	return g.height * 4
}

// dotLocation converts coordinates in the full dot field into a Braille
// cell coordinate and the bit index for the corresponding dot.
func (g *Grid) dotLocation(x, y int) (cellX, cellY, dot int, ok bool) {
	if x < 0 || x >= g.DotWidth() || y < 0 || y >= g.DotHeight() {
		return 0, 0, 0, false
	}

	cellX = x / 2
	cellY = y / 4

	// Braille bit positions are arranged differently from reading order:
	//       left  right
	// row 0   0      3
	// row 1   1      4
	// row 2   2      5
	// row 3   6      7
	dot = int([4][2]uint8{
		{0, 3},
		{1, 4},
		{2, 5},
		{6, 7},
	}[y%4][x%2])

	return cellX, cellY, dot, true
}

// DotAt reports whether the individual dot at (x, y) is raised.
// Coordinates are measured across the entire dot field, not in Braille cells.
// It returns false for coordinates outside the grid.
func (g *Grid) DotAt(x, y int) bool {
	cellX, cellY, dot, ok := g.dotLocation(x, y)
	if !ok {
		return false
	}

	return g.At(cellX, cellY)&(1<<dot) != 0
}

// SetDotAt sets the individual dot at (x, y) to value.
// It returns false if the coordinates are outside the grid.
func (g *Grid) SetDotAt(x, y int, value bool) bool {
	cellX, cellY, dot, ok := g.dotLocation(x, y)
	if !ok {
		return false
	}

	g.SetDot(cellX, cellY, dot, value)

	return true
}

// ToggleDotAt flips the individual dot at (x, y).
// It returns false if the coordinates are outside the grid.
func (g *Grid) ToggleDotAt(x, y int) bool {
	cellX, cellY, dot, ok := g.dotLocation(x, y)
	if !ok {
		return false
	}

	g.ToggleDot(cellX, cellY, dot)
	return true
}

func (g *Grid) Width() int {
	return g.width
}

func (g *Grid) Height() int {
	return g.height
}
