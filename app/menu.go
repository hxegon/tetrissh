package app

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuSelectMsg struct {
	model tea.Model
}

type MenuItem struct {
	newModel func() tea.Model
	cmd      func() tea.Cmd
	title    string
	desc     string
}

func (m MenuItem) Title() string       { return m.title }
func (m MenuItem) Description() string { return m.desc }
func (m MenuItem) FilterValue() string { return m.title }

func (m MenuItem) SelectCmd() tea.Cmd {
	return func() tea.Msg {
		return MenuSelectMsg{model: m.newModel()}
	}
}

type MenuModel struct {
	list  list.Model
	style lipgloss.Style
}

func NewMenuModel() MenuModel {
	options := []list.Item{
		MenuItem{
			title: "Start",
			desc:  "Single player",
			newModel: func() tea.Model {
				return NewSinglePlayer()
			},
		}, MenuItem{
			title: "VS",
			desc:  "Multiplayer",
			newModel: func() tea.Model {
				return NewMultiplayer()
			},
		},
	}

	list := list.New(options, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Menu"

	return MenuModel{
		list:  list,
		style: lipgloss.NewStyle(),
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := m.style.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected := m.list.SelectedItem().(MenuItem)
			cmd = selected.SelectCmd()

			return m, cmd
		case "q", "ctrl+c":
			return m, DeactivateCmd
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m MenuModel) View() string {
	// FIXME: Styling is center aligned with m.style.Render()
	// return m.style.Render(m.list.View())
	return m.list.View()
}
