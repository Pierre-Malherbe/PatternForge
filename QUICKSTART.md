# 🚀 Quick Start - Testing PatternForge

## Step-by-Step Setup (macOS)

### 1. Prerequisites Check

```bash
# Check Go version (need 1.21+)
go version

# If not installed:
brew install go
```

### 2. Install Claude Code

```bash
# Install via pip
pip install claude-code

# Authenticate
claude auth login

# Verify it works
claude --version
```

### 3. Build PatternForge

```bash
# Navigate to the project
cd patternforge

# Install dependencies
make deps

# Build the binary
make build

# You should see: ✅ Build complete: bin/patternforge
```

### 4. Run It!

```bash
# Run from the project directory
./bin/patternforge

# Or install system-wide
make install

# Then run from anywhere
patternforge
```

## First Run Experience

When you run PatternForge for the first time:

```
✨ PatternForge initialized!
✨ Created default pattern: mermaid-map.md
```

You'll see the pattern selection screen:

```
╔══════════════════════════════════════════╗
║  🎨 PatternForge - Select a pattern      ║
╠══════════════════════════════════════════╣
║  > 📊 Mermaid Map                [1/3]  ║
║    Create detailed mind maps            ║
║                                          ║
║    📝 Summarize                          ║
║    💻 Explain Code                       ║
╚══════════════════════════════════════════╝
```

## Testing the Workflow

### Test 1: Create a Mermaid Map

1. **Select pattern**: Press `enter` on "Mermaid Map"
2. **Paste content**:
   ```
   Modern web development:
   - Frontend: React, Vue, Svelte
   - Backend: Node.js, Python, Go
   - DevOps: Docker, Kubernetes
   ```
3. **Process**: Press `Ctrl+D`
4. **Wait**: Spinner shows while Claude processes (~2-5 seconds)
5. **Results**: See your Mermaid diagram with token stats!
6. **Copy**: Press `y` to copy to clipboard

### Test 2: Create Your Own Pattern

1. **In selection screen**: Press `n`
2. **Vi opens**: Edit the template
   ```markdown
   # 🎯 My Pattern
   
   > What it does
   
   ## Prompt
   
   Do something with:
   {{input}}
   
   And provide...
   ```
3. **Save**: Type `:wq` in Vi
4. **See it appear**: Your pattern is now in the list!
5. **Use it**: Select and test

### Test 3: Modify Existing Pattern

1. **Select a pattern**: Navigate with `j/k`
2. **Press `m`**: Vi opens the pattern file
3. **Edit**: Change description, emoji, or prompt
4. **Save**: `:wq`
5. **Reloaded**: Changes appear immediately

## Troubleshooting

### Build Fails

```bash
# Clean and retry
make clean
rm go.sum
make deps
make build
```

### "claude: command not found"

```bash
# Check if pip installed correctly
which claude

# Try reinstalling
pip uninstall claude-code
pip install claude-code

# Check PATH
echo $PATH | grep -o '[^:]*bin'
```

### Import Errors

```bash
# The project uses internal/ directory
# Make sure you're building from project root
pwd  # Should end in /patternforge

# Module name must match
head -1 go.mod  # Should show: module github.com/Pierre-Malherbe/patternforge
```

### Vi/Vim Issues

```bash
# Set your preferred editor
export EDITOR=nano  # or vim, vi, code, etc.

# Or install vim
brew install vim
```

## Verification Checklist

Before considering it "working":

- [ ] `make build` completes without errors
- [ ] `./bin/patternforge` launches the TUI
- [ ] Can navigate patterns with `j/k`
- [ ] Can select pattern with `enter`
- [ ] Can input text and process with `Ctrl+D`
- [ ] See results with token stats
- [ ] Can copy with `y`
- [ ] Can create new pattern with `n`
- [ ] Can modify pattern with `m`
- [ ] Can quit with `q`

## Performance Notes

**Expected Performance:**
- **Startup**: Instant (<100ms)
- **Pattern loading**: Instant (even with 100+ patterns)
- **Claude processing**: 2-10 seconds depending on complexity
- **UI responsiveness**: Buttery smooth (60fps)

**Token Estimation:**
- Input/output tokens are estimated (1 token ≈ 4 chars)
- Actual tokens may vary slightly
- Duration is real execution time

## Next Steps

Once everything works:

1. **Create patterns** for your common tasks
2. **Organize** patterns in categories (use prefixes)
3. **Share** your best patterns with the team
4. **Backup** your patterns/ directory (git recommended)
5. **Customize** styles in `internal/ui/styles/styles.go`

## Development Mode

For rapid iteration during development:

```bash
# Run without building
make run

# Watch for changes (requires entr)
ls **/*.go | entr -r make run

# Build for macOS (both architectures)
make build-macos
```

## Getting Help

If stuck:
1. Check README.md for detailed docs
2. Review pattern examples in `patterns/`
3. Check Go logs: `go build -v ./...`
4. Verify Claude Code: `claude ask "test"`

---

**Happy pattern crafting! 🎨**
