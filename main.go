// Its Braille not Brail

package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	width  = 32 / 2
	height = 32 / 4
)

var (
	controlStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Width(25)

	displayStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			Padding(0)
)

type Vector2 struct {
	x int
	y int
}

type model struct {
	cursor Vector2
	frame  int
	grid   Grid
}

func initialModel() model {
	return model{
		cursor: Vector2{x: 0, y: 0},
		frame:  1,
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

		case "d":
			frameCountStr := strconv.Itoa(m.frame)
			file := ReadImgFromFile("frames/" + frameCountStr + ".brlimg")

			img := ParseImage(file, 32, 32)
			m.grid.DrawImage(img)

			if m.frame < 3 {
				m.frame++
			}

		case "f":
			m.grid = NewGrid(32/2, 32/4, 0xFF)

		case "up":
			m.cursor.y--
		case "down":
			m.cursor.y++
		case "left":
			m.cursor.x--
		case "right":
			m.cursor.x++

		case " ":
			m.grid.SetDotAt(m.cursor.x, m.cursor.y, true)

		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	curStr := "cursor: " + strconv.Itoa(m.cursor.x) + "," + strconv.Itoa(m.cursor.y)
	curStr += " frame: " + strconv.Itoa(m.frame)
	controls := "Press q to quit"
	display := RenderGrid(&m.grid)

	menu := lipgloss.JoinVertical(
		lipgloss.Top,
		controlStyle.Render(curStr),
		controlStyle.Render(controls),
	)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		displayStyle.Render(display),
		menu,
	)
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
