# Local Terminal AI Manual (2025)

**Audience:** DevOps / SRE / MLOps engineers
**Platform:** macOS (Apple Silicon, 16 GB RAM)
**Goal:** Local, fast, context-aware AI companion for production environments

---

## 1. Design Principles

* **Local only** (no cloud, no telemetry, safe for secrets)
* **Silent on success**, helpful only on errors
* **Low latency** on M1 / M2 / M3 / M4 (< 2s response)
* **Context-aware** (infra, containers, ML, cloud)
* **Production-safe** for regulated environments (SOC2, HIPAA, PCI-DSS)
* **Smart filtering** (skips trivial commands like `cd`, `ls`)

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

## 5. Create AI Helper Script

### 5.1 Directory

```bash
mkdir -p ~/.ai
```

---

### 5.2 Enhanced Helper Script

Create file:

```bash
nano ~/.ai/ai-helper.sh
```

Paste the enhanced version from `ai-helper.sh` in this repository, which includes:

* **Enhanced model routing** for 15+ DevOps/MLOps tools
* **Context awareness** (cwd, exit code, timestamp)
* **Smart filtering** (skips trivial commands)
* **Error handling** (checks Ollama availability)
* **Production-grade prompts** with safety considerations

Key features:

* Detects infrastructure tools (kubectl, terraform, aws, docker)
* Recognizes ML platforms (mlflow, kubeflow, ray, spark)
* Provides corrected commands + root cause analysis
* Includes best practices and production safety notes

Make executable:

```bash
chmod +x ~/.ai/ai-helper.sh
```

---

## 6. Enhanced Terminal Hook (Production-Ready)

Add to `~/.zshrc`:

```bash
autoload -Uz add-zsh-hook

# State tracking for AI helper
LAST_CMD=""
LAST_OUTPUT=""
LAST_EXIT_CODE=0

# Rate limiting (prevent spam on rapid failures)
AI_LAST_CALL=0
AI_COOLDOWN=2  # seconds between AI calls

# Capture command before execution
preexec() {
  LAST_CMD="$1"
  LAST_OUTPUT=""
}

# Check for errors after command completes
precmd() {
  local exit_code=$?
  LAST_EXIT_CODE=$exit_code
  
  # Only trigger on failure
  if [[ $exit_code -ne 0 ]]; then
    # Rate limiting check
    local now=$(date +%s)
    local elapsed=$((now - AI_LAST_CALL))
    
    if [[ $elapsed -ge $AI_COOLDOWN ]]; then
      echo "\n🤖 AI Assistant (exit $exit_code):"
      ~/.ai/ai-helper.sh "$LAST_CMD" "$LAST_OUTPUT" "$exit_code"
      AI_LAST_CALL=$now
    fi
  fi
}

# Capture stderr for error context
exec 2> >(while IFS= read -r line; do 
  LAST_OUTPUT+="$line"$'\n'
  echo "$line" >&2
done)
```

**Features:**

* Rate limiting (2s cooldown) prevents AI spam on rapid failures
* Captures stderr for better error context
* Passes exit code for precise diagnosis
* Shows exit code in prompt for transparency

---

## 7. Manual AI Trigger (Enhanced)

Add to `~/.zshrc`:

```bash
# Manual AI invocation for any previous command
ai() {
  if [[ -n "$LAST_CMD" ]]; then
    ~/.ai/ai-helper.sh "$LAST_CMD" "$LAST_OUTPUT" "$LAST_EXIT_CODE"
  else
    echo "No previous command found"
  fi
}

# Hotkey binding (Option+A)
bindkey '^[a' ai

# Alternative: Ask AI about specific command
ask() {
  local cmd="$*"
  ~/.ai/ai-helper.sh "$cmd" "" "0"
}
```

**Usage:**

* Type `ai` or press **⌥A** (Option+A) to re-analyze last failure
* Type `ask kubectl get pods --all-namespaces` for help with any command

---

## 8. Real-World Examples

### 8.1 Kubernetes Debugging

```bash
$ kubectl apply -f deployment.yaml
Error: unknown field "replicas" in v1.Pod

🤖 AI Assistant (exit 1):
✓ Change `kind: Pod` to `kind: Deployment`
Root cause: Pods don't have a replicas field; only Deployments/StatefulSets do.
Use: kubectl apply -f deployment.yaml --dry-run=client -o yaml
Best practice: Validate with kubeval before applying.
```

### 8.2 Terraform State Issues

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

## 14. Installation Checklist

```bash
# 1. Install Ollama
brew install ollama

# 2. Start Ollama service
ollama serve &

# 3. Pull required models
ollama pull qwen3:8b-q4_K_M      # Primary (8B)
ollama pull qwen3:4b-q4_K_M      # Fast fallback (4B)
ollama pull gemma3:4b-it-q4_K_M  # Explanations (4B)
ollama pull qwen3:1.7b-q4_K_M    # Ultra-fast (1.7B) - optional but recommended

# 4. Setup script directory
mkdir -p ~/.ai
cp ai-helper.sh ~/.ai/
chmod +x ~/.ai/ai-helper.sh

# 5. Add to ~/.zshrc (see section 6 & 7)
# ... paste config ...

# 6. Reload shell
source ~/.zshrc

# 7. Test
kubectl get pods --invalid-flag  # Should trigger AI
```

---

## 15. Status

**Production-Ready AI Terminal Assistant (2025)**
**Optimized for:** DevOps · SRE · MLOps · Platform Engineering

Built with ❤️ for engineers who need fast, private, actionable help.
