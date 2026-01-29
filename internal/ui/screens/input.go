package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/patternforge/patternforge/internal/pattern"
	"github.com/patternforge/patternforge/internal/ui/styles"
)

// InputScreen allows user to paste/type their content
type InputScreen struct {
	textarea textarea.Model
	pattern  pattern.Pattern
	width    int
	height   int
}

// NewInputScreen creates the input screen for a pattern
func NewInputScreen(p pattern.Pattern) InputScreen {
	ta := textarea.New()
	ta.Placeholder = "Paste your content here..."
	ta.Focus()
	ta.CharLimit = 50000

	return InputScreen{
		textarea: ta,
		pattern:  p,
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
		s.textarea.SetWidth(msg.Width - 6)
		s.textarea.SetHeight(msg.Height - 12)
	}

	s.textarea, cmd = s.textarea.Update(msg)
	return s, cmd
}

func (s InputScreen) View() string {
	title := styles.Title.Render(s.pattern.DisplayTitle())
	desc := styles.Subtitle.Render(s.pattern.Description)
	charCount := styles.Subtitle.Render(fmt.Sprintf("Characters: %d/50000", len(s.textarea.Value())))
	help := styles.Help.Render("\nCtrl+D: process • Esc: back • m: edit pattern")

	return fmt.Sprintf("%s\n%s\n\n%s\n%s\n%s",
		title,
		desc,
		s.textarea.View(),
		charCount,
		help,
	)
}

// Value returns the user's input
func (s InputScreen) Value() string {
	return s.textarea.Value()
}
