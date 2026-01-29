package styles

import "github.com/charmbracelet/lipgloss"

// Color palette
const (
	ColorPrimary   = lipgloss.Color("205") // Pink
	ColorSecondary = lipgloss.Color("170") // Purple
	ColorSuccess   = lipgloss.Color("120") // Green
	ColorInfo      = lipgloss.Color("63")  // Blue
	ColorMuted     = lipgloss.Color("241") // Gray
	ColorBorder    = lipgloss.Color("238") // Dark gray
)

// Title is for main headings
var Title = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorPrimary).
	MarginLeft(2)

// Subtitle is for secondary text
var Subtitle = lipgloss.NewStyle().
	Foreground(ColorMuted).
	MarginLeft(2)

// Selected item style
var Selected = lipgloss.NewStyle().
	Foreground(ColorSecondary).
	Bold(true)

// Help text at bottom of screens
var Help = lipgloss.NewStyle().
	Foreground(ColorMuted).
	Italic(true)

// Stats box for token info
var Stats = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(ColorInfo).
	Padding(0, 1).
	MarginTop(1)

// Result output box
var Result = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(ColorSuccess).
	Padding(1).
	MarginTop(1)

// Error messages
var Error = lipgloss.NewStyle().
	Foreground(lipgloss.Color("196")). // Red
	Bold(true)

// Success messages
var Success = lipgloss.NewStyle().
	Foreground(ColorSuccess).
	Bold(true)
