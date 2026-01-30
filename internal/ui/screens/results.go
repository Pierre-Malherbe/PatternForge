package screens

import (
	"fmt"

	"github.com/Pierre-Malherbe/patternforge/internal/claude"
	"github.com/Pierre-Malherbe/patternforge/internal/ui/styles"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// ResultsScreen displays Claude's output with stats
type ResultsScreen struct {
	viewport viewport.Model
	result   string
	stats    claude.Stats
	width    int
	height   int
	copied    bool   // Show "copied!" message
	saved     bool   // Show "saved!" message
	savedPath string // Path of saved file
}

// NewResultsScreen creates the results display screen
func NewResultsScreen(result string, stats claude.Stats) ResultsScreen {
	vp := viewport.New(0, 0)
	vp.SetContent(result)

	return ResultsScreen{
		viewport: vp,
		result:   result,
		stats:    stats,
	}
}

func (s ResultsScreen) Init() tea.Cmd {
	return nil
}

func (s ResultsScreen) Update(msg tea.Msg) (ResultsScreen, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.viewport.Width = msg.Width - 6
		s.viewport.Height = msg.Height - 15
		// Re-set content after resize to ensure it's rendered
		s.viewport.SetContent(s.result)
	}

	s.viewport, cmd = s.viewport.Update(msg)
	return s, cmd
}

func (s ResultsScreen) View() string {
	title := styles.Title.Render("✨ Result")

	statsContent := fmt.Sprintf(
		"📊 Tokens: %d input + %d output = %d total | ⏱️  Duration: %s",
		s.stats.InputTokens,
		s.stats.OutputTokens,
		s.stats.TotalTokens,
		s.stats.Duration.Round(1),
	)
	statsBox := styles.Stats.Render(statsContent)

	help := styles.Help.Render("\n↑/↓ or j/k: scroll • y: copy • s: save • Esc: new pattern • q: quit")

	statusMsg := ""
	if s.copied {
		statusMsg = styles.Success.Render("✅ Copied to clipboard!\n")
	} else if s.saved {
		statusMsg = styles.Success.Render(fmt.Sprintf("✅ Saved to %s!\n", s.savedPath))
	}

	return fmt.Sprintf("%s\n\n%s\n%s%s%s",
		title,
		styles.Result.Render(s.viewport.View()),
		statsBox,
		statusMsg,
		help,
	)
}

// Result returns the output text
func (s ResultsScreen) Result() string {
	return s.result
}

// SetCopied shows the "copied" confirmation message
func (s *ResultsScreen) SetCopied(copied bool) {
	s.copied = copied
	s.saved = false
}

// SetSaved shows the "saved" confirmation message with the file path
func (s *ResultsScreen) SetSaved(path string) {
	s.saved = true
	s.savedPath = path
	s.copied = false
}
