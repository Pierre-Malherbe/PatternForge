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
SelectionView → (enter) → InputView → (Ctrl+D) → ProcessingView → (auto) → ResultsView
     ↑              ↑                                                           |
     |              └── (Esc) ─────────────────────────────────────────────────┘
     └── (m) → Vi edit → ReloadPatternsMsg ─┐
     └── (n) → Vi create → ReloadPatternsMsg ┘
```

### Package Layout

- **`cmd/patternforge/`** — Application layer. `main.go` handles CLI startup (checks Claude availability, sets up patterns dir). `model.go` is the Bubbletea Model with view enum (`SelectionView`, `InputView`, `ProcessingView`, `ResultsView`), keyboard routing, and view transitions.
- **`internal/pattern/`** — Domain layer. Loads `.md` files from the patterns directory, parses markdown structure (emoji, name, description, prompt body), and handles `{{input}}` template substitution. If `{{input}}` is absent, user content is appended.
- **`internal/claude/`** — Service layer. Executes `claude ask --print` as a subprocess with the rendered prompt on stdin. Returns output plus stats (duration, estimated tokens at ~4 chars/token).
- **`internal/ui/screens/`** — Each screen (`selection.go`, `input.go`, `processing.go`, `results.go`) is a self-contained Bubbletea component with `Init()`, `Update()`, `View()` methods.
- **`internal/ui/styles/`** — Centralized Lipgloss style definitions.
- **`patterns/`** — User-facing markdown pattern files loaded at startup.

### Pattern File Format

```markdown
# [EMOJI] [TITLE]
> [One-line description]

## Prompt

[Instructions with optional {{input}} placeholder]
```

## Key Bubbletea Patterns Used

- **`tea.ExecProcess`** — Used to shell out to Vi for pattern editing, then sends `ReloadPatternsMsg` on return.
- **`tea.Batch`** — Used in `startProcessing()` to run the spinner animation and Claude execution concurrently.
- **`WindowSizeMsg` forwarding** — Must be passed to `ResultsScreen` on creation (not just on resize) to avoid empty viewport on first load.

## Known Pitfalls

1. **Viewport empty on first load**: `ResultsScreen` must receive a `WindowSizeMsg` immediately after creation in the `ProcessedMsg` handler, not just from window resize events.
2. **Pattern list stale after Vi edit**: The `ReloadPatternsMsg` handler must both reload patterns and re-send `WindowSizeMsg` to the new `SelectionScreen`.
3. **New pattern filename collisions**: New patterns use `new-pattern-{unix_timestamp}.md` to avoid overwrites.