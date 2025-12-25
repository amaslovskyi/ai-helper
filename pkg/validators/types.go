package validators

import "fmt"

// Validator interface for command validation
type Validator interface {
	// Validate checks if a command is valid
	Validate(command string) error
	
	// CanValidate returns true if this validator can handle the command
	CanValidate(command string) bool
}

// ValidationError represents a validation failure
type ValidationError struct {
	Command string
	Reason  string
	Hint    string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid command: %s (hint: %s)", e.Reason, e.Hint)
}

// NewValidationError creates a new validation error
func NewValidationError(command, reason, hint string) *ValidationError {
	return &ValidationError{
		Command: command,
		Reason:  reason,
		Hint:    hint,
	}
}

