package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AppModel struct {
	modeModel tea.Model
	style     lipgloss.Style
	bgColor   lipgloss.Color
	height    int
	width     int
}

func NewAppModel(r *lipgloss.Renderer) AppModel {
	bgColor := lipgloss.Color("232")
	var menu tea.Model = NewMenuModel()
	return AppModel{
		style: r.NewStyle().
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Background(bgColor),
		// Margin(1), // Why does this cause a "ghosting" issue?
		modeModel: menu,
	}
}

func (a AppModel) Init() tea.Cmd {
	return nil
}

func (a AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h := msg.Height
		w := msg.Width
		a.style.Height(h)
		a.style.Width(w)
		a.height = h
		a.width = w
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit
		case "q":
			// TODO: This should probably be handled through the standard mode activation
			// so an app specific version of GameModel that handles "q" and can send back a score
			// like for a game over message? or to record? Or maybe GameModel should have a gameover screen.
			a.modeModel = NewMenuModel()
			cmd = func() tea.Msg {
				return tea.WindowSizeMsg{
					Height: a.height,
					Width:  a.width,
				}
			}

			return a, cmd
		}
	case ModeActivateMsg:
		a.modeModel = msg.newModel
		cmd = msg.cmd
		return a, cmd
	}

	// By default, pass through the message to modeModel
	a.modeModel, cmd = a.modeModel.Update(msg)
	return a, cmd
}

func (a AppModel) View() string {
	var view string
	switch model := a.modeModel.(type) {
	// case MenuModel:
	default:
		view = model.View()
	}
	return a.style.Render(view)
}
