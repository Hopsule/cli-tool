package commands

import (
	"fmt"

	"github.com/Cagangedik/cli-tool/internal/api"
	"github.com/Cagangedik/cli-tool/internal/config"
	"github.com/spf13/cobra"
)

func NewWhoamiCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "Show current user information",
		Long: `Display information about the currently logged in user.

Shows:
  • User name and email
  • Organizations you belong to
  • Projects you have access to`,
		RunE: runWhoami,
	}

	cmd.Flags().Bool("json", false, "Output as JSON")

	return cmd
}

func runWhoami(cmd *cobra.Command, args []string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if logged in
	if !cfg.IsAuthenticated() {
		fmt.Println("Not logged in.")
		fmt.Println()
		fmt.Println("Run 'hopsule login' to sign in.")
		return nil
	}

	// Show cached info immediately
	fmt.Println("┌─────────────────────────────────────────┐")
	fmt.Println("│           Current User                  │")
	fmt.Println("└─────────────────────────────────────────┘")
	fmt.Println()

	if cfg.User != nil {
		fmt.Printf("  Name:  %s\n", cfg.User.Name)
		fmt.Printf("  Email: %s\n", cfg.User.Email)
		fmt.Printf("  ID:    %s\n", cfg.User.ID)
	}

	// Try to fetch fresh data from API
	client := api.NewClient(cfg)
	meResp, err := client.GetMe()
	if err != nil {
		fmt.Println()
		fmt.Printf("  (Could not fetch latest data: %v)\n", err)
		return nil
	}

	// Update user info from API
	if meResp.User != nil {
		fmt.Println()
		fmt.Println("┌─────────────────────────────────────────┐")
		fmt.Println("│           Organizations                 │")
		fmt.Println("└─────────────────────────────────────────┘")
		fmt.Println()

		if len(meResp.Organizations) == 0 {
			fmt.Println("  No organizations")
		} else {
			for _, org := range meResp.Organizations {
				fmt.Printf("  • %s (@%s)\n", org.Name, org.Slug)
			}
		}

		fmt.Println()
		fmt.Println("┌─────────────────────────────────────────┐")
		fmt.Println("│           Projects                      │")
		fmt.Println("└─────────────────────────────────────────┘")
		fmt.Println()

		if len(meResp.Projects) == 0 {
			fmt.Println("  No projects")
		} else {
			for _, proj := range meResp.Projects {
				desc := proj.Description
				if len(desc) > 40 {
					desc = desc[:37] + "..."
				}
				if desc != "" {
					fmt.Printf("  • %s - %s\n", proj.Name, desc)
				} else {
					fmt.Printf("  • %s\n", proj.Name)
				}
			}
		}
	}

	fmt.Println()
	fmt.Println("Config file: ~/.decision-cli/config.yaml")

	return nil
}
