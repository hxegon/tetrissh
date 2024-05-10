package app

import (
	"tetrissh/tetris"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Core bubbletea model that wraps the tetris game as thinly as possible.
type GameModel struct {
	*tetris.Game
}

func NewGameModel() GameModel {
	t := tetris.NewGame(20, 10, tetris.RandomPiece())

	return GameModel{Game: &t}
}

type FallMsg struct{}

func FallTickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return FallMsg{} })
}

func (m GameModel) Init() tea.Cmd { return FallTickCmd() }

func (m GameModel) Update(msg tea.Msg) (GameModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case FallMsg:
		m.Fall()
		if !m.GameOver {
			cmd = FallTickCmd()
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			m.Act(tetris.ActionLeft)
		case "l", "right":
			m.Act(tetris.ActionRight)
		case "j", "down":
			m.Act(tetris.ActionDown)
		case "k", "r", "up":
			m.Act(tetris.ActionRotate)
		case " ":
			m.Act(tetris.ActionDrop)
		}
	}

	return m, cmd
}

func (m GameModel) View() string {
	board := BoardView(m)
	score := ScoreView(m)

	return lipgloss.JoinVertical(lipgloss.Center, score, board)
}
