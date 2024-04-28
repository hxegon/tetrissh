package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	bgColor  = lipgloss.Color("232")
	appStyle = lipgloss.NewStyle().
			Background(bgColor)
)

type DeactivateMsg struct{}

func DeactivateCmd() tea.Msg {
	return DeactivateMsg{}
}

type AppModel struct {
	menu          tea.Model
	selectedModel tea.Model
}

func NewAppModel(r *lipgloss.Renderer) AppModel {
	return AppModel{
		menu: NewMenuModel(),
	}
}

func (a AppModel) Init() tea.Cmd {
	return nil
}

func (a AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Could maybe deal with propogating window sizes to components through a pointer?
		appStyle.Height(msg.Height)
		appStyle.Width(msg.Width)
		a.menu, _ = a.menu.Update(msg)
	case MenuSelectMsg:
		a.selectedModel = msg.model
		return a, a.selectedModel.Init()
	case DeactivateMsg:
		if a.selectedModel != nil {
			a.selectedModel = nil // drop *tea.Model contents
		} else {
			// Else we are completely closing the app
			return a, tea.Quit
		}
	}

	// Route update to menu unless there's a selectedModel
	if a.selectedModel == nil {
		a.menu, cmd = a.menu.Update(msg)
	} else {
		a.selectedModel, cmd = a.selectedModel.Update(msg)
	}

	return a, cmd
}

func (a AppModel) View() string {
	if a.selectedModel == nil { // No model selected
		return a.menu.View()
	}
	return a.selectedModel.View()
}
