package screens

import (
	"fmt"

	"github.com/Pierre-Malherbe/patternforge/internal/config"
	"github.com/Pierre-Malherbe/patternforge/internal/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// SetupCompleteMsg is sent when the user finishes first-launch setup.
type SetupCompleteMsg struct {
	Config config.Config
}

// SetupScreen is shown on first launch to configure the save directory.
type SetupScreen struct {
	textInput textinput.Model
	width     int
	height    int
}

// NewSetupScreen creates the first-launch setup screen.
func NewSetupScreen() SetupScreen {
	ti := textinput.New()
	ti.Placeholder = config.DefaultSaveDirectory()
	ti.SetValue(config.DefaultSaveDirectory())
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60

	return SetupScreen{
		textInput: ti,
	}
}

func (s SetupScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (s SetupScreen) Update(msg tea.Msg) (SetupScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.textInput.Width = msg.Width - 10
		if s.textInput.Width > 80 {
			s.textInput.Width = 80
		}
	case tea.KeyMsg:
		if msg.String() == "enter" {
			val := s.textInput.Value()
			if val == "" {
				val = config.DefaultSaveDirectory()
			}
			return s, func() tea.Msg {
				return SetupCompleteMsg{Config: config.Config{SaveDirectory: val}}
			}
		}
	}

	var cmd tea.Cmd
	s.textInput, cmd = s.textInput.Update(msg)
	return s, cmd
}

func (s SetupScreen) View() string {
	title := styles.Title.Render("Welcome to PatternForge!")
	subtitle := styles.Subtitle.Render("Where should results be saved?")
	help := styles.Help.Render("\nenter: confirm")

	return fmt.Sprintf("\n\n%s\n\n%s\n\n  %s\n%s",
		title,
		subtitle,
		s.textInput.View(),
		help,
	)
}
