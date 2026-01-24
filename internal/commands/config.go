package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Cagangedik/cli-tool/internal/config"
	"github.com/spf13/cobra"
)

func NewConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure CLI settings",
		Long:  "Interactively configure API URL, token, and default project",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.GetConfig()
			if err != nil {
				cfg = &config.Config{}
			}

			reader := bufio.NewReader(os.Stdin)

			fmt.Printf("API URL [%s]: ", cfg.APIURL)
			apiURL, _ := reader.ReadString('\n')
			apiURL = strings.TrimSpace(apiURL)
			if apiURL != "" {
				cfg.APIURL = apiURL
			}
			if cfg.APIURL == "" {
				cfg.APIURL = "http://localhost:8080"
			}

			fmt.Printf("Token [%s]: ", maskToken(cfg.Token))
			token, _ := reader.ReadString('\n')
			token = strings.TrimSpace(token)
			if token != "" {
				cfg.Token = token
			}

			fmt.Printf("Default Project ID [%s]: ", cfg.Project)
			project, _ := reader.ReadString('\n')
			project = strings.TrimSpace(project)
			if project != "" {
				cfg.Project = project
			}

			if err := config.SaveConfig(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Println("\nConfiguration saved successfully!")
			return nil
		},
	}

	return cmd
}

func maskToken(token string) string {
	if token == "" {
		return ""
	}
	if len(token) <= 8 {
		return "***"
	}
	return token[:4] + "..." + token[len(token)-4:]
}
