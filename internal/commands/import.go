package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Cagangedik/cli-tool/internal/api"
	"github.com/Cagangedik/cli-tool/internal/config"
	"github.com/spf13/cobra"
)

func NewImportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import [file]",
		Short: "Import decisions from a file",
		Long: `Import decisions from a JSON file into decision-api.

This command reads decisions from a file and submits them to decision-api
for processing. The file should contain decisions in JSON format.

Note: This command does NOT bypass decision-api authority. All imported
decisions are submitted through the proper API channels.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]

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

			// Read the import file
			data, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", filePath, err)
			}

			// Parse the import data
			var importData ImportData
			if err := json.Unmarshal(data, &importData); err != nil {
				return fmt.Errorf("failed to parse import file: %w", err)
			}

			client := api.NewClient(cfg).
				WithBaseURL(apiURL).
				WithToken(token)

			// Import each decision
			imported := 0
			for _, decision := range importData.Decisions {
				_, err := client.CreateDecision(projectID, api.CreateDecisionRequest{
					Statement: decision.Statement,
					Rationale: decision.Rationale,
					Tags:      decision.Tags,
				})
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to import decision '%s': %v\n", decision.Statement, err)
					continue
				}
				imported++
			}

			fmt.Printf("Successfully imported %d of %d decisions.\n", imported, len(importData.Decisions))
			return nil
		},
	}

	return cmd
}

// ImportData represents the structure of an import file
type ImportData struct {
	Decisions []ImportDecision `json:"decisions"`
}

// ImportDecision represents a single decision to import
type ImportDecision struct {
	Statement string   `json:"statement"`
	Rationale string   `json:"rationale"`
	Tags      []string `json:"tags,omitempty"`
}
