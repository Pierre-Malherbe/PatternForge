package repository

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Repository represents a pattern repository configuration.
type Repository struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Enabled    bool   `json:"enabled"`
	LastUpdate string `json:"last_update,omitempty"`
}

// DefaultOfficialRepository returns the default official repository config.
func DefaultOfficialRepository() Repository {
	return Repository{
		Name:    "official",
		URL:     "https://github.com/Pierre-Malherbe/patternforge-patterns",
		Enabled: true,
	}
}

// Manager handles git operations for pattern repositories.
type Manager struct {
	reposDir string
	repos    []Repository
}

// NewManager creates a new repository manager.
func NewManager(reposDir string, repos []Repository) *Manager {
	return &Manager{
		reposDir: reposDir,
		repos:    repos,
	}
}

// IsGitAvailable checks if git is installed and accessible.
func IsGitAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// Clone clones a repository to the repos directory.
// Uses shallow clone (--depth 1) to minimize data transfer.
func (m *Manager) Clone(url, name string) error {
	if !IsGitAvailable() {
		return fmt.Errorf("git is not installed")
	}

	repoPath := filepath.Join(m.reposDir, name)

	// Check if already exists
	if _, err := os.Stat(repoPath); err == nil {
		return fmt.Errorf("repository %q already exists", name)
	}

	// Create parent directory if needed
	if err := os.MkdirAll(m.reposDir, 0755); err != nil {
		return fmt.Errorf("failed to create repositories directory: %w", err)
	}

	// Clone with depth 1 for shallow clone
	cmd := exec.Command("git", "clone", "--depth", "1", url, repoPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %s", strings.TrimSpace(string(output)))
	}

	return nil
}

// Pull updates a repository to the latest version.
func (m *Manager) Pull(name string) error {
	if !IsGitAvailable() {
		return fmt.Errorf("git is not installed")
	}

	repoPath := filepath.Join(m.reposDir, name)

	// Check if exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return fmt.Errorf("repository %q does not exist", name)
	}

	cmd := exec.Command("git", "-C", repoPath, "pull", "--ff-only")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git pull failed: %s", strings.TrimSpace(string(output)))
	}

	return nil
}

// GetPatternsDir returns the patterns directory for a repository.
func (m *Manager) GetPatternsDir(name string) string {
	return filepath.Join(m.reposDir, name, "patterns")
}

// GetRepoDir returns the base directory for a repository.
func (m *Manager) GetRepoDir(name string) string {
	return filepath.Join(m.reposDir, name)
}

// Exists checks if a repository has been cloned.
func (m *Manager) Exists(name string) bool {
	repoPath := filepath.Join(m.reposDir, name)
	_, err := os.Stat(repoPath)
	return err == nil
}

// Remove deletes a cloned repository.
func (m *Manager) Remove(name string) error {
	repoPath := filepath.Join(m.reposDir, name)

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return fmt.Errorf("repository %q does not exist", name)
	}

	return os.RemoveAll(repoPath)
}

// UpdateResult represents the result of updating a single repository.
type UpdateResult struct {
	Name    string
	Success bool
	Message string
}

// UpdateAll pulls the latest changes for all enabled repositories.
// Returns results for each repository.
func (m *Manager) UpdateAll() []UpdateResult {
	var results []UpdateResult

	for _, repo := range m.repos {
		if !repo.Enabled {
			continue
		}

		result := UpdateResult{Name: repo.Name}

		if !m.Exists(repo.Name) {
			// Clone if not exists
			if err := m.Clone(repo.URL, repo.Name); err != nil {
				result.Success = false
				result.Message = err.Error()
			} else {
				result.Success = true
				result.Message = "Cloned successfully"
			}
		} else {
			// Pull if exists
			if err := m.Pull(repo.Name); err != nil {
				result.Success = false
				result.Message = err.Error()
			} else {
				result.Success = true
				result.Message = "Updated successfully"
			}
		}

		results = append(results, result)
	}

	return results
}

// GetEnabledRepos returns a list of enabled repositories.
func (m *Manager) GetEnabledRepos() []Repository {
	var enabled []Repository
	for _, repo := range m.repos {
		if repo.Enabled {
			enabled = append(enabled, repo)
		}
	}
	return enabled
}

// GetRepos returns all configured repositories.
func (m *Manager) GetRepos() []Repository {
	return m.repos
}

// FormatLastUpdate returns a formatted timestamp string.
func FormatLastUpdate() string {
	return time.Now().Format(time.RFC3339)
}
