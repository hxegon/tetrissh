package tetris

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GameModel struct {
	filledStyle lipgloss.Style
	emptyStyle  lipgloss.Style
	game        Game
}

func NewGameModel() GameModel {
	return GameModel{
		game: NewGame(20, 10, RandomPiece()),
		emptyStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("240")).
			Foreground(lipgloss.Color("238")),
		filledStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("240")).
			Foreground(lipgloss.Color("112")),
	}
}

func (m GameModel) Init() tea.Cmd {
	return nil
}

func (m GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m GameModel) View() string {
	var sb strings.Builder

	b := m.game.GetBoard()

	for y := range b {
		for x := range b[y] {
			if b[y][x] == 0 {
				sb.WriteString(m.emptyStyle.Render("░░"))
			} else {
				sb.WriteString(m.filledStyle.Render("🮑🮒"))
			}
		}

		if y != len(b)-1 { // Don't add a newline at the bottom
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
