package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/patternforge/patternforge/internal/claude"
	"github.com/patternforge/patternforge/internal/pattern"
	"github.com/patternforge/patternforge/internal/ui/screens"
	"github.com/patternforge/patternforge/internal/ui/styles"
)

// View represents the current screen state
type View int

const (
	SelectionView View = iota
	InputView
	ProcessingView
	ResultsView
)

// Model is the main application state
type Model struct {
	view        View
	selection   screens.SelectionScreen
	input       screens.InputScreen
	processing  screens.ProcessingScreen
	results     screens.ResultsScreen
	patterns    []pattern.Pattern
	patternsDir string
	err         error
	width       int
	height      int
}

// ProcessedMsg is sent when Claude finishes processing
type ProcessedMsg struct {
	result string
	stats  claude.Stats
	err    error
}

// ReloadPatternsMsg is sent after editing/creating patterns to reload the list
type ReloadPatternsMsg struct{}

// NewModel initializes the application
func NewModel(patternsDir string) (Model, error) {
	// Load patterns from directory
	patterns, err := pattern.LoadAll(patternsDir)
	if err != nil {
		return Model{}, fmt.Errorf("failed to load patterns: %w", err)
	}

	if len(patterns) == 0 {
		return Model{}, fmt.Errorf("no patterns found in %s", patternsDir)
	}

	return Model{
		view:        SelectionView,
		selection:   screens.NewSelectionScreen(patterns),
		patterns:    patterns,
		patternsDir: patternsDir,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case ReloadPatternsMsg:
		// Reload patterns from disk after Vi edit
		patterns, _ := pattern.LoadAll(m.patternsDir)
		m.patterns = patterns
		m.selection = screens.NewSelectionScreen(patterns)
		// Update dimensions of the new selection screen
		if m.width > 0 && m.height > 0 {
			m.selection, _ = m.selection.Update(tea.WindowSizeMsg{
				Width:  m.width,
				Height: m.height,
			})
		}
		return m, nil

	case ProcessedMsg:
		// Claude finished processing
		if msg.err != nil {
			m.err = msg.err
			m.view = SelectionView
			return m, nil
		}
		m.results = screens.NewResultsScreen(msg.result, msg.stats)
		// Initialize viewport with current window dimensions
		m.results, _ = m.results.Update(tea.WindowSizeMsg{
			Width:  m.width,
			Height: m.height,
		})
		m.view = ResultsView
		return m, m.results.Init()
	}

	// Delegate to current screen
	return m.updateCurrentScreen(msg)
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		if m.view == SelectionView {
			return m, tea.Quit
		}
		if m.view == ResultsView {
			return m, tea.Quit
		}

	case "esc":
		if m.view == InputView {
			m.view = SelectionView
			return m, nil
		}
		if m.view == ResultsView {
			m.view = SelectionView
			m.input = screens.InputScreen{} // Reset input
			return m, nil
		}

	case "enter":
		if m.view == SelectionView {
			p, ok := m.selection.SelectedPattern()
			if ok {
				m.input = screens.NewInputScreen(p)
				m.view = InputView
				return m, m.input.Init()
			}
		}

	case "ctrl+d":
		if m.view == InputView {
			return m.startProcessing()
		}

	case "m":
		if m.view == SelectionView {
			return m.editPattern()
		}

	case "n":
		if m.view == SelectionView {
			return m.createPattern()
		}

	case "y":
		if m.view == ResultsView {
			return m.copyToClipboard()
		}
	}

	return m.updateCurrentScreen(msg)
}

func (m Model) updateCurrentScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.view {
	case SelectionView:
		m.selection, cmd = m.selection.Update(msg)
	case InputView:
		m.input, cmd = m.input.Update(msg)
	case ProcessingView:
		m.processing, cmd = m.processing.Update(msg)
	case ResultsView:
		m.results, cmd = m.results.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	if m.err != nil {
		return styles.Error.Render(fmt.Sprintf("Error: %v\n\nPress q to quit", m.err))
	}

	switch m.view {
	case SelectionView:
		return m.selection.View()
	case InputView:
		return m.input.View()
	case ProcessingView:
		return m.processing.View()
	case ResultsView:
		return m.results.View()
	}

	return ""
}

// startProcessing switches to processing view and calls Claude
func (m Model) startProcessing() (tea.Model, tea.Cmd) {
	p, _ := m.selection.SelectedPattern()
	userInput := m.input.Value()

	m.processing = screens.NewProcessingScreen()
	m.view = ProcessingView

	return m, tea.Batch(
		m.processing.Init(),
		m.callClaude(p, userInput),
	)
}

// callClaude executes the prompt via Claude Code
func (m Model) callClaude(p pattern.Pattern, userInput string) tea.Cmd {
	return func() tea.Msg {
		fullPrompt := p.Render(userInput)
		result, stats, err := claude.Execute(fullPrompt)
		return ProcessedMsg{
			result: result,
			stats:  stats,
			err:    err,
		}
	}
}

// copyToClipboard copies result to system clipboard
func (m Model) copyToClipboard() (tea.Model, tea.Cmd) {
	if err := clipboard.WriteAll(m.results.Result()); err != nil {
		m.err = err
		return m, nil
	}
	m.results.SetCopied(true)
	return m, nil
}

// editPattern opens the selected pattern in Vi
func (m Model) editPattern() (tea.Model, tea.Cmd) {
	p, ok := m.selection.SelectedPattern()
	if !ok {
		return m, nil
	}

	return m, tea.ExecProcess(exec.Command("vi", p.FilePath), func(err error) tea.Msg {
		return ReloadPatternsMsg{}
	})
}

// createPattern creates a new pattern file and opens it in Vi
func (m Model) createPattern() (tea.Model, tea.Cmd) {
	template := `# 📝 New Pattern

> Short description of your pattern

## Prompt

Your instructions here.

Use {{input}} to insert user content at specific position.

Instructions after input...
`

	filename := fmt.Sprintf("%s/new-pattern-%d.md", m.patternsDir, time.Now().Unix())
	if err := os.WriteFile(filename, []byte(template), 0644); err != nil {
		m.err = err
		return m, nil
	}

	return m, tea.ExecProcess(exec.Command("vi", filename), func(err error) tea.Msg {
		return ReloadPatternsMsg{}
	})
}
