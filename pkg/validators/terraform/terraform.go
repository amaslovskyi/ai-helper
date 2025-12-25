package terraform

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator implements the Validator interface for terraform commands.
type Validator struct {
	validSubcommands []string
	dangerousOps     []string
}

// NewValidator creates a new terraform validator.
func NewValidator() *Validator {
	return &Validator{
		validSubcommands: []string{
			"init", "plan", "apply", "destroy", "validate", "fmt", "force-unlock",
			"get", "graph", "import", "login", "logout", "output", "providers",
			"refresh", "show", "state", "taint", "untaint", "version", "workspace",
			"console", "test", "metadata",
		},
		dangerousOps: []string{
			"destroy", "force-unlock", "taint",
		},
	}
}

// CanValidate returns true if this validator can handle the command.
func (v *Validator) CanValidate(command string) bool {
	return strings.HasPrefix(command, "terraform") || strings.HasPrefix(command, "tf ")
}

// Validate checks if a terraform command is valid.
func (v *Validator) Validate(command string) error {
	// Handle alias 'tf' for 'terraform'
	if strings.HasPrefix(command, "tf ") {
		command = "terraform" + command[2:]
	}
	
	if !strings.HasPrefix(command, "terraform") {
		return nil // Not a terraform command
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

// checkHallucinatedFlags checks for commonly hallucinated terraform flags.
func (v *Validator) checkHallucinatedFlags(command string) error {
	hallucinatedPatterns := map[string]string{
		`plan --apply`:           "terraform plan does not have an --apply flag. Run 'terraform apply' separately.",
		`apply --plan`:           "terraform apply does not take a --plan flag. Use 'terraform apply plan.tfplan' to apply a saved plan.",
		`--force-yes`:            "terraform does not have a --force-yes flag. Use -auto-approve instead.",
		`--skip-validation`:      "terraform does not have a --skip-validation flag.",
		`apply --target-all`:     "terraform apply does not have a --target-all flag. Omit --target to apply all resources.",
		`destroy --force`:        "terraform destroy does not have a --force flag. Use -auto-approve instead.",
		`init --upgrade-modules`: "terraform init does not have --upgrade-modules. Use -upgrade instead.",
		`plan --save`:            "terraform plan does not have a --save flag. Use -out=FILE to save the plan.",
		`apply --dry-run`:        "terraform apply does not have a --dry-run flag. Use 'terraform plan' instead.",
	}

	for pattern, message := range hallucinatedPatterns {
		if matched, _ := regexp.MatchString(pattern, command); matched {
			return fmt.Errorf("%s", message)
		}
	}

	return nil
}

// checkSubcommand validates the terraform subcommand.
func (v *Validator) checkSubcommand(command string) error {
	parts := strings.Fields(command)
	if len(parts) < 2 {
		return fmt.Errorf("incomplete terraform command")
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
		return fmt.Errorf("'%s' is not a valid terraform subcommand. Use 'terraform --help' to see valid commands.", subcommand)
	}

	return nil
}

// checkDangerousOps checks for dangerous terraform operations.
func (v *Validator) checkDangerousOps(command string) string {
	for _, danger := range v.dangerousOps {
		if strings.Contains(command, danger) {
			return fmt.Sprintf("Dangerous operation detected: 'terraform %s'. This will modify or destroy infrastructure.", danger)
		}
	}
	return ""
}

// checkCommonMistakes checks for common terraform command mistakes.
func (v *Validator) checkCommonMistakes(command string) error {
	// Check for missing -auto-approve on apply/destroy
	if strings.Contains(command, "apply") && strings.Contains(command, "auto-approve") {
		if !strings.Contains(command, "-auto-approve") && !strings.Contains(command, "--auto-approve") {
			// They probably meant to use the flag correctly
			return fmt.Errorf("did you mean -auto-approve? (note the single dash)")
		}
	}

	// Check for -target without value
	if matched, _ := regexp.MatchString(`-target\s*$`, command); matched {
		return fmt.Errorf("-target flag requires a value. Example: -target=aws_instance.example")
	}

	// Check for -var without value
	if matched, _ := regexp.MatchString(`-var\s*$`, command); matched {
		return fmt.Errorf("-var flag requires a value. Example: -var=\"key=value\"")
	}

	return nil
}

