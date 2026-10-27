package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// ---- domain ----

type conversion struct {
	title   string
	convert func(float64) float64
}

type category struct {
	title       string
	conversions []conversion
}

var categories = []category{
	{title: "Temperature", conversions: []conversion{
		{"Fahrenheit to Celsius", func(f float64) float64 { return (f - 32) * 5 / 9 }},
		{"Celsius to Fahrenheit", func(c float64) float64 { return c*9/5 + 32 }},
	}},
}

// ---- navigation messages (children -> root) ----

type categoryChosenMsg struct{ category category }
type conversionChosenMsg struct{ conversion conversion }
type backMsg struct{}

func send(msg tea.Msg) tea.Cmd { return func() tea.Msg { return msg } }

// ---- root model: owns screen state and routing ----

type screen int

const (
	screenCategories screen = iota
	screenConversions
	screenInput
)

type model struct {
	screen      screen
	categories  menuModel
	conversions menuModel
	input       inputModel
}

func newModel() model {
	return model{categories: newCategoryMenu(categories)}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case categoryChosenMsg:
		m.conversions = newConversionMenu(msg.category)
		m.screen = screenConversions
		return m, nil
	case conversionChosenMsg:
		m.input = newInputModel(msg.conversion)
		m.screen = screenInput
		return m, m.input.Init() // the runtime only calls the *root* Init
	case backMsg:
		if m.screen > screenCategories {
			m.screen--
		}
		return m, nil
	}

	// Everything else goes to the active screen only.
	var cmd tea.Cmd
	switch m.screen {
	case screenCategories:
		m.categories, cmd = m.categories.Update(msg)
	case screenConversions:
		m.conversions, cmd = m.conversions.Update(msg)
	case screenInput:
		m.input, cmd = m.input.Update(msg)
	}
	return m, cmd
}

func (m model) View() tea.View {
	switch m.screen {
	case screenConversions:
		return tea.NewView(m.conversions.View())
	case screenInput:
		return tea.NewView(m.input.View())
	default:
		return tea.NewView(m.categories.View())
	}
}

// ---- menu: one generic list screen, reused for both menus ----

type menuModel struct {
	title    string
	labels   []string
	onSelect func(i int) tea.Msg
	canBack  bool
}

func newCategoryMenu(cs []category) menuModel {
	labels := make([]string, len(cs))
	for i, c := range cs {
		labels[i] = c.title
	}
	return menuModel{
		title:    "What would you like to convert?",
		labels:   labels,
		onSelect: func(i int) tea.Msg { return categoryChosenMsg{cs[i]} },
	}
}

func newConversionMenu(c category) menuModel {
	labels := make([]string, len(c.conversions))
	for i, cv := range c.conversions {
		labels[i] = cv.title
	}
	return menuModel{
		title:    c.title + " Conversions",
		labels:   labels,
		onSelect: func(i int) tea.Msg { return conversionChosenMsg{c.conversions[i]} },
		canBack:  true,
	}
}

func (m menuModel) Update(msg tea.Msg) (menuModel, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	if key.String() == "esc" && m.canBack {
		return m, send(backMsg{})
	}
	if n, err := strconv.Atoi(key.String()); err == nil && n >= 1 && n <= len(m.labels) {
		return m, send(m.onSelect(n - 1))
	}
	return m, nil
}

func (m menuModel) View() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n\n", m.title)
	for i, l := range m.labels {
		fmt.Fprintf(&b, "[%d] %s\n", i+1, l)
	}
	return b.String()
}

// ---- input screen ----

type inputModel struct {
	conversion conversion
	input      textinput.Model
	result     string
	err        error
}

func newInputModel(c conversion) inputModel {
	ti := textinput.New()
	ti.Placeholder = "Enter a number"
	ti.Focus()
	return inputModel{conversion: c, input: ti}
}

func (m inputModel) Init() tea.Cmd { return textinput.Blink }

func (m inputModel) Update(msg tea.Msg) (inputModel, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "esc":
			return m, send(backMsg{})
		case "enter":
			v, err := strconv.ParseFloat(strings.TrimSpace(m.input.Value()), 64)
			m.err = err
			if err == nil {
				m.result = strconv.FormatFloat(m.conversion.convert(v), 'f', 2, 64)
			}
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	s := m.conversion.title + "\n\n" + m.input.View() + "\n\n"
	switch {
	case m.err != nil:
		s += "Not a number: " + m.err.Error() + "\n"
	case m.result != "":
		s += "Result: " + m.result + "\n"
	}
	return s + "\n(esc to go back)"
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
