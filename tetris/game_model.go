package tetris

import (
	"fmt"
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

type GameOverMsg struct {
	Score int
}

func GameOverCmd(score int) tea.Cmd {
	return func() tea.Msg {
		return GameOverMsg{Score: score}
	}
}

type GameModel struct {
	filledStyle lipgloss.Style
	emptyStyle  lipgloss.Style
	scoreStyle  lipgloss.Style
	game        Game
}

func NewGameModel() GameModel {
	return GameModel{
		game:        NewGame(20, 10, RandomPiece()),
		emptyStyle:  lipgloss.NewStyle().Background(lipgloss.Color("233")),
		filledStyle: lipgloss.NewStyle().Background(lipgloss.Color("240")),
		scoreStyle: lipgloss.NewStyle().
			Width(18).
			Border(lipgloss.NormalBorder(), true).
			AlignHorizontal(lipgloss.Center),
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
			cmd = GameOverCmd(m.game.score)
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

func toColor(ci int) lipgloss.Color {
	c := Color(ci)
	var code string

	switch c {
	case ColorEmpty:
		code = "240"
	case ColorRed:
		code = "001"
	case ColorBlue:
		code = "004"
	case ColorGreen:
		code = "002"
	case ColorOrange:
		code = "202"
	case ColorPurple:
		code = "129"
	case ColorYellow:
		code = "011"
	default:
		panic("Trying to convert an int to a color but there were no matching colors")
	}
	return lipgloss.Color(code)
}

func (m GameModel) View() string {
	var sb strings.Builder

	b := m.game.GetBoard()

	for y := range b {
		for x := range b[y] {
			val := b[y][x]
			color := lipgloss.Color(toColor(val))
			if val == 0 {
				sb.WriteString(m.emptyStyle.Foreground(color).Render("░░"))
			} else {
				sb.WriteString(m.filledStyle.Foreground(color).Render("🮑🮒"))
			}
		}

		if y != len(b)-1 { // Don't add a newline at the bottom
			sb.WriteString("\n")
		}
	}

	score := m.scoreStyle.Render(fmt.Sprintf("Score: %v", m.game.score))
	board := sb.String()
	return lipgloss.JoinVertical(lipgloss.Center, score, board)
}
