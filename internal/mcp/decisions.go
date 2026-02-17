package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Hopsule/cli-tool/internal/api"
	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerDecisionTools(server *gomcp.Server, sctx *ServerContext) {
	server.AddTool(listDecisionsTool(), listDecisionsHandler(sctx))
	server.AddTool(getDecisionTool(), getDecisionHandler(sctx))
	server.AddTool(searchDecisionsTool(), searchDecisionsHandler(sctx))
	server.AddTool(createDecisionTool(), createDecisionHandler(sctx))
	server.AddTool(acceptDecisionTool(), acceptDecisionHandler(sctx))
	server.AddTool(deprecateDecisionTool(), deprecateDecisionHandler(sctx))
}

// ============================================================================
// list-decisions
// ============================================================================

func listDecisionsTool() *gomcp.Tool {
	return &gomcp.Tool{
		Name:        "list-decisions",
		Description: "List decisions in the current project. Decisions are explicit, reviewed, and versioned commitments that define 'what MUST be followed'. States: DRAFT → PENDING → ACCEPTED → DEPRECATED.",
		InputSchema: inputSchema(map[string]interface{}{
			"status": prop("string", "Filter by status. Use 'accepted' for accepted-only. Leave empty for all."),
		}, nil),
	}
}

func listDecisionsHandler(sctx *ServerContext) gomcp.ToolHandler {
	return func(ctx context.Context, req *gomcp.CallToolRequest) (*gomcp.CallToolResult, error) {
		args := parseArgs(req.Params.Arguments)
		status := getStringArg(args, "status")

		var decisions []api.Decision
		var err error

		switch status {
		case "accepted":
			decisions, err = sctx.Client.ListAcceptedDecisions(sctx.ProjectID)
		default:
			decisions, err = sctx.Client.ListDecisions(sctx.ProjectID)
		}

		if err != nil {
			return toolError(fmt.Sprintf("Failed to list decisions: %v", err)), nil
		}

		data, _ := json.MarshalIndent(decisions, "", "  ")
		return toolResult(string(data)), nil
	}
}

// ============================================================================
// get-decision
// ============================================================================

func getDecisionTool() *gomcp.Tool {
	return &gomcp.Tool{
		Name:        "get-decision",
		Description: "Get a single decision by ID with full details including statement, rationale, status, and health.",
		InputSchema: inputSchema(map[string]interface{}{
			"id": prop("string", "The decision ID"),
		}, []string{"id"}),
	}
}

func getDecisionHandler(sctx *ServerContext) gomcp.ToolHandler {
	return func(ctx context.Context, req *gomcp.CallToolRequest) (*gomcp.CallToolResult, error) {
		args := parseArgs(req.Params.Arguments)
		id := getStringArg(args, "id")
		if id == "" {
			return toolError("'id' is required"), nil
		}

		decision, err := sctx.Client.GetDecision(sctx.ProjectID, id)
		if err != nil {
			return toolError(fmt.Sprintf("Failed to get decision: %v", err)), nil
		}

		data, _ := json.MarshalIndent(decision, "", "  ")
		return toolResult(string(data)), nil
	}
}

// ============================================================================
// search-decisions
// ============================================================================

func searchDecisionsTool() *gomcp.Tool {
	return &gomcp.Tool{
		Name:        "search-decisions",
		Description: "Semantic search across decisions. Uses hybrid vector + full-text search with reranking.",
		InputSchema: inputSchema(map[string]interface{}{
			"query": prop("string", "Search query (natural language)"),
			"limit": prop("number", "Maximum results (default: 10, max: 100)"),
			"tags":  arrayProp("Filter results by tags"),
		}, []string{"query"}),
	}
}

func searchDecisionsHandler(sctx *ServerContext) gomcp.ToolHandler {
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
			EntityTypes:    []string{"decision"},
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
// create-decision (write delegation)
// ============================================================================

func createDecisionTool() *gomcp.Tool {
	return &gomcp.Tool{
		Name:        "create-decision",
		Description: "[WRITE — delegates to decision-api] Create a new decision draft. This does NOT accept the decision — it creates a DRAFT that requires human review. The decision-api is the single authority.",
		InputSchema: inputSchema(map[string]interface{}{
			"statement": prop("string", "The decision statement — what MUST be followed"),
			"rationale": prop("string", "Why this decision is being made"),
			"tags":      arrayProp("Optional tags for categorization"),
		}, []string{"statement", "rationale"}),
	}
}

func createDecisionHandler(sctx *ServerContext) gomcp.ToolHandler {
	return func(ctx context.Context, req *gomcp.CallToolRequest) (*gomcp.CallToolResult, error) {
		args := parseArgs(req.Params.Arguments)
		statement := getStringArg(args, "statement")
		rationale := getStringArg(args, "rationale")
		if statement == "" || rationale == "" {
			return toolError("'statement' and 'rationale' are required"), nil
		}

		createReq := api.CreateDecisionRequest{
			Statement: statement,
			Rationale: rationale,
			Tags:      getStringSliceArg(args, "tags"),
		}

		decision, err := sctx.Client.CreateDecision(sctx.ProjectID, createReq)
		if err != nil {
			return toolError(fmt.Sprintf("Failed to create decision: %v", err)), nil
		}

		data, _ := json.MarshalIndent(decision, "", "  ")
		return toolResult(fmt.Sprintf("Decision draft created (state: DRAFT). It requires human review to be accepted.\n\n%s", string(data))), nil
	}
}

// ============================================================================
// accept-decision (write delegation)
// ============================================================================

func acceptDecisionTool() *gomcp.Tool {
	return &gomcp.Tool{
		Name:        "accept-decision",
		Description: "[WRITE — delegates to decision-api] Accept a decision, making it an authoritative constraint. Requires DECIDER authority role.",
		InputSchema: inputSchema(map[string]interface{}{
			"id":          prop("string", "The decision ID to accept"),
			"accepted_by": prop("string", "Who is accepting (defaults to 'MCP Agent')"),
		}, []string{"id"}),
	}
}

func acceptDecisionHandler(sctx *ServerContext) gomcp.ToolHandler {
	return func(ctx context.Context, req *gomcp.CallToolRequest) (*gomcp.CallToolResult, error) {
		args := parseArgs(req.Params.Arguments)
		id := getStringArg(args, "id")
		if id == "" {
			return toolError("'id' is required"), nil
		}
		acceptedBy := getStringArg(args, "accepted_by")
		if acceptedBy == "" {
			acceptedBy = "MCP Agent"
		}

		decision, err := sctx.Client.AcceptDecision(sctx.ProjectID, id, acceptedBy)
		if err != nil {
			return toolError(fmt.Sprintf("Failed to accept decision: %v", err)), nil
		}

		data, _ := json.MarshalIndent(decision, "", "  ")
		return toolResult(fmt.Sprintf("Decision accepted. It is now an authoritative constraint.\n\n%s", string(data))), nil
	}
}

// ============================================================================
// deprecate-decision (write delegation)
// ============================================================================

func deprecateDecisionTool() *gomcp.Tool {
	return &gomcp.Tool{
		Name:        "deprecate-decision",
		Description: "[WRITE — delegates to decision-api] Deprecate a decision, removing it from active enforcement. Requires DECIDER authority role.",
		InputSchema: inputSchema(map[string]interface{}{
			"id": prop("string", "The decision ID to deprecate"),
		}, []string{"id"}),
	}
}

func deprecateDecisionHandler(sctx *ServerContext) gomcp.ToolHandler {
	return func(ctx context.Context, req *gomcp.CallToolRequest) (*gomcp.CallToolResult, error) {
		args := parseArgs(req.Params.Arguments)
		id := getStringArg(args, "id")
		if id == "" {
			return toolError("'id' is required"), nil
		}

		decision, err := sctx.Client.DeprecateDecision(sctx.ProjectID, id)
		if err != nil {
			return toolError(fmt.Sprintf("Failed to deprecate decision: %v", err)), nil
		}

		data, _ := json.MarshalIndent(decision, "", "  ")
		return toolResult(fmt.Sprintf("Decision deprecated.\n\n%s", string(data))), nil
	}
}
