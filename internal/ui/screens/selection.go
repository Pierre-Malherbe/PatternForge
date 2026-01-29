package screens

import (
	"github.com/Pierre-Malherbe/patternforge/internal/pattern"
	"github.com/Pierre-Malherbe/patternforge/internal/ui/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// PatternItem wraps pattern for the list component
type PatternItem struct {
	pattern pattern.Pattern
}

func (i PatternItem) Title() string       { return i.pattern.DisplayTitle() }
func (i PatternItem) Description() string { return i.pattern.Description }
func (i PatternItem) FilterValue() string { return i.pattern.Name }

// SelectionScreen shows the list of available patterns
type SelectionScreen struct {
	list     list.Model
	patterns []pattern.Pattern
	width    int
	height   int
}

// NewSelectionScreen creates the pattern selection screen
func NewSelectionScreen(patterns []pattern.Pattern) SelectionScreen {
	items := make([]list.Item, len(patterns))
	for i, p := range patterns {
		items[i] = PatternItem{pattern: p}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "🎨 PatternForge - Select a pattern"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	return SelectionScreen{
		list:     l,
		patterns: patterns,
	}
}

func (s SelectionScreen) Init() tea.Cmd {
	return nil
}

func (s SelectionScreen) Update(msg tea.Msg) (SelectionScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.list.SetSize(msg.Width-4, msg.Height-4)
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s SelectionScreen) View() string {
	help := styles.Help.Render("↑/↓ or j/k: navigate • enter: select • m: modify • n: new • q: quit")
	return s.list.View() + "\n" + help
}

// SelectedPattern returns the currently selected pattern
func (s SelectionScreen) SelectedPattern() (pattern.Pattern, bool) {
	item := s.list.SelectedItem()
	if item == nil {
		return pattern.Pattern{}, false
	}
	return item.(PatternItem).pattern, true
}
