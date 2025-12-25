#!/usr/bin/env bash

# AI Terminal Helper for DevOps/SRE/MLOps
# Provides context-aware AI assistance for failed commands
# Safe for production environments (fully local)

CMD="$1"
ERR="$2"
EXIT_CODE="$3"

# Capture additional context for better assistance
CWD="${PWD}"

# Skip AI help for sensitive/fast-failing commands
SKIP_PATTERNS="(^cd |^ls |^pwd|^echo |^cat |^history)"
if [[ "$CMD" =~ $SKIP_PATTERNS ]]; then
  exit 0
fi

# Enhanced model routing for DevOps/SRE/MLOps workflows
select_model() {
  local cmd="$1"
  
  # Infrastructure & Orchestration (8B - needs deep reasoning)
  if [[ "$cmd" =~ (kubectl|helm|terraform|terragrunt|aws|gcloud|azure) ]]; then
    echo "qwen3:8b-q4_K_M"
  
  # Container & CI/CD (8B - complex configs)
  elif [[ "$cmd" =~ (docker|podman|buildah|gitlab-ci|jenkins|circleci) ]]; then
    echo "qwen3:8b-q4_K_M"
  
  # Monitoring & Observability (8B - log analysis)
  elif [[ "$cmd" =~ (prometheus|grafana|datadog|kubectl logs|stern) ]]; then
    echo "qwen3:8b-q4_K_M"

  # Config management (8B - needs deep reasoning)
  elif [[ "$cmd" =~ (ansible|salt|puppet) ]]; then
    echo "qwen3:8b-q4_K_M"
  
  # ML/Data Engineering (4B instruction-tuned - stack traces)
  elif [[ "$cmd" =~ (python|pip|conda|poetry|jupyter|mlflow|kubeflow|ray|spark) ]]; then
    echo "gemma3:4b-it-q4_K_M"
  
  # Ultra-fast for trivial shell errors (1B - instant response)
  elif [[ "$cmd" =~ ^(cp|mv|rm|mkdir|touch|grep|find|awk|sed) ]]; then
    # Check if we have qwen3 1.7B (best quality), fallback to gemma3 1B
    if ollama list 2>/dev/null | grep -q "qwen3:1.7b-q4_K_M"; then
      echo "qwen3:1.7b-q4_K_M"
    else
      echo "gemma3:1b-it-q4_K_M"
    fi
  
  # Simple shell commands (4B - fast response)
  else
    echo "qwen3:4b-q4_K_M"
  fi
}

MODEL=$(select_model "$CMD")

# Build concise prompt (no verbose thinking, direct answer only)
PROMPT=$(cat <<EOF
You are a senior DevOps/SRE. Fix this failed command.

CRITICAL RULES:
1. DO NOT output "Thinking..." or any reasoning process
2. DO NOT start with "Okay," "Let me," "Wait," or any explanation
3. START IMMEDIATELY with ✓ followed by the corrected command
4. NO thinking blocks, NO verbose reasoning, NO process explanation

Command: $CMD
Error: $ERR
Exit: ${EXIT_CODE:-?}
Dir: $CWD

REQUIRED OUTPUT FORMAT (start immediately, no preamble):
✓ [corrected command]
Root: [1 sentence why it failed]
Tip: [optional best practice]

Your first line MUST be: ✓ [command]
EOF
)

# Run AI with error handling
if ! command -v ollama &> /dev/null; then
  echo "⚠️  Ollama not found. Install: https://ollama.ai"
  exit 1
fi

# Check if model is available
if ! ollama list | grep -q "$MODEL"; then
  echo "⚠️  Model $MODEL not found. Run: ollama pull $MODEL"
  exit 1
fi

# Execute AI query with ULTRA-AGGRESSIVE filtering (zero thinking output)
# Use OLLAMA_THINKING=false environment variable to suppress thinking at source
# Then filter output to ensure clean results
OLLAMA_THINKING=false ollama run "$MODEL" "$PROMPT" 2>/dev/null | awk '
  BEGIN { 
    answer_started = 0
    line_count = 0
  }
  
  # Skip ALL lines until we see an answer marker (✓ or Root:/Tip:/etc)
  !answer_started {
    # Answer markers that signal the start of useful output
    if (/^✓/ || /^Root:/ || /^Tip:/ || /^Check:/ || /^Fix:/ || /^Note:/ || /^Error:/) {
      answer_started = 1
      print $0
      line_count = 1
      next
    }
    # Skip everything else (thinking, reasoning, etc.)
    next
  }
  
  # After answer started, only output answer format lines
  answer_started {
    # Answer format markers
    if (/^Root:/ || /^Tip:/ || /^Check:/ || /^Fix:/ || /^Note:/ || /^Error:/) {
      print $0
      line_count++
      next
    }
    
    # Skip verbose reasoning patterns (even after answer started)
    if (/^[Tt]hinking/ || /done [Tt]hinking/ || /^[Oo]kay/ || /^[Ll]et me/ || \
        /^[Ww]ait/ || /^[Aa]lternatively/ || /^[Hh]owever/ || /^[Bb]ut / || \
        /^[Ss]o / || /^[Ii]n / || /^[Ww]hen / || /^[Tt]he / || /^[Tt]his / || \
        /^[Tt]hat / || /^[Mm]aybe/ || /^[Ff]irst/ || /^[Aa]nother/ || /^[Oo]r / || \
        /^[Aa]lso/ || /^[Ii]f / || /^[Ss]hould/ || /^[Mm]ight/ || /^[Pp]erhaps/) {
      next
    }
    
    # Allow short continuation lines (max 4 lines total)
    if (line_count < 4 && length($0) < 100 && !/^[A-Z][a-z]+ [a-z]/) {
      print $0
      line_count++
      next
    }
    
    # Stop after reasonable output
    if (line_count >= 4) {
      exit 0
    }
    
    next
  }
  
  END { 
    if (!answer_started) {
      exit 1
    }
  }
' || {
  echo "⚠️  AI helper failed. Check Ollama status: ollama ps"
  exit 1
}