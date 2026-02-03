package screens

import (
	"fmt"
	"strings"

	"github.com/Pierre-Malherbe/patternforge/internal/pattern"
	"github.com/Pierre-Malherbe/patternforge/internal/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// VariablesScreen shows a form for configuring pattern variables
type VariablesScreen struct {
	pattern      pattern.Pattern
	fields       []variableField
	currentField int
	values       map[string]string
	width        int
	height       int
}

type variableField struct {
	variable     pattern.Variable
	textInput    textinput.Model
	selectIndex  int
	isSelect     bool
}

// NewVariablesScreen creates the variables configuration screen
func NewVariablesScreen(p pattern.Pattern) VariablesScreen {
	fields := make([]variableField, len(p.Variables))
	values := make(map[string]string)

	for i, v := range p.Variables {
		field := variableField{
			variable: v,
			isSelect: v.Type == "select" && len(v.Options) > 0,
		}

		if field.isSelect {
			// Find default option index
			field.selectIndex = 0
			for j, opt := range v.Options {
				if opt == v.Default {
					field.selectIndex = j
					break
				}
			}
			values[v.Name] = v.Options[field.selectIndex]
		} else {
			ti := textinput.New()
			ti.Placeholder = v.Placeholder
			if v.Default != "" {
				ti.SetValue(v.Default)
			}
			ti.CharLimit = 500
			ti.Width = 50
			if i == 0 {
				ti.Focus()
			}
			field.textInput = ti
			values[v.Name] = v.Default
		}

		fields[i] = field
	}

	return VariablesScreen{
		pattern:      p,
		fields:       fields,
		currentField: 0,
		values:       values,
	}
}

func (s VariablesScreen) Init() tea.Cmd {
	if len(s.fields) > 0 && !s.fields[0].isSelect {
		return textinput.Blink
	}
	return nil
}

func (s VariablesScreen) Update(msg tea.Msg) (VariablesScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		for i := range s.fields {
			if !s.fields[i].isSelect {
				s.fields[i].textInput.Width = msg.Width - 20
				if s.fields[i].textInput.Width > 60 {
					s.fields[i].textInput.Width = 60
				}
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			// Move to next field
			s.fields[s.currentField].textInput.Blur()
			s.currentField = (s.currentField + 1) % len(s.fields)
			if !s.fields[s.currentField].isSelect {
				s.fields[s.currentField].textInput.Focus()
				return s, textinput.Blink
			}
			return s, nil

		case "shift+tab", "up":
			// Move to previous field
			s.fields[s.currentField].textInput.Blur()
			s.currentField = (s.currentField - 1 + len(s.fields)) % len(s.fields)
			if !s.fields[s.currentField].isSelect {
				s.fields[s.currentField].textInput.Focus()
				return s, textinput.Blink
			}
			return s, nil

		case "left":
			// For select fields, move to previous option
			if s.fields[s.currentField].isSelect {
				field := &s.fields[s.currentField]
				field.selectIndex = (field.selectIndex - 1 + len(field.variable.Options)) % len(field.variable.Options)
				s.values[field.variable.Name] = field.variable.Options[field.selectIndex]
			}
			return s, nil

		case "right":
			// For select fields, move to next option
			if s.fields[s.currentField].isSelect {
				field := &s.fields[s.currentField]
				field.selectIndex = (field.selectIndex + 1) % len(field.variable.Options)
				s.values[field.variable.Name] = field.variable.Options[field.selectIndex]
			}
			return s, nil
		}
	}

	// Update current text input
	if len(s.fields) > 0 && !s.fields[s.currentField].isSelect {
		var cmd tea.Cmd
		s.fields[s.currentField].textInput, cmd = s.fields[s.currentField].textInput.Update(msg)
		s.values[s.fields[s.currentField].variable.Name] = s.fields[s.currentField].textInput.Value()
		return s, cmd
	}

	return s, nil
}

func (s VariablesScreen) View() string {
	title := styles.Title.Render(fmt.Sprintf("⚙️  Configure: %s", s.pattern.Name))
	desc := styles.Subtitle.Render(s.pattern.Description)

	var fieldsView strings.Builder
	fieldsView.WriteString("\n")

	for i, field := range s.fields {
		cursor := "  "
		if i == s.currentField {
			cursor = "▸ "
		}

		label := field.variable.Label
		if field.variable.Required {
			label += " *"
		}

		if field.isSelect {
			// Render select field
			var options strings.Builder
			for j, opt := range field.variable.Options {
				if j == field.selectIndex {
					options.WriteString(styles.Selected.Render(fmt.Sprintf("[%s]", opt)))
				} else {
					options.WriteString(fmt.Sprintf(" %s ", opt))
				}
			}
			fieldsView.WriteString(fmt.Sprintf("%s%s: %s\n\n", cursor, label, options.String()))
		} else {
			// Render text input
			fieldsView.WriteString(fmt.Sprintf("%s%s:\n   %s\n\n", cursor, label, field.textInput.View()))
		}
	}

	help := styles.Help.Render("\n↑/↓ or Tab: navigate fields • ←/→: change selection • enter: continue • esc: back")

	return fmt.Sprintf("%s\n%s\n%s\n%s",
		title,
		desc,
		fieldsView.String(),
		help,
	)
}

// IsComplete returns true if all required fields are filled
func (s VariablesScreen) IsComplete() bool {
	for _, field := range s.fields {
		if field.variable.Required {
			value := s.values[field.variable.Name]
			if strings.TrimSpace(value) == "" {
				return false
			}
		}
	}
	return true
}

// Values returns the configured variable values
func (s VariablesScreen) Values() map[string]string {
	return s.values
}

// Pattern returns the pattern being configured
func (s VariablesScreen) Pattern() pattern.Pattern {
	return s.pattern
}
