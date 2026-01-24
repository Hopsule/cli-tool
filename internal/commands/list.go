package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Cagangedik/cli-tool/internal/api"
	"github.com/Cagangedik/cli-tool/internal/config"
	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all decisions",
		Long:  "List all decisions for the current project",
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

			decisions, err := client.ListDecisions(projectID)
			if err != nil {
				return fmt.Errorf("failed to list decisions: %w", err)
			}

			if len(decisions) == 0 {
				fmt.Println("No decisions found.")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "ID\tTITLE\tSTATUS\tCREATED")
			fmt.Fprintln(w, "---\t-----\t------\t-------")

			for _, d := range decisions {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					truncate(d.ID, 12),
					truncate(d.Statement, 40),
					d.Status,
					truncate(d.CreatedAt, 20),
				)
			}

			return w.Flush()
		},
	}

	return cmd
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
