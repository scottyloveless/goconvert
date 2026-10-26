package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type item struct {
	num   int
	title string
}

type mainMenuModel struct {
	Options  []item
	TempMenu tempMenuModel
}

func initialModel() mainMenuModel {
	return mainMenuModel{
		Options: []item{
			{num: 1, title: "Temperature"},
		},
	}
}

func (m mainMenuModel) Init() tea.Cmd {
	return nil
}

func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "1":
			return initialTempModel(), nil // temp model goes here
		}

		return m, nil
	}
	return m, nil
}

func (m mainMenuModel) View() tea.View {
	s := "What would you like to convert?\n\n"

	for _, choice := range m.Options {
		s += fmt.Sprintf("[%d] %s", choice.num, choice.title)
	}
	return tea.NewView(s)
}

type tempMenuModel struct {
	Options []item
}

func initialTempModel() tea.Model {
	return tempMenuModel{
		Options: []item{
			{num: 1, title: "Farenheit to Celcius"},
		},
	}
}

func (m tempMenuModel) Init() tea.Cmd {
	return nil
}

func (m tempMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			return initialModel(), nil
		case "1":
			return initialInputModel(), nil
		}
		return m, nil
	}
	return m, nil
}

func (m tempMenuModel) View() tea.View {
	s := "Temperature Conversions\n\n"

	for _, choice := range m.Options {
		s += fmt.Sprintf("[%d] %s", choice.num, choice.title)
	}

	return tea.NewView(s)
}

type inputModel struct {
	textarea textarea.Model
	err      error
}

func initialInputModel() tea.Model {
	ta := textarea.New()
	ta.Placeholder = "Enter number as float"
	ta.SetStyles(textarea.DefaultStyles(true))
	ta.Focus()

	return inputModel{
		textarea: ta,
		err:      nil,
	}
}

func (m inputModel) headerView() string {
	return "Please input value as float: \n\n"
}

func (i inputModel) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, tea.RequestBackgroundColor)
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		m.textarea.SetStyles(textarea.DefaultStyles(msg.IsDark()))
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			return initialTempModel(), nil
		default:
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m inputModel) View() tea.View {
	var c *tea.Cursor
	if !m.textarea.VirtualCursor() {
		c = m.textarea.Cursor()

		if c != nil {
			// Set the y offset of the cursor based on the position of the textarea
			// in the application.
			offset := lipgloss.Height(m.headerView())
			c.Y += offset
		}
	}

	f := strings.Join([]string{
		m.headerView(),
		m.textarea.View(),
	}, "\n")

	v := tea.NewView(f)
	v.Cursor = c

	return v
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
