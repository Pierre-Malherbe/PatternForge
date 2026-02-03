package screens

import (
	"fmt"

	"github.com/Pierre-Malherbe/patternforge/internal/pattern"
	"github.com/Pierre-Malherbe/patternforge/internal/ui/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// SelectionMode determines what the selection screen is showing
type SelectionMode int

const (
	ModeCategories SelectionMode = iota
	ModePatterns
)

// CategoryItem wraps a category for the list component
type CategoryItem struct {
	name  string
	count int
}

func (i CategoryItem) Title() string {
	return fmt.Sprintf("📁 %s", i.name)
}
func (i CategoryItem) Description() string {
	return fmt.Sprintf("%d patterns", i.count)
}
func (i CategoryItem) FilterValue() string { return i.name }

// PatternItem wraps pattern for the list component
type PatternItem struct {
	pattern pattern.Pattern
}

func (i PatternItem) Title() string {
	title := i.pattern.DisplayTitle()
	if i.pattern.HasVariables() {
		title += " ⚙️"
	}
	return title
}
func (i PatternItem) Description() string {
	desc := i.pattern.Description
	if i.pattern.Source != "" && i.pattern.Source != "local" {
		desc += " [" + i.pattern.Source + "]"
	}
	return desc
}
func (i PatternItem) FilterValue() string { return i.pattern.Name }

// SelectionScreen shows the list of available patterns with category navigation
type SelectionScreen struct {
	mode             SelectionMode
	categoryList     list.Model
	patternList      list.Model
	allPatterns      []pattern.Pattern
	filteredPatterns []pattern.Pattern
	categories       []string
	currentCategory  string
	width            int
	height           int
}

// NewSelectionScreen creates the pattern selection screen
func NewSelectionScreen(patterns []pattern.Pattern) SelectionScreen {
	categories := pattern.GetCategories(patterns)

	s := SelectionScreen{
		mode:        ModeCategories,
		allPatterns: patterns,
		categories:  categories,
	}

	// If only one category or no categories, skip category view
	if len(categories) <= 1 {
		s.mode = ModePatterns
		s.filteredPatterns = patterns
		s.currentCategory = ""
		s.patternList = s.createPatternList(patterns, "All Patterns")
	} else {
		s.categoryList = s.createCategoryList(patterns, categories)
	}

	return s
}

func (s SelectionScreen) createCategoryList(patterns []pattern.Pattern, categories []string) list.Model {
	items := make([]list.Item, len(categories)+1)

	// Add "All" option first
	items[0] = CategoryItem{name: "All", count: len(patterns)}

	for i, cat := range categories {
		count := len(pattern.FilterByCategory(patterns, cat))
		items[i+1] = CategoryItem{name: cat, count: count}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "🎨 PatternForge - Select a category"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	return l
}

func (s SelectionScreen) createPatternList(patterns []pattern.Pattern, title string) list.Model {
	items := make([]list.Item, len(patterns))
	for i, p := range patterns {
		items[i] = PatternItem{pattern: p}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = fmt.Sprintf("🎨 %s", title)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	return l
}

func (s SelectionScreen) Init() tea.Cmd {
	return nil
}

func (s SelectionScreen) Update(msg tea.Msg) (SelectionScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		// Only set size on the list that's initialized
		if len(s.categories) > 1 {
			s.categoryList.SetSize(msg.Width-4, msg.Height-4)
		}
		if s.mode == ModePatterns {
			s.patternList.SetSize(msg.Width-4, msg.Height-4)
		}
	}

	var cmd tea.Cmd
	if s.mode == ModeCategories && len(s.categories) > 1 {
		s.categoryList, cmd = s.categoryList.Update(msg)
	} else if s.mode == ModePatterns {
		s.patternList, cmd = s.patternList.Update(msg)
	}
	return s, cmd
}

func (s SelectionScreen) View() string {
	var help string
	if s.mode == ModeCategories {
		help = styles.Help.Render("↑/↓: navigate • enter: select • U: upgrade • S: settings • q: quit")
		return s.categoryList.View() + "\n" + help
	}

	// Pattern mode
	backHint := ""
	if len(s.categories) > 1 {
		backHint = "esc: back • "
	}
	help = styles.Help.Render(fmt.Sprintf("↑/↓: navigate • enter: select • %sm: modify • n: new • U: upgrade • S: settings • q: quit", backHint))
	return s.patternList.View() + "\n" + help
}

// HandleEnter processes the enter key and returns true if a category was selected
func (s *SelectionScreen) HandleEnter() bool {
	if s.mode == ModeCategories {
		item := s.categoryList.SelectedItem()
		if item == nil {
			return true
		}

		catItem := item.(CategoryItem)
		s.currentCategory = catItem.name

		if catItem.name == "All" {
			s.filteredPatterns = s.allPatterns
			s.patternList = s.createPatternList(s.allPatterns, "All Patterns")
		} else {
			s.filteredPatterns = pattern.FilterByCategory(s.allPatterns, catItem.name)
			s.patternList = s.createPatternList(s.filteredPatterns, catItem.name)
		}

		// Update list size
		if s.width > 0 && s.height > 0 {
			s.patternList.SetSize(s.width-4, s.height-4)
		}

		s.mode = ModePatterns
		return true
	}
	return false
}

// HandleEsc processes the esc key and returns true if going back to categories
func (s *SelectionScreen) HandleEsc() bool {
	if s.mode == ModePatterns && len(s.categories) > 1 {
		s.mode = ModeCategories
		return true
	}
	return false
}

// SelectedPattern returns the currently selected pattern
func (s SelectionScreen) SelectedPattern() (pattern.Pattern, bool) {
	if s.mode == ModeCategories {
		return pattern.Pattern{}, false
	}

	item := s.patternList.SelectedItem()
	if item == nil {
		return pattern.Pattern{}, false
	}
	return item.(PatternItem).pattern, true
}

// IsInPatternMode returns true if viewing patterns (not categories)
func (s SelectionScreen) IsInPatternMode() bool {
	return s.mode == ModePatterns
}

// IsInCategoryMode returns true if viewing categories
func (s SelectionScreen) IsInCategoryMode() bool {
	return s.mode == ModeCategories
}
