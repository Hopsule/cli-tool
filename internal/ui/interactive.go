package ui

import (
	"encoding/json"
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
// STYLES - Adaptive colors for dark/light terminal support
// ============================================================================

var (
	cyanColor    = lipgloss.AdaptiveColor{Light: "30", Dark: "51"}
	greenColor   = lipgloss.AdaptiveColor{Light: "28", Dark: "46"}
	yellowColor  = lipgloss.AdaptiveColor{Light: "136", Dark: "226"}
	magentaColor = lipgloss.AdaptiveColor{Light: "127", Dark: "201"}
	purpleColor  = lipgloss.AdaptiveColor{Light: "93", Dark: "141"}
	grayColor    = lipgloss.AdaptiveColor{Light: "244", Dark: "244"}
	dimColor     = lipgloss.AdaptiveColor{Light: "244", Dark: "240"}
	blueColor    = lipgloss.AdaptiveColor{Light: "25", Dark: "39"}
	redColor     = lipgloss.AdaptiveColor{Light: "160", Dark: "196"}

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "255"}).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "255"}).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "240", Dark: "250"})

	logoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "238", Dark: "252"})

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "244", Dark: "244"})

	accentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "238", Dark: "252"})

	versionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "248", Dark: "240"})

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "244", Dark: "244"})

	statusOnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "28", Dark: "46"}).
			Bold(true)

	cardBorderColor         = lipgloss.AdaptiveColor{Light: "240", Dark: "240"}
	cardSelectedBorderColor = lipgloss.AdaptiveColor{Light: "28", Dark: "46"}

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(cardBorderColor).
			Padding(1, 2).
			Width(32).
			MarginRight(2)

	selectedCardStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(cardSelectedBorderColor).
				Padding(1, 2).
				Width(32).
				MarginRight(2)

	cardTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "255"}).
			Bold(true)

	cardDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "244", Dark: "244"})

	breadcrumbStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "238", Dark: "252"})

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "244", Dark: "244"})
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
// CHAT SESSION
// ============================================================================

type chatSession struct {
	ID       string
	Topic    string
	Messages []api.ChatMessage
	Time     string
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

	// Chat history (sidebar sessions)
	chatSessions       []chatSession
	chatSessionIdx     int
	hopperSidebarFocus bool // true = sidebar focused, false = chat focused
	chatScroll         int  // scroll offset for messages

	loginDeviceCode string
	loginAuthURL    string
	loginPolling    bool
	loginSpinnerIdx int
	loginPollCount  int

	selected     int
	scrollOffset int

	searchMode  bool
	searchQuery string
	detailView  bool
	detailScroll int

	// Status filter for decisions: "" = ALL, "DRAFT", "PENDING", "ACCEPTED", "DEPRECATED"
	statusFilter string

	// Tasks view mode: "list" or "kanban"
	tasksViewMode string

	// Delete confirmation
	confirmDelete   bool
	confirmDeleteID string
	confirmDeleteType string // "memory", "task", "decision"

	// Remember selected index when navigating into sub-views
	menuSelectedIdx int
	projSelectedIdx int
	orgSelectedIdx  int

	inputMode     bool
	inputPrompt   string
	inputValue    string
	inputAction   string
	editingItemID string

	width    int
	height   int
	loading  bool
	errorMsg string

	executeCmd string
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
	user    *config.User
	token   string
	err     error
}

type loginInitMsg struct {
	deviceCode string
	authURL    string
	err        error
}

type loginPollMsg struct {
	status string
	token  string
	user   *config.User
	err    error
}

type loginTickMsg struct{}

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

type chatStreamChunkMsg struct {
	chunk string
}

type chatStreamDoneMsg struct {
	err error
}

type hopperContextLoadedMsg struct {
	decisions    []api.Decision
	memories     []*api.Memory
	chatHistory  []api.ChatHistoryListItem
	err          error
}

type chatHistorySavedMsg struct {
	chatID string
	err    error
}

type chatHistoryLoadedMsg struct {
	sessionID string
	messages  []api.ChatMessage
	err       error
}

// ============================================================================
// INIT & DATA LOADERS
// ============================================================================

func NewInteractiveModel(cfg *config.Config) model {
	isLoggedIn := cfg != nil && cfg.IsAuthenticated()
	m := model{
		cfg:             cfg,
		selected:        0,
		hopperSessionID: fmt.Sprintf("cli-%d", time.Now().UnixNano()),
		tasksViewMode:   "list",
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
	openBrowser(authURL)
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

func (m *model) hopperNewChat() {
	// Save current chat to sessions if it has messages
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

// ============================================================================
// UPDATE
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
// KEY HANDLING
// ============================================================================

func (m model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.errorMsg = ""

	// Handle delete confirmation first
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

	// Hopper chat
	if m.currentView == viewHopper {
		key := msg.String()

		// Sidebar focused
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

		// Chat focused
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

	// Helper: go back from current view, restoring previous selection
	goBack := func() (tea.Model, tea.Cmd) {
		switch m.currentView {
		case viewDashboard, viewDecisions, viewMemories, viewCapsules, viewTasks, viewHopper:
			m.currentView = viewProjectMenu
			m.selected = m.menuSelectedIdx
			m.searchQuery = ""
			m.statusFilter = ""
			m.scrollOffset = 0
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
		// q goes back in feature/menu views, quits at top level
		if m.currentView == viewOrganizations || m.currentView == viewLogin {
			return m, tea.Quit
		}
		return goBack()

	case "esc":
		if m.searchQuery != "" {
			m.searchQuery = ""
			m.selected = 0
			m.scrollOffset = 0
			return m, nil
		}
		if m.currentView == viewOrganizations || m.currentView == viewLogin {
			return m, nil
		}
		return goBack()

	case "up", "k":
		if m.currentView == viewDecisions || m.currentView == viewMemories || m.currentView == viewCapsules {
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
		if m.currentView == viewDecisions || m.currentView == viewMemories || m.currentView == viewCapsules {
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
		if m.currentView == viewDecisions || m.currentView == viewMemories || m.currentView == viewCapsules {
			if m.selected > 0 {
				m.selected--
				cols := m.getGridCols()
				row := m.selected / cols
				if row < m.scrollOffset {
					m.scrollOffset = row
				}
			}
		} else if m.currentView == viewProjects && m.selected > 0 {
			m.selected--
		}

	case "right", "l":
		if m.currentView == viewDecisions || m.currentView == viewMemories || m.currentView == viewCapsules {
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

// handleConfirmDelete handles y/n confirmation for delete
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
		// Any other key cancels
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
		return m, nil
	case tea.KeyEnter:
		m.searchMode = false
		m.selected = 0
		m.scrollOffset = 0
		return m, nil
	case tea.KeyBackspace, tea.KeyDelete:
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			m.selected = 0
			m.scrollOffset = 0
		}
		return m, nil
	case tea.KeySpace:
		m.searchQuery += " "
		m.selected = 0
		m.scrollOffset = 0
		return m, nil
	case tea.KeyRunes:
		m.searchQuery += string(msg.Runes)
		m.selected = 0
		m.scrollOffset = 0
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
			m.projSelectedIdx = m.selected
			m.currentProj = orgProjects[m.selected]
			m.menuItems = []menuItem{
				{"●", "Dashboard", "Project overview & stats", "dashboard"},
				{"●", "Decisions", "View & manage decisions", "decisions"},
				{"●", "Memories", "Project memories & context", "memories"},
				{"●", "Capsules", "Context packs", "capsules"},
				{"●", "Tasks", "Task management", "tasks"},
				{"●", "Hopper", "AI Assistant", "hopper"},
				{"", "", "", ""},
				{"<", "Back", "Return to projects", "back"},
			}
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
				m.loading = true
				return m, m.loadDecisions
			case "memories":
				m.currentView = viewMemories
				m.selected = 0
				m.scrollOffset = 0
				m.searchQuery = ""
				m.loading = true
				return m, m.loadMemories
			case "capsules":
				m.currentView = viewCapsules
				m.selected = 0
				m.scrollOffset = 0
				m.searchQuery = ""
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
			}
		}
	}
	return m, nil
}

// ============================================================================
// SELECTION HELPERS - use filtered lists for CRUD correctness
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

// getVisibleCards returns visible items for card-style views (3 lines per card)
func (m model) getGridCols() int {
	if m.currentView != viewDecisions && m.currentView != viewMemories && m.currentView != viewCapsules {
		return 1
	}
	cardMinWidth := 30
	available := m.width - 6
	if available < cardMinWidth {
		return 1
	}
	cols := available / cardMinWidth
	if cols > 3 {
		cols = 3
	}
	if cols < 1 {
		cols = 1
	}
	return cols
}

func (m model) getCardWidth() int {
	cols := m.getGridCols()
	available := m.width - 6
	if available < 30 {
		available = 30
	}
	gap := 2
	cardW := (available - (cols-1)*gap) / cols
	if cardW < 24 {
		cardW = 24
	}
	return cardW
}

func (m model) getVisibleRows() int {
	cardHeight := 6
	available := m.height - 12
	if available < cardHeight {
		return 1
	}
	rows := available / cardHeight
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

// ============================================================================
// VIEW - Unified layout
// ============================================================================

func (m model) View() string {
	if m.currentView == viewLogin {
		return m.renderLoginView()
	}
	return m.renderPageView()
}

func (m model) renderLoginView() string {
	var s strings.Builder
	border := versionStyle.Render("──────────────────────────────────────────────────────────────────")

	s.WriteString("\n  " + border + "\n\n")

	// Compact logo
	s.WriteString(logoStyle.Render(`    ⢠⣶⣶⣶⣶⣶⣶⣶⣆          ⣴⣶⣶⣶⣶⣶⣶⣶⡄
    ⢸⣿⣯⣷⣿⢿⣾⣟⣿          ⣿⣿⣽⣾⡿⣷⡿⣯⡇
    ⢸⣿⣾⣟⣿⣟⣯⣿⣿          ⣿⣷⣿⣻⡿⣟⣿⣿⡇
    ⠘⢿⣻⣿⣽⣾⣿⣻⣽          ⣿⣷⣿⣿⢿⣷⣿⠿⠃
             ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⡀
             ⣿⣿⣻⣷⣿⣿⣻⣿⣯⣿⡿⣯⣿⣻⣿⣄
    ⢠⣾⣿⣿⢿⣿⡿⣿⣿          ⣿⣽⣿⣽⣷⣿⡿⣟⡇
    ⢸⣿⢿⣻⣾⣿⣻⣿⣻          ⣿⣷⣿⣯⣿⣷⣿⣾⡇
    ⠘⠿⠿⠻⠿⠻⠽⠿⠾          ⠻⠷⠿⠿⠻⠾⠟⠿⠁`) + "\n\n")

	s.WriteString("  " + titleStyle.Render("Hopsule") + "  " + dimStyle.Render("Decision & Memory Layer") + "\n")
	s.WriteString("  " + dimStyle.Render("for AI teams & coding tools") + "  " + versionStyle.Render("v0.8.0") + "\n\n")

	if m.loading && !m.loginPolling {
		s.WriteString("  " + infoStyle.Render("Loading...") + "\n")
	} else if m.errorMsg != "" {
		s.WriteString("  " + lipgloss.NewStyle().Foreground(redColor).Render("Error: "+m.errorMsg) + "\n\n")
	}

	if m.loginPolling {
		spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		s.WriteString("  " + accentStyle.Render(spinner[m.loginSpinnerIdx]) + " Waiting for browser authentication\n\n")
		if m.loginDeviceCode != "" {
			s.WriteString("  Device Code: " + selectedStyle.Render(m.loginDeviceCode) + "\n\n")
		}
		if m.loginAuthURL != "" {
			s.WriteString("  " + dimStyle.Render("If browser didn't open:") + "\n")
			s.WriteString("  " + logoStyle.Render(m.loginAuthURL) + "\n\n")
		}
		s.WriteString("  " + dimStyle.Render(fmt.Sprintf("Polling... %ds", m.loginPollCount*2)) + "\n")
	} else {
		if m.selected == 0 {
			s.WriteString("  " + selectedStyle.Render("> Login") + "  " + infoStyle.Render("(authenticate)") + "\n")
		} else {
			s.WriteString("    " + normalStyle.Render("Login") + "  " + infoStyle.Render("(authenticate)") + "\n")
		}
		s.WriteString("\n  " + dimStyle.Render("Press Enter to continue") + "\n")
	}

	s.WriteString("\n  " + border + "\n")
	s.WriteString("  " + helpStyle.Render("enter login  •  q quit") + "\n")
	return s.String()
}

func (m model) renderPageView() string {
	var s strings.Builder
	border := versionStyle.Render("──────────────────────────────────────────────────────────────────")

	s.WriteString("  " + border + "\n")

	// Breadcrumb header
	s.WriteString(m.renderBreadcrumb())
	s.WriteString("  " + border + "\n\n")

	// Delete confirmation overlay
	if m.confirmDelete {
		action := "delete"
		if m.confirmDeleteType == "decision" {
			action = "deprecate"
		}
		s.WriteString("  " + lipgloss.NewStyle().Foreground(redColor).Bold(true).Render("Are you sure you want to "+action+" this "+m.confirmDeleteType+"?") + "\n\n")
		s.WriteString("  " + dimStyle.Render("Press [y] to confirm, any other key to cancel") + "\n")
	} else if m.inputMode {
		s.WriteString("  " + titleStyle.Render(m.inputPrompt) + "\n\n")
		s.WriteString("  " + m.inputValue + "\u2588\n\n")
		s.WriteString("  " + dimStyle.Render("Enter to confirm  •  Esc to cancel") + "\n")
	} else if m.loading {
		s.WriteString("  " + infoStyle.Render("Loading...") + "\n")
	} else {
		if m.errorMsg != "" {
			s.WriteString("  " + lipgloss.NewStyle().Foreground(redColor).Render("Error: "+m.errorMsg) + "\n")
		}
		// View content
		switch m.currentView {
		case viewOrganizations:
			s.WriteString(m.renderOrganizationsContent())
		case viewProjects:
			s.WriteString(m.renderProjectsContent())
		case viewProjectMenu:
			s.WriteString(m.renderProjectMenuContent())
		case viewDashboard:
			s.WriteString(m.renderDashboardContent())
		case viewDecisions:
			s.WriteString(m.renderDecisionsContent())
		case viewMemories:
			s.WriteString(m.renderMemoriesContent())
		case viewCapsules:
			s.WriteString(m.renderCapsulesContent())
		case viewTasks:
			s.WriteString(m.renderTasksContent())
		case viewHopper:
			s.WriteString(m.renderHopperContent())
		}
	}

	s.WriteString("\n  " + border + "\n")
	s.WriteString(m.renderFooter())
	return s.String()
}

func (m model) renderBreadcrumb() string {
	var parts []string

	switch m.currentView {
	case viewOrganizations:
		parts = []string{"Hopsule", "Organizations"}
	case viewProjects:
		orgName := "..."
		if m.currentOrg != nil {
			orgName = m.currentOrg.Name
		}
		parts = []string{orgName, "Projects"}
	case viewProjectMenu:
		orgName := "..."
		projName := "..."
		if m.currentOrg != nil {
			orgName = m.currentOrg.Name
		}
		if m.currentProj != nil {
			projName = m.currentProj.Name
		}
		parts = []string{orgName, projName}
	default:
		orgName := "..."
		projName := "..."
		pageName := ""
		if m.currentOrg != nil {
			orgName = m.currentOrg.Name
		}
		if m.currentProj != nil {
			projName = m.currentProj.Name
		}
		switch m.currentView {
		case viewDashboard:
			pageName = "Dashboard"
		case viewDecisions:
			pageName = "Decisions"
		case viewMemories:
			pageName = "Memories"
		case viewCapsules:
			pageName = "Capsules"
		case viewTasks:
			pageName = "Tasks"
		case viewHopper:
			pageName = "Hopper"
		}
		parts = []string{orgName, projName, pageName}
	}

	// Build breadcrumb: dim parts except last one bold
	var s strings.Builder
	s.WriteString("  ")
	for i, part := range parts {
		if i > 0 {
			s.WriteString(dimStyle.Render(" / "))
		}
		if i == len(parts)-1 {
			s.WriteString(titleStyle.Render(part))
		} else {
			s.WriteString(dimStyle.Render(part))
		}
	}

	// User info on the right side
	if m.cfg != nil && m.cfg.IsAuthenticated() && m.cfg.User != nil {
		s.WriteString("  " + dimStyle.Render("|") + "  " + statusOnStyle.Render("●") + " " + dimStyle.Render(m.cfg.User.Name))
	}

	s.WriteString("\n")
	return s.String()
}

// ============================================================================
// CONTENT VIEWS
// ============================================================================

func (m model) renderOrganizationsContent() string {
	var s strings.Builder
	if len(m.organizations) == 0 {
		s.WriteString("  " + dimStyle.Render("No organizations found.") + "\n")
		s.WriteString("  " + dimStyle.Render("Create one at app.hopsule.com") + "\n")
	} else {
		for i, org := range m.organizations {
			if i == m.selected {
				s.WriteString("  " + selectedStyle.Render("> "+org.Name) + " " + dimStyle.Render("@"+org.Slug) + "\n")
			} else {
				s.WriteString("    " + normalStyle.Render(org.Name) + " " + dimStyle.Render("@"+org.Slug) + "\n")
			}
		}
		s.WriteString("\n")
		if m.selected == len(m.organizations) {
			s.WriteString("  " + selectedStyle.Render("> Logout") + "\n")
		} else {
			s.WriteString("    " + dimStyle.Render("Logout") + "\n")
		}
	}
	return s.String()
}

func (m model) renderProjectsContent() string {
	var s strings.Builder
	projects := m.getOrgProjects()
	if len(projects) == 0 {
		s.WriteString("  " + dimStyle.Render("No projects found.") + "\n")
		s.WriteString("  " + dimStyle.Render("Create one at app.hopsule.com") + "\n")
		return s.String()
	}
	for i, proj := range projects {
		desc := proj.Description
		if desc == "" {
			desc = "No description"
		}
		if len(desc) > 40 {
			desc = desc[:37] + "..."
		}
		if i == m.selected {
			s.WriteString("  " + selectedStyle.Render("> "+proj.Name) + "  " + dimStyle.Render(desc) + "\n")
		} else {
			s.WriteString("    " + normalStyle.Render(proj.Name) + "  " + dimStyle.Render(desc) + "\n")
		}
	}
	return s.String()
}

func (m model) renderProjectMenuContent() string {
	var s strings.Builder
	for i, item := range m.menuItems {
		if item.name == "" {
			s.WriteString("\n")
			continue
		}
		name := item.name
		for len(name) < 12 {
			name += " "
		}
		if i == m.selected {
			s.WriteString("  " + selectedStyle.Render("> "+item.icon+" "+name) + "  " + dimStyle.Render(item.description) + "\n")
		} else {
			s.WriteString("    " + normalStyle.Render(item.icon+" "+name) + "  " + dimStyle.Render(item.description) + "\n")
		}
	}
	return s.String()
}

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

func (m model) renderDecisionsContent() string {
	var s strings.Builder

	if m.detailView {
		d := m.getSelectedDecision()
		if d != nil {
			return m.renderDecisionDetail(*d)
		}
		m.detailView = false
	}

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

	cols := m.getGridCols()
	cardW := m.getCardWidth()
	visibleRows := m.getVisibleRows()
	totalRows := (len(filtered) + cols - 1) / cols

	if m.scrollOffset > totalRows-visibleRows {
		m.scrollOffset = max(0, totalRows-visibleRows)
	}
	startRow := m.scrollOffset
	endRow := min(startRow+visibleRows, totalRows)

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
		for _, line := range strings.Split(joined, "\n") {
			s.WriteString("  " + line + "\n")
		}
	}

	shown := min(endRow*cols, len(filtered))
	if totalRows > visibleRows {
		s.WriteString("  " + dimStyle.Render(fmt.Sprintf("Showing %d of %d  (scroll with j/k)", shown, len(filtered))) + "\n")
	}
	return s.String()
}

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

	innerW := cardW - 4 // border(2) + space padding(2)
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

func (m model) renderMemoriesContent() string {
	var s strings.Builder

	if m.detailView {
		mem := m.getSelectedMemory()
		if mem != nil {
			return m.renderMemoryDetail(mem)
		}
		m.detailView = false
	}

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

	cols := m.getGridCols()
	cardW := m.getCardWidth()
	visibleRows := m.getVisibleRows()
	totalRows := (len(filtered) + cols - 1) / cols

	if m.scrollOffset > totalRows-visibleRows {
		m.scrollOffset = max(0, totalRows-visibleRows)
	}
	startRow := m.scrollOffset
	endRow := min(startRow+visibleRows, totalRows)

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
		for _, line := range strings.Split(joined, "\n") {
			s.WriteString("  " + line + "\n")
		}
	}

	shown := min(endRow*cols, len(filtered))
	if totalRows > visibleRows {
		s.WriteString("  " + dimStyle.Render(fmt.Sprintf("Showing %d of %d  (scroll with j/k)", shown, len(filtered))) + "\n")
	}
	return s.String()
}

func (m model) renderMemoryCard(mem *api.Memory, idx, cardW int) string {
	innerW := cardW - 4 // border(2) + space padding(2)
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
		// Split by visual width, not rune count, to keep card lines stable.
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

func (m model) renderCapsulesContent() string {
	var s strings.Builder

	if m.detailView {
		c := m.getSelectedCapsule()
		if c != nil {
			return m.renderCapsuleDetail(c)
		}
		m.detailView = false
	}

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

	cols := m.getGridCols()
	cardW := m.getCardWidth()
	visibleRows := m.getVisibleRows()
	totalRows := (len(filtered) + cols - 1) / cols

	if m.scrollOffset > totalRows-visibleRows {
		m.scrollOffset = max(0, totalRows-visibleRows)
	}
	startRow := m.scrollOffset
	endRow := min(startRow+visibleRows, totalRows)

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
		for _, line := range strings.Split(joined, "\n") {
			s.WriteString("  " + line + "\n")
		}
	}

	shown := min(endRow*cols, len(filtered))
	if totalRows > visibleRows {
		s.WriteString("  " + dimStyle.Render(fmt.Sprintf("Showing %d of %d  (scroll with j/k)", shown, len(filtered))) + "\n")
	}
	return s.String()
}

func (m model) renderCapsuleCard(cap *api.Capsule, idx, cardW int) string {
	innerW := cardW - 4 // border(2) + space padding(2)
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

func (m model) renderTasksContent() string {
	var s strings.Builder

	// View mode toggle
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

func (m model) renderHopperContent() string {
	if m.loading && !m.hopperContextLoaded {
		return "  " + dimStyle.Render("Loading project context...") + "\n"
	}

	contentH := m.height - 8
	if contentH < 10 {
		contentH = 10
	}

	// Sidebar width
	sidebarW := 28
	if m.width < 80 {
		sidebarW = 22
	}
	if m.width < 60 {
		sidebarW = 0 // hide sidebar on very small terminals
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

	// Join sidebar and chat panel side by side
	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, chatPanel)
}

func (m model) renderHopperSidebar(w, h int) string {
	sidebarBorderColor := lipgloss.AdaptiveColor{Light: "240", Dark: "237"}
	if m.hopperSidebarFocus {
		sidebarBorderColor = purpleColor
	}

	var sb strings.Builder

	// Header
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

	// Pad to fill height
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

	// Messages area
	var msgLines []string

	if len(m.chatMessages) == 0 {
		// Welcome screen
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
				// User message - right aligned with color
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
				// Parse structured data markers
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

	// Scroll logic
	totalLines := len(msgLines)
	if totalLines > msgAreaH {
		// Auto-scroll to bottom, allow manual scroll up
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

	// Separator
	panel.WriteString(lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "240", Dark: "237"}).Render(strings.Repeat("─", w-4)) + "\n")

	// Input area
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

func (m model) renderHopperAssistantContent(content string, maxW int) string {
	var result strings.Builder
	bubbleW := maxW
	if bubbleW > 72 {
		bubbleW = 72
	}

	// Parse :::DECISIONS, :::MEMORIES, :::CAPSULES, :::TASKS blocks
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

		// Render inline status badges: [STATUS]
		line = renderInlineBadges(line)

		// Format markdown
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

	// Headers
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

	// Code blocks - style with dim
	if strings.HasPrefix(trimmed, "```") {
		return dimStyle.Render("───")
	}

	// Bullet points
	if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
		text := trimmed[2:]
		text = strings.ReplaceAll(text, "**", "")
		text = strings.ReplaceAll(text, "*", "")
		wrapped := wordWrap(text, maxW-4)
		var lines []string
		for i, l := range strings.Split(wrapped, "\n") {
			if i == 0 {
				lines = append(lines, dimStyle.Render("  - ")+normalStyle.Render(l))
			} else {
				lines = append(lines, "    "+normalStyle.Render(l))
			}
		}
		return strings.Join(lines, "\n")
	}

	// Numbered lists
	if len(trimmed) > 2 && trimmed[0] >= '0' && trimmed[0] <= '9' && (trimmed[1] == '.' || (len(trimmed) > 2 && trimmed[1] >= '0' && trimmed[1] <= '9' && trimmed[2] == '.')) {
		text := strings.ReplaceAll(trimmed, "**", "")
		wrapped := wordWrap(text, maxW-2)
		return normalStyle.Render(wrapped)
	}

	// Blockquote
	if strings.HasPrefix(trimmed, "> ") {
		text := strings.TrimPrefix(trimmed, "> ")
		return lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "240", Dark: "242"}).Italic(true).Render("  " + text)
	}

	// Regular text: strip markdown formatting
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

	// Try proper JSON parsing first
	var rawItems []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &rawItems); err == nil {
		var items []hopperBlockItem
		for _, raw := range rawItems {
			item := hopperBlockItem{}
			// Name: try statement > name > title > content
			for _, key := range []string{"statement", "name", "title", "content"} {
				if val, ok := raw[key]; ok {
					if s, ok := val.(string); ok && s != "" {
						item.Name = s
						break
					}
				}
			}
			// Status
			for _, key := range []string{"status", "sourceStatus"} {
				if val, ok := raw[key]; ok {
					if s, ok := val.(string); ok && s != "" {
						item.Status = s
						break
					}
				}
			}
			// Content: try rationale > description > content (if != name)
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

	// Fallback: treat as single text item
	return []hopperBlockItem{{Name: truncateString(jsonData, 80)}}
}

// ============================================================================
// FOOTER
// ============================================================================

func (m model) renderFooter() string {
	if m.confirmDelete {
		action := "delete"
		if m.confirmDeleteType == "decision" {
			action = "deprecate"
		}
		return "  " + helpStyle.Render("[y] confirm "+action+"  •  any key to cancel") + "\n"
	}

	if m.inputMode {
		return "  " + helpStyle.Render("type  •  enter confirm  •  esc cancel") + "\n"
	}

	var help string
	switch m.currentView {
	case viewLogin:
		help = "enter login  •  q quit"
	case viewOrganizations:
		help = "j/k navigate  •  enter select  •  q quit"
	case viewProjects:
		help = "j/k navigate  •  enter select  •  q/esc back"
	case viewProjectMenu:
		help = "j/k navigate  •  enter select  •  q/esc back"
	case viewDashboard:
		help = "q/esc back"
	case viewDecisions:
		if m.detailView {
			help = "[a]ccept  •  [x]deprecate  •  [d]eprec  •  q/esc back"
		} else if m.searchMode {
			help = "type to search  •  enter confirm  •  esc cancel"
		} else {
			help = "hjkl navigate  •  enter detail  •  [/]search  •  tab filter  •  [n]ew  •  [a]ccept  •  [d]eprec  •  q/esc back"
		}
	case viewMemories:
		if m.detailView {
			help = "[e]dit  •  [d]elete  •  q/esc back"
		} else if m.searchMode {
			help = "type to search  •  enter confirm  •  esc cancel"
		} else {
			help = "hjkl navigate  •  enter detail  •  [/]search  •  [n]ew  •  [e]dit  •  [d]elete  •  q/esc back"
		}
	case viewCapsules:
		if m.detailView {
			help = "q/esc back"
		} else if m.searchMode {
			help = "type to search  •  enter confirm  •  esc cancel"
		} else {
			help = "hjkl navigate  •  enter detail  •  [/]search  •  q/esc back"
		}
	case viewTasks:
		if m.searchMode {
			help = "type to search  •  enter confirm  •  esc cancel"
		} else {
			help = "j/k  •  tab view  •  [/]search  •  [n]ew  •  [e]dit  •  [t]oggle  •  [d]elete  •  q/esc back"
		}
	case viewHopper:
		if m.hopperSidebarFocus {
			help = "j/k navigate  •  enter open  •  ctrl+n new chat  •  d delete  •  tab chat  •  esc back"
		} else {
			help = "type message  •  enter send  •  tab history  •  ctrl+n new  •  ctrl+u/d scroll  •  esc back"
		}
	}
	return "  " + helpStyle.Render(help) + "\n"
}

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

func max(a, b int) int {
	if a > b {
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

func wordWrap(s string, width int) string {
	if len(s) <= width {
		return s
	}
	var result strings.Builder
	words := strings.Fields(s)
	lineLen := 0
	for i, word := range words {
		if i > 0 && lineLen+len(word)+1 > width {
			result.WriteString("\n  ")
			lineLen = 0
		} else if i > 0 {
			result.WriteString(" ")
			lineLen++
		}
		result.WriteString(word)
		lineLen += len(word)
	}
	return result.String()
}

func sanitizeMarkdownForTerminal(s string) string {
	// Remove streaming metadata markers
	if idx := strings.Index(s, "__USAGE__"); idx != -1 {
		s = s[:idx]
	}
	if idx := strings.Index(s, "__CONTENT_END__"); idx != -1 {
		s = s[:idx]
	}
	if idx := strings.Index(s, "\"completion_tokens\""); idx != -1 {
		s = s[:idx]
	}

	// Remove progress metadata like {"step":"analyzing","status":"complete"}
	lines := strings.Split(s, "\n")
	var cleaned []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "{\"step\"") || strings.HasPrefix(trimmed, "{\"status\"") {
			continue
		}
		// Also handle inline progress markers
		if strings.Contains(trimmed, "{\"step\":") && strings.Contains(trimmed, "\"status\":") {
			// Remove the JSON part
			if idx := strings.Index(line, "{\"step\":"); idx >= 0 {
				endIdx := strings.Index(line[idx:], "}")
				if endIdx >= 0 {
					line = line[:idx] + line[idx+endIdx+1:]
				}
			}
		}
		cleaned = append(cleaned, line)
	}
	s = strings.Join(cleaned, "\n")

	s = strings.TrimRight(s, " \n\r\t")

	// Strip markdown bold
	s = strings.ReplaceAll(s, "**", "")

	return s
}

func formatMarkdownForTerminal(s string) string {
	lines := strings.Split(s, "\n")
	var result []string
	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			continue
		}
		if strings.HasPrefix(line, "# ") {
			line = strings.TrimPrefix(line, "# ")
		} else if strings.HasPrefix(line, "## ") {
			line = strings.TrimPrefix(line, "## ")
		} else if strings.HasPrefix(line, "### ") {
			line = strings.TrimPrefix(line, "### ")
		}
		if strings.HasPrefix(line, "- ") {
			line = "- " + strings.TrimPrefix(line, "- ")
		}
		if strings.HasPrefix(line, "* ") {
			line = "- " + strings.TrimPrefix(line, "* ")
		}
		line = strings.ReplaceAll(line, "**", "")
		line = strings.ReplaceAll(line, "__", "")
		line = strings.ReplaceAll(line, "`", "")
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

// ============================================================================
// RUN INTERACTIVE
// ============================================================================

func RunInteractive() (string, error) {
	cfg, _ := config.LoadConfig()
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
// EXECUTE LOGIN (standalone, kept for backward compat)
// ============================================================================

func ExecuteLogin(cfg *config.Config) error {
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
	fmt.Println()
	fmt.Println("  Initializing login...")
	initResp, err := client.DeviceAuthInit("CLI")
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
	for attempt := 0; attempt < 300; attempt++ {
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
			fmt.Printf("\r  %s Authentication complete!                    \n\n", statusOnStyle.Render("v"))
			cfg.Token = resp.Token
			cfg.User = &config.User{ID: resp.UserID, Email: resp.Email, Name: resp.Name, AvatarURL: resp.AvatarURL}
			if err := config.SaveConfig(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			fmt.Printf("  Signed in as: %s (%s)\n\n", titleStyle.Render(resp.Name), dimStyle.Render(resp.Email))
			return nil
		case "expired":
			fmt.Println()
			return fmt.Errorf("login session expired - please try again")
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

func ShowOrganizations(cfg *config.Config) {}
func ShowProjects(cfg *config.Config)      {}
