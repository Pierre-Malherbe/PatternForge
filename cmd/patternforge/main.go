package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Pierre-Malherbe/patternforge/internal/claude"
	"github.com/Pierre-Malherbe/patternforge/internal/config"
	"github.com/Pierre-Malherbe/patternforge/internal/pattern"
	"github.com/Pierre-Malherbe/patternforge/internal/repository"
	"github.com/Pierre-Malherbe/patternforge/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

const version = "0.0.3"

func main() {
	// Handle subcommands first
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Printf("PatternForge v%s\n", version)
			os.Exit(0)
		case "--help", "-h":
			printHelp()
			os.Exit(0)
		case "upgrade":
			runUpgrade()
			os.Exit(0)
		case "repo":
			runRepoCommand()
			os.Exit(0)
		}
	}

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

// runUpgrade syncs all configured repositories.
func runUpgrade() {
	cfg, err := config.Load()
	if err != nil {
		if errors.Is(err, config.ErrNotFound) {
			fmt.Println(styles.Error.Render("❌ No configuration found. Run patternforge first to set up."))
			return
		}
		fmt.Println(styles.Error.Render(fmt.Sprintf("❌ %v", err)))
		return
	}

	if len(cfg.Repositories) == 0 {
		fmt.Println("No repositories configured.")
		fmt.Println("Run 'patternforge repo add <url>' to add a repository.")
		return
	}

	if !repository.IsGitAvailable() {
		fmt.Println(styles.Error.Render("❌ Git is not installed. Cannot upgrade repositories."))
		return
	}

	// Convert config.Repository to repository.Repository
	repos := make([]repository.Repository, len(cfg.Repositories))
	for i, r := range cfg.Repositories {
		repos[i] = repository.Repository{
			Name:       r.Name,
			URL:        r.URL,
			Enabled:    r.Enabled,
			LastUpdate: r.LastUpdate,
		}
	}

	mgr := repository.NewManager(config.RepositoriesDir(), repos)
	results := mgr.UpdateAll()

	fmt.Println(styles.Title.Render("Upgrading repositories..."))
	fmt.Println()

	hasUpdates := false
	for _, r := range results {
		if r.Success {
			fmt.Printf("  %s %s: %s\n", styles.Success.Render("✓"), r.Name, r.Message)
			hasUpdates = true
		} else {
			fmt.Printf("  %s %s: %s\n", styles.Error.Render("✗"), r.Name, r.Message)
		}
	}

	if hasUpdates {
		// Update last_update timestamps
		for i := range cfg.Repositories {
			for _, r := range results {
				if cfg.Repositories[i].Name == r.Name && r.Success {
					cfg.Repositories[i].LastUpdate = repository.FormatLastUpdate()
				}
			}
		}
		config.Save(cfg)
	}

	fmt.Println()
	fmt.Println("Done!")
}

// runRepoCommand handles repo subcommands.
func runRepoCommand() {
	if len(os.Args) < 3 {
		printRepoHelp()
		return
	}

	switch os.Args[2] {
	case "list":
		runRepoList()
	case "add":
		if len(os.Args) < 4 {
			fmt.Println(styles.Error.Render("❌ Usage: patternforge repo add <url> [name]"))
			return
		}
		name := ""
		if len(os.Args) >= 5 {
			name = os.Args[4]
		}
		runRepoAdd(os.Args[3], name)
	case "remove":
		if len(os.Args) < 4 {
			fmt.Println(styles.Error.Render("❌ Usage: patternforge repo remove <name>"))
			return
		}
		runRepoRemove(os.Args[3])
	default:
		printRepoHelp()
	}
}

func runRepoList() {
	cfg, err := config.Load()
	if err != nil {
		if errors.Is(err, config.ErrNotFound) {
			fmt.Println("No repositories configured.")
			return
		}
		fmt.Println(styles.Error.Render(fmt.Sprintf("❌ %v", err)))
		return
	}

	if len(cfg.Repositories) == 0 {
		fmt.Println("No repositories configured.")
		fmt.Println("Run 'patternforge repo add <url>' to add a repository.")
		return
	}

	fmt.Println(styles.Title.Render("Configured repositories:"))
	fmt.Println()

	for i, repo := range cfg.Repositories {
		status := "enabled"
		if !repo.Enabled {
			status = "disabled"
		}
		fmt.Printf("  %d. %s (%s)\n", i+1, repo.Name, status)
		fmt.Printf("     URL: %s\n", repo.URL)
		if repo.LastUpdate != "" {
			fmt.Printf("     Last update: %s\n", repo.LastUpdate)
		}
		fmt.Println()
	}
}

func runRepoAdd(url, name string) {
	cfg, err := config.Load()
	if err != nil && !errors.Is(err, config.ErrNotFound) {
		fmt.Println(styles.Error.Render(fmt.Sprintf("❌ %v", err)))
		return
	}
	if errors.Is(err, config.ErrNotFound) {
		cfg = config.Config{SaveDirectory: config.DefaultSaveDirectory()}
	}

	// Generate name from URL if not provided
	if name == "" {
		name = extractRepoName(url)
	}

	// Check for duplicates
	for _, repo := range cfg.Repositories {
		if repo.Name == name {
			fmt.Println(styles.Error.Render(fmt.Sprintf("❌ Repository %q already exists", name)))
			return
		}
		if repo.URL == url {
			fmt.Println(styles.Error.Render(fmt.Sprintf("❌ Repository with URL %q already exists as %q", url, repo.Name)))
			return
		}
	}

	cfg.Repositories = append(cfg.Repositories, config.Repository{
		Name:    name,
		URL:     url,
		Enabled: true,
	})

	if err := config.Save(cfg); err != nil {
		fmt.Println(styles.Error.Render(fmt.Sprintf("❌ %v", err)))
		return
	}

	fmt.Printf("%s Added repository %q\n", styles.Success.Render("✓"), name)
	fmt.Println("Run 'patternforge upgrade' to clone the repository.")
}

func runRepoRemove(nameOrIndex string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(styles.Error.Render(fmt.Sprintf("❌ %v", err)))
		return
	}

	// Try to parse as index first
	if idx, err := strconv.Atoi(nameOrIndex); err == nil {
		if idx < 1 || idx > len(cfg.Repositories) {
			fmt.Println(styles.Error.Render(fmt.Sprintf("❌ Invalid index: %d", idx)))
			return
		}
		name := cfg.Repositories[idx-1].Name
		cfg.Repositories = append(cfg.Repositories[:idx-1], cfg.Repositories[idx:]...)

		// Also remove the cloned repo if it exists
		repos := make([]repository.Repository, len(cfg.Repositories))
		for i, r := range cfg.Repositories {
			repos[i] = repository.Repository{Name: r.Name, URL: r.URL, Enabled: r.Enabled}
		}
		mgr := repository.NewManager(config.RepositoriesDir(), repos)
		mgr.Remove(name)

		if err := config.Save(cfg); err != nil {
			fmt.Println(styles.Error.Render(fmt.Sprintf("❌ %v", err)))
			return
		}
		fmt.Printf("%s Removed repository %q\n", styles.Success.Render("✓"), name)
		return
	}

	// Find by name
	found := -1
	for i, repo := range cfg.Repositories {
		if repo.Name == nameOrIndex {
			found = i
			break
		}
	}

	if found == -1 {
		fmt.Println(styles.Error.Render(fmt.Sprintf("❌ Repository %q not found", nameOrIndex)))
		return
	}

	name := cfg.Repositories[found].Name
	cfg.Repositories = append(cfg.Repositories[:found], cfg.Repositories[found+1:]...)

	// Also remove the cloned repo if it exists
	repos := make([]repository.Repository, len(cfg.Repositories))
	for i, r := range cfg.Repositories {
		repos[i] = repository.Repository{Name: r.Name, URL: r.URL, Enabled: r.Enabled}
	}
	mgr := repository.NewManager(config.RepositoriesDir(), repos)
	mgr.Remove(name)

	if err := config.Save(cfg); err != nil {
		fmt.Println(styles.Error.Render(fmt.Sprintf("❌ %v", err)))
		return
	}
	fmt.Printf("%s Removed repository %q\n", styles.Success.Render("✓"), name)
}

func extractRepoName(url string) string {
	// Extract name from URL like "https://github.com/user/repo-name"
	parts := filepath.Base(url)
	// Remove .git suffix if present
	if len(parts) > 4 && parts[len(parts)-4:] == ".git" {
		parts = parts[:len(parts)-4]
	}
	return parts
}

func printRepoHelp() {
	help := `Repository management commands:

  patternforge repo list              List configured repositories
  patternforge repo add <url> [name]  Add a new repository
  patternforge repo remove <name|#>   Remove a repository by name or number
`
	fmt.Print(help)
}

// LoadPatternsWithRepos loads patterns from local dir and all enabled repositories.
// Local patterns take precedence over repository patterns.
func LoadPatternsWithRepos(localDir string, cfg config.Config) ([]pattern.Pattern, error) {
	// Load local patterns first (highest priority)
	localPatterns, err := pattern.LoadAll(localDir)
	if err != nil {
		return nil, err
	}

	// If no repos configured or git not available, return local only
	if len(cfg.Repositories) == 0 || !repository.IsGitAvailable() {
		return localPatterns, nil
	}

	// Collect all pattern sources (local first for priority)
	sources := [][]pattern.Pattern{localPatterns}

	// Load patterns from each enabled repository
	repos := make([]repository.Repository, len(cfg.Repositories))
	for i, r := range cfg.Repositories {
		repos[i] = repository.Repository{Name: r.Name, URL: r.URL, Enabled: r.Enabled}
	}
	mgr := repository.NewManager(config.RepositoriesDir(), repos)

	for _, repo := range cfg.Repositories {
		if !repo.Enabled {
			continue
		}

		patternsDir := mgr.GetPatternsDir(repo.Name)
		if _, err := os.Stat(patternsDir); os.IsNotExist(err) {
			// Repo not cloned yet, skip
			continue
		}

		repoPatterns, err := pattern.LoadAllWithSource(patternsDir, repo.Name)
		if err != nil {
			// Skip repos that fail to load
			continue
		}
		sources = append(sources, repoPatterns)
	}

	// Merge with local taking precedence
	return pattern.MergePatterns(sources...), nil
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
  patternforge upgrade
  patternforge repo <command>

COMMANDS:
  upgrade              Sync all configured repositories
  repo list            List configured repositories
  repo add <url>       Add a new pattern repository
  repo remove <name>   Remove a repository

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
  Upgrade repos:  U (syncs community patterns)
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

  # Add official patterns repository
  patternforge repo add https://github.com/Pierre-Malherbe/patternforge-patterns

  # Sync repositories
  patternforge upgrade

For more info: https://github.com/Pierre-Malherbe/patternforge
`
	fmt.Print(help)
}
