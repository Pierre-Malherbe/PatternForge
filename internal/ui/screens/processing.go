package screens

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/patternforge/patternforge/internal/ui/styles"
)

// ProcessingScreen shows a spinner while Claude processes
type ProcessingScreen struct {
	spinner string
	frame   int
}

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// NewProcessingScreen creates the processing/loading screen
func NewProcessingScreen() ProcessingScreen {
	return ProcessingScreen{
		spinner: spinnerFrames[0],
		frame:   0,
	}
}

// TickMsg is sent periodically to animate the spinner
type TickMsg time.Time

func (s ProcessingScreen) Init() tea.Cmd {
	return tick()
}

func (s ProcessingScreen) Update(msg tea.Msg) (ProcessingScreen, tea.Cmd) {
	switch msg.(type) {
	case TickMsg:
		s.frame = (s.frame + 1) % len(spinnerFrames)
		s.spinner = spinnerFrames[s.frame]
		return s, tick()
	}
	return s, nil
}

func (s ProcessingScreen) View() string {
	title := styles.Title.Render("⚙️  Processing with Claude Code...")
	message := fmt.Sprintf("\n   %s  Please wait...", s.spinner)
	return title + message
}

func tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
