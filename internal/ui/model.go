package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Hopsule/cli-tool/internal/api"
	"github.com/Hopsule/cli-tool/internal/config"
	tea "github.com/charmbracelet/bubbletea"
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
	viewHopper
)

// ============================================================================
// VIEW MODE — card grid vs datatable list
// ============================================================================

type viewMode int

const (
	viewModeCard viewMode = iota
	viewModeList
)

// ============================================================================
// CHAT SESSION
// ============================================================================

type chatSession struct {
	ID       string
	Topic    string
	Messages []api.ChatMessage
	Time     string
}

// ============================================================================
// MENU ITEM
// ============================================================================

type menuItem struct {
	icon        string
	name        string
	description string
	action      string
}

// ============================================================================
// MODEL
// ============================================================================

type model struct {
	cfg         *config.Config
	client      *api.Client
	currentView viewType

	organizations []*api.Organization
	projects      []*api.Project
	currentOrg    *api.Organization
	currentProj   *api.Project

	menuItems []menuItem

	decisions  []api.Decision
	memories   []*api.Memory
	tasks      []*api.Task
	capsules   []*api.Capsule
	graphStats *api.GraphStats

	chatMessages        []api.ChatMessage
	chatInput           string
	chatStreaming        bool
	streamingContent    string
	hopperDecisions     []api.Decision
	hopperMemories      []*api.Memory
	hopperContextLoaded bool
	hopperSessionID     string

	chatSessions       []chatSession
	chatSessionIdx     int
	hopperSidebarFocus bool
	chatScroll         int

	loginDeviceCode string
	loginAuthURL    string
	loginPolling    bool
	loginSpinnerIdx int
	loginPollCount  int

	selected     int
	scrollOffset int

	searchMode   bool
	searchQuery  string
	detailView   bool
	detailScroll int

	// Status filter for decisions
	statusFilter string

	// Tasks view mode: "list" or "kanban"
	tasksViewMode string

	// Card vs List toggle for decisions/memories/capsules
	currentViewMode viewMode

	// List view pagination
	listPage     int
	listPageSize int

	// Delete confirmation
	confirmDelete     bool
	confirmDeleteID   string
	confirmDeleteType string

	// Remembered selection indices for navigation
	menuSelectedIdx int
	projSelectedIdx int
	orgSelectedIdx  int

	// Input mode
	inputMode     bool
	inputPrompt   string
	inputValue    string
	inputAction   string
	editingItemID string

	// Terminal dimensions
	width  int
	height int

	loading  bool
	errorMsg string

	executeCmd string

	// Tracks the project ID that is already initialized in the current directory (.hopsule)
	initializedProjectID string

	// Current working directory info for project matching
	cwdInGitRepo bool
	cwdDirName   string
}

// ============================================================================
// CONSTRUCTOR
// ============================================================================

// NewInteractiveModel creates a new TUI model.
func NewInteractiveModel(cfg *config.Config) model {
	isLoggedIn := cfg != nil && cfg.IsAuthenticated()
	m := model{
		cfg:             cfg,
		selected:        0,
		hopperSessionID: fmt.Sprintf("cli-%d", time.Now().UnixNano()),
		tasksViewMode:   "list",
		currentViewMode: viewModeList,
		listPage:        0,
		listPageSize:    20,
	}
	if isLoggedIn {
		m.client = api.NewClient(cfg)
		m.currentView = viewOrganizations
		m.loading = true
	} else {
		m.currentView = viewLogin
	}

	// Detect current directory info
	if cwd, err := os.Getwd(); err == nil {
		m.cwdDirName = filepath.Base(cwd)

		// Check if inside a git repo
		cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
		cmd.Dir = cwd
		if out, err := cmd.Output(); err == nil && strings.TrimSpace(string(out)) == "true" {
			m.cwdInGitRepo = true
		}
	}

	// Check if there's an initialized project in the current directory
	if projCfg, _, err := config.LoadProjectConfig(); err == nil && projCfg != nil {
		m.initializedProjectID = projCfg.Project.ID
	}

	return m
}

// Init implements tea.Model.
func (m model) Init() tea.Cmd {
	if m.loading {
		return m.loadData
	}
	return nil
}

// GetSelectedCommand returns the command to execute after TUI exit.
func (m model) GetSelectedCommand() string {
	return m.executeCmd
}

// ============================================================================
// DATA LOADERS
// ============================================================================

func (m model) loadData() tea.Msg {
	if m.client == nil {
		return dataLoadedMsg{err: fmt.Errorf("not authenticated")}
	}
	meResp, err := m.client.GetMe()
	if err != nil {
		return dataLoadedMsg{err: err}
	}
	return dataLoadedMsg{organizations: meResp.Organizations, projects: meResp.Projects}
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

func (m model) loadDashboardData() tea.Msg {
	if m.client == nil || m.currentProj == nil {
		return dashboardLoadedMsg{err: fmt.Errorf("not authenticated or no project selected")}
	}
	decisions, _ := m.client.ListDecisions(m.currentProj.ID)
	memories, _ := m.client.ListMemories(m.currentProj.ID)
	tasks, _ := m.client.ListTasks(m.currentProj.ID)
	capsules, _ := m.client.ListCapsules(m.currentProj.ID)
	return dashboardLoadedMsg{decisions: decisions, memories: memories, tasks: tasks, capsules: capsules}
}

func (m model) loadHopperContext() tea.Msg {
	if m.client == nil || m.currentProj == nil {
		return hopperContextLoadedMsg{err: fmt.Errorf("not authenticated or no project selected")}
	}
	decisions, err := m.client.ListDecisions(m.currentProj.ID)
	if err != nil {
		return hopperContextLoadedMsg{err: fmt.Errorf("failed to load decisions: %w", err)}
	}
	memories, err := m.client.ListMemories(m.currentProj.ID)
	if err != nil {
		return hopperContextLoadedMsg{err: fmt.Errorf("failed to load memories: %w", err)}
	}
	chatHistory, _ := m.client.ListChatHistory(m.currentProj.ID, 20)
	return hopperContextLoadedMsg{decisions: decisions, memories: memories, chatHistory: chatHistory}
}

// ============================================================================
// HOPPER CHAT HELPERS
// ============================================================================

func (m *model) hopperNewChat() {
	if len(m.chatMessages) > 0 {
		m.hopperSaveCurrentChat()
	}
	m.chatMessages = nil
	m.chatInput = ""
	m.streamingContent = ""
	m.chatScroll = 0
	m.hopperSessionID = fmt.Sprintf("hopper-%d-%d", time.Now().UnixNano(), time.Now().UnixMilli()%1000)
}

func (m *model) hopperSaveCurrentChat() {
	if len(m.chatMessages) == 0 {
		return
	}
	topic := "New Chat"
	for _, msg := range m.chatMessages {
		if msg.Role == "user" {
			topic = truncateString(msg.Content, 40)
			break
		}
	}
	for i, s := range m.chatSessions {
		if s.ID == m.hopperSessionID {
			m.chatSessions[i].Messages = make([]api.ChatMessage, len(m.chatMessages))
			copy(m.chatSessions[i].Messages, m.chatMessages)
			m.chatSessions[i].Topic = topic
			return
		}
	}
	msgs := make([]api.ChatMessage, len(m.chatMessages))
	copy(msgs, m.chatMessages)
	sess := chatSession{
		ID:       m.hopperSessionID,
		Topic:    topic,
		Messages: msgs,
		Time:     time.Now().Format("15:04"),
	}
	m.chatSessions = append([]chatSession{sess}, m.chatSessions...)
	m.chatSessionIdx = 0
}

func (m model) saveChatHistoryToAPI() tea.Msg {
	if m.client == nil || m.currentProj == nil || len(m.chatMessages) == 0 {
		return chatHistorySavedMsg{}
	}
	topic := "New Chat"
	for _, msg := range m.chatMessages {
		if msg.Role == "user" {
			topic = truncateString(msg.Content, 40)
			break
		}
	}
	isLocalID := strings.HasPrefix(m.hopperSessionID, "hopper-") || strings.HasPrefix(m.hopperSessionID, "cli-")
	if isLocalID {
		entry, err := m.client.CreateChatHistory(m.currentProj.ID, &api.CreateChatHistoryRequest{
			Topic:    topic,
			Messages: m.chatMessages,
		})
		if err != nil {
			return chatHistorySavedMsg{err: err}
		}
		return chatHistorySavedMsg{chatID: entry.ID}
	}
	_, err := m.client.UpdateChatHistory(m.currentProj.ID, m.hopperSessionID, &api.UpdateChatHistoryRequest{
		Topic:    topic,
		Messages: m.chatMessages,
	})
	if err != nil {
		return chatHistorySavedMsg{err: err}
	}
	return chatHistorySavedMsg{chatID: m.hopperSessionID}
}

func (m model) loadChatHistoryFromAPI(sessionID string) tea.Cmd {
	return func() tea.Msg {
		if m.client == nil || m.currentProj == nil {
			return chatHistoryLoadedMsg{sessionID: sessionID, err: fmt.Errorf("not authenticated")}
		}
		entry, err := m.client.GetChatHistory(m.currentProj.ID, sessionID)
		if err != nil {
			return chatHistoryLoadedMsg{sessionID: sessionID, err: err}
		}
		var msgs []api.ChatMessage
		for _, em := range entry.Messages {
			msgs = append(msgs, api.ChatMessage{Role: em.Role, Content: em.Content})
		}
		return chatHistoryLoadedMsg{sessionID: sessionID, messages: msgs}
	}
}

func (m model) sendHopperMessage() (tea.Model, tea.Cmd) {
	if m.client == nil || m.currentProj == nil {
		m.errorMsg = "Not authenticated or no project selected"
		return m, nil
	}
	userMessage := m.chatInput
	m.chatInput = ""
	m.chatStreaming = true
	m.streamingContent = ""
	m.chatMessages = append(m.chatMessages, api.ChatMessage{Role: "user", Content: userMessage})
	m.chatScroll = max(0, len(m.chatMessages)-4)

	var taggedItems []api.TaggedItem
	decisionLimit := min(10, len(m.hopperDecisions))
	for i := 0; i < decisionLimit; i++ {
		d := m.hopperDecisions[i]
		taggedItems = append(taggedItems, api.TaggedItem{
			ID: d.ID, Type: "decision",
			Statement: truncateString(d.Statement, 160),
			Content:   truncateString(d.Rationale, 240),
		})
	}
	memoryLimit := min(15, len(m.hopperMemories))
	for i := 0; i < memoryLimit; i++ {
		mem := m.hopperMemories[i]
		taggedItems = append(taggedItems, api.TaggedItem{
			ID: mem.ID, Type: "memory", Content: truncateString(mem.Content, 240),
		})
	}

	client := m.client
	projectID := m.currentProj.ID
	projectName := m.currentProj.Name
	history := make([]api.ChatMessage, len(m.chatMessages)-1)
	copy(history, m.chatMessages[:len(m.chatMessages)-1])

	return m, func() tea.Msg {
		var fullResponse string
		req := &api.ChatRequest{
			Message: userMessage, ConversationHistory: history,
			TaggedItems: taggedItems, Stream: true,
			SessionID: m.hopperSessionID, ProjectName: projectName,
		}
		err := client.SendChatMessage(projectID, req, func(chunk string) {
			fullResponse += chunk
		})
		if err != nil {
			return chatStreamDoneMsg{err: err}
		}
		return chatStreamChunkMsg{chunk: fullResponse}
	}
}

// ============================================================================
// SELECTION HELPERS
// ============================================================================

func (m model) getSelectedDecision() *api.Decision {
	filtered := m.getFilteredDecisions()
	if m.selected >= 0 && m.selected < len(filtered) {
		d := filtered[m.selected]
		return &d
	}
	return nil
}

func (m model) getSelectedMemory() *api.Memory {
	filtered := m.getFilteredMemories()
	if m.selected >= 0 && m.selected < len(filtered) {
		return filtered[m.selected]
	}
	return nil
}

func (m model) getSelectedTask() *api.Task {
	filtered := m.getFilteredTasks()
	if m.selected >= 0 && m.selected < len(filtered) {
		return filtered[m.selected]
	}
	return nil
}

func (m model) getSelectedCapsule() *api.Capsule {
	filtered := m.getFilteredCapsules()
	if m.selected >= 0 && m.selected < len(filtered) {
		return filtered[m.selected]
	}
	return nil
}

func (m model) getMaxSelection() int {
	switch m.currentView {
	case viewLogin:
		return 1
	case viewOrganizations:
		return len(m.organizations) + 1
	case viewProjects:
		return len(m.getOrgProjects())
	case viewProjectMenu:
		return len(m.menuItems)
	case viewDecisions:
		return len(m.getFilteredDecisions())
	case viewMemories:
		return len(m.getFilteredMemories())
	case viewTasks:
		return len(m.getFilteredTasks())
	case viewCapsules:
		return len(m.getFilteredCapsules())
	case viewDashboard:
		return 0
	}
	return 0
}

func (m model) getVisibleItems() int {
	available := m.height - 8
	if available < 5 {
		return 5
	}
	if available > 30 {
		return 30
	}
	return available
}

func (m model) getGridCols() int {
	if m.currentView != viewDecisions && m.currentView != viewMemories && m.currentView != viewCapsules {
		return 1
	}
	// Responsive: 3 cols if terminal wide enough, fall back to 2 or 1
	available := m.width - 8 // account for scroll area border + padding
	if available >= 90 {
		return 3
	}
	if available >= 56 {
		return 2
	}
	return 1
}

func (m model) getCardWidth() int {
	cols := m.getGridCols()
	available := m.width - 10 // border(2) + padding(4) + margin(4)
	if available < 28 {
		available = 28
	}
	gap := 2
	cardW := (available - (cols-1)*gap) / cols
	if cardW < 24 {
		cardW = 24
	}
	return cardW
}

// getVisibleRows returns how many card rows fit in the scroll viewport.
// Target is 3 rows for 3x3, but adapts to small terminals.
func (m model) getVisibleRows() int {
	cardHeight := 8 // each card ~8 lines including border
	available := m.height - 14 // header/breadcrumb/footer/filter bars
	if available < cardHeight {
		return 1
	}
	rows := available / cardHeight
	if rows > 3 {
		rows = 3
	}
	if rows < 1 {
		rows = 1
	}
	return rows
}

func (m model) getVisibleCards() int {
	if m.currentView == viewDecisions || m.currentView == viewMemories || m.currentView == viewCapsules {
		return m.getVisibleRows() * m.getGridCols()
	}
	return m.getVisibleItems()
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
// FILTER FUNCTIONS
// ============================================================================

func (m model) getFilteredDecisions() []api.Decision {
	var filtered []api.Decision
	for _, d := range m.decisions {
		if m.statusFilter != "" && d.Status != m.statusFilter {
			continue
		}
		if m.searchQuery != "" {
			query := strings.ToLower(m.searchQuery)
			if !strings.Contains(strings.ToLower(d.Statement), query) &&
				!strings.Contains(strings.ToLower(d.Rationale), query) {
				continue
			}
		}
		filtered = append(filtered, d)
	}
	if filtered == nil {
		return []api.Decision{}
	}
	return filtered
}

func (m model) getFilteredMemories() []*api.Memory {
	if m.searchQuery == "" {
		return m.memories
	}
	query := strings.ToLower(m.searchQuery)
	var filtered []*api.Memory
	for _, mem := range m.memories {
		if strings.Contains(strings.ToLower(mem.Content), query) {
			filtered = append(filtered, mem)
		}
	}
	return filtered
}

func (m model) getFilteredTasks() []*api.Task {
	if m.searchQuery == "" {
		return m.tasks
	}
	query := strings.ToLower(m.searchQuery)
	var filtered []*api.Task
	for _, t := range m.tasks {
		if strings.Contains(strings.ToLower(t.Title), query) ||
			strings.Contains(strings.ToLower(t.Description), query) {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

func (m model) getFilteredCapsules() []*api.Capsule {
	if m.searchQuery == "" {
		return m.capsules
	}
	query := strings.ToLower(m.searchQuery)
	var filtered []*api.Capsule
	for _, c := range m.capsules {
		if strings.Contains(strings.ToLower(c.Name), query) ||
			strings.Contains(strings.ToLower(c.Description), query) {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

// isCardView returns true if the current entity view is in card mode.
func (m model) isCardView() bool {
	return m.currentViewMode == viewModeCard
}

// getListPageItems returns the page-slice indices for the current list view page.
func (m model) getListPageItems(total int) (start, end int) {
	if m.listPageSize <= 0 {
		m.listPageSize = 20
	}
	totalPages := (total + m.listPageSize - 1) / m.listPageSize
	if m.listPage >= totalPages {
		m.listPage = max(0, totalPages-1)
	}
	start = m.listPage * m.listPageSize
	end = min(start+m.listPageSize, total)
	return start, end
}
