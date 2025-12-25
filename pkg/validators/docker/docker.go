package docker

import (
	"strings"

	"github.com/yourusername/ai-helper/pkg/validators"
)

// Validator validates Docker commands
type Validator struct {
	validPsFlags     map[string]bool
	validStatsFlags  map[string]bool
	validRunFlags    map[string]bool
}

// NewValidator creates a new Docker validator
func NewValidator() *Validator {
	return &Validator{
		validPsFlags: map[string]bool{
			"-a":         true,
			"--all":      true,
			"-f":         true,
			"--filter":   true,
			"--format":   true,
			"-n":         true,
			"--last":     true,
			"-l":         true,
			"--latest":   true,
			"--no-trunc": true,
			"-q":         true,
			"--quiet":    true,
			"-s":         true,
			"--size":     true,
		},
		validStatsFlags: map[string]bool{
			"-a":           true,
			"--all":        true,
			"--format":     true,
			"--no-stream":  true,
			"--no-trunc":   true,
		},
		validRunFlags: map[string]bool{
			"-d":            true,
			"--detach":      true,
			"-e":            true,
			"--env":         true,
			"-p":            true,
			"--publish":     true,
			"-v":            true,
			"--volume":      true,
			"--name":        true,
			"--rm":          true,
			"-i":            true,
			"--interactive": true,
			"-t":            true,
			"--tty":         true,
			"--network":     true,
			"--restart":     true,
			"-w":            true,
			"--workdir":     true,
			"-u":            true,
			"--user":        true,
		},
	}
}

// CanValidate returns true if this is a Docker command
func (v *Validator) CanValidate(command string) bool {
	return strings.HasPrefix(strings.TrimSpace(command), "docker ")
}

// Validate validates a Docker command
func (v *Validator) Validate(command string) error {
	command = strings.TrimSpace(command)
	parts := strings.Fields(command)

	if len(parts) < 2 {
		return validators.NewValidationError(
			command,
			"incomplete docker command",
			"docker commands need a subcommand (e.g., 'docker ps')",
		)
	}

	subCommand := parts[1]

	switch subCommand {
	case "ps":
		return v.validatePs(parts[2:])
	case "stats":
		return v.validateStats(parts[2:])
	case "run":
		return v.validateRun(parts[2:])
	}

	// For other subcommands, just do basic validation
	return nil
}

// validatePs validates 'docker ps' command
func (v *Validator) validatePs(args []string) error {
	// Check for common hallucinations
	for _, arg := range args {
		if strings.HasPrefix(arg, "--sort") {
			return validators.NewValidationError(
				"docker ps",
				"docker ps does not have a --sort flag",
				"use 'docker stats --no-stream | sort' or format with --format and pipe to sort",
			)
		}
		
		// Check if it's a flag (starts with -)
		if strings.HasPrefix(arg, "-") {
			flag := strings.Split(arg, "=")[0] // Handle --flag=value
			if !v.validPsFlags[flag] {
				return validators.NewValidationError(
					"docker ps "+arg,
					"invalid flag for docker ps",
					"run 'docker ps --help' for valid flags",
				)
			}
		}
	}

	return nil
}

// validateStats validates 'docker stats' command
func (v *Validator) validateStats(args []string) error {
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			flag := strings.Split(arg, "=")[0]
			if !v.validStatsFlags[flag] {
				return validators.NewValidationError(
					"docker stats "+arg,
					"invalid flag for docker stats",
					"run 'docker stats --help' for valid flags",
				)
			}
		}
	}

	return nil
}

// validateRun validates 'docker run' command
func (v *Validator) validateRun(args []string) error {
	// Check for port conflicts and other common issues
	for i, arg := range args {
		if strings.HasPrefix(arg, "-") {
			flag := strings.Split(arg, "=")[0]
			
			// Skip known flags
			if v.validRunFlags[flag] {
				continue
			}
			
			// Check for common typos
			if strings.HasPrefix(flag, "--port") {
				return validators.NewValidationError(
					"docker run "+flag,
					"invalid flag --port, did you mean -p or --publish?",
					"use -p host:container or --publish host:container",
				)
			}
		}

		// Check for port format
		if (args[i-1] == "-p" || args[i-1] == "--publish") && i > 0 {
			if !strings.Contains(arg, ":") {
				return validators.NewValidationError(
					"docker run -p "+arg,
					"port mapping must be in format host:container",
					"example: -p 8080:80",
				)
			}
		}
	}

	return nil
}

