package app

import (
	"fmt"
	"strings"
	"tetrissh/tetris"

	"github.com/charmbracelet/lipgloss"
)

var (
	emptyStyle = lipgloss.NewStyle().Background(lipgloss.Color("233"))
	blockStyle = lipgloss.NewStyle().Background(lipgloss.Color("240"))
	scoreStyle = lipgloss.NewStyle().
			Width(18).
			Border(lipgloss.NormalBorder(), true).
			AlignHorizontal(lipgloss.Center)
)

type GameInfo interface {
	Board() [][]int
	Score() int
	// TODO: GameState?
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
		// TODO: Only panic in development mode
		panic("Trying to convert an int to a color but there were no matching colors")
	}
	return lipgloss.Color(code)
}

func BoardView(g GameInfo) string {
	var sb strings.Builder

	b := g.Board()

	for y := range b {
		for x := range b[y] {
			val := b[y][x]
			color := lipgloss.Color(toColor(val))
			if val == 0 {
				sb.WriteString(emptyStyle.Foreground(color).Render("░░"))
			} else {
				sb.WriteString(blockStyle.Foreground(color).Render("🮑🮒"))
			}
		}

		if y != len(b)-1 { // Don't add a newline at the bottom
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func ScoreView(g GameInfo) string {
	return scoreStyle.Render(fmt.Sprintf("Score: %v", g.Score()))
}
