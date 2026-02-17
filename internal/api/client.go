package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Hopsule/cli-tool/internal/config"
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
func (c *Client) AcceptDecision(projectID, decisionID, acceptedBy string) (*Decision, error) {
	if acceptedBy == "" {
		acceptedBy = "CLI User"
	}
	req := AcceptDecisionRequest{
		ID:         decisionID,
		AcceptedBy: acceptedBy,
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

// ListAcceptedDecisions lists only accepted decisions for a project
func (c *Client) ListAcceptedDecisions(projectID string) ([]Decision, error) {
	resp, err := c.doRequest("GET", "/decisions/accepted", nil, projectID)
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

// DoRawRequest exposes doRequest for MCP tools that need low-level access
func (c *Client) DoRawRequest(method, path string, body interface{}, projectID string) (*http.Response, error) {
	return c.doRequest(method, path, body, projectID)
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
	ID         string   `json:"id"`
	Statement  string   `json:"statement"`
	Rationale  string   `json:"rationale"`
	Status     string   `json:"status"` // DRAFT, PENDING, ACCEPTED, REJECTED, DEPRECATED
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
	AcceptedAt *string  `json:"accepted_at,omitempty"`
	AcceptedBy *string  `json:"accepted_by,omitempty"`
	Tags       []string `json:"tags,omitempty"`
}

type CreateDecisionRequest struct {
	Statement string   `json:"statement"`
	Rationale string   `json:"rationale"`
	ScopeKey  *string  `json:"scope_key,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

type AcceptDecisionRequest struct {
	ID             string `json:"id"`
	AcceptedBy     string `json:"accepted_by"`
	AcceptanceNote string `json:"acceptance_note"`
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

// GetMemory retrieves a single memory by ID
func (c *Client) GetMemory(projectID, memoryID string) (*Memory, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/memories/%s", memoryID), nil, projectID)
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

// ============================================================================
// SEARCH TYPES & METHODS
// ============================================================================

// SearchRequest is the request body for POST /search
type SearchRequest struct {
	Query          string   `json:"query"`
	Mode           string   `json:"mode,omitempty"`
	EntityTypes    []string `json:"entity_types,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	MinSimilarity  float64  `json:"min_similarity,omitempty"`
	Limit          int      `json:"limit,omitempty"`
	IncludeContent bool     `json:"include_content,omitempty"`
	UseReranking   bool     `json:"use_reranking,omitempty"`
}

// SearchResult represents a single search result
type SearchResult struct {
	ID         string   `json:"id"`
	EntityType string   `json:"entity_type"`
	Content    string   `json:"content,omitempty"`
	Statement  string   `json:"statement,omitempty"`
	Status     string   `json:"status,omitempty"`
	Score      float64  `json:"score,omitempty"`
	Tags       []string `json:"tags,omitempty"`
}

// SearchResponse is the response from POST /search
type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
}

// Search performs a semantic search across entities
func (c *Client) Search(projectID string, req SearchRequest) (*SearchResponse, error) {
	resp, err := c.doRequest("POST", "/search", req, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var result SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// ============================================================================
// CAPSULE EXTENDED METHODS
// ============================================================================

// MaterializedCapsule represents a capsule with full decision and memory objects
type MaterializedCapsule struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Status      string     `json:"status"`
	Decisions   []Decision `json:"decisions,omitempty"`
	Memories    []*Memory  `json:"memories,omitempty"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	FrozenAt    *string    `json:"frozen_at,omitempty"`
	IsActive    bool       `json:"is_active,omitempty"`
}

// GetCapsuleMaterialized retrieves a capsule with full decisions and memories
func (c *Client) GetCapsuleMaterialized(projectID, capsuleID string) (*MaterializedCapsule, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/capsules/%s/materialize", capsuleID), nil, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var result MaterializedCapsule
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// ActiveContextResponse is the response from GET /capsules/active
type ActiveContextResponse struct {
	Status string   `json:"status"` // "none", "active", "ambiguous"
	Pack   *Capsule `json:"pack,omitempty"`
	Count  *int     `json:"count,omitempty"`
}

// GetActiveContext retrieves the active context pack for a project
func (c *Client) GetActiveContext(projectID string) (*ActiveContextResponse, error) {
	resp, err := c.doRequest("GET", "/capsules/active", nil, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var result ActiveContextResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetTask retrieves a single task by ID
func (c *Client) GetTask(projectID, taskID string) (*Task, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/tasks/%s", taskID), nil, projectID)
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

// ============================================================================
// TASK TYPES & METHODS
// ============================================================================

// Task represents a project task
type Task struct {
	ID                 string   `json:"id"`
	Title              string   `json:"title"`
	Description        string   `json:"description,omitempty"`
	Status             string   `json:"status"`   // TODO, IN_PROGRESS, REVIEW, DONE
	Priority           string   `json:"priority"` // LOW, MEDIUM, HIGH
	CreatedAt          string   `json:"created_at"`
	UpdatedAt          string   `json:"updated_at"`
	CompletedAt        *string  `json:"completed_at,omitempty"`
	OwnerID            string   `json:"owner_id,omitempty"`
	OwnerName          string   `json:"owner_name,omitempty"`
	RelatedDecisionIds []string `json:"related_decision_ids,omitempty"`
	RelatedMemoryIds   []string `json:"related_memory_ids,omitempty"`
}

// CreateTaskRequest is the request body for creating a task
type CreateTaskRequest struct {
	Title              string   `json:"title"`
	Description        string   `json:"description,omitempty"`
	Priority           string   `json:"priority,omitempty"`
	RelatedDecisionIds []string `json:"related_decision_ids,omitempty"`
	RelatedMemoryIds   []string `json:"related_memory_ids,omitempty"`
}

// UpdateTaskRequest is the request body for updating a task
type UpdateTaskRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	Priority    string `json:"priority,omitempty"`
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
	SessionID           string        `json:"sessionId,omitempty"`
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

	// Read streaming response with larger buffer and accumulator for split markers
	buf := make([]byte, 4096)
	inContent := false
	var pending string

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			pending += string(buf[:n])

			// Check for content start marker
			if !inContent {
				if idx := findIndex(pending, "__CONTENT_START__"); idx != -1 {
					inContent = true
					pending = pending[idx+len("__CONTENT_START__"):]
				} else {
					// Could be partial marker, keep accumulating
					if len(pending) > 100 {
						pending = pending[len(pending)-30:]
					}
					if err != nil {
						break
					}
					continue
				}
			}

			if inContent {
				// Check for usage/end markers
				endMarkers := []string{"__USAGE__", "__CONTENT_END__", "\"completion_tokens\""}
				contentEnd := -1
				for _, marker := range endMarkers {
					if idx := findIndex(pending, marker); idx != -1 {
						if contentEnd == -1 || idx < contentEnd {
							contentEnd = idx
						}
					}
				}

				if contentEnd != -1 {
					toSend := pending[:contentEnd]
					// Remove progress markers
					toSend = removeMarkers(toSend)
					if toSend != "" {
						onChunk(toSend)
					}
					break
				}

				// Check for partial markers at the end - keep last 30 chars as pending
				safeEnd := len(pending) - 30
				if safeEnd > 0 {
					toSend := pending[:safeEnd]
					pending = pending[safeEnd:]
					toSend = removeMarkers(toSend)
					if toSend != "" {
						onChunk(toSend)
					}
				}
			}
		}

		if err != nil {
			if err == io.EOF {
				// Flush remaining content
				if inContent && pending != "" {
					// Clean any trailing metadata
					for _, marker := range []string{"__USAGE__", "__CONTENT_END__", "\"completion_tokens\""} {
						if idx := findIndex(pending, marker); idx != -1 {
							pending = pending[:idx]
						}
					}
					pending = removeMarkers(pending)
					if pending != "" {
						onChunk(pending)
					}
				}
				break
			}
			return fmt.Errorf("error reading response: %w", err)
		}
	}

	return nil
}

func removeMarkers(s string) string {
	result := s
	markers := []string{"__PROGRESS__", "__END_PROGRESS__"}
	for _, m := range markers {
		result = strings.ReplaceAll(result, m, "")
	}
	return result
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

// ============================================================================
// CHAT HISTORY TYPES & METHODS
// ============================================================================

// ChatHistoryEntry represents a single chat history item (full)
type ChatHistoryEntry struct {
	ID        string        `json:"id"`
	ProjectID string        `json:"project_id"`
	UserID    string        `json:"user_id"`
	Topic     string        `json:"topic"`
	Messages  []ChatMessage `json:"messages"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
}

// ChatHistoryListItem is a lightweight version for listing
type ChatHistoryListItem struct {
	ID        string `json:"id"`
	Topic     string `json:"topic"`
	UpdatedAt string `json:"updated_at"`
}

// ListChatHistoryResponse is the response from GET /chat-history
type ListChatHistoryResponse struct {
	ChatHistory []ChatHistoryListItem `json:"chatHistory"`
	Total       int                   `json:"total"`
}

// CreateChatHistoryRequest is the request body for creating chat history
type CreateChatHistoryRequest struct {
	Topic    string        `json:"topic"`
	Messages []ChatMessage `json:"messages"`
}

// UpdateChatHistoryRequest is the request body for updating chat history
type UpdateChatHistoryRequest struct {
	Topic    string        `json:"topic,omitempty"`
	Messages []ChatMessage `json:"messages,omitempty"`
}

// ListChatHistory retrieves chat history list for the current project
func (c *Client) ListChatHistory(projectID string, limit int) ([]ChatHistoryListItem, error) {
	path := fmt.Sprintf("/chat-history?limit=%d", limit)
	resp, err := c.doRequest("GET", path, nil, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var result ListChatHistoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.ChatHistory, nil
}

// GetChatHistory retrieves a single chat history by ID
func (c *Client) GetChatHistory(projectID, chatID string) (*ChatHistoryEntry, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/chat-history/%s", chatID), nil, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var entry ChatHistoryEntry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &entry, nil
}

// CreateChatHistory creates a new chat history entry
func (c *Client) CreateChatHistory(projectID string, req *CreateChatHistoryRequest) (*ChatHistoryEntry, error) {
	resp, err := c.doRequest("POST", "/chat-history", req, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var entry ChatHistoryEntry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &entry, nil
}

// UpdateChatHistory updates an existing chat history entry
func (c *Client) UpdateChatHistory(projectID, chatID string, req *UpdateChatHistoryRequest) (*ChatHistoryEntry, error) {
	resp, err := c.doRequest("PUT", fmt.Sprintf("/chat-history/%s", chatID), req, projectID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var entry ChatHistoryEntry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &entry, nil
}
