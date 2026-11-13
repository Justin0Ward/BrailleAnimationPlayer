// Its Braille not Brail

package main

import (
	"fmt"
	"math/rand/v2"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	width  = 16
	height = 8
)

type Vector2 struct {
	x int
	y int
}

type model struct {
	cursor Vector2
	grid   Grid
}

func initialModel() model {
	return model{
		cursor: Vector2{x: 0, y: 0},
		grid:   NewGrid(width, height, 0x00),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			dwidth := m.grid.DotWidth()
			dheight := m.grid.DotHeight()

			dotx := rand.IntN(dwidth)
			doty := rand.IntN(dheight)

			m.grid.ToggleDotAt(dotx, doty)

		case "up":
			m.cursor.y--
		case "down":
			m.cursor.y++
		case "left":
			m.cursor.x--
		case "right":
			m.cursor.x++

		case " ":
			m.grid.ToggleDotAt(m.cursor.x, m.cursor.y)

		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	return RenderGrid(&m.grid) + "\nPress q to quit"
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
