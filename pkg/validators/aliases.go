package validators

import (
	"strings"
)

// AliasMapper handles command alias resolution
type AliasMapper struct {
	aliases map[string]string
}

// NewAliasMapper creates a new alias mapper with common DevOps aliases
func NewAliasMapper() *AliasMapper {
	return &AliasMapper{
		aliases: map[string]string{
			// Kubernetes
			"k": "kubectl",

			// Terraform
			"tf": "terraform",

			// Terragrunt
			"tg": "terragrunt",

			// Helm
			"h": "helm",

			// Docker
			"d":  "docker",
			"dc": "docker-compose",

			// Oh My Zsh Git Aliases (most common ones)
			// Checkout
			"gco":  "git checkout",
			"gcb":  "git checkout -b",
			"gcm":  "git checkout master",
			"gcd":  "git checkout develop",
			"gcmg": "git checkout main",

			// Add/Commit
			"ga":    "git add",
			"gaa":   "git add --all",
			"gc":    "git commit -v",
			"gc!":   "git commit -v --amend",
			"gcmsg": "git commit -m",
			"gca":   "git commit -v -a",
			"gca!":  "git commit -v -a --amend",
			"gcam":  "git commit -a -m",

			// Branch
			"gb":  "git branch",
			"gba": "git branch -a",
			"gbd": "git branch -d",
			"gbD": "git branch -D",

			// Status/Diff
			"gst":  "git status",
			"gss":  "git status -s",
			"gd":   "git diff",
			"gdca": "git diff --cached",

			// Push/Pull/Fetch
			"gp":   "git push",
			"gpf":  "git push --force",
			"gpf!": "git push --force",
			"gl":   "git pull",
			"ggl":  "git pull origin",
			"ggp":  "git push origin",
			"gf":   "git fetch",
			"gfa":  "git fetch --all",

			// Log
			"glog":  "git log --oneline --decorate",
			"glol":  "git log --graph --pretty='%Cred%h%Creset -%C(auto)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset'",
			"glola": "git log --graph --pretty='%Cred%h%Creset -%C(auto)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --all",

			// Merge/Rebase
			"gm":   "git merge",
			"grb":  "git rebase",
			"grbi": "git rebase -i",
			"grbc": "git rebase --continue",
			"grba": "git rebase --abort",

			// Stash
			"gsta": "git stash",
			"gstp": "git stash pop",
			"gstl": "git stash list",

			// Remote
			"gr":   "git remote",
			"gra":  "git remote add",
			"grv":  "git remote -v",
			"grmv": "git remote rename",
			"grrm": "git remote remove",

			// Clone
			"gcl": "git clone",

			// Reset/Clean
			"grh":   "git reset",
			"grhh":  "git reset --hard",
			"gclean": "git clean -fd",

			// AWS CLI
			"a": "aws",

			// Google Cloud
			"g": "gcloud",

			// Azure CLI
			"az": "az",
		},
	}
}

// ResolveAlias converts an aliased command to its full form
func (am *AliasMapper) ResolveAlias(command string) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return command
	}

	firstWord := parts[0]

	// Check if it's a known alias
	if fullCmd, exists := am.aliases[firstWord]; exists {
		// Replace the alias with the full command
		remaining := ""
		if len(parts) > 1 {
			remaining = " " + strings.Join(parts[1:], " ")
		}
		return fullCmd + remaining
	}

	return command
}

// GetToolName extracts the base tool name from a command (handling aliases)
func (am *AliasMapper) GetToolName(command string) string {
	resolved := am.ResolveAlias(command)
	parts := strings.Fields(resolved)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

// IsAlias checks if a string is a known alias
func (am *AliasMapper) IsAlias(cmd string) bool {
	_, exists := am.aliases[cmd]
	return exists
}

