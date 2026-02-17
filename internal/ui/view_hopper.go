package ui

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ============================================================================
// HOPPER (AI CHAT) VIEW
// ============================================================================

func (m model) renderHopperContent() string {
	if m.loading && !m.hopperContextLoaded {
		return "  " + dimStyle.Render("Loading project context...") + "\n"
	}

	contentH := m.height - 8
	if contentH < 10 {
		contentH = 10
	}

	sidebarW := 28
	if m.width < 80 {
		sidebarW = 22
	}
	if m.width < 60 {
		sidebarW = 0
	}

	chatW := m.width - sidebarW - 4
	if chatW < 30 {
		chatW = 30
	}

	sidebar := m.renderHopperSidebar(sidebarW, contentH)
	chatPanel := m.renderHopperChatPanel(chatW, contentH)

	if sidebarW == 0 {
		return chatPanel
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, chatPanel)
}

func (m model) renderHopperSidebar(w, h int) string {
	sidebarBorderColor := lipgloss.AdaptiveColor{Light: "240", Dark: "237"}
	if m.hopperSidebarFocus {
		sidebarBorderColor = purpleColor
	}

	var sb strings.Builder

	newChatLabel := " + New Chat"
	if m.hopperSidebarFocus {
		newChatLabel = lipgloss.NewStyle().Foreground(greenColor).Render(" + New Chat")
	}
	sb.WriteString(newChatLabel + "\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "240", Dark: "237"}).Render(strings.Repeat("─", w-4)) + "\n")

	if len(m.chatSessions) == 0 {
		sb.WriteString("\n" + dimStyle.Render(" No chats yet"))
	} else {
		listH := h - 5
		if listH < 1 {
			listH = 1
		}
		startIdx := 0
		if m.chatSessionIdx >= listH {
			startIdx = m.chatSessionIdx - listH + 1
		}
		endIdx := min(startIdx+listH, len(m.chatSessions))

		for i := startIdx; i < endIdx; i++ {
			sess := m.chatSessions[i]
			topic := truncateString(sess.Topic, w-8)
			if i == m.chatSessionIdx && m.hopperSidebarFocus {
				sb.WriteString(lipgloss.NewStyle().Foreground(purpleColor).Bold(true).Render(" > " + topic) + "\n")
				sb.WriteString("   " + dimStyle.Render(sess.Time) + "\n")
			} else if sess.ID == m.hopperSessionID {
				sb.WriteString(lipgloss.NewStyle().Foreground(greenColor).Render(" * " + topic) + "\n")
				sb.WriteString("   " + dimStyle.Render(sess.Time) + "\n")
			} else {
				sb.WriteString(" " + dimStyle.Render("  "+topic) + "\n")
				sb.WriteString("   " + dimStyle.Render(sess.Time) + "\n")
			}
		}
	}

	lines := strings.Count(sb.String(), "\n")
	for i := lines; i < h-2; i++ {
		sb.WriteString("\n")
	}

	style := lipgloss.NewStyle().
		Width(w - 2).
		Height(h - 2).
		MaxHeight(h - 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(sidebarBorderColor).
		Padding(0, 0)

	return style.Render(sb.String())
}

func (m model) renderHopperChatPanel(w, h int) string {
	chatBorderColor := lipgloss.AdaptiveColor{Light: "240", Dark: "237"}
	if !m.hopperSidebarFocus {
		chatBorderColor = purpleColor
	}

	var panel strings.Builder
	msgW := w - 6
	if msgW < 20 {
		msgW = 20
	}

	inputH := 3
	msgAreaH := h - inputH - 4
	if msgAreaH < 4 {
		msgAreaH = 4
	}

	var msgLines []string

	if len(m.chatMessages) == 0 {
		welcomeLines := []string{
			"",
			lipgloss.NewStyle().Bold(true).Foreground(purpleColor).Render("  Hopper AI Assistant"),
			"",
			dimStyle.Render("  I have context about your project's"),
			dimStyle.Render("  decisions, memories and capsules."),
			"",
			dimStyle.Render("  Try asking:"),
			normalStyle.Render("  > What are the key decisions?"),
			normalStyle.Render("  > Summarize recent changes"),
			normalStyle.Render("  > What should I work on next?"),
			"",
		}
		msgLines = welcomeLines
	} else {
		for _, msg := range m.chatMessages {
			if msg.Role == "user" {
				msgLines = append(msgLines, "")
				userLabel := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.AdaptiveColor{Light: "25", Dark: "75"}).Render("You")
				msgLines = append(msgLines, "  "+userLabel)

				userBubbleW := msgW - 4
				if userBubbleW > 60 {
					userBubbleW = 60
				}
				wrapped := wordWrap(msg.Content, userBubbleW)
				bubbleStyle := lipgloss.NewStyle().
					Foreground(lipgloss.AdaptiveColor{Light: "15", Dark: "255"}).
					Background(lipgloss.AdaptiveColor{Light: "25", Dark: "57"}).
					Padding(0, 1)

				for _, line := range strings.Split(wrapped, "\n") {
					msgLines = append(msgLines, "  "+bubbleStyle.Render(line))
				}
			} else {
				msgLines = append(msgLines, "")
				hopperLabel := statusOnStyle.Render("●") + " " + lipgloss.NewStyle().Bold(true).Foreground(purpleColor).Render("Hopper")
				msgLines = append(msgLines, "  "+hopperLabel)

				content := msg.Content
				rendered := m.renderHopperAssistantContent(content, msgW-4)
				for _, line := range strings.Split(rendered, "\n") {
					msgLines = append(msgLines, "  "+line)
				}
			}
		}
	}

	if m.chatStreaming {
		msgLines = append(msgLines, "")
		hopperLabel := statusOnStyle.Render("●") + " " + lipgloss.NewStyle().Bold(true).Foreground(purpleColor).Render("Hopper")
		msgLines = append(msgLines, "  "+hopperLabel)
		msgLines = append(msgLines, "  "+dimStyle.Render("thinking..."))
	}

	totalLines := len(msgLines)
	if totalLines > msgAreaH {
		maxScroll := totalLines - msgAreaH
		scrollPos := maxScroll
		if m.chatScroll < maxScroll {
			scrollPos = max(0, m.chatScroll)
		}
		msgLines = msgLines[scrollPos : scrollPos+msgAreaH]
		if scrollPos > 0 {
			msgLines[0] = dimStyle.Render(fmt.Sprintf("  ^ %d more lines (ctrl+u scroll up)", scrollPos))
		}
	}
	for i := len(msgLines); i < msgAreaH; i++ {
		msgLines = append(msgLines, "")
	}
	for _, line := range msgLines {
		panel.WriteString(line + "\n")
	}

	panel.WriteString(lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "240", Dark: "237"}).Render(strings.Repeat("─", w-4)) + "\n")

	inputPrefix := dimStyle.Render("> ")
	if m.chatStreaming {
		inputPrefix = dimStyle.Render("  ")
	}
	panel.WriteString(inputPrefix + m.chatInput)
	if !m.chatStreaming {
		panel.WriteString("\u2588")
	}
	panel.WriteString("\n")

	style := lipgloss.NewStyle().
		Width(w - 2).
		Height(h - 2).
		MaxHeight(h - 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(chatBorderColor).
		Padding(0, 0)

	return style.Render(panel.String())
}

// ============================================================================
// HOPPER CONTENT PARSING
// ============================================================================

func (m model) renderHopperAssistantContent(content string, maxW int) string {
	var result strings.Builder
	bubbleW := maxW
	if bubbleW > 72 {
		bubbleW = 72
	}

	lines := strings.Split(content, "\n")
	inBlock := false
	blockType := ""
	var blockLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, ":::DECISIONS") || strings.HasPrefix(trimmed, ":::MEMORIES") || strings.HasPrefix(trimmed, ":::CAPSULES") || strings.HasPrefix(trimmed, ":::TASKS") {
			inBlock = true
			blockType = strings.TrimPrefix(trimmed, ":::")
			if idx := strings.Index(blockType, "\n"); idx >= 0 {
				blockType = blockType[:idx]
			}
			blockLines = nil
			continue
		}
		if trimmed == ":::END" && inBlock {
			inBlock = false
			result.WriteString(m.renderHopperDataBlock(blockType, strings.Join(blockLines, "\n"), bubbleW))
			blockLines = nil
			continue
		}
		if inBlock {
			blockLines = append(blockLines, line)
			continue
		}

		line = renderInlineBadges(line)
		formatted := formatMarkdownLineForHopper(line, bubbleW)
		result.WriteString(formatted + "\n")
	}

	return strings.TrimRight(result.String(), "\n")
}

func renderInlineBadges(line string) string {
	badges := map[string]lipgloss.AdaptiveColor{
		"[ACCEPTED]":    greenColor,
		"[PENDING]":     lipgloss.AdaptiveColor{Light: "214", Dark: "214"},
		"[DRAFT]":       lipgloss.AdaptiveColor{Light: "240", Dark: "240"},
		"[REJECTED]":    lipgloss.AdaptiveColor{Light: "196", Dark: "196"},
		"[DEPRECATED]":  lipgloss.AdaptiveColor{Light: "196", Dark: "196"},
		"[TODO]":        lipgloss.AdaptiveColor{Light: "214", Dark: "214"},
		"[IN_PROGRESS]": lipgloss.AdaptiveColor{Light: "75", Dark: "75"},
		"[DONE]":        greenColor,
		"[HIGH]":        lipgloss.AdaptiveColor{Light: "196", Dark: "196"},
		"[MEDIUM]":      lipgloss.AdaptiveColor{Light: "214", Dark: "214"},
		"[LOW]":         lipgloss.AdaptiveColor{Light: "240", Dark: "240"},
	}
	for badge, color := range badges {
		if strings.Contains(line, badge) {
			styled := lipgloss.NewStyle().Foreground(color).Bold(true).Render(badge)
			line = strings.ReplaceAll(line, badge, styled)
		}
	}
	return line
}

func formatMarkdownLineForHopper(line string, maxW int) string {
	trimmed := strings.TrimSpace(line)

	if strings.HasPrefix(trimmed, "### ") {
		text := strings.TrimPrefix(trimmed, "### ")
		return lipgloss.NewStyle().Bold(true).Foreground(purpleColor).Render(text)
	}
	if strings.HasPrefix(trimmed, "## ") {
		text := strings.TrimPrefix(trimmed, "## ")
		return lipgloss.NewStyle().Bold(true).Foreground(purpleColor).Render(text)
	}
	if strings.HasPrefix(trimmed, "# ") {
		text := strings.TrimPrefix(trimmed, "# ")
		return lipgloss.NewStyle().Bold(true).Foreground(purpleColor).Render(text)
	}

	if strings.HasPrefix(trimmed, "```") {
		return dimStyle.Render("───")
	}

	if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
		text := trimmed[2:]
		text = strings.ReplaceAll(text, "**", "")
		text = strings.ReplaceAll(text, "*", "")
		wrapped := wordWrap(text, maxW-4)
		var bulletLines []string
		for i, l := range strings.Split(wrapped, "\n") {
			if i == 0 {
				bulletLines = append(bulletLines, dimStyle.Render("  - ")+normalStyle.Render(l))
			} else {
				bulletLines = append(bulletLines, "    "+normalStyle.Render(l))
			}
		}
		return strings.Join(bulletLines, "\n")
	}

	if len(trimmed) > 2 && trimmed[0] >= '0' && trimmed[0] <= '9' && (trimmed[1] == '.' || (len(trimmed) > 2 && trimmed[1] >= '0' && trimmed[1] <= '9' && trimmed[2] == '.')) {
		text := strings.ReplaceAll(trimmed, "**", "")
		wrapped := wordWrap(text, maxW-2)
		return normalStyle.Render(wrapped)
	}

	if strings.HasPrefix(trimmed, "> ") {
		text := strings.TrimPrefix(trimmed, "> ")
		return lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "240", Dark: "242"}).Italic(true).Render("  " + text)
	}

	text := strings.ReplaceAll(trimmed, "**", "")
	text = strings.ReplaceAll(text, "*", "")
	if text == "" {
		return ""
	}
	wrapped := wordWrap(text, maxW)
	return normalStyle.Render(wrapped)
}

func (m model) renderHopperDataBlock(blockType, jsonData string, maxW int) string {
	var result strings.Builder

	typeColor := purpleColor
	typeLabel := blockType
	switch blockType {
	case "DECISIONS":
		typeColor = lipgloss.AdaptiveColor{Light: "25", Dark: "75"}
		typeLabel = "Decisions"
	case "MEMORIES":
		typeColor = purpleColor
		typeLabel = "Memories"
	case "CAPSULES":
		typeColor = greenColor
		typeLabel = "Capsules"
	case "TASKS":
		typeColor = lipgloss.AdaptiveColor{Light: "214", Dark: "214"}
		typeLabel = "Tasks"
	}

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(typeColor)
	result.WriteString("\n" + headerStyle.Render("  "+typeLabel) + "\n")

	cardBorderColor := lipgloss.AdaptiveColor{Light: "240", Dark: "237"}

	items := parseHopperBlockItems(jsonData)
	if len(items) == 0 {
		result.WriteString("  " + dimStyle.Render("(no items)") + "\n")
		return result.String()
	}

	cols := 3
	if maxW < 80 {
		cols = 2
	}
	if maxW < 50 {
		cols = 1
	}

	cardW := (maxW - 4) / cols
	if cardW < 20 {
		cardW = 20
	}
	contentW := cardW - 4
	if contentW < 10 {
		contentW = 10
	}

	for i := 0; i < len(items); i += cols {
		var rowCards []string
		for col := 0; col < cols; col++ {
			idx := i + col
			if idx >= len(items) {
				rowCards = append(rowCards, strings.Repeat(" ", cardW))
				continue
			}
			item := items[idx]

			var cardLines []string
			nameStyle := lipgloss.NewStyle().Bold(true)
			cardLines = append(cardLines, nameStyle.Render(fixedLine(item.Name, contentW)))

			if item.Status != "" {
				statusCol := dimStyle
				switch strings.ToUpper(item.Status) {
				case "ACCEPTED", "DONE":
					statusCol = lipgloss.NewStyle().Foreground(greenColor)
				case "PENDING", "IN_PROGRESS", "TODO":
					statusCol = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "214", Dark: "214"})
				case "DRAFT":
					statusCol = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "75", Dark: "75"})
				case "REJECTED", "DEPRECATED":
					statusCol = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "196", Dark: "196"})
				}
				cardLines = append(cardLines, statusCol.Render(fixedLine(item.Status, contentW)))
			}
			if item.Content != "" {
				cardLines = append(cardLines, dimStyle.Render(fixedLine(item.Content, contentW)))
			}

			rowCards = append(rowCards, manualCard(cardLines, contentW, cardBorderColor))
		}
		joined := lipgloss.JoinHorizontal(lipgloss.Top, rowCards...)
		for _, line := range strings.Split(joined, "\n") {
			result.WriteString("  " + line + "\n")
		}
	}

	return result.String()
}

type hopperBlockItem struct {
	Name    string
	Status  string
	Content string
}

func parseHopperBlockItems(jsonData string) []hopperBlockItem {
	jsonData = strings.TrimSpace(jsonData)
	if jsonData == "" {
		return nil
	}

	var rawItems []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &rawItems); err == nil {
		var items []hopperBlockItem
		for _, raw := range rawItems {
			item := hopperBlockItem{}
			for _, key := range []string{"statement", "name", "title", "content"} {
				if val, ok := raw[key]; ok {
					if s, ok := val.(string); ok && s != "" {
						item.Name = s
						break
					}
				}
			}
			for _, key := range []string{"status", "sourceStatus"} {
				if val, ok := raw[key]; ok {
					if s, ok := val.(string); ok && s != "" {
						item.Status = s
						break
					}
				}
			}
			for _, key := range []string{"rationale", "description", "content"} {
				if val, ok := raw[key]; ok {
					if s, ok := val.(string); ok && s != "" && s != item.Name {
						item.Content = s
						break
					}
				}
			}
			if item.Name != "" {
				items = append(items, item)
			}
		}
		return items
	}

	return []hopperBlockItem{{Name: truncateString(jsonData, 80)}}
}
