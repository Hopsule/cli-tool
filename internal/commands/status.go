package commands

import (
	"encoding/json"
	"fmt"

	"github.com/Cagangedik/cli-tool/internal/api"
	"github.com/Cagangedik/cli-tool/internal/config"
	"github.com/spf13/cobra"
)

func NewStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show current project status",
		Long:  "Display statistics about decisions in the current project",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			projectID, _ := cmd.Flags().GetString("project")
			if projectID == "" {
				projectID = cfg.Project
			}
			if projectID == "" {
				return fmt.Errorf("project ID is required (use --project or set in config)")
			}

			apiURL, _ := cmd.Flags().GetString("api-url")
			if apiURL == "" {
				apiURL = cfg.APIURL
			}
			if apiURL == "" {
				return fmt.Errorf("API URL is required (use --api-url or set in config)")
			}

			token, _ := cmd.Flags().GetString("token")
			if token == "" {
				token = cfg.Token
			}

			client := api.NewClient(cfg).
				WithBaseURL(apiURL).
				WithToken(token)

			status, err := client.GetProjectStatus(projectID)
			if err != nil {
				return fmt.Errorf("failed to get status: %w", err)
			}

			output, _ := cmd.Flags().GetString("output")
			if output == "json" {
				jsonData, _ := json.MarshalIndent(status, "", "  ")
				fmt.Println(string(jsonData))
			} else {
				fmt.Printf("Project: %s\n\n", status.ProjectID)
				fmt.Printf("Total Decisions: %d\n", status.TotalDecisions)
				fmt.Printf("  Accepted:   %d\n", status.Accepted)
				fmt.Printf("  Pending:   %d\n", status.Pending)
				fmt.Printf("  Draft:     %d\n", status.Draft)
				fmt.Printf("  Deprecated: %d\n", status.Deprecated)
			}

			return nil
		},
	}

	cmd.Flags().StringP("output", "o", "text", "Output format (text, json)")

	return cmd
}
