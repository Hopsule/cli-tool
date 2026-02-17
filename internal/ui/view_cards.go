package ui

import (
	"fmt"
	"strings"

	"github.com/Hopsule/cli-tool/internal/api"
	"github.com/charmbracelet/lipgloss"
)

// ============================================================================
// DECISIONS — card grid + detail
// ============================================================================

func (m model) renderDecisionsContent() string {
	if m.detailView {
		d := m.getSelectedDecision()
		if d != nil {
			return m.renderDecisionDetail(*d)
		}
		m.detailView = false
	}

	var s strings.Builder

	// View mode toggle + status filter
	s.WriteString(m.renderViewModeHeader("Decisions"))

	// Status filter tabs
	filters := []struct{ label, value string }{
		{"ALL", ""}, {"DRAFT", "DRAFT"}, {"PENDING", "PENDING"}, {"ACCEPTED", "ACCEPTED"}, {"DEPRECATED", "DEPRECATED"},
	}
	s.WriteString("  ")
	for _, f := range filters {
		if f.value == m.statusFilter {
			s.WriteString(selectedStyle.Render("["+f.label+"]") + " ")
		} else {
			s.WriteString(dimStyle.Render(" "+f.label+" ") + " ")
		}
	}
	s.WriteString("\n")

	// Search
	if m.searchMode {
		s.WriteString("  Search: " + m.searchQuery + "\u2588\n")
	} else if m.searchQuery != "" {
		s.WriteString("  " + dimStyle.Render("Filter: "+m.searchQuery) + "\n")
	}

	filtered := m.getFilteredDecisions()
	s.WriteString("  " + dimStyle.Render(fmt.Sprintf("%d items", len(filtered))) + "\n\n")

	if len(filtered) == 0 {
		if m.searchQuery != "" || m.statusFilter != "" {
			s.WriteString("  " + dimStyle.Render("No matching decisions.") + "\n")
		} else {
			s.WriteString("  " + dimStyle.Render("No decisions yet. Press [n] to create one.") + "\n")
		}
		return s.String()
	}

	if m.currentViewMode == viewModeList {
		s.WriteString(m.renderDecisionsList(filtered))
	} else {
		s.WriteString(m.renderDecisionsCards(filtered))
	}

	return s.String()
}

func (m model) renderDecisionsCards(filtered []api.Decision) string {
	cols := m.getGridCols()
	cardW := m.getCardWidth()
	visibleRows := m.getVisibleRows()
	totalRows := (len(filtered) + cols - 1) / cols

	scrollOff := m.scrollOffset
	if scrollOff > totalRows-visibleRows {
		scrollOff = max(0, totalRows-visibleRows)
	}
	startRow := scrollOff
	endRow := min(startRow+visibleRows, totalRows)

	var gridLines []string
	for row := startRow; row < endRow; row++ {
		var rowCards []string
		for col := 0; col < cols; col++ {
			idx := row*cols + col
			if idx >= len(filtered) {
				rowCards = append(rowCards, strings.Repeat(" ", cardW))
				continue
			}
			d := filtered[idx]
			rowCards = append(rowCards, m.renderDecisionCard(d, idx, cardW))
		}
		joined := lipgloss.JoinHorizontal(lipgloss.Top, rowCards...)
		gridLines = append(gridLines, strings.Split(joined, "\n")...)
	}

	return m.wrapCardScrollArea(gridLines, scrollOff, totalRows, visibleRows, len(filtered))
}

func (m model) renderDecisionCard(d api.Decision, idx, cardW int) string {
	statusColor := dimColor
	switch d.Status {
	case "ACCEPTED":
		statusColor = greenColor
	case "PENDING":
		statusColor = yellowColor
	case "DEPRECATED":
		statusColor = redColor
	case "DRAFT":
		statusColor = blueColor
	}

	innerW := cardW - 4
	if innerW < 10 {
		innerW = 10
	}

	created := ""
	if d.CreatedAt != "" && len(d.CreatedAt) >= 10 {
		created = d.CreatedAt[:10]
	}

	header := strings.TrimSpace(strings.TrimSpace(d.Status) + " " + created)
	if header == "" {
		header = "DRAFT"
	}
	title := d.Statement
	if strings.TrimSpace(title) == "" {
		title = "(no statement)"
	}
	subtitle := d.Rationale
	if strings.TrimSpace(subtitle) == "" {
		subtitle = "No rationale"
	}

	borderColor := lipgloss.AdaptiveColor{Light: "240", Dark: "237"}
	if idx == m.selected {
		borderColor = greenColor
	}
	lines := []string{
		fixedLine("", innerW),
		styledFixedLine(header, innerW, lipgloss.NewStyle().Foreground(statusColor)),
		styledFixedLine(title, innerW, lipgloss.NewStyle().Bold(true)),
		styledFixedLine(subtitle, innerW, dimStyle),
		fixedLine("", innerW),
		fixedLine("", innerW),
	}
	return manualCard(lines, innerW, borderColor)
}

func (m model) renderDecisionDetail(d api.Decision) string {
	var s strings.Builder
	statusIcon := "○"
	statusColor := dimColor
	switch d.Status {
	case "ACCEPTED":
		statusIcon = "●"
		statusColor = greenColor
	case "PENDING":
		statusIcon = "◐"
		statusColor = yellowColor
	case "DRAFT":
		statusIcon = "○"
		statusColor = blueColor
	case "DEPRECATED":
		statusIcon = "x"
		statusColor = redColor
	}
	s.WriteString("  " + lipgloss.NewStyle().Foreground(statusColor).Bold(true).Render(statusIcon+" "+d.Status) + "\n\n")
	s.WriteString("  " + titleStyle.Render("Statement") + "\n")
	s.WriteString("  " + wordWrap(d.Statement, 60) + "\n\n")
	if d.Rationale != "" {
		s.WriteString("  " + titleStyle.Render("Rationale") + "\n")
		s.WriteString("  " + wordWrap(d.Rationale, 60) + "\n\n")
	}
	s.WriteString("  " + dimStyle.Render(fmt.Sprintf("ID: %s", d.ID)) + "\n")
	if d.CreatedAt != "" && len(d.CreatedAt) >= 10 {
		s.WriteString("  " + dimStyle.Render(fmt.Sprintf("Created: %s", d.CreatedAt[:10])) + "\n")
	}
	s.WriteString("\n  " + dimStyle.Render("[a] accept  [x] deprecate  [d] deprecate  [esc] back") + "\n")
	return s.String()
}

// ============================================================================
// MEMORIES — card grid + detail
// ============================================================================

func (m model) renderMemoriesContent() string {
	if m.detailView {
		mem := m.getSelectedMemory()
		if mem != nil {
			return m.renderMemoryDetail(mem)
		}
		m.detailView = false
	}

	var s strings.Builder
	s.WriteString(m.renderViewModeHeader("Memories"))

	if m.searchMode {
		s.WriteString("  Search: " + m.searchQuery + "\u2588\n")
	} else if m.searchQuery != "" {
		s.WriteString("  " + dimStyle.Render("Filter: "+m.searchQuery) + "\n")
	}

	filtered := m.getFilteredMemories()
	s.WriteString("  " + dimStyle.Render(fmt.Sprintf("%d items", len(filtered))) + "\n\n")

	if len(filtered) == 0 {
		if m.searchQuery != "" {
			s.WriteString("  " + dimStyle.Render("No matching memories.") + "\n")
		} else {
			s.WriteString("  " + dimStyle.Render("No memories yet. Press [n] to create one.") + "\n")
		}
		return s.String()
	}

	if m.currentViewMode == viewModeList {
		s.WriteString(m.renderMemoriesList(filtered))
	} else {
		s.WriteString(m.renderMemoriesCards(filtered))
	}

	return s.String()
}

func (m model) renderMemoriesCards(filtered []*api.Memory) string {
	cols := m.getGridCols()
	cardW := m.getCardWidth()
	visibleRows := m.getVisibleRows()
	totalRows := (len(filtered) + cols - 1) / cols

	scrollOff := m.scrollOffset
	if scrollOff > totalRows-visibleRows {
		scrollOff = max(0, totalRows-visibleRows)
	}
	startRow := scrollOff
	endRow := min(startRow+visibleRows, totalRows)

	var gridLines []string
	for row := startRow; row < endRow; row++ {
		var rowCards []string
		for col := 0; col < cols; col++ {
			idx := row*cols + col
			if idx >= len(filtered) {
				rowCards = append(rowCards, strings.Repeat(" ", cardW))
				continue
			}
			mem := filtered[idx]
			rowCards = append(rowCards, m.renderMemoryCard(mem, idx, cardW))
		}
		joined := lipgloss.JoinHorizontal(lipgloss.Top, rowCards...)
		gridLines = append(gridLines, strings.Split(joined, "\n")...)
	}

	return m.wrapCardScrollArea(gridLines, scrollOff, totalRows, visibleRows, len(filtered))
}

func (m model) renderMemoryCard(mem *api.Memory, idx, cardW int) string {
	innerW := cardW - 4
	if innerW < 10 {
		innerW = 10
	}

	created := ""
	if mem.CreatedAt != "" && len(mem.CreatedAt) >= 10 {
		created = mem.CreatedAt[:10]
	}

	content := mem.Content
	if strings.TrimSpace(content) == "" {
		content = "(empty memory)"
	}
	first := content
	second := ""
	r := []rune(content)
	if len(r) > 0 {
		buf := make([]rune, 0, len(r))
		for _, rr := range r {
			candidate := string(append(buf, rr))
			if lipgloss.Width(candidate) > innerW {
				break
			}
			buf = append(buf, rr)
		}
		first = string(buf)
		if len(buf) < len(r) {
			second = string(r[len(buf):])
		}
	}
	if strings.TrimSpace(second) == "" {
		second = " "
	}

	borderColor := lipgloss.AdaptiveColor{Light: "240", Dark: "237"}
	if idx == m.selected {
		borderColor = greenColor
	}
	lines := []string{
		fixedLine("", innerW),
		styledFixedLine("Memory "+created, innerW, lipgloss.NewStyle().Foreground(greenColor)),
		styledFixedLine(first, innerW, normalStyle),
		styledFixedLine(second, innerW, dimStyle),
		fixedLine("", innerW),
		fixedLine("", innerW),
	}
	return manualCard(lines, innerW, borderColor)
}

func (m model) renderMemoryDetail(mem *api.Memory) string {
	var s strings.Builder
	s.WriteString("  " + titleStyle.Render("Memory") + "\n\n")
	s.WriteString("  " + wordWrap(mem.Content, 60) + "\n\n")
	s.WriteString("  " + dimStyle.Render(fmt.Sprintf("ID: %s", mem.ID)) + "\n")
	if mem.CreatedAt != "" && len(mem.CreatedAt) >= 10 {
		s.WriteString("  " + dimStyle.Render(fmt.Sprintf("Created: %s", mem.CreatedAt[:10])) + "\n")
	}
	s.WriteString("\n  " + dimStyle.Render("[esc] back  [e] edit  [d] delete") + "\n")
	return s.String()
}

// ============================================================================
// CAPSULES — card grid + detail
// ============================================================================

func (m model) renderCapsulesContent() string {
	if m.detailView {
		c := m.getSelectedCapsule()
		if c != nil {
			return m.renderCapsuleDetail(c)
		}
		m.detailView = false
	}

	var s strings.Builder
	s.WriteString(m.renderViewModeHeader("Capsules"))

	if m.searchMode {
		s.WriteString("  Search: " + m.searchQuery + "\u2588\n")
	} else if m.searchQuery != "" {
		s.WriteString("  " + dimStyle.Render("Filter: "+m.searchQuery) + "\n")
	}

	filtered := m.getFilteredCapsules()
	s.WriteString("  " + dimStyle.Render(fmt.Sprintf("%d items", len(filtered))) + "\n\n")

	if len(filtered) == 0 {
		s.WriteString("  " + dimStyle.Render("No capsules found.") + "\n")
		return s.String()
	}

	if m.currentViewMode == viewModeList {
		s.WriteString(m.renderCapsulesList(filtered))
	} else {
		s.WriteString(m.renderCapsulesCards(filtered))
	}

	return s.String()
}

func (m model) renderCapsulesCards(filtered []*api.Capsule) string {
	cols := m.getGridCols()
	cardW := m.getCardWidth()
	visibleRows := m.getVisibleRows()
	totalRows := (len(filtered) + cols - 1) / cols

	scrollOff := m.scrollOffset
	if scrollOff > totalRows-visibleRows {
		scrollOff = max(0, totalRows-visibleRows)
	}
	startRow := scrollOff
	endRow := min(startRow+visibleRows, totalRows)

	var gridLines []string
	for row := startRow; row < endRow; row++ {
		var rowCards []string
		for col := 0; col < cols; col++ {
			idx := row*cols + col
			if idx >= len(filtered) {
				rowCards = append(rowCards, strings.Repeat(" ", cardW))
				continue
			}
			c := filtered[idx]
			rowCards = append(rowCards, m.renderCapsuleCard(c, idx, cardW))
		}
		joined := lipgloss.JoinHorizontal(lipgloss.Top, rowCards...)
		gridLines = append(gridLines, strings.Split(joined, "\n")...)
	}

	return m.wrapCardScrollArea(gridLines, scrollOff, totalRows, visibleRows, len(filtered))
}

func (m model) renderCapsuleCard(cap *api.Capsule, idx, cardW int) string {
	innerW := cardW - 4
	if innerW < 10 {
		innerW = 10
	}
	statusColor := dimColor
	if cap.Status == "FROZEN" {
		statusColor = cyanColor
	} else if cap.Status == "HISTORICAL" {
		statusColor = magentaColor
	}

	created := ""
	if cap.CreatedAt != "" && len(cap.CreatedAt) >= 10 {
		created = cap.CreatedAt[:10]
	}
	nDec := len(cap.DecisionIds)
	nMem := len(cap.MemoryIds)

	raw1 := fixedLine(cap.Status+" "+created, innerW)
	raw2 := fixedLine(cap.Name, innerW)
	raw3Src := cap.Description
	if raw3Src == "" {
		raw3Src = fmt.Sprintf("%dd  %dm", nDec, nMem)
	}
	raw3 := fixedLine(raw3Src, innerW)

	line1 := lipgloss.NewStyle().Foreground(statusColor).Render(raw1)
	line2 := lipgloss.NewStyle().Bold(true).Render(raw2)
	line3 := dimStyle.Render(raw3)

	borderColor := lipgloss.AdaptiveColor{Light: "240", Dark: "237"}
	if idx == m.selected {
		borderColor = greenColor
	}
	return manualCard([]string{line1, line2, line3}, innerW, borderColor)
}

func (m model) renderCapsuleDetail(c *api.Capsule) string {
	var s strings.Builder
	icon := "○"
	if c.Status == "FROZEN" {
		icon = "*"
	} else if c.Status == "HISTORICAL" {
		icon = "~"
	}
	s.WriteString("  " + titleStyle.Render(icon+" "+c.Name) + "  " + dimStyle.Render(c.Status) + "\n\n")
	if c.Description != "" {
		s.WriteString("  " + titleStyle.Render("Description") + "\n")
		s.WriteString("  " + wordWrap(c.Description, 60) + "\n\n")
	}
	s.WriteString("  " + titleStyle.Render("Contents") + "\n")
	s.WriteString("  " + dimStyle.Render(fmt.Sprintf("Decisions: %d", len(c.DecisionIds))) + "\n")
	s.WriteString("  " + dimStyle.Render(fmt.Sprintf("Memories: %d", len(c.MemoryIds))) + "\n\n")
	s.WriteString("  " + dimStyle.Render(fmt.Sprintf("ID: %s", c.ID)) + "\n")
	if c.CreatedAt != "" && len(c.CreatedAt) >= 10 {
		s.WriteString("  " + dimStyle.Render(fmt.Sprintf("Created: %s", c.CreatedAt[:10])) + "\n")
	}
	if c.FrozenAt != nil && *c.FrozenAt != "" && len(*c.FrozenAt) >= 10 {
		s.WriteString("  " + dimStyle.Render(fmt.Sprintf("Frozen: %s", (*c.FrozenAt)[:10])) + "\n")
	}
	s.WriteString("\n  " + dimStyle.Render("[esc/q] back") + "\n")
	return s.String()
}

// ============================================================================
// VIEW MODE HEADER — shared Card/List toggle display
// ============================================================================

func (m model) renderViewModeHeader(title string) string {
	var s strings.Builder
	cardLabel := "Cards"
	listLabel := "List"
	if m.currentViewMode == viewModeCard {
		s.WriteString("  " + selectedStyle.Render("["+cardLabel+"]") + " " + dimStyle.Render(listLabel))
	} else {
		s.WriteString("  " + dimStyle.Render(cardLabel) + " " + selectedStyle.Render("["+listLabel+"]"))
	}
	s.WriteString("  " + dimStyle.Render("(press [v] to toggle)") + "\n")
	return s.String()
}

// ============================================================================
// CARD HELPERS
// ============================================================================

func fixedLine(text string, w int) string {
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")
	text = strings.ReplaceAll(text, "\t", " ")
	text = strings.Join(strings.Fields(text), " ")
	tw := lipgloss.Width(text)
	if tw > w {
		r := []rune(text)
		for len(r) > 0 && lipgloss.Width(string(r))+3 > w {
			r = r[:len(r)-1]
		}
		text = string(r) + "..."
		tw = lipgloss.Width(text)
	}
	if tw < w {
		return text + strings.Repeat(" ", w-tw)
	}
	return text
}

func styledFixedLine(text string, w int, style lipgloss.Style) string {
	return style.Render(fixedLine(text, w))
}

func manualCard(lines []string, innerW int, borderColor lipgloss.TerminalColor) string {
	bc := lipgloss.NewStyle().Foreground(borderColor)
	top := bc.Render("╭" + strings.Repeat("─", innerW+2) + "╮")
	bot := bc.Render("╰" + strings.Repeat("─", innerW+2) + "╯")
	vl := bc.Render("│")

	vlW := lipgloss.Width(vl)
	totalW := lipgloss.Width(top)
	contentArea := totalW - (vlW * 2) - 2

	var sb strings.Builder
	sb.WriteString(top + "\n")
	for _, line := range lines {
		lineW := lipgloss.Width(line)
		pad := contentArea - lineW
		if pad < 0 {
			pad = 0
		}
		sb.WriteString(vl + " " + line + strings.Repeat(" ", pad) + " " + vl + "\n")
	}
	sb.WriteString(bot)
	return sb.String()
}

// ============================================================================
// SCROLL AREA CONTAINER
// ============================================================================

// wrapCardScrollArea wraps grid lines in a bordered scroll viewport with a
// scrollbar on the right side. The viewport shows only the visible rows of
// cards, and a proportional thumb on the scrollbar indicates position.
func (m model) wrapCardScrollArea(gridLines []string, scrollOff, totalRows, visibleRows, totalItems int) string {
	if len(gridLines) == 0 {
		return ""
	}

	var s strings.Builder

	// Compute viewport height (number of terminal lines for visible rows)
	vpHeight := len(gridLines)

	// Compute content width from the widest grid line
	contentW := 0
	for _, line := range gridLines {
		w := lipgloss.Width(line)
		if w > contentW {
			contentW = w
		}
	}

	// Available area for scroll area border
	areaW := m.width - 6
	if areaW < contentW+4 {
		areaW = contentW + 4
	}
	innerW := areaW - 4 // border(2) + scrollbar(1) + gap(1)

	// Build the scrollbar
	canScroll := totalRows > visibleRows
	scrollBarHeight := vpHeight
	if scrollBarHeight < 1 {
		scrollBarHeight = 1
	}

	thumbSize := 1
	thumbPos := 0
	if canScroll && totalRows > 0 {
		thumbSize = max(1, (visibleRows*scrollBarHeight)/totalRows)
		maxOff := totalRows - visibleRows
		if maxOff > 0 {
			thumbPos = (scrollOff * (scrollBarHeight - thumbSize)) / maxOff
		}
		if thumbPos+thumbSize > scrollBarHeight {
			thumbPos = scrollBarHeight - thumbSize
		}
		if thumbPos < 0 {
			thumbPos = 0
		}
	}

	trackStyle := lipgloss.NewStyle().Foreground(scrollBarTrackColor)
	thumbStyle := lipgloss.NewStyle().Foreground(scrollBarActiveColor)

	// Border chars
	borderStyle := lipgloss.NewStyle().Foreground(scrollAreaBorderColor)
	topBorder := borderStyle.Render("╭" + strings.Repeat("─", innerW+2) + "─╮")
	botBorder := borderStyle.Render("╰" + strings.Repeat("─", innerW+2) + "─╯")
	vBar := borderStyle.Render("│")

	s.WriteString("  " + topBorder + "\n")

	for i := 0; i < vpHeight; i++ {
		line := ""
		if i < len(gridLines) {
			line = gridLines[i]
		}

		// Pad line to innerW
		lineW := lipgloss.Width(line)
		if lineW < innerW {
			line = line + strings.Repeat(" ", innerW-lineW)
		}

		// Scrollbar character for this row
		var scrollChar string
		if canScroll {
			if i >= thumbPos && i < thumbPos+thumbSize {
				scrollChar = thumbStyle.Render("█")
			} else {
				scrollChar = trackStyle.Render("░")
			}
		} else {
			scrollChar = " "
		}

		s.WriteString("  " + vBar + " " + line + " " + scrollChar + vBar + "\n")
	}

	s.WriteString("  " + botBorder + "\n")

	// Scroll info line
	shown := min((scrollOff+visibleRows)*m.getGridCols(), totalItems)
	if shown > totalItems {
		shown = totalItems
	}
	if canScroll {
		s.WriteString("  " + dimStyle.Render(fmt.Sprintf(
			"Showing %d of %d  ·  ↑↓ j/k scroll  ·  Row %d-%d of %d",
			shown, totalItems, scrollOff+1, min(scrollOff+visibleRows, totalRows), totalRows,
		)) + "\n")
	} else {
		s.WriteString("  " + dimStyle.Render(fmt.Sprintf("%d items", totalItems)) + "\n")
	}

	return s.String()
}
