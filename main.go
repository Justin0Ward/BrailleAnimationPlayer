// Its Braille not Brail

package main

import (
	"fmt"
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
			Width(35).
			Foreground(lipgloss.Color("118"))

	displayStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			Padding(0).
			Foreground(lipgloss.Color("118"))
)

type model struct {
	frame int
	grid  Grid
}

func initialModel() model {
	return model{
		frame: 1,
		grid:  NewGrid(width, height, 0x00),
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
			m.frame = 1

		case "right":
			frameCountStr := strconv.Itoa(m.frame)
			file := ReadImgFromFile("frames/spin/frame_00" + frameCountStr + ".txt")

			img := ParseImage(file)
			m.grid.DrawImage(img)

			if m.frame < 6 {
				m.frame++
			} else if m.frame >= 6 {
				m.frame = 1
			}

		case "c":
			m.grid = NewGrid(width, height, 0x00)

		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	controls := "Press -> for next frame\n" +
		"Press c to clear display\n" +
		"Press r to set frame to first\n" +
		"Press q to quit"
	display := RenderGrid(&m.grid)

	menu := lipgloss.JoinVertical(
		lipgloss.Top,
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
