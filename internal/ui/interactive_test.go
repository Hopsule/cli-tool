package ui

import (
	"strings"
	"testing"

	"github.com/Cagangedik/cli-tool/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

func TestNewInteractiveModel(t *testing.T) {
	cfg := &config.Config{
		APIURL:       "http://localhost:8080",
		Project:      "test",
		Organization: "test-org",
	}

	m := NewInteractiveModel(cfg)

	if len(m.commands) != 4 {
		t.Errorf("Expected 4 commands, got %d", len(m.commands))
	}

	if m.selected != 0 {
		t.Errorf("Expected initial selection to be 0, got %d", m.selected)
	}

	if m.showHelp {
		t.Errorf("Expected showHelp to be false initially")
	}
}

func TestModelUpdate_Navigation(t *testing.T) {
	cfg := &config.Config{}
	m := NewInteractiveModel(cfg)

	// Test down navigation
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updatedModel.(model)
	if m.selected != 1 {
		t.Errorf("Expected selection to be 1 after down, got %d", m.selected)
	}

	// Test up navigation
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updatedModel.(model)
	if m.selected != 0 {
		t.Errorf("Expected selection to be 0 after up, got %d", m.selected)
	}

	// Test down at boundary
	m = NewInteractiveModel(cfg)
	m.selected = 3
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updatedModel.(model)
	if m.selected != 3 {
		t.Errorf("Expected selection to stay at 3 at boundary, got %d", m.selected)
	}

	// Test up at boundary
	m = NewInteractiveModel(cfg)
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updatedModel.(model)
	if m.selected != 0 {
		t.Errorf("Expected selection to stay at 0 at boundary, got %d", m.selected)
	}
}

func TestModelUpdate_Help(t *testing.T) {
	cfg := &config.Config{}
	m := NewInteractiveModel(cfg)

	// Toggle help on
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m = updatedModel.(model)
	if !m.showHelp {
		t.Errorf("Expected showHelp to be true after ?")
	}

	// Toggle help off
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m = updatedModel.(model)
	if m.showHelp {
		t.Errorf("Expected showHelp to be false after second ?")
	}
}

func TestModelUpdate_Quit(t *testing.T) {
	cfg := &config.Config{}
	m := NewInteractiveModel(cfg)

	// Test quit with 'q'
	updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m = updatedModel.(model)
	if m.selected != -1 {
		t.Errorf("Expected selection to be -1 after quit, got %d", m.selected)
	}
	if cmd == nil {
		t.Errorf("Expected quit command to be returned")
	}
}

func TestModelView(t *testing.T) {
	cfg := &config.Config{}
	m := NewInteractiveModel(cfg)

	view := m.View()

	// Check for logo
	if !strings.Contains(view, "▟███▙") {
		t.Errorf("Expected view to contain logo")
	}

	// Check for title
	if !strings.Contains(view, "Hopsule") {
		t.Errorf("Expected view to contain title")
	}

	// Check for commands
	if !strings.Contains(view, "hopsule init") {
		t.Errorf("Expected view to contain hopsule init command")
	}
}

func TestLogoView(t *testing.T) {
	cfg := &config.Config{}
	m := NewInteractiveModel(cfg)

	logo := m.logoView()

	if !strings.Contains(logo, "▟███▙") {
		t.Errorf("Expected logo to contain geometric characters")
	}
}

func TestInfoView(t *testing.T) {
	cfg := &config.Config{}
	m := NewInteractiveModel(cfg)

	info := m.infoView()

	if !strings.Contains(info, "Hopsule") {
		t.Errorf("Expected info to contain Hopsule title")
	}

	if !strings.Contains(info, "Decision & Memory Layer") {
		t.Errorf("Expected info to contain subtitle")
	}

	if !strings.Contains(info, "v0.4.3") {
		t.Errorf("Expected info to contain version")
	}

	if !strings.Contains(info, "hopsule init") {
		t.Errorf("Expected info to contain init command")
	}
}

func TestSideBySide(t *testing.T) {
	cfg := &config.Config{}
	m := NewInteractiveModel(cfg)

	left := "line1\nline2"
	right := "right1\nright2"

	result := m.sideBySide(left, right)

	lines := strings.Split(result, "\n")
	if len(lines) < 2 {
		t.Errorf("Expected at least 2 lines in side by side output, got %d", len(lines))
	}

	// Check that left and right are combined
	if !strings.Contains(result, "line1") || !strings.Contains(result, "right1") {
		t.Errorf("Expected side by side to contain both left and right content")
	}
}

func TestHelpView(t *testing.T) {
	cfg := &config.Config{}
	m := NewInteractiveModel(cfg)

	help := m.helpView()

	if !strings.Contains(help, "Keyboard Shortcuts") {
		t.Errorf("Expected help to contain Keyboard Shortcuts")
	}

	if !strings.Contains(help, "Navigation") {
		t.Errorf("Expected help to contain Navigation section")
	}

	if !strings.Contains(help, "Actions") {
		t.Errorf("Expected help to contain Actions section")
	}
}
