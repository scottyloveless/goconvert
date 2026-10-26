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
	activeModel      tea.Model
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

func (m appModel) Update(tea.Msg) (tea.Model, tea.Cmd) {
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
