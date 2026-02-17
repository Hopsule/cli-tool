package ui

import (
	"fmt"
	"strings"

	"github.com/Hopsule/cli-tool/internal/api"
	"github.com/charmbracelet/lipgloss"
)

// ============================================================================
// LIST / DATATABLE VIEW
// ============================================================================
//
// Paginated table with search, column headers, alternating row colours,
// and page navigation (h/l = prev/next page).
// Modelled after the web-app DataTable pattern.
// ============================================================================

// ── Decisions list ──────────────────────────────────────────────────────────

func (m model) renderDecisionsList(filtered []api.Decision) string {
	var s strings.Builder
	total := len(filtered)
	start, end := m.getListPageItems(total)
	totalPages := (total + m.listPageSize - 1) / m.listPageSize
	pageSlice := filtered[start:end]

	// Column widths (dynamic based on terminal width)
	availW := m.width - 8
	if availW < 60 {
		availW = 60
	}
	statusW := 12
	dateW := 12
	idW := 14
	stmtW := availW - statusW - dateW - idW - 6
	if stmtW < 20 {
		stmtW = 20
	}

	// Header
	header := fmt.Sprintf("  %-*s %-*s %-*s %-*s",
		statusW, "STATUS", stmtW, "STATEMENT", dateW, "CREATED", idW, "ID")
	s.WriteString(tableHeaderStyle.Render(header) + "\n")

	// Rows
	for i, d := range pageSlice {
		globalIdx := start + i
		status := padRight(d.Status, statusW)
		stmt := padRight(truncateString(d.Statement, stmtW), stmtW)
		created := ""
		if d.CreatedAt != "" && len(d.CreatedAt) >= 10 {
			created = d.CreatedAt[:10]
		}
		created = padRight(created, dateW)
		id := padRight(truncateString(d.ID, idW), idW)

		row := fmt.Sprintf("  %s %s %s %s", status, stmt, created, id)

		statusStyle := lipgloss.NewStyle().Foreground(decisionStatusColor(d.Status))

		if globalIdx == m.selected {
			s.WriteString(tableSelectedRowStyle.Render("> "+row) + "\n")
		} else {
			styledStatus := statusStyle.Render(padRight(d.Status, statusW))
			rest := fmt.Sprintf(" %s %s %s", stmt, created, id)
			prefix := "  "
			if i%2 == 1 {
				s.WriteString(prefix + styledStatus + dimStyle.Render(rest) + "\n")
			} else {
				s.WriteString(prefix + styledStatus + normalStyle.Render(rest) + "\n")
			}
		}
	}

	// Pagination footer
	s.WriteString("\n")
	s.WriteString(m.renderPaginationFooter(start, end, total, totalPages))
	return s.String()
}

// ── Memories list ───────────────────────────────────────────────────────────

func (m model) renderMemoriesList(filtered []*api.Memory) string {
	var s strings.Builder
	total := len(filtered)
	start, end := m.getListPageItems(total)
	totalPages := (total + m.listPageSize - 1) / m.listPageSize
	pageSlice := filtered[start:end]

	availW := m.width - 8
	if availW < 60 {
		availW = 60
	}
	dateW := 12
	idW := 14
	contentW := availW - dateW - idW - 4
	if contentW < 20 {
		contentW = 20
	}

	header := fmt.Sprintf("  %-*s %-*s %-*s",
		contentW, "CONTENT", dateW, "CREATED", idW, "ID")
	s.WriteString(tableHeaderStyle.Render(header) + "\n")

	for i, mem := range pageSlice {
		globalIdx := start + i
		content := padRight(truncateString(strings.ReplaceAll(mem.Content, "\n", " "), contentW), contentW)
		created := ""
		if mem.CreatedAt != "" && len(mem.CreatedAt) >= 10 {
			created = mem.CreatedAt[:10]
		}
		created = padRight(created, dateW)
		id := padRight(truncateString(mem.ID, idW), idW)

		row := fmt.Sprintf("  %s %s %s", content, created, id)

		if globalIdx == m.selected {
			s.WriteString(tableSelectedRowStyle.Render("> "+row) + "\n")
		} else if i%2 == 1 {
			s.WriteString("  " + dimStyle.Render(row) + "\n")
		} else {
			s.WriteString("  " + normalStyle.Render(row) + "\n")
		}
	}

	s.WriteString("\n")
	s.WriteString(m.renderPaginationFooter(start, end, total, totalPages))
	return s.String()
}

// ── Capsules list ───────────────────────────────────────────────────────────

func (m model) renderCapsulesList(filtered []*api.Capsule) string {
	var s strings.Builder
	total := len(filtered)
	start, end := m.getListPageItems(total)
	totalPages := (total + m.listPageSize - 1) / m.listPageSize
	pageSlice := filtered[start:end]

	availW := m.width - 8
	if availW < 60 {
		availW = 60
	}
	statusW := 12
	dateW := 12
	decW := 5
	memW := 5
	nameW := availW - statusW - dateW - decW - memW - 8
	if nameW < 16 {
		nameW = 16
	}

	header := fmt.Sprintf("  %-*s %-*s %-*s %-*s %-*s",
		statusW, "STATUS", nameW, "NAME", decW, "DEC", memW, "MEM", dateW, "CREATED")
	s.WriteString(tableHeaderStyle.Render(header) + "\n")

	for i, c := range pageSlice {
		globalIdx := start + i
		status := padRight(c.Status, statusW)
		name := padRight(truncateString(c.Name, nameW), nameW)
		nDec := fmt.Sprintf("%-*d", decW, len(c.DecisionIds))
		nMem := fmt.Sprintf("%-*d", memW, len(c.MemoryIds))
		created := ""
		if c.CreatedAt != "" && len(c.CreatedAt) >= 10 {
			created = c.CreatedAt[:10]
		}
		created = padRight(created, dateW)

		row := fmt.Sprintf("  %s %s %s %s %s", status, name, nDec, nMem, created)

		if globalIdx == m.selected {
			s.WriteString(tableSelectedRowStyle.Render("> "+row) + "\n")
		} else {
			statusStyle := lipgloss.NewStyle().Foreground(capsuleStatusColor(c.Status))
			styledStatus := statusStyle.Render(padRight(c.Status, statusW))
			rest := fmt.Sprintf(" %s %s %s %s", name, nDec, nMem, created)
			prefix := "  "
			if i%2 == 1 {
				s.WriteString(prefix + styledStatus + dimStyle.Render(rest) + "\n")
			} else {
				s.WriteString(prefix + styledStatus + normalStyle.Render(rest) + "\n")
			}
		}
	}

	s.WriteString("\n")
	s.WriteString(m.renderPaginationFooter(start, end, total, totalPages))
	return s.String()
}

// ============================================================================
// PAGINATION FOOTER
// ============================================================================

func (m model) renderPaginationFooter(start, end, total, totalPages int) string {
	if total == 0 {
		return ""
	}
	page := m.listPage + 1

	var nav strings.Builder
	nav.WriteString("  ")

	if page > 1 {
		nav.WriteString(normalStyle.Render("< h prev"))
	} else {
		nav.WriteString(dimStyle.Render("< h prev"))
	}

	nav.WriteString(dimStyle.Render(fmt.Sprintf("   Page %d / %d   (%d-%d of %d)   ", page, totalPages, start+1, end, total)))

	if page < totalPages {
		nav.WriteString(normalStyle.Render("next l >"))
	} else {
		nav.WriteString(dimStyle.Render("next l >"))
	}

	nav.WriteString("\n")
	return nav.String()
}

// ============================================================================
// HELPERS
// ============================================================================

func padRight(s string, w int) string {
	if len(s) >= w {
		return s[:w]
	}
	return s + strings.Repeat(" ", w-len(s))
}

func decisionStatusColor(status string) lipgloss.AdaptiveColor {
	switch status {
	case "ACCEPTED":
		return greenColor
	case "PENDING":
		return yellowColor
	case "DEPRECATED":
		return redColor
	case "DRAFT":
		return blueColor
	default:
		return dimColor
	}
}

func capsuleStatusColor(status string) lipgloss.AdaptiveColor {
	switch status {
	case "FROZEN":
		return cyanColor
	case "HISTORICAL":
		return magentaColor
	default:
		return dimColor
	}
}
