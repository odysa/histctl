package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Browser colors
	SafariColor  = lipgloss.Color("#0A84FF")
	ChromeColor  = lipgloss.Color("#FF5F00")
	EdgeColor    = lipgloss.Color("#0078D4")
	FirefoxColor = lipgloss.Color("#FF7139")

	Subtle  = lipgloss.Color("#6C6C6C")
	Accent  = lipgloss.Color("#7D56F4")
	Danger  = lipgloss.Color("#FF4D4F")
	Success = lipgloss.Color("#52C41A")
	Muted   = lipgloss.Color("#4A4A4A")

	// Header
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(Accent).
			Padding(0, 1)

	// Browser pills
	BrowserPillBase = lipgloss.NewStyle().
			Padding(0, 1).
			MarginRight(1).
			Bold(true)

	BrowserPillActive = func(color lipgloss.Color) lipgloss.Style {
		return BrowserPillBase.
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(color)
	}

	BrowserPillInactive = func(color lipgloss.Color) lipgloss.Style {
		return BrowserPillBase.
			Foreground(color)
	}

	// Search bar
	SearchBarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Accent).
			Padding(0, 1).
			MarginTop(1).
			MarginBottom(1)

	SearchBarActiveStyle = SearchBarStyle.
				BorderForeground(lipgloss.Color("#FF79C6"))

	SearchLabelStyle = lipgloss.NewStyle().
				Foreground(Accent).
				Bold(true)

	SearchPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF79C6"))

	// Status bar
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(Subtle).
			MarginTop(1)

	SelectedCountStyle = lipgloss.NewStyle().
				Foreground(Danger).
				Bold(true)

	// Help
	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(Subtle)

	HelpSepStyle = lipgloss.NewStyle().
			Foreground(Muted)

	// Error
	ErrorStyle = lipgloss.NewStyle().
			Foreground(Danger).
			Bold(true)

	// Dialog
	DialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Danger).
			Padding(1, 3).
			Background(lipgloss.Color("#1A1A2E")).
			Foreground(lipgloss.Color("#FAFAFA"))
)

// BrowserColor returns the color for a browser name.
func BrowserColor(name string) lipgloss.Color {
	switch name {
	case "safari":
		return SafariColor
	case "chrome":
		return ChromeColor
	case "edge":
		return EdgeColor
	case "firefox":
		return FirefoxColor
	default:
		return Subtle
	}
}
