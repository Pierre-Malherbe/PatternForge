package screens

import (
	"fmt"
	"strings"

	"github.com/Pierre-Malherbe/patternforge/internal/pattern"
	"github.com/Pierre-Malherbe/patternforge/internal/ui/styles"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

// InputScreen allows user to paste/type their content
type InputScreen struct {
	textarea  textarea.Model
	pattern   pattern.Pattern
	variables map[string]string
	width     int
	height    int
}

// NewInputScreen creates the input screen for a pattern
func NewInputScreen(p pattern.Pattern) InputScreen {
	ta := textarea.New()
	ta.Placeholder = "Paste your content here..."
	ta.Focus()
	ta.CharLimit = 50000

	return InputScreen{
		textarea:  ta,
		pattern:   p,
		variables: make(map[string]string),
	}
}

// NewInputScreenWithVariables creates the input screen with pre-filled variables
func NewInputScreenWithVariables(p pattern.Pattern, variables map[string]string) InputScreen {
	ta := textarea.New()
	ta.Placeholder = "Paste your content here..."
	ta.Focus()
	ta.CharLimit = 50000

	return InputScreen{
		textarea:  ta,
		pattern:   p,
		variables: variables,
	}
}

func (s InputScreen) Init() tea.Cmd {
	return textarea.Blink
}

func (s InputScreen) Update(msg tea.Msg) (InputScreen, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		// Maximize textarea size - compact layout uses minimal chrome
		s.textarea.SetWidth(msg.Width - 2)
		// Reserve: 1 header line + 1 blank + 1 footer line = 3 lines
		s.textarea.SetHeight(msg.Height - 4)
	}

	s.textarea, cmd = s.textarea.Update(msg)
	return s, cmd
}

func (s InputScreen) View() string {
	// Compact header: title + description on same area
	header := styles.Title.Render(s.pattern.DisplayTitle())
	if s.pattern.Description != "" {
		header += " " + styles.Subtitle.Render("- "+s.pattern.Description)
	}

	// Show configured variables inline if any
	if len(s.variables) > 0 {
		var varParts []string
		for name, value := range s.variables {
			varParts = append(varParts, fmt.Sprintf("%s:%s", name, value))
		}
		header += " " + styles.Subtitle.Render("["+strings.Join(varParts, " | ")+"]")
	}

	// Compact footer
	charCount := styles.Subtitle.Render(fmt.Sprintf("%d chars", len(s.textarea.Value())))
	help := styles.Help.Render("Ctrl+D: process • Esc: back")

	return fmt.Sprintf("%s\n\n%s\n%s  %s",
		header,
		s.textarea.View(),
		charCount,
		help,
	)
}

// Value returns the user's input
func (s InputScreen) Value() string {
	return s.textarea.Value()
}

// Variables returns the configured variables
func (s InputScreen) Variables() map[string]string {
	return s.variables
}
