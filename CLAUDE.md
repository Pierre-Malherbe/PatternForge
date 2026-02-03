# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

PatternForge is a Go TUI application that provides Fabric-style reusable prompt patterns for Claude Code. Users select a markdown pattern template, input content, and the app sends the composed prompt to `claude ask --print` and displays results.

## Build & Development Commands

```bash
make build          # Build binary to bin/patternforge
make run            # Run without building (go run)
make test           # Run tests (go test -v ./...)
make deps           # Install/tidy Go dependencies
make build-macos    # Cross-compile for ARM64 + AMD64
make install        # Build and install to /usr/local/bin
make clean          # Remove bin/ directory
```

## Architecture

The app uses the **Elm Architecture** via Bubbletea. The central state machine is in `cmd/patternforge/model.go`.

### State Machine Flow

```
CategoryView → (enter) → SelectionView → (enter) → [VariablesView →] InputView → (Ctrl+D) → ProcessingView → ResultsView
     ↑              ↑            ↑                                                                              |
     |              |            └── (Esc) ────────────────────────────────────────────────────────────────────┘
     |              └── (Esc) ─────────────────────────────────────────────────────────────────────────────────┘
     └── (m) → Vi edit → ReloadPatternsMsg
     └── (n) → Vi create → ReloadPatternsMsg
     └── (U) → UpgradeView → (Esc) ──────────────────────────────────────────────────────────────────────────┘
```

### Package Layout

- **`cmd/patternforge/`** — Application layer. `main.go` handles CLI startup, subcommands (`upgrade`, `repo`), and pattern loading. `model.go` is the Bubbletea Model with view enum, keyboard routing, and view transitions.
- **`internal/pattern/`** — Domain layer. Loads `.md` files, parses markdown structure (emoji, name, description, category, variables, prompt), handles `{{input}}` and `{{var:name}}` substitution.
- **`internal/repository/`** — Git operations for cloning/pulling pattern repositories.
- **`internal/config/`** — Configuration management including repositories.
- **`internal/claude/`** — Service layer. Executes `claude ask --print` as a subprocess.
- **`internal/ui/screens/`** — Each screen is a self-contained Bubbletea component.
- **`internal/ui/styles/`** — Centralized Lipgloss style definitions.
- **`patterns/`** — User-facing markdown pattern files loaded at startup.

### Pattern File Format

#### Basic Pattern
```markdown
# [EMOJI] [TITLE]
> [One-line description]

[Category: CategoryName]

## Prompt

[Instructions with optional {{input}} placeholder]
```

#### Pattern with Variables
```markdown
# [EMOJI] [TITLE]
> [One-line description]

[Category: CategoryName]

## Variables

- name: Label | type | options | default | placeholder
- priority: Priority | select | low,medium,high | medium
- env: Environment | select | dev,staging,prod | dev

## Prompt

Your prompt using {{var:priority}} and {{var:env}} placeholders.

{{input}}
```

Variable types:
- `text` - Simple text input (default)
- `select` - Dropdown with options (comma-separated)
- `multiline` - Multi-line text area

### Repository System

Patterns can be loaded from:
1. Local directory (highest priority)
2. Community repositories (cloned via git)

Commands:
- `patternforge upgrade` - Sync all repositories
- `patternforge repo list` - List configured repositories
- `patternforge repo add <url>` - Add a repository
- `patternforge repo remove <name>` - Remove a repository

## Key Bubbletea Patterns Used

- **`tea.ExecProcess`** — Used to shell out to Vi for pattern editing, then sends `ReloadPatternsMsg` on return.
- **`tea.Batch`** — Used in `startProcessing()` to run the spinner animation and Claude execution concurrently.
- **`WindowSizeMsg` forwarding** — Must be passed to screens on creation (not just on resize) to avoid empty viewport on first load.

## Known Pitfalls

1. **Viewport empty on first load**: Screens must receive a `WindowSizeMsg` immediately after creation, not just from window resize events.
2. **Pattern list stale after Vi edit**: The `ReloadPatternsMsg` handler must both reload patterns and re-send `WindowSizeMsg` to the new `SelectionScreen`.
3. **New pattern filename collisions**: New patterns use `new-pattern-{unix_timestamp}.md` to avoid overwrites.
4. **Category view skip**: If only one category exists, the app skips directly to pattern list.
5. **Variables view conditional**: Only shown for patterns with `## Variables` section.
