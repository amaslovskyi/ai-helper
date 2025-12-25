package security

import (
	"fmt"
	"strings"
)

// DangerousPattern represents a dangerous command pattern
type DangerousPattern struct {
	Pattern     string
	Description string
	Severity    string
}

// Scanner scans commands for dangerous patterns
type Scanner struct {
	patterns []DangerousPattern
}

// NewScanner creates a new security scanner
func NewScanner() *Scanner {
	return &Scanner{
		patterns: []DangerousPattern{
			{"rm -rf /", "recursive deletion of root", "CRITICAL"},
			{"rm -rf *", "recursive deletion of all files", "CRITICAL"},
			{"rm -rf ~", "recursive deletion of home directory", "CRITICAL"},
			{"rm -rf $HOME", "recursive deletion of home directory", "CRITICAL"},
			{"> /dev/sda", "writing to disk device", "CRITICAL"},
			{"dd if=/dev/zero", "disk overwrite", "CRITICAL"},
			{"mkfs.", "filesystem formatting", "CRITICAL"},
			{"DROP DATABASE", "database deletion", "CRITICAL"},
			{"DROP TABLE", "table deletion", "CRITICAL"},
			{"TRUNCATE", "data truncation", "HIGH"},
			{"chmod -R 777", "insecure permissions", "HIGH"},
			{"chmod 777", "insecure permissions", "MEDIUM"},
			{"chown -R", "ownership change (use with caution)", "MEDIUM"},
			{":(){ :|:& };:", "fork bomb", "CRITICAL"},
			{"--no-preserve-root", "disables root protection", "CRITICAL"},
			{"mv .* /dev/null", "move to null device", "HIGH"},
			{"kubectl delete", "kubernetes resource deletion", "MEDIUM"},
			{"terraform destroy", "infrastructure destruction", "MEDIUM"},
		},
	}
}

// Scan checks if a command contains dangerous patterns
func (s *Scanner) Scan(command string) (*DangerResult, error) {
	cmdLower := strings.ToLower(command)

	for _, pattern := range s.patterns {
		if strings.Contains(cmdLower, strings.ToLower(pattern.Pattern)) {
			return &DangerResult{
				IsDangerous: true,
				Pattern:     pattern.Pattern,
				Description: pattern.Description,
				Severity:    pattern.Severity,
				Command:     command,
			}, nil
		}
	}

	return &DangerResult{IsDangerous: false}, nil
}

// DangerResult represents the result of a security scan
type DangerResult struct {
	IsDangerous bool
	Pattern     string
	Description string
	Severity    string
	Command     string
}

// Warning returns a formatted warning message
func (r *DangerResult) Warning() string {
	if !r.IsDangerous {
		return ""
	}

	return fmt.Sprintf(`🚨 DANGER: Command contains potentially destructive pattern: %s
⚠️  This could cause data loss or system damage! (%s)
📋 Command: %s

If you're ABSOLUTELY SURE this is safe, you can:
  1. Review the command carefully
  2. Test in a safe environment first
  3. Execute manually after verification`,
		r.Pattern,
		r.Description,
		r.Command,
	)
}

