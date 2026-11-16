//Braille Animation Player
// by Justin Ward w/ ChatGPT

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
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

type tickMsg struct {
	id uint64
}

func tick(id uint64) tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{id: id}
	})
}

type model struct {
	grid  Grid
	frame int

	directory    string
	files        []string
	loadedFrames []Image

	picker    filepicker.Model
	selecting bool

	playing    bool
	playbackID uint64

	error string
}

func initialModel() model {
	fp := filepicker.New()

	fp.DirAllowed = true
	fp.FileAllowed = true
	fp.ShowSize = false
	fp.CurrentDirectory, _ = os.Getwd()
	fp.SetHeight(12)

	return model{
		frame:     0,
		grid:      NewGrid(width, height, 0x00),
		picker:    fp,
		selecting: false,
		playing:   false,
	}
}

func (m model) Init() tea.Cmd {
	return m.picker.Init()
}

func (m *model) LoadAllFrames() {
	if len(m.files) == 0 {
		return
	}

	m.loadedFrames = []Image{}

	for f := range m.files {
		file := ReadImgFromFile(m.files[f])
		img := ParseImage(file)
		m.loadedFrames = append(m.loadedFrames, img)
	}
}

func (m *model) LoadFrame() {
	if len(m.loadedFrames) == 0 {
		return
	}

	m.grid.DrawImage(m.loadedFrames[m.frame])
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// Handle filepicker separately when selecting
	if m.selecting {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.selecting = false
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			}
		}

		var cmd tea.Cmd
		m.picker, cmd = m.picker.Update(msg)

		// Did we select a directory?
		if didSelect, path := m.picker.DidSelectFile(msg); didSelect {
			info, err := os.Stat(path)
			if err != nil {
				m.error = err.Error()
				return m, cmd
			}

			if !info.IsDir() {
				m.error = "Please select a directory, not a file"
				return m, cmd
			}

			files, err := LoadDirectory(path)
			if err != nil {
				m.error = err.Error()
				return m, cmd
			}

			if len(files) == 0 {
				m.error = "No .txt frames found in directory"
				return m, cmd
			}

			m.files = files
			m.directory = path
			m.frame = 0
			m.playing = false
			m.playbackID++
			m.error = ""

			m.LoadAllFrames()
			m.LoadFrame()
			m.selecting = false
		}

		return m, cmd
	}

	// Normal application controls
	switch msg := msg.(type) {
	case tickMsg:
		// Ignore outdated timers
		if msg.id != m.playbackID || !m.playing {
			return m, nil
		}

		// Advance to next frame
		if len(m.files) > 0 {
			m.frame = (m.frame + 1) % len(m.files)
			m.LoadFrame()
		}

		// Schedule the next tick
		return m, tick(msg.id)
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			if len(m.files) == 0 {
				break
			}

			m.playing = !m.playing
			m.playbackID++

			if m.playing {
				return m, tick(m.playbackID)
			}
		case "o":
			m.playing = false
			m.playbackID++

			fp := filepicker.New()
			fp.DirAllowed = true
			fp.FileAllowed = true
			fp.ShowSize = false
			fp.SetHeight(12)

			fp.CurrentDirectory, _ = os.Getwd()

			m.picker = fp
			m.error = ""
			m.selecting = true

			return m, m.picker.Init()

		case "right":
			if len(m.files) > 0 {
				m.frame = (m.frame + 1) % len(m.files)
				m.LoadFrame()
			}

		case "left":
			if len(m.files) > 0 {
				m.frame = (m.frame - 1 + len(m.files)) % len(m.files)
				m.LoadFrame()
			}

		case "r":
			m.frame = 0
			m.LoadFrame()

		case "c":
			m.grid = NewGrid(width, height, 0x00)

		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {

	if m.selecting {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			"Choose an animation directory:",
			"Current directory: "+m.picker.CurrentDirectory,
			"",
			m.picker.View(),
			"Press esc to cancel",
			m.error,
		)
	}

	playStatus := "Paused"

	if m.playing {
		playStatus = "Playing"
	}

	controls := "o      Open directory\n" +
		"space  Play / Pause\n" +
		"← →    Previous / Next frame\n" +
		"r      Restart animation\n" +
		"c      Clear display\n" +
		"q      Quit\n\n" +
		"Status: " + playStatus +
		"\nFrames: " + strconv.Itoa(len(m.files)) +
		"\nFrame: " + strconv.Itoa(m.frame+1)

	if m.directory != "" {
		controls += "\n\nDirectory: " + m.directory
	}

	if m.error != "" {
		controls += "\n\nError: " + m.error
	}

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
