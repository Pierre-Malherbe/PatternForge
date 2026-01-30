package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Pierre-Malherbe/patternforge/internal/claude"
	"github.com/Pierre-Malherbe/patternforge/internal/config"
	"github.com/Pierre-Malherbe/patternforge/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

const version = "0.0.2"

func main() {
	// Check if Claude Code is available
	if err := claude.IsAvailable(); err != nil {
		fmt.Println(styles.Error.Render("❌ " + err.Error()))
		fmt.Println("\nInstall Claude Code:")
		fmt.Println("  pip install claude-code")
		fmt.Println("  claude auth login")
		os.Exit(1)
	}

	// Determine patterns directory
	patternsDir := "./patterns"
	if len(os.Args) > 1 {
		if os.Args[1] == "--version" || os.Args[1] == "-v" {
			fmt.Printf("PatternForge v%s\n", version)
			os.Exit(0)
		}
		if os.Args[1] == "--help" || os.Args[1] == "-h" {
			printHelp()
			os.Exit(0)
		}
		patternsDir = os.Args[1]
	}

	// Ensure patterns directory exists
	if err := ensurePatternsDir(patternsDir); err != nil {
		fmt.Println(styles.Error.Render(fmt.Sprintf("❌ %v", err)))
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load()
	firstLaunch := false
	if err != nil {
		if errors.Is(err, config.ErrNotFound) {
			firstLaunch = true
			cfg = config.Config{SaveDirectory: config.DefaultSaveDirectory()}
		} else {
			fmt.Println(styles.Error.Render(fmt.Sprintf("❌ %v", err)))
			os.Exit(1)
		}
	}

	// Initialize application
	model, err := NewModel(patternsDir, cfg, firstLaunch)
	if err != nil {
		fmt.Println(styles.Error.Render(fmt.Sprintf("❌ %v", err)))
		os.Exit(1)
	}

	// Run TUI
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// ensurePatternsDir creates the patterns directory and adds a default pattern if empty
func ensurePatternsDir(dir string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create patterns directory: %w", err)
	}

	// Check if any .md files exist
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read patterns directory: %w", err)
	}

	hasPatterns := false
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
			hasPatterns = true
			break
		}
	}

	// Create default pattern if none exist
	if !hasPatterns {
		defaultPattern := `# 📊 Mermaid Map

> Create detailed mind maps in Mermaid markdown

## Prompt

Create a detailed mind map in Mermaid markdown format from the following content:

{{input}}

Use:
- Clear hierarchical nodes
- Logical relationships between concepts
- Colors with style (fill:#FF6B6B, #4ECDC4, #45B7D1, etc.)
- Relevant emojis in nodes
- Syntax: graph TD for top-down

Structure:
1. One main central node
2. Main branches (3-5)
3. Sub-branches for details
4. CSS styles for each level

The result must be in a ` + "```mermaid ... ```" + ` block ready to display.
`
		path := filepath.Join(dir, "mermaid-map.md")
		if err := os.WriteFile(path, []byte(defaultPattern), 0644); err != nil {
			return fmt.Errorf("failed to create default pattern: %w", err)
		}

		fmt.Println(styles.Success.Render("✨ Created default pattern: mermaid-map.md"))
	}

	return nil
}

func printHelp() {
	help := `PatternForge - Dynamic prompt patterns powered by Claude Code

USAGE:
  patternforge [patterns-directory]

OPTIONS:
  -h, --help      Show this help
  -v, --version   Show version

PATTERNS DIRECTORY:
  Default: ./patterns
  Contains .md files with your prompt patterns

KEYBOARD SHORTCUTS:
  Navigation:     ↑/↓ or j/k
  Select:         enter
  Modify pattern: m (opens Vi)
  New pattern:    n (opens Vi)
  Process:        Ctrl+D (in input view)
  Copy output:    y (in results view)
  Save output:    s (in results view)
  Settings:       S (in selection view)
  Back:           Esc
  Quit:           q

EXAMPLES:
  # Use default patterns directory
  patternforge

  # Use custom directory
  patternforge ~/my-patterns

  # Check version
  patternforge --version

For more info: https://github.com/Pierre-Malherbe/patternforge
`
	fmt.Print(help)
}
