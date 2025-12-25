package argocd

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator implements the Validator interface for argocd CLI commands.
type Validator struct {
	validSubcommands []string
	dangerousOps     []string
}

// NewValidator creates a new argocd validator.
func NewValidator() *Validator {
	return &Validator{
		validSubcommands: []string{
			"account", "admin", "app", "appset", "cert", "cluster", "completion",
			"context", "gpg", "login", "logout", "proj", "relogin", "repo",
			"repocreds", "version",
		},
		dangerousOps: []string{
			"app delete", "app terminate-op", "admin", "cluster rm",
		},
	}
}

// CanValidate returns true if this validator can handle the command.
func (v *Validator) CanValidate(command string) bool {
	return strings.HasPrefix(command, "argocd")
}

// Validate checks if an argocd command is valid.
func (v *Validator) Validate(command string) error {
	if !strings.HasPrefix(command, "argocd") {
		return nil // Not an argocd command
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
	if warning := v.checkDangerousOps(command); warning != "" {
		return fmt.Errorf("⚠️  Warning: %s", warning)
	}

	return nil
}

// checkHallucinatedFlags checks for commonly hallucinated argocd flags.
func (v *Validator) checkHallucinatedFlags(command string) error {
	hallucinatedPatterns := map[string]string{
		`app create --auto-sync`:    "use --sync-policy automated (not --auto-sync).",
		`app sync --wait`:           "use --timeout instead of --wait.",
		`app delete --force`:        "use --cascade (not --force) to control deletion behavior.",
		`--namespace`:               "use --dest-namespace for app destination, or --app-namespace for app placement.",
		`app rollback`:              "argocd doesn't have 'app rollback'. Use 'app sync --revision <version>' instead.",
		`--auto-approve`:            "argocd does not have --auto-approve.",
		`app deploy`:                "argocd doesn't have 'app deploy'. Use 'app sync' instead.",
		`app list --sort`:           "use --sort-by (not --sort).",
	}

	for pattern, message := range hallucinatedPatterns {
		if matched, _ := regexp.MatchString(pattern, command); matched {
			return fmt.Errorf("%s", message)
		}
	}

	return nil
}

// checkSubcommand validates the argocd subcommand.
func (v *Validator) checkSubcommand(command string) error {
	parts := strings.Fields(command)
	if len(parts) < 2 {
		return fmt.Errorf("incomplete argocd command")
	}

	subcommand := parts[1]
	
	// Check if it's a valid subcommand
	isValid := false
	for _, valid := range v.validSubcommands {
		if subcommand == valid {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("'%s' is not a valid argocd subcommand. Use 'argocd --help' to see valid commands.", subcommand)
	}

	return nil
}

// checkDangerousOps checks for dangerous argocd operations.
func (v *Validator) checkDangerousOps(command string) string {
	for _, danger := range v.dangerousOps {
		if strings.Contains(command, danger) {
			switch {
			case strings.Contains(danger, "delete"):
				return "Dangerous: This will delete the ArgoCD application and potentially the deployed resources!"
			case strings.Contains(danger, "admin"):
				return "Admin commands can modify ArgoCD's core configuration. Be careful!"
			case strings.Contains(danger, "terminate-op"):
				return "Terminating an operation might leave the application in an inconsistent state."
			case strings.Contains(danger, "cluster rm"):
				return "This will remove cluster from ArgoCD, affecting all apps deployed to it!"
			}
		}
	}
	return ""
}

