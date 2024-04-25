package tetris

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FallMsg struct{}

func FallTickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return FallMsg{}
	})
}

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

func (m GameModel) Update(msg tea.Msg) (GameModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case FallMsg:
		m.game.Fall()
		if m.game.GameOver {
			cmd = tea.Quit
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			m.game.Act(ActionLeft)
		case "l", "right":
			m.game.Act(ActionRight)
		case "j", "down":
			m.game.Act(ActionDown)
		case "k", "r", "up":
			m.game.Act(ActionRotate)
		case " ":
			m.game.Act(ActionDrop)
		}
	}

	return m, cmd
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
