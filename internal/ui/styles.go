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
	grayColor    = lipgloss.AdaptiveColor{Light: "244", Dark: "244"}
	dimColor     = lipgloss.AdaptiveColor{Light: "244", Dark: "240"}
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

	breadcrumbStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "238", Dark: "252"})
)

// ============================================================================
// CARD STYLES
// ============================================================================

var (
	cardBorderColor         = lipgloss.AdaptiveColor{Light: "240", Dark: "240"}
	cardSelectedBorderColor = lipgloss.AdaptiveColor{Light: "28", Dark: "46"}

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(cardBorderColor).
			Padding(1, 2).
			Width(32).
			MarginRight(2)

	selectedCardStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(cardSelectedBorderColor).
				Padding(1, 2).
				Width(32).
				MarginRight(2)

	cardTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "255"}).
			Bold(true)

	cardDescStyle = lipgloss.NewStyle().
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

	tableRowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "240", Dark: "250"})

	tableSelectedRowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "255"}).
				Bold(true)

	tableRowAltBg = lipgloss.AdaptiveColor{Light: "254", Dark: "235"}
)

// ============================================================================
// SCROLL AREA STYLES
// ============================================================================

var (
	scrollAreaBorderColor = lipgloss.AdaptiveColor{Light: "240", Dark: "238"}

	scrollAreaStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(scrollAreaBorderColor).
			Padding(0, 1)

	scrollBarTrackColor  = lipgloss.AdaptiveColor{Light: "254", Dark: "236"}
	scrollBarThumbColor  = lipgloss.AdaptiveColor{Light: "244", Dark: "245"}
	scrollBarActiveColor = lipgloss.AdaptiveColor{Light: "28", Dark: "46"}
)
