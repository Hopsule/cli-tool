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

// resolveProjectID determines the project ID using this priority:
//  1. Explicit projectID argument (from --project-id flag)
//  2. HOPSULE_PROJECT_ID environment variable
//  3. .hopsule file in projectDir (from --project-dir flag)
//  4. .hopsule file in cwd or parent directories
func resolveProjectID(projectID, projectDir string) (string, error) {
	if projectID != "" {
		return projectID, nil
	}

	if envID := os.Getenv("HOPSULE_PROJECT_ID"); envID != "" {
		return envID, nil
	}

	var projectCfg *config.ProjectConfig
	var err error
	if projectDir != "" {
		projectCfg, _, err = config.LoadProjectConfigFrom(projectDir)
	} else {
		projectCfg, _, err = config.LoadProjectConfig()
	}
	if err == nil && projectCfg != nil {
		return projectCfg.Project.ID, nil
	}

	return "", fmt.Errorf(
		"project not identified — set HOPSULE_PROJECT_ID env var, " +
			"use --project-id flag, or run 'hopsule init' in your project directory",
	)
}

// NewMCPServer creates and configures the MCP server with all tools registered.
// It resolves the project via env var, flag, or .hopsule file.
func NewMCPServer(projectID, projectDir string) (*gomcp.Server, *ServerContext, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}
	if !cfg.IsAuthenticated() {
		return nil, nil, fmt.Errorf("not authenticated — run 'hopsule login' first")
	}

	pid, err := resolveProjectID(projectID, projectDir)
	if err != nil {
		return nil, nil, err
	}

	client := api.NewClient(cfg)

	sctx := &ServerContext{
		Client:    client,
		ProjectID: pid,
	}

	server := gomcp.NewServer(
		&gomcp.Implementation{
			Name:    "hopsule",
			Version: "0.9.8",
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
func Run(ctx context.Context, projectID, projectDir string) error {
	server, _, err := NewMCPServer(projectID, projectDir)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Hopsule MCP server starting (stdio transport)...\n")
	return server.Run(ctx, &gomcp.StdioTransport{})
}
