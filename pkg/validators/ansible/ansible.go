package ansible

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator implements the Validator interface for ansible/ansible-playbook commands.
type Validator struct {
	validCommands    []string
	dangerousModules []string
}

// NewValidator creates a new ansible validator.
func NewValidator() *Validator {
	return &Validator{
		validCommands: []string{
			"ansible", "ansible-playbook", "ansible-vault", "ansible-galaxy",
			"ansible-config", "ansible-console", "ansible-doc", "ansible-inventory",
			"ansible-pull",
		},
		dangerousModules: []string{
			"shell", "command", "raw", "file", "copy", "template",
		},
	}
}

// CanValidate returns true if this validator can handle the command.
func (v *Validator) CanValidate(command string) bool {
	for _, cmd := range v.validCommands {
		if strings.HasPrefix(command, cmd) {
			return true
		}
	}
	return false
}

// Validate checks if an ansible command is valid.
func (v *Validator) Validate(command string) error {
	// Check for common hallucinated flags
	if err := v.checkHallucinatedFlags(command); err != nil {
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

// checkHallucinatedFlags checks for commonly hallucinated ansible flags.
func (v *Validator) checkHallucinatedFlags(command string) error {
	hallucinatedPatterns := map[string]string{
		`--force`:          "ansible-playbook does not have --force. Use --check for dry-run.",
		`--yes`:            "ansible does not have --yes. It doesn't prompt by default.",
		`--auto-approve`:   "ansible does not have --auto-approve.",
		`--dry-run`:        "use --check (not --dry-run) for ansible dry-run mode.",
		`--parallel`:       "use -f or --forks (not --parallel) to control parallelism.",
		`--inventory-file`: "use -i or --inventory (not --inventory-file).",
		`--playbook`:       "ansible-playbook doesn't need --playbook flag. Just provide the playbook file.",
		`--sudo`:           "--sudo is deprecated. Use --become instead.",
	}

	for pattern, message := range hallucinatedPatterns {
		if matched, _ := regexp.MatchString(pattern, command); matched {
			return fmt.Errorf("%s", message)
		}
	}

	return nil
}

// checkDangerousOps checks for dangerous ansible operations.
func (v *Validator) checkDangerousOps(command string) string {
	// Check for dangerous modules in ad-hoc commands
	if strings.Contains(command, "ansible ") && !strings.Contains(command, "ansible-") {
		for _, module := range v.dangerousModules {
			if strings.Contains(command, "-m "+module) || strings.Contains(command, "-module-name="+module) {
				return fmt.Sprintf("Using module '%s' in ad-hoc command. This can be dangerous. Consider using a playbook for better control and logging.", module)
			}
		}
	}

	// Check for become without limit
	if (strings.Contains(command, "--become") || strings.Contains(command, "-b")) &&
		!strings.Contains(command, "--limit") && !strings.Contains(command, "-l") {
		return "Running with elevated privileges (--become) without --limit. This affects ALL hosts in the inventory!"
	}

	// Check for command module with rm -rf or similar
	if strings.Contains(command, "-m shell") || strings.Contains(command, "-m command") {
		if strings.Contains(command, "rm -rf") || strings.Contains(command, "rm -fr") {
			return "Dangerous command detected: rm -rf in ansible shell/command module. This could delete critical files!"
		}
	}

	return ""
}

// checkCommonMistakes checks for common ansible command mistakes.
func (v *Validator) checkCommonMistakes(command string) error {
	// Check for missing inventory
	if strings.Contains(command, "ansible-playbook") &&
		!strings.Contains(command, "-i ") && !strings.Contains(command, "--inventory") {
		return fmt.Errorf("⚠️  No inventory specified. Ansible will use default /etc/ansible/hosts. Is this intentional?")
	}

	// Check for --syntax-check with other flags that won't work
	if strings.Contains(command, "--syntax-check") {
		if strings.Contains(command, "--check") || strings.Contains(command, "--diff") {
			return fmt.Errorf("--syntax-check only validates syntax. It doesn't run with --check or --diff.")
		}
	}

	return nil
}

