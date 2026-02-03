package screens

import (
	"fmt"
	"strings"

	"github.com/Pierre-Malherbe/patternforge/internal/config"
	"github.com/Pierre-Malherbe/patternforge/internal/repository"
	"github.com/Pierre-Malherbe/patternforge/internal/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// SetupCompleteMsg is sent when the user finishes first-launch setup.
type SetupCompleteMsg struct {
	Config          config.Config
	EnableCommunity bool
}

// SetupStep represents the current step in the setup wizard.
type SetupStep int

const (
	StepSaveDirectory SetupStep = iota
	StepCommunityPatterns
	StepRepoChoice
	StepCustomRepo
)

// SetupScreen is shown on first launch to configure settings.
type SetupScreen struct {
	step             SetupStep
	textInput        textinput.Model
	saveDirectory    string
	communityEnabled bool
	useOfficialRepo  bool
	customRepoURL    string
	width            int
	height           int
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
		step:      StepSaveDirectory,
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
		switch s.step {
		case StepSaveDirectory:
			if msg.String() == "enter" {
				val := s.textInput.Value()
				if val == "" {
					val = config.DefaultSaveDirectory()
				}
				s.saveDirectory = val

				// Check if git is available for community patterns
				if repository.IsGitAvailable() {
					s.step = StepCommunityPatterns
					return s, nil
				}

				// Skip community step if git not available
				return s, s.completeSetup()
			}

		case StepCommunityPatterns:
			switch msg.String() {
			case "y", "Y", "enter":
				s.communityEnabled = true
				s.step = StepRepoChoice
				return s, nil
			case "n", "N":
				s.communityEnabled = false
				return s, s.completeSetup()
			}

		case StepRepoChoice:
			switch msg.String() {
			case "o", "O", "enter":
				s.useOfficialRepo = true
				return s, s.completeSetup()
			case "c", "C":
				s.useOfficialRepo = false
				s.step = StepCustomRepo
				// Reset text input for custom URL
				s.textInput.SetValue("")
				s.textInput.Placeholder = "https://github.com/user/patterns-repo"
				s.textInput.Focus()
				return s, textinput.Blink
			}

		case StepCustomRepo:
			if msg.String() == "enter" {
				url := strings.TrimSpace(s.textInput.Value())
				if url != "" {
					s.customRepoURL = url
					return s, s.completeSetup()
				}
			}
			if msg.String() == "esc" {
				s.step = StepRepoChoice
				return s, nil
			}
		}
	}

	if s.step == StepSaveDirectory || s.step == StepCustomRepo {
		var cmd tea.Cmd
		s.textInput, cmd = s.textInput.Update(msg)
		return s, cmd
	}

	return s, nil
}

func (s SetupScreen) completeSetup() tea.Cmd {
	return func() tea.Msg {
		cfg := config.Config{SaveDirectory: s.saveDirectory}

		if s.communityEnabled {
			if s.useOfficialRepo {
				cfg.Repositories = []config.Repository{
					config.DefaultOfficialRepository(),
				}
			} else if s.customRepoURL != "" {
				// Extract name from URL
				name := extractRepoNameFromURL(s.customRepoURL)
				cfg.Repositories = []config.Repository{
					{
						Name:    name,
						URL:     s.customRepoURL,
						Enabled: true,
					},
				}
			}
		}

		return SetupCompleteMsg{
			Config:          cfg,
			EnableCommunity: s.communityEnabled,
		}
	}
}

func extractRepoNameFromURL(url string) string {
	// Extract name from URL like "https://github.com/user/repo-name"
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		name := parts[len(parts)-1]
		// Remove .git suffix if present
		if strings.HasSuffix(name, ".git") {
			name = name[:len(name)-4]
		}
		return name
	}
	return "custom"
}

func (s SetupScreen) View() string {
	title := styles.Title.Render("Welcome to PatternForge!")

	switch s.step {
	case StepSaveDirectory:
		subtitle := styles.Subtitle.Render("Where should results be saved?")
		help := styles.Help.Render("\nenter: confirm")

		return fmt.Sprintf("\n\n%s\n\n%s\n\n  %s\n%s",
			title,
			subtitle,
			s.textInput.View(),
			help,
		)

	case StepCommunityPatterns:
		subtitle := styles.Subtitle.Render("Enable community patterns?")
		description := styles.Subtitle.Render(
			"Community patterns are shared templates from git repositories.\n" +
				"You can add the official repository or your own custom repo.")
		help := styles.Help.Render("\nY: yes (recommended)  N: no")

		return fmt.Sprintf("\n\n%s\n\n%s\n\n%s\n%s",
			title,
			subtitle,
			description,
			help,
		)

	case StepRepoChoice:
		subtitle := styles.Subtitle.Render("Which repository to use?")
		description := styles.Subtitle.Render(
			"Official: https://github.com/Pierre-Malherbe/patternforge-patterns\n\n" +
				"You can always add more repositories later with:\n" +
				"  patternforge repo add <url>")
		help := styles.Help.Render("\nO: official (recommended)  C: custom URL")

		return fmt.Sprintf("\n\n%s\n\n%s\n\n%s\n%s",
			title,
			subtitle,
			description,
			help,
		)

	case StepCustomRepo:
		subtitle := styles.Subtitle.Render("Enter your repository URL:")
		help := styles.Help.Render("\nenter: confirm  esc: back")

		return fmt.Sprintf("\n\n%s\n\n%s\n\n  %s\n%s",
			title,
			subtitle,
			s.textInput.View(),
			help,
		)
	}

	return ""
}
