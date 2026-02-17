package ui

import (
	"fmt"
	"strings"
)

// ============================================================================
// DASHBOARD VIEW
// ============================================================================

func (m model) renderDashboardContent() string {
	var s strings.Builder
	dCount := len(m.decisions)
	mCount := len(m.memories)
	tCount := len(m.tasks)
	cCount := len(m.capsules)

	statsLine := fmt.Sprintf("  Decisions: %d  |  Memories: %d  |  Tasks: %d  |  Capsules: %d", dCount, mCount, tCount, cCount)
	s.WriteString(dimStyle.Render(statsLine) + "\n\n")

	if len(m.tasks) > 0 {
		s.WriteString("  " + titleStyle.Render("Tasks") + "\n\n")
		kanban := TaskKanban(m.tasks, 4)
		for _, line := range strings.Split(kanban, "\n") {
			s.WriteString("  " + line + "\n")
		}
	} else {
		s.WriteString("  " + dimStyle.Render("No tasks yet. Go to Tasks to create one.") + "\n\n")
	}

	if len(m.decisions) > 0 {
		s.WriteString("\n  " + titleStyle.Render("Recent Decisions") + "\n")
		for i, d := range m.decisions {
			if i >= 3 {
				break
			}
			icon := "○"
			if d.Status == "ACCEPTED" {
				icon = "●"
			} else if d.Status == "PENDING" {
				icon = "◐"
			}
			s.WriteString(fmt.Sprintf("  %s %s\n", icon, truncateString(d.Statement, 50)))
		}
	}
	return s.String()
}

// ============================================================================
// TASKS VIEW
// ============================================================================

func (m model) renderTasksContent() string {
	var s strings.Builder

	if m.tasksViewMode == "kanban" {
		s.WriteString("  " + dimStyle.Render("List") + " " + selectedStyle.Render("[Kanban]") + "\n\n")
	} else {
		s.WriteString("  " + selectedStyle.Render("[List]") + " " + dimStyle.Render("Kanban") + "\n\n")
	}

	if m.tasksViewMode == "kanban" {
		if len(m.tasks) == 0 {
			s.WriteString("  " + dimStyle.Render("No tasks yet. Press [n] to create one.") + "\n")
		} else {
			kanban := TaskKanban(m.tasks, 6)
			for _, line := range strings.Split(kanban, "\n") {
				s.WriteString("  " + line + "\n")
			}
		}
		return s.String()
	}

	// List view
	if m.searchMode {
		s.WriteString("  Search: " + m.searchQuery + "\u2588\n")
	} else if m.searchQuery != "" {
		s.WriteString("  " + dimStyle.Render("Filter: "+m.searchQuery) + "\n")
	}

	filtered := m.getFilteredTasks()
	s.WriteString("  " + dimStyle.Render(fmt.Sprintf("%d items", len(filtered))) + "\n\n")

	if len(filtered) == 0 {
		if m.searchQuery != "" {
			s.WriteString("  " + dimStyle.Render("No matching tasks.") + "\n")
		} else {
			s.WriteString("  " + dimStyle.Render("No tasks yet. Press [n] to create one.") + "\n")
		}
		return s.String()
	}

	visibleItems := m.getVisibleItems()
	startIdx := m.scrollOffset
	if startIdx > len(filtered)-visibleItems {
		startIdx = max(0, len(filtered)-visibleItems)
	}
	endIdx := min(startIdx+visibleItems, len(filtered))

	for i := startIdx; i < endIdx; i++ {
		t := filtered[i]
		icon := "○"
		switch t.Status {
		case "DONE":
			icon = "v"
		case "IN_PROGRESS":
			icon = ">"
		case "REVIEW":
			icon = "?"
		}
		title := truncateString(t.Title, 45)
		status := dimStyle.Render(t.Status)
		if i == m.selected {
			s.WriteString("  " + selectedStyle.Render("> "+icon+" "+title) + " " + status + "\n")
		} else {
			s.WriteString("    " + icon + " " + normalStyle.Render(title) + " " + status + "\n")
		}
	}

	if len(filtered) > visibleItems {
		s.WriteString("\n  " + dimStyle.Render(fmt.Sprintf("Showing %d-%d of %d", startIdx+1, endIdx, len(filtered))) + "\n")
	}
	return s.String()
}
