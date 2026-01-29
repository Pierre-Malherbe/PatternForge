package pattern

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Pattern represents a reusable prompt template
type Pattern struct {
	Name        string // Display name
	Emoji       string // Visual identifier
	Description string // Short one-liner
	Prompt      string // Template with optional {{input}}
	FilePath    string // Where it's stored
}

// LoadAll reads all .md files from a directory
func LoadAll(dir string) ([]Pattern, error) {
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
		patterns = append(patterns, p)
	}

	return patterns, nil
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
//   # [emoji] [name]
//   > [description]
//   ## Prompt
//   [prompt content with optional {{input}}]
func parse(content, path string) Pattern {
	p := Pattern{
		FilePath: path,
		Emoji:    "📄", // Default emoji
	}

	lines := strings.Split(content, "\n")
	inPrompt := false
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
		} else if trimmed == "## Prompt" {
			// Start of prompt section
			inPrompt = true
		} else if inPrompt {
			if strings.HasPrefix(trimmed, "##") {
				// New section, stop collecting prompt
				inPrompt = false
			} else if trimmed != "" && !strings.HasPrefix(trimmed, "```") {
				// Collect prompt lines (skip code fences)
				promptLines = append(promptLines, line)
			}
		}
	}

	p.Prompt = strings.TrimSpace(strings.Join(promptLines, "\n"))
	return p
}

// Render replaces {{input}} with actual user content
func (p Pattern) Render(userInput string) string {
	if strings.Contains(p.Prompt, "{{input}}") {
		return strings.ReplaceAll(p.Prompt, "{{input}}", userInput)
	}
	// If no {{input}}, append user content at the end
	return p.Prompt + "\n\n---\n\n" + userInput
}

// DisplayTitle returns emoji + name for UI
func (p Pattern) DisplayTitle() string {
	return p.Emoji + " " + p.Name
}
