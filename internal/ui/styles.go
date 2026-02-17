package ui

import "github.com/charmbracelet/lipgloss"

// ============================================================================
// ADAPTIVE COLORS — work on both light and dark terminals
// ============================================================================

var (
	cyanColor    = lipgloss.AdaptiveColor{Light: "30", Dark: "51"}
	greenColor   = lipgloss.AdaptiveColor{Light: "28", Dark: "46"}
	yellowColor  = lipgloss.AdaptiveColor{Light: "136", Dark: "226"}
	magentaColor = lipgloss.AdaptiveColor{Light: "127", Dark: "201"}
	purpleColor  = lipgloss.AdaptiveColor{Light: "93", Dark: "141"}
	dimColor = lipgloss.AdaptiveColor{Light: "244", Dark: "240"}
	blueColor    = lipgloss.AdaptiveColor{Light: "25", Dark: "39"}
	redColor     = lipgloss.AdaptiveColor{Light: "160", Dark: "196"}
)

// ============================================================================
// TEXT STYLES
// ============================================================================

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "255"}).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "255"}).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "240", Dark: "250"})

	logoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "238", Dark: "252"})

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "244", Dark: "244"})

	accentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "238", Dark: "252"})

	versionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "248", Dark: "240"})

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "244", Dark: "244"})

	statusOnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "28", Dark: "46"}).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "244", Dark: "244"})

)

// ============================================================================
// TABLE / LIST STYLES
// ============================================================================

var (
	tableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "255"}).
				BorderStyle(lipgloss.NormalBorder()).
				BorderBottom(true).
				BorderForeground(lipgloss.AdaptiveColor{Light: "240", Dark: "237"}).
				PaddingRight(2)

	tableSelectedRowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "255"}).
				Bold(true)
)

// ============================================================================
// SCROLL AREA STYLES
// ============================================================================

var (
	scrollAreaBorderColor = lipgloss.AdaptiveColor{Light: "240", Dark: "238"}

	scrollBarTrackColor  = lipgloss.AdaptiveColor{Light: "254", Dark: "236"}
	scrollBarActiveColor = lipgloss.AdaptiveColor{Light: "28", Dark: "46"}
)
