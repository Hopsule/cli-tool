package ui

import (
	"github.com/Hopsule/cli-tool/internal/api"
	"github.com/Hopsule/cli-tool/internal/config"
)

// ============================================================================
// TEA MESSAGES — all message types used by the Bubble Tea update loop
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
	decisions   []api.Decision
	memories    []*api.Memory
	chatHistory []api.ChatHistoryListItem
	err         error
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
