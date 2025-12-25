package git

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator implements the Validator interface for git commands.
type Validator struct {
	validSubcommands []string
	dangerousOps     []DangerousOp
}

// DangerousOp represents a dangerous git operation.
type DangerousOp struct {
	Pattern     *regexp.Regexp
	Warning     string
	Suggestion  string
	IsBlocking  bool // If true, block the command; if false, just warn
}

// NewValidator creates a new git validator.
func NewValidator() *Validator {
	return &Validator{
		validSubcommands: []string{
			"add", "bisect", "branch", "checkout", "cherry-pick", "clean", "clone",
			"commit", "config", "diff", "fetch", "grep", "init", "log", "merge",
			"mv", "pull", "push", "rebase", "remote", "reset", "restore", "revert",
			"rm", "show", "stash", "status", "switch", "tag", "worktree",
		},
		dangerousOps: []DangerousOp{
			{
				Pattern:    regexp.MustCompile(`push\s+(-f|--force)\s+.*\s+(main|master)$`),
				Warning:    "Force pushing to main/master branch is dangerous and can cause data loss for other team members.",
				Suggestion: "Use --force-with-lease instead, or push to a feature branch first.",
				IsBlocking: true,
			},
			{
				Pattern:    regexp.MustCompile(`reset\s+--hard`),
				Warning:    "git reset --hard will permanently delete uncommitted changes.",
				Suggestion: "Consider using 'git stash' to save changes, or 'git reset --soft' to keep changes staged.",
				IsBlocking: false,
			},
			{
				Pattern:    regexp.MustCompile(`clean\s+-[dfx]+`),
				Warning:    "git clean will permanently delete untracked files.",
				Suggestion: "Run with -n (--dry-run) first to see what will be deleted.",
				IsBlocking: false,
			},
			{
				Pattern:    regexp.MustCompile(`push\s+(-f|--force)`),
				Warning:    "Force pushing can overwrite remote history and cause issues for collaborators.",
				Suggestion: "Use --force-with-lease instead, which is safer.",
				IsBlocking: false,
			},
			{
				Pattern:    regexp.MustCompile(`rebase\s+.*\s+(main|master)`),
				Warning:    "Rebasing main/master can cause issues if others are working on it.",
				Suggestion: "Consider merging instead, or ensure you're not rewriting shared history.",
				IsBlocking: false,
			},
			{
				Pattern:    regexp.MustCompile(`branch\s+-D`),
				Warning:    "git branch -D will force delete a branch, even with unmerged changes.",
				Suggestion: "Use -d instead to safely delete only merged branches.",
				IsBlocking: false,
			},
		},
	}
}

// CanValidate returns true if this validator can handle the command.
func (v *Validator) CanValidate(command string) bool {
	// Check for git or Oh My Zsh git aliases
	if strings.HasPrefix(command, "git") {
		return true
	}
	
	// Common Oh My Zsh git aliases
	gitAliases := []string{
		"gco", "gcb", "gcm", "gcd", "gcmg", "ga", "gaa", "gc", "gc!", "gcmsg",
		"gca", "gca!", "gcam", "gb", "gba", "gbd", "gbD", "gst", "gss", "gd",
		"gdca", "gp", "gpf", "gpf!", "gl", "ggl", "ggp", "gf", "gfa", "glog",
		"glol", "glola", "gm", "grb", "grbi", "grbc", "grba", "gsta", "gstp",
		"gstl", "gr", "gra", "grv", "grmv", "grrm", "gcl", "grh", "grhh", "gclean",
	}
	
	firstWord := strings.Fields(command)[0]
	for _, alias := range gitAliases {
		if firstWord == alias {
			return true
		}
	}
	
	return false
}

// Validate checks if a git command is valid.
func (v *Validator) Validate(command string) error {
	// Resolve Oh My Zsh aliases to git commands
	command = v.resolveGitAlias(command)
	
	if !strings.HasPrefix(command, "git") {
		return nil // Not a git command
	}

	// Check for common hallucinated flags
	if err := v.checkHallucinatedFlags(command); err != nil {
		return err
	}

	// Check for invalid subcommands
	if err := v.checkSubcommand(command); err != nil {
		return err
	}

	// Check for dangerous operations
	if err := v.checkDangerousOps(command); err != nil {
		return err
	}

	return nil
}

// checkHallucinatedFlags checks for commonly hallucinated git flags.
func (v *Validator) checkHallucinatedFlags(command string) error {
	hallucinatedPatterns := map[string]string{
		`push --force-all`:     "git push does not have a --force-all flag. Use --force or --force-with-lease.",
		`commit --push`:        "git commit does not have a --push flag. Run 'git push' separately.",
		`pull --commit`:        "git pull does not have a --commit flag. Commits are created automatically during merge.",
		`log --sort`:           "git log does not have a --sort flag. Use --author-date-order or --date-order instead.",
		`branch --rename-all`:  "git branch does not have a --rename-all flag. Use -m to rename a single branch.",
		`merge --force-merge`:  "git merge does not have a --force-merge flag.",
		`checkout --create`:    "git checkout does not have a --create flag. Use -b to create a new branch.",
		`stash --list`:         "Use 'git stash list' (without --list flag).",
		`rebase --interactive`: "Use 'git rebase -i' (not --interactive).",
	}

	for pattern, message := range hallucinatedPatterns {
		if matched, _ := regexp.MatchString(pattern, command); matched {
			return fmt.Errorf("%s", message)
		}
	}

	return nil
}

// checkSubcommand validates the git subcommand.
func (v *Validator) checkSubcommand(command string) error {
	parts := strings.Fields(command)
	if len(parts) < 2 {
		return fmt.Errorf("incomplete git command")
	}

	subcommand := parts[1]
	
	// Check if it's a valid subcommand (or a flag like --version)
	if strings.HasPrefix(subcommand, "-") {
		return nil // It's a flag, not a subcommand
	}
	
	// Check if it's a valid subcommand
	isValid := false
	for _, valid := range v.validSubcommands {
		if subcommand == valid {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("'%s' is not a valid git subcommand. Use 'git --help' to see valid commands.", subcommand)
	}

	return nil
}

// checkDangerousOps checks for dangerous git operations.
func (v *Validator) checkDangerousOps(command string) error {
	for _, danger := range v.dangerousOps {
		if danger.Pattern.MatchString(command) {
			if danger.IsBlocking {
				return fmt.Errorf("🚨 BLOCKED: %s\n💡 Suggestion: %s", danger.Warning, danger.Suggestion)
			}
			return fmt.Errorf("⚠️  Warning: %s\n💡 Suggestion: %s", danger.Warning, danger.Suggestion)
		}
	}
	return nil
}

// resolveGitAlias converts Oh My Zsh git aliases to full git commands.
func (v *Validator) resolveGitAlias(command string) string {
	aliases := map[string]string{
		"gco":  "git checkout",
		"gcb":  "git checkout -b",
		"gcm":  "git checkout master",
		"gcd":  "git checkout develop",
		"gcmg": "git checkout main",
		"ga":    "git add",
		"gaa":   "git add --all",
		"gc":    "git commit -v",
		"gc!":   "git commit -v --amend",
		"gcmsg": "git commit -m",
		"gca":   "git commit -v -a",
		"gca!":  "git commit -v -a --amend",
		"gcam":  "git commit -a -m",
		"gb":  "git branch",
		"gba": "git branch -a",
		"gbd": "git branch -d",
		"gbD": "git branch -D",
		"gst":  "git status",
		"gss":  "git status -s",
		"gd":   "git diff",
		"gdca": "git diff --cached",
		"gp":   "git push",
		"gpf":  "git push --force",
		"gpf!": "git push --force",
		"gl":   "git pull",
		"ggl":  "git pull origin",
		"ggp":  "git push origin",
		"gf":   "git fetch",
		"gfa":  "git fetch --all",
		"glog":  "git log --oneline --decorate",
		"glol":  "git log --graph --pretty='%Cred%h%Creset -%C(auto)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset'",
		"glola": "git log --graph --pretty='%Cred%h%Creset -%C(auto)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --all",
		"gm":   "git merge",
		"grb":  "git rebase",
		"grbi": "git rebase -i",
		"grbc": "git rebase --continue",
		"grba": "git rebase --abort",
		"gsta": "git stash",
		"gstp": "git stash pop",
		"gstl": "git stash list",
		"gr":   "git remote",
		"gra":  "git remote add",
		"grv":  "git remote -v",
		"grmv": "git remote rename",
		"grrm": "git remote remove",
		"gcl": "git clone",
		"grh":   "git reset",
		"grhh":  "git reset --hard",
		"gclean": "git clean -fd",
	}
	
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return command
	}
	
	firstWord := parts[0]
	if fullCmd, exists := aliases[firstWord]; exists {
		remaining := ""
		if len(parts) > 1 {
			remaining = " " + strings.Join(parts[1:], " ")
		}
		return fullCmd + remaining
	}
	
	return command
}

