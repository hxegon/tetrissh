package main

import (
	"fmt"
	"os"
	"tetrissh/tetris"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/* Just start a tetris game */

type Model struct {
	game  tea.Model
	style lipgloss.Style
}

func (m Model) Init() tea.Cmd {
	return tetris.FallTickCmd()
}

var score int

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tetris.FallMsg:
		m.game, cmd = m.game.Update(msg)
		return m, tea.Batch(cmd, tetris.FallTickCmd())
	case tea.WindowSizeMsg:
		h := msg.Height
		v := msg.Width
		m.style.Height(h)
		m.style.Width(v)

		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tetris.GameOverMsg:
		score = msg.Score
		return m, tea.Quit
	}

	m.game, cmd = m.game.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.style.Render(m.game.View())
}

func initialModel() Model {
	return Model{
		game: tetris.NewGameModel(),
		style: lipgloss.NewStyle().
			Padding(1).
			Background(lipgloss.Color("240")).
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Center),
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error entcountered: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Your score was %v!", score)
}
