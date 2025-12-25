# AI Terminal Helper - Quick Start Guide (Go Version)

**Get started in 5 minutes!** 🚀

---

## Prerequisites

### 1. Install Ollama

```bash
# macOS
brew install ollama

# Linux
curl -fsSL https://ollama.com/install.sh | sh

# Verify installation
ollama --version
```

### 2. Start Ollama Server

```bash
# Start in background
ollama serve &

# Or use system service (macOS)
brew services start ollama
```

### 3. Pull Required Models

```bash
# Essential models (download all 3 - REQUIRED)
ollama pull qwen3:8b-q4_K_M      # Primary (K8s, Terraform, AWS, Docker)
ollama pull gemma3:4b-it-q4_K_M  # Python/ML errors
ollama pull qwen3:4b-q4_K_M      # Fast fallback (default for unknown commands)

# Optional (for even faster responses on simple shell commands)
ollama pull qwen3:1.7b-q4_K_M    # Ultra-fast for cp, mv, rm, grep, find, etc.
# ⚠️  Only used for trivial shell commands (cp, mv, rm, mkdir, touch, grep, find, awk, sed)
# ⚠️  If not installed, falls back to qwen3:4b-q4_K_M (still fast!)

# Verify models
ollama list
```

**Expected output (minimum required):**
```
NAME                    ID              SIZE    MODIFIED
qwen3:8b-q4_K_M        abc123...       4.8 GB  2 minutes ago
gemma3:4b-it-q4_K_M    def456...       2.4 GB  3 minutes ago
qwen3:4b-q4_K_M        ghi789...       2.3 GB  4 minutes ago
```

**With optional 1.7B model:**
```
NAME                    ID              SIZE    MODIFIED
qwen3:8b-q4_K_M        abc123...       4.8 GB  2 minutes ago
gemma3:4b-it-q4_K_M    def456...       2.4 GB  3 minutes ago
qwen3:4b-q4_K_M        ghi789...       2.3 GB  4 minutes ago
qwen3:1.7b-q4_K_M      jkl012...       1.0 GB  5 minutes ago  ← Optional
```

### 📊 Model Selection Logic

The AI helper **automatically** selects the best model based on your command:

| Command Type                   | Model Used            | Why?                                | Example                       |
| ------------------------------ | --------------------- | ----------------------------------- | ----------------------------- |
| **Kubernetes, Terraform, AWS** | `qwen3:8b-q4_K_M`     | Complex configs need deep reasoning | `kubectl`, `terraform`, `aws` |
| **Docker, CI/CD**              | `qwen3:8b-q4_K_M`     | Complex orchestration               | `docker`, `jenkins`           |
| **Python, ML/Data**            | `gemma3:4b-it-q4_K_M` | Best for stack traces               | `python`, `pip`, `jupyter`    |
| **Simple shell commands**      | `qwen3:1.7b-q4_K_M`   | Ultra-fast (if installed)           | `cp`, `mv`, `rm`, `grep`      |
| **Unknown commands**           | `qwen3:4b-q4_K_M`     | Fast fallback                       | Anything else                 |
| **Proactive mode (`ask`)**     | `qwen3:8b-q4_K_M`     | Always use best quality             | All `ask` queries             |

**Key Points:**
- ✅ **3 models required:** `8b`, `4b-it`, `4b` (9.5 GB total)
- ⚠️ **1.7B is optional:** Only for trivial shell commands (saves ~0.5s per query)
- 🎯 **If 1.7B not installed:** Falls back to `4b` (still fast!)
- 🚀 **Proactive mode always uses 8B:** Better quality for generating commands

**Recommendation:**
- **Start with 3 required models** (skip 1.7B)
- **Add 1.7B later** if you want even faster responses for `cp`, `mv`, `rm`, etc.

---

## Installation

### Step 1: Clone Repository

```bash
cd ~/Ai  # or your preferred directory
git clone https://github.com/amaslovskyi/ai-helper.git
cd ai-helper
```

### Step 2: Build & Install

```bash
# Build and install (cleans old bash files automatically)
make install
```

**What this does:**
- ✅ Builds the Go binary (`bin/ai-helper`)
- ✅ Installs to `~/.ai/ai-helper`
- ✅ Copies ZSH integration to `~/.ai/ai-helper.zsh`
- ✅ Removes old bash scripts (`ai-helper.sh`, `cache-manager.sh`, `zsh-integration.sh`)
- ✅ Sets executable permissions

### Step 3: Add to Shell

```bash
# Add to ~/.zshrc (one-time setup)
echo 'source ~/.ai/ai-helper.zsh' >> ~/.zshrc

# Reload shell
source ~/.zshrc
```

**You should see:**
```
✅ AI Terminal Helper v2.0 (Go) Loaded!

Quick Commands:
  ai          - Re-analyze last failed command
  ask <query> - Generate command from natural language
  kask/dask/task/gask - Tool-specific queries

Features:
  🔒 Security scanning    ⏱️  Smart rate limiting
  ✅ Command validation   💾 Offline caching
  🚀 Fast Go binary       🎯 Fixes hallucinations
```

---

## Verify Installation

### Test 1: Check Binary

```bash
which ai-helper
# Expected: /Users/yourusername/.ai/ai-helper

ai-helper version
# Expected: AI Terminal Helper 2.0.0-go
```

### Test 2: Test Validation (The Big Feature!)

```bash
# This will catch AI hallucination!
ask list docker containers sorted by memory
```

**Expected output:**
```
🤖 Generating command for: list docker containers sorted by memory
⚠️  Validation failed: docker ps does not have a --sort flag
ℹ️  Querying AI again with validation context...
✓ docker stats --no-stream | sort -k2 -hr
Root: docker ps doesn't support sorting, use docker stats instead
Tip: Add --format for custom output columns
```

**Without validation, you'd get:**
```
✓ docker ps --sort memusage  # ❌ This doesn't exist!
```

### Test 3: Test Reactive Mode

```bash
# Run a command that will fail
kubectl get pods --invalid-flag
```

**Expected output:**
```
unknown flag: --invalid-flag

🤖 AI Assistant (exit 1):
✓ kubectl get pods
Root: Invalid flag removed
Tip: Use kubectl get pods --help for valid flags
```

### Test 4: Test Proactive Mode

```bash
ask how do I list all running pods
```

**Expected output:**
```
🤖 Generating command for: how do I list all running pods
✓ kubectl get pods --field-selector=status.phase=Running
Root: Lists only pods in Running state
Tip: Add -A for all namespaces
```

---

## What Gets Installed

### Files in `~/.ai/`

After installation, you should have **only these files**:

```bash
~/.ai/
├── ai-helper          # Go binary (8MB)
├── ai-helper.zsh      # ZSH integration (~100 lines)
├── cache.json         # Cache (created on first use)
├── history.log        # History (created on first use)
└── rate_limit.log     # Rate limiting (created on first use)
```

**Old bash files are automatically removed:**
- ❌ `ai-helper.sh` (replaced by Go binary)
- ❌ `cache-manager.sh` (built into Go binary)
- ❌ `zsh-integration.sh` (replaced by `ai-helper.zsh`)

---

## Usage Examples

### Reactive Mode (Automatic)

AI triggers automatically when commands fail:

```bash
# Example 1: Kubernetes
$ kubectl get pods --invalid-flag
unknown flag: --invalid-flag

🤖 AI Assistant (exit 1):
✓ kubectl get pods
Root: Invalid flag removed
Tip: Use kubectl get pods --help for valid flags

# Example 2: Docker
$ docker ps --sort memusage
unknown flag: --sort

🤖 AI Assistant (exit 125):
⚠️  Validation failed: docker ps does not have a --sort flag
ℹ️  Querying AI again with validation context...
✓ docker stats --no-stream | sort -k2 -hr
Root: docker ps doesn't support sorting, use docker stats instead
Tip: Add --format for custom output columns

# Example 3: Git
$ git push --force-with-lease=wrong
fatal: invalid argument to --force-with-lease

🤖 AI Assistant (exit 128):
✓ git push --force-with-lease
Root: --force-with-lease doesn't take arguments in this context
Tip: Use --force-with-lease=<refname>:<expect> for specific refs
```

### Proactive Mode (Ask Before Running)

Generate commands from natural language:

```bash
# General query
$ ask how do I list all running containers
🤖 Generating command for: how do I list all running containers
✓ docker ps
Root: Lists all running Docker containers
Tip: Add -a to see stopped containers too

# Kubernetes query
$ kask show pods in production
🤖 Generating command for: kubernetes: show pods in production
✓ kubectl get pods -n production
Root: Lists all pods in the production namespace
Tip: Add -o wide for more details

# Docker query
$ dask list containers by memory
🤖 Generating command for: docker: list containers by memory
⚠️  Validation failed: docker ps does not have a --sort flag
ℹ️  Querying AI again with validation context...
✓ docker stats --no-stream | sort -k2 -hr
Root: docker stats shows memory usage, sorted by column 2
Tip: Use --format to customize output

# Terraform query
$ task how do I plan changes
🤖 Generating command for: terraform: how do I plan changes
✓ terraform plan
Root: Shows execution plan without applying changes
Tip: Use -out=plan.tfplan to save the plan

# Git query
$ gask how do I undo last commit
🤖 Generating command for: git: how do I undo last commit
✓ git reset --soft HEAD~1
Root: Undoes last commit but keeps changes staged
Tip: Use --hard to discard changes completely
```

### Manual Trigger

```bash
# Re-analyze last failed command
$ ai

# Show cache stats
$ ai-stats

# Clear cache
$ ai-clear

# Show version
$ ai-version
```

### Hotkeys

- **⌥A (Option+A)** - Re-analyze last failure
- **⌥K (Option+K)** - Quick ask mode (opens `ask ` prompt)

---

## Configuration

### Environment Variables

```bash
# Optional: Customize Ollama API endpoint
export OLLAMA_HOST="http://localhost:11434"

# Optional: Increase timeout for slow models
export AI_HELPER_TIMEOUT=10  # seconds (default: 5)
```

### Cache Settings

Cache is stored in `~/.ai/cache.json` with:
- **Max entries:** 100 (LRU eviction)
- **Format:** JSON
- **Size:** ~50KB typical

```bash
# View cache stats
ai-stats

# Clear cache
ai-clear
```

---

## Troubleshooting

### Issue 1: `command not found: ai-helper`

**Cause:** Binary not in PATH or not installed

**Fix:**
```bash
# Reinstall
cd /path/to/ai-helper
make install

# Reload shell
source ~/.zshrc

# Verify
which ai-helper
# Should show: /Users/yourusername/.ai/ai-helper
```

### Issue 2: `Ollama server not running`

**Cause:** Ollama service not started

**Fix:**
```bash
# Start Ollama
ollama serve &

# Or use system service
brew services start ollama

# Verify
curl http://localhost:11434
# Should return: Ollama is running
```

### Issue 3: `Model not found`

**Cause:** Required model not pulled

**Fix:**
```bash
# Pull missing model
ollama pull qwen3:8b-q4_K_M

# Verify
ollama list
```

### Issue 4: Old bash files still present

**Cause:** Previous installation not cleaned

**Fix:**
```bash
# Clean old files
rm -f ~/.ai/ai-helper.sh ~/.ai/cache-manager.sh ~/.ai/zsh-integration.sh

# Reinstall
cd /path/to/ai-helper
make install
```

### Issue 5: Validation not working

**Cause:** Using old bash version

**Fix:**
```bash
# Check version
ai-helper version
# Should show: 2.0.0-go (not 2.0.0-bash)

# If wrong version, reinstall
cd /path/to/ai-helper
make install
source ~/.zshrc
```

---

## ✨ Key Benefits

### 🎯 Command Validation
- **Catches AI hallucinations** before they reach you
- **Automatic re-querying** when validation fails
- **Prevents invalid commands** like `docker ps --sort`

### 🚀 Performance
- **10x faster** than bash version
- **5ms startup** time
- **0.5ms cache lookups**
- **Single 5.5MB binary**

### 🔒 Privacy & Security
- **100% local** - No cloud, no telemetry
- **Safe for secrets** - AWS keys, k8s tokens, DB passwords
- **18 dangerous patterns** detected and blocked
- **Air-gapped friendly** - Works completely offline

### 🛠️ Developer Experience
- **Easy to test** - `go test` for unit tests
- **Type-safe** - Compile-time error checking
- **Easy to distribute** - Single binary, no dependencies
- **Extensible** - Add new validators easily

### 💡 Smart Features
- **Intelligent model routing** - Picks best model for each command
- **Offline caching** - 40-60% faster for repeated errors
- **Rate limiting** - Prevents AI spam
- **Proactive mode** - Generate commands from natural language
- **Colorful output** - Easy to read terminal output
- **Hotkey support** - ⌥A and ⌥K for quick access

---

## 📊 Performance

| Metric            | Go Binary | Bash Scripts  | Improvement    |
| ----------------- | --------- | ------------- | -------------- |
| **Startup**       | ~5ms      | ~50ms         | **10x faster** |
| **Cache lookup**  | ~0.5ms    | ~5ms          | **10x faster** |
| **Security scan** | ~1ms      | ~10ms         | **10x faster** |
| **Validation**    | ~1ms      | ❌ N/A         | **NEW!**       |
| **Binary size**   | 8.0 MB    | N/A (5 files) | Single file    |

---

## Next Steps

### 1. Read Full Documentation
- [README.md](README.md) - Complete overview
- [ROADMAP.md](ROADMAP.md) - Future features
- [GO-MIGRATION-GUIDE.md](GO-MIGRATION-GUIDE.md) - Migrate from bash

### 2. Explore Advanced Features
```bash
# Test security scanning
ask how do I delete all files recursively
# Should warn about dangerous command

# Test rate limiting
# Run same failing command 3+ times quickly
# Should pause AI suggestions

# Test caching
# Run same failing command twice
# Second time should show "💾 [Cached]"
```

### 3. Customize
- Add custom validators (see [README.md](README.md))
- Adjust cache size in `pkg/cache/cache.go`
- Add more dangerous patterns in `pkg/security/scanner.go`

---

## Uninstall

```bash
# Remove installation
make uninstall

# Remove from ~/.zshrc
# Delete this line: source ~/.ai/ai-helper.zsh

# Remove data (optional)
rm -rf ~/.ai/
```

---

## Support

- **Issues:** [GitHub Issues](https://github.com/amaslovskyi/ai-helper/issues)
- **Discussions:** [GitHub Discussions](https://github.com/amaslovskyi/ai-helper/discussions)
- **Documentation:** [README.md](README.md)

---

**Built with ❤️ in Go. No more hallucinations!** 🚀
