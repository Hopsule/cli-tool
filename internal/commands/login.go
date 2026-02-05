package commands

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/Cagangedik/cli-tool/internal/api"
	"github.com/Cagangedik/cli-tool/internal/config"
	"github.com/spf13/cobra"
)

const (
	pollInterval    = 2 * time.Second
	maxPollAttempts = 300 // 10 minutes / 2 seconds
)

func NewLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Sign in to Hopsule",
		Long: `Sign in to Hopsule using your browser.

This command will:
1. Open your browser to authenticate
2. Wait for you to complete sign-in
3. Save your credentials locally

After signing in, you can use all authenticated CLI commands.`,
		RunE: runLogin,
	}

	cmd.Flags().String("api-url", "", "Override API URL")
	cmd.Flags().String("web-url", "", "Override Web URL")
	cmd.Flags().Bool("no-browser", false, "Don't open browser automatically")

	return cmd
}

func runLogin(cmd *cobra.Command, args []string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		cfg = &config.Config{}
	}

	// Check if already logged in
	if cfg.IsAuthenticated() {
		fmt.Printf("Already logged in as %s (%s)\n", cfg.User.Name, cfg.User.Email)
		fmt.Println("Use 'hopsule logout' to sign out first.")
		return nil
	}

	// Override URLs if provided
	apiURL, _ := cmd.Flags().GetString("api-url")
	if apiURL != "" {
		cfg.APIURL = apiURL
	}
	if cfg.APIURL == "" {
		cfg.APIURL = "http://localhost:8080"
	}

	webURL, _ := cmd.Flags().GetString("web-url")
	if webURL != "" {
		cfg.WebURL = webURL
	}
	if cfg.WebURL == "" {
		cfg.WebURL = "http://localhost:3000"
	}

	noBrowser, _ := cmd.Flags().GetBool("no-browser")

	// Create API client
	client := api.NewClient(cfg)

	// Get device name
	deviceName := getDeviceName()

	fmt.Println("┌─────────────────────────────────────────┐")
	fmt.Println("│           Hopsule CLI Login             │")
	fmt.Println("└─────────────────────────────────────────┘")
	fmt.Println()

	// Step 1: Initialize device auth
	fmt.Print("Initializing login... ")
	initResp, err := client.DeviceAuthInit(deviceName)
	if err != nil {
		fmt.Println("✗")
		return fmt.Errorf("failed to initialize login: %w", err)
	}
	fmt.Println("✓")

	// Step 2: Open browser
	authURL := fmt.Sprintf("%s/auth/device?code=%s", cfg.WebURL, initResp.Code)

	fmt.Println()
	fmt.Printf("Device Code: %s\n", initResp.Code)
	fmt.Println()

	if !noBrowser {
		fmt.Println("Opening browser to complete sign-in...")
		if err := openBrowser(authURL); err != nil {
			fmt.Printf("Could not open browser automatically.\n")
		}
	}

	fmt.Println()
	fmt.Println("If the browser doesn't open, visit this URL:")
	fmt.Printf("  %s\n", authURL)
	fmt.Println()

	// Step 3: Poll for completion
	fmt.Println("Waiting for authentication...")
	fmt.Println("(Press Ctrl+C to cancel)")
	fmt.Println()

	token, userInfo, err := pollForCompletion(client, initResp.Code)
	if err != nil {
		return err
	}

	// Step 4: Save config
	cfg.Token = token
	cfg.User = &config.User{
		ID:        userInfo.UserID,
		Email:     userInfo.Email,
		Name:      userInfo.Name,
		AvatarURL: userInfo.AvatarURL,
	}

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println()
	fmt.Println("┌─────────────────────────────────────────┐")
	fmt.Println("│         ✓ Login Successful!             │")
	fmt.Println("└─────────────────────────────────────────┘")
	fmt.Println()
	fmt.Printf("Signed in as: %s (%s)\n", userInfo.Name, userInfo.Email)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  • Run 'hopsule whoami' to see your account info")
	fmt.Println("  • Run 'hopsule orgs' to list your organizations")
	fmt.Println("  • Run 'hopsule projects' to list your projects")
	fmt.Println("  • Run 'hopsule init' to connect this directory to a project")

	return nil
}

func pollForCompletion(client *api.Client, code string) (string, *api.DeviceAuthPollResponse, error) {
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinnerIdx := 0

	for attempt := 0; attempt < maxPollAttempts; attempt++ {
		// Show spinner
		fmt.Printf("\r%s Waiting for browser authentication... (%ds)", spinner[spinnerIdx], attempt*2)
		spinnerIdx = (spinnerIdx + 1) % len(spinner)

		resp, err := client.DeviceAuthPoll(code)
		if err != nil {
			fmt.Println()
			return "", nil, fmt.Errorf("failed to check login status: %w", err)
		}

		switch resp.Status {
		case "complete":
			fmt.Printf("\r✓ Authentication complete!                    \n")
			return resp.Token, resp, nil
		case "expired":
			fmt.Println()
			return "", nil, fmt.Errorf("login session expired - please try again")
		case "pending":
			// Continue polling
		default:
			fmt.Println()
			return "", nil, fmt.Errorf("unexpected status: %s", resp.Status)
		}

		time.Sleep(pollInterval)
	}

	fmt.Println()
	return "", nil, fmt.Errorf("login timed out after %d minutes", maxPollAttempts*int(pollInterval.Seconds())/60)
}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}

func getDeviceName() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	osName := runtime.GOOS
	switch osName {
	case "darwin":
		osName = "macOS"
	case "linux":
		osName = "Linux"
	case "windows":
		osName = "Windows"
	}

	return fmt.Sprintf("CLI on %s (%s)", osName, hostname)
}
