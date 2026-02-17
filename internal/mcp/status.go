package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerStatusTools(server *gomcp.Server, sctx *ServerContext) {
	server.AddTool(getProjectStatusTool(), getProjectStatusHandler(sctx))
}

// ============================================================================
// get-project-status
// ============================================================================

func getProjectStatusTool() *gomcp.Tool {
	return &gomcp.Tool{
		Name:        "get-project-status",
		Description: "Get a comprehensive overview of the current project: total decisions, accepted decisions, memories, tasks, active context, and conflicts.",
		InputSchema: inputSchema(map[string]interface{}{}, nil),
	}
}

type projectStatus struct {
	ProjectID         string `json:"project_id"`
	TotalDecisions    int    `json:"total_decisions"`
	AcceptedDecisions int    `json:"accepted_decisions"`
	TotalMemories     int    `json:"total_memories"`
	TotalTasks        int    `json:"total_tasks"`
	ActiveContext     string `json:"active_context"`
	ActivePackName    string `json:"active_pack_name,omitempty"`
	TotalConflicts    int    `json:"total_conflicts"`
}

func getProjectStatusHandler(sctx *ServerContext) gomcp.ToolHandler {
	return func(ctx context.Context, req *gomcp.CallToolRequest) (*gomcp.CallToolResult, error) {
		status := &projectStatus{
			ProjectID: sctx.ProjectID,
		}

		var wg sync.WaitGroup
		var mu sync.Mutex
		var errs []string

		wg.Add(1)
		go func() {
			defer wg.Done()
			decisions, err := sctx.Client.ListDecisions(sctx.ProjectID)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errs = append(errs, fmt.Sprintf("decisions: %v", err))
				return
			}
			status.TotalDecisions = len(decisions)
			for _, d := range decisions {
				if d.Status == "ACCEPTED" {
					status.AcceptedDecisions++
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			memories, err := sctx.Client.ListMemories(sctx.ProjectID)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errs = append(errs, fmt.Sprintf("memories: %v", err))
				return
			}
			status.TotalMemories = len(memories)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			tasks, err := sctx.Client.ListTasks(sctx.ProjectID)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errs = append(errs, fmt.Sprintf("tasks: %v", err))
				return
			}
			status.TotalTasks = len(tasks)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			activeCtx, err := sctx.Client.GetActiveContext(sctx.ProjectID)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				status.ActiveContext = "unknown"
				errs = append(errs, fmt.Sprintf("active context: %v", err))
				return
			}
			status.ActiveContext = activeCtx.Status
			if activeCtx.Pack != nil {
				status.ActivePackName = activeCtx.Pack.Name
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := sctx.Client.DoRawRequest("GET", "/conflicts", nil, sctx.ProjectID)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errs = append(errs, fmt.Sprintf("conflicts: %v", err))
				return
			}
			defer resp.Body.Close()
			var conflicts []json.RawMessage
			if err := json.NewDecoder(resp.Body).Decode(&conflicts); err != nil {
				errs = append(errs, fmt.Sprintf("conflicts decode: %v", err))
				return
			}
			status.TotalConflicts = len(conflicts)
		}()

		wg.Wait()

		data, _ := json.MarshalIndent(status, "", "  ")
		summary := string(data)
		if len(errs) > 0 {
			summary += fmt.Sprintf("\n\n(Partial results due to errors: %v)", errs)
		}
		return toolResult(summary), nil
	}
}
