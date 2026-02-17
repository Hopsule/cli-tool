package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Hopsule/cli-tool/internal/api"
	"github.com/Hopsule/cli-tool/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

// ============================================================================
// UPDATE — main Bubble Tea update loop
// ============================================================================

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case dataLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.organizations = msg.organizations
			m.projects = msg.projects

			// Auto-match: compare local git remote "owner/repo" with API's github_repo_full_name
			if m.initializedProjectID == "" && m.cwdInGitRepo {
				matched := false

				// Priority 1: git remote full name (most reliable)
				if m.cwdGitRemoteFullName != "" {
					for _, proj := range m.projects {
						if proj.GitHubRepoFullName != "" && strings.EqualFold(proj.GitHubRepoFullName, m.cwdGitRemoteFullName) {
							for _, org := range m.organizations {
								if org.ID == proj.OrganizationID {
									m.initProjectLocally(proj, org)
									matched = true
									break
								}
							}
							break
						}
					}
				}

				// Priority 2: fallback to directory name match
				if !matched && m.cwdDirName != "" {
					dirLower := strings.ToLower(m.cwdDirName)
					for _, proj := range m.projects {
						if strings.ToLower(proj.Slug) == dirLower || strings.ToLower(proj.Name) == dirLower {
							for _, org := range m.organizations {
								if org.ID == proj.OrganizationID {
									m.initProjectLocally(proj, org)
									break
								}
							}
							break
						}
					}
				}
			}

			// Auto-navigate to initialized project if .hopsule exists (or was just created)
			if m.initializedProjectID != "" {
				for _, proj := range m.projects {
					if proj.ID == m.initializedProjectID {
						for oi, org := range m.organizations {
							if org.ID == proj.OrganizationID {
								m.orgSelectedIdx = oi
								m.currentOrg = org
								m.currentView = viewProjects

								orgProjects := m.getOrgProjects()
								for pi, op := range orgProjects {
									if op.ID == proj.ID {
										m.projSelectedIdx = pi
										m.currentProj = op
										m.menuItems = buildProjectMenuItems()
										m.currentView = viewProjectMenu
										m.selected = 0
										break
									}
								}
								break
							}
						}
						break
					}
				}
			}
		}
		return m, nil
	case loginCompleteMsg:
		m.loginPolling = false
		m.loginDeviceCode = ""
		m.loginAuthURL = ""
		if msg.success {
			if m.cfg == nil {
				m.cfg = &config.Config{}
			}
			m.cfg.Token = msg.token
			m.cfg.User = msg.user
			if m.cfg.APIURL == "" {
				m.cfg.APIURL = "https://api.hopsule.com"
			}
			if m.cfg.WebURL == "" {
				m.cfg.WebURL = "https://app.hopsule.com"
			}
			config.SaveConfig(m.cfg)
			m.client = api.NewClient(m.cfg)
			m.currentView = viewOrganizations
			m.loading = true
			m.selected = 0
			return m, m.loadData
		} else if msg.err != nil {
			m.errorMsg = msg.err.Error()
		}
		return m, nil
	case loginInitMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			m.loginPolling = false
			return m, nil
		}
		m.loginDeviceCode = msg.deviceCode
		m.loginAuthURL = msg.authURL
		m.loginPolling = true
		m.loginPollCount = 0
		return m, loginTickCmd()
	case loginPollMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			m.loginPolling = false
			return m, nil
		}
		switch msg.status {
		case "complete":
			return m, func() tea.Msg {
				return loginCompleteMsg{success: true, token: msg.token, user: msg.user}
			}
		case "expired":
			m.loginPolling = false
			m.errorMsg = "Login session expired - press Enter to try again"
			return m, nil
		default:
			return m, loginTickCmd()
		}
	case loginTickMsg:
		if !m.loginPolling {
			return m, nil
		}
		m.loginSpinnerIdx = (m.loginSpinnerIdx + 1) % 10
		m.loginPollCount++
		if m.loginPollCount > 150 {
			m.loginPolling = false
			m.errorMsg = "Login timed out - press Enter to try again"
			return m, nil
		}
		return m, m.pollLogin
	case decisionsLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.decisions = msg.decisions
		}
		return m, nil
	case memoriesLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.memories = msg.memories
		}
		return m, nil
	case tasksLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.tasks = msg.tasks
		}
		return m, nil
	case capsulesLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.capsules = msg.capsules
		}
		return m, nil
	case brainStatsLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.graphStats = msg.stats
		}
		return m, nil
	case dashboardLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.decisions = msg.decisions
			m.memories = msg.memories
			m.tasks = msg.tasks
			m.capsules = msg.capsules
		}
		return m, nil
	case chatStreamChunkMsg:
		m.chatStreaming = false
		if msg.chunk != "" {
			content := sanitizeMarkdownForTerminal(msg.chunk)
			m.chatMessages = append(m.chatMessages, api.ChatMessage{
				Role: "assistant", Content: content,
			})
			m.chatScroll = max(0, len(m.chatMessages)-4)
			m.hopperSaveCurrentChat()
		}
		return m, m.saveChatHistoryToAPI
	case chatStreamDoneMsg:
		m.chatStreaming = false
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		}
		return m, nil
	case hopperContextLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.hopperDecisions = msg.decisions
			m.hopperMemories = msg.memories
			m.hopperContextLoaded = true
			if len(msg.chatHistory) > 0 {
				for _, ch := range msg.chatHistory {
					exists := false
					for _, s := range m.chatSessions {
						if s.ID == ch.ID {
							exists = true
							break
						}
					}
					if !exists {
						updatedAt := ch.UpdatedAt
						if len(updatedAt) >= 16 {
							updatedAt = updatedAt[11:16]
						}
						m.chatSessions = append(m.chatSessions, chatSession{
							ID:    ch.ID,
							Topic: ch.Topic,
							Time:  updatedAt,
						})
					}
				}
			}
		}
		return m, nil
	case chatHistorySavedMsg:
		if msg.err == nil && msg.chatID != "" {
			for i, s := range m.chatSessions {
				if s.ID == m.hopperSessionID {
					m.chatSessions[i].ID = msg.chatID
					break
				}
			}
			m.hopperSessionID = msg.chatID
		}
		return m, nil
	case chatHistoryLoadedMsg:
		if msg.err == nil {
			for i, s := range m.chatSessions {
				if s.ID == msg.sessionID {
					m.chatSessions[i].Messages = msg.messages
					break
				}
			}
			if m.hopperSessionID == msg.sessionID {
				m.chatMessages = msg.messages
				m.chatScroll = max(0, len(m.chatMessages)*3)
			}
		}
		return m, nil
	}
	return m, nil
}

// ============================================================================
// KEY HANDLING — delegates to sub-handlers
// ============================================================================

func (m model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.errorMsg = ""

	if m.confirmDelete {
		return m.handleConfirmDelete(msg)
	}
	if m.inputMode {
		return m.handleInputMode(msg)
	}
	if m.searchMode {
		return m.handleSearchMode(msg)
	}
	if m.detailView {
		return m.handleDetailView(msg)
	}

	// Hopper chat has its own key handler
	if m.currentView == viewHopper {
		return m.handleHopperKeys(msg)
	}

	goBack := func() (tea.Model, tea.Cmd) {
		switch m.currentView {
		case viewDashboard, viewDecisions, viewMemories, viewCapsules, viewTasks, viewHopper, viewMCP:
			m.currentView = viewProjectMenu
			m.selected = m.menuSelectedIdx
			m.searchQuery = ""
			m.statusFilter = ""
			m.scrollOffset = 0
			m.listPage = 0
			return m, nil
		case viewProjectMenu:
			m.currentView = viewProjects
			m.currentProj = nil
			m.selected = m.projSelectedIdx
			return m, nil
		case viewProjects:
			m.currentView = viewOrganizations
			m.currentOrg = nil
			m.selected = m.orgSelectedIdx
			return m, nil
		}
		return m, nil
	}

	switch msg.String() {
	case "ctrl+c":
		if m.currentView == viewOrganizations || m.currentView == viewLogin {
			return m, tea.Quit
		}
		return goBack()

	case "q":
		if m.currentView == viewOrganizations || m.currentView == viewLogin {
			return m, tea.Quit
		}
		return goBack()

	case "esc":
		if m.searchQuery != "" {
			m.searchQuery = ""
			m.selected = 0
			m.scrollOffset = 0
			m.listPage = 0
			return m, nil
		}
		if m.currentView == viewOrganizations || m.currentView == viewLogin {
			return m, nil
		}
		return goBack()

	case "up", "k":
		if m.isCardView() && (m.currentView == viewDecisions || m.currentView == viewMemories || m.currentView == viewCapsules) {
			cols := m.getGridCols()
			if m.selected >= cols {
				m.selected -= cols
				row := m.selected / cols
				if row < m.scrollOffset {
					m.scrollOffset = row
				}
			}
		} else {
			if m.selected > 0 {
				m.selected--
				if m.selected < m.scrollOffset {
					m.scrollOffset = m.selected
				}
			}
		}

	case "down", "j":
		maxSel := m.getMaxSelection() - 1
		if m.isCardView() && (m.currentView == viewDecisions || m.currentView == viewMemories || m.currentView == viewCapsules) {
			cols := m.getGridCols()
			if m.selected+cols <= maxSel {
				m.selected += cols
			} else if m.selected < maxSel {
				m.selected = maxSel
			}
			row := m.selected / cols
			visibleRows := m.getVisibleRows()
			if row >= m.scrollOffset+visibleRows {
				m.scrollOffset = row - visibleRows + 1
			}
		} else {
			if m.selected < maxSel {
				m.selected++
				visibleItems := m.getVisibleCards()
				if m.selected >= m.scrollOffset+visibleItems {
					m.scrollOffset = m.selected - visibleItems + 1
				}
			}
		}

	case "left", "h":
		if m.isCardView() && (m.currentView == viewDecisions || m.currentView == viewMemories || m.currentView == viewCapsules) {
			if m.selected > 0 {
				m.selected--
				cols := m.getGridCols()
				row := m.selected / cols
				if row < m.scrollOffset {
					m.scrollOffset = row
				}
			}
		} else if !m.isCardView() && (m.currentView == viewDecisions || m.currentView == viewMemories || m.currentView == viewCapsules) {
			// List mode: previous page
			if m.listPage > 0 {
				m.listPage--
				m.selected = 0
				m.scrollOffset = 0
			}
		} else if m.currentView == viewProjects && m.selected > 0 {
			m.selected--
		}

	case "right", "l":
		if m.isCardView() && (m.currentView == viewDecisions || m.currentView == viewMemories || m.currentView == viewCapsules) {
			maxSel := m.getMaxSelection() - 1
			if m.selected < maxSel {
				m.selected++
				cols := m.getGridCols()
				vRows := m.getVisibleRows()
				row := m.selected / cols
				if row >= m.scrollOffset+vRows {
					m.scrollOffset = row - vRows + 1
				}
			}
		} else if !m.isCardView() && (m.currentView == viewDecisions || m.currentView == viewMemories || m.currentView == viewCapsules) {
			// List mode: next page
			filtered := m.getFilteredCount()
			totalPages := (filtered + m.listPageSize - 1) / m.listPageSize
			if m.listPage < totalPages-1 {
				m.listPage++
				m.selected = 0
				m.scrollOffset = 0
			}
		} else if m.currentView == viewProjects {
			m.selected = min(m.selected+1, m.getMaxSelection()-1)
		}

	case "tab":
		if m.currentView == viewTasks {
			if m.tasksViewMode == "list" {
				m.tasksViewMode = "kanban"
			} else {
				m.tasksViewMode = "list"
			}
			return m, nil
		}
		if m.currentView == viewDecisions {
			filters := []string{"", "DRAFT", "PENDING", "ACCEPTED", "DEPRECATED"}
			idx := 0
			for i, f := range filters {
				if f == m.statusFilter {
					idx = i
					break
				}
			}
			m.statusFilter = filters[(idx+1)%len(filters)]
			m.selected = 0
			m.scrollOffset = 0
			m.listPage = 0
			return m, nil
		}

	case "v":
		// Toggle card/list view for decisions, memories, capsules
		if m.currentView == viewDecisions || m.currentView == viewMemories || m.currentView == viewCapsules {
			if m.currentViewMode == viewModeCard {
				m.currentViewMode = viewModeList
			} else {
				m.currentViewMode = viewModeCard
			}
			m.selected = 0
			m.scrollOffset = 0
			m.listPage = 0
			return m, nil
		}

	case "n":
		if m.currentView == viewDecisions {
			m.inputMode = true
			m.inputPrompt = "New Decision Statement:"
			m.inputAction = "create_decision"
			m.inputValue = ""
			return m, nil
		} else if m.currentView == viewMemories {
			m.inputMode = true
			m.inputPrompt = "New Memory Content:"
			m.inputAction = "create_memory"
			m.inputValue = ""
			return m, nil
		} else if m.currentView == viewTasks {
			m.inputMode = true
			m.inputPrompt = "New Task Title:"
			m.inputAction = "create_task"
			m.inputValue = ""
			return m, nil
		}

	case "e":
		if m.currentView == viewMemories {
			mem := m.getSelectedMemory()
			if mem != nil {
				m.inputMode = true
				m.inputPrompt = "Edit Memory:"
				m.inputAction = "edit_memory"
				m.inputValue = mem.Content
				m.editingItemID = mem.ID
				return m, nil
			}
		} else if m.currentView == viewTasks {
			task := m.getSelectedTask()
			if task != nil {
				m.inputMode = true
				m.inputPrompt = "Edit Task Title:"
				m.inputAction = "edit_task"
				m.inputValue = task.Title
				m.editingItemID = task.ID
				return m, nil
			}
		}

	case "a":
		if m.currentView == viewDecisions {
			d := m.getSelectedDecision()
			if d != nil && (d.Status == "DRAFT" || d.Status == "PENDING") {
				userName := "CLI User"
				if m.cfg != nil && m.cfg.User != nil && m.cfg.User.Name != "" {
					userName = m.cfg.User.Name
				}
				_, err := m.client.AcceptDecision(m.currentProj.ID, d.ID, userName)
				if err != nil {
					m.errorMsg = fmt.Sprintf("Failed to accept: %v", err)
				} else {
					m.loading = true
					return m, m.loadDecisions
				}
			} else if d != nil {
				m.errorMsg = "Can only accept DRAFT or PENDING decisions"
			}
		}

	case "x":
		if m.currentView == viewDecisions {
			d := m.getSelectedDecision()
			if d != nil && d.Status == "ACCEPTED" {
				_, err := m.client.DeprecateDecision(m.currentProj.ID, d.ID)
				if err != nil {
					m.errorMsg = fmt.Sprintf("Failed to deprecate: %v", err)
				} else {
					m.loading = true
					return m, m.loadDecisions
				}
			} else if d != nil {
				m.errorMsg = "Can only deprecate ACCEPTED decisions"
			}
		}

	case "d":
		if m.currentView == viewDecisions {
			d := m.getSelectedDecision()
			if d != nil {
				m.confirmDelete = true
				m.confirmDeleteID = d.ID
				m.confirmDeleteType = "decision"
				return m, nil
			}
		} else if m.currentView == viewMemories {
			mem := m.getSelectedMemory()
			if mem != nil {
				m.confirmDelete = true
				m.confirmDeleteID = mem.ID
				m.confirmDeleteType = "memory"
				return m, nil
			}
		} else if m.currentView == viewTasks {
			task := m.getSelectedTask()
			if task != nil {
				m.confirmDelete = true
				m.confirmDeleteID = task.ID
				m.confirmDeleteType = "task"
				return m, nil
			}
		}

	case "t":
		if m.currentView == viewTasks {
			task := m.getSelectedTask()
			if task != nil {
				newStatus := "DONE"
				if task.Status == "DONE" {
					newStatus = "TODO"
				} else if task.Status == "TODO" {
					newStatus = "IN_PROGRESS"
				} else if task.Status == "IN_PROGRESS" {
					newStatus = "DONE"
				}
				_, err := m.client.UpdateTask(m.currentProj.ID, task.ID, api.UpdateTaskRequest{Status: newStatus})
				if err != nil {
					m.errorMsg = fmt.Sprintf("Failed to update: %v", err)
				} else {
					m.loading = true
					return m, m.loadTasks
				}
			}
		}

	case "/":
		if m.currentView == viewDecisions || m.currentView == viewMemories ||
			m.currentView == viewTasks || m.currentView == viewCapsules {
			m.searchMode = true
			m.searchQuery = ""
			return m, nil
		}

	case "enter", " ":
		if m.currentView == viewMCP {
			if m.selected < len(m.mcpIDEs) {
				ide := m.mcpIDEs[m.selected]
				if !ide.detected {
					m.errorMsg = ide.name + " not detected on this system"
					return m, nil
				}
				hopsuleBin, err := findHopsuleBinaryPath()
				if err != nil {
					m.errorMsg = "Could not find hopsule binary"
					return m, nil
				}
				if err := mcpInstallForIDE(ide.key, hopsuleBin); err != nil {
					m.errorMsg = "Install failed: " + err.Error()
					return m, nil
				}
				m.mcpIDEs = detectMCPStatus()
				m.errorMsg = ""
			}
			return m, nil
		}
		if m.currentView == viewDecisions {
			d := m.getSelectedDecision()
			if d != nil {
				m.detailView = true
				m.detailScroll = 0
				return m, nil
			}
		}
		if m.currentView == viewMemories {
			mem := m.getSelectedMemory()
			if mem != nil {
				m.detailView = true
				m.detailScroll = 0
				return m, nil
			}
		}
		if m.currentView == viewCapsules {
			c := m.getSelectedCapsule()
			if c != nil {
				m.detailView = true
				m.detailScroll = 0
				return m, nil
			}
		}
		return m.handleSelect()
	}

	return m, nil
}

// ============================================================================
// HOPPER CHAT KEY HANDLER
// ============================================================================

func (m model) handleHopperKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if m.hopperSidebarFocus {
		switch key {
		case "esc", "q":
			m.currentView = viewProjectMenu
			m.selected = m.menuSelectedIdx
			return m, nil
		case "tab":
			m.hopperSidebarFocus = false
			return m, nil
		case "j", "down":
			if m.chatSessionIdx < len(m.chatSessions)-1 {
				m.chatSessionIdx++
			}
			return m, nil
		case "k", "up":
			if m.chatSessionIdx > 0 {
				m.chatSessionIdx--
			}
			return m, nil
		case "enter":
			if len(m.chatSessions) > 0 && m.chatSessionIdx < len(m.chatSessions) {
				sess := m.chatSessions[m.chatSessionIdx]
				if len(m.chatMessages) > 0 {
					m.hopperSaveCurrentChat()
				}
				m.hopperSessionID = sess.ID
				m.hopperSidebarFocus = false
				if len(sess.Messages) > 0 {
					m.chatMessages = sess.Messages
					m.chatScroll = max(0, len(m.chatMessages)-6)
					return m, nil
				}
				m.chatMessages = nil
				m.chatScroll = 0
				return m, m.loadChatHistoryFromAPI(sess.ID)
			}
			return m, nil
		case "ctrl+n":
			m.hopperNewChat()
			m.hopperSidebarFocus = false
			return m, nil
		case "d":
			if len(m.chatSessions) > 0 && m.chatSessionIdx < len(m.chatSessions) {
				m.chatSessions = append(m.chatSessions[:m.chatSessionIdx], m.chatSessions[m.chatSessionIdx+1:]...)
				if m.chatSessionIdx >= len(m.chatSessions) && m.chatSessionIdx > 0 {
					m.chatSessionIdx--
				}
			}
			return m, nil
		}
		return m, nil
	}

	if m.chatStreaming {
		if key == "esc" {
			m.chatStreaming = false
			return m, nil
		}
		return m, nil
	}

	switch key {
	case "esc":
		m.currentView = viewProjectMenu
		m.selected = m.menuSelectedIdx
		m.chatInput = ""
		m.streamingContent = ""
		return m, nil
	case "tab":
		m.hopperSidebarFocus = true
		return m, nil
	case "ctrl+n":
		m.hopperNewChat()
		return m, nil
	case "ctrl+u":
		m.chatScroll = max(0, m.chatScroll-3)
		return m, nil
	case "ctrl+d":
		m.chatScroll = min(max(0, len(m.chatMessages)-1), m.chatScroll+3)
		return m, nil
	case "enter":
		if m.chatInput != "" {
			return m.sendHopperMessage()
		}
		return m, nil
	case "backspace", "delete":
		if len(m.chatInput) > 0 {
			m.chatInput = m.chatInput[:len(m.chatInput)-1]
		}
		return m, nil
	default:
		if msg.Type == tea.KeyRunes {
			m.chatInput += string(msg.Runes)
		} else if msg.Type == tea.KeySpace || key == " " {
			m.chatInput += " "
		}
		return m, nil
	}
}

// ============================================================================
// SUB-HANDLERS
// ============================================================================

func (m model) handleConfirmDelete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	switch key {
	case "y", "Y":
		m.confirmDelete = false
		switch m.confirmDeleteType {
		case "memory":
			err := m.client.DeleteMemory(m.currentProj.ID, m.confirmDeleteID)
			if err != nil {
				m.errorMsg = fmt.Sprintf("Failed to delete: %v", err)
			} else {
				m.loading = true
				m.confirmDeleteID = ""
				m.confirmDeleteType = ""
				if m.detailView {
					m.detailView = false
				}
				return m, m.loadMemories
			}
		case "task":
			err := m.client.DeleteTask(m.currentProj.ID, m.confirmDeleteID)
			if err != nil {
				m.errorMsg = fmt.Sprintf("Failed to delete: %v", err)
			} else {
				m.loading = true
				m.confirmDeleteID = ""
				m.confirmDeleteType = ""
				return m, m.loadTasks
			}
		case "decision":
			_, err := m.client.DeprecateDecision(m.currentProj.ID, m.confirmDeleteID)
			if err != nil {
				m.errorMsg = fmt.Sprintf("Failed to deprecate: %v", err)
			} else {
				m.loading = true
				m.confirmDeleteID = ""
				m.confirmDeleteType = ""
				if m.detailView {
					m.detailView = false
				}
				return m, m.loadDecisions
			}
		}
		m.confirmDeleteID = ""
		m.confirmDeleteType = ""
		return m, nil
	default:
		m.confirmDelete = false
		m.confirmDeleteID = ""
		m.confirmDeleteType = ""
		return m, nil
	}
}

func (m model) handleSearchMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.searchMode = false
		m.searchQuery = ""
		m.selected = 0
		m.scrollOffset = 0
		m.listPage = 0
		return m, nil
	case tea.KeyEnter:
		m.searchMode = false
		m.selected = 0
		m.scrollOffset = 0
		m.listPage = 0
		return m, nil
	case tea.KeyBackspace, tea.KeyDelete:
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			m.selected = 0
			m.scrollOffset = 0
			m.listPage = 0
		}
		return m, nil
	case tea.KeySpace:
		m.searchQuery += " "
		m.selected = 0
		m.scrollOffset = 0
		m.listPage = 0
		return m, nil
	case tea.KeyRunes:
		m.searchQuery += string(msg.Runes)
		m.selected = 0
		m.scrollOffset = 0
		m.listPage = 0
		return m, nil
	default:
		return m, nil
	}
}

func (m model) handleDetailView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.detailView = false
		m.detailScroll = 0
		return m, nil
	case "j", "down":
		m.detailScroll++
		return m, nil
	case "k", "up":
		if m.detailScroll > 0 {
			m.detailScroll--
		}
		return m, nil
	case "a":
		if m.currentView == viewDecisions {
			d := m.getSelectedDecision()
			if d != nil && (d.Status == "DRAFT" || d.Status == "PENDING") {
				userName := "CLI User"
				if m.cfg != nil && m.cfg.User != nil && m.cfg.User.Name != "" {
					userName = m.cfg.User.Name
				}
				_, err := m.client.AcceptDecision(m.currentProj.ID, d.ID, userName)
				if err != nil {
					m.errorMsg = fmt.Sprintf("Failed to accept: %v", err)
				} else {
					m.detailView = false
					m.loading = true
					return m, m.loadDecisions
				}
			}
		}
	case "x":
		if m.currentView == viewDecisions {
			d := m.getSelectedDecision()
			if d != nil && d.Status == "ACCEPTED" {
				_, err := m.client.DeprecateDecision(m.currentProj.ID, d.ID)
				if err != nil {
					m.errorMsg = fmt.Sprintf("Failed to deprecate: %v", err)
				} else {
					m.detailView = false
					m.loading = true
					return m, m.loadDecisions
				}
			}
		}
	case "e":
		if m.currentView == viewMemories {
			mem := m.getSelectedMemory()
			if mem != nil {
				m.detailView = false
				m.inputMode = true
				m.inputPrompt = "Edit Memory:"
				m.inputAction = "edit_memory"
				m.inputValue = mem.Content
				m.editingItemID = mem.ID
				return m, nil
			}
		}
	case "d":
		if m.currentView == viewDecisions {
			d := m.getSelectedDecision()
			if d != nil {
				m.confirmDelete = true
				m.confirmDeleteID = d.ID
				m.confirmDeleteType = "decision"
				return m, nil
			}
		} else if m.currentView == viewMemories {
			mem := m.getSelectedMemory()
			if mem != nil {
				m.confirmDelete = true
				m.confirmDeleteID = mem.ID
				m.confirmDeleteType = "memory"
				return m, nil
			}
		}
	}
	return m, nil
}

func (m model) handleInputMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.inputMode = false
		m.inputValue = ""
		m.inputAction = ""
		m.inputPrompt = ""
		m.editingItemID = ""
		return m, nil
	case tea.KeyEnter:
		if m.inputValue != "" {
			return m.executeInputAction()
		}
		return m, nil
	case tea.KeyBackspace, tea.KeyDelete:
		if len(m.inputValue) > 0 {
			m.inputValue = m.inputValue[:len(m.inputValue)-1]
		}
		return m, nil
	case tea.KeySpace:
		m.inputValue += " "
		return m, nil
	case tea.KeyRunes:
		m.inputValue += string(msg.Runes)
		return m, nil
	default:
		return m, nil
	}
}

func (m model) executeInputAction() (tea.Model, tea.Cmd) {
	m.inputMode = false
	switch m.inputAction {
	case "create_task":
		_, err := m.client.CreateTask(m.currentProj.ID, api.CreateTaskRequest{Title: m.inputValue})
		if err != nil {
			m.errorMsg = fmt.Sprintf("Failed to create task: %v", err)
		} else {
			m.loading = true
			m.inputValue = ""
			m.inputAction = ""
			return m, m.loadTasks
		}
	case "create_memory":
		_, err := m.client.CreateMemory(m.currentProj.ID, api.CreateMemoryRequest{Content: m.inputValue})
		if err != nil {
			m.errorMsg = fmt.Sprintf("Failed to create memory: %v", err)
		} else {
			m.loading = true
			m.inputValue = ""
			m.inputAction = ""
			return m, m.loadMemories
		}
	case "create_decision":
		_, err := m.client.CreateDecision(m.currentProj.ID, api.CreateDecisionRequest{Statement: m.inputValue})
		if err != nil {
			m.errorMsg = fmt.Sprintf("Failed to create decision: %v", err)
		} else {
			m.loading = true
			m.inputValue = ""
			m.inputAction = ""
			return m, m.loadDecisions
		}
	case "edit_task":
		if m.editingItemID != "" {
			_, err := m.client.UpdateTask(m.currentProj.ID, m.editingItemID, api.UpdateTaskRequest{Title: m.inputValue})
			if err != nil {
				m.errorMsg = fmt.Sprintf("Failed to update task: %v", err)
			} else {
				m.loading = true
				m.inputValue = ""
				m.inputAction = ""
				m.editingItemID = ""
				return m, m.loadTasks
			}
		}
	case "edit_memory":
		if m.editingItemID != "" {
			_, err := m.client.UpdateMemory(m.currentProj.ID, m.editingItemID, api.UpdateMemoryRequest{Content: m.inputValue})
			if err != nil {
				m.errorMsg = fmt.Sprintf("Failed to update memory: %v", err)
			} else {
				m.loading = true
				m.inputValue = ""
				m.inputAction = ""
				m.editingItemID = ""
				return m, m.loadMemories
			}
		}
	}
	m.inputValue = ""
	m.inputAction = ""
	m.editingItemID = ""
	return m, nil
}

func (m model) handleSelect() (tea.Model, tea.Cmd) {
	switch m.currentView {
	case viewLogin:
		if !m.loginPolling {
			m.loading = true
			return m, m.initLogin
		}
		return m, nil
	case viewOrganizations:
		if m.selected < len(m.organizations) {
			m.orgSelectedIdx = m.selected
			m.currentOrg = m.organizations[m.selected]
			m.currentView = viewProjects
			m.selected = 0
		} else if m.selected == len(m.organizations) {
			m.executeCmd = "logout"
			return m, tea.Quit
		}
	case viewProjects:
		orgProjects := m.getOrgProjects()
		if m.selected < len(orgProjects) {
			selectedProj := orgProjects[m.selected]

			if selectedProj.ID == m.initializedProjectID {
				// Already initialized — proceed
			} else if m.cwdInGitRepo {
				// In a git repo — init this project here
				m.initProjectLocally(selectedProj, m.currentOrg)
			} else {
				// Not in a git repo — block
				m.errorMsg = "Run inside a git repository: cd <project-dir> && hopsule"
				return m, nil
			}

			m.projSelectedIdx = m.selected
			m.currentProj = selectedProj
			m.menuItems = buildProjectMenuItems()
			m.currentView = viewProjectMenu
			m.selected = 0
		}
	case viewProjectMenu:
		if m.selected < len(m.menuItems) {
			m.menuSelectedIdx = m.selected
			item := m.menuItems[m.selected]
			switch item.action {
			case "back":
				m.currentView = viewProjects
				m.currentProj = nil
				m.selected = m.projSelectedIdx
			case "dashboard":
				m.currentView = viewDashboard
				m.selected = 0
				m.loading = true
				return m, m.loadDashboardData
			case "decisions":
				m.currentView = viewDecisions
				m.selected = 0
				m.scrollOffset = 0
				m.searchQuery = ""
				m.statusFilter = ""
				m.listPage = 0
				m.loading = true
				return m, m.loadDecisions
			case "memories":
				m.currentView = viewMemories
				m.selected = 0
				m.scrollOffset = 0
				m.searchQuery = ""
				m.listPage = 0
				m.loading = true
				return m, m.loadMemories
			case "capsules":
				m.currentView = viewCapsules
				m.selected = 0
				m.scrollOffset = 0
				m.searchQuery = ""
				m.listPage = 0
				m.loading = true
				return m, m.loadCapsules
			case "tasks":
				m.currentView = viewTasks
				m.selected = 0
				m.scrollOffset = 0
				m.searchQuery = ""
				m.loading = true
				return m, m.loadTasks
			case "hopper":
				m.currentView = viewHopper
				m.selected = 0
				m.hopperSidebarFocus = false
				m.chatInput = ""
				m.streamingContent = ""
				if m.hopperSessionID == "" {
					m.hopperSessionID = fmt.Sprintf("hopper-%d-%d", time.Now().UnixNano(), time.Now().UnixMilli()%1000)
				}
				if !m.hopperContextLoaded {
					m.loading = true
					return m, m.loadHopperContext
				}
				return m, nil
			case "mcp":
				m.currentView = viewMCP
				m.selected = 0
				m.mcpIDEs = detectMCPStatus()
				return m, nil
			}
		}
	}
	return m, nil
}

// ============================================================================
// LOGIN HELPERS
// ============================================================================

func (m model) initLogin() tea.Msg {
	cfg := m.cfg
	if cfg == nil {
		cfg = &config.Config{}
	}
	if cfg.APIURL == "" {
		cfg.APIURL = "https://api.hopsule.com"
	}
	if cfg.WebURL == "" {
		cfg.WebURL = "https://app.hopsule.com"
	}
	client := api.NewClient(cfg)
	initResp, err := client.DeviceAuthInit("CLI")
	if err != nil {
		return loginInitMsg{err: err}
	}
	authURL := fmt.Sprintf("%s/auth/device?code=%s", cfg.WebURL, initResp.Code)
	_ = openBrowser(authURL)
	return loginInitMsg{deviceCode: initResp.Code, authURL: authURL}
}

func (m model) pollLogin() tea.Msg {
	cfg := m.cfg
	if cfg == nil {
		cfg = &config.Config{}
	}
	if cfg.APIURL == "" {
		cfg.APIURL = "https://api.hopsule.com"
	}
	client := api.NewClient(cfg)
	resp, err := client.DeviceAuthPoll(m.loginDeviceCode)
	if err != nil {
		return loginPollMsg{err: err}
	}
	switch resp.Status {
	case "complete":
		return loginPollMsg{
			status: "complete",
			token:  resp.Token,
			user: &config.User{
				ID: resp.UserID, Email: resp.Email,
				Name: resp.Name, AvatarURL: resp.AvatarURL,
			},
		}
	case "expired":
		return loginPollMsg{status: "expired", err: fmt.Errorf("login session expired")}
	default:
		return loginPollMsg{status: "pending"}
	}
}

func loginTickCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg { return loginTickMsg{} })
}

// getFilteredCount returns the count of filtered items for the current view.
func (m model) getFilteredCount() int {
	switch m.currentView {
	case viewDecisions:
		return len(m.getFilteredDecisions())
	case viewMemories:
		return len(m.getFilteredMemories())
	case viewCapsules:
		return len(m.getFilteredCapsules())
	case viewTasks:
		return len(m.getFilteredTasks())
	}
	return 0
}

// initProjectLocally saves the .hopsule file and updates .gitignore for the selected project.
func (m *model) initProjectLocally(proj *api.Project, org *api.Organization) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	projectCfg := &config.ProjectConfig{
		Version: config.HopsuleFileVersion,
		Project: config.ProjectInfo{
			ID:   proj.ID,
			Slug: proj.Slug,
			Name: proj.Name,
			Organization: config.OrganizationInfo{
				ID:   org.ID,
				Slug: org.Slug,
				Name: org.Name,
			},
		},
	}

	if err := config.SaveProjectConfig(cwd, projectCfg); err != nil {
		return
	}

	m.initializedProjectID = proj.ID

	// Add .hopsule to .gitignore if inside a git repo
	if isGitRepo(cwd) {
		if gitRoot := gitRootDir(cwd); gitRoot != "" {
			addHopsuleToGitignore(gitRoot)
		}
	}

	// Update global config
	if m.cfg != nil {
		m.cfg.Project = proj.ID
		m.cfg.Organization = org.ID
		config.SaveConfig(m.cfg)
	}
}

func isGitRepo(dir string) bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}

func gitRootDir(dir string) string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func addHopsuleToGitignore(gitRoot string) {
	gitignorePath := filepath.Join(gitRoot, ".gitignore")

	content, err := os.ReadFile(gitignorePath)
	if err == nil {
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == ".hopsule" || trimmed == ".hopsule/" {
				return
			}
		}
	}

	f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	if len(content) > 0 && content[len(content)-1] != '\n' {
		f.WriteString("\n")
	}
	f.WriteString("\n# Hopsule project config (local, not shared)\n.hopsule\n")
}

// ============================================================================
// PROJECT MENU BUILDER
// ============================================================================

func buildProjectMenuItems() []menuItem {
	return []menuItem{
		{"●", "Dashboard", "Project overview & stats", "dashboard"},
		{"●", "Decisions", "View & manage decisions", "decisions"},
		{"●", "Memories", "Project memories & context", "memories"},
		{"●", "Capsules", "Context packs", "capsules"},
		{"●", "Tasks", "Task management", "tasks"},
		{"◆", "MCP", "AI tool connections", "mcp"},
		{"●", "Hopper", "AI Assistant", "hopper"},
		{"", "", "", ""},
		{"<", "Back", "Return to projects", "back"},
	}
}

// ============================================================================
// MCP IDE DETECTION & INSTALL (used by TUI)
// ============================================================================

func detectMCPStatus() []mcpIDEStatus {
	ides := []mcpIDEStatus{
		{name: "Cursor", key: "cursor"},
		{name: "Claude Desktop", key: "claude-desktop"},
		{name: "Claude Code", key: "claude-code"},
	}

	for i := range ides {
		switch ides[i].key {
		case "cursor":
			if _, err := os.Stat(".cursor"); err == nil {
				ides[i].detected = true
				ides[i].configured = isCursorMCPConfigured()
			}
		case "claude-desktop":
			configDir := claudeDesktopConfigDir()
			if configDir != "" {
				if _, err := os.Stat(configDir); err == nil {
					ides[i].detected = true
					ides[i].configured = isClaudeDesktopMCPConfigured(configDir)
				}
			}
		case "claude-code":
			if _, err := exec.LookPath("claude"); err == nil {
				ides[i].detected = true
				ides[i].configured = isClaudeCodeMCPConfigured()
			}
		}
	}

	return ides
}

func isCursorMCPConfigured() bool {
	data, err := os.ReadFile(filepath.Join(".cursor", "mcp.json"))
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "hopsule")
}

func claudeDesktopConfigDir() string {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Claude")
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "Claude")
	case "linux":
		return filepath.Join(home, ".config", "Claude")
	}
	return ""
}

func isClaudeDesktopMCPConfigured(configDir string) bool {
	data, err := os.ReadFile(filepath.Join(configDir, "claude_desktop_config.json"))
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "hopsule")
}

func isClaudeCodeMCPConfigured() bool {
	cmd := exec.Command("claude", "mcp", "list")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "hopsule")
}

func findHopsuleBinaryPath() (string, error) {
	execPath, err := os.Executable()
	if err == nil {
		return execPath, nil
	}
	path, err := exec.LookPath("hopsule")
	if err == nil {
		return path, nil
	}
	return "", fmt.Errorf("hopsule binary not found")
}

func mcpInstallForIDE(ide, hopsuleBin string) error {
	switch ide {
	case "cursor":
		return mcpInstallCursor(hopsuleBin)
	case "claude-desktop":
		return mcpInstallClaudeDesktop(hopsuleBin)
	case "claude-code":
		return mcpInstallClaudeCode(hopsuleBin)
	}
	return fmt.Errorf("unsupported IDE: %s", ide)
}

func mcpInstallCursor(hopsuleBin string) error {
	configDir := ".cursor"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	configPath := filepath.Join(configDir, "mcp.json")

	cwd, _ := os.Getwd()

	var config map[string]interface{}
	if data, err := os.ReadFile(configPath); err == nil {
		json.Unmarshal(data, &config)
	}
	if config == nil {
		config = make(map[string]interface{})
	}
	servers, ok := config["mcpServers"].(map[string]interface{})
	if !ok {
		servers = make(map[string]interface{})
	}
	serverCfg := map[string]interface{}{
		"command": hopsuleBin,
		"args":    []string{"mcp", "serve"},
	}
	if cwd != "" {
		serverCfg["cwd"] = cwd
	}
	servers["hopsule"] = serverCfg
	config["mcpServers"] = servers

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func mcpInstallClaudeDesktop(hopsuleBin string) error {
	configDir := claudeDesktopConfigDir()
	if configDir == "" {
		return fmt.Errorf("unsupported OS")
	}
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	configPath := filepath.Join(configDir, "claude_desktop_config.json")

	var config map[string]interface{}
	if data, err := os.ReadFile(configPath); err == nil {
		json.Unmarshal(data, &config)
	}
	if config == nil {
		config = make(map[string]interface{})
	}
	servers, ok := config["mcpServers"].(map[string]interface{})
	if !ok {
		servers = make(map[string]interface{})
	}
	servers["hopsule"] = map[string]interface{}{
		"command": hopsuleBin,
		"args":    []string{"mcp", "serve"},
	}
	config["mcpServers"] = servers

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func mcpInstallClaudeCode(hopsuleBin string) error {
	claudePath, err := exec.LookPath("claude")
	if err != nil {
		return fmt.Errorf("'claude' command not found")
	}
	cmd := exec.Command(claudePath, "mcp", "add", "hopsule", "--", hopsuleBin, "mcp", "serve")
	return cmd.Run()
}
