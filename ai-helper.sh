#!/usr/bin/env bash

# AI Terminal Helper for DevOps/SRE/MLOps v2.0
# Provides context-aware AI assistance for failed commands
# Safe for production environments (fully local)
#
# Features:
# - Security scanning for dangerous commands
# - Smart rate limiting to prevent AI spam
# - Proactive mode for natural language queries
# - Command history learning
# - Production-safe defaults

# ============================================================================
# COLOR CONFIGURATION
# ============================================================================

# ANSI color codes for beautiful output
COLOR_RESET='\033[0m'
COLOR_BOLD='\033[1m'
COLOR_DIM='\033[2m'

# Foreground colors
COLOR_RED='\033[0;31m'
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW='\033[0;33m'
COLOR_BLUE='\033[0;34m'
COLOR_MAGENTA='\033[0;35m'
COLOR_CYAN='\033[0;36m'
COLOR_WHITE='\033[0;37m'

# Bold colors
COLOR_RED_BOLD='\033[1;31m'
COLOR_GREEN_BOLD='\033[1;32m'
COLOR_YELLOW_BOLD='\033[1;33m'
COLOR_BLUE_BOLD='\033[1;34m'
COLOR_MAGENTA_BOLD='\033[1;35m'
COLOR_CYAN_BOLD='\033[1;36m'

# Background colors
COLOR_BG_RED='\033[41m'
COLOR_BG_YELLOW='\033[43m'
COLOR_BG_GREEN='\033[42m'

# Disable colors if NO_COLOR is set or not in terminal
if [[ -n "$NO_COLOR" ]] || [[ ! -t 1 ]]; then
  COLOR_RESET=''
  COLOR_BOLD=''
  COLOR_DIM=''
  COLOR_RED=''
  COLOR_GREEN=''
  COLOR_YELLOW=''
  COLOR_BLUE=''
  COLOR_MAGENTA=''
  COLOR_CYAN=''
  COLOR_WHITE=''
  COLOR_RED_BOLD=''
  COLOR_GREEN_BOLD=''
  COLOR_YELLOW_BOLD=''
  COLOR_BLUE_BOLD=''
  COLOR_MAGENTA_BOLD=''
  COLOR_CYAN_BOLD=''
  COLOR_BG_RED=''
  COLOR_BG_YELLOW=''
  COLOR_BG_GREEN=''
fi

# ============================================================================
# CONFIGURATION
# ============================================================================

# Directory for AI helper data
AI_DIR="${HOME}/.ai"
RATE_LIMIT_FILE="${AI_DIR}/rate_limit.log"
HISTORY_FILE="${AI_DIR}/history.log"
CACHE_FILE="${AI_DIR}/cache.json"

# Rate limiting configuration (prevent AI spam)
RATE_LIMIT_WINDOW=10  # seconds
RATE_LIMIT_MAX_SAME_CMD=3  # max same command failures in window

# Create AI directory if it doesn't exist
mkdir -p "$AI_DIR"

# ============================================================================
# INPUT HANDLING
# ============================================================================

CMD="$1"
ERR="$2"
EXIT_CODE="$3"
CONTEXT="${4:-}"  # Optional context for proactive mode

# Capture additional context for better assistance
CWD="${PWD}"

# Detect if this is a proactive query (not an error)
PROACTIVE_MODE=false
if [[ "$CMD" == "PROACTIVE_QUERY" ]]; then
  PROACTIVE_MODE=true
  CMD="$ERR"  # In proactive mode, query is in ERR field
  ERR=""
  EXIT_CODE="0"
fi

# ============================================================================
# SECURITY SCANNING
# ============================================================================

# Dangerous command patterns that require extra caution
# These patterns are checked in AI suggestions to prevent catastrophic mistakes
DANGEROUS_PATTERNS=(
  "rm -rf /"
  "rm -rf \*"
  "rm -rf ~"
  "rm -rf \$HOME"
  "> /dev/sda"
  "dd if=/dev/zero"
  "mkfs\."
  "DROP DATABASE"
  "DROP TABLE"
  "TRUNCATE"
  "chmod -R 777"
  "chown -R"
  ":(){ :\|:& };:"  # Fork bomb
  "--no-preserve-root"
  "mv .* /dev/null"
)

# Check if command contains dangerous patterns
check_dangerous_command() {
  local cmd="$1"
  
  for pattern in "${DANGEROUS_PATTERNS[@]}"; do
    if [[ "$cmd" =~ $pattern ]]; then
      echo -e "${COLOR_RED_BOLD}${COLOR_BG_RED} DANGER ${COLOR_RESET} ${COLOR_RED_BOLD}Command contains potentially destructive pattern:${COLOR_RESET} ${COLOR_YELLOW}$pattern${COLOR_RESET}"
      echo -e "${COLOR_RED}⚠️  This could cause data loss or system damage!${COLOR_RESET}"
      echo -e "${COLOR_DIM}📋 Command: ${COLOR_WHITE}$cmd${COLOR_RESET}"
      echo ""
      echo -e "${COLOR_YELLOW}If you're ABSOLUTELY SURE this is safe, you can:${COLOR_RESET}"
      echo -e "  ${COLOR_CYAN}1.${COLOR_RESET} Review the command carefully"
      echo -e "  ${COLOR_CYAN}2.${COLOR_RESET} Test in a safe environment first"
      echo -e "  ${COLOR_CYAN}3.${COLOR_RESET} Execute manually after verification"
      return 1
    fi
  done
  
  return 0
}

# ============================================================================
# SMART RATE LIMITING
# ============================================================================

# Check if we're spamming AI with the same failed command
check_rate_limit() {
  local cmd="$1"
  local now
  now=$(date +%s)
  
  # Create rate limit file if it doesn't exist
  touch "$RATE_LIMIT_FILE"
  
  # Count recent failures of the same command
  local recent_count
  recent_count=$(grep -c "^${cmd}|" "$RATE_LIMIT_FILE" 2>/dev/null || echo "0")
  
  # Clean old entries (older than RATE_LIMIT_WINDOW seconds)
  if [[ -f "$RATE_LIMIT_FILE" ]]; then
    while IFS='|' read -r logged_cmd timestamp; do
      local age=$((now - timestamp))
      if [[ $age -lt $RATE_LIMIT_WINDOW ]]; then
        echo "${logged_cmd}|${timestamp}"
      fi
    done < "$RATE_LIMIT_FILE" > "${RATE_LIMIT_FILE}.tmp"
    mv "${RATE_LIMIT_FILE}.tmp" "$RATE_LIMIT_FILE"
  fi
  
  # Re-count after cleanup
  recent_count=$(grep -c "^${cmd}|" "$RATE_LIMIT_FILE" 2>/dev/null || echo "0")
  # Remove any whitespace/newlines
  recent_count=$(echo "$recent_count" | tr -d '[:space:]')
  
  # If same command failed too many times, skip AI
  if [[ "$recent_count" -ge "$RATE_LIMIT_MAX_SAME_CMD" ]]; then
    echo -e "${COLOR_YELLOW}⚠️  Same command failed ${COLOR_BOLD}${recent_count}x${COLOR_RESET}${COLOR_YELLOW} in ${RATE_LIMIT_WINDOW}s${COLOR_RESET}"
    echo -e "${COLOR_CYAN}💡 Tip:${COLOR_RESET} Review the command syntax or try ${COLOR_GREEN}man $(echo "$cmd" | awk '{print $1}')${COLOR_RESET}"
    echo -e "${COLOR_BLUE}🔄 AI suggestions paused${COLOR_RESET} to prevent spam. Wait ${RATE_LIMIT_WINDOW}s or fix manually."
    return 1
  fi
  
  # Log this attempt
  echo "${cmd}|${now}" >> "$RATE_LIMIT_FILE"
  return 0
}

# ============================================================================
# SKIP PATTERNS
# ============================================================================

# Skip AI help for trivial/sensitive commands (only in error mode)
if [[ "$PROACTIVE_MODE" == false ]]; then
  SKIP_PATTERNS="(^cd |^ls |^pwd|^echo |^cat |^history)"
  if [[ "$CMD" =~ $SKIP_PATTERNS ]]; then
    exit 0
  fi
  
  # Check rate limit for repeated failures
  if ! check_rate_limit "$CMD"; then
    exit 0
  fi
fi

# ============================================================================
# MODEL SELECTION
# ============================================================================

# Enhanced model routing for DevOps/SRE/MLOps workflows
select_model() {
  local cmd="$1"
  
  # Proactive mode: Use 8B for better quality suggestions
  if [[ "$PROACTIVE_MODE" == true ]]; then
    echo "qwen3:8b-q4_K_M"
    return
  fi
  
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

# ============================================================================
# PROMPT GENERATION
# ============================================================================

# Build prompt based on mode (proactive vs error-fixing)
if [[ "$PROACTIVE_MODE" == true ]]; then
  # Proactive mode: Natural language to command translation
  PROMPT=$(cat <<EOF
You are a senior DevOps/SRE. Convert this natural language query to a command.

CRITICAL RULES:
1. DO NOT output "Thinking..." or any reasoning process
2. START IMMEDIATELY with ✓ followed by the command
3. NO verbose reasoning, NO process explanation

Query: $CMD
Context: $CONTEXT
Dir: $CWD

REQUIRED OUTPUT FORMAT (start immediately):
✓ [command]
Root: [1 sentence what this does]
Tip: [optional safety note or best practice]

Your first line MUST be: ✓ [command]
EOF
)
else
  # Error-fixing mode: Fix failed command
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
fi

# ============================================================================
# CACHE LOOKUP (Fast path - skip AI if cached)
# ============================================================================

# Try cache first for instant responses (only in error mode)
if [[ "$PROACTIVE_MODE" == false ]]; then
  # Source cache manager functions
  if [[ -f "${AI_DIR}/cache-manager.sh" ]]; then
    source "${AI_DIR}/cache-manager.sh"
    
    # Try to find cached response
    CACHED_RESPONSE=$(lookup_cache "$CMD" "$ERR" 2>/dev/null)
    if [[ -n "$CACHED_RESPONSE" ]]; then
      echo "$CACHED_RESPONSE"
      exit 0
    fi
  fi
fi

# ============================================================================
# AI EXECUTION
# ============================================================================

# Check Ollama availability
if ! command -v ollama &> /dev/null; then
  echo -e "${COLOR_RED}⚠️  Ollama not found.${COLOR_RESET} Install: ${COLOR_BLUE}https://ollama.ai${COLOR_RESET}"
  exit 1
fi

# Check if model is available
if ! ollama list | grep -q "$MODEL"; then
  echo -e "${COLOR_YELLOW}⚠️  Model ${COLOR_BOLD}$MODEL${COLOR_RESET}${COLOR_YELLOW} not found.${COLOR_RESET}"
  echo -e "Run: ${COLOR_GREEN}ollama pull $MODEL${COLOR_RESET}"
  exit 1
fi

# Execute AI query and capture output
AI_RESPONSE=$(OLLAMA_THINKING=false ollama run "$MODEL" "$PROMPT" 2>/dev/null | awk '
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
')

# Check if AI response was successful
if [[ -z "$AI_RESPONSE" ]]; then
  echo -e "${COLOR_RED}⚠️  AI helper failed.${COLOR_RESET} Check Ollama status: ${COLOR_GREEN}ollama ps${COLOR_RESET}"
  exit 1
fi

# Extract suggested command from AI response (line starting with ✓)
SUGGESTED_CMD=$(echo "$AI_RESPONSE" | grep "^✓" | sed 's/^✓ //')

# ============================================================================
# SECURITY CHECK ON AI SUGGESTIONS
# ============================================================================

# Check if AI suggested a dangerous command
if [[ -n "$SUGGESTED_CMD" ]]; then
  if ! check_dangerous_command "$SUGGESTED_CMD"; then
    echo ""
    echo "🤖 AI suggested:"
    echo "$AI_RESPONSE"
    exit 1
  fi
fi

# ============================================================================
# OUTPUT & HISTORY
# ============================================================================

# Format and colorize AI response
format_ai_response() {
  local response="$1"
  
  echo "$response" | while IFS= read -r line; do
    case "$line" in
      ✓*)
        # Green for corrected commands
        echo -e "${COLOR_GREEN_BOLD}$line${COLOR_RESET}"
        ;;
      Root:*)
        # Cyan for root cause
        echo -e "${COLOR_CYAN}$line${COLOR_RESET}"
        ;;
      Tip:*|Check:*|Fix:*|Note:*)
        # Yellow for tips and notes
        echo -e "${COLOR_YELLOW}$line${COLOR_RESET}"
        ;;
      Error:*)
        # Red for errors
        echo -e "${COLOR_RED}$line${COLOR_RESET}"
        ;;
      💾*)
        # Magenta for cached responses
        echo -e "${COLOR_MAGENTA_BOLD}$line${COLOR_RESET}"
        ;;
      *)
        # Default color
        echo -e "${COLOR_WHITE}$line${COLOR_RESET}"
        ;;
    esac
  done
}

# Display formatted AI response
format_ai_response "$AI_RESPONSE"

# Log successful fix to history and cache (for future learning)
if [[ "$PROACTIVE_MODE" == false ]] && [[ -n "$SUGGESTED_CMD" ]]; then
  # Save to history file
  echo "$(date +%s)|$CMD|$SUGGESTED_CMD|$EXIT_CODE" >> "$HISTORY_FILE"
  
  # Save to cache for instant future responses
  if [[ -f "${AI_DIR}/cache-manager.sh" ]]; then
    source "${AI_DIR}/cache-manager.sh"
    save_to_cache "$CMD" "$ERR" "$AI_RESPONSE" 2>/dev/null
  fi
fi