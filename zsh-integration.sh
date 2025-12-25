# ============================================================================
# AI Terminal Helper - ZSH Integration
# Add this to your ~/.zshrc file
# ============================================================================

# Colors for output
COLOR_RESET='\033[0m'
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW='\033[0;33m'
COLOR_BLUE='\033[0;34m'
COLOR_CYAN='\033[0;36m'
COLOR_RED='\033[0;31m'
COLOR_BOLD='\033[1m'

# Load ZSH hooks if not already loaded
autoload -Uz add-zsh-hook

# ============================================================================
# STATE TRACKING
# ============================================================================

# Track last command and output for AI assistance
LAST_CMD=""
LAST_OUTPUT=""
LAST_EXIT_CODE=0

# Rate limiting (prevent spam on rapid failures)
AI_LAST_CALL=0
AI_COOLDOWN=2  # seconds between AI calls

# ============================================================================
# COMMAND EXECUTION HOOKS
# ============================================================================

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
    # Rate limiting check
    local now=$(date +%s)
    local elapsed=$((now - AI_LAST_CALL))
    
    if [[ $elapsed -ge $AI_COOLDOWN ]]; then
      echo -e "\n${COLOR_CYAN}${COLOR_BOLD}🤖 AI Assistant${COLOR_RESET} ${COLOR_YELLOW}(exit $exit_code)${COLOR_RESET}:"
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

# ============================================================================
# PROACTIVE AI COMMANDS
# ============================================================================

# Manual AI invocation for any previous command
ai() {
  if [[ -n "$LAST_CMD" ]]; then
    echo -e "${COLOR_CYAN}${COLOR_BOLD}🤖 Analyzing last command:${COLOR_RESET}"
    ~/.ai/ai-helper.sh "$LAST_CMD" "$LAST_OUTPUT" "$LAST_EXIT_CODE"
  else
    echo -e "${COLOR_YELLOW}⚠️  No previous command found${COLOR_RESET}"
  fi
}

# Proactive mode: Ask AI to generate command from natural language
# Usage: ask "how do I list all pods in production namespace"
ask() {
  local query="$*"
  
  if [[ -z "$query" ]]; then
    echo -e "${COLOR_BOLD}Usage:${COLOR_RESET} ${COLOR_GREEN}ask${COLOR_RESET} ${COLOR_CYAN}<natural language query>${COLOR_RESET}"
    echo ""
    echo -e "${COLOR_BOLD}Examples:${COLOR_RESET}"
    echo -e "  ${COLOR_GREEN}ask${COLOR_RESET} ${COLOR_CYAN}how do I list all running pods${COLOR_RESET}"
    echo -e "  ${COLOR_GREEN}ask${COLOR_RESET} ${COLOR_CYAN}show me disk usage of all directories${COLOR_RESET}"
    echo -e "  ${COLOR_GREEN}ask${COLOR_RESET} ${COLOR_CYAN}find all python files modified today${COLOR_RESET}"
    echo -e "  ${COLOR_GREEN}ask${COLOR_RESET} ${COLOR_CYAN}deploy my app to kubernetes${COLOR_RESET}"
    return 1
  fi
  
  # Detect context for better suggestions
  local context=""
  if [[ "$query" =~ (kubernetes|k8s|pod|deployment|service) ]]; then
    context="k8s_context=$(kubectl config current-context 2>/dev/null || echo 'none')"
  elif [[ "$query" =~ (terraform|tf) ]]; then
    context="tf_workspace=$(terraform workspace show 2>/dev/null || echo 'none')"
  elif [[ "$query" =~ (docker|container) ]]; then
    context="docker_info=$(docker info --format '{{.Name}}' 2>/dev/null || echo 'none')"
  elif [[ "$query" =~ (git) ]]; then
    context="git_branch=$(git branch --show-current 2>/dev/null || echo 'none')"
  fi
  
  echo -e "${COLOR_CYAN}${COLOR_BOLD}🤖 Generating command for:${COLOR_RESET} ${COLOR_YELLOW}$query${COLOR_RESET}"
  ~/.ai/ai-helper.sh "PROACTIVE_QUERY" "$query" "0" "$context"
}

# Quick kubectl question helper
kask() {
  ask "kubernetes: $*"
}

# Quick docker question helper
dask() {
  ask "docker: $*"
}

# Quick terraform question helper
task() {
  ask "terraform: $*"
}

# Quick git question helper
gask() {
  ask "git: $*"
}

# ============================================================================
# HOTKEY BINDINGS
# ============================================================================

# Bind Option+A (⌥A) to re-analyze last failure
bindkey -s '^[a' 'ai\n'

# Bind Option+K (⌥K) for quick ask
bindkey -s '^[k' 'ask '

# ============================================================================
# ALIASES
# ============================================================================

# Quick access to AI stats and management
alias ai-stats='cat ~/.ai/history.log | wc -l | xargs echo "Total AI assists:"'
alias ai-clear='rm -f ~/.ai/rate_limit.log && echo "Rate limit cleared"'
alias ai-history='tail -20 ~/.ai/history.log'
alias ai-cache='ls -lh ~/.ai/cache.json 2>/dev/null || echo "No cache yet"'

# ============================================================================
# USAGE TIPS
# ============================================================================

# Show tips on first load (comment out after reading)
echo -e "${COLOR_GREEN}${COLOR_BOLD}✅ AI Terminal Helper v2.0 Loaded!${COLOR_RESET}"
echo ""
echo -e "${COLOR_CYAN}${COLOR_BOLD}Quick Commands:${COLOR_RESET}"
echo -e "  ${COLOR_GREEN}ai${COLOR_RESET}          - Re-analyze last failed command"
echo -e "  ${COLOR_GREEN}ask${COLOR_RESET} ${COLOR_YELLOW}<query>${COLOR_RESET} - Generate command from natural language"
echo -e "  ${COLOR_GREEN}kask${COLOR_RESET} ${COLOR_YELLOW}<q>${COLOR_RESET}    - Kubernetes-specific query"
echo -e "  ${COLOR_GREEN}dask${COLOR_RESET} ${COLOR_YELLOW}<q>${COLOR_RESET}    - Docker-specific query"
echo -e "  ${COLOR_GREEN}task${COLOR_RESET} ${COLOR_YELLOW}<q>${COLOR_RESET}    - Terraform-specific query"
echo -e "  ${COLOR_GREEN}gask${COLOR_RESET} ${COLOR_YELLOW}<q>${COLOR_RESET}    - Git-specific query"
echo ""
echo -e "${COLOR_CYAN}${COLOR_BOLD}Hotkeys:${COLOR_RESET}"
echo -e "  ${COLOR_YELLOW}⌥A (Option+A)${COLOR_RESET} - Re-analyze last failure"
echo -e "  ${COLOR_YELLOW}⌥K (Option+K)${COLOR_RESET} - Quick ask mode"
echo ""
echo -e "${COLOR_CYAN}${COLOR_BOLD}Examples:${COLOR_RESET}"
echo -e "  ${COLOR_GREEN}ask${COLOR_RESET} how do I list all running pods"
echo -e "  ${COLOR_GREEN}kask${COLOR_RESET} show pods in production namespace"
echo -e "  ${COLOR_GREEN}dask${COLOR_RESET} list all containers using more than 1GB RAM"
echo ""
echo -e "${COLOR_CYAN}${COLOR_BOLD}Features:${COLOR_RESET}"
echo -e "  🔒 Security scanning for dangerous commands"
echo -e "  ⏱️  Smart rate limiting (prevents AI spam)"
echo -e "  🧠 Command history learning"
echo -e "  🚀 Context-aware suggestions"
echo ""

