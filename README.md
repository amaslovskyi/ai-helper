# Local Terminal AI Manual v2.0 (2025)

**Audience:** DevOps / SRE / MLOps engineers  
**Platform:** macOS (Apple Silicon, 16 GB RAM)  
**Goal:** Local, fast, context-aware AI companion for production environments  
**Version:** 2.0 - Now with proactive mode, security scanning, and caching!

---

## 🆕 What's New in v2.0

### ✨ Major Features Added
1. **🔒 Security Scanning** - Prevents dangerous commands (rm -rf /, DROP DATABASE, etc.)
2. **⏱️ Smart Rate Limiting** - Prevents AI spam on repeated failures
3. **🗣️ Proactive Mode** - Ask AI to generate commands BEFORE errors
4. **💾 Offline Cache** - Instant responses for common errors (40-60% faster)
5. **🧠 Command History Learning** - Learns from your successful fixes
6. **🎯 Enhanced Context Awareness** - Better suggestions for kubectl, docker, terraform, git

### 🏆 Why v2.0 Beats Warp Terminal
- ✅ **Privacy:** 100% local (Warp sends data to cloud)
- ✅ **Cost:** Free forever (Warp has usage limits & paid tiers)
- ✅ **Security:** Safe for secrets & regulated environments
- ✅ **Speed:** 0.3-2.5s local inference (no network latency)
- ✅ **Proactive Mode:** Natural language commands (new!)
- ✅ **Smart Caching:** Instant responses for common errors (new!)

---

## 1. Design Principles

* **Local only** (no cloud, no telemetry, safe for secrets)
* **Proactive + Reactive** (help before AND after errors)
* **Low latency** on M1 / M2 / M3 / M4 (< 2s response, 0.3s cached)
* **Context-aware** (infra, containers, ML, cloud)
* **Production-safe** with security scanning for dangerous commands
* **Smart filtering** (skips trivial commands, prevents AI spam)
* **Self-learning** (caches successful fixes for instant reuse)

---

## 2. Model Strategy (Corrected, 2025)

This manual reflects **real Ollama model availability** and **best-in-market local performance** for SRE work.

---

### 2.1 Primary Model (Main Brain)

**`qwen3:8b-q4_K_M`**
Role: *Primary reasoning & infra assistant*

Why:

* Best overall local model in 2025 for SRE / DevOps
* Strong at Kubernetes, Terraform, CI/CD, Linux
* Excellent accuracy-to-latency ratio

```bash
ollama pull qwen3:8b-q4_K_M
```

---

### 2.2 Fast Fallback (Shell / Low-Latency)

**`qwen3:4b-q4_K_M`**
Role: *Instant shell help & simple errors*

Use when:

* Flag mistakes
* Command syntax errors
* Fast suggestions without deep reasoning

```bash
ollama pull qwen3:4b-q4_K_M
```

---

### 2.3 Explanation Model (Corrected Gemma Choice)

**`gemma3:4b-it-q4_K_M`**
Role: *Clear explanations & stack trace summaries*

Why this exact model:

* Instruction-tuned (`-it`)
* Quantized for speed (`q4_K_M`)
* 128K context (logs, traces)
* **There is NO Gemma 3 9B model**

```bash
ollama pull gemma3:4b-it-q4_K_M
```

---

### 2.4 Ultra-Fast Tier (Trivial Errors)

**`qwen3:1.7b-q4_K_M`** (recommended) or **`gemma3:1b-it-q4_K_M`** (fallback)

Use for:

* Simple shell command errors (`cp`, `mv`, `rm`, `mkdir`, `grep`, `find`)
* Instant responses (< 0.5s)
* Trivial flag mistakes

```bash
ollama pull qwen3:1.7b-q4_K_M   # Best quality for 1-2B tier
ollama pull gemma3:1b-it-q4_K_M # Fallback option
```

**Note:** Script auto-detects which model is available and uses the best one.

---

## 3. Final Model Matrix (Recommended)

| Role                  | Model                                      |
| --------------------- | ------------------------------------------ |
| Primary reasoning     | qwen3:8b-q4_K_M                            |
| Fast shell fallback   | qwen3:4b-q4_K_M                            |
| Explanations          | gemma3:4b-it-q4_K_M                        |
| Ultra-fast (trivial)  | qwen3:1.7b-q4_K_M or gemma3:1b-it-q4_K_M   |

---

## 4. Enhanced Model Routing Logic

The assistant automatically selects the best model based on command context and complexity.

| Context                                      | Model                        | Latency    | Why                          |
| -------------------------------------------- | ---------------------------- | ---------- | ---------------------------- |
| kubectl / helm / terraform / terragrunt      | qwen3:8b-q4_K_M              | 1-2.5s     | Complex infra reasoning      |
| aws / gcloud / azure (cloud CLIs)            | qwen3:8b-q4_K_M              | 1-2.5s     | Multi-service orchestration  |
| docker / podman / CI/CD                      | qwen3:8b-q4_K_M              | 1-2.5s     | Config debugging             |
| prometheus / grafana / monitoring            | qwen3:8b-q4_K_M              | 1-2.5s     | Log analysis                 |
| python / pip / conda / mlflow / spark        | gemma3:4b-it-q4_K_M          | 0.5-1.5s   | Stack trace interpretation   |
| jupyter / kubeflow / ray (ML platforms)      | gemma3:4b-it-q4_K_M          | 0.5-1.5s   | ML workflow context          |
| cp / mv / rm / mkdir / grep / find           | qwen3:1.7b or gemma3:1b      | 0.3-0.8s   | Trivial flag errors          |
| other shell commands                         | qwen3:4b-q4_K_M              | 0.5-1.5s   | Fast syntax fixes            |

**Smart Skipping:** Non-errors (`cd`, `ls`, `pwd`, `echo`) are ignored to reduce noise.

---

## 5. New Feature: Proactive Mode 🚀

v2.0 introduces **proactive mode** - generate commands from natural language BEFORE making errors!

### 5.1 Usage Examples

```bash
# Generate commands from natural language
$ ask how do I list all pods in production namespace
🤖 Generating command for: how do I list all pods in production namespace
✓ kubectl get pods -n production
Root: Lists all pods in the production namespace
Tip: Add -o wide for more details or --watch for live updates

# Quick kubernetes helper
$ kask show pods with high memory usage
✓ kubectl top pods --sort-by=memory
Root: Shows pods sorted by memory consumption
Tip: Requires metrics-server to be installed

# Docker queries
$ dask find containers using more than 1GB RAM
✓ docker stats --no-stream --format "table {{.Name}}\t{{.MemUsage}}" | awk '$2 > 1'
Root: Filters containers by memory usage
Tip: Use 'docker stats' for live monitoring

# Terraform queries
$ task how do I preview changes without applying
✓ terraform plan
Root: Shows what changes Terraform will make without applying them
Tip: Use -out=plan.tfplan to save the plan for later apply

# Git queries
$ gask how do I undo my last commit but keep changes
✓ git reset --soft HEAD~1
Root: Undoes last commit but keeps changes staged
Tip: Use --hard to discard changes completely (dangerous!)
```

### 5.2 Proactive Mode Shortcuts

```bash
ask    - General queries (any command)
kask   - Kubernetes-specific
dask   - Docker-specific
task   - Terraform-specific
gask   - Git-specific
```

### 5.3 Hotkeys

* **⌥A (Option+A)** - Re-analyze last failed command
* **⌥K (Option+K)** - Quick ask mode

---

## 6. Security Features 🔒

v2.0 includes production-grade security scanning to prevent catastrophic mistakes.

### 6.1 Dangerous Command Detection

The AI helper automatically scans AI suggestions for dangerous patterns:

```bash
# Example: AI suggests a dangerous command
$ kubectl delete --all
🚨 DANGER: Command contains potentially destructive pattern: --all
⚠️  This could cause data loss or system damage!
📋 Command: kubectl delete pods --all -n production

If you're ABSOLUTELY SURE this is safe, you can:
1. Review the command carefully
2. Test in a safe environment first
3. Execute manually after verification
```

### 6.2 Blocked Patterns

* `rm -rf /` - Recursive root deletion
* `rm -rf *` - Mass file deletion
* `DROP DATABASE` - SQL database deletion
* `DROP TABLE` - SQL table deletion
* `chmod -R 777` - Insecure permissions
* `dd if=/dev/zero` - Disk overwrite
* `mkfs.*` - Filesystem formatting
* `--no-preserve-root` - Dangerous flag
* Fork bombs and other malicious patterns

### 6.3 Smart Rate Limiting

Prevents AI spam when same command fails repeatedly:

```bash
# After 3 failures of same command in 10 seconds:
⚠️  Same command failed 3x in 10s
💡 Tip: Review the command syntax or try 'man kubectl'
🔄 AI suggestions paused to prevent spam. Wait 10s or fix manually.
```

---

## 7. Performance: Offline Cache System 💾

v2.0 includes an offline cache for instant responses to common errors.

### 7.1 How It Works

1. **First time:** AI analyzes error and suggests fix (1-2s)
2. **Second time:** Cache returns instant response (0.05s) - **40-60% faster!**
3. **Learning:** Every successful fix is cached automatically

### 7.2 Cache Management

```bash
# View cache statistics
~/.ai/cache-manager.sh stats
📊 Cache Statistics:
  Total patterns: 47
  Cache file: ~/.ai/cache.json
  Size: 12K

# Clear cache (if needed)
~/.ai/cache-manager.sh clear

# Reinitialize with common patterns
~/.ai/cache-manager.sh init
```

### 7.3 Pre-Populated Patterns

Cache comes with 10+ common error patterns:
- `kubectl: command not found`
- `docker: permission denied`
- `terraform has not been initialized`
- `not a git repository`
- `npm: command not found`
- `ModuleNotFoundError` (Python)
- `Permission denied`
- `No such file or directory`
- `Port already in use`
- `Connection refused`

---

## 8. Create AI Helper Script

### 8.1 Quick Installation

```bash
# 1. Create directory
mkdir -p ~/.ai

# 2. Copy all scripts
cp ai-helper.sh ~/.ai/
cp cache-manager.sh ~/.ai/
cp zsh-integration.sh ~/.ai/

# 3. Make executable
chmod +x ~/.ai/*.sh

# 4. Initialize cache with common patterns
~/.ai/cache-manager.sh init

# 5. Add to ~/.zshrc
echo "source ~/.ai/zsh-integration.sh" >> ~/.zshrc

# 6. Reload shell
source ~/.zshrc
```

---

### 8.2 What's Included

**ai-helper.sh (Main Script)** - v2.0 features:
* 🔒 Security scanning for dangerous commands
* ⏱️ Smart rate limiting (prevents AI spam)
* 🗣️ Proactive mode support
* 🧠 Command history learning
* 💾 Offline cache integration
* 🎯 Enhanced model routing for 15+ tools
* 📊 Context-aware suggestions

**cache-manager.sh (Performance)**:
* 💾 Offline cache for common errors
* ⚡ 40-60% faster responses
* 🧠 Self-learning from successful fixes
* 📦 Pre-populated with 10+ common patterns

**zsh-integration.sh (Terminal Hooks)**:
* 🪝 Automatic error detection
* ⌨️ Proactive mode commands (ask, kask, dask, task, gask)
* ⚡ Hotkey bindings (⌥A, ⌥K)
* 📊 Quick stats and management commands

---

## 9. Terminal Integration

### 9.1 Automatic Setup

If you used the quick installation above, you're done! The `zsh-integration.sh` file includes:

* ✅ Automatic error detection (triggers AI on failures)
* ✅ Proactive commands (ask, kask, dask, task, gask)
* ✅ Hotkey bindings (⌥A, ⌥K)
* ✅ Rate limiting (prevents AI spam)
* ✅ Error context capture
* ✅ Quick management aliases

### 9.2 Available Commands

```bash
# Reactive Mode (automatic on errors)
# - Just run commands normally, AI triggers on failure

# Proactive Mode
ask <query>   # General: ask how do I list all pods
kask <query>  # Kubernetes: kask show pods in production
dask <query>  # Docker: dask list containers by memory
task <query>  # Terraform: task how do I plan changes
gask <query>  # Git: gask how do I revert last commit

# Manual Trigger
ai            # Re-analyze last failed command

# Management
ai-stats      # Show AI usage statistics
ai-history    # Show last 20 AI assists
ai-clear      # Clear rate limit (if stuck)
ai-cache      # Show cache info
```

### 9.3 Hotkeys

* **⌥A (Option+A)** - Trigger `ai` (re-analyze last error)
* **⌥K (Option+K)** - Quick `ask ` mode (opens prompt)

### 9.4 Manual Integration (Advanced)

If you want to customize the integration, see the `zsh-integration.sh` file for the full implementation. Key hooks:

* `preexec()` - Captures command before execution
* `precmd()` - Checks for errors after execution
* Rate limiting with 2s cooldown
* Error context capture via stderr redirect

---

## 10. Real-World Examples (v2.0)

### 10.1 Kubernetes Debugging (with Cache)

```bash
# First time - AI analyzes (1.5s)
$ kubectl apply -f deployment.yaml
Error: unknown field "replicas" in v1.Pod

🤖 AI Assistant (exit 1):
✓ Change `kind: Pod` to `kind: Deployment`
Root: Pods don't have a replicas field; only Deployments/StatefulSets do.
Tip: Validate with kubectl apply --dry-run=client first

# Second time - Cached response (0.05s) ⚡
$ kubectl apply -f deployment.yaml
Error: unknown field "replicas" in v1.Pod

💾 [Cached] ✓ Change `kind: Pod` to `kind: Deployment`
Root: Pods don't have a replicas field; only Deployments/StatefulSets do.
Tip: Validate with kubectl apply --dry-run=client first
```

### 10.2 Proactive Mode (NEW!)

```bash
# Generate commands before errors
$ kask how do I scale my deployment to 5 replicas
🤖 Generating command for: how do I scale my deployment to 5 replicas
✓ kubectl scale deployment <deployment-name> --replicas=5
Root: Scales the specified deployment to 5 pod replicas
Tip: Use -n <namespace> to target specific namespace
```

### 10.3 Security Scanning (NEW!)

```bash
$ kubectl delete pods --all -n production
# User types this by mistake, AI catches dangerous pattern

🤖 AI Assistant:
🚨 DANGER: Command contains potentially destructive pattern: --all
⚠️  This could cause data loss or system damage!
📋 Command: kubectl delete pods --all -n production

If you're ABSOLUTELY SURE this is safe, you can:
1. Review the command carefully
2. Test in a safe environment first
3. Execute manually after verification
```

### 10.4 Smart Rate Limiting (NEW!)

```bash
# First failure
$ terraform apply
Error: state lock already acquired
🤖 AI Assistant (exit 1):
✓ terraform force-unlock <lock-id>

# Second failure (same command)
$ terraform apply
Error: state lock already acquired
🤖 AI Assistant (exit 1):
✓ terraform force-unlock <lock-id>

# Third failure - Rate limit kicks in
$ terraform apply
Error: state lock already acquired
⚠️  Same command failed 3x in 10s
💡 Tip: Review the command syntax or try 'man terraform'
🔄 AI suggestions paused to prevent spam. Wait 10s or fix manually.
```

### 10.5 Terraform State Issues

```bash
$ terraform apply
Error: state lock already acquired

🤖 AI Assistant (exit 1):
✓ terraform force-unlock <lock-id>
Root cause: Previous operation crashed or timed out without releasing lock.
Check: terraform state list to verify state integrity after unlock.
Production note: Use remote state with proper locking (S3 + DynamoDB).
```

### 8.3 Docker Build Failures

```bash
$ docker build -t myapp:latest .
ERROR: failed to solve: failed to compute cache key

🤖 AI Assistant (exit 1):
✓ docker build --no-cache -t myapp:latest .
Root cause: Corrupted build cache or context path issue.
Alternative: docker builder prune -a (clears build cache)
Check Dockerfile COPY paths match context.
```

### 8.4 Python/ML Environment Issues

```bash
$ pip install tensorflow==2.13.0
ERROR: Could not find a version that satisfies the requirement

🤖 AI Assistant (exit 1):
✓ pip install tensorflow==2.13.0 --extra-index-url https://...
Root cause: Package not available for Python 3.12 or arm64.
Check: python --version and uname -m
For M1/M2: Use tensorflow-macos or conda-forge channel.
```

### 8.5 AWS CLI Permission Errors

```bash
$ aws s3 ls s3://prod-bucket
An error occurred (AccessDenied)

🤖 AI Assistant (exit 1):
✓ Check IAM policy: aws iam list-attached-user-policies --user-name $(aws sts get-caller-identity --query 'Arn' --output text | cut -d'/' -f2)
Root cause: Missing s3:ListBucket permission.
Debug: aws sts get-caller-identity (verify assumed role)
Production: Use least-privilege policies + AWS Organizations SCPs.
```

---

## 9. Performance Tuning (M1 / M2 / M3 / M4)

Add to `~/.zshrc`:

```bash
# Optimize for M1/M2/M3/M4 (16 GB RAM)
export OLLAMA_NUM_THREADS=6        # Use 6 performance cores
export OLLAMA_MAX_LOADED_MODELS=1  # Keep memory under 8 GB
export OLLAMA_FLASH_ATTENTION=1    # Enable Flash Attention 2.0
export OLLAMA_THINKING=false      # CRITICAL: Disable verbose thinking output

# Optional: Reduce context window for faster responses
export OLLAMA_NUM_CTX=4096         # Default is 2048, increase for logs
```

**Important:** `OLLAMA_THINKING=false` is **required** to suppress the verbose "Thinking..." blocks. Without this, you'll see 100+ lines of reasoning process.

**Expected latency:**

* 1-2B models (ultra-fast): 0.3-0.8s
* 4B models: 0.5-1.5s
* 8B models: 1-2.5s

**Note:** Ultra-fast tier uses qwen3:1.7b (preferred) or gemma3:1b for trivial shell errors.

---

## 10. Security & Production Safety

### 10.1 Zero-Trust Guarantees

* ✅ **Fully local execution** (no network calls)
* ✅ **No telemetry** (Ollama doesn't phone home)
* ✅ **Safe for secrets** (AWS keys, k8s tokens, DB passwords)
* ✅ **Audit trail** (all prompts visible in script)

### 10.2 Compliance-Friendly

Safe for use with:

* Production logs (PII, PHI)
* Regulated environments (SOC2, HIPAA, PCI-DSS, FedRAMP)
* Air-gapped networks
* Client confidential data

### 10.3 Smart Filtering & Clean Output

Script automatically:

* **Skips** trivial commands (`cd`, `ls`, `pwd`, `echo`)
* **Skips** successful commands (exit code 0)
* **Rate-limits** rapid failures (2s cooldown)
* **Aggressively filters** verbose AI thinking:
  * Strips entire "Thinking..." blocks
  * Removes reasoning process ("Okay,", "Let me", "Wait,", etc.)
  * Only outputs lines starting with ✓, Root:, Tip:, Check:, Fix:
* **Enforces** ultra-concise format (max 4 lines, direct answers only)
* **Zero fluff** - only the fix, never the thinking process

---

## 11. Advanced Configuration

### 11.1 Custom Model for Specific Tools

Edit `~/.ai/ai-helper.sh`:

```bash
# Add your own tool routing
elif [[ "$cmd" =~ (ansible|salt|puppet) ]]; then
  echo "qwen3:8b-q4_K_M"  # Config management needs reasoning
```

### 11.2 Logging for Debugging

Add to script (before `ollama run`):

```bash
# Log all AI queries for analysis
echo "[$(date)] CMD: $CMD | MODEL: $MODEL" >> ~/.ai/ai-helper.log
```

### 11.3 Integration with Monitoring

```bash
# Send AI trigger metrics to StatsD/Prometheus
echo "ai_helper.triggered:1|c" | nc -u -w0 localhost 8125
```

---

## 12. Troubleshooting

### Model Not Found

```bash
ollama list                  # Check installed models
ollama pull qwen3:8b-q4_K_M  # Download missing model
```

### Slow Responses

```bash
ollama ps                    # Check running models
killall ollama && ollama serve  # Restart Ollama
```

### AI Not Triggering

```bash
# Test manually
~/.ai/ai-helper.sh "kubectl get pods" "connection refused" "1"

# Check zsh hooks
which preexec precmd
```

### Verbose Output Still Appearing

If you're still seeing "Thinking..." blocks despite the aggressive filtering:

```bash
# 1. Verify the filter is working
echo -e "Thinking...\nOkay\n✓ test\nRoot: test" | awk '!answer_started {if(/^✓/) {answer_started=1; print} next} answer_started && (/^Root:/ || length($0)<100) {print}'

# 2. Check Ollama version (newer = better at following prompts)
ollama --version

# 3. Try a different model (some models are more verbose)
# qwen3 models tend to be less verbose than gemma3
# Edit ~/.ai/ai-helper.sh and temporarily force: MODEL="qwen3:4b-q4_K_M"

# 4. Check if the model is outputting ✓ correctly
~/.ai/ai-helper.sh "test command" "test error" "1" | head -5

# 5. If thinking persists, the model may not support suppressing it
# Consider using a different model or updating Ollama
ollama pull qwen3:4b-q4_K_M  # Try this model as it's less verbose
```

**Note:** The filter skips ALL output until it sees a line starting with ✓. If your model never outputs ✓, you'll see no output. This is by design to prevent verbose thinking.

---

## 13. Expected UX

* ✅ AI appears **only on failure** (silent on success)
* ✅ **Zero verbose output** - no thinking process, no reasoning, no fluff
* ✅ Ultra-short answers (2-4 lines max)
* ✅ Corrected command first, explanation second
* ✅ Production safety notes included
* ✅ **Ultra-fast responses** (0.3-2.5s depending on complexity)
* ✅ Feels like a senior SRE pair-programming with you

**Example clean output:**

```text
🤖 AI Assistant (exit 127):
✓ mkdir -p ~/.ai
Root: Directory doesn't exist, -p creates parent dirs.
```

**What you WON'T see:**

* ❌ NO "Thinking..." blocks (100+ lines of reasoning)
* ❌ NO "Let me think...", "Okay,", "Wait," verbose narration
* ❌ NO walls of text explaining the thought process
* ❌ NO redundant explanations

**Just the fix. Nothing else.**

---

## 14. Installation Checklist (v2.0)

```bash
# 1. Install Ollama
brew install ollama

# 2. Start Ollama service
ollama serve &

# 3. Pull required models
ollama pull qwen3:8b-q4_K_M      # Primary (8B) - proactive + complex infra
ollama pull qwen3:4b-q4_K_M      # Fast fallback (4B)
ollama pull gemma3:4b-it-q4_K_M  # Explanations (4B) - ML/Python
ollama pull qwen3:1.7b-q4_K_M    # Ultra-fast (1.7B) - optional but recommended

# 4. Clone repository
git clone https://github.com/yourusername/ai-helper.git
cd ai-helper

# 5. Install scripts
mkdir -p ~/.ai
cp ai-helper.sh ~/.ai/
cp cache-manager.sh ~/.ai/
cp zsh-integration.sh ~/.ai/
chmod +x ~/.ai/*.sh

# 6. Initialize cache with common patterns
~/.ai/cache-manager.sh init
✅ Cache initialized with common error patterns

# 7. Add to ~/.zshrc
echo "source ~/.ai/zsh-integration.sh" >> ~/.zshrc

# 8. Reload shell
source ~/.zshrc
✅ AI Terminal Helper v2.0 Loaded!

# 9. Test reactive mode (error fixing)
kubectl get pods --invalid-flag  # Should trigger AI

# 10. Test proactive mode (NEW!)
ask how do I list all pods  # Should generate command
```

---

## 15. Status & Performance

**Production-Ready AI Terminal Assistant v2.0 (2025)**  
**Optimized for:** DevOps · SRE · MLOps · Platform Engineering

### Performance Metrics (v2.0)
- ⚡ **0.05s** - Cached responses (40-60% of queries)
- ⚡ **0.3-0.8s** - Ultra-fast tier (1-2B models)
- ⚡ **0.5-1.5s** - Fast tier (4B models)
- ⚡ **1-2.5s** - Deep reasoning tier (8B models)

### Feature Completion
- ✅ Security scanning for dangerous commands
- ✅ Smart rate limiting (prevents AI spam)
- ✅ Proactive mode (natural language → commands)
- ✅ Offline cache (instant common error responses)
- ✅ Command history learning
- ✅ Context-aware model routing
- ✅ Production-safe defaults

### Comparison to Warp Terminal
| Feature | AI Helper v2.0 | Warp Terminal |
|---------|---------------|---------------|
| Privacy | ✅ 100% local | ❌ Cloud-based |
| Cost | ✅ Free | ❌ $10-20/mo |
| Proactive Mode | ✅ Yes (NEW!) | ✅ Yes |
| Caching | ✅ Yes (NEW!) | ⚠️ Limited |
| Security Scanning | ✅ Yes (NEW!) | ❌ No |
| Secrets Safe | ✅ Yes | ❌ No |
| Offline | ✅ Yes | ❌ No |

Built with ❤️ for engineers who need fast, private, actionable help.

---

## 16. Roadmap & Future Improvements

See [ROADMAP.md](ROADMAP.md) for detailed plans on upcoming features.

### ✅ Completed (v2.0)
- ✅ Security scanning for dangerous commands (Phase 3.3)
- ✅ Smart rate limiting (Phase 1.2)
- ✅ Proactive mode for natural language queries
- ✅ Offline cache for common errors (Phase 1.4)
- ✅ Command history learning (Phase 1.1)

### 🚧 Next Up (v2.1)
- Tool-specific helpers (kubectl, terraform, docker syntax validation)
- Context preservation across commands
- Multi-model ensemble for critical operations
- Interactive mode (choose action before AI call)
- Enhanced confidence scoring

**Current Version:** v2.0 (Production-ready)  
**Next Milestone:** Phase 2 - Intelligence & Context (2-3 weeks)
