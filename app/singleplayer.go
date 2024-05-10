package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type SinglePlayer struct {
	gm *GameModel
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
	if s.gm.GameOver {
		return fmt.Sprintf("Your final score is %v! Press q to go back to menu", s.gm.Score())
	} else {
		return s.gm.View()
	}
}
