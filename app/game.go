package app

import (
	"fmt"
	"tetrissh/tetris"
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
	boardStyle boardStyle
	scoreStyle lipgloss.Style
	game       tetris.Game
}

func NewGameModel() GameModel {
	return GameModel{
		game:       tetris.NewGame(20, 10, tetris.RandomPiece()),
		boardStyle: defaultBoardStyle(),
		scoreStyle: lipgloss.NewStyle().
			Width(18).
			Border(lipgloss.NormalBorder(), true).
			AlignHorizontal(lipgloss.Center),
	}
}

func (m GameModel) Init() tea.Cmd {
	return FallTickCmd()
}

func (m GameModel) Update(msg tea.Msg) (GameModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case FallMsg:
		m.game.Fall()
		if m.game.GameOver {
			cmd = GameOverCmd(m.game.Score)
		} else {
			cmd = FallTickCmd()
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			m.game.Act(tetris.ActionLeft)
		case "l", "right":
			m.game.Act(tetris.ActionRight)
		case "j", "down":
			m.game.Act(tetris.ActionDown)
		case "k", "r", "up":
			m.game.Act(tetris.ActionRotate)
		case " ":
			m.game.Act(tetris.ActionDrop)
		}
	}

	return m, cmd
}

func toColor(ci int) lipgloss.Color {
	c := tetris.Color(ci)
	var code string

	switch c {
	case tetris.ColorEmpty:
		code = "240"
	case tetris.ColorRed:
		code = "001"
	case tetris.ColorBlue:
		code = "004"
	case tetris.ColorGreen:
		code = "002"
	case tetris.ColorOrange:
		code = "202"
	case tetris.ColorPurple:
		code = "129"
	case tetris.ColorYellow:
		code = "011"
	default:
		// TODO: Don't use panic for error handling, maybe an Error()?
		panic("Trying to convert an int to a color but there were no matching colors")
	}
	return lipgloss.Color(code)
}

func (m GameModel) View() string {
	board := RenderBoard(m.game.Board(), m.boardStyle)

	score := m.scoreStyle.Render(fmt.Sprintf("Score: %v", m.game.Score))
	return lipgloss.JoinVertical(lipgloss.Center, score, board)
}
