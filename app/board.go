package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type boardStyle struct {
	emptyStyle lipgloss.Style
	blockStyle lipgloss.Style
}

func defaultBoardStyle() boardStyle {
	return boardStyle{
		emptyStyle: lipgloss.NewStyle().Background(lipgloss.Color("233")),
		blockStyle: lipgloss.NewStyle().Background(lipgloss.Color("240")),
	}
}

// TODO: Extract board type [][]int
func RenderBoard(b [][]int, s boardStyle) string {
	var sb strings.Builder

	for y := range b {
		for x := range b[y] {
			val := b[y][x]
			color := lipgloss.Color(toColor(val))
			if val == 0 {
				sb.WriteString(s.emptyStyle.Foreground(color).Render("░░"))
			} else {
				sb.WriteString(s.blockStyle.Foreground(color).Render("🮑🮒"))
			}
		}

		if y != len(b)-1 { // Don't add a newline at the bottom
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
