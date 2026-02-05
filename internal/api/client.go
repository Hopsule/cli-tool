package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Cagangedik/cli-tool/internal/config"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		baseURL: cfg.APIURL,
		token:   cfg.Token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) WithToken(token string) *Client {
	return &Client{
		baseURL:    c.baseURL,
		token:      token,
		httpClient: c.httpClient,
	}
}

func (c *Client) WithBaseURL(url string) *Client {
	return &Client{
		baseURL:    url,
		token:      c.token,
		httpClient: c.httpClient,
	}
}

func (c *Client) doRequest(method, path string, body interface{}, projectID string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}
	if projectID != "" {
		req.Header.Set("X-Project-ID", projectID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// ListDecisionsResponse is the response from GET /decisions
type ListDecisionsResponse struct {
	Decisions []Decision `json:"decisions"`
	Total     int        `json:"total"`
}

// ListDecisions lists all decisions for a project
func (c *Client) ListDecisions(projectID string) ([]Decision, error) {
	resp, err := c.doRequest("GET", "/decisions", nil, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var result ListDecisionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Decisions, nil
}

// GetDecision retrieves a specific decision
func (c *Client) GetDecision(projectID, decisionID string) (*Decision, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/decisions/%s", decisionID), nil, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var decision Decision
	if err := json.NewDecoder(resp.Body).Decode(&decision); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &decision, nil
}

// CreateDecision creates a new decision draft
func (c *Client) CreateDecision(projectID string, req CreateDecisionRequest) (*Decision, error) {
	resp, err := c.doRequest("POST", "/decisions/draft", req, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var decision Decision
	if err := json.NewDecoder(resp.Body).Decode(&decision); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &decision, nil
}

// AcceptDecision accepts a decision (moves from DRAFT/PENDING to ACCEPTED)
func (c *Client) AcceptDecision(projectID, decisionID string) (*Decision, error) {
	req := AcceptDecisionRequest{
		ID: decisionID,
	}
	resp, err := c.doRequest("POST", "/decisions/accept", req, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var decision Decision
	if err := json.NewDecoder(resp.Body).Decode(&decision); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &decision, nil
}

// DeprecateDecision deprecates a decision (moves to DEPRECATED)
func (c *Client) DeprecateDecision(projectID, decisionID string) (*Decision, error) {
	req := DeprecateDecisionRequest{
		ID: decisionID,
	}
	resp, err := c.doRequest("POST", "/decisions/deprecate", req, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var decision Decision
	if err := json.NewDecoder(resp.Body).Decode(&decision); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &decision, nil
}

// GetProjectStatus retrieves project status
func (c *Client) GetProjectStatus(projectID string) (*ProjectStatus, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/api/v1/projects/%s/status", projectID), nil, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var status ProjectStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &status, nil
}

// Types matching decision-api schema
type Decision struct {
	ID          string   `json:"id"`
	Statement   string   `json:"statement"`
	Rationale   string   `json:"rationale"`
	Status      string   `json:"status"` // DRAFT, PENDING, ACCEPTED, REJECTED, DEPRECATED
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	AcceptedAt  *string  `json:"accepted_at,omitempty"`
	AcceptedBy  *string  `json:"accepted_by,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type CreateDecisionRequest struct {
	Statement string   `json:"statement"`
	Rationale string   `json:"rationale"`
	ScopeKey  *string  `json:"scope_key,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

type AcceptDecisionRequest struct {
	ID             string `json:"id"`
	AcceptedBy     string `json:"accepted_by,omitempty"`
	AcceptanceNote string `json:"acceptance_note,omitempty"`
}

type DeprecateDecisionRequest struct {
	ID string `json:"id"`
}

type ProjectStatus struct {
	ProjectID      string `json:"project_id"`
	TotalDecisions int    `json:"total_decisions"`
	Accepted       int    `json:"accepted"`
	Pending        int    `json:"pending"`
	Draft          int    `json:"draft"`
	Deprecated     int    `json:"deprecated"`
}

// ============================================================================
// DEVICE AUTH TYPES & METHODS
// ============================================================================

// DeviceAuthInitResponse is returned when initiating device auth
type DeviceAuthInitResponse struct {
	Code      string `json:"code"`
	ExpiresAt string `json:"expires_at"`
	ExpiresIn int    `json:"expires_in"`
}

// DeviceAuthPollResponse is returned when polling for completion
type DeviceAuthPollResponse struct {
	Status    string `json:"status"` // "pending", "complete", "expired"
	UserID    string `json:"user_id,omitempty"`
	Email     string `json:"email,omitempty"`
	Name      string `json:"name,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Token     string `json:"token,omitempty"`
}

// DeviceAuthInit starts the device code flow
func (c *Client) DeviceAuthInit(deviceName string) (*DeviceAuthInitResponse, error) {
	body := map[string]string{"device_name": deviceName}
	resp, err := c.doRequest("POST", "/auth/device/init", body, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var result DeviceAuthInitResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// DeviceAuthPoll polls for device auth completion
func (c *Client) DeviceAuthPoll(code string) (*DeviceAuthPollResponse, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/auth/device/%s/poll", code), nil, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 404 or 410 means code is invalid/expired
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone {
		return &DeviceAuthPollResponse{Status: "expired"}, nil
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var result DeviceAuthPollResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// ============================================================================
// USER & ORGANIZATION TYPES & METHODS
// ============================================================================

// User represents the authenticated user
type User struct {
	ID        string `json:"id"`
	ClerkID   string `json:"clerk_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// Organization represents an organization
type Organization struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// Project represents a project
type Project struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	Description    string `json:"description,omitempty"`
	OrganizationID string `json:"organization_id"`
}

// MeResponse is the response from GET /me
type MeResponse struct {
	User          *User           `json:"user"`
	Organizations []*Organization `json:"organizations"`
	Projects      []*Project      `json:"projects"`
}

// GetMe retrieves the current user's info
func (c *Client) GetMe() (*MeResponse, error) {
	resp, err := c.doRequest("GET", "/me", nil, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var result MeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// ListOrganizations lists the user's organizations
func (c *Client) ListOrganizations() ([]*Organization, error) {
	resp, err := c.doRequest("GET", "/organizations", nil, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var result []*Organization
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// ListProjects lists the user's projects
func (c *Client) ListProjects() ([]*Project, error) {
	resp, err := c.doRequest("GET", "/projects", nil, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var result []*Project
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// ============================================================================
// MEMORY TYPES & METHODS
// ============================================================================

// Memory represents a project memory
type Memory struct {
	ID                 string   `json:"id"`
	Content            string   `json:"content"`
	CreatedAt          string   `json:"created_at"`
	UpdatedAt          string   `json:"updated_at,omitempty"`
	Tags               []string `json:"tags,omitempty"`
	RelatedDecisionIds []string `json:"related_decision_ids,omitempty"`
	CreatedByUserID    string   `json:"created_by_user_id,omitempty"`
	CreatedByName      string   `json:"created_by_name,omitempty"`
}

// CreateMemoryRequest is the request body for creating a memory
type CreateMemoryRequest struct {
	Content            string   `json:"content"`
	Tags               []string `json:"tags,omitempty"`
	RelatedDecisionIds []string `json:"related_decision_ids,omitempty"`
}

// UpdateMemoryRequest is the request body for updating a memory
type UpdateMemoryRequest struct {
	Content            string   `json:"content,omitempty"`
	Tags               []string `json:"tags,omitempty"`
	RelatedDecisionIds []string `json:"related_decision_ids,omitempty"`
}

// ListMemoriesResponse is the response from GET /memories
type ListMemoriesResponse struct {
	Memories []*Memory `json:"memories"`
	Total    int       `json:"total"`
}

// ListMemories lists all memories for a project
func (c *Client) ListMemories(projectID string) ([]*Memory, error) {
	resp, err := c.doRequest("GET", "/memories?limit=100", nil, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var result ListMemoriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Memories, nil
}

// CreateMemory creates a new memory
func (c *Client) CreateMemory(projectID string, req CreateMemoryRequest) (*Memory, error) {
	resp, err := c.doRequest("POST", "/memories", req, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var result struct {
		Memory *Memory `json:"memory"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Memory, nil
}

// UpdateMemory updates a memory
func (c *Client) UpdateMemory(projectID, memoryID string, req UpdateMemoryRequest) (*Memory, error) {
	resp, err := c.doRequest("PATCH", fmt.Sprintf("/memories/%s", memoryID), req, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var memory Memory
	if err := json.NewDecoder(resp.Body).Decode(&memory); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &memory, nil
}

// DeleteMemory deletes a memory
func (c *Client) DeleteMemory(projectID, memoryID string) error {
	resp, err := c.doRequest("DELETE", fmt.Sprintf("/memories/%s", memoryID), nil, projectID)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

// ============================================================================
// TASK TYPES & METHODS
// ============================================================================

// Task represents a project task
type Task struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Status      string   `json:"status"` // TODO, IN_PROGRESS, REVIEW, DONE
	Priority    string   `json:"priority"` // LOW, MEDIUM, HIGH
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	CompletedAt *string  `json:"completed_at,omitempty"`
	OwnerID     string   `json:"owner_id,omitempty"`
	OwnerName   string   `json:"owner_name,omitempty"`
	RelatedDecisionIds []string `json:"related_decision_ids,omitempty"`
	RelatedMemoryIds   []string `json:"related_memory_ids,omitempty"`
}

// CreateTaskRequest is the request body for creating a task
type CreateTaskRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Priority    string   `json:"priority,omitempty"`
	RelatedDecisionIds []string `json:"related_decision_ids,omitempty"`
	RelatedMemoryIds   []string `json:"related_memory_ids,omitempty"`
}

// UpdateTaskRequest is the request body for updating a task
type UpdateTaskRequest struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Status      string   `json:"status,omitempty"`
	Priority    string   `json:"priority,omitempty"`
}

// ListTasksResponse is the response from GET /tasks
type ListTasksResponse struct {
	Tasks []*Task `json:"tasks"`
	Total int     `json:"total"`
}

// ListTasks lists all tasks for a project
func (c *Client) ListTasks(projectID string) ([]*Task, error) {
	resp, err := c.doRequest("GET", "/tasks", nil, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var result ListTasksResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Tasks, nil
}

// CreateTask creates a new task
func (c *Client) CreateTask(projectID string, req CreateTaskRequest) (*Task, error) {
	resp, err := c.doRequest("POST", "/tasks", req, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &task, nil
}

// UpdateTask updates a task
func (c *Client) UpdateTask(projectID, taskID string, req UpdateTaskRequest) (*Task, error) {
	resp, err := c.doRequest("PUT", fmt.Sprintf("/tasks/%s", taskID), req, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &task, nil
}

// DeleteTask deletes a task
func (c *Client) DeleteTask(projectID, taskID string) error {
	resp, err := c.doRequest("DELETE", fmt.Sprintf("/tasks/%s", taskID), nil, projectID)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

// ============================================================================
// CAPSULE TYPES & METHODS
// ============================================================================

// Capsule represents a context pack
type Capsule struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Status      string   `json:"status"` // DRAFT, FROZEN, HISTORICAL
	DecisionIds []string `json:"decision_ids,omitempty"`
	MemoryIds   []string `json:"memory_ids,omitempty"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	FrozenAt    *string  `json:"frozen_at,omitempty"`
	IsActive    bool     `json:"is_active,omitempty"`
}

// ListCapsulesResponse is the response from GET /capsules
type ListCapsulesResponse struct {
	Capsules []*Capsule `json:"capsules"`
	Total    int        `json:"total"`
}

// ListCapsules lists all capsules for a project
func (c *Client) ListCapsules(projectID string) ([]*Capsule, error) {
	resp, err := c.doRequest("GET", "/capsules", nil, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var result ListCapsulesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Capsules, nil
}

// ============================================================================
// GRAPH TYPES & METHODS
// ============================================================================

// GraphStats represents graph statistics
type GraphStats struct {
	NodeCount   int            `json:"nodeCount"`
	EdgeCount   int            `json:"edgeCount"`
	NodesByType map[string]int `json:"nodesByType"`
	EdgesByType map[string]int `json:"edgesByType,omitempty"`
}

// GetGraphStats retrieves graph statistics for a project
func (c *Client) GetGraphStats(projectID string) (*GraphStats, error) {
	resp, err := c.doRequest("GET", "/graph/stats", nil, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var stats GraphStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &stats, nil
}

// ============================================================================
// HOPPER AI CHAT TYPES & METHODS
// ============================================================================

// ChatMessage represents a message in the conversation
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// TaggedItem represents a decision, memory, capsule or task to include as context
type TaggedItem struct {
	ID          string `json:"id"`
	Type        string `json:"type"` // "decision", "memory", "capsule", "task"
	Name        string `json:"name,omitempty"`
	Content     string `json:"content,omitempty"`
	Statement   string `json:"statement,omitempty"`
	Description string `json:"description,omitempty"`
}

// ChatRequest is the request body for Hopper chat
type ChatRequest struct {
	Message             string        `json:"message"`
	ConversationHistory []ChatMessage `json:"conversationHistory,omitempty"`
	TaggedItems         []TaggedItem  `json:"taggedItems,omitempty"`
	Stream              bool          `json:"stream"`
	ProjectName         string        `json:"projectName,omitempty"`
}

// SendChatMessage sends a message to Hopper and streams the response
// The callback is called with each chunk of the response
func (c *Client) SendChatMessage(projectID string, req *ChatRequest, onChunk func(string)) error {
	// Create longer timeout client for AI chat
	chatClient := &http.Client{
		Timeout: 120 * time.Second,
	}

	url := fmt.Sprintf("%s/ai/hopper/chat", c.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	if c.token != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}
	if projectID != "" {
		httpReq.Header.Set("X-Project-Id", projectID)
	}

	resp, err := chatClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("chat request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("chat API error: %d - %s", resp.StatusCode, string(body))
	}

	// Read streaming response
	buf := make([]byte, 128)
	inContent := false
	
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			chunk := string(buf[:n])
			
			// Check for content start marker
			if !inContent {
				if idx := findIndex(chunk, "__CONTENT_START__"); idx != -1 {
					inContent = true
					chunk = chunk[idx+len("__CONTENT_START__"):]
				}
			}
			
			if inContent {
				// Check for usage marker (end of content)
				if idx := findIndex(chunk, "__USAGE__"); idx != -1 {
					chunk = chunk[:idx]
				}
				
				// Skip progress markers
				if !containsMarker(chunk, "__PROGRESS__") && !containsMarker(chunk, "__END_PROGRESS__") {
					if chunk != "" {
						onChunk(chunk)
					}
				}
			}
		}
		
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading response: %w", err)
		}
	}

	return nil
}

// Helper function to find index of substring
func findIndex(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// Helper function to check if string contains marker
func containsMarker(s, marker string) bool {
	return findIndex(s, marker) != -1
}
