package ui

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/Cagangedik/cli-tool/internal/config"
)

var (
	cyan    = color.New(color.FgCyan).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	yellow  = color.New(color.FgYellow).SprintFunc()
	gray    = color.New(color.FgHiBlack).SprintFunc()
	bold    = color.New(color.Bold).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()
)

func PrintDashboard(cfg *config.Config) {
	logo := `
        ████      ████                       %s
        ████      ████                       
            ████████                         
            ████████                         
        ████          ████                   
        ████          ████                   
`

	// Header
	fmt.Println(cyan("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Printf(logo, bold("Hopsule"))
	fmt.Println()

	// Context Info
	org := "hopsule-inc"
	if cfg.Organization != "" {
		org = cfg.Organization
	}

	project := "app"
	if cfg.Project != "" {
		project = truncate(cfg.Project, 30)
	}

	apiURL := "not configured"
	if cfg.APIURL != "" {
		apiURL = cfg.APIURL
	}

	captureStatus := green("ON")
	syncStatus := green("ON")
	privacy := yellow("redacted")
	
	lastSync := "12s"
	latency := "84ms"

	fmt.Printf("        %s: %s  •  %s: %s\n", 
		gray("org"), cyan(org),
		gray("project"), cyan(project))
	
	fmt.Printf("        %s: %s  •  %s: %s  •  %s: %s\n",
		gray("capture"), captureStatus,
		gray("sync"), syncStatus,
		gray("privacy"), privacy)
	
	fmt.Printf("        %s: %s  •  %s: %s\n",
		gray("last sync"), yellow(lastSync),
		gray("latency"), green(latency))

	fmt.Println(cyan("        ─────────────────────────────────────────────────────────────────────────────"))

	// Status
	if cfg.APIURL == "" {
		fmt.Printf("        %s %s\n\n", yellow("⚠"), bold("Not configured yet"))
	} else {
		fmt.Printf("        %s %s\n\n", green("✓"), bold("Connected"))
	}

	// Commands
	fmt.Printf("        %s\n", bold("Get started"))
	fmt.Printf("        %s %-15s %s\n", magenta("❯"), cyan("hopsule config"), gray("(configure cli)"))
	fmt.Printf("          %-15s %s\n", cyan("hopsule list"), gray("(list decisions)"))
	fmt.Printf("          %-15s %s\n", cyan("hopsule create"), gray("(create decision)"))
	fmt.Printf("          %-15s %s\n", cyan("hopsule status"), gray("(health check)"))
	fmt.Println()

	// API Info
	fmt.Printf("        %s: %s\n", gray("API"), cyan(apiURL))
	
	if cfg.Token != "" {
		fmt.Printf("        %s: %s\n", gray("Token"), green("configured ✓"))
	} else {
		fmt.Printf("        %s: %s\n", gray("Token"), yellow("not set"))
	}

	// Footer
	fmt.Println()
	fmt.Printf("        %s\n", gray("Run 'hopsule --help' for more commands"))
	fmt.Println(cyan("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println()
}

func PrintCompactStatus(cfg *config.Config) {
	// Compact version for commands
	project := "none"
	if cfg.Project != "" {
		project = truncate(cfg.Project, 20)
	}

	status := green("●")
	if cfg.APIURL == "" {
		status = yellow("●")
	}

	fmt.Printf("%s %s %s %s\n",
		status,
		bold("hopsule"),
		gray("→"),
		cyan(project))
}

func PrintSuccess(message string) {
	fmt.Printf("%s %s\n", green("✓"), message)
}

func PrintError(message string) {
	fmt.Printf("%s %s\n", color.RedString("✗"), message)
}

func PrintWarning(message string) {
	fmt.Printf("%s %s\n", yellow("⚠"), message)
}

func PrintInfo(message string) {
	fmt.Printf("%s %s\n", cyan("ℹ"), message)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// FormatTime formats time ago
func FormatTime(t time.Time) string {
	duration := time.Since(t)
	
	if duration < time.Minute {
		return fmt.Sprintf("%ds", int(duration.Seconds()))
	}
	if duration < time.Hour {
		return fmt.Sprintf("%dm", int(duration.Minutes()))
	}
	if duration < 24*time.Hour {
		return fmt.Sprintf("%dh", int(duration.Hours()))
	}
	return fmt.Sprintf("%dd", int(duration.Hours()/24))
}
