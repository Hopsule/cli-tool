package commands

import (
	"fmt"

	"github.com/Cagangedik/cli-tool/internal/config"
	"github.com/spf13/cobra"
)

func NewLogoutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Sign out from Hopsule",
		Long: `Sign out from Hopsule and clear stored credentials.

This will remove your authentication token from local storage.
You will need to run 'hopsule login' again to use authenticated commands.`,
		RunE: runLogout,
	}

	return cmd
}

func runLogout(cmd *cobra.Command, args []string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if logged in
	if !cfg.IsAuthenticated() {
		fmt.Println("You are not currently logged in.")
		return nil
	}

	// Get user info before clearing
	userName := "Unknown"
	userEmail := ""
	if cfg.User != nil {
		userName = cfg.User.Name
		userEmail = cfg.User.Email
	}

	// Clear auth data
	if err := config.ClearAuth(cfg); err != nil {
		return fmt.Errorf("failed to clear credentials: %w", err)
	}

	fmt.Println("┌─────────────────────────────────────────┐")
	fmt.Println("│         ✓ Logout Successful!            │")
	fmt.Println("└─────────────────────────────────────────┘")
	fmt.Println()
	if userEmail != "" {
		fmt.Printf("Signed out from: %s (%s)\n", userName, userEmail)
	} else {
		fmt.Printf("Signed out from: %s\n", userName)
	}
	fmt.Println()
	fmt.Println("Your credentials have been removed from this device.")
	fmt.Println("Run 'hopsule login' to sign in again.")

	return nil
}
