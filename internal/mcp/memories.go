package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Hopsule/cli-tool/internal/api"
	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerMemoryTools(server *gomcp.Server, sctx *ServerContext) {
	server.AddTool(listMemoriesTool(), listMemoriesHandler(sctx))
	server.AddTool(getMemoryTool(), getMemoryHandler(sctx))
	server.AddTool(searchMemoriesTool(), searchMemoriesHandler(sctx))
	server.AddTool(createMemoryTool(), createMemoryHandler(sctx))
}

// ============================================================================
// list-memories
// ============================================================================

func listMemoriesTool() *gomcp.Tool {
	return &gomcp.Tool{
		Name:        "list-memories",
		Description: "List memories in the current project. Memories are persistent, append-only records that preserve reasoning, intent, history, and lessons learned. They explain WHY things are the way they are. Memories are NOT decisions.",
		InputSchema: inputSchema(map[string]interface{}{}, nil),
	}
}

func listMemoriesHandler(sctx *ServerContext) gomcp.ToolHandler {
	return func(ctx context.Context, req *gomcp.CallToolRequest) (*gomcp.CallToolResult, error) {
		memories, err := sctx.Client.ListMemories(sctx.ProjectID)
		if err != nil {
			return toolError(fmt.Sprintf("Failed to list memories: %v", err)), nil
		}

		data, _ := json.MarshalIndent(memories, "", "  ")
		return toolResult(string(data)), nil
	}
}

// ============================================================================
// get-memory
// ============================================================================

func getMemoryTool() *gomcp.Tool {
	return &gomcp.Tool{
		Name:        "get-memory",
		Description: "Get a single memory by ID with full details.",
		InputSchema: inputSchema(map[string]interface{}{
			"id": prop("string", "The memory ID"),
		}, []string{"id"}),
	}
}

func getMemoryHandler(sctx *ServerContext) gomcp.ToolHandler {
	return func(ctx context.Context, req *gomcp.CallToolRequest) (*gomcp.CallToolResult, error) {
		args := parseArgs(req.Params.Arguments)
		id := getStringArg(args, "id")
		if id == "" {
			return toolError("'id' is required"), nil
		}

		memory, err := sctx.Client.GetMemory(sctx.ProjectID, id)
		if err != nil {
			return toolError(fmt.Sprintf("Failed to get memory: %v", err)), nil
		}

		data, _ := json.MarshalIndent(memory, "", "  ")
		return toolResult(string(data)), nil
	}
}

// ============================================================================
// search-memories
// ============================================================================

func searchMemoriesTool() *gomcp.Tool {
	return &gomcp.Tool{
		Name:        "search-memories",
		Description: "Semantic search across memories. Returns memories matching the query using hybrid vector + full-text search with reranking.",
		InputSchema: inputSchema(map[string]interface{}{
			"query": prop("string", "Search query (natural language)"),
			"limit": prop("number", "Maximum results (default: 10, max: 100)"),
			"tags":  arrayProp("Filter results by tags"),
		}, []string{"query"}),
	}
}

func searchMemoriesHandler(sctx *ServerContext) gomcp.ToolHandler {
	return func(ctx context.Context, req *gomcp.CallToolRequest) (*gomcp.CallToolResult, error) {
		args := parseArgs(req.Params.Arguments)
		query := getStringArg(args, "query")
		if query == "" {
			return toolError("'query' is required"), nil
		}

		limit := getIntArg(args, "limit", 10)
		tags := getStringSliceArg(args, "tags")

		searchReq := api.SearchRequest{
			Query:          query,
			Mode:           "hybrid",
			EntityTypes:    []string{"memory"},
			Limit:          limit,
			Tags:           tags,
			IncludeContent: true,
			UseReranking:   true,
		}

		result, err := sctx.Client.Search(sctx.ProjectID, searchReq)
		if err != nil {
			return toolError(fmt.Sprintf("Search failed: %v", err)), nil
		}

		data, _ := json.MarshalIndent(result, "", "  ")
		return toolResult(string(data)), nil
	}
}

// ============================================================================
// create-memory (write delegation)
// ============================================================================

func createMemoryTool() *gomcp.Tool {
	return &gomcp.Tool{
		Name:        "create-memory",
		Description: "[WRITE — delegates to decision-api] Create a new memory. Memories preserve reasoning, intent, history, and lessons learned. They do NOT create authority. The decision-api is the single authority.",
		InputSchema: inputSchema(map[string]interface{}{
			"content": prop("string", "The memory content — reasoning, context, or lesson to preserve"),
			"tags":    arrayProp("Optional tags for categorization"),
		}, []string{"content"}),
	}
}

func createMemoryHandler(sctx *ServerContext) gomcp.ToolHandler {
	return func(ctx context.Context, req *gomcp.CallToolRequest) (*gomcp.CallToolResult, error) {
		args := parseArgs(req.Params.Arguments)
		content := getStringArg(args, "content")
		if content == "" {
			return toolError("'content' is required"), nil
		}

		createReq := api.CreateMemoryRequest{
			Content: content,
			Tags:    getStringSliceArg(args, "tags"),
		}

		memory, err := sctx.Client.CreateMemory(sctx.ProjectID, createReq)
		if err != nil {
			return toolError(fmt.Sprintf("Failed to create memory: %v", err)), nil
		}

		data, _ := json.MarshalIndent(memory, "", "  ")
		return toolResult(fmt.Sprintf("Memory created.\n\n%s", string(data))), nil
	}
}
