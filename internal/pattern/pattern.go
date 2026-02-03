package pattern

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Pattern represents a reusable prompt template
type Pattern struct {
	Name        string     // Display name
	Emoji       string     // Visual identifier
	Description string     // Short one-liner
	Category    string     // Category for grouping (e.g., "Bug", "Tickets", "Refacto")
	Prompt      string     // Template with optional {{input}}
	FilePath    string     // Where it's stored
	Source      string     // Origin: "local", "official", or custom repo name
	Variables   []Variable // Optional variables for form-based input
}

// Variable represents a configurable field in a pattern
type Variable struct {
	Name        string   // Variable name (used in template as {{var:name}})
	Label       string   // Display label for the form
	Type        string   // Type: "text", "select", "multiline"
	Options     []string // Options for select type
	Default     string   // Default value
	Placeholder string   // Placeholder text
	Required    bool     // Whether the variable is required
}

// DefaultCategory is used when no category is specified
const DefaultCategory = "General"

// LoadAll reads all .md files from a directory (Source defaults to "local")
func LoadAll(dir string) ([]Pattern, error) {
	return LoadAllWithSource(dir, "local")
}

// LoadAllWithSource reads all .md files from a directory with a specific source tag
func LoadAllWithSource(dir, source string) ([]Pattern, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read patterns directory: %w", err)
	}

	var patterns []Pattern
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		p, err := Load(path)
		if err != nil {
			// Skip invalid patterns, don't fail the whole load
			continue
		}
		p.Source = source
		patterns = append(patterns, p)
	}

	return patterns, nil
}

// MergePatterns merges multiple pattern sources with earlier sources taking precedence.
// Patterns are identified by their Name; duplicates from later sources are ignored.
func MergePatterns(sources ...[]Pattern) []Pattern {
	seen := make(map[string]bool)
	var result []Pattern

	for _, patterns := range sources {
		for _, p := range patterns {
			if !seen[p.Name] {
				seen[p.Name] = true
				result = append(result, p)
			}
		}
	}

	return result
}

// Load reads and parses a single pattern file
func Load(path string) (Pattern, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Pattern{}, fmt.Errorf("failed to read pattern: %w", err)
	}

	return parse(string(content), path), nil
}

// parse extracts pattern metadata from markdown
// Format:
//
//	# [emoji] [name]
//	> [description]
//	[Category: CategoryName]
//	## Variables (optional)
//	- name: Label | type | options | default | placeholder
//	## Prompt
//	[prompt content with optional {{input}} and {{var:name}}]
func parse(content, path string) Pattern {
	p := Pattern{
		FilePath: path,
		Emoji:    "📄",           // Default emoji
		Category: DefaultCategory, // Default category
	}

	lines := strings.Split(content, "\n")
	inPrompt := false
	inVariables := false
	var promptLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "# ") {
			// Title line: "# 📊 Mermaid Map"
			title := strings.TrimPrefix(trimmed, "# ")
			parts := strings.SplitN(title, " ", 2)
			if len(parts) == 2 {
				p.Emoji = parts[0]
				p.Name = parts[1]
			} else {
				p.Name = title
			}
		} else if strings.HasPrefix(trimmed, "> ") {
			// Description line
			p.Description = strings.TrimPrefix(trimmed, "> ")
		} else if strings.HasPrefix(trimmed, "[Category:") && strings.HasSuffix(trimmed, "]") {
			// Category line: [Category: Bug]
			cat := strings.TrimPrefix(trimmed, "[Category:")
			cat = strings.TrimSuffix(cat, "]")
			p.Category = strings.TrimSpace(cat)
		} else if trimmed == "## Variables" {
			// Start of variables section
			inVariables = true
			inPrompt = false
		} else if trimmed == "## Prompt" {
			// Start of prompt section
			inPrompt = true
			inVariables = false
		} else if inVariables {
			if strings.HasPrefix(trimmed, "##") {
				inVariables = false
			} else if strings.HasPrefix(trimmed, "- ") {
				// Parse variable: - name: Label | type | options | default | placeholder
				v := parseVariable(strings.TrimPrefix(trimmed, "- "))
				if v.Name != "" {
					p.Variables = append(p.Variables, v)
				}
			}
		} else if inPrompt {
			if strings.HasPrefix(trimmed, "##") {
				// New section, stop collecting prompt
				inPrompt = false
			} else {
				// Collect prompt lines (include empty lines for formatting)
				promptLines = append(promptLines, line)
			}
		}
	}

	p.Prompt = strings.TrimSpace(strings.Join(promptLines, "\n"))
	return p
}

// parseVariable parses a variable definition line
// Format: name: Label | type | options | default | placeholder
// Examples:
//   - priority: Priority | select | low,medium,high,critical | medium
//   - description: Description | multiline | | | Enter details...
//   - env: Environment | select | dev,staging,prod | dev
func parseVariable(line string) Variable {
	v := Variable{
		Type:     "text",
		Required: false,
	}

	// Split by colon for name
	colonIdx := strings.Index(line, ":")
	if colonIdx == -1 {
		return v
	}

	v.Name = strings.TrimSpace(line[:colonIdx])
	rest := strings.TrimSpace(line[colonIdx+1:])

	// Split by pipe for other fields
	parts := strings.Split(rest, "|")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}

	if len(parts) >= 1 && parts[0] != "" {
		v.Label = parts[0]
	} else {
		v.Label = v.Name // Default label to name
	}

	if len(parts) >= 2 && parts[1] != "" {
		v.Type = parts[1]
	}

	if len(parts) >= 3 && parts[2] != "" {
		v.Options = strings.Split(parts[2], ",")
		for i := range v.Options {
			v.Options[i] = strings.TrimSpace(v.Options[i])
		}
	}

	if len(parts) >= 4 {
		v.Default = parts[3]
	}

	if len(parts) >= 5 {
		v.Placeholder = parts[4]
	}

	// Check for required marker (asterisk in name)
	if strings.HasSuffix(v.Name, "*") {
		v.Name = strings.TrimSuffix(v.Name, "*")
		v.Required = true
	}

	return v
}

// HasVariables returns true if the pattern has configurable variables
func (p Pattern) HasVariables() bool {
	return len(p.Variables) > 0
}

// Render replaces {{input}} with actual user content
func (p Pattern) Render(userInput string) string {
	if strings.Contains(p.Prompt, "{{input}}") {
		return strings.ReplaceAll(p.Prompt, "{{input}}", userInput)
	}
	// If no {{input}}, append user content at the end
	return p.Prompt + "\n\n---\n\n" + userInput
}

// RenderWithVariables replaces {{input}} and {{var:name}} with provided values
func (p Pattern) RenderWithVariables(userInput string, variables map[string]string) string {
	result := p.Prompt

	// Replace variable placeholders
	for name, value := range variables {
		placeholder := "{{var:" + name + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// Replace {{input}} placeholder
	if strings.Contains(result, "{{input}}") {
		result = strings.ReplaceAll(result, "{{input}}", userInput)
	} else if userInput != "" {
		// If no {{input}}, append user content at the end
		result = result + "\n\n---\n\n" + userInput
	}

	return result
}

// GetCategories returns a sorted unique list of categories from patterns
func GetCategories(patterns []Pattern) []string {
	categoryMap := make(map[string]bool)
	for _, p := range patterns {
		categoryMap[p.Category] = true
	}

	var categories []string
	// Put "General" first if it exists
	if categoryMap[DefaultCategory] {
		categories = append(categories, DefaultCategory)
		delete(categoryMap, DefaultCategory)
	}

	// Add remaining categories sorted
	for cat := range categoryMap {
		categories = append(categories, cat)
	}

	return categories
}

// FilterByCategory returns patterns that match the given category
func FilterByCategory(patterns []Pattern, category string) []Pattern {
	if category == "" {
		return patterns
	}

	var filtered []Pattern
	for _, p := range patterns {
		if p.Category == category {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

// DisplayTitle returns emoji + name for UI
func (p Pattern) DisplayTitle() string {
	return p.Emoji + " " + p.Name
}
