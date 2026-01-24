package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Cagangedik/cli-tool/internal/config"
)

type command struct {
	name        string
	description string
	execute     string
}

type model struct {
	commands []command
	selected int
	width    int
	height   int
	cfg      *config.Config
	showHelp bool
}

var (
	// Adaptive colors for dark/light terminal support
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "255"}).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "255"}).
			Bold(true).
			PaddingLeft(2)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "240", Dark: "250"}).
			PaddingLeft(4)

	logoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "238", Dark: "252"})

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "244", Dark: "244"})

	accentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "238", Dark: "252"})

	versionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "248", Dark: "240"})

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "255"}).
			BorderForeground(lipgloss.AdaptiveColor{Light: "240", Dark: "248"}).
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2)
)

func NewInteractiveModel(cfg *config.Config) model {
	commands := []command{
		{"hopsule config", "Configure CLI settings", "config"},
		{"hopsule list", "List all decisions", "list"},
		{"hopsule create", "Create a new decision", "create"},
		{"hopsule get <id>", "Get decision details", "get"},
		{"hopsule accept <id>", "Accept a decision", "accept"},
		{"hopsule deprecate <id>", "Deprecate a decision", "deprecate"},
		{"hopsule status", "Show project status", "status"},
		{"hopsule sync", "Sync with decision-api", "sync"},
	}

	return model{
		commands: commands,
		selected: 0,
		cfg:      cfg,
		showHelp: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.showHelp {
				m.showHelp = false
				return m, nil
			}
			return m, tea.Quit

		case "?":
			m.showHelp = !m.showHelp
			return m, nil

		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}

		case "down", "j":
			if m.selected < len(m.commands)-1 {
				m.selected++
			}

		case "enter":
			// Return selected command to be executed
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	if m.showHelp {
		return m.helpView()
	}

	var s strings.Builder

	// Top border - using a dedicated border style
	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "240", Dark: "248"})
	
	s.WriteString(borderStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") + "\n\n")

	// Logo and info side by side
	logo := m.logoView()
	info := m.infoView()

	// Split into left (logo) and right (info)
	s.WriteString(m.sideBySide(logo, info))

	s.WriteString("\n")
	s.WriteString(borderStyle.Render("        ─────────────────────────────────────────────────────────────────────────────") + "\n\n")

	// Status
	status := "        ✓ " + titleStyle.Render("Connected") + "\n\n"
	if m.cfg.APIURL == "" {
		status = "        ⚠ " + versionStyle.Render("Not configured yet") + "\n\n"
	}
	s.WriteString(status)

	// Commands section
	s.WriteString("        " + titleStyle.Render("Get started") + "\n\n")

	for i, cmd := range m.commands {
		if i == m.selected {
			s.WriteString(selectedStyle.Render(fmt.Sprintf("❯ %-25s %s", cmd.name, infoStyle.Render("("+cmd.description+")"))) + "\n")
		} else {
			s.WriteString(normalStyle.Render(fmt.Sprintf("  %-25s %s", cmd.name, infoStyle.Render("("+cmd.description+")"))) + "\n")
		}
	}

	s.WriteString("\n")

	// Footer
	apiURL := "http://localhost:8080"
	if m.cfg.APIURL != "" {
		apiURL = m.cfg.APIURL
	}

	tokenStatus := versionStyle.Render("not set")
	if m.cfg.Token != "" {
		tokenStatus = accentStyle.Render("configured ✓")
	}

	s.WriteString(fmt.Sprintf("        %s: %s\n", infoStyle.Render("API"), accentStyle.Render(apiURL)))
	s.WriteString(fmt.Sprintf("        %s: %s\n\n", infoStyle.Render("Token"), tokenStatus))

	// Keybinds
	s.WriteString(infoStyle.Render("        ↑/↓: navigate  •  Enter: execute  •  q: quit  •  ?: help") + "\n")

	s.WriteString(borderStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") + "\n")

	return s.String()
}

func (m model) logoView() string {
	logo := logoStyle.Render(`
          ███  ███
          ████████
          ████████
          ███  ███
`)
	return logo
}

func (m model) infoView() string {
	org := "hopsule-inc"
	if m.cfg.Organization != "" {
		org = m.cfg.Organization
	}

	project := "app"
	if m.cfg.Project != "" {
		project = truncate(m.cfg.Project, 30)
	}

	var s strings.Builder

	s.WriteString("        " + titleStyle.Render("Hopsule") + "\n")
	s.WriteString("        " + versionStyle.Render("The future of dev governance") + "\n\n")

	s.WriteString(fmt.Sprintf("        %s: %s  •  %s: %s\n",
		infoStyle.Render("org"), accentStyle.Render(org),
		infoStyle.Render("project"), accentStyle.Render(project)))

	captureStatus := accentStyle.Render("ON")
	syncStatus := accentStyle.Render("ON")
	privacy := versionStyle.Render("redacted")

	s.WriteString(fmt.Sprintf("        %s: %s  •  %s: %s  •  %s: %s\n",
		infoStyle.Render("capture"), captureStatus,
		infoStyle.Render("sync"), syncStatus,
		infoStyle.Render("privacy"), privacy))

	s.WriteString(fmt.Sprintf("        %s: %s  •  %s: %s\n",
		infoStyle.Render("last sync"), versionStyle.Render("12s"),
		infoStyle.Render("latency"), accentStyle.Render("84ms")))

	s.WriteString("\n        " + versionStyle.Render("v0.1.1") + "\n")

	return s.String()
}

func (m model) sideBySide(left, right string) string {
	leftLines := strings.Split(left, "\n")
	rightLines := strings.Split(right, "\n")

	maxLines := len(leftLines)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}

	var result strings.Builder
	for i := 0; i < maxLines; i++ {
		var leftLine, rightLine string

		if i < len(leftLines) {
			leftLine = leftLines[i]
		}
		if i < len(rightLines) {
			rightLine = rightLines[i]
		}

		// Use lipgloss Width to get actual display width (handles ANSI codes)
		leftWidth := lipgloss.Width(leftLine)
		paddingNeeded := 40 - leftWidth
		
		// Ensure non-negative padding
		if paddingNeeded < 0 {
			paddingNeeded = 0
		}
		
		leftPadded := leftLine + strings.Repeat(" ", paddingNeeded)
		result.WriteString(leftPadded + rightLine + "\n")
	}

	return result.String()
}

func (m model) helpView() string {
	help := `
Hopsule CLI - Keyboard Shortcuts

Navigation:
  ↑/k         Move selection up
  ↓/j         Move selection down

Actions:
  Enter       Execute selected command
  q           Quit
  ?           Toggle this help

Commands:
  hopsule config      Configure CLI settings
  hopsule list        List all decisions
  hopsule create      Create a new decision
  hopsule get <id>    Get decision details
  hopsule accept <id> Accept a decision
  hopsule deprecate   Deprecate a decision
  hopsule status      Show project status
  hopsule sync        Sync with decision-api

Configuration:
  Config file: ~/.decision-cli/config.yaml
  Environment: DECISION_API_URL, DECISION_PROJECT, DECISION_TOKEN

Press ? to close this help
`

	return helpStyle.Render(help)
}

// stripAnsi is no longer needed - using lipgloss.Width instead

func RunInteractive(cfg *config.Config) (string, error) {
	p := tea.NewProgram(NewInteractiveModel(cfg))
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	m := finalModel.(model)
	if m.selected >= 0 && m.selected < len(m.commands) {
		return m.commands[m.selected].execute, nil
	}

	return "", nil
}
