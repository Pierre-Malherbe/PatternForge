package screens

import (
	"fmt"

	"github.com/Pierre-Malherbe/patternforge/internal/config"
	"github.com/Pierre-Malherbe/patternforge/internal/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// SettingsSavedMsg is sent when the user saves settings.
type SettingsSavedMsg struct {
	Config config.Config
}

// SettingsCancelledMsg is sent when the user cancels settings.
type SettingsCancelledMsg struct{}

// SettingsScreen allows the user to edit application settings.
type SettingsScreen struct {
	textInput textinput.Model
	width     int
	height    int
}

// NewSettingsScreen creates the settings screen with current config values.
func NewSettingsScreen(cfg config.Config) SettingsScreen {
	ti := textinput.New()
	ti.Placeholder = config.DefaultSaveDirectory()
	ti.SetValue(cfg.SaveDirectory)
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60

	return SettingsScreen{
		textInput: ti,
	}
}

func (s SettingsScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (s SettingsScreen) Update(msg tea.Msg) (SettingsScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.textInput.Width = msg.Width - 10
		if s.textInput.Width > 80 {
			s.textInput.Width = 80
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			val := s.textInput.Value()
			if val == "" {
				val = config.DefaultSaveDirectory()
			}
			return s, func() tea.Msg {
				return SettingsSavedMsg{Config: config.Config{SaveDirectory: val}}
			}
		case "esc":
			return s, func() tea.Msg {
				return SettingsCancelledMsg{}
			}
		}
	}

	var cmd tea.Cmd
	s.textInput, cmd = s.textInput.Update(msg)
	return s, cmd
}

func (s SettingsScreen) View() string {
	title := styles.Title.Render("Settings")
	label := styles.Subtitle.Render("Save directory:")
	help := styles.Help.Render("\nenter: save • esc: cancel")

	return fmt.Sprintf("\n\n%s\n\n%s\n\n  %s\n%s",
		title,
		label,
		s.textInput.View(),
		help,
	)
}
