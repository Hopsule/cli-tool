package commands

import (
	"fmt"

	"github.com/Cagangedik/cli-tool/internal/api"
	"github.com/Cagangedik/cli-tool/internal/config"
	"github.com/spf13/cobra"
)

func NewDeprecateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deprecate <decision-id>",
		Short: "Deprecate a decision",
		Long:  "Deprecate a decision, moving it to DEPRECATED status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			decisionID := args[0]

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

			decision, err := client.DeprecateDecision(projectID, decisionID)
			if err != nil {
				return fmt.Errorf("failed to deprecate decision: %w", err)
			}

			fmt.Printf("Decision deprecated successfully!\n")
			fmt.Printf("ID: %s\n", decision.ID)
			fmt.Printf("Status: %s\n", decision.Status)

			return nil
		},
	}

	return cmd
}
