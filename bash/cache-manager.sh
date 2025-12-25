#!/usr/bin/env bash

# ============================================================================
# AI Terminal Helper - Cache Manager
# Provides instant responses for common errors without AI calls
# ============================================================================

# Colors
COLOR_RESET='\033[0m'
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW='\033[0;33m'
COLOR_BLUE='\033[0;34m'
COLOR_MAGENTA='\033[0;35m'
COLOR_CYAN='\033[0;36m'
COLOR_RED='\033[0;31m'
COLOR_BOLD='\033[1m'
COLOR_MAGENTA_BOLD='\033[1;35m'

CACHE_FILE="${HOME}/.ai/cache.json"

# ============================================================================
# CACHE LOOKUP
# ============================================================================

# Check if a command+error pattern exists in cache
# Returns: cached fix if found, empty if not found
lookup_cache() {
  local cmd="$1"
  local error="$2"
  
  # Create cache file if it doesn't exist
  if [[ ! -f "$CACHE_FILE" ]]; then
    init_cache
  fi
  
  # Create cache key (hash of command + first line of error)
  local error_first_line=$(echo "$error" | head -1 | tr -d '\n')
  local cache_key=$(echo "${cmd}::${error_first_line}" | md5sum | cut -d' ' -f1)
  
  # Look up in cache using jq if available, fallback to grep
  if command -v jq &> /dev/null; then
    local cached=$(jq -r ".\"$cache_key\".fix // empty" "$CACHE_FILE" 2>/dev/null)
    if [[ -n "$cached" ]]; then
      echo -e "${COLOR_MAGENTA_BOLD}💾 [Cached]${COLOR_RESET} $cached"
      return 0
    fi
  fi
  
  return 1
}

# ============================================================================
# CACHE UPDATE
# ============================================================================

# Save a successful fix to cache for future use
save_to_cache() {
  local cmd="$1"
  local error="$2"
  local fix="$3"
  
  # Create cache key
  local error_first_line=$(echo "$error" | head -1 | tr -d '\n')
  local cache_key=$(echo "${cmd}::${error_first_line}" | md5sum | cut -d' ' -f1)
  local timestamp=$(date +%s)
  
  # Save to cache (create simple JSON structure)
  if command -v jq &> /dev/null; then
    # Use jq if available for proper JSON
    local temp_file=$(mktemp)
    jq --arg key "$cache_key" \
       --arg cmd "$cmd" \
       --arg err "$error_first_line" \
       --arg fix "$fix" \
       --arg ts "$timestamp" \
       '.[$key] = {cmd: $cmd, error: $err, fix: $fix, timestamp: $ts, hits: 0}' \
       "$CACHE_FILE" > "$temp_file" 2>/dev/null
    mv "$temp_file" "$CACHE_FILE"
  else
    # Fallback: simple append (not proper JSON, but works)
    echo "{\"$cache_key\": {\"cmd\": \"$cmd\", \"error\": \"$error_first_line\", \"fix\": \"$fix\", \"timestamp\": $timestamp}}" >> "$CACHE_FILE"
  fi
}

# ============================================================================
# CACHE INITIALIZATION
# ============================================================================

# Initialize cache with common error patterns
init_cache() {
  mkdir -p "$(dirname "$CACHE_FILE")"
  
  # Create cache with pre-populated common errors
  cat > "$CACHE_FILE" << 'EOF'
{
  "common_kubectl_not_found": {
    "cmd": "kubectl get pods",
    "error": "kubectl: command not found",
    "fix": "✓ brew install kubectl\nRoot: kubectl CLI not installed\nTip: Verify installation with: kubectl version",
    "timestamp": 0
  },
  "common_docker_permission": {
    "cmd": "docker ps",
    "error": "permission denied",
    "fix": "✓ sudo docker ps\nRoot: Docker daemon requires root or docker group membership\nTip: Add user to docker group: sudo usermod -aG docker $USER",
    "timestamp": 0
  },
  "common_terraform_init": {
    "cmd": "terraform apply",
    "error": "terraform has not been initialized",
    "fix": "✓ terraform init\nRoot: Terraform working directory not initialized\nTip: Always run 'terraform init' in new directories",
    "timestamp": 0
  },
  "common_git_not_repo": {
    "cmd": "git status",
    "error": "not a git repository",
    "fix": "✓ git init\nRoot: Current directory is not a Git repository\nTip: Use 'git clone' for existing repos or 'git init' for new ones",
    "timestamp": 0
  },
  "common_npm_not_found": {
    "cmd": "npm install",
    "error": "npm: command not found",
    "fix": "✓ brew install node\nRoot: Node.js and npm not installed\nTip: Verify with: node --version && npm --version",
    "timestamp": 0
  },
  "common_python_module": {
    "cmd": "python script.py",
    "error": "ModuleNotFoundError",
    "fix": "✓ pip install <module_name>\nRoot: Python module not installed\nTip: Use 'pip list' to see installed packages",
    "timestamp": 0
  },
  "common_permission_denied": {
    "cmd": "rm file.txt",
    "error": "Permission denied",
    "fix": "✓ sudo rm file.txt\nRoot: Insufficient permissions\nTip: Check ownership with 'ls -l' before using sudo",
    "timestamp": 0
  },
  "common_no_such_file": {
    "cmd": "cat file.txt",
    "error": "No such file or directory",
    "fix": "✓ ls -la\nRoot: File doesn't exist at specified path\nTip: Verify file path and spelling",
    "timestamp": 0
  },
  "common_port_in_use": {
    "cmd": "docker run -p 8080:80",
    "error": "address already in use",
    "fix": "✓ lsof -i :8080\nRoot: Port 8080 is already occupied\nTip: Kill process or use different port: -p 8081:80",
    "timestamp": 0
  },
  "common_connection_refused": {
    "cmd": "curl localhost:8080",
    "error": "Connection refused",
    "fix": "✓ Check if service is running: ps aux | grep <service>\nRoot: Target service is not running or unreachable\nTip: Verify service status and network connectivity",
    "timestamp": 0
  }
}
EOF
  
  echo -e "${COLOR_GREEN}✅ Cache initialized with common error patterns${COLOR_RESET}"
}

# ============================================================================
# CACHE MANAGEMENT
# ============================================================================

# Show cache statistics
show_cache_stats() {
  if [[ ! -f "$CACHE_FILE" ]]; then
    echo -e "${COLOR_RED}❌ No cache file found${COLOR_RESET}"
    return 1
  fi
  
  if command -v jq &> /dev/null; then
    local total=$(jq 'length' "$CACHE_FILE" 2>/dev/null || echo "0")
    echo -e "${COLOR_CYAN}${COLOR_BOLD}📊 Cache Statistics:${COLOR_RESET}"
    echo -e "  ${COLOR_YELLOW}Total patterns:${COLOR_RESET} ${COLOR_GREEN}$total${COLOR_RESET}"
    echo -e "  ${COLOR_YELLOW}Cache file:${COLOR_RESET} ${COLOR_BLUE}$CACHE_FILE${COLOR_RESET}"
    echo -e "  ${COLOR_YELLOW}Size:${COLOR_RESET} ${COLOR_GREEN}$(du -h "$CACHE_FILE" | cut -f1)${COLOR_RESET}"
  else
    echo -e "${COLOR_CYAN}📊 Cache exists at:${COLOR_RESET} ${COLOR_BLUE}$CACHE_FILE${COLOR_RESET}"
    echo -e "  ${COLOR_YELLOW}Install 'jq' for detailed statistics:${COLOR_RESET} ${COLOR_GREEN}brew install jq${COLOR_RESET}"
  fi
}

# Clear cache
clear_cache() {
  if [[ -f "$CACHE_FILE" ]]; then
    rm "$CACHE_FILE"
    echo -e "${COLOR_GREEN}✅ Cache cleared${COLOR_RESET}"
    init_cache
  else
    echo -e "${COLOR_YELLOW}❌ No cache to clear${COLOR_RESET}"
  fi
}

# Export functions for use in ai-helper.sh
export -f lookup_cache
export -f save_to_cache
export -f init_cache

# ============================================================================
# CLI INTERFACE (when run directly)
# ============================================================================

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  case "${1:-}" in
    init)
      init_cache
      ;;
    stats)
      show_cache_stats
      ;;
    clear)
      clear_cache
      ;;
    *)
      echo "Usage: $0 {init|stats|clear}"
      echo ""
      echo "Commands:"
      echo "  init   - Initialize cache with common patterns"
      echo "  stats  - Show cache statistics"
      echo "  clear  - Clear and reinitialize cache"
      ;;
  esac
fi

