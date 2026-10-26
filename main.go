package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

type appModel struct {
	header           string
	activeFormula    formula
	mostRecentResult result
	previousResults  []result
	footer           string
}

type formula struct {
	name     string
	category string
	function func(float32) float32
}

type result struct {
	timestamp time.Time
	value     float32
	formula   formula
}

func (m appModel) Init() tea.Cmd {
	return nil
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case switchToMainMenu:
		m.activeModel = initialMainMenu()
		return m, m.activeModel.Init()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m appModel) View() tea.View {
	s := strings.Join([]string{
		m.header,
		m.activeModel.View().Content,
		m.footer,
	}, "\n")

	return tea.NewView(s)
}

func initialModel() appModel {
	return appModel{
		header:      "[Go Convert]\n\n",
		activeModel: MainMenuModel{},
		footer:      "Ctrl+C to quit",
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
