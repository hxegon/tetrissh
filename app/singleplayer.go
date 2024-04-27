package app

import (
	"tetrissh/tetris"

	tea "github.com/charmbracelet/bubbletea"
)

type SinglePlayer struct {
	*tetris.GameModel
	height, width int
}

func NewSinglePlayer() SinglePlayer {
	gm := tetris.NewGameModel()

	return SinglePlayer{
		GameModel: &gm,
	}
}

func (s SinglePlayer) Update(msg tea.Msg) (m tea.Model, cmd tea.Cmd) {
	if _, ok := msg.(tetris.GameOverMsg); ok {
		return NewMenuModel(), nil
	}

	newm, cmd := s.GameModel.Update(msg)
	if newg, ok := newm.(tetris.GameModel); ok {
		s.GameModel = &newg
		return s, cmd
	}

	// Should never happen
	panic("Couldn't coerce GameModel update value to GameModel???")
}
