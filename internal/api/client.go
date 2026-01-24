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

func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
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

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// ListDecisions lists all decisions for a project
func (c *Client) ListDecisions(projectID string) ([]Decision, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/api/v1/projects/%s/decisions", projectID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var decisions []Decision
	if err := json.NewDecoder(resp.Body).Decode(&decisions); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return decisions, nil
}

// GetDecision retrieves a specific decision
func (c *Client) GetDecision(projectID, decisionID string) (*Decision, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/api/v1/projects/%s/decisions/%s", projectID, decisionID), nil)
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

// CreateDecision creates a new decision
func (c *Client) CreateDecision(projectID string, req CreateDecisionRequest) (*Decision, error) {
	resp, err := c.doRequest("POST", fmt.Sprintf("/api/v1/projects/%s/decisions", projectID), req)
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
	resp, err := c.doRequest("POST", fmt.Sprintf("/api/v1/projects/%s/decisions/%s/accept", projectID, decisionID), nil)
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
	resp, err := c.doRequest("POST", fmt.Sprintf("/api/v1/projects/%s/decisions/%s/deprecate", projectID, decisionID), nil)
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
	resp, err := c.doRequest("GET", fmt.Sprintf("/api/v1/projects/%s/status", projectID), nil)
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

type ProjectStatus struct {
	ProjectID      string `json:"project_id"`
	TotalDecisions int    `json:"total_decisions"`
	Accepted       int    `json:"accepted"`
	Pending        int    `json:"pending"`
	Draft          int    `json:"draft"`
	Deprecated     int    `json:"deprecated"`
}
