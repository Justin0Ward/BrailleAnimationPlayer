// Its Braille not Brail

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

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

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type model struct {
	grid      Grid
	frame     int
	directory string
	files     []string
	playing   bool
}

func initialModel() model {
	return model{
		grid:    NewGrid(width, height, 0x00),
		playing: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if !m.playing || len(m.files) == 0 {
			return m, nil
		}

		m.frame = (m.frame + 1) % len(m.files)

		file := ReadImgFromFile(m.files[m.frame])
		img := ParseImage(file)
		m.grid = NewGrid(width, height, 0x00)
		m.grid.DrawImage(img)

		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			m.playing = !m.playing

			if m.playing {
				return m, tick()
			}

		case "t":
			files, err := LoadDirectory("frames/spin")
			if err != nil {
				return m, nil
			}

			m.directory = "frames/spin"
			m.files = files
			m.frame = 0

			file := ReadImgFromFile(m.files[m.frame])
			img := ParseImage(file)
			m.grid.DrawImage(img)

		case "r":
			m.frame = 0

		case "right":
			frameCountStr := strconv.Itoa(m.frame + 1)
			file := ReadImgFromFile("frames/spin/frame_00" + frameCountStr + ".txt")

			img := ParseImage(file)
			m.grid.DrawImage(img)

			if m.frame < 5 {
				m.frame++
			} else if m.frame >= 5 {
				m.frame = 0
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
