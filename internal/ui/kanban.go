package ui

import (
	"fmt"
	"strings"

	"github.com/Cagangedik/cli-tool/internal/api"
	"github.com/charmbracelet/lipgloss"
)

// Kanban-style display for terminal
// Uses box-drawing characters for a clean look

var (
	kanbanHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "255"}).
				BorderStyle(lipgloss.NormalBorder()).
				BorderBottom(true).
				BorderForeground(lipgloss.AdaptiveColor{Light: "240", Dark: "248"}).
				PaddingLeft(1).
				PaddingRight(1)

	kanbanItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "252"}).
			PaddingLeft(1)

	kanbanSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "255"}).
				Background(lipgloss.AdaptiveColor{Light: "252", Dark: "238"}).
				PaddingLeft(1).
				PaddingRight(1)

	kanbanDimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "244", Dark: "244"})

	kanbanBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "240", Dark: "248"}).
			Padding(0, 1)
)

// OrganizationKanban displays organizations in a kanban-style layout
func OrganizationKanban(orgs []*api.Organization, selected int) string {
	if len(orgs) == 0 {
		return kanbanBoxStyle.Render("No organizations found")
	}

	var sb strings.Builder
	
	sb.WriteString(kanbanHeaderStyle.Render("Organizations"))
	sb.WriteString("\n\n")

	for i, org := range orgs {
		marker := "  "
		style := kanbanItemStyle
		if i == selected {
			marker = "> "
			style = kanbanSelectedStyle
		}

		name := truncateStr(org.Name, 30)
		slug := "@" + org.Slug
		
		line := fmt.Sprintf("%s%s %s", marker, name, kanbanDimStyle.Render(slug))
		sb.WriteString(style.Render(line))
		sb.WriteString("\n")
	}

	return sb.String()
}

// ProjectKanban displays projects grouped by status or organization
func ProjectKanban(projects []*api.Project, orgNames map[string]string, selected int) string {
	if len(projects) == 0 {
		return kanbanBoxStyle.Render("No projects found")
	}

	var sb strings.Builder
	
	sb.WriteString(kanbanHeaderStyle.Render("Projects"))
	sb.WriteString("\n\n")

	for i, proj := range projects {
		marker := "  "
		style := kanbanItemStyle
		if i == selected {
			marker = "> "
			style = kanbanSelectedStyle
		}

		name := truncateStr(proj.Name, 25)
		orgName := orgNames[proj.OrganizationID]
		if orgName == "" {
			orgName = "Unknown"
		}
		orgName = truncateStr(orgName, 15)

		line := fmt.Sprintf("%s%-25s %s", marker, name, kanbanDimStyle.Render(orgName))
		sb.WriteString(style.Render(line))
		sb.WriteString("\n")
	}

	return sb.String()
}

// DecisionKanban displays decisions grouped by status
func DecisionKanban(decisions []api.Decision) string {
	if len(decisions) == 0 {
		return kanbanBoxStyle.Render("No decisions found")
	}

	// Group by status
	groups := map[string][]api.Decision{
		"DRAFT":      {},
		"PENDING":    {},
		"ACCEPTED":   {},
		"DEPRECATED": {},
	}

	for _, d := range decisions {
		status := d.Status
		if _, ok := groups[status]; ok {
			groups[status] = append(groups[status], d)
		}
	}

	// Build columns
	columns := []string{}
	columnOrder := []string{"DRAFT", "PENDING", "ACCEPTED", "DEPRECATED"}
	
	for _, status := range columnOrder {
		items := groups[status]
		col := buildKanbanColumn(status, items)
		columns = append(columns, col)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, columns...)
}

func buildKanbanColumn(title string, decisions []api.Decision) string {
	var sb strings.Builder
	
	// Header with count
	header := fmt.Sprintf("%s (%d)", title, len(decisions))
	sb.WriteString(kanbanHeaderStyle.Render(header))
	sb.WriteString("\n")

	if len(decisions) == 0 {
		sb.WriteString(kanbanDimStyle.Render("  (empty)"))
		sb.WriteString("\n")
	} else {
		for _, d := range decisions {
			stmt := truncateStr(d.Statement, 20)
			sb.WriteString(kanbanItemStyle.Render("• " + stmt))
			sb.WriteString("\n")
		}
	}

	// Add some padding
	for sb.Len() < 200 {
		sb.WriteString("\n")
	}

	return kanbanBoxStyle.Width(25).Render(sb.String())
}

// StatusBoard shows a quick status overview
func StatusBoard(stats *api.ProjectStatus) string {
	var sb strings.Builder
	
	sb.WriteString(kanbanHeaderStyle.Render("Project Status"))
	sb.WriteString("\n\n")

	// Stats in a grid
	statItems := []struct {
		label string
		value int
		emoji string
	}{
		{"Total", stats.TotalDecisions, "📊"},
		{"Accepted", stats.Accepted, "✓"},
		{"Pending", stats.Pending, "⏳"},
		{"Draft", stats.Draft, "📝"},
		{"Deprecated", stats.Deprecated, "🚫"},
	}

	for _, item := range statItems {
		line := fmt.Sprintf("  %s %-12s %d", item.emoji, item.label+":", item.value)
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	return kanbanBoxStyle.Render(sb.String())
}

// UserCard displays user info in a card format
func UserCard(user *api.User) string {
	var sb strings.Builder
	
	sb.WriteString(kanbanHeaderStyle.Render("Current User"))
	sb.WriteString("\n\n")
	
	sb.WriteString(fmt.Sprintf("  Name:  %s\n", user.Name))
	sb.WriteString(fmt.Sprintf("  Email: %s\n", user.Email))
	sb.WriteString(fmt.Sprintf("  ID:    %s\n", truncateStr(user.ID, 12)))

	return kanbanBoxStyle.Render(sb.String())
}

func truncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
