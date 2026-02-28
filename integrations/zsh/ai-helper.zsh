# ============================================================================
# AI Terminal Helper - Minimal ZSH Integration (Go Binary)
# Add this to your ~/.zshrc: source ~/.ai/ai-helper.zsh
# ============================================================================

# Ensure ~/.ai is in PATH (for ai-helper binary)
export PATH="$HOME/.ai:$PATH"

# Load ZSH hooks
autoload -Uz add-zsh-hook

# State tracking
LAST_CMD=""
LAST_OUTPUT=""
LAST_EXIT_CODE=0

# Rate limiting
AI_LAST_CALL=0
AI_COOLDOWN=2

# Session-level disable flag (set when user picks "Disable AI for session")
# Persists until terminal is restarted
AI_SESSION_DISABLED=0

# Capture command before execution
preexec() {
  LAST_CMD="$1"
  LAST_OUTPUT=""
}

# Check for errors after command completes
precmd() {
  local exit_code=$?
  LAST_EXIT_CODE=$exit_code

  # Skip AI for signal-terminated commands (Ctrl+C = 130, SIGTERM = 143, SIGKILL = 137)
  # These are user-initiated interruptions, not actual command failures
  if [[ $exit_code -eq 130 ]] || [[ $exit_code -eq 143 ]] || [[ $exit_code -eq 137 ]]; then
    return 0
  fi

  # Skip if AI was disabled for this session (user chose option 4)
  if [[ $AI_SESSION_DISABLED -eq 1 ]]; then
    return 0
  fi

  # Only trigger AI on actual failures
  if [[ $exit_code -ne 0 ]]; then
    local now=$(date +%s)
    local elapsed=$((now - AI_LAST_CALL))
    
    if [[ $elapsed -ge $AI_COOLDOWN ]]; then
      echo -e "\n\033[1;36m🤖 AI Assistant\033[0m \033[0;33m(exit $exit_code)\033[0m:"
      ai-helper analyze "$LAST_CMD" "$exit_code" "$LAST_OUTPUT"
      local ai_exit=$?
      AI_LAST_CALL=$now

      # Exit code 7 = user chose "Disable AI for session"
      if [[ $ai_exit -eq 7 ]]; then
        AI_SESSION_DISABLED=1
      fi
    fi
  fi
}

# Capture stderr
exec 2> >(while IFS= read -r line; do 
  LAST_OUTPUT+="$line"$'\n'
  echo "$line" >&2
done)

# ============================================================================
# PROACTIVE COMMANDS
# ============================================================================

# Manual AI invocation
ai() {
  if [[ -n "$LAST_CMD" ]]; then
    echo -e "\033[1;36m🤖 Analyzing last command:\033[0m"
    ai-helper analyze "$LAST_CMD" "$LAST_EXIT_CODE" "$LAST_OUTPUT"
  else
    echo -e "\033[0;33m⚠️  No previous command found\033[0m"
  fi
}

# Proactive mode: Ask AI to generate command
ask() {
  if [[ -z "$*" ]]; then
    echo -e "\033[1mUsage:\033[0m \033[0;32mask\033[0m \033[0;36m<natural language query>\033[0m"
    echo ""
    echo -e "\033[1mExamples:\033[0m"
    echo -e "  \033[0;32mask\033[0m \033[0;36mhow do I list all running pods\033[0m"
    echo -e "  \033[0;32mask\033[0m \033[0;36mshow me disk usage of all directories\033[0m"
    return 1
  fi
  
  ai-helper proactive "$*"
}

# Tool-specific helpers
kask() { ai-helper proactive "kubernetes: $*"; }
dask() { ai-helper proactive "docker: $*"; }
task() { ai-helper proactive "terraform: $*"; }
gask() { ai-helper proactive "git: $*"; }
hask() { ai-helper proactive "helm: $*"; }
tgask() { ai-helper proactive "terragrunt: $*"; }
aask() { ai-helper proactive "ansible: $*"; }
arask() { ai-helper proactive "argocd: $*"; }

# Hotkeys
bindkey -s '^[a' 'ai\n'
bindkey -s '^[k' 'ask '

# Aliases
alias ai-stats='ai-helper cache-stats'
alias ai-clear='ai-helper cache-clear'
alias ai-version='ai-helper version'

# Welcome message
echo -e "\033[1;32m✅ AI Terminal Helper Loaded\033[0m"
echo ""
echo -e "\033[1;36mCommands:\033[0m"
echo -e "  \033[0;32mai\033[0m          - Re-analyze last failed command"
echo -e "  \033[0;32mask\033[0m \033[0;33m<query>\033[0m - Generate command from natural language"
echo ""
echo -e "\033[1;36mTool-specific:\033[0m"
echo -e "  \033[0;32mkask\033[0m  - kubectl  \033[0;32mdask\033[0m  - docker   \033[0;32mtask\033[0m  - terraform"
echo -e "  \033[0;32mgask\033[0m  - git      \033[0;32mhask\033[0m  - helm     \033[0;32mtgask\033[0m - terragrunt"
echo -e "  \033[0;32maask\033[0m  - ansible  \033[0;32marask\033[0m - argocd"
echo ""

