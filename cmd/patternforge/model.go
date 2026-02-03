package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Pierre-Malherbe/patternforge/internal/claude"
	"github.com/Pierre-Malherbe/patternforge/internal/config"
	"github.com/Pierre-Malherbe/patternforge/internal/pattern"
	"github.com/Pierre-Malherbe/patternforge/internal/repository"
	"github.com/Pierre-Malherbe/patternforge/internal/ui/screens"
	"github.com/Pierre-Malherbe/patternforge/internal/ui/styles"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

// View represents the current screen state
type View int

const (
	SetupView View = iota
	SelectionView
	InputView
	ProcessingView
	ResultsView
	SettingsView
	UpgradeView
	VariablesView
)

// Model is the main application state
type Model struct {
	view        View
	selection   screens.SelectionScreen
	input       screens.InputScreen
	processing  screens.ProcessingScreen
	results     screens.ResultsScreen
	setup       screens.SetupScreen
	settings    screens.SettingsScreen
	upgrade     screens.UpgradeScreen
	variables   screens.VariablesScreen
	config      config.Config
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
func NewModel(patternsDir string, cfg config.Config, firstLaunch bool) (Model, error) {
	// Load patterns from directory and repositories
	patterns, err := LoadPatternsWithRepos(patternsDir, cfg)
	if err != nil {
		return Model{}, fmt.Errorf("failed to load patterns: %w", err)
	}

	if len(patterns) == 0 {
		return Model{}, fmt.Errorf("no patterns found in %s", patternsDir)
	}

	m := Model{
		view:        SelectionView,
		selection:   screens.NewSelectionScreen(patterns),
		patterns:    patterns,
		patternsDir: patternsDir,
		config:      cfg,
	}

	if firstLaunch {
		m.view = SetupView
		m.setup = screens.NewSetupScreen()
	}

	return m, nil
}

func (m Model) Init() tea.Cmd {
	if m.view == SetupView {
		return m.setup.Init()
	}
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
		patterns, _ := LoadPatternsWithRepos(m.patternsDir, m.config)
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

	case screens.SetupCompleteMsg:
		if err := config.Save(msg.Config); err != nil {
			m.err = err
			return m, nil
		}
		m.config = msg.Config

		// If community patterns enabled, trigger upgrade to clone repos
		if msg.EnableCommunity && len(m.config.Repositories) > 0 {
			m.upgrade = screens.NewUpgradeScreen()
			m.view = UpgradeView
			return m, m.upgrade.Init()
		}

		m.view = SelectionView
		if m.width > 0 && m.height > 0 {
			m.selection, _ = m.selection.Update(tea.WindowSizeMsg{
				Width:  m.width,
				Height: m.height,
			})
		}
		return m, nil

	case screens.SettingsSavedMsg:
		if err := config.Save(msg.Config); err != nil {
			m.err = err
			return m, nil
		}
		m.config = msg.Config
		m.view = SelectionView
		return m, nil

	case screens.SettingsCancelledMsg:
		m.view = SelectionView
		return m, nil

	case screens.UpgradeStartMsg:
		return m, m.runUpgrade()

	case screens.UpgradeCompleteMsg:
		m.upgrade, _ = m.upgrade.Update(msg)
		// Update config with new timestamps
		for i := range m.config.Repositories {
			for _, r := range msg.Results {
				if m.config.Repositories[i].Name == r.Name && r.Success {
					m.config.Repositories[i].LastUpdate = repository.FormatLastUpdate()
				}
			}
		}
		config.Save(m.config)
		return m, nil
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
		if m.view == SelectionView {
			// Try to go back to category view
			if m.selection.HandleEsc() {
				return m, nil
			}
		}
		if m.view == InputView {
			m.view = SelectionView
			return m, nil
		}
		if m.view == VariablesView {
			m.view = SelectionView
			return m, nil
		}
		if m.view == ResultsView {
			m.view = SelectionView
			m.input = screens.InputScreen{} // Reset input
			return m, nil
		}
		if m.view == SettingsView {
			m.view = SelectionView
			return m, nil
		}
		if m.view == UpgradeView && m.upgrade.IsCompleted() {
			// Reload patterns after upgrade
			patterns, _ := LoadPatternsWithRepos(m.patternsDir, m.config)
			m.patterns = patterns
			m.selection = screens.NewSelectionScreen(patterns)
			if m.width > 0 && m.height > 0 {
				m.selection, _ = m.selection.Update(tea.WindowSizeMsg{
					Width:  m.width,
					Height: m.height,
				})
			}
			m.view = SelectionView
			return m, nil
		}

	case "enter":
		if m.view == SelectionView {
			// First check if we're in category mode
			if m.selection.HandleEnter() {
				// Category was selected, stay in selection view
				return m, nil
			}

			// We're in pattern mode, select the pattern
			p, ok := m.selection.SelectedPattern()
			if ok {
				if p.HasVariables() {
					// Pattern has variables, show variable form
					m.variables = screens.NewVariablesScreen(p)
					if m.width > 0 && m.height > 0 {
						m.variables, _ = m.variables.Update(tea.WindowSizeMsg{
							Width:  m.width,
							Height: m.height,
						})
					}
					m.view = VariablesView
					return m, m.variables.Init()
				}
				m.input = screens.NewInputScreen(p)
				// Forward window size to input screen
				if m.width > 0 && m.height > 0 {
					m.input, _ = m.input.Update(tea.WindowSizeMsg{
						Width:  m.width,
						Height: m.height,
					})
				}
				m.view = InputView
				return m, m.input.Init()
			}
		}

		if m.view == VariablesView {
			// Check if variables form is complete
			if m.variables.IsComplete() {
				p := m.variables.Pattern()
				vars := m.variables.Values()
				m.input = screens.NewInputScreenWithVariables(p, vars)
				// Forward window size to input screen
				if m.width > 0 && m.height > 0 {
					m.input, _ = m.input.Update(tea.WindowSizeMsg{
						Width:  m.width,
						Height: m.height,
					})
				}
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

	case "s":
		if m.view == ResultsView {
			return m.saveToFile()
		}

	case "S":
		if m.view == SelectionView {
			m.settings = screens.NewSettingsScreen(m.config)
			m.view = SettingsView
			return m, m.settings.Init()
		}

	case "U":
		if m.view == SelectionView {
			return m.startUpgrade()
		}
	}

	return m.updateCurrentScreen(msg)
}

func (m Model) updateCurrentScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.view {
	case SetupView:
		m.setup, cmd = m.setup.Update(msg)
	case SelectionView:
		m.selection, cmd = m.selection.Update(msg)
	case InputView:
		m.input, cmd = m.input.Update(msg)
	case ProcessingView:
		m.processing, cmd = m.processing.Update(msg)
	case ResultsView:
		m.results, cmd = m.results.Update(msg)
	case SettingsView:
		m.settings, cmd = m.settings.Update(msg)
	case UpgradeView:
		m.upgrade, cmd = m.upgrade.Update(msg)
	case VariablesView:
		m.variables, cmd = m.variables.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	if m.err != nil {
		return styles.Error.Render(fmt.Sprintf("Error: %v\n\nPress q to quit", m.err))
	}

	switch m.view {
	case SetupView:
		return m.setup.View()
	case SelectionView:
		return m.selection.View()
	case InputView:
		return m.input.View()
	case ProcessingView:
		return m.processing.View()
	case ResultsView:
		return m.results.View()
	case SettingsView:
		return m.settings.View()
	case UpgradeView:
		return m.upgrade.View()
	case VariablesView:
		return m.variables.View()
	}

	return ""
}

// startProcessing switches to processing view and calls Claude
func (m Model) startProcessing() (tea.Model, tea.Cmd) {
	p, _ := m.selection.SelectedPattern()
	userInput := m.input.Value()
	variables := m.input.Variables()

	m.processing = screens.NewProcessingScreen()
	m.view = ProcessingView

	return m, tea.Batch(
		m.processing.Init(),
		m.callClaude(p, userInput, variables),
	)
}

// callClaude executes the prompt via Claude Code
func (m Model) callClaude(p pattern.Pattern, userInput string, variables map[string]string) tea.Cmd {
	return func() tea.Msg {
		var fullPrompt string
		if len(variables) > 0 {
			fullPrompt = p.RenderWithVariables(userInput, variables)
		} else {
			fullPrompt = p.Render(userInput)
		}
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

// saveToFile writes the result to a timestamped markdown file in the configured save directory
func (m Model) saveToFile() (tea.Model, tea.Cmd) {
	saveDir := m.config.SaveDirectory
	if saveDir == "" {
		saveDir = config.DefaultSaveDirectory()
	}

	if err := os.MkdirAll(saveDir, 0755); err != nil {
		m.err = fmt.Errorf("failed to create save directory: %w", err)
		return m, nil
	}

	filename := fmt.Sprintf("result-%s.md", time.Now().Format("2006-01-02-150405"))
	fullPath := filepath.Join(saveDir, filename)

	if err := os.WriteFile(fullPath, []byte(m.results.Result()), 0644); err != nil {
		m.err = err
		return m, nil
	}
	m.results.SetSaved(fullPath)
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

// startUpgrade switches to upgrade view and starts the sync process
func (m Model) startUpgrade() (tea.Model, tea.Cmd) {
	if len(m.config.Repositories) == 0 {
		m.err = fmt.Errorf("no repositories configured. Use 'patternforge repo add <url>' to add one")
		return m, nil
	}

	if !repository.IsGitAvailable() {
		m.err = fmt.Errorf("git is not installed. Cannot upgrade repositories")
		return m, nil
	}

	m.upgrade = screens.NewUpgradeScreen()
	m.view = UpgradeView

	return m, m.upgrade.Init()
}

// runUpgrade executes the repository sync in the background
func (m Model) runUpgrade() tea.Cmd {
	return func() tea.Msg {
		// Convert config.Repository to repository.Repository
		repos := make([]repository.Repository, len(m.config.Repositories))
		for i, r := range m.config.Repositories {
			repos[i] = repository.Repository{
				Name:       r.Name,
				URL:        r.URL,
				Enabled:    r.Enabled,
				LastUpdate: r.LastUpdate,
			}
		}

		mgr := repository.NewManager(config.RepositoriesDir(), repos)
		results := mgr.UpdateAll()

		return screens.UpgradeCompleteMsg{Results: results}
	}
}
