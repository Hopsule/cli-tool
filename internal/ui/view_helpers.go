package ui

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/Hopsule/cli-tool/internal/api"
	"github.com/Hopsule/cli-tool/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

// ============================================================================
// STRING HELPERS
// ============================================================================

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func wordWrap(s string, width int) string {
	if len(s) <= width {
		return s
	}
	var result strings.Builder
	words := strings.Fields(s)
	lineLen := 0
	for i, word := range words {
		if i > 0 && lineLen+len(word)+1 > width {
			result.WriteString("\n  ")
			lineLen = 0
		} else if i > 0 {
			result.WriteString(" ")
			lineLen++
		}
		result.WriteString(word)
		lineLen += len(word)
	}
	return result.String()
}

// ============================================================================
// MARKDOWN SANITIZATION
// ============================================================================

func sanitizeMarkdownForTerminal(s string) string {
	if idx := strings.Index(s, "__USAGE__"); idx != -1 {
		s = s[:idx]
	}
	if idx := strings.Index(s, "__CONTENT_END__"); idx != -1 {
		s = s[:idx]
	}
	if idx := strings.Index(s, "\"completion_tokens\""); idx != -1 {
		s = s[:idx]
	}

	lines := strings.Split(s, "\n")
	var cleaned []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "{\"step\"") || strings.HasPrefix(trimmed, "{\"status\"") {
			continue
		}
		if strings.Contains(trimmed, "{\"step\":") && strings.Contains(trimmed, "\"status\":") {
			if idx := strings.Index(line, "{\"step\":"); idx >= 0 {
				endIdx := strings.Index(line[idx:], "}")
				if endIdx >= 0 {
					line = line[:idx] + line[idx+endIdx+1:]
				}
			}
		}
		cleaned = append(cleaned, line)
	}
	s = strings.Join(cleaned, "\n")
	s = strings.TrimRight(s, " \n\r\t")
	s = strings.ReplaceAll(s, "**", "")
	return s
}

func formatMarkdownForTerminal(s string) string {
	lines := strings.Split(s, "\n")
	var result []string
	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			continue
		}
		if strings.HasPrefix(line, "# ") {
			line = strings.TrimPrefix(line, "# ")
		} else if strings.HasPrefix(line, "## ") {
			line = strings.TrimPrefix(line, "## ")
		} else if strings.HasPrefix(line, "### ") {
			line = strings.TrimPrefix(line, "### ")
		}
		if strings.HasPrefix(line, "- ") {
			line = "- " + strings.TrimPrefix(line, "- ")
		}
		if strings.HasPrefix(line, "* ") {
			line = "- " + strings.TrimPrefix(line, "* ")
		}
		line = strings.ReplaceAll(line, "**", "")
		line = strings.ReplaceAll(line, "__", "")
		line = strings.ReplaceAll(line, "`", "")
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

// ============================================================================
// BROWSER HELPERS
// ============================================================================

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

// ============================================================================
// RUN INTERACTIVE — entry point
// ============================================================================

func RunInteractive() (string, error) {
	cfg, _ := config.LoadConfig()
	if cfg == nil {
		cfg = &config.Config{}
	}
	p := tea.NewProgram(NewInteractiveModel(cfg), tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}
	m, ok := finalModel.(model)
	if !ok {
		return "", nil
	}
	return m.GetSelectedCommand(), nil
}

// ============================================================================
// EXECUTE LOGIN (standalone, backward compat)
// ============================================================================

func ExecuteLogin(cfg *config.Config) error {
	if cfg == nil {
		cfg = &config.Config{}
	}
	if cfg.APIURL == "" {
		cfg.APIURL = "https://api.hopsule.com"
	}
	if cfg.WebURL == "" {
		cfg.WebURL = "https://app.hopsule.com"
	}
	client := api.NewClient(cfg)
	fmt.Println()
	fmt.Println("  Initializing login...")
	initResp, err := client.DeviceAuthInit("CLI")
	if err != nil {
		return fmt.Errorf("failed to initialize login: %w", err)
	}
	authURL := fmt.Sprintf("%s/auth/device?code=%s", cfg.WebURL, initResp.Code)
	fmt.Printf("  Device Code: %s\n\n", accentStyle.Render(initResp.Code))
	fmt.Println("  Opening browser to complete sign-in...")
	if err := openBrowser(authURL); err != nil {
		fmt.Println("  Could not open browser automatically.")
	}
	fmt.Println()
	fmt.Println("  If the browser doesn't open, visit:")
	fmt.Printf("  %s\n\n", logoStyle.Render(authURL))
	fmt.Println("  Waiting for authentication...")
	fmt.Println("  (Press Ctrl+C to cancel)")
	fmt.Println()

	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinnerIdx := 0
	for attempt := 0; attempt < 300; attempt++ {
		fmt.Printf("\r  %s Waiting for browser authentication... (%ds)",
			logoStyle.Render(spinner[spinnerIdx]), attempt*2)
		spinnerIdx = (spinnerIdx + 1) % len(spinner)
		resp, err := client.DeviceAuthPoll(initResp.Code)
		if err != nil {
			fmt.Println()
			return fmt.Errorf("failed to check login status: %w", err)
		}
		switch resp.Status {
		case "complete":
			fmt.Printf("\r  %s Authentication complete!                    \n\n", statusOnStyle.Render("v"))
			cfg.Token = resp.Token
			cfg.User = &config.User{ID: resp.UserID, Email: resp.Email, Name: resp.Name, AvatarURL: resp.AvatarURL}
			if err := config.SaveConfig(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			fmt.Printf("  Signed in as: %s (%s)\n\n", titleStyle.Render(resp.Name), dimStyle.Render(resp.Email))
			return nil
		case "expired":
			fmt.Println()
			return fmt.Errorf("login session expired - please try again")
		}
		time.Sleep(2 * time.Second)
	}
	fmt.Println()
	return fmt.Errorf("login timed out")
}

// Stubs so kanban.go public functions remain accessible.
func ShowOrganizations(cfg *config.Config) {}
func ShowProjects(cfg *config.Config)      {}
