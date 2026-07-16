// Package tui implements the interactive `fullstack-template init` wizard.
// It only collects a scaffold.Config — generation itself happens in the
// caller, so all generator logic stays testable in internal/scaffold.
package tui

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jaireddjawed/fullstack-template-go-pocketbase/internal/scaffold"
)

type step int

const (
	stepName step = iota
	stepModule
	stepFrontend
	stepAuth
	stepDatabase
	stepExtras
	stepConfirm
	stepDone
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true)
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	disabledStyle = lipgloss.NewStyle().Faint(true)
	detailStyle   = lipgloss.NewStyle().Faint(true)
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	helpStyle     = lipgloss.NewStyle().Faint(true).MarginTop(1)
	summaryKey    = lipgloss.NewStyle().Faint(true).Width(10)
)

type model struct {
	step step

	nameInput   textinput.Model
	moduleInput textinput.Model

	cursor  int
	checked map[string]bool // extras toggles

	cfg     scaffold.Config
	errMsg  string
	aborted bool
}

// Run shows the wizard and returns the collected config.
// ok is false when the user aborted.
func Run(defaults scaffold.Config) (cfg scaffold.Config, ok bool, err error) {
	name := textinput.New()
	name.Placeholder = "my-app"
	name.SetValue(defaults.Name)
	name.Focus()

	module := textinput.New()
	module.Placeholder = "github.com/you/my-app"
	module.SetValue(defaults.Module)

	m := model{
		nameInput:   name,
		moduleInput: module,
		checked:     map[string]bool{},
		cfg:         defaults,
	}

	final, err := tea.NewProgram(m).Run()
	if err != nil {
		return scaffold.Config{}, false, err
	}

	result := final.(model)
	return result.cfg, !result.aborted && result.step == stepDone, nil
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// choices returns the option list for the current selection step.
func (m model) choices() []scaffold.Choice {
	switch m.step {
	case stepFrontend:
		return scaffold.FrontendChoices()
	case stepAuth:
		return scaffold.AuthChoices(m.cfg.Frontend)
	case stepDatabase:
		return scaffold.DatabaseChoices()
	case stepExtras:
		return scaffold.ExtraChoices(m.cfg.Frontend, m.cfg.Auth)
	default:
		return nil
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, isKey := msg.(tea.KeyMsg)
	if isKey {
		switch key.String() {
		case "ctrl+c", "q":
			if m.step != stepName && m.step != stepModule || key.String() == "ctrl+c" {
				m.aborted = true
				return m, tea.Quit
			}
		case "esc":
			if m.step > stepName {
				m.step--
				m.cursor = 0
				m.errMsg = ""
				m.syncFocus()
			}
			return m, nil
		}
	}

	switch m.step {
	case stepName, stepModule:
		return m.updateTextStep(msg)
	case stepFrontend, stepAuth, stepDatabase, stepExtras:
		if isKey {
			return m.updateChoiceStep(key), nil
		}
	case stepConfirm:
		if isKey && key.String() == "enter" {
			m.step = stepDone
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) updateTextStep(msg tea.Msg) (tea.Model, tea.Cmd) {
	input := &m.nameInput
	if m.step == stepModule {
		input = &m.moduleInput
	}

	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		value := strings.TrimSpace(input.Value())
		if value == "" {
			value = input.Placeholder
		}

		if m.step == stepName {
			m.cfg.Name = value
			if m.moduleInput.Value() == "" {
				m.moduleInput.SetValue("github.com/you/" + value)
			}
		} else {
			m.cfg.Module = value
		}

		m.step++
		m.errMsg = ""
		m.syncFocus()
		return m, nil
	}

	var cmd tea.Cmd
	*input, cmd = input.Update(msg)
	return m, cmd
}

func (m *model) syncFocus() {
	m.nameInput.Blur()
	m.moduleInput.Blur()
	switch m.step {
	case stepName:
		m.nameInput.Focus()
	case stepModule:
		m.moduleInput.Focus()
	}
}

func (m model) updateChoiceStep(key tea.KeyMsg) model {
	choices := m.choices()

	switch key.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(choices)-1 {
			m.cursor++
		}
	case " ":
		if m.step == stepExtras {
			choice := choices[m.cursor]
			if choice.Disabled {
				m.errMsg = choice.Label + " is not available for this stack"
			} else {
				m.checked[choice.Value] = !m.checked[choice.Value]
				m.errMsg = ""
			}
		}
	case "enter":
		choice := choices[m.cursor]

		if m.step != stepExtras && choice.Disabled {
			m.errMsg = choice.Label + " is not available yet"
			return m
		}

		m.errMsg = ""
		switch m.step {
		case stepFrontend:
			m.cfg.Frontend = scaffold.Frontend(choice.Value)
			// A frontend change can invalidate dependent selections.
			m.cfg.Auth = ""
			m.checked = map[string]bool{}
		case stepAuth:
			m.cfg.Auth = scaffold.Auth(choice.Value)
		case stepDatabase:
			m.cfg.Database = scaffold.Database(choice.Value)
		case stepExtras:
			m.cfg.Extras = nil
			for _, c := range scaffold.ExtraChoices(m.cfg.Frontend, m.cfg.Auth) {
				if m.checked[c.Value] {
					m.cfg.Extras = append(m.cfg.Extras, scaffold.Extra(c.Value))
				}
			}
		}

		m.step++
		m.cursor = 0
	}

	return m
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("fullstack-template init") + "\n\n")

	switch m.step {
	case stepName:
		b.WriteString("Project name:\n\n" + m.nameInput.View())
	case stepModule:
		b.WriteString("Go module path:\n\n" + m.moduleInput.View())
	case stepFrontend:
		m.renderChoices(&b, "Choose a frontend:")
	case stepAuth:
		m.renderChoices(&b, "Choose auth:")
	case stepDatabase:
		m.renderChoices(&b, "Choose a database:")
	case stepExtras:
		m.renderChoices(&b, "Choose extras (space to toggle):")
	case stepConfirm, stepDone:
		b.WriteString("Ready to generate:\n\n")
		b.WriteString(m.summary())
	}

	if m.errMsg != "" {
		b.WriteString("\n" + errorStyle.Render(m.errMsg))
	}

	b.WriteString(helpStyle.Render("\n" + m.helpLine()))
	return b.String() + "\n"
}

func (m model) renderChoices(b *strings.Builder, title string) {
	b.WriteString(title + "\n\n")

	for i, choice := range m.choices() {
		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("> ")
		}

		label := choice.Label
		if m.step == stepExtras {
			mark := "[ ]"
			if m.checked[choice.Value] {
				mark = "[x]"
			}
			label = mark + " " + label
		}

		line := label + "  " + detailStyle.Render(choice.Detail)
		switch {
		case choice.Disabled:
			line = disabledStyle.Render(label + "  " + choice.Detail)
		case i == m.cursor:
			line = selectedStyle.Render(label) + "  " + detailStyle.Render(choice.Detail)
		}

		fmt.Fprintf(b, "%s%s\n", cursor, line)
	}
}

func (m model) summary() string {
	extras := "none"
	if len(m.cfg.Extras) > 0 {
		parts := make([]string, len(m.cfg.Extras))
		for i, e := range m.cfg.Extras {
			parts[i] = string(e)
		}
		extras = strings.Join(parts, ", ")
	}

	branch, _ := m.cfg.Branch()

	rows := [][2]string{
		{"name", m.cfg.Name},
		{"module", m.cfg.Module},
		{"frontend", string(m.cfg.Frontend)},
		{"auth", string(m.cfg.Auth)},
		{"database", string(m.cfg.Database)},
		{"extras", extras},
		{"branch", branch},
	}

	var b strings.Builder
	for _, row := range rows {
		b.WriteString(summaryKey.Render(row[0]) + " " + row[1] + "\n")
	}
	return b.String()
}

func (m model) helpLine() string {
	base := []string{"esc back", "ctrl+c quit"}
	switch m.step {
	case stepName, stepModule:
		base = slices.Insert(base, 0, "enter next")
	case stepExtras:
		base = slices.Insert(base, 0, "↑/↓ move", "space toggle", "enter next")
	case stepConfirm:
		base = slices.Insert(base, 0, "enter generate")
	default:
		base = slices.Insert(base, 0, "↑/↓ move", "enter select")
	}
	return strings.Join(base, " · ")
}
