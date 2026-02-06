package ui

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/Cagangedik/cli-tool/internal/api"
	"github.com/Cagangedik/cli-tool/internal/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ============================================================================
// STYLES
// ============================================================================

var (
	// Colors
	cyanColor    = lipgloss.Color("51")
	greenColor   = lipgloss.Color("46")
	yellowColor  = lipgloss.Color("226")
	magentaColor = lipgloss.Color("201")
	grayColor    = lipgloss.Color("244")
	whiteColor   = lipgloss.Color("255")
	dimColor     = lipgloss.Color("240")
	blueColor    = lipgloss.Color("39")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Foreground(whiteColor).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(greenColor).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(whiteColor)

	logoStyle = lipgloss.NewStyle().
			Foreground(cyanColor)

	dimStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	accentStyle = lipgloss.NewStyle().
			Foreground(yellowColor)

	statusOnStyle = lipgloss.NewStyle().
			Foreground(greenColor).
			Bold(true)

	// Card styles for kanban
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(grayColor).
			Padding(1, 2).
			Width(32).
			MarginRight(2)

	selectedCardStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(greenColor).
				Padding(1, 2).
				Width(32).
				MarginRight(2)

	cardTitleStyle = lipgloss.NewStyle().
			Foreground(whiteColor).
			Bold(true)

	cardDescStyle = lipgloss.NewStyle().
			Foreground(grayColor)

	breadcrumbStyle = lipgloss.NewStyle().
			Foreground(cyanColor)

	helpStyle = lipgloss.NewStyle().
			Foreground(dimColor)
)

// ============================================================================
// VIEW TYPES
// ============================================================================

type viewType int

const (
	viewLogin viewType = iota
	viewOrganizations
	viewProjects
	viewProjectMenu
	viewDashboard
	viewDecisions
	viewMemories
	viewCapsules
	viewTasks
	viewBrain
	viewHopper
)

// ============================================================================
// MODEL
// ============================================================================

type model struct {
	cfg           *config.Config
	client        *api.Client
	currentView   viewType
	
	// Data
	organizations []*api.Organization
	projects      []*api.Project
	currentOrg    *api.Organization
	currentProj   *api.Project
	
	// Project menu items
	menuItems     []menuItem
	
	// Feature data
	decisions     []api.Decision
	memories      []*api.Memory
	tasks         []*api.Task
	capsules      []*api.Capsule
	graphStats    *api.GraphStats
	
	// Hopper chat state
	chatMessages      []api.ChatMessage
	chatInput         string
	chatStreaming     bool
	streamingContent  string
	hopperDecisions   []api.Decision  // RAG context for Hopper
	hopperMemories    []*api.Memory   // RAG context for Hopper
	hopperContextLoaded bool
	hopperSessionID   string
	
	// Selection
	selected      int
	scrollOffset  int  // For list pagination
	
	// UI state
	width         int
	height        int
	loading       bool
	errorMsg      string
	
	// Actions
	executeCmd    string
}

type menuItem struct {
	icon        string
	name        string
	description string
	action      string
}

// ============================================================================
// MESSAGES
// ============================================================================

type dataLoadedMsg struct {
	organizations []*api.Organization
	projects      []*api.Project
	err           error
}

type loginCompleteMsg struct {
	success bool
	err     error
}

type decisionsLoadedMsg struct {
	decisions []api.Decision
	err       error
}

type memoriesLoadedMsg struct {
	memories []*api.Memory
	err      error
}

type tasksLoadedMsg struct {
	tasks []*api.Task
	err   error
}

type capsulesLoadedMsg struct {
	capsules []*api.Capsule
	err      error
}

type brainStatsLoadedMsg struct {
	stats *api.GraphStats
	err   error
}

type dashboardLoadedMsg struct {
	decisions []api.Decision
	memories  []*api.Memory
	tasks     []*api.Task
	capsules  []*api.Capsule
	err       error
}

// Hopper chat messages
type chatStreamChunkMsg struct {
	chunk string
}

type chatStreamDoneMsg struct {
	err error
}

type hopperContextLoadedMsg struct {
	decisions []api.Decision
	memories  []*api.Memory
	err       error
}

// ============================================================================
// INIT & UPDATE
// ============================================================================

func NewInteractiveModel(cfg *config.Config) model {
	isLoggedIn := cfg != nil && cfg.IsAuthenticated()
	
	m := model{
		cfg:      cfg,
		selected: 0,
		hopperSessionID: fmt.Sprintf("cli-%d", time.Now().UnixNano()),
	}
	
	if isLoggedIn {
		m.client = api.NewClient(cfg)
		m.currentView = viewOrganizations
		m.loading = true
	} else {
		m.currentView = viewLogin
	}
	
	return m
}

func (m model) Init() tea.Cmd {
	if m.loading {
		return m.loadData
	}
	return nil
}

func (m model) loadData() tea.Msg {
	if m.client == nil {
		return dataLoadedMsg{err: fmt.Errorf("not authenticated")}
	}
	
	meResp, err := m.client.GetMe()
	if err != nil {
		return dataLoadedMsg{err: err}
	}
	
	return dataLoadedMsg{
		organizations: meResp.Organizations,
		projects:      meResp.Projects,
	}
}

func (m model) loadDecisions() tea.Msg {
	if m.client == nil || m.currentProj == nil {
		return decisionsLoadedMsg{err: fmt.Errorf("not authenticated or no project selected")}
	}
	decisions, err := m.client.ListDecisions(m.currentProj.ID)
	if err != nil {
		return decisionsLoadedMsg{err: err}
	}
	return decisionsLoadedMsg{decisions: decisions}
}

func (m model) loadMemories() tea.Msg {
	if m.client == nil || m.currentProj == nil {
		return memoriesLoadedMsg{err: fmt.Errorf("not authenticated or no project selected")}
	}
	memories, err := m.client.ListMemories(m.currentProj.ID)
	if err != nil {
		return memoriesLoadedMsg{err: err}
	}
	return memoriesLoadedMsg{memories: memories}
}

func (m model) loadTasks() tea.Msg {
	if m.client == nil || m.currentProj == nil {
		return tasksLoadedMsg{err: fmt.Errorf("not authenticated or no project selected")}
	}
	tasks, err := m.client.ListTasks(m.currentProj.ID)
	if err != nil {
		return tasksLoadedMsg{err: err}
	}
	return tasksLoadedMsg{tasks: tasks}
}

func (m model) loadCapsules() tea.Msg {
	if m.client == nil || m.currentProj == nil {
		return capsulesLoadedMsg{err: fmt.Errorf("not authenticated or no project selected")}
	}
	capsules, err := m.client.ListCapsules(m.currentProj.ID)
	if err != nil {
		return capsulesLoadedMsg{err: err}
	}
	return capsulesLoadedMsg{capsules: capsules}
}

func (m model) loadBrainStats() tea.Msg {
	if m.client == nil || m.currentProj == nil {
		return brainStatsLoadedMsg{err: fmt.Errorf("not authenticated or no project selected")}
	}
	stats, err := m.client.GetGraphStats(m.currentProj.ID)
	if err != nil {
		return brainStatsLoadedMsg{err: err}
	}
	return brainStatsLoadedMsg{stats: stats}
}

// loadHopperContext loads decisions and memories for RAG context
func (m model) loadHopperContext() tea.Msg {
	if m.client == nil || m.currentProj == nil {
		return hopperContextLoadedMsg{err: fmt.Errorf("not authenticated or no project selected")}
	}
	
	// Load decisions
	decisions, err := m.client.ListDecisions(m.currentProj.ID)
	if err != nil {
		return hopperContextLoadedMsg{err: fmt.Errorf("failed to load decisions: %w", err)}
	}
	
	// Load memories
	memories, err := m.client.ListMemories(m.currentProj.ID)
	if err != nil {
		return hopperContextLoadedMsg{err: fmt.Errorf("failed to load memories: %w", err)}
	}
	
	return hopperContextLoadedMsg{
		decisions: decisions,
		memories:  memories,
	}
}

// sendHopperMessage sends a chat message to Hopper AI with RAG context
func (m model) sendHopperMessage() (tea.Model, tea.Cmd) {
	if m.client == nil || m.currentProj == nil {
		m.errorMsg = "Not authenticated or no project selected"
		return m, nil
	}
	
	userMessage := m.chatInput
	m.chatInput = ""
	m.chatStreaming = true
	m.streamingContent = ""
	
	// Add user message to history
	m.chatMessages = append(m.chatMessages, api.ChatMessage{
		Role:    "user",
		Content: userMessage,
	})
	
	// Build TaggedItems from decisions and memories for RAG context
	var taggedItems []api.TaggedItem
	
	// Add decisions as context (limit to 10 for speed)
	decisionLimit := 10
	if len(m.hopperDecisions) < decisionLimit {
		decisionLimit = len(m.hopperDecisions)
	}
	for i := 0; i < decisionLimit; i++ {
		d := m.hopperDecisions[i]
		taggedItems = append(taggedItems, api.TaggedItem{
			ID:        d.ID,
			Type:      "decision",
			Statement: truncateString(d.Statement, 160),
			Content:   truncateString(d.Rationale, 240),
		})
	}
	
	// Add memories as context (limit to 15 for speed)
	memoryLimit := 15
	if len(m.hopperMemories) < memoryLimit {
		memoryLimit = len(m.hopperMemories)
	}
	for i := 0; i < memoryLimit; i++ {
		mem := m.hopperMemories[i]
		taggedItems = append(taggedItems, api.TaggedItem{
			ID:      mem.ID,
			Type:    "memory",
			Content: truncateString(mem.Content, 240),
		})
	}
	
	// Create async command that collects the full response
	client := m.client
	projectID := m.currentProj.ID
	projectName := m.currentProj.Name
	history := make([]api.ChatMessage, len(m.chatMessages)-1)
	copy(history, m.chatMessages[:len(m.chatMessages)-1])
	
	return m, func() tea.Msg {
		var fullResponse string
		
		req := &api.ChatRequest{
			Message:             userMessage,
			ConversationHistory: history,
			TaggedItems:         taggedItems,
			Stream:              true,
			SessionID:           m.hopperSessionID,
			ProjectName:         projectName,
		}
		
		err := client.SendChatMessage(projectID, req, func(chunk string) {
			fullResponse += chunk
		})
		
		if err != nil {
			return chatStreamDoneMsg{err: err}
		}
		
		// Return the full response as a chunk message first
		return chatStreamChunkMsg{chunk: fullResponse}
	}
}

func (m model) loadDashboardData() tea.Msg {
	if m.client == nil || m.currentProj == nil {
		return dashboardLoadedMsg{err: fmt.Errorf("not authenticated or no project selected")}
	}
	
	decisions, _ := m.client.ListDecisions(m.currentProj.ID)
	memories, _ := m.client.ListMemories(m.currentProj.ID)
	tasks, _ := m.client.ListTasks(m.currentProj.ID)
	capsules, _ := m.client.ListCapsules(m.currentProj.ID)
	
	return dashboardLoadedMsg{
		decisions: decisions,
		memories:  memories,
		tasks:     tasks,
		capsules:  capsules,
	}
}

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
		}
		return m, nil
		
	case loginCompleteMsg:
		if msg.success {
			// Reload config and data
			m.cfg, _ = config.GetConfig()
			m.client = api.NewClient(m.cfg)
			m.currentView = viewOrganizations
			m.loading = true
			m.selected = 0
			return m, m.loadData
		} else if msg.err != nil {
			m.errorMsg = msg.err.Error()
		}
		return m, nil
		
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
		// Add assistant response to history
		if msg.chunk != "" {
			m.chatMessages = append(m.chatMessages, api.ChatMessage{
				Role:    "assistant",
				Content: sanitizeMarkdownForTerminal(msg.chunk),
			})
		}
		return m, nil
		
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
		}
		return m, nil
	}
	
	return m, nil
}

func (m model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Clear error on any key press
	m.errorMsg = ""
	
	// Handle Hopper chat input specially
	if m.currentView == viewHopper && !m.chatStreaming {
		key := msg.String()
		switch key {
		case "esc":
			// ESC exits Hopper
			m.currentView = viewProjectMenu
			m.selected = 0
			m.chatMessages = nil
			m.chatInput = ""
			m.streamingContent = ""
			return m, nil
		case "ctrl+c":
			// Ctrl+C also exits
			m.currentView = viewProjectMenu
			m.selected = 0
			return m, nil
		case "enter":
			// Send message
			if m.chatInput != "" {
				return m.sendHopperMessage()
			}
			return m, nil
		case "backspace":
			if len(m.chatInput) > 0 {
				m.chatInput = m.chatInput[:len(m.chatInput)-1]
			}
			return m, nil
		default:
			// Add character to input
			if len(key) == 1 {
				m.chatInput += key
			} else if key == "space" {
				m.chatInput += " "
			}
			return m, nil
		}
	}
	
	switch msg.String() {
	case "ctrl+c", "q":
		// Feature views go back to project menu
		if m.currentView == viewDashboard || m.currentView == viewDecisions ||
			m.currentView == viewMemories || m.currentView == viewCapsules ||
			m.currentView == viewTasks || m.currentView == viewBrain {
			m.currentView = viewProjectMenu
			m.selected = 0
				return m, nil
			}
		if m.currentView == viewProjectMenu {
			m.currentView = viewProjects
			m.currentProj = nil
			m.selected = 0
			return m, nil
		}
		if m.currentView == viewProjects {
			m.currentView = viewOrganizations
			m.currentOrg = nil
			m.selected = 0
			return m, nil
		}
			return m, tea.Quit

	case "esc":
		// Feature views go back to project menu
		if m.currentView == viewDashboard || m.currentView == viewDecisions ||
			m.currentView == viewMemories || m.currentView == viewCapsules ||
			m.currentView == viewTasks || m.currentView == viewBrain {
			m.currentView = viewProjectMenu
			m.selected = 0
			return m, nil
		}
		if m.currentView == viewProjectMenu {
			m.currentView = viewProjects
			m.currentProj = nil
			m.selected = 0
			return m, nil
		}
		if m.currentView == viewProjects {
			m.currentView = viewOrganizations
			m.currentOrg = nil
			m.selected = 0
			return m, nil
		}
			return m, nil

		case "up", "k":
			if m.selected > 0 {
				m.selected--
			// Scroll up if selection goes above visible area
			if m.selected < m.scrollOffset {
				m.scrollOffset = m.selected
			}
			}

		case "down", "j":
		maxSel := m.getMaxSelection() - 1
		if m.selected < maxSel {
				m.selected++
			// Scroll down if selection goes below visible area (10 items visible)
			visibleItems := 10
			if m.selected >= m.scrollOffset+visibleItems {
				m.scrollOffset = m.selected - visibleItems + 1
			}
		}
		
	case "left", "h":
		if m.currentView == viewProjects && m.selected > 0 {
			m.selected--
		}
		
	case "right", "l":
		if m.currentView == viewProjects {
			m.selected = min(m.selected+1, m.getMaxSelection()-1)
		}
	
	// CRUD shortcuts for Decisions
	case "n":
		if m.currentView == viewDecisions {
			// TODO: Open create decision dialog
			m.errorMsg = "Create decision: Coming soon! Use web app for now."
		} else if m.currentView == viewMemories {
			// TODO: Open create memory dialog
			m.errorMsg = "Create memory: Coming soon! Use web app for now."
		} else if m.currentView == viewTasks {
			// TODO: Open create task dialog
			m.errorMsg = "Create task: Coming soon! Use web app for now."
		}
		
	case "a":
		if m.currentView == viewDecisions && len(m.decisions) > 0 && m.selected < len(m.decisions) {
			d := m.decisions[m.selected]
			if d.Status == "DRAFT" || d.Status == "PENDING" {
				// Accept decision
				_, err := m.client.AcceptDecision(m.currentProj.ID, d.ID)
				if err != nil {
					m.errorMsg = fmt.Sprintf("Failed to accept: %v", err)
				} else {
					m.loading = true
					return m, m.loadDecisions
				}
			} else {
				m.errorMsg = "Can only accept DRAFT or PENDING decisions"
			}
		}
		
	case "x":
		if m.currentView == viewDecisions && len(m.decisions) > 0 && m.selected < len(m.decisions) {
			d := m.decisions[m.selected]
			if d.Status == "ACCEPTED" {
				// Deprecate decision
				_, err := m.client.DeprecateDecision(m.currentProj.ID, d.ID)
				if err != nil {
					m.errorMsg = fmt.Sprintf("Failed to deprecate: %v", err)
				} else {
					m.loading = true
					return m, m.loadDecisions
				}
			} else {
				m.errorMsg = "Can only deprecate ACCEPTED decisions"
			}
		}
		
	case "d":
		if m.currentView == viewMemories && len(m.memories) > 0 && m.selected < len(m.memories) {
			mem := m.memories[m.selected]
			err := m.client.DeleteMemory(m.currentProj.ID, mem.ID)
			if err != nil {
				m.errorMsg = fmt.Sprintf("Failed to delete: %v", err)
			} else {
				m.loading = true
				return m, m.loadMemories
			}
		} else if m.currentView == viewTasks && len(m.tasks) > 0 && m.selected < len(m.tasks) {
			task := m.tasks[m.selected]
			err := m.client.DeleteTask(m.currentProj.ID, task.ID)
			if err != nil {
				m.errorMsg = fmt.Sprintf("Failed to delete: %v", err)
			} else {
				m.loading = true
				return m, m.loadTasks
			}
		}
		
	case "t":
		if m.currentView == viewTasks && len(m.tasks) > 0 && m.selected < len(m.tasks) {
			task := m.tasks[m.selected]
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
		
	case "enter", " ":
		return m.handleSelect()
	}
	
	return m, nil
}

func (m model) handleSelect() (tea.Model, tea.Cmd) {
	switch m.currentView {
	case viewLogin:
		// Start login flow
		m.executeCmd = "login"
		return m, tea.Quit
		
	case viewOrganizations:
		if m.selected < len(m.organizations) {
			m.currentOrg = m.organizations[m.selected]
			m.currentView = viewProjects
			m.selected = 0
		} else if m.selected == len(m.organizations) {
			// Logout option
			m.executeCmd = "logout"
			return m, tea.Quit
		}

	case viewProjects:
		orgProjects := m.getOrgProjects()
		if m.selected < len(orgProjects) {
			// Open project menu
			m.currentProj = orgProjects[m.selected]
			m.menuItems = []menuItem{
				{"📊", "Dashboard", "Project overview & stats", "dashboard"},
				{"📋", "Decisions", "View & manage decisions", "decisions"},
				{"💾", "Memories", "Project memories & context", "memories"},
				{"📦", "Capsules", "Context packs", "capsules"},
				{"✅", "Tasks", "Task management", "tasks"},
				{"🧠", "Brain", "Knowledge graph", "brain"},
				{"🤖", "Hopper", "AI Assistant", "hopper"},
				{"", "", "", ""},
				{"🔙", "Back", "Return to projects", "back"},
			}
			m.currentView = viewProjectMenu
			m.selected = 0
		}
		
	case viewProjectMenu:
		if m.selected < len(m.menuItems) {
			item := m.menuItems[m.selected]
			switch item.action {
			case "back":
				m.currentView = viewProjects
				m.currentProj = nil
				m.selected = 0
			case "dashboard":
				m.currentView = viewDashboard
				m.selected = 0
				m.loading = true
				return m, m.loadDashboardData
			case "decisions":
				m.currentView = viewDecisions
				m.selected = 0
				m.scrollOffset = 0
				m.loading = true
				return m, m.loadDecisions
			case "memories":
				m.currentView = viewMemories
				m.selected = 0
				m.scrollOffset = 0
				m.loading = true
				return m, m.loadMemories
			case "capsules":
				m.currentView = viewCapsules
				m.selected = 0
				m.scrollOffset = 0
				m.loading = true
				return m, m.loadCapsules
			case "tasks":
				m.currentView = viewTasks
				m.selected = 0
				m.scrollOffset = 0
				m.loading = true
				return m, m.loadTasks
			case "brain":
				m.currentView = viewBrain
				m.selected = 0
				m.loading = true
				return m, m.loadBrainStats
			case "hopper":
				m.currentView = viewHopper
				m.selected = 0
				m.chatMessages = nil
				m.chatInput = ""
				m.streamingContent = ""
				m.loading = true
				// Load decisions and memories for RAG context
				return m, m.loadHopperContext
			}
		}
	}

	return m, nil
}

func (m model) getMaxSelection() int {
	switch m.currentView {
	case viewLogin:
		return 1
	case viewOrganizations:
		return len(m.organizations) + 1 // +1 for logout
	case viewProjects:
		return len(m.getOrgProjects())
	case viewProjectMenu:
		return len(m.menuItems)
	case viewDecisions:
		return len(m.decisions)
	case viewMemories:
		return len(m.memories)
	case viewTasks:
		return len(m.tasks)
	case viewCapsules:
		return len(m.capsules)
	case viewDashboard, viewBrain:
		return 0 // Read-only views
	}
	return 0
}

func (m model) getOrgProjects() []*api.Project {
	if m.currentOrg == nil {
		return nil
	}
	var projects []*api.Project
	for _, p := range m.projects {
		if p.OrganizationID == m.currentOrg.ID {
			projects = append(projects, p)
		}
	}
	return projects
}

// ============================================================================
// VIEW
// ============================================================================

func (m model) View() string {
	var s string
	
	// Header
	s += m.renderHeader()
	s += "\n"
	
	// Loading state
	if m.loading {
		s += "\n  " + dimStyle.Render("Loading...") + "\n"
		return s
	}
	
	// Error state
	if m.errorMsg != "" {
		s += "\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("Error: "+m.errorMsg) + "\n"
	}
	
	// Main content
	switch m.currentView {
	case viewLogin:
		s += m.renderLoginView()
	case viewOrganizations:
		s += m.renderOrganizationsView()
	case viewProjects:
		s += m.renderProjectsView()
	case viewProjectMenu:
		s += m.renderProjectMenuView()
	case viewDashboard:
		s += m.renderDashboardView()
	case viewDecisions:
		s += m.renderDecisionsView()
	case viewMemories:
		s += m.renderMemoriesView()
	case viewCapsules:
		s += m.renderCapsulesView()
	case viewTasks:
		s += m.renderTasksView()
	case viewBrain:
		s += m.renderBrainView()
	case viewHopper:
		s += m.renderHopperView()
	}
	
	// Footer
	s += "\n"
	s += m.renderFooter()
	
	return s
}

func (m model) renderHeader() string {
	var s string
	
	s += "\n"
	
	// Mini logo + title
	logo := logoStyle.Render("◆")
	title := titleStyle.Render("Hopsule")
	version := dimStyle.Render("v0.7.1")
	
	s += fmt.Sprintf("  %s %s %s", logo, title, version)
	
	// User info if logged in
	if m.cfg != nil && m.cfg.User != nil {
		s += "  " + dimStyle.Render("•") + "  "
		s += statusOnStyle.Render("●") + " "
		s += normalStyle.Render(m.cfg.User.Name)
	}
	
	s += "\n"
	
	// Breadcrumb
	if m.currentView == viewProjects && m.currentOrg != nil {
		s += "  " + breadcrumbStyle.Render("Organizations") + dimStyle.Render(" / ") + accentStyle.Render(m.currentOrg.Name)
		s += "\n"
	}
	
	return s
}

func (m model) renderLoginView() string {
	var s string
	
	s += "\n"
	s += m.renderLogo()
	s += "\n\n"
	
	s += "  " + titleStyle.Render("Welcome to Hopsule") + "\n"
	s += "  " + dimStyle.Render("Decision & Memory Layer for AI teams") + "\n\n"
	
	// Login button
	if m.selected == 0 {
		s += "  " + selectedStyle.Render("> ") + accentStyle.Render("Login with Browser") + "\n"
	} else {
		s += "    " + normalStyle.Render("Login with Browser") + "\n"
	}
	
	s += "\n"
	s += "  " + dimStyle.Render("Press Enter to open browser and sign in") + "\n"
	
	return s
}

func (m model) renderOrganizationsView() string {
	var s string
	
	s += "\n"
	s += "  " + titleStyle.Render("Organizations") + "\n"
	s += "  " + dimStyle.Render("Select an organization to view projects") + "\n\n"
	
	if len(m.organizations) == 0 {
		s += "  " + dimStyle.Render("No organizations found.") + "\n"
		s += "  " + dimStyle.Render("Create one at https://hopsule.com") + "\n"
	} else {
		for i, org := range m.organizations {
		if i == m.selected {
				s += "  " + selectedStyle.Render("> ")
				s += accentStyle.Render(org.Name)
				s += " " + dimStyle.Render("@"+org.Slug)
		} else {
				s += "    "
				s += normalStyle.Render(org.Name)
				s += " " + dimStyle.Render("@"+org.Slug)
			}
			s += "\n"
		}
		
		s += "\n"
		
		// Logout option
		if m.selected == len(m.organizations) {
			s += "  " + selectedStyle.Render("> ") + dimStyle.Render("Logout")
		} else {
			s += "    " + dimStyle.Render("Logout")
		}
		s += "\n"
	}
	
	return s
}

func (m model) renderProjectsView() string {
	var s string
	
	s += "\n"
	s += "  " + titleStyle.Render("Projects") + "\n"
	s += "  " + dimStyle.Render("Select a project to work with") + "\n\n"
	
	projects := m.getOrgProjects()
	
	if len(projects) == 0 {
		s += "  " + dimStyle.Render("No projects in this organization.") + "\n"
		s += "  " + dimStyle.Render("Create one at https://hopsule.com") + "\n"
		return s
	}
	
	// Render projects as cards in a grid
	cardsPerRow := 2
	if m.width > 120 {
		cardsPerRow = 3
	}
	
	for i := 0; i < len(projects); i += cardsPerRow {
		var rowCards []string
		
		for j := 0; j < cardsPerRow && i+j < len(projects); j++ {
			idx := i + j
			proj := projects[idx]
			card := m.renderProjectCard(proj, idx == m.selected)
			rowCards = append(rowCards, card)
		}
		
		row := lipgloss.JoinHorizontal(lipgloss.Top, rowCards...)
		s += lipgloss.NewStyle().MarginLeft(2).Render(row) + "\n"
	}
	
	return s
}

func (m model) renderProjectCard(proj *api.Project, selected bool) string {
	style := cardStyle
	if selected {
		style = selectedCardStyle
	}
	
	title := cardTitleStyle.Render(proj.Name)
	
	desc := proj.Description
	if desc == "" {
		desc = "No description"
	}
	if len(desc) > 40 {
		desc = desc[:37] + "..."
	}
	description := cardDescStyle.Render(desc)
	
	// Card content
	content := title + "\n\n" + description
	
	if selected {
		content += "\n\n" + selectedStyle.Render("▶ Select")
	}
	
	return style.Render(content)
}

func (m model) renderLogo() string {
	return logoStyle.Render(`                      ▟███████▙      ▟███████▙
                      █████████      █████████
                      █████████      █████████
                      █████████      █████████
                      ▝▀▀▀▀▀▀▀▀▄▄▄▄▄▄▛▀▀▀▀▀▀▀▘
                           ██████████████▄
                           ██████████████████▖
                      ▄▄▄▄▄▀▀▀▀▀▀████████████▙
                      █████████      █████████
                      █████████      █████████
                      █████████      █████████
                      ▜███████▛      ▜███████▛`)
}

func (m model) renderProjectMenuView() string {
	var s string
	
	// Project header
	if m.currentProj != nil {
		projName := titleStyle.Render(m.currentProj.Name)
		orgName := ""
		if m.currentOrg != nil {
			orgName = dimStyle.Render(" @ " + m.currentOrg.Name)
		}
		s += "\n"
		s += fmt.Sprintf("  📁 %s%s\n", projName, orgName)
		s += "\n"
	}
	
	// Breadcrumb
	breadcrumb := dimStyle.Render("  Organizations")
	breadcrumb += dimStyle.Render(" › ")
	if m.currentOrg != nil {
		breadcrumb += dimStyle.Render(m.currentOrg.Name)
	}
	breadcrumb += dimStyle.Render(" › ")
	if m.currentProj != nil {
		breadcrumb += lipgloss.NewStyle().Foreground(greenColor).Render(m.currentProj.Name)
	}
	s += breadcrumb + "\n\n"
	
	// Divider
	s += "  " + dimStyle.Render("─────────────────────────────────────────") + "\n\n"
	
	// Menu items with nice styling
	menuBox := lipgloss.NewStyle().
		Padding(0, 2)
	
	var menuContent string
	for i, item := range m.menuItems {
		if item.name == "" {
			// Separator
			menuContent += "\n"
			continue
		}
		
		prefix := "  "
		if i == m.selected {
			prefix = accentStyle.Render("▸ ")
		}
		
		icon := item.icon
		if icon != "" {
			icon += " "
		}
		
		var line string
		if i == m.selected {
			line = prefix + icon + selectedStyle.Render(item.name)
			if item.description != "" {
				line += "  " + dimStyle.Render(item.description)
			}
		} else {
			line = prefix + icon + normalStyle.Render(item.name)
			if item.description != "" {
				line += "  " + dimStyle.Render(item.description)
			}
		}
		
		menuContent += line + "\n"
	}
	
	s += menuBox.Render(menuContent)
	s += "\n"
	
	return s
}

func (m model) renderFooter() string {
	var help string
	
	switch m.currentView {
	case viewLogin:
		help = "enter login • q quit"
	case viewOrganizations:
		help = "↑↓ navigate • enter select • q quit"
	case viewProjects:
		help = "←→↑↓ navigate • enter select • esc back • q quit"
	case viewProjectMenu:
		help = "↑↓ navigate • enter select • esc back • q quit"
	case viewDashboard:
		help = "esc back • q quit"
	case viewDecisions:
		help = "↑↓ navigate • [n]ew • [a]ccept • [d]eprecate • esc back • q quit"
	case viewMemories:
		help = "↑↓ navigate • [n]ew • [d]elete • esc back • q quit"
	case viewCapsules:
		help = "↑↓ navigate • esc back • q quit"
	case viewTasks:
		help = "↑↓ navigate • [n]ew • [t]oggle • [d]elete • esc back • q quit"
	case viewBrain:
		help = "esc back • q quit"
	case viewHopper:
		help = "Type your message • enter send • esc back"
	}
	
	return "  " + helpStyle.Render(help) + "\n"
}

// GetSelectedCommand returns the action to execute
func (m model) GetSelectedCommand() string {
	return m.executeCmd
}

// ============================================================================
// HELPERS
// ============================================================================

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// sanitizeMarkdownForTerminal removes basic markdown markers for CLI display.
func sanitizeMarkdownForTerminal(s string) string {
	replacer := strings.NewReplacer(
		"**", "",
		"__", "",
		"`", "",
	)
	return replacer.Replace(s)
}

// ============================================================================
// FEATURE VIEWS
// ============================================================================

func (m model) renderDashboardView() string {
	var s string
	
	// Title
	s += "\n"
	s += "  " + titleStyle.Render("📊 Dashboard") + "\n"
	if m.currentProj != nil {
		s += "  " + dimStyle.Render(m.currentProj.Name) + "\n"
	}
	s += "\n"
	
	// KPI Cards - fixed width with margin
	kpiStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(grayColor).
		Padding(0, 2).
		Width(18).
		MarginRight(1)
	
	// Count stats
	decisionCount := len(m.decisions)
	acceptedCount := 0
	pendingCount := 0
	for _, d := range m.decisions {
		if d.Status == "ACCEPTED" {
			acceptedCount++
		} else if d.Status == "PENDING" {
			pendingCount++
		}
	}
	
	memoryCount := len(m.memories)
	activeTaskCount := 0
	for _, t := range m.tasks {
		if t.Status == "TODO" || t.Status == "IN_PROGRESS" {
			activeTaskCount++
		}
	}
	capsuleCount := len(m.capsules)
	
	// Row 1: KPIs
	kpi1 := kpiStyle.Render(fmt.Sprintf("📋 Decisions\n%s", 
		lipgloss.NewStyle().Bold(true).Foreground(cyanColor).Render(fmt.Sprintf("%d", decisionCount))))
	kpi2 := kpiStyle.Render(fmt.Sprintf("💾 Memories\n%s", 
		lipgloss.NewStyle().Bold(true).Foreground(magentaColor).Render(fmt.Sprintf("%d", memoryCount))))
	kpi3 := kpiStyle.Render(fmt.Sprintf("✅ Tasks\n%s", 
		lipgloss.NewStyle().Bold(true).Foreground(greenColor).Render(fmt.Sprintf("%d active", activeTaskCount))))
	kpi4 := kpiStyle.Render(fmt.Sprintf("📦 Capsules\n%s", 
		lipgloss.NewStyle().Bold(true).Foreground(yellowColor).Render(fmt.Sprintf("%d", capsuleCount))))
	
	row := lipgloss.JoinHorizontal(lipgloss.Top, kpi1, kpi2, kpi3, kpi4)
	s += lipgloss.NewStyle().MarginLeft(2).Render(row) + "\n\n"
	
	// Stats summary
	s += "  " + dimStyle.Render("─────────────────────────────────────────────────────────────────────────") + "\n\n"
	
	// Decision stats
	s += "  " + titleStyle.Render("Decision Status") + "\n"
	s += fmt.Sprintf("  %s %d  %s %d  %s %d\n\n",
		lipgloss.NewStyle().Foreground(greenColor).Render("● Accepted:"), acceptedCount,
		lipgloss.NewStyle().Foreground(yellowColor).Render("● Pending:"), pendingCount,
		lipgloss.NewStyle().Foreground(grayColor).Render("● Draft:"), decisionCount-acceptedCount-pendingCount)
	
	// Recent decisions
	s += "  " + titleStyle.Render("Recent Decisions") + "\n"
	if len(m.decisions) == 0 {
		s += "  " + dimStyle.Render("No decisions yet") + "\n"
	} else {
		for i, d := range m.decisions {
			if i >= 3 {
				break
			}
			statusIcon := "○"
			statusColor := grayColor
			switch d.Status {
			case "ACCEPTED":
				statusIcon = "●"
				statusColor = greenColor
			case "PENDING":
				statusIcon = "◐"
				statusColor = yellowColor
			}
			s += fmt.Sprintf("  %s %s\n", 
				lipgloss.NewStyle().Foreground(statusColor).Render(statusIcon),
				truncateString(d.Statement, 60))
		}
	}
	s += "\n"
	
	// Recent tasks
	s += "  " + titleStyle.Render("Recent Tasks") + "\n"
	if len(m.tasks) == 0 {
		s += "  " + dimStyle.Render("No tasks yet") + "\n"
	} else {
		for i, t := range m.tasks {
			if i >= 3 {
				break
			}
			statusIcon := "○"
			switch t.Status {
			case "DONE":
				statusIcon = "✓"
			case "IN_PROGRESS":
				statusIcon = "▶"
			}
			priorityColor := grayColor
			switch t.Priority {
			case "HIGH":
				priorityColor = lipgloss.Color("196")
			case "MEDIUM":
				priorityColor = yellowColor
			}
			s += fmt.Sprintf("  %s %s %s\n", 
				statusIcon,
				truncateString(t.Title, 50),
				lipgloss.NewStyle().Foreground(priorityColor).Render(t.Priority))
		}
	}
	
	return s
}

func (m model) renderDecisionsView() string {
	var s string
	
	// Title
	s += "  " + titleStyle.Render("📋 Decisions") + "\n"
	if m.currentProj != nil {
		s += "  " + dimStyle.Render(m.currentProj.Name) + "\n"
	}
	s += "\n"
	
	// Header
	headerStyle := lipgloss.NewStyle().Foreground(dimColor).Bold(true)
	s += "  " + headerStyle.Render(fmt.Sprintf("%-6s %-50s %-12s %s", "Status", "Statement", "Date", "Tags")) + "\n"
	s += "  " + dimStyle.Render("─────────────────────────────────────────────────────────────────────────────────") + "\n"
	
	if len(m.decisions) == 0 {
		s += "\n  " + dimStyle.Render("No decisions found. Press [n] to create one.") + "\n"
	} else {
		// Selected row style with background
		selectedRowStyle := lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("51")).
			Bold(true).
			Width(85)
		
		// Pagination: show only 10 items at a time
		visibleItems := 10
		startIdx := m.scrollOffset
		endIdx := startIdx + visibleItems
		if endIdx > len(m.decisions) {
			endIdx = len(m.decisions)
		}
		
		// Show scroll indicator if there are items above
		if startIdx > 0 {
			s += "  " + dimStyle.Render(fmt.Sprintf("  ↑ %d more above", startIdx)) + "\n"
		}
		
		for i := startIdx; i < endIdx; i++ {
			d := m.decisions[i]
			statusIcon := "○"
			statusColor := grayColor
			switch d.Status {
			case "ACCEPTED":
				statusIcon = "✓"
				statusColor = greenColor
			case "PENDING":
				statusIcon = "◐"
				statusColor = yellowColor
			case "REJECTED":
				statusIcon = "✗"
				statusColor = lipgloss.Color("196")
			case "DEPRECATED":
				statusIcon = "◌"
				statusColor = dimColor
			}
			
			statement := truncateString(d.Statement, 45)
			date := ""
			if len(d.CreatedAt) >= 10 {
				date = d.CreatedAt[:10]
			}
			
			tags := ""
			if len(d.Tags) > 0 {
				tags = truncateString(fmt.Sprintf("%v", d.Tags), 15)
			}
			
			line := fmt.Sprintf("%s %-47s %-12s %s", lipgloss.NewStyle().Foreground(statusColor).Render(statusIcon), statement, date, tags)
			
			if i == m.selected {
				// Highlighted row with background color
				s += "  " + selectedRowStyle.Render("▸ "+line) + "\n"
			} else {
				s += "    " + line + "\n"
			}
		}
		
		// Show scroll indicator if there are items below
		remaining := len(m.decisions) - endIdx
		if remaining > 0 {
			s += "  " + dimStyle.Render(fmt.Sprintf("  ↓ %d more below", remaining)) + "\n"
		}
	}
	
	s += "\n"
	s += "  " + dimStyle.Render(fmt.Sprintf("Total: %d decisions | [a]ccept [x]deprecate [n]ew | ↑↓ scroll", len(m.decisions))) + "\n"
	
	return s
}

func (m model) renderMemoriesView() string {
	var s string
	
	// Title
	s += "\n"
	s += "  " + titleStyle.Render("💾 Memories") + "\n"
	if m.currentProj != nil {
		s += "  " + dimStyle.Render(m.currentProj.Name) + "\n"
	}
	s += "\n"
	
	// Header
	headerStyle := lipgloss.NewStyle().Foreground(dimColor).Bold(true)
	s += "  " + headerStyle.Render(fmt.Sprintf("%-60s %-12s %s", "Content", "Date", "Tags")) + "\n"
	s += "  " + dimStyle.Render("─────────────────────────────────────────────────────────────────────────────────") + "\n"
	
	if len(m.memories) == 0 {
		s += "\n  " + dimStyle.Render("No memories found. Press [n] to create one.") + "\n"
	} else {
		// Selected row style with background
		selectedRowStyle := lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("51")).
			Bold(true).
			Width(85)
		
		// Pagination: show only 10 items at a time
		visibleItems := 10
		startIdx := m.scrollOffset
		endIdx := startIdx + visibleItems
		if endIdx > len(m.memories) {
			endIdx = len(m.memories)
		}
		
		// Show scroll indicator if there are items above
		if startIdx > 0 {
			s += "  " + dimStyle.Render(fmt.Sprintf("  ↑ %d more above", startIdx)) + "\n"
		}
		
		for i := startIdx; i < endIdx; i++ {
			mem := m.memories[i]
			content := truncateString(mem.Content, 55)
			date := ""
			if len(mem.CreatedAt) >= 10 {
				date = mem.CreatedAt[:10]
			}
			
			tags := ""
			if len(mem.Tags) > 0 {
				tags = truncateString(fmt.Sprintf("%v", mem.Tags), 15)
			}
			
			line := fmt.Sprintf("%-57s %-12s %s", content, date, tags)
			
			if i == m.selected {
				// Highlighted row with background color
				s += "  " + selectedRowStyle.Render("▸ "+line) + "\n"
			} else {
				s += "    " + line + "\n"
			}
		}
		
		// Show scroll indicator if there are items below
		remaining := len(m.memories) - endIdx
		if remaining > 0 {
			s += "  " + dimStyle.Render(fmt.Sprintf("  ↓ %d more below", remaining)) + "\n"
		}
	}
	
	s += "\n"
	s += "  " + dimStyle.Render(fmt.Sprintf("Total: %d memories | [d]elete [n]ew | ↑↓ scroll", len(m.memories))) + "\n"
	
	return s
}

func (m model) renderCapsulesView() string {
	var s string
	
	// Title
	s += "\n"
	s += "  " + titleStyle.Render("📦 Capsules") + "\n"
	if m.currentProj != nil {
		s += "  " + dimStyle.Render(m.currentProj.Name) + "\n"
	}
	s += "\n"
	
	// Header
	headerStyle := lipgloss.NewStyle().Foreground(dimColor).Bold(true)
	s += "  " + headerStyle.Render(fmt.Sprintf("%-8s %-30s %-10s %-10s", "Status", "Name", "Decisions", "Memories")) + "\n"
	s += "  " + dimStyle.Render("─────────────────────────────────────────────────────────────────") + "\n"
	
	if len(m.capsules) == 0 {
		s += "\n  " + dimStyle.Render("No capsules found.") + "\n"
	} else {
		for i, c := range m.capsules {
			prefix := "  "
			if i == m.selected {
				prefix = accentStyle.Render("▸ ")
			}
			
			statusIcon := "○"
			statusColor := grayColor
			switch c.Status {
			case "FROZEN":
				statusIcon = "❄"
				statusColor = cyanColor
			case "DRAFT":
				statusIcon = "◐"
				statusColor = yellowColor
			}
			
			name := truncateString(c.Name, 28)
			decCount := len(c.DecisionIds)
			memCount := len(c.MemoryIds)
			
			line := fmt.Sprintf("%-30s %-10d %-10d", name, decCount, memCount)
			if i == m.selected {
				line = selectedStyle.Render(line)
			}
			
			s += fmt.Sprintf("%s%s %s\n", prefix, lipgloss.NewStyle().Foreground(statusColor).Render(statusIcon), line)
		}
	}
	
	s += "\n"
	s += "  " + dimStyle.Render(fmt.Sprintf("Total: %d capsules", len(m.capsules))) + "\n"
	
	return s
}

func (m model) renderTasksView() string {
	var s string
	
	// Title
	s += "\n"
	s += "  " + titleStyle.Render("✅ Tasks") + "\n"
	if m.currentProj != nil {
		s += "  " + dimStyle.Render(m.currentProj.Name) + "\n"
	}
	s += "\n"
	
	// Header
	headerStyle := lipgloss.NewStyle().Foreground(dimColor).Bold(true)
	s += "  " + headerStyle.Render(fmt.Sprintf("%-4s %-50s %-12s %-8s", "Done", "Title", "Status", "Priority")) + "\n"
	s += "  " + dimStyle.Render("─────────────────────────────────────────────────────────────────────────────") + "\n"
	
	if len(m.tasks) == 0 {
		s += "\n  " + dimStyle.Render("No tasks found. Press [n] to create one.") + "\n"
	} else {
		// Selected row style with background
		selectedRowStyle := lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("51")).
			Bold(true).
			Width(85)
		
		// Pagination: show only 10 items at a time
		visibleItems := 10
		startIdx := m.scrollOffset
		endIdx := startIdx + visibleItems
		if endIdx > len(m.tasks) {
			endIdx = len(m.tasks)
		}
		
		// Show scroll indicator if there are items above
		if startIdx > 0 {
			s += "  " + dimStyle.Render(fmt.Sprintf("  ↑ %d more above", startIdx)) + "\n"
		}
		
		for i := startIdx; i < endIdx; i++ {
			t := m.tasks[i]
			checkbox := "[ ]"
			if t.Status == "DONE" {
				checkbox = "[✓]"
			}
			
			title := truncateString(t.Title, 45)
			
			statusColor := grayColor
			switch t.Status {
			case "IN_PROGRESS":
				statusColor = cyanColor
			case "DONE":
				statusColor = greenColor
			case "REVIEW":
				statusColor = magentaColor
			}
			
			priorityColor := grayColor
			switch t.Priority {
			case "HIGH":
				priorityColor = lipgloss.Color("196")
			case "MEDIUM":
				priorityColor = yellowColor
			case "LOW":
				priorityColor = dimColor
			}
			
			statusStr := lipgloss.NewStyle().Foreground(statusColor).Render(fmt.Sprintf("%-12s", t.Status))
			priorityStr := lipgloss.NewStyle().Foreground(priorityColor).Render(t.Priority)
			line := fmt.Sprintf("%s %-47s %s %s", checkbox, title, statusStr, priorityStr)
			
			if i == m.selected {
				// Highlighted row with background color
				s += "  " + selectedRowStyle.Render("▸ "+line) + "\n"
			} else {
				s += "    " + line + "\n"
			}
		}
		
		// Show scroll indicator if there are items below
		remaining := len(m.tasks) - endIdx
		if remaining > 0 {
			s += "  " + dimStyle.Render(fmt.Sprintf("  ↓ %d more below", remaining)) + "\n"
		}
	}
	
	s += "\n"
	
	// Stats
	todoCount := 0
	inProgressCount := 0
	doneCount := 0
	for _, t := range m.tasks {
		switch t.Status {
		case "TODO":
			todoCount++
		case "IN_PROGRESS":
			inProgressCount++
		case "DONE":
			doneCount++
		}
	}
	s += "  " + dimStyle.Render(fmt.Sprintf("Total: %d | Todo: %d | In Progress: %d | Done: %d | [t]oggle [d]elete [n]ew", 
		len(m.tasks), todoCount, inProgressCount, doneCount)) + "\n"
	
	return s
}

func (m model) renderBrainView() string {
	var s string
	
	// Title
	s += "\n"
	s += "  " + titleStyle.Render("🧠 Brain") + "\n"
	if m.currentProj != nil {
		s += "  " + dimStyle.Render(m.currentProj.Name) + "\n"
	}
	s += "\n"
	
	if m.graphStats == nil {
		s += "  " + dimStyle.Render("Loading brain stats...") + "\n"
		return s
	}
	
	// Stats cards - fixed width with margin
	statStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(grayColor).
		Padding(1, 3).
		Width(25).
		MarginRight(2)
	
	nodeCard := statStyle.Render(fmt.Sprintf("🔵 Nodes\n%s",
		lipgloss.NewStyle().Bold(true).Foreground(cyanColor).Render(fmt.Sprintf("%d", m.graphStats.NodeCount))))
	
	edgeCard := statStyle.Render(fmt.Sprintf("🔗 Connections\n%s",
		lipgloss.NewStyle().Bold(true).Foreground(magentaColor).Render(fmt.Sprintf("%d", m.graphStats.EdgeCount))))
	
	row := lipgloss.JoinHorizontal(lipgloss.Top, nodeCard, edgeCard)
	s += lipgloss.NewStyle().MarginLeft(2).Render(row) + "\n\n"
	
	// Node types breakdown
	s += "  " + dimStyle.Render("─────────────────────────────────────────────────────") + "\n\n"
	s += "  " + titleStyle.Render("Nodes by Type") + "\n\n"
	
	if m.graphStats.NodesByType != nil {
		for nodeType, count := range m.graphStats.NodesByType {
			icon := "○"
			color := grayColor
			switch nodeType {
			case "decision":
				icon = "📋"
				color = cyanColor
			case "memory":
				icon = "💾"
				color = magentaColor
			case "task":
				icon = "✅"
				color = greenColor
			case "capsule":
				icon = "📦"
				color = yellowColor
			case "code_chunk":
				icon = "📄"
				color = blueColor
			}
			s += fmt.Sprintf("  %s %s: %s\n", icon, nodeType, 
				lipgloss.NewStyle().Foreground(color).Bold(true).Render(fmt.Sprintf("%d", count)))
		}
	} else {
		s += "  " + dimStyle.Render("No node data available") + "\n"
	}
	
	return s
}

func (m model) renderHopperView() string {
	var s string
	
	// Title with cute hopper emoji
	s += "  " + lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true).
		Render("🐰 Hopper AI") + "\n"
	if m.currentProj != nil {
		s += "  " + dimStyle.Render(m.currentProj.Name) + "\n"
	}
	
	// Show RAG context status
	if m.hopperContextLoaded {
		contextInfo := fmt.Sprintf("📚 Context: %d decisions, %d memories loaded", 
			len(m.hopperDecisions), len(m.hopperMemories))
		s += "  " + lipgloss.NewStyle().
			Foreground(lipgloss.Color("71")).
			Render(contextInfo) + "\n"
	} else if m.loading {
		s += "  " + lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Render("⏳ Loading project context...") + "\n"
	}
	s += "\n"
	
	// If still loading, show loading message
	if m.loading && !m.hopperContextLoaded {
		loadingStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Italic(true)
		s += "  " + loadingStyle.Render("Hopper is loading project knowledge...") + "\n"
		return s
	}
	
	// Chat area
	chatAreaStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1, 2).
		Width(80)
	
	var chatContent string
	
	if len(m.chatMessages) == 0 && !m.chatStreaming {
		// Welcome message with context-aware info
		welcomeStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("251")).
			Italic(true)
		
		welcomeText := "Hi! I'm Hopper, your project assistant. 🐰\n\n"
		if m.hopperContextLoaded {
			welcomeText += fmt.Sprintf("I've loaded %d decisions and %d memories from your project.\n", 
				len(m.hopperDecisions), len(m.hopperMemories))
			welcomeText += "Ask me anything about your project!\n\n"
		}
		welcomeText += "Examples:\n"
		welcomeText += "• \"What decisions do we have about authentication?\"\n"
		welcomeText += "• \"Summarize our coding standards\"\n"
		welcomeText += "• \"What do our memories say about API design?\""
		
		chatContent = welcomeStyle.Render(welcomeText)
	} else {
		// Show conversation with pagination (last 6 messages)
		startIdx := 0
		if len(m.chatMessages) > 6 {
			startIdx = len(m.chatMessages) - 6
			chatContent += dimStyle.Render(fmt.Sprintf("... %d earlier messages ...\n\n", startIdx))
		}
		
		for i := startIdx; i < len(m.chatMessages); i++ {
			msg := m.chatMessages[i]
			if msg.Role == "user" {
				// User message style
				userStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("117")).
					Bold(true)
				chatContent += userStyle.Render("You: ") + msg.Content + "\n\n"
			} else {
				// Assistant message style  
				assistantStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("213")).
					Bold(true)
				// Word wrap long responses
				content := sanitizeMarkdownForTerminal(msg.Content)
				if len(content) > 300 {
					content = content[:300] + "..."
				}
				chatContent += assistantStyle.Render("Hopper: ") + content + "\n\n"
			}
		}
		
		// Show streaming indicator
		if m.chatStreaming {
			chatContent += lipgloss.NewStyle().
				Foreground(lipgloss.Color("213")).
				Blink(true).
				Render("Hopper is thinking... 🐰")
		}
	}
	
	s += chatAreaStyle.Render(chatContent) + "\n\n"
	
	// Input area
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("117")).
		Padding(0, 1).
		Width(80)
	
	inputPrompt := lipgloss.NewStyle().
		Foreground(lipgloss.Color("117")).
		Bold(true).
		Render("> ")
	
	cursorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("117")).
		Blink(true)
	
	inputDisplay := m.chatInput
	if !m.chatStreaming {
		inputDisplay += cursorStyle.Render("█")
	}
	
	s += inputStyle.Render(inputPrompt + inputDisplay) + "\n"
	
	// Error message if any
	if m.errorMsg != "" {
		s += "\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("Error: "+m.errorMsg) + "\n"
	}
	
	return s
}

// ============================================================================
// RUN INTERACTIVE
// ============================================================================

func RunInteractive() (string, error) {
	cfg, _ := config.GetConfig()
	if cfg == nil {
		cfg = &config.Config{}
	}
	
	p := tea.NewProgram(NewInteractiveModel(cfg), tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	m, ok := finalModel.(model)
	if !ok {
		return "", nil
	}

	return m.GetSelectedCommand(), nil
}

// ============================================================================
// EXECUTE LOGIN
// ============================================================================

func ExecuteLogin(cfg *config.Config) error {
	if cfg == nil {
		cfg = &config.Config{}
	}

	if cfg.APIURL == "" {
		cfg.APIURL = "http://localhost:8080"
	}
	if cfg.WebURL == "" {
		cfg.WebURL = "http://localhost:3000"
	}

	client := api.NewClient(cfg)
	deviceName := "CLI"

	fmt.Println()
	fmt.Println("  Initializing login...")

	initResp, err := client.DeviceAuthInit(deviceName)
	if err != nil {
		return fmt.Errorf("failed to initialize login: %w", err)
	}

	authURL := fmt.Sprintf("%s/auth/device?code=%s", cfg.WebURL, initResp.Code)

	fmt.Printf("  Device Code: %s\n\n", accentStyle.Render(initResp.Code))
	fmt.Println("  Opening browser to complete sign-in...")

	if err := openBrowser(authURL); err != nil {
		fmt.Println("  Could not open browser automatically.")
	}

	fmt.Println()
	fmt.Println("  If the browser doesn't open, visit:")
	fmt.Printf("  %s\n\n", logoStyle.Render(authURL))

	fmt.Println("  Waiting for authentication...")
	fmt.Println("  (Press Ctrl+C to cancel)")
	fmt.Println()

	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinnerIdx := 0
	maxAttempts := 300

	for attempt := 0; attempt < maxAttempts; attempt++ {
		fmt.Printf("\r  %s Waiting for browser authentication... (%ds)",
			logoStyle.Render(spinner[spinnerIdx]), attempt*2)
		spinnerIdx = (spinnerIdx + 1) % len(spinner)

		resp, err := client.DeviceAuthPoll(initResp.Code)
		if err != nil {
			fmt.Println()
			return fmt.Errorf("failed to check login status: %w", err)
		}

		switch resp.Status {
		case "complete":
			fmt.Printf("\r  %s Authentication complete!                    \n\n",
				statusOnStyle.Render("✓"))

			cfg.Token = resp.Token
			cfg.User = &config.User{
				ID:        resp.UserID,
				Email:     resp.Email,
				Name:      resp.Name,
				AvatarURL: resp.AvatarURL,
			}

			if err := config.SaveConfig(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("  Signed in as: %s (%s)\n\n",
				titleStyle.Render(resp.Name),
				dimStyle.Render(resp.Email))

			return nil

		case "expired":
			fmt.Println()
			return fmt.Errorf("login session expired - please try again")

		case "pending":
			// Continue polling
		}

		time.Sleep(2 * time.Second)
	}

	fmt.Println()
	return fmt.Errorf("login timed out")
}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}

// ShowOrganizations - kept for backwards compatibility
func ShowOrganizations(cfg *config.Config) {
	// Now handled in TUI
}

// ShowProjects - kept for backwards compatibility  
func ShowProjects(cfg *config.Config) {
	// Now handled in TUI
}
