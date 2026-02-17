package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Hopsule/cli-tool/internal/api"
	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerContextTools(server *gomcp.Server, sctx *ServerContext) {
	server.AddTool(getContextPackTool(), getContextPackHandler(sctx))
	server.AddTool(getProjectContextTool(), getProjectContextHandler(sctx))
	server.AddTool(getRelevantContextTool(), getRelevantContextHandler(sctx))
}

// ============================================================================
// get-context-pack
// ============================================================================

func getContextPackTool() *gomcp.Tool {
	return &gomcp.Tool{
		Name:        "get-context-pack",
		Description: "Get a context pack (capsule) by ID, materialized with full decisions and memories. Context packs are portable bundles of decisions + memories.",
		InputSchema: inputSchema(map[string]interface{}{
			"id": prop("string", "The context pack (capsule) ID"),
		}, []string{"id"}),
	}
}

func getContextPackHandler(sctx *ServerContext) gomcp.ToolHandler {
	return func(ctx context.Context, req *gomcp.CallToolRequest) (*gomcp.CallToolResult, error) {
		args := parseArgs(req.Params.Arguments)
		id := getStringArg(args, "id")
		if id == "" {
			return toolError("'id' is required"), nil
		}

		capsule, err := sctx.Client.GetCapsuleMaterialized(sctx.ProjectID, id)
		if err != nil {
			return toolError(fmt.Sprintf("Failed to get context pack: %v", err)), nil
		}

		data, _ := json.MarshalIndent(capsule, "", "  ")
		return toolResult(string(data)), nil
	}
}

// ============================================================================
// get-project-context
// ============================================================================

func getProjectContextTool() *gomcp.Tool {
	return &gomcp.Tool{
		Name:        "get-project-context",
		Description: "Get the active context for the current project. Returns the currently active context pack if one is set.",
		InputSchema: inputSchema(map[string]interface{}{}, nil),
	}
}

func getProjectContextHandler(sctx *ServerContext) gomcp.ToolHandler {
	return func(ctx context.Context, req *gomcp.CallToolRequest) (*gomcp.CallToolResult, error) {
		activeCtx, err := sctx.Client.GetActiveContext(sctx.ProjectID)
		if err != nil {
			return toolError(fmt.Sprintf("Failed to get active context: %v", err)), nil
		}

		switch activeCtx.Status {
		case "none":
			return toolResult("No active context pack for this project. Use the web UI or CLI to set one."), nil
		case "ambiguous":
			count := 0
			if activeCtx.Count != nil {
				count = *activeCtx.Count
			}
			return toolResult(fmt.Sprintf("Ambiguous: %d context packs found. Use the web UI to select one active pack.", count)), nil
		case "active":
			if activeCtx.Pack == nil {
				return toolResult("Active context pack exists but details are unavailable."), nil
			}
			materialized, err := sctx.Client.GetCapsuleMaterialized(sctx.ProjectID, activeCtx.Pack.ID)
			if err != nil {
				data, _ := json.MarshalIndent(activeCtx.Pack, "", "  ")
				return toolResult(string(data)), nil
			}
			data, _ := json.MarshalIndent(materialized, "", "  ")
			return toolResult(string(data)), nil
		default:
			data, _ := json.MarshalIndent(activeCtx, "", "  ")
			return toolResult(string(data)), nil
		}
	}
}

// ============================================================================
// get-relevant-context (smart tool)
// ============================================================================

func getRelevantContextTool() *gomcp.Tool {
	return &gomcp.Tool{
		Name:        "get-relevant-context",
		Description: "SMART TOOL: Get relevant decisions, memories, and tasks for a query. Performs parallel semantic searches across all entity types. Use this as the primary entry point to understand project context before making changes.",
		InputSchema: inputSchema(map[string]interface{}{
			"query":     prop("string", "What you want to understand — natural language query about the codebase, feature, or decision area"),
			"file_path": prop("string", "Optional file path to contextualize the search (e.g. 'src/auth/middleware.ts')"),
			"limit":     prop("number", "Maximum results per entity type (default: 5)"),
		}, []string{"query"}),
	}
}

type relevantContextResult struct {
	Query     string              `json:"query"`
	FilePath  string              `json:"file_path,omitempty"`
	Decisions *api.SearchResponse `json:"decisions,omitempty"`
	Memories  *api.SearchResponse `json:"memories,omitempty"`
	Tasks     []*api.Task         `json:"tasks,omitempty"`
	Summary   string              `json:"summary"`
}

func getRelevantContextHandler(sctx *ServerContext) gomcp.ToolHandler {
	return func(ctx context.Context, req *gomcp.CallToolRequest) (*gomcp.CallToolResult, error) {
		args := parseArgs(req.Params.Arguments)
		query := getStringArg(args, "query")
		if query == "" {
			return toolError("'query' is required"), nil
		}

		filePath := getStringArg(args, "file_path")
		limit := getIntArg(args, "limit", 5)

		searchQuery := query
		if filePath != "" {
			searchQuery = fmt.Sprintf("%s (related to file: %s)", query, filePath)
		}

		result := &relevantContextResult{
			Query:    query,
			FilePath: filePath,
		}

		var wg sync.WaitGroup
		var mu sync.Mutex
		var errs []string

		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := sctx.Client.Search(sctx.ProjectID, api.SearchRequest{
				Query: searchQuery, Mode: "hybrid", EntityTypes: []string{"decision"},
				Limit: limit, IncludeContent: true, UseReranking: true,
			})
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errs = append(errs, fmt.Sprintf("decision search: %v", err))
				return
			}
			result.Decisions = resp
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := sctx.Client.Search(sctx.ProjectID, api.SearchRequest{
				Query: searchQuery, Mode: "hybrid", EntityTypes: []string{"memory"},
				Limit: limit, IncludeContent: true, UseReranking: true,
			})
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errs = append(errs, fmt.Sprintf("memory search: %v", err))
				return
			}
			result.Memories = resp
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			tasks, err := sctx.Client.ListTasks(sctx.ProjectID)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errs = append(errs, fmt.Sprintf("task list: %v", err))
				return
			}
			if len(tasks) > limit {
				tasks = tasks[:limit]
			}
			result.Tasks = tasks
		}()

		wg.Wait()

		decisionCount, memoryCount := 0, 0
		if result.Decisions != nil {
			decisionCount = len(result.Decisions.Results)
		}
		if result.Memories != nil {
			memoryCount = len(result.Memories.Results)
		}

		result.Summary = fmt.Sprintf(
			"Found %d relevant decisions, %d relevant memories, and %d tasks for query: '%s'",
			decisionCount, memoryCount, len(result.Tasks), query,
		)
		if len(errs) > 0 {
			result.Summary += fmt.Sprintf(" (partial results due to errors: %v)", errs)
		}

		data, _ := json.MarshalIndent(result, "", "  ")
		return toolResult(string(data)), nil
	}
}
