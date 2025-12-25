package terragrunt

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator implements the Validator interface for terragrunt commands.
type Validator struct {
	validSubcommands []string
	dangerousOps     []string
}

// NewValidator creates a new terragrunt validator.
func NewValidator() *Validator {
	return &Validator{
		validSubcommands: []string{
			// Terragrunt-specific commands
			"run-all", "apply-all", "destroy-all", "plan-all", "output-all",
			"validate-all", "graph-dependencies", "hclfmt", "aws-provider-patch",
			"render-json", "validate-inputs", "graph",
			
			// Terraform passthrough commands
			"init", "plan", "apply", "destroy", "validate", "fmt", "force-unlock",
			"get", "import", "login", "logout", "output", "providers",
			"refresh", "show", "state", "taint", "untaint", "version", "workspace",
			"console", "test", "metadata",
		},
		dangerousOps: []string{
			"destroy", "destroy-all", "force-unlock", "apply-all",
		},
	}
}

// CanValidate returns true if this validator can handle the command.
func (v *Validator) CanValidate(command string) bool {
	return strings.HasPrefix(command, "terragrunt") || strings.HasPrefix(command, "tg ")
}

// Validate checks if a terragrunt command is valid.
func (v *Validator) Validate(command string) error {
	if !strings.HasPrefix(command, "terragrunt") && !strings.HasPrefix(command, "tg ") {
		return nil // Not a terragrunt command
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

// checkHallucinatedFlags checks for commonly hallucinated terragrunt flags.
func (v *Validator) checkHallucinatedFlags(command string) error {
	hallucinatedPatterns := map[string]string{
		`--all-modules`:          "terragrunt does not have --all-modules. Use run-all subcommand instead.",
		`--recurse`:              "terragrunt does not have --recurse. Use run-all subcommand instead.",
		`--force-yes`:            "terragrunt does not have --force-yes. Use -auto-approve for terraform commands.",
		`--skip-validation`:      "terragrunt does not have --skip-validation.",
		`apply --target-all`:     "terragrunt apply does not have --target-all. Use apply-all or run-all apply.",
		`destroy --force`:        "terragrunt destroy does not have --force. Use -auto-approve instead.",
		`--skip-dependencies`:    "Use --terragrunt-ignore-dependency-errors instead.",
		`--parallel`:             "Use --terragrunt-parallelism instead.",
		`run-all --auto-approve`: "Use --terragrunt-non-interactive instead for run-all commands.",
	}

	for pattern, message := range hallucinatedPatterns {
		if matched, _ := regexp.MatchString(pattern, command); matched {
			return fmt.Errorf("%s", message)
		}
	}

	return nil
}

// checkSubcommand validates the terragrunt subcommand.
func (v *Validator) checkSubcommand(command string) error {
	parts := strings.Fields(command)
	if len(parts) < 2 {
		return fmt.Errorf("incomplete terragrunt command")
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
		return fmt.Errorf("'%s' is not a valid terragrunt subcommand. Use 'terragrunt --help' to see valid commands.", subcommand)
	}

	return nil
}

// checkDangerousOps checks for dangerous terragrunt operations.
func (v *Validator) checkDangerousOps(command string) string {
	for _, danger := range v.dangerousOps {
		if strings.Contains(command, danger) {
			if strings.Contains(danger, "-all") {
				return fmt.Sprintf("EXTREMELY DANGEROUS: 'terragrunt %s' will affect ALL modules in the dependency tree. This can destroy entire environments!", danger)
			}
			return fmt.Sprintf("Dangerous operation detected: 'terragrunt %s'. This will modify or destroy infrastructure.", danger)
		}
	}
	return ""
}

// checkCommonMistakes checks for common terragrunt command mistakes.
func (v *Validator) checkCommonMistakes(command string) error {
	// Check for using terraform flags that don't work with terragrunt
	if strings.Contains(command, "run-all") && strings.Contains(command, "-target") {
		return fmt.Errorf("-target doesn't work with run-all. Use it with individual apply/plan commands.")
	}

	// Check for missing terragrunt-specific prefix
	if matched, _ := regexp.MatchString(`--working-dir\s`, command); matched {
		return fmt.Errorf("did you mean --terragrunt-working-dir?")
	}

	if matched, _ := regexp.MatchString(`--config\s`, command); matched {
		return fmt.Errorf("did you mean --terragrunt-config?")
	}

	// Warn about run-all without proper flags
	if strings.Contains(command, "run-all apply") && !strings.Contains(command, "--terragrunt-non-interactive") {
		return fmt.Errorf("⚠️  run-all apply without --terragrunt-non-interactive will prompt for each module. Consider adding this flag.")
	}

	return nil
}

