package tui

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jaireddjawed/fullstack-template-go-pocketbase/internal/scaffold"
)

func newTestModel() model {
	name := textinput.New()
	name.Focus()
	return model{
		nameInput:   name,
		moduleInput: textinput.New(),
		checked:     map[string]bool{},
		cfg:         scaffold.Config{Database: scaffold.DatabasePocketBase},
	}
}

func press(m model, keys ...string) model {
	for _, k := range keys {
		var msg tea.Msg
		switch k {
		case "enter":
			msg = tea.KeyMsg{Type: tea.KeyEnter}
		case "down":
			msg = tea.KeyMsg{Type: tea.KeyDown}
		case "space":
			msg = tea.KeyMsg{Type: tea.KeySpace}
		case "esc":
			msg = tea.KeyMsg{Type: tea.KeyEsc}
		default:
			msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)}
		}
		next, _ := m.Update(msg)
		m = next.(model)
	}
	return m
}

// Walks the whole wizard: name → module → Next.js → Clerk → PocketBase db →
// shadcn extra → confirm.
func TestWizardFullWalk(t *testing.T) {
	m := newTestModel()

	m = press(m, "m", "y", "-", "a", "p", "p", "enter") // name
	if m.step != stepModule {
		t.Fatalf("after name, step = %d, want %d", m.step, stepModule)
	}
	if m.cfg.Name != "my-app" {
		t.Fatalf("name = %q", m.cfg.Name)
	}

	m = press(m, "enter") // module (defaulted from name)
	if m.cfg.Module != "github.com/you/my-app" {
		t.Fatalf("module = %q", m.cfg.Module)
	}

	m = press(m, "down", "enter") // frontend: next
	if m.cfg.Frontend != scaffold.FrontendNext {
		t.Fatalf("frontend = %q", m.cfg.Frontend)
	}

	m = press(m, "down", "enter") // auth: clerk
	if m.cfg.Auth != scaffold.AuthClerk {
		t.Fatalf("auth = %q", m.cfg.Auth)
	}

	m = press(m, "enter") // database: pocketbase

	m = press(m, "down", "space", "enter") // extras: toggle shadcn, continue
	if len(m.cfg.Extras) != 1 || m.cfg.Extras[0] != scaffold.ExtraShadcn {
		t.Fatalf("extras = %v", m.cfg.Extras)
	}

	if m.step != stepConfirm {
		t.Fatalf("step = %d, want confirm", m.step)
	}

	m = press(m, "enter")
	if m.step != stepDone {
		t.Fatalf("step = %d, want done", m.step)
	}

	if err := m.cfg.Validate(); err != nil {
		t.Fatalf("collected config invalid: %v", err)
	}
	if branch, _ := m.cfg.Branch(); branch != "nextjs-clerk" {
		t.Fatalf("branch = %q, want nextjs-clerk", branch)
	}
}

func TestInertiaShadcnKeepsReactInertiaBranch(t *testing.T) {
	m := newTestModel()
	m.step = stepFrontend

	m = press(m, "down", "down", "enter") // frontend: React + Inertia
	m = press(m, "enter")                 // auth: PocketBase
	m = press(m, "enter")                 // database: PocketBase
	m = press(m, "down", "space", "enter")

	if m.cfg.Frontend != scaffold.FrontendInertia {
		t.Fatalf("frontend = %q, want inertia", m.cfg.Frontend)
	}
	if len(m.cfg.Extras) != 1 || m.cfg.Extras[0] != scaffold.ExtraShadcn {
		t.Fatalf("extras = %v, want shadcn", m.cfg.Extras)
	}
	if branch, err := m.cfg.Branch(); err != nil || branch != "react-inertia" {
		t.Fatalf("branch = %q, err = %v; want react-inertia", branch, err)
	}
}

func TestDisabledOptionCannotBeSelected(t *testing.T) {
	m := newTestModel()
	m.step = stepAuth
	m.cfg.Frontend = scaffold.FrontendInertia

	m = press(m, "down", "enter") // clerk is disabled for inertia
	if m.step != stepAuth {
		t.Fatal("selecting a disabled option should not advance the wizard")
	}
	if m.errMsg == "" {
		t.Fatal("expected an error message for a disabled option")
	}
}

func TestEscGoesBack(t *testing.T) {
	m := newTestModel()
	m.step = stepDatabase

	m = press(m, "esc")
	if m.step != stepAuth {
		t.Fatalf("esc: step = %d, want %d", m.step, stepAuth)
	}
}
