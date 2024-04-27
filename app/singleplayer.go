package app

import (
	"fmt"
	"tetrissh/tetris"

	tea "github.com/charmbracelet/bubbletea"
)

type SinglePlayer struct {
	gm *tetris.GameModel
	// A 0 value is potentially valid for a score, so use a pointer instead
	finalScore    *int
	height, width int
}

func NewSinglePlayer() SinglePlayer {
	gm := tetris.NewGameModel()

	return SinglePlayer{
		gm: &gm,
	}
}

func (s SinglePlayer) Init() tea.Cmd {
	return tetris.FallTickCmd()
}

func (s SinglePlayer) Update(msg tea.Msg) (m tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tetris.GameOverMsg:
		s.finalScore = &msg.Score
		return s, nil
	case tea.KeyMsg:
		if msg.String() == "q" {
			return s, DeactivateCmd
		}
	}

	newm, cmd := s.gm.Update(msg)
	if newg, ok := newm.(tetris.GameModel); ok {
		s.gm = &newg
		return s, cmd
	}

	// TODO: Replace with proper error handling
	// Should never happen
	panic("Couldn't coerce GameModel update value to GameModel???")
}

func (s SinglePlayer) View() string {
	if s.finalScore == nil {
		return s.gm.View()
	} else {
		return fmt.Sprintf("Your final score is %v! Press q to go back to menu", *s.finalScore)
	}
}
