package screens

import (
	"fmt"
	"time"

	"github.com/Pierre-Malherbe/patternforge/internal/repository"
	"github.com/Pierre-Malherbe/patternforge/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// UpgradeStartMsg triggers the upgrade process.
type UpgradeStartMsg struct{}

// UpgradeCompleteMsg is sent when upgrade finishes.
type UpgradeCompleteMsg struct {
	Results []repository.UpdateResult
}

// UpgradeScreen shows sync progress for repositories.
type UpgradeScreen struct {
	spinner   string
	frame     int
	results   []repository.UpdateResult
	completed bool
	width     int
	height    int
}

var upgradeSpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// NewUpgradeScreen creates the upgrade/sync screen.
func NewUpgradeScreen() UpgradeScreen {
	return UpgradeScreen{
		spinner: upgradeSpinnerFrames[0],
		frame:   0,
	}
}

// UpgradeTickMsg is sent periodically to animate the spinner.
type UpgradeTickMsg time.Time

func (s UpgradeScreen) Init() tea.Cmd {
	return tea.Batch(upgradeTick(), func() tea.Msg { return UpgradeStartMsg{} })
}

func (s UpgradeScreen) Update(msg tea.Msg) (UpgradeScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
	case UpgradeTickMsg:
		if !s.completed {
			s.frame = (s.frame + 1) % len(upgradeSpinnerFrames)
			s.spinner = upgradeSpinnerFrames[s.frame]
			return s, upgradeTick()
		}
	case UpgradeCompleteMsg:
		s.completed = true
		s.results = msg.Results
	}
	return s, nil
}

func (s UpgradeScreen) View() string {
	if !s.completed {
		title := styles.Title.Render("⬆️  Upgrading Repositories...")
		message := fmt.Sprintf("\n   %s  Syncing patterns...", s.spinner)
		return title + message
	}

	title := styles.Title.Render("⬆️  Upgrade Complete")
	var results string

	if len(s.results) == 0 {
		results = "\n   No repositories configured."
	} else {
		results = "\n"
		for _, r := range s.results {
			if r.Success {
				results += fmt.Sprintf("   %s %s: %s\n", styles.Success.Render("✓"), r.Name, r.Message)
			} else {
				results += fmt.Sprintf("   %s %s: %s\n", styles.Error.Render("✗"), r.Name, r.Message)
			}
		}
	}

	help := styles.Help.Render("\n   Press Esc to return")
	return title + results + help
}

// IsCompleted returns whether the upgrade has finished.
func (s UpgradeScreen) IsCompleted() bool {
	return s.completed
}

// Results returns the upgrade results.
func (s UpgradeScreen) Results() []repository.UpdateResult {
	return s.results
}

func upgradeTick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return UpgradeTickMsg(t)
	})
}
