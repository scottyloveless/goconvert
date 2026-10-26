package main

import (
	tea "charm.land/bubbletea/v2"
)

type MainMenuModel struct {
	list []string
}

type switchToMainMenu struct{}

func (m MainMenuModel) Init() tea.Cmd {
	return nil
}

func (m MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m MainMenuModel) View() tea.View {
	s := "Test Main Menu"
	return tea.NewView(s)
}

func initialMainMenu() tea.Model {
	return MainMenuModel{
		list: []string{
			"test",
			"again",
			"main menu",
		},
	}
}
