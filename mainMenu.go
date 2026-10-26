package main

import (
	tea "charm.land/bubbletea/v2"
)

type MainMenuModel struct {
	list []string
}

func (m MainMenuModel) Init() tea.Cmd {
	return nil
}

func (m MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m MainMenuModel) View() tea.View {
	s := "Test Main Menu"
	return tea.NewView(s)
}
