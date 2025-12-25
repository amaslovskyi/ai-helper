# AI Terminal Helper v2.0 - Quick Start Guide

## 🚀 5-Minute Setup

### Step 1: Install Ollama & Models (2 min)

```bash
# Install Ollama
brew install ollama

# Start Ollama service (run in background)
ollama serve &

# Pull models (this takes time, but only once)
ollama pull qwen3:8b-q4_K_M      # Primary - 4.5GB
ollama pull qwen3:4b-q4_K_M      # Fast - 2.5GB
ollama pull gemma3:4b-it-q4_K_M  # Explanations - 2.5GB
ollama pull qwen3:1.7b-q4_K_M    # Ultra-fast (optional) - 1GB
```

### Step 2: Install AI Helper (1 min)

```bash
# Create directory
mkdir -p ~/.ai

# Navigate to repo directory
cd /Users/amaslovs/Ai/ai-helper

# Copy scripts
cp ai-helper.sh cache-manager.sh zsh-integration.sh ~/.ai/
chmod +x ~/.ai/*.sh

# Initialize cache with common patterns
~/.ai/cache-manager.sh init
```

### Step 3: Add to Shell (1 min)

```bash
# Add integration to ~/.zshrc
echo "source ~/.ai/zsh-integration.sh" >> ~/.zshrc

# Reload shell
source ~/.zshrc
```

### Step 4: Test It! (1 min)

```bash
# Test reactive mode (error fixing)
kubectl get pods --invalid-flag
# Should show: 🤖 AI Assistant (exit 1):
#              ✓ kubectl get pods

# Test proactive mode (NEW!)
ask how do I list all pods
# Should show: 🤖 Generating command for: how do I list all pods
#              ✓ kubectl get pods

# Test kubernetes shortcut
kask show pods in default namespace
# Should show: ✓ kubectl get pods -n default

# Test cache (run same error twice)
docker ps --invalid
docker ps --invalid
# Second time should show: 💾 [Cached] ✓ docker ps
```

---

## 🎯 What You Get

### 1. Reactive Mode (Automatic Error Fixing)
Just run commands normally. AI triggers automatically on failures.

```bash
$ terraform apply
Error: terraform has not been initialized

🤖 AI Assistant (exit 1):
✓ terraform init
Root: Terraform working directory not initialized
Tip: Always run 'terraform init' in new directories
```

### 2. Proactive Mode (NEW! 🚀)
Generate commands from natural language BEFORE errors.

```bash
$ ask how do I list all docker containers sorted by memory
✓ docker stats --no-stream --format "table {{.Name}}\t{{.MemUsage}}" | sort -k2 -h
Root: Shows containers sorted by memory usage
Tip: Use 'docker stats' for live monitoring
```

### 3. Tool-Specific Shortcuts (NEW! 🚀)

```bash
kask <query>  # Kubernetes: kask show failing pods
dask <query>  # Docker: dask list containers using port 80
task <query>  # Terraform: task how do I destroy resources
gask <query>  # Git: gask how do I undo last commit
```

### 4. Security Scanning (NEW! 🔒)
Prevents catastrophic mistakes.

```bash
$ kubectl delete pods --all -n production
🚨 DANGER: Command contains potentially destructive pattern: --all
⚠️  This could cause data loss or system damage!
```

### 5. Smart Rate Limiting (NEW! ⏱️)
Prevents AI spam on repeated failures.

```bash
# After 3 failures of same command:
⚠️  Same command failed 3x in 10s
🔄 AI suggestions paused to prevent spam. Wait 10s or fix manually.
```

### 6. Offline Cache (NEW! 💾)
Instant responses for common errors.

```bash
# First time: AI call (1.5s)
$ kubectl version
Error: command not found
✓ brew install kubectl

# Second time: Cached (0.05s) ⚡
$ kubectl version
Error: command not found
💾 [Cached] ✓ brew install kubectl
```

---

## ⌨️ Hotkeys

- **⌥A (Option+A)** - Re-analyze last failed command
- **⌥K (Option+K)** - Quick ask mode

---

## 📊 Management Commands

```bash
ai-stats    # Show total AI assists
ai-history  # Show last 20 AI suggestions
ai-clear    # Clear rate limit
ai-cache    # Show cache statistics
ai          # Manually trigger AI on last command
```

---

## 🔍 Example Workflows

### DevOps: Kubernetes Debugging

```bash
# Reactive: Fix errors automatically
$ kubectl apply -f broken.yaml
Error: unknown field "replicas" in v1.Pod
🤖 ✓ Change `kind: Pod` to `kind: Deployment`

# Proactive: Generate commands
$ kask show pods with high CPU usage
✓ kubectl top pods --sort-by=cpu
```

### SRE: Docker Troubleshooting

```bash
# Reactive
$ docker run -p 8080:80 nginx
Error: port 8080 already in use
🤖 ✓ lsof -i :8080  # Find process using port

# Proactive
$ dask find containers using more than 1GB RAM
✓ docker stats --format "{{.Name}}: {{.MemUsage}}" | awk '$2 > 1'
```

### MLOps: Python Environment

```bash
# Reactive
$ pip install tensorflow
Error: Could not find a version that satisfies
🤖 ✓ pip install tensorflow-macos  # For M1/M2

# Proactive
$ ask how do I create a virtual environment
✓ python -m venv venv && source venv/bin/activate
```

---

## 🆚 Comparison: AI Helper v2.0 vs Warp Terminal

| Feature | AI Helper v2.0 | Warp Terminal |
|---------|---------------|---------------|
| **Privacy** | ✅ 100% local | ❌ Cloud-based |
| **Cost** | ✅ Free | ❌ $10-20/mo |
| **Proactive Mode** | ✅ Yes | ✅ Yes |
| **Security Scanning** | ✅ Yes | ❌ No |
| **Offline Cache** | ✅ Yes | ⚠️ Limited |
| **Rate Limiting** | ✅ Yes | ❌ No |
| **Secrets Safe** | ✅ Yes | ❌ No |
| **Air-gapped** | ✅ Yes | ❌ No |
| **SOC2/HIPAA Safe** | ✅ Yes | ❌ No |
| **Speed (cached)** | ✅ 0.05s | ⚠️ Network latency |
| **Speed (AI call)** | ✅ 0.3-2.5s | ⚠️ 1-5s+ |

**Verdict:** AI Helper v2.0 is better for security, privacy, cost, and speed. Warp is better for UI/UX polish.

---

## 🔧 Troubleshooting

### AI Not Triggering

```bash
# Check if scripts are loaded
which ai
# Should show: ai () { ... }

# Manually trigger
~/.ai/ai-helper.sh "kubectl get pods" "connection refused" "1"
```

### Model Not Found

```bash
# Check installed models
ollama list

# Pull missing model
ollama pull qwen3:8b-q4_K_M
```

### Slow Responses

```bash
# Check Ollama status
ollama ps

# Restart Ollama
killall ollama && ollama serve &
```

### Cache Issues

```bash
# View cache stats
~/.ai/cache-manager.sh stats

# Reinitialize cache
~/.ai/cache-manager.sh clear
~/.ai/cache-manager.sh init
```

---

## 📚 Learn More

- **Full Documentation:** [README.md](README.md)
- **Roadmap:** [ROADMAP.md](ROADMAP.md)
- **Feature Details:** See README.md sections 5-7

---

## 🎉 You're Ready!

Your terminal now has a senior SRE sitting next to you, ready to help 24/7.

**Next Steps:**
1. Try breaking some commands to see reactive mode
2. Use `ask` to explore proactive mode
3. Check `ai-stats` after a few days to see impact

**Happy debugging!** 🚀

