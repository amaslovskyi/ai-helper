# Changelog

All notable changes to AI Terminal Helper will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [2.0.0] - 2025-12-25

### 🎉 Major Release: Complete Go Rewrite

This is a **complete rewrite** from bash scripts to Go, providing significant improvements in performance, reliability, and features.

### ✨ Added

#### Core Features
- **Command Validation** - Prevents AI hallucinations by validating suggested commands
  - Docker validator catches non-existent flags (e.g., `docker ps --sort`)
  - Automatic re-querying when validation fails
  - Extensible validator architecture for adding more tools
- **Security Scanning** - 18 dangerous command patterns detected and blocked
  - Prevents destructive commands (`rm -rf /`, `DROP DATABASE`, etc.)
  - Warns about insecure permissions (`chmod 777`)
  - Blocks fork bombs and other malicious patterns
- **Smart Caching** - JSON-based offline cache with LRU eviction
  - 100 entry limit with automatic cleanup
  - Hit tracking and statistics
  - 40-60% faster responses for repeated errors
- **Rate Limiting** - Prevents AI spam on repeated failures
  - Configurable cooldown period (default: 2s)
  - Per-command tracking
  - Helpful tips when rate limit triggered
- **Proactive Mode** - Generate commands from natural language
  - `ask` - General queries
  - `kask` - Kubernetes-specific queries
  - `dask` - Docker-specific queries
  - `task` - Terraform-specific queries
  - `gask` - Git-specific queries
- **Colorful Output** - ANSI color support for better readability
  - Green for suggestions
  - Cyan for root causes
  - Yellow for tips
  - Red for errors and warnings
- **Hotkey Support** - Quick access to AI features
  - ⌥A (Option+A) - Re-analyze last failure
  - ⌥K (Option+K) - Quick ask mode

#### Model Routing
- **Intelligent Model Selection** - Automatically chooses best model based on command context
  - `qwen3:8b-q4_K_M` - Kubernetes, Terraform, AWS, Docker (complex configs)
  - `gemma3:4b-it-q4_K_M` - Python, ML/Data (stack traces)
  - `qwen3:4b-q4_K_M` - Fast fallback for unknown commands
  - `qwen3:1.7b-q4_K_M` - Optional ultra-fast for simple shell commands
- **Proactive Mode Always Uses 8B** - Best quality for generating commands

#### Architecture
- **Modular Design** - Clean separation of concerns
  - `pkg/llm/` - Ollama integration and model routing
  - `pkg/validators/` - Command validators (extensible)
  - `pkg/security/` - Security scanning
  - `pkg/cache/` - Cache management
  - `pkg/ui/` - Terminal UI and colors
- **Single Binary** - 5.5 MB compiled Go binary
- **Minimal Integration** - ~110 line ZSH integration script
- **Cross-platform Support** - macOS (amd64, arm64), Linux (amd64, arm64)

#### Build System
- **Makefile** - Automated build and installation
  - `make build` - Build binary
  - `make install` - Install to `~/.ai/` with automatic cleanup
  - `make uninstall` - Clean removal
  - `make test` - Run tests
  - `make build-all` - Cross-compile for all platforms
  - `make clean` - Remove build artifacts

#### Documentation
- **Comprehensive Docs** - Complete documentation suite
  - `README.md` - Main documentation with feature highlights
  - `QUICKSTART.md` - 5-minute setup guide
  - `ROADMAP.md` - Future plans and feature matrix
  - `GO-MIGRATION-GUIDE.md` - Migration guide from bash
  - `bash/README.md` - Archive notice for bash version

### 🚀 Performance Improvements

- **10x Faster Startup** - 5ms vs 50ms (bash)
- **10x Faster Cache Lookups** - 0.5ms vs 5ms (bash)
- **10x Faster Security Scanning** - 1ms vs 10ms (bash)
- **Instant Command Validation** - <1ms (new feature)

### 🔒 Security Improvements

- **Command Validation** - Prevents execution of hallucinated commands
- **18 Dangerous Patterns** - Comprehensive security scanning
- **100% Local** - No cloud, no telemetry, no data leakage
- **Safe for Secrets** - Works with AWS keys, k8s tokens, DB passwords

### 🛠️ Changed

- **Installation Location** - Now installs to `~/.ai/` (was scattered)
- **Binary Name** - `ai-helper` (was `ai-helper.sh`)
- **Integration File** - `ai-helper.zsh` (was `zsh-integration.sh`)
- **Cache Format** - JSON (was custom format)
- **Automatic Cleanup** - Old bash files removed during installation

### 🗑️ Removed

- **Bash Scripts** - Replaced with Go binary
  - `ai-helper.sh` → Go binary
  - `cache-manager.sh` → Built into Go binary
  - `zsh-integration.sh` → Minimal `ai-helper.zsh`
- **Subprocess Overhead** - No more bash subprocess calls
- **Script Dependencies** - Self-contained binary

### 🐛 Fixed

- **AI Hallucinations** - Command validation prevents invalid suggestions
- **Rate Limit Bugs** - Proper tracking and cleanup
- **Cache Corruption** - Robust JSON parsing with error handling
- **PATH Issues** - Automatic PATH management in ZSH integration

### 📊 Comparison: Go vs Bash

| Aspect | Go Binary | Bash Scripts |
|--------|-----------|--------------|
| **Hallucination Prevention** | ✅ Yes | ❌ No |
| **Performance** | ✅ 10x faster | ⚠️ Slow |
| **Distribution** | ✅ Single 5.5MB binary | ⚠️ 5 files |
| **Testing** | ✅ Easy (`go test`) | ❌ Difficult |
| **Parsing** | ✅ Real parsers | ❌ Regex only |
| **Type Safety** | ✅ Compile-time | ❌ Runtime |
| **Maintainability** | ✅ Clean architecture | ⚠️ Complex bash |

### 🔄 Migration from Bash

See [GO-MIGRATION-GUIDE.md](GO-MIGRATION-GUIDE.md) for detailed migration instructions.

**Quick Migration:**
```bash
cd /path/to/ai-helper
make install  # Automatically removes old bash files
source ~/.zshrc
```

### 📦 Installation

**Requirements:**
- Go 1.21+ (for building)
- Ollama with models:
  - `qwen3:8b-q4_K_M` (required)
  - `gemma3:4b-it-q4_K_M` (required)
  - `qwen3:4b-q4_K_M` (required)
  - `qwen3:1.7b-q4_K_M` (optional)

**Install:**
```bash
git clone https://github.com/yourusername/ai-helper.git
cd ai-helper
make install
echo 'source ~/.ai/ai-helper.zsh' >> ~/.zshrc
source ~/.zshrc
```

### 🙏 Acknowledgments

- Built with [Ollama](https://ollama.ai) for local LLM inference
- Inspired by Warp Terminal's AI features
- Designed for DevOps/SRE/MLOps professionals

---

## [1.0.0] - 2024-XX-XX (Bash Version - Archived)

### Initial Release (Bash Implementation)

The original bash implementation has been archived to `bash/` folder.
See `bash/CHANGELOG-bash.md` for bash version history.

**Key Features (Bash):**
- Reactive mode (automatic error fixing)
- Basic caching
- ZSH integration
- Model routing
- Security scanning

**Why We Moved to Go:**
- ❌ No command validation (AI hallucinations reached users)
- ❌ Slow performance (subprocess overhead)
- ❌ Hard to test and maintain
- ❌ Limited parsing capabilities

---

## Versioning

We use [Semantic Versioning](https://semver.org/):
- **MAJOR** version for incompatible API changes
- **MINOR** version for new functionality (backwards-compatible)
- **PATCH** version for backwards-compatible bug fixes

---

## Links

- **Repository:** https://github.com/yourusername/ai-helper
- **Issues:** https://github.com/yourusername/ai-helper/issues
- **Discussions:** https://github.com/yourusername/ai-helper/discussions
- **Documentation:** [README.md](README.md)
- **Quick Start:** [QUICKSTART.md](QUICKSTART.md)
- **Roadmap:** [ROADMAP.md](ROADMAP.md)

---

**Built with ❤️ in Go. No more hallucinations!** 🚀

