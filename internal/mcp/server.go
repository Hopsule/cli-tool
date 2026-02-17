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
// If projectDir is non-empty, it looks for .hopsule in that directory instead of cwd.
func NewMCPServer(projectDir string) (*gomcp.Server, *ServerContext, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}
	if !cfg.IsAuthenticated() {
		return nil, nil, fmt.Errorf("not authenticated — run 'hopsule login' first")
	}

	var projectCfg *config.ProjectConfig
	if projectDir != "" {
		projectCfg, _, err = config.LoadProjectConfigFrom(projectDir)
	} else {
		projectCfg, _, err = config.LoadProjectConfig()
	}
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
			Version: "0.9.7",
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
// projectDir overrides the working directory for .hopsule lookup.
func Run(ctx context.Context, projectDir string) error {
	server, _, err := NewMCPServer(projectDir)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Hopsule MCP server starting (stdio transport)...\n")
	return server.Run(ctx, &gomcp.StdioTransport{})
}
