package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Cagangedik/cli-tool/internal/api"
	"github.com/Cagangedik/cli-tool/internal/config"
	"github.com/spf13/cobra"
)

func NewGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <decision-id>",
		Short: "Get decision details",
		Long:  "Retrieve detailed information about a specific decision",
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

			decision, err := client.GetDecision(projectID, decisionID)
			if err != nil {
				return fmt.Errorf("failed to get decision: %w", err)
			}

			output, _ := cmd.Flags().GetString("output")
			if output == "json" {
				jsonData, _ := json.MarshalIndent(decision, "", "  ")
				fmt.Println(string(jsonData))
			} else {
				fmt.Printf("ID: %s\n", decision.ID)
				fmt.Printf("Statement: %s\n", decision.Statement)
				fmt.Printf("Status: %s\n", decision.Status)
				fmt.Printf("Created: %s\n", decision.CreatedAt)
				fmt.Printf("Updated: %s\n", decision.UpdatedAt)
				if decision.AcceptedAt != nil {
					fmt.Printf("Accepted: %s", *decision.AcceptedAt)
					if decision.AcceptedBy != nil {
						fmt.Printf(" by %s", *decision.AcceptedBy)
					}
					fmt.Println()
				}
				if len(decision.Tags) > 0 {
					fmt.Printf("Tags: %s\n", strings.Join(decision.Tags, ", "))
				}
				fmt.Printf("\nRationale:\n%s\n", decision.Rationale)
			}

			return nil
		},
	}

	cmd.Flags().StringP("output", "o", "text", "Output format (text, json)")

	return cmd
}
