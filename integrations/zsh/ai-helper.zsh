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

# Capture command before execution
preexec() {
  LAST_CMD="$1"
  LAST_OUTPUT=""
}

# Check for errors after command completes
precmd() {
  local exit_code=$?
  LAST_EXIT_CODE=$exit_code
  
  # Only trigger AI on failure
  if [[ $exit_code -ne 0 ]]; then
    local now=$(date +%s)
    local elapsed=$((now - AI_LAST_CALL))
    
    if [[ $elapsed -ge $AI_COOLDOWN ]]; then
      echo -e "\n\033[1;36m🤖 AI Assistant\033[0m \033[0;33m(exit $exit_code)\033[0m:"
      ai-helper analyze "$LAST_CMD" "$exit_code" "$LAST_OUTPUT"
      AI_LAST_CALL=$now
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

# Quick helpers
kask() { ai-helper proactive "kubernetes: $*"; }
dask() { ai-helper proactive "docker: $*"; }
task() { ai-helper proactive "terraform: $*"; }
gask() { ai-helper proactive "git: $*"; }

# Hotkeys
bindkey -s '^[a' 'ai\n'
bindkey -s '^[k' 'ask '

# Aliases
alias ai-stats='ai-helper cache-stats'
alias ai-clear='ai-helper cache-clear'
alias ai-version='ai-helper version'

# Welcome message
echo -e "\033[1;32m✅ AI Terminal Helper v2.1 (Go) Loaded!\033[0m"
echo ""
echo -e "\033[1;36mQuick Commands:\033[0m"
echo -e "  \033[0;32mai\033[0m          - Re-analyze last failed command"
echo -e "  \033[0;32mask\033[0m \033[0;33m<query>\033[0m - Generate command from natural language"
echo -e "  \033[0;32mkask\033[0m/\033[0;32mdask\033[0m/\033[0;32mtask\033[0m/\033[0;32mgask\033[0m - Tool-specific queries"
echo ""
echo -e "\033[1;36mNew in v2.1:\033[0m"
echo -e "  🎯 8 validators (kubectl, terraform, git, helm, terragrunt, ansible, argocd, docker)"
echo -e "  🔤 50+ alias support (k, tf, tg, h, gco, gp, etc.)"
echo -e "  📊 Confidence scoring (High/Medium/Low)"
echo ""
echo -e "\033[1;36mFeatures:\033[0m"
echo -e "  🔒 Security scanning    ⏱️  Smart rate limiting"
echo -e "  ✅ Command validation   💾 Offline caching"
echo -e "  🚀 Fast Go binary       🎯 Fixes hallucinations"
echo ""

