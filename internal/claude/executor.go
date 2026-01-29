package claude

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"
)

// Stats represents execution metrics
type Stats struct {
	InputTokens  int
	OutputTokens int
	TotalTokens  int
	Duration     time.Duration
}

// IsAvailable checks if Claude Code CLI is installed and working
func IsAvailable() error {
	path, err := exec.LookPath("claude")
	if err != nil {
		return fmt.Errorf("claude CLI not found in PATH. Install: pip install claude-code")
	}

	// Try a simple help command to verify it works
	cmd := exec.Command(path, "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("claude CLI found at %s but not functional: %w", path, err)
	}

	return nil
}

// Execute runs the prompt through Claude Code and returns the response
func Execute(prompt string) (string, Stats, error) {
	start := time.Now()

	// Execute: claude ask with --print for non-interactive output
	// -p/--print: Print response and exit (batch mode, no TUI)
	cmd := exec.Command("claude", "ask", "--print")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = bytes.NewBufferString(prompt)

	if err := cmd.Run(); err != nil {
		return "", Stats{}, fmt.Errorf("claude execution failed: %w\nStderr: %s", err, stderr.String())
	}

	duration := time.Since(start)
	result := stdout.String()

	// Calculate stats (rough estimation since Claude Code doesn't return token counts)
	stats := Stats{
		InputTokens:  estimateTokens(prompt),
		OutputTokens: estimateTokens(result),
		Duration:     duration,
	}
	stats.TotalTokens = stats.InputTokens + stats.OutputTokens

	return result, stats, nil
}

// estimateTokens approximates token count (1 token ≈ 4 characters)
func estimateTokens(text string) int {
	return len(text) / 4
}
