# 🎨 PatternForge

**Dynamic prompt patterns powered by Claude Code**

PatternForge is a blazing-fast TUI (Terminal User Interface) that brings the power of [Fabric](https://github.com/danielmiessler/fabric)-style patterns to Claude Code. Create, edit, and manage reusable AI prompts with an elegant interface and Vim-style keybindings.

![Version](https://img.shields.io/badge/version-0.0.1-blue)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

## ✨ Features

- 🎯 **Dynamic Pattern Loading** - Patterns are markdown files you can edit anytime
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

# 2. Select a pattern (j/k to navigate)
> 📊 Mermaid Map

# 3. Paste your content, press Ctrl+D

# 4. Get beautiful results with stats
✨ Result
[Your Mermaid diagram here]
📊 172 tokens input + 612 output = 784 total | 3.2s

# 5. Press 'y' to copy output
```

## 📦 Installation

### Prerequisites

1. **Go 1.21+** - [Install Go](https://go.dev/doc/install)
2. **Claude Code** - Required for AI processing
   ```bash
   pip install claude-code
   claude auth login
   ```
3. **Vi/Vim** - For editing patterns (usually pre-installed on macOS/Linux)

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

### macOS Binary (Coming Soon)

Pre-built binaries for Apple Silicon and Intel Macs will be available on the releases page.

## 🚀 Usage

### Basic Usage

```bash
# Use default patterns directory (./patterns)
patternforge

# Use custom patterns directory
patternforge ~/my-patterns
```

### First Run

On first run, PatternForge will:
1. ✅ Check if Claude Code is installed and functional
2. 📁 Create a `patterns/` directory
3. 📝 Add a default Mermaid Map pattern
4. 🎯 Launch the TUI

### Workflow

1. **Select Pattern** - Use `j/k` or arrow keys, press `enter`
2. **Input Content** - Paste or type your content
3. **Process** - Press `Ctrl+D` to send to Claude Code
4. **View Results** - Scroll through output, see token stats
5. **Copy/Save** - Press `y` to copy, or `s` to save to file

## ⌨️  Keyboard Shortcuts

### Pattern Selection Screen

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `enter` | Select pattern |
| `m` | Modify selected pattern (opens Vi) |
| `n` | Create new pattern (opens Vi) |
| `/` | Search patterns |
| `q` | Quit |

### Input Screen

| Key | Action |
|-----|--------|
| `Ctrl+D` | Process with Claude Code |
| `Esc` | Back to pattern selection |
| `m` | Edit current pattern |

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

### Pattern Template

```markdown
# [EMOJI] [PATTERN NAME]

> [One-line description]

## Prompt

[Your prompt instructions]

{{input}}

[More instructions after user content]
```

### The {{input}} Variable

Control where user content is inserted:

```markdown
## Prompt

Analyze this code:

{{input}}

Provide:
1. Bugs found
2. Improvements
3. Security issues
```

**Without {{input}}**: User content is appended at the end.

### Creating a New Pattern

**Method 1: Via TUI**
```bash
# In PatternForge, press 'n'
# Vi opens with template
# Edit, save (:wq), done!
```

**Method 2: Manually**
```bash
# Create a file in patterns/
vi patterns/my-pattern.md

# Write your pattern
# Save and reload PatternForge
```

### Example Patterns

**Code Review Pattern:**
```markdown
# 🔍 Code Review

> Comprehensive code review with scoring

## Prompt

Perform a professional code review:

{{input}}

Provide:
1. Overall score (X/10)
2. Strengths (3-5 points)
3. Issues (with severity levels)
4. Refactored version
5. Learning points
```

**Summarize Pattern:**
```markdown
# 📝 Summarize

> Create concise summaries

## Prompt

Summarize this content concisely:

{{input}}

Max 150 words, keep essential points only.
```

## 🏗️  Architecture

PatternForge follows **Clean Architecture** principles for maintainability:

```
PatternForge/
├── cmd/patternforge/     # Entrypoint, main application
│   ├── main.go           # CLI initialization
│   └── model.go          # Main TUI model (state machine)
├── internal/
│   ├── pattern/          # Pattern loading & parsing
│   ├── claude/           # Claude Code executor
│   └── ui/
│       ├── screens/      # Individual screens (selection, input, results)
│       ├── components/   # Reusable UI components
│       └── styles/       # Centralized lipgloss styles
└── patterns/             # User's pattern files
```

**Why this structure?**
- **Testable**: Each package has single responsibility
- **Maintainable**: Clear separation of concerns
- **Extensible**: Easy to add new features or patterns
- **Clean**: No circular dependencies

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

### Project Structure Explained

**`internal/pattern`**: Core domain logic
- Loads .md files from disk
- Parses markdown structure
- Handles {{input}} templating

**`internal/claude`**: External service integration
- Executes Claude Code CLI
- Manages temp files
- Estimates token usage

**`internal/ui/screens`**: View layer (Bubbletea screens)
- `selection.go`: Pattern list
- `input.go`: Content input
- `processing.go`: Loading spinner
- `results.go`: Output display

**`cmd/patternforge`**: Application layer
- Ties everything together
- State machine (view transitions)
- Keyboard shortcut handling

## 🎯 Common Use Cases

### Development

- **Code Review**: Analyze code for bugs and improvements
- **Explain Code**: Break down complex algorithms
- **Generate Docs**: Create documentation from code
- **Write Tests**: Generate unit tests

### Content Creation

- **Blog Posts**: Draft articles with AI assistance
- **Social Media**: Create engaging posts
- **Email Templates**: Professional email composition
- **SEO Content**: Optimize content for search

### Data & Analysis

- **Summarize Reports**: Condense long documents
- **Extract Insights**: Pull key findings from data
- **Create Visualizations**: Generate Mermaid diagrams
- **Data Stories**: Transform numbers into narratives

## 🐛 Troubleshooting

**Q: "claude CLI not found"**
```bash
# Install Claude Code
pip install claude-code

# Verify installation
claude --version
```

**Q: "Pattern not showing up"**
- Check file is in `patterns/` directory
- Ensure filename ends with `.md`
- Restart PatternForge or reload patterns

**Q: "Vi/Vim not found"**
```bash
# macOS (should be pre-installed)
which vim

# If missing, install via Homebrew
brew install vim
```

**Q: "Copy (y) not working"**
- Ensure clipboard tools are installed
- macOS: Uses `pbcopy` (built-in)
- Linux: Install `xclip` or `xsel`

**Q: "{{input}} not replaced"**
- Check spelling: must be exactly `{{input}}`
- Case-sensitive
- No spaces inside braces

## 📚 Resources

- **Claude Code Docs**: https://docs.anthropic.com/claude/docs/claude-code
- **Fabric Project**: https://github.com/danielmiessler/fabric
- **Bubbletea**: https://github.com/charmbracelet/bubbletea
- **Pattern Examples**: See `patterns/` directory

## 🤝 Contributing

Contributions welcome! Areas for improvement:

- [ ] More default patterns
- [ ] Pattern validation
- [ ] Export results (PDF, HTML)
- [ ] Pattern marketplace/sharing
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
