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

type model struct {
	grid Grid
}

func initialModel() model {
	return model{
		grid: NewGrid(width, height, 0xFF),
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
			x := rand.IntN(m.grid.Width())
			y := rand.IntN(m.grid.Height())
			dot := rand.IntN(7)

			m.grid.ToggleDot(x, y, dot)

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
