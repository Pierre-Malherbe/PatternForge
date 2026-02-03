# 🎨 PatternForge

**Dynamic prompt patterns powered by Claude Code**

PatternForge is a blazing-fast TUI (Terminal User Interface) that brings the power of [Fabric](https://github.com/danielmiessler/fabric)-style patterns to Claude Code. Create, edit, and manage reusable AI prompts with an elegant interface and Vim-style keybindings.

![Version](https://img.shields.io/badge/version-0.0.3-blue)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

## ✨ Features

- 🎯 **Dynamic Pattern Loading** - Patterns are markdown files you can edit anytime
- 📁 **Category Navigation** - Organize patterns by category (Bug, Tickets, Refactor, etc.)
- ⚙️ **Pattern Variables** - Forms with configurable fields (priority, environment, etc.)
- 🌐 **Community Patterns** - Clone and sync patterns from git repositories
- 📝 **{{input}} Templating** - Control exactly where user content goes
- ⚡ **Vim Keybindings** - j/k navigation, natural for terminal users
- 🔧 **Live Editing** - Press `m` to edit patterns in Vi, changes apply instantly
- 📊 **Token Stats** - See input/output tokens and execution time
- 📋 **Copy to Clipboard** - One key press to copy results
- 🎨 **Beautiful TUI** - Built with Bubbletea and Lipgloss
- 🚀 **Fast & Lightweight** - Pure Go, single binary, no dependencies

## 🎬 Quick Demo

```bash
# 1. Launch PatternForge
$ patternforge

# 2. Select a category, then a pattern
> 📁 Tickets
  > 🎫 Create Ticket ⚙️

# 3. Fill in variables (if pattern has them)
Priority: [low] medium [high]
Type: [bug] feature [task]

# 4. Paste your content, press Ctrl+D

# 5. Get beautiful results with stats
✨ Result
[Your formatted ticket here]
📊 172 tokens input + 612 output = 784 total | 3.2s

# 6. Press 'y' to copy output
```

## 📦 Installation

### Prerequisites

1. **Go 1.21+** - [Install Go](https://go.dev/doc/install)
2. **Claude Code** - Required for AI processing
   ```bash
   pip install claude-code
   claude auth login
   ```
3. **Git** - For community patterns (optional)
4. **Vi/Vim** - For editing patterns (usually pre-installed on macOS/Linux)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/Pierre-Malherbe/PatternForge.git
cd patternforge

# Install dependencies
make deps

# Build
make build

# Install system-wide (optional)
make install
```

## 🚀 Usage

### Basic Usage

```bash
# Use default patterns directory (./patterns)
patternforge

# Use custom patterns directory
patternforge ~/my-patterns
```

### CLI Commands

```bash
# Sync community patterns from repositories
patternforge upgrade

# List configured repositories
patternforge repo list

# Add a pattern repository
patternforge repo add https://github.com/user/patterns-repo

# Remove a repository
patternforge repo remove <name>
```

### First Run

On first run, PatternForge will:
1. ✅ Check if Claude Code is installed
2. 📁 Ask where to save results
3. 🌐 Offer to enable community patterns (clones official repo)
4. 🎯 Launch the TUI

### Workflow

1. **Select Category** - Browse patterns by category
2. **Select Pattern** - Use `j/k` or arrow keys, press `enter`
3. **Configure Variables** - Fill in form fields (if pattern has variables)
4. **Input Content** - Paste or type your content
5. **Process** - Press `Ctrl+D` to send to Claude Code
6. **View Results** - Scroll through output, see token stats
7. **Copy/Save** - Press `y` to copy, or `s` to save to file

## ⌨️ Keyboard Shortcuts

### Category/Pattern Selection Screen

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `enter` | Select category/pattern |
| `esc` | Back to categories |
| `m` | Modify selected pattern (opens Vi) |
| `n` | Create new pattern (opens Vi) |
| `U` | Upgrade/sync repositories |
| `S` | Settings |
| `/` | Search |
| `q` | Quit |

### Variables Screen

| Key | Action |
|-----|--------|
| `↑` / `↓` / `Tab` | Navigate fields |
| `←` / `→` | Change selection (for select fields) |
| `enter` | Continue to input |
| `esc` | Back |

### Input Screen

| Key | Action |
|-----|--------|
| `Ctrl+D` | Process with Claude Code |
| `Esc` | Back to pattern selection |

### Results Screen

| Key | Action |
|-----|--------|
| `j/k` or `↑/↓` | Scroll output |
| `y` | Copy output to clipboard |
| `s` | Save output to file |
| `Esc` | New pattern |
| `q` | Quit |

## 📝 Creating Patterns

Patterns are simple Markdown files in the `patterns/` directory.

### Basic Pattern

```markdown
# 🔍 Code Review

> Review code for bugs and improvements

[Category: Review]

## Prompt

Analyze this code and provide feedback:

{{input}}

Include:
1. Bugs found
2. Improvements
3. Security issues
```

### Pattern with Variables

Add configurable fields that appear as a form before input:

```markdown
# 🎫 Create Ticket

> Generate a well-structured ticket

[Category: Tickets]

## Variables

- priority: Priority | select | low,medium,high,critical | medium
- type: Type | select | bug,feature,task | task
- env: Environment | select | dev,staging,prod | prod

## Prompt

Create a ticket with:
- Priority: {{var:priority}}
- Type: {{var:type}}
- Environment: {{var:env}}

Description:
{{input}}
```

### Variable Definition Format

```
- name: Label | type | options | default | placeholder
```

**Types:**
- `text` - Simple text input (default)
- `select` - Dropdown with options (comma-separated)
- `multiline` - Multi-line text area

**Examples:**
```markdown
- priority: Priority | select | low,medium,high | medium
- description: Description | text | | | Enter details...
- notes: Notes | multiline
```

### The {{input}} Variable

Control where user content is inserted:

```markdown
## Prompt

Analyze this code:

{{input}}

Provide detailed feedback.
```

**Without {{input}}**: User content is appended at the end.

## 🌐 Community Patterns

### Official Repository

Add the official patterns repository:

```bash
patternforge repo add https://github.com/Pierre-Malherbe/patternforge-patterns
patternforge upgrade
```

Or enable during first launch setup.

### Pattern Sources

Patterns are loaded from multiple sources with priority:
1. **Local patterns** (`./patterns`) - Highest priority
2. **Community patterns** (from repositories) - Merged, local wins on conflict

### Managing Repositories

```bash
# List all configured repos
patternforge repo list

# Add a new repository
patternforge repo add https://github.com/user/patterns

# Remove a repository
patternforge repo remove patterns

# Sync all repositories (also available via 'U' key in TUI)
patternforge upgrade
```

## 🏗️ Architecture

PatternForge follows **Clean Architecture** principles:

```
PatternForge/
├── cmd/patternforge/     # Application entrypoint
│   ├── main.go           # CLI & subcommands
│   └── model.go          # Main TUI model (state machine)
├── internal/
│   ├── pattern/          # Pattern loading, parsing, variables
│   ├── repository/       # Git operations for community patterns
│   ├── config/           # Settings & repository configuration
│   ├── claude/           # Claude Code executor
│   └── ui/
│       ├── screens/      # TUI screens (selection, variables, input, results)
│       └── styles/       # Centralized lipgloss styles
└── patterns/             # User's local pattern files
```

## 🔧 Development

### Running Tests

```bash
make test
```

### Building for macOS

```bash
# Build for both Apple Silicon and Intel
make build-macos

# Outputs:
# bin/patternforge-darwin-arm64
# bin/patternforge-darwin-amd64
```

## 🐛 Troubleshooting

**Q: "claude CLI not found"**
```bash
pip install claude-code
claude --version
```

**Q: "Pattern not showing up"**
- Check file is in `patterns/` directory
- Ensure filename ends with `.md`
- Press `U` to reload or restart PatternForge

**Q: "Variables not working"**
- Check `## Variables` section exists
- Format: `- name: Label | type | options | default`
- Use `{{var:name}}` in prompt

**Q: "Repository sync failed"**
- Ensure git is installed: `git --version`
- Check internet connection
- Verify repository URL is correct

## 📚 Resources

- **Claude Code Docs**: https://docs.anthropic.com/claude/docs/claude-code
- **Fabric Project**: https://github.com/danielmiessler/fabric
- **Official Patterns**: https://github.com/Pierre-Malherbe/patternforge-patterns
- **Bubbletea**: https://github.com/charmbracelet/bubbletea

## 🤝 Contributing

Contributions welcome! Areas for improvement:

- [ ] More community patterns
- [ ] Pattern validation
- [ ] Export results (PDF, HTML)
- [ ] Custom themes
- [ ] Windows support testing

## 📄 License

MIT License - See [LICENSE](LICENSE) file

## 🙏 Credits

- Inspired by [Fabric](https://github.com/danielmiessler/fabric) by Daniel Miessler
- Built with [Charm](https://charm.sh/) tools (Bubbletea, Lipgloss, Bubbles)

---

**Built with ❤️ in 🇫🇷 for the Claude Code community**

[⭐ Star on GitHub](https://github.com/Pierre-Malherbe/PatternForge) | [🐛 Report Bug](https://github.com/Pierre-Malherbe/PatternForge/issues) | [💡 Request Feature](https://github.com/Pierre-Malherbe/PatternForge/issues)
