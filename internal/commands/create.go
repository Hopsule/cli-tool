package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Cagangedik/cli-tool/internal/api"
	"github.com/Cagangedik/cli-tool/internal/config"
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new decision",
		Long:  "Interactively create a new decision (will be in DRAFT status)",
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

			// Interactive prompts
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Statement: ")
			statement, _ := reader.ReadString('\n')
			statement = strings.TrimSpace(statement)
			if statement == "" {
				return fmt.Errorf("statement is required")
			}

			fmt.Print("Rationale (multi-line, end with empty line):\n")
			var rationaleLines []string
			for {
				line, _ := reader.ReadString('\n')
				line = strings.TrimSpace(line)
				if line == "" && len(rationaleLines) > 0 {
					break
				}
				if line != "" {
					rationaleLines = append(rationaleLines, line)
				}
			}
			rationale := strings.Join(rationaleLines, "\n")

			client := api.NewClient(cfg).
				WithBaseURL(apiURL).
				WithToken(token)

			req := api.CreateDecisionRequest{
				Statement: statement,
				Rationale: rationale,
			}

			decision, err := client.CreateDecision(projectID, req)
			if err != nil {
				return fmt.Errorf("failed to create decision: %w", err)
			}

			fmt.Printf("\nDecision created successfully!\n")
			fmt.Printf("ID: %s\n", decision.ID)
			fmt.Printf("Status: %s\n", decision.Status)

			return nil
		},
	}

	return cmd
}
