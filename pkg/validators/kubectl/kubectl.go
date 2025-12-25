package kubectl

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Validator implements the Validator interface for kubectl commands.
type Validator struct {
	validSubcommands []string
	dangerousOps     []string
}

// NewValidator creates a new kubectl validator.
func NewValidator() *Validator {
	return &Validator{
		validSubcommands: []string{
			"get", "describe", "logs", "exec", "apply", "create", "delete",
			"edit", "replace", "patch", "scale", "rollout", "expose",
			"port-forward", "proxy", "cp", "attach", "run", "explain",
			"drain", "cordon", "uncordon", "taint", "label", "annotate",
			"config", "cluster-info", "top", "api-resources", "api-versions",
		},
		dangerousOps: []string{
			"delete", "drain", "delete --all", "delete namespace",
		},
	}
}

// CanValidate returns true if this validator can handle the command.
func (v *Validator) CanValidate(command string) bool {
	return strings.HasPrefix(command, "kubectl") || strings.HasPrefix(command, "k ")
}

// Validate checks if a kubectl command is valid.
func (v *Validator) Validate(command string) error {
	// Handle alias 'k' for 'kubectl'
	if strings.HasPrefix(command, "k ") {
		command = "kubectl" + command[1:]
	}
	
	if !strings.HasPrefix(command, "kubectl") {
		return nil // Not a kubectl command
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
		// Don't block, just warn
		return fmt.Errorf("⚠️  Warning: %s", warning)
	}

	// Check for YAML syntax if apply/create with -f
	if strings.Contains(command, " apply ") || strings.Contains(command, " create ") {
		if strings.Contains(command, " -f ") || strings.Contains(command, "--filename") {
			// We can't validate file content from command string alone
			// This would require reading the file, which we'll skip for now
		}
	}

	return nil
}

// checkHallucinatedFlags checks for commonly hallucinated kubectl flags.
func (v *Validator) checkHallucinatedFlags(command string) error {
	hallucinatedPatterns := map[string]string{
		`--sort[\s=]`:          "kubectl get does not have a --sort flag. Use --sort-by instead.",
		`--filter[\s=]`:        "kubectl does not have a --filter flag. Use -l (label selector) or field selectors.",
		`--format[\s=]`:        "kubectl does not have a --format flag. Use -o or --output instead.",
		`--limit[\s=]`:         "kubectl does not have a --limit flag. Use --field-selector or pipe to head.",
		`--where[\s=]`:         "kubectl does not have a --where flag. Use --field-selector instead.",
		`--order-by[\s=]`:      "kubectl does not have an --order-by flag. Use --sort-by instead.",
		`get pods --memory`:    "kubectl get pods does not show memory directly. Use 'kubectl top pods' instead.",
		`get pods --cpu`:       "kubectl get pods does not show CPU directly. Use 'kubectl top pods' instead.",
		`logs --grep`:          "kubectl logs does not have a --grep flag. Pipe to grep instead.",
		`apply --force-delete`: "kubectl apply does not have --force-delete. Use 'kubectl delete --force' separately.",
	}

	for pattern, message := range hallucinatedPatterns {
		if matched, _ := regexp.MatchString(pattern, command); matched {
			return fmt.Errorf("%s", message)
		}
	}

	return nil
}

// checkSubcommand validates the kubectl subcommand.
func (v *Validator) checkSubcommand(command string) error {
	parts := strings.Fields(command)
	if len(parts) < 2 {
		return fmt.Errorf("incomplete kubectl command")
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
		return fmt.Errorf("'%s' is not a valid kubectl subcommand. Use 'kubectl --help' to see valid commands.", subcommand)
	}

	return nil
}

// checkDangerousOps checks for dangerous kubectl operations.
func (v *Validator) checkDangerousOps(command string) string {
	for _, danger := range v.dangerousOps {
		if strings.Contains(command, danger) {
			return fmt.Sprintf("Dangerous operation detected: '%s'. Ensure this is intentional.", danger)
		}
	}
	return ""
}

// ValidateYAML validates Kubernetes YAML content.
func (v *Validator) ValidateYAML(content string) error {
	var data interface{}
	if err := yaml.Unmarshal([]byte(content), &data); err != nil {
		return fmt.Errorf("invalid YAML syntax: %w", err)
	}

	// Basic structure validation
	if m, ok := data.(map[string]interface{}); ok {
		// Check for required Kubernetes fields
		if _, hasAPIVersion := m["apiVersion"]; !hasAPIVersion {
			return fmt.Errorf("missing required field: apiVersion")
		}
		if _, hasKind := m["kind"]; !hasKind {
			return fmt.Errorf("missing required field: kind")
		}
		if _, hasMetadata := m["metadata"]; !hasMetadata {
			return fmt.Errorf("missing required field: metadata")
		}
	}

	return nil
}

