package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Hopsule/cli-tool/internal/config"
	mcpserver "github.com/Hopsule/cli-tool/internal/mcp"
	"github.com/spf13/cobra"
)

// NewMCPCommand creates the top-level 'mcp' command with subcommands
func NewMCPCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "MCP (Model Context Protocol) server management",
		Long: `Manage the embedded MCP server for AI tool integration.

The MCP server provides decisions, memories, tasks, and context packs
to AI agents (Cursor, Claude Desktop, Claude Code, etc.) via the
Model Context Protocol.

Available subcommands:
  serve    - Start the MCP server (stdio transport)
  install  - Configure AI tools to use this MCP server`,
	}

	cmd.AddCommand(newMCPServeCommand())
	cmd.AddCommand(newMCPInstallCommand())

	return cmd
}

// ============================================================================
// hopsule mcp serve
// ============================================================================

func newMCPServeCommand() *cobra.Command {
	var projectID string
	var projectDir string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the MCP server (stdio transport)",
		Long: `Start the Hopsule MCP server using stdio transport.

This command is typically invoked by AI tools (Cursor, Claude Desktop, etc.)
rather than run manually. It communicates via stdin/stdout using the
Model Context Protocol.

Project identification (in priority order):
  1. --project-id flag (direct project UUID)
  2. HOPSULE_PROJECT_ID environment variable
  3. --project-dir flag (directory containing .hopsule)
  4. .hopsule file in current directory or parents

Prerequisites:
  - Run 'hopsule login' first`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			return mcpserver.Run(ctx, projectID, projectDir)
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&projectID, "project-id", "", "Project ID (UUID) to serve")
	cmd.Flags().StringVar(&projectDir, "project-dir", "", "Project directory containing .hopsule file")

	return cmd
}

// ============================================================================
// hopsule mcp install
// ============================================================================

func newMCPInstallCommand() *cobra.Command {
	var targetIDE string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Configure AI tools to use the Hopsule MCP server",
		Long: `Detect and configure supported AI tools to connect to the Hopsule MCP server.

Supported AI tools:
  - cursor        Cursor IDE (writes to .cursor/mcp.json)
  - claude-desktop  Claude Desktop (writes to claude_desktop_config.json)
  - claude-code   Claude Code CLI (uses 'claude mcp add' command)

Without --ide flag, auto-detects installed tools and prompts for selection.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMCPInstall(targetIDE)
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&targetIDE, "ide", "", "Target IDE: cursor, claude-desktop, claude-code")

	return cmd
}

func runMCPInstall(targetIDE string) error {
	hopsuleBin, err := findHopsuleBinary()
	if err != nil {
		return fmt.Errorf("could not find hopsule binary: %w", err)
	}

	if targetIDE != "" {
		return installForIDE(targetIDE, hopsuleBin)
	}

	// Auto-detect and install for all available IDEs
	installed := detectInstalledIDEs()
	if len(installed) == 0 {
		fmt.Println("No supported AI tools detected.")
		fmt.Println("\nSupported tools: Cursor, Claude Desktop, Claude Code")
		fmt.Println("Use --ide flag to install manually: hopsule mcp install --ide cursor")
		return nil
	}

	for _, ide := range installed {
		fmt.Printf("Configuring %s...\n", ide)
		if err := installForIDE(ide, hopsuleBin); err != nil {
			fmt.Fprintf(os.Stderr, "  Warning: failed to configure %s: %v\n", ide, err)
		} else {
			fmt.Printf("  ✓ %s configured successfully\n", ide)
		}
	}

	return nil
}

func findHopsuleBinary() (string, error) {
	execPath, err := os.Executable()
	if err == nil {
		return execPath, nil
	}

	path, err := exec.LookPath("hopsule")
	if err == nil {
		return path, nil
	}

	return "", fmt.Errorf("hopsule binary not found in PATH")
}

func detectInstalledIDEs() []string {
	var detected []string

	// Cursor: check for .cursor directory in cwd or home
	if _, err := os.Stat(".cursor"); err == nil {
		detected = append(detected, "cursor")
	}

	// Claude Desktop: check for config directory
	configDir := getClaudeDesktopConfigDir()
	if configDir != "" {
		if _, err := os.Stat(configDir); err == nil {
			detected = append(detected, "claude-desktop")
		}
	}

	// Claude Code: check if 'claude' command exists
	if _, err := exec.LookPath("claude"); err == nil {
		detected = append(detected, "claude-code")
	}

	return detected
}

func installForIDE(ide, hopsuleBin string) error {
	switch strings.ToLower(ide) {
	case "cursor":
		return installForCursor(hopsuleBin)
	case "claude-desktop":
		return installForClaudeDesktop(hopsuleBin)
	case "claude-code":
		return installForClaudeCode(hopsuleBin)
	default:
		return fmt.Errorf("unsupported IDE: %s (supported: cursor, claude-desktop, claude-code)", ide)
	}
}

// ============================================================================
// Cursor
// ============================================================================

func getProjectIDForMCP() (string, error) {
	projectCfg, _, err := config.LoadProjectConfig()
	if err != nil {
		return "", fmt.Errorf("no .hopsule file found — run 'hopsule init' first")
	}
	return projectCfg.Project.ID, nil
}

func installForCursor(hopsuleBin string) error {
	projectID, err := getProjectIDForMCP()
	if err != nil {
		return err
	}

	configDir := ".cursor"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create .cursor directory: %w", err)
	}

	configPath := filepath.Join(configDir, "mcp.json")

	var mcpConfig map[string]interface{}
	if data, err := os.ReadFile(configPath); err == nil {
		json.Unmarshal(data, &mcpConfig)
	}
	if mcpConfig == nil {
		mcpConfig = make(map[string]interface{})
	}

	servers, ok := mcpConfig["mcpServers"].(map[string]interface{})
	if !ok {
		servers = make(map[string]interface{})
	}

	servers["hopsule"] = map[string]interface{}{
		"command": hopsuleBin,
		"args":    []string{"mcp", "serve"},
		"env": map[string]string{
			"HOPSULE_PROJECT_ID": projectID,
		},
	}

	mcpConfig["mcpServers"] = servers

	data, err := json.MarshalIndent(mcpConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}

// ============================================================================
// Claude Desktop
// ============================================================================

func getClaudeDesktopConfigDir() string {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Claude")
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "Claude")
	case "linux":
		return filepath.Join(home, ".config", "Claude")
	}
	return ""
}

func installForClaudeDesktop(hopsuleBin string) error {
	projectID, err := getProjectIDForMCP()
	if err != nil {
		return err
	}

	configDir := getClaudeDesktopConfigDir()
	if configDir == "" {
		return fmt.Errorf("unsupported OS for Claude Desktop: %s", runtime.GOOS)
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "claude_desktop_config.json")

	var config map[string]interface{}
	if data, err := os.ReadFile(configPath); err == nil {
		json.Unmarshal(data, &config)
	}
	if config == nil {
		config = make(map[string]interface{})
	}

	servers, ok := config["mcpServers"].(map[string]interface{})
	if !ok {
		servers = make(map[string]interface{})
	}

	servers["hopsule"] = map[string]interface{}{
		"command": hopsuleBin,
		"args":    []string{"mcp", "serve"},
		"env": map[string]string{
			"HOPSULE_PROJECT_ID": projectID,
		},
	}

	config["mcpServers"] = servers

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}

// ============================================================================
// Claude Code
// ============================================================================

func installForClaudeCode(hopsuleBin string) error {
	projectID, err := getProjectIDForMCP()
	if err != nil {
		return err
	}

	claudePath, err := exec.LookPath("claude")
	if err != nil {
		return fmt.Errorf("'claude' command not found — install Claude Code CLI first")
	}

	cmd := exec.Command(claudePath, "mcp", "add", "hopsule",
		"-e", fmt.Sprintf("HOPSULE_PROJECT_ID=%s", projectID),
		"--", hopsuleBin, "mcp", "serve")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
