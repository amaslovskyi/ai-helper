package helm

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator implements the Validator interface for helm commands.
type Validator struct {
	validSubcommands []string
	dangerousOps     []string
}

// NewValidator creates a new helm validator.
func NewValidator() *Validator {
	return &Validator{
		validSubcommands: []string{
			"install", "upgrade", "uninstall", "rollback", "list", "history",
			"status", "get", "create", "dependency", "env", "lint", "package",
			"plugin", "pull", "push", "registry", "repo", "search", "show",
			"template", "test", "verify", "version",
		},
		dangerousOps: []string{
			"uninstall", "delete", "rollback",
		},
	}
}

// CanValidate returns true if this validator can handle the command.
func (v *Validator) CanValidate(command string) bool {
	return strings.HasPrefix(command, "helm")
}

// Validate checks if a helm command is valid.
func (v *Validator) Validate(command string) error {
	if !strings.HasPrefix(command, "helm") {
		return nil // Not a helm command
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

	// Check for common mistakes
	if err := v.checkCommonMistakes(command); err != nil {
		return err
	}

	return nil
}

// checkHallucinatedFlags checks for commonly hallucinated helm flags.
func (v *Validator) checkHallucinatedFlags(command string) error {
	hallucinatedPatterns := map[string]string{
		`install --update`:        "helm install does not have --update. Use 'helm upgrade --install' instead.",
		`upgrade --force-install`: "helm upgrade does not have --force-install. Use 'helm upgrade --install' instead.",
		`--auto-approve`:          "helm does not have --auto-approve. Helm operations proceed without confirmation by default.",
		`list --sort`:             "helm list does not have --sort. Use --date or --reverse instead.",
		`--force-yes`:             "helm does not have --force-yes.",
		`install --dry-run`:       "use --dry-run=client or --dry-run=server, not just --dry-run.",
		`--no-hooks`:              "use --no-hooks (correct), but be aware this skips pre/post hooks.",
		`repo add --update`:       "helm repo add does not have --update. Run 'helm repo update' separately.",
		`--version latest`:        "helm doesn't support 'latest' as a version. Omit --version to get the latest.",
		`install --replace`:       "helm install does not have --replace. Use 'helm upgrade --install' instead.",
	}

	for pattern, message := range hallucinatedPatterns {
		if matched, _ := regexp.MatchString(pattern, command); matched {
			return fmt.Errorf("%s", message)
		}
	}

	return nil
}

// checkSubcommand validates the helm subcommand.
func (v *Validator) checkSubcommand(command string) error {
	parts := strings.Fields(command)
	if len(parts) < 2 {
		return fmt.Errorf("incomplete helm command")
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
		return fmt.Errorf("'%s' is not a valid helm subcommand. Use 'helm --help' to see valid commands.", subcommand)
	}

	return nil
}

// checkDangerousOps checks for dangerous helm operations.
func (v *Validator) checkDangerousOps(command string) string {
	for _, danger := range v.dangerousOps {
		if strings.Contains(command, danger) {
			return fmt.Sprintf("Dangerous operation detected: 'helm %s'. This will modify or remove deployed applications.", danger)
		}
	}

	// Check for uninstall --purge (deprecated in Helm 3)
	if strings.Contains(command, "uninstall") && strings.Contains(command, "--purge") {
		return "The --purge flag is deprecated in Helm 3. Uninstall now purges by default."
	}

	return ""
}

// checkCommonMistakes checks for common helm command mistakes.
func (v *Validator) checkCommonMistakes(command string) error {
	// Check for helm 2 vs helm 3 differences
	if strings.Contains(command, " delete ") {
		return fmt.Errorf("'helm delete' is deprecated in Helm 3. Use 'helm uninstall' instead.")
	}

	// Check for missing namespace in helm 3
	if (strings.Contains(command, "install") || strings.Contains(command, "upgrade")) &&
		!strings.Contains(command, "-n ") && !strings.Contains(command, "--namespace") {
		return fmt.Errorf("⚠️  No namespace specified. In Helm 3, releases are namespaced. Add -n <namespace> or use --namespace.")
	}

	// Check for install without release name
	if strings.Contains(command, "helm install") {
		parts := strings.Fields(command)
		if len(parts) < 4 {
			return fmt.Errorf("helm install requires: helm install [NAME] [CHART] [flags]")
		}
	}

	// Check for missing repo name format
	if strings.Contains(command, "install") || strings.Contains(command, "upgrade") {
		if matched, _ := regexp.MatchString(`(install|upgrade)\s+\S+\s+\S+/\S+`, command); !matched {
			// Check if it's a local chart (.)
			if !strings.Contains(command, " ./") && !strings.Contains(command, " .") {
				return fmt.Errorf("⚠️  Chart should be in format: repo/chart or ./local-path. Did you forget 'helm repo add'?")
			}
		}
	}

	return nil
}

