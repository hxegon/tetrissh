package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type SinglePlayer struct {
	gm *GameModel
	// A 0 value is a valid for a score, so use a pointer instead
	finalScore    *int // TODO: Move this into GameModel, can be used with multiplayer
	height, width int
}

func NewSinglePlayer() SinglePlayer {
	gm := NewGameModel()

	return SinglePlayer{
		gm: &gm,
	}
}

func (s SinglePlayer) Init() tea.Cmd {
	return FallTickCmd()
}

func (s SinglePlayer) Update(msg tea.Msg) (m tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case GameOverMsg:
		s.finalScore = &msg.Score
		return s, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return s, DeactivateCmd
		}
	}

	*s.gm, cmd = s.gm.Update(msg)
	return s, cmd
}

func (s SinglePlayer) View() string {
	if s.finalScore == nil {
		return s.gm.View()
	} else {
		return fmt.Sprintf("Your final score is %v! Press q to go back to menu", *s.finalScore)
	}
}
