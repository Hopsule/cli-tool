package main

import (
	"fmt"
	"os"

	"github.com/Cagangedik/cli-tool/internal/commands"
	"github.com/Cagangedik/cli-tool/internal/config"
	"github.com/Cagangedik/cli-tool/internal/ui"
	"github.com/spf13/cobra"
)

// Version information (set by goreleaser)
var (
	version = "0.7.1"
	commit  = "none"
	date    = "unknown"
)

func runInteractiveTUI() {
	for {
		action, err := ui.RunInteractive()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if action == "" {
			// User quit
			return
		}

		cfg, _ := config.GetConfig()
		if cfg == nil {
			cfg = &config.Config{}
		}

		// Execute the selected action
		switch {
		case action == "login":
			if err := ui.ExecuteLogin(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "\n  Login failed: %v\n\n", err)
				fmt.Println("  Press Enter to continue...")
				fmt.Scanln()
			}
			// Restart TUI after login
			continue

		case action == "logout":
			if cfg.IsAuthenticated() {
				userName := ""
				if cfg.User != nil {
					userName = cfg.User.Name
				}
				if err := config.ClearAuth(cfg); err != nil {
					fmt.Fprintf(os.Stderr, "\n  Logout failed: %v\n\n", err)
				} else {
					fmt.Printf("\n  ✓ Logged out from %s\n\n", userName)
					fmt.Println("  Press Enter to continue...")
					fmt.Scanln()
				}
			}
			// Restart TUI after logout
			continue

		case len(action) > 8 && action[:8] == "project:":
			// Project selected
			projectID := action[8:]
			fmt.Printf("\n  Selected project: %s\n", projectID)
			fmt.Println("  Use 'hopsule init' to connect this directory to the project.")
			fmt.Println()
			return

		default:
			return
		}
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "hopsule",
		Short: "Hopsule CLI - Decision-first project memory",
		Long: `Hopsule CLI - Decision-first, context-aware, portable memory system.

Hopsule helps you track architectural decisions, project context, 
and team knowledge in a portable, AI-friendly format.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
		Run: func(cmd *cobra.Command, args []string) {
			// Run interactive TUI when no subcommand is provided
			runInteractiveTUI()
		},
	}

	// Persistent flags (available to all commands)
	rootCmd.PersistentFlags().String("api-url", "", "Override API URL")
	rootCmd.PersistentFlags().String("token", "", "Override authentication token")
	rootCmd.PersistentFlags().String("project", "", "Override project ID")

	// ========================================================================
	// AUTH COMMANDS
	// ========================================================================
	rootCmd.AddCommand(commands.NewLoginCommand())
	rootCmd.AddCommand(commands.NewLogoutCommand())
	rootCmd.AddCommand(commands.NewWhoamiCommand())

	// ========================================================================
	// ORGANIZATION & PROJECT COMMANDS
	// ========================================================================
	rootCmd.AddCommand(commands.NewOrgsCommand())
	rootCmd.AddCommand(commands.NewProjectsCommand())
	rootCmd.AddCommand(commands.NewInitCommand())

	// ========================================================================
	// DECISION COMMANDS
	// ========================================================================
	rootCmd.AddCommand(commands.NewListCommand())
	rootCmd.AddCommand(commands.NewGetCommand())
	rootCmd.AddCommand(commands.NewCreateCommand())
	rootCmd.AddCommand(commands.NewAcceptCommand())
	rootCmd.AddCommand(commands.NewDeprecateCommand())

	// ========================================================================
	// UTILITY COMMANDS
	// ========================================================================
	rootCmd.AddCommand(commands.NewConfigCommand())
	rootCmd.AddCommand(commands.NewStatusCommand())
	rootCmd.AddCommand(commands.NewSyncCommand())

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
