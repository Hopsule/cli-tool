package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Hopsule/cli-tool/internal/api"
	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerTaskTools(server *gomcp.Server, sctx *ServerContext) {
	server.AddTool(listTasksTool(), listTasksHandler(sctx))
	server.AddTool(createTaskTool(), createTaskHandler(sctx))
	server.AddTool(updateTaskTool(), updateTaskHandler(sctx))
}

// ============================================================================
// list-tasks
// ============================================================================

func listTasksTool() *gomcp.Tool {
	return &gomcp.Tool{
		Name:        "list-tasks",
		Description: "List tasks in the current project. Tasks track actionable work items derived from decisions and context.",
		InputSchema: inputSchema(map[string]interface{}{}, nil),
	}
}

func listTasksHandler(sctx *ServerContext) gomcp.ToolHandler {
	return func(ctx context.Context, req *gomcp.CallToolRequest) (*gomcp.CallToolResult, error) {
		tasks, err := sctx.Client.ListTasks(sctx.ProjectID)
		if err != nil {
			return toolError(fmt.Sprintf("Failed to list tasks: %v", err)), nil
		}

		data, _ := json.MarshalIndent(tasks, "", "  ")
		return toolResult(string(data)), nil
	}
}

// ============================================================================
// create-task (write delegation)
// ============================================================================

func createTaskTool() *gomcp.Tool {
	return &gomcp.Tool{
		Name:        "create-task",
		Description: "[WRITE — delegates to decision-api] Create a new task. The decision-api is the single authority for this mutation.",
		InputSchema: inputSchema(map[string]interface{}{
			"title":                prop("string", "Task title"),
			"description":         prop("string", "Task description with details"),
			"priority":            prop("string", "Priority: 'low', 'medium', 'high', 'critical'"),
			"related_decision_ids": arrayProp("Decision IDs related to this task"),
			"related_memory_ids":   arrayProp("Memory IDs related to this task"),
		}, []string{"title"}),
	}
}

func createTaskHandler(sctx *ServerContext) gomcp.ToolHandler {
	return func(ctx context.Context, req *gomcp.CallToolRequest) (*gomcp.CallToolResult, error) {
		args := parseArgs(req.Params.Arguments)
		title := getStringArg(args, "title")
		if title == "" {
			return toolError("'title' is required"), nil
		}

		createReq := api.CreateTaskRequest{
			Title:              title,
			Description:        getStringArg(args, "description"),
			Priority:           getStringArg(args, "priority"),
			RelatedDecisionIds: getStringSliceArg(args, "related_decision_ids"),
			RelatedMemoryIds:   getStringSliceArg(args, "related_memory_ids"),
		}

		task, err := sctx.Client.CreateTask(sctx.ProjectID, createReq)
		if err != nil {
			return toolError(fmt.Sprintf("Failed to create task: %v", err)), nil
		}

		data, _ := json.MarshalIndent(task, "", "  ")
		return toolResult(fmt.Sprintf("Task created.\n\n%s", string(data))), nil
	}
}

// ============================================================================
// update-task (write delegation)
// ============================================================================

func updateTaskTool() *gomcp.Tool {
	return &gomcp.Tool{
		Name:        "update-task",
		Description: "[WRITE — delegates to decision-api] Update an existing task (status, description, priority). The decision-api is the single authority.",
		InputSchema: inputSchema(map[string]interface{}{
			"id":          prop("string", "The task ID to update"),
			"title":       prop("string", "New task title"),
			"description": prop("string", "New task description"),
			"status":      prop("string", "New status (e.g. 'open', 'in_progress', 'done')"),
			"priority":    prop("string", "New priority: 'low', 'medium', 'high', 'critical'"),
		}, []string{"id"}),
	}
}

func updateTaskHandler(sctx *ServerContext) gomcp.ToolHandler {
	return func(ctx context.Context, req *gomcp.CallToolRequest) (*gomcp.CallToolResult, error) {
		args := parseArgs(req.Params.Arguments)
		id := getStringArg(args, "id")
		if id == "" {
			return toolError("'id' is required"), nil
		}

		updateReq := api.UpdateTaskRequest{
			Title:       getStringArg(args, "title"),
			Description: getStringArg(args, "description"),
			Status:      getStringArg(args, "status"),
			Priority:    getStringArg(args, "priority"),
		}

		task, err := sctx.Client.UpdateTask(sctx.ProjectID, id, updateReq)
		if err != nil {
			return toolError(fmt.Sprintf("Failed to update task: %v", err)), nil
		}

		data, _ := json.MarshalIndent(task, "", "  ")
		return toolResult(fmt.Sprintf("Task updated.\n\n%s", string(data))), nil
	}
}
