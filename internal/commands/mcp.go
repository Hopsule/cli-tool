package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Hopsule/cli-tool/internal/config"
	"github.com/spf13/cobra"
)

const hostedMCPURL = "https://mcp.hopsule.com"

// NewMCPCommand creates the top-level 'mcp' command with subcommands
func NewMCPCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "MCP (Model Context Protocol) server management",
		Long: `Manage MCP connections for AI tool integration.

The Hopsule MCP server provides decisions, memories, tasks, and context
to AI agents (Cursor, Claude Desktop, Claude Code, etc.) via the
Model Context Protocol.

Available subcommands:
  install  - Configure AI tools to use the Hosted MCP server
  token    - Manage MCP tokens`,
	}

	cmd.AddCommand(newMCPInstallCommand())
	cmd.AddCommand(newMCPTokenCommand())

	return cmd
}

// ============================================================================
// hopsule mcp install
// ============================================================================

func newMCPInstallCommand() *cobra.Command {
	var targetIDE string
	var token string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Configure AI tools to use the Hosted Hopsule MCP server",
		Long: `Configure supported AI tools to connect to the Hosted Hopsule MCP server.

This uses the hosted MCP endpoint at ` + hostedMCPURL + `
You need an MCP token — generate one at app.hopsule.com → Settings → MCP,
or use 'hopsule mcp token create'.

Supported AI tools:
  - cursor          Cursor IDE (writes to .cursor/mcp.json)
  - claude-desktop  Claude Desktop (writes to claude_desktop_config.json)
  - claude-code     Claude Code CLI (uses 'claude mcp add' command)

Without --ide flag, auto-detects installed tools.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMCPInstall(targetIDE, token)
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&targetIDE, "ide", "", "Target IDE: cursor, claude-desktop, claude-code")
	cmd.Flags().StringVar(&token, "token", "", "MCP token (or set HOPSULE_MCP_TOKEN env var)")

	return cmd
}

func newMCPTokenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Manage MCP tokens",
		Long: `Create and manage MCP tokens for the Hosted MCP server.

MCP tokens authenticate your IDE with the Hopsule MCP endpoint.
You can also manage tokens from the web dashboard at
app.hopsule.com → Settings → MCP.`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "create [name]",
		Short: "Create a new MCP token",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.GetConfig()
			if err != nil || !cfg.IsAuthenticated() {
				return fmt.Errorf("not authenticated — run 'hopsule login' first")
			}

			name := "CLI MCP Token"
			if len(args) > 0 {
				name = args[0]
			}

			fmt.Printf("Creating MCP token '%s'...\n", name)
			fmt.Println("\nTo create tokens, visit: https://app.hopsule.com → Settings → MCP")
			fmt.Println("Or use the API: POST /auth/mcp/tokens")
			return nil
		},
	})

	return cmd
}

func runMCPInstall(targetIDE, token string) error {
	if token == "" {
		token = os.Getenv("HOPSULE_MCP_TOKEN")
	}
	if token == "" {
		return fmt.Errorf("MCP token required. Set --token flag or HOPSULE_MCP_TOKEN env var.\nGenerate a token at: https://app.hopsule.com → Settings → MCP")
	}

	if targetIDE != "" {
		return installHostedForIDE(targetIDE, token)
	}

	installed := detectInstalledIDEs()
	if len(installed) == 0 {
		fmt.Println("No supported AI tools detected.")
		fmt.Println("\nSupported tools: Cursor, Claude Desktop, Claude Code")
		fmt.Println("Use --ide flag to install manually: hopsule mcp install --ide cursor --token <token>")
		return nil
	}

	for _, ide := range installed {
		fmt.Printf("Configuring %s with Hosted MCP...\n", ide)
		if err := installHostedForIDE(ide, token); err != nil {
			fmt.Fprintf(os.Stderr, "  Warning: failed to configure %s: %v\n", ide, err)
		} else {
			fmt.Printf("  ✓ %s configured with hosted MCP\n", ide)
		}
	}

	return nil
}

func installHostedForIDE(ide, token string) error {
	switch strings.ToLower(ide) {
	case "cursor":
		return installHostedForCursor(token)
	case "claude-desktop":
		return installHostedForClaudeDesktop(token)
	case "claude-code":
		return installHostedForClaudeCode(token)
	default:
		return fmt.Errorf("unsupported IDE: %s (supported: cursor, claude-desktop, claude-code)", ide)
	}
}


// ============================================================================
// Hosted MCP Install Functions
// ============================================================================

func hostedMCPServerConfig(token string) map[string]interface{} {
	return map[string]interface{}{
		"url": hostedMCPURL + "/mcp",
		"headers": map[string]string{
			"Authorization": "Bearer " + token,
		},
	}
}

func installHostedForCursor(token string) error {
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

	servers["hopsule"] = hostedMCPServerConfig(token)
	mcpConfig["mcpServers"] = servers

	data, err := json.MarshalIndent(mcpConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}

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

func installHostedForClaudeDesktop(token string) error {
	configDir := getClaudeDesktopConfigDir()
	if configDir == "" {
		return fmt.Errorf("unsupported OS for Claude Desktop: %s", runtime.GOOS)
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "claude_desktop_config.json")

	var cfgMap map[string]interface{}
	if data, err := os.ReadFile(configPath); err == nil {
		json.Unmarshal(data, &cfgMap)
	}
	if cfgMap == nil {
		cfgMap = make(map[string]interface{})
	}

	servers, ok := cfgMap["mcpServers"].(map[string]interface{})
	if !ok {
		servers = make(map[string]interface{})
	}

	servers["hopsule"] = hostedMCPServerConfig(token)
	cfgMap["mcpServers"] = servers

	data, err := json.MarshalIndent(cfgMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}

func installHostedForClaudeCode(token string) error {
	claudePath, err := exec.LookPath("claude")
	if err != nil {
		return fmt.Errorf("'claude' command not found — install Claude Code CLI first")
	}

	cmd := exec.Command(claudePath, "mcp", "add", "hopsule",
		"--transport", "streamable-http",
		hostedMCPURL+"/mcp",
		"--header", "Authorization: Bearer "+token)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func detectInstalledIDEs() []string {
	var detected []string

	if _, err := os.Stat(".cursor"); err == nil {
		detected = append(detected, "cursor")
	}

	configDir := getClaudeDesktopConfigDir()
	if configDir != "" {
		if _, err := os.Stat(configDir); err == nil {
			detected = append(detected, "claude-desktop")
		}
	}

	if _, err := exec.LookPath("claude"); err == nil {
		detected = append(detected, "claude-code")
	}

	return detected
}
