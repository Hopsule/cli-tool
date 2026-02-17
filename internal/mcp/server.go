// Package mcp implements the Hopsule MCP (Model Context Protocol) server.
//
// ARCHITECTURAL CONSTRAINTS (NON-NEGOTIABLE):
//   - This MCP server delegates all mutations to decision-api.
//   - decision-api is the SINGLE AUTHORITY for state changes.
//   - AI agents using these tools are ADVISORY — they do not decide.
//   - Write tools explicitly delegate to decision-api; they do not create authority.
package mcp

import (
	"context"
	"fmt"
	"os"

	"github.com/Hopsule/cli-tool/internal/api"
	"github.com/Hopsule/cli-tool/internal/config"
	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// ServerContext holds the shared state for all MCP tool handlers.
type ServerContext struct {
	Client    *api.Client
	ProjectID string
}

// NewMCPServer creates and configures the MCP server with all tools registered.
// It reads the project config (.hopsule) and global auth config to bootstrap.
func NewMCPServer() (*gomcp.Server, *ServerContext, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}
	if !cfg.IsAuthenticated() {
		return nil, nil, fmt.Errorf("not authenticated — run 'hopsule login' first")
	}

	projectCfg, _, err := config.LoadProjectConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("no .hopsule file found — run 'hopsule init' in your project directory first")
	}

	client := api.NewClient(cfg)

	sctx := &ServerContext{
		Client:    client,
		ProjectID: projectCfg.Project.ID,
	}

	server := gomcp.NewServer(
		&gomcp.Implementation{
			Name:    "hopsule",
			Version: "0.9.2",
		},
		nil,
	)

	registerDecisionTools(server, sctx)
	registerMemoryTools(server, sctx)
	registerTaskTools(server, sctx)
	registerContextTools(server, sctx)
	registerStatusTools(server, sctx)

	return server, sctx, nil
}

// Run starts the MCP server with stdio transport.
func Run(ctx context.Context) error {
	server, _, err := NewMCPServer()
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Hopsule MCP server starting (stdio transport)...\n")
	return server.Run(ctx, &gomcp.StdioTransport{})
}
