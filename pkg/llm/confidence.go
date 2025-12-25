package llm

import (
	"strings"
)

// ConfidenceLevel represents the confidence in an AI suggestion.
type ConfidenceLevel string

const (
	HighConfidence   ConfidenceLevel = "High"   // 90%+ confidence
	MediumConfidence ConfidenceLevel = "Medium" // 70-90% confidence
	LowConfidence    ConfidenceLevel = "Low"    // <70% confidence
)

// CalculateConfidence determines confidence level based on multiple factors.
func CalculateConfidence(response *Response, validationError error, commandComplexity int) (ConfidenceLevel, int) {
	score := 100 // Start with 100%

	// Factor 1: Validation result (40% weight)
	if validationError != nil {
		errorMsg := validationError.Error()
		if strings.Contains(errorMsg, "🚨 BLOCKED") {
			score -= 50 // Critical validation failure
		} else if strings.Contains(errorMsg, "⚠️") {
			score -= 20 // Warning
		} else {
			score -= 30 // General validation failure
		}
	}

	// Factor 2: Command structure quality (30% weight)
	if response.Suggestion == "" {
		score -= 50
	} else {
		// Check if command looks well-formed
		if !strings.Contains(response.Suggestion, " ") {
			score -= 10 // Single-word command (might be incomplete)
		}
		if strings.Contains(response.Suggestion, "...") ||
			strings.Contains(response.Suggestion, "<") ||
			strings.Contains(response.Suggestion, "[") {
			score -= 20 // Contains placeholder-like patterns
		}
	}

	// Factor 3: Root cause presence (15% weight)
	if response.RootCause == "" {
		score -= 15
	} else if len(response.RootCause) < 10 {
		score -= 10 // Very short explanation
	}

	// Factor 4: Command complexity (15% weight)
	// More complex commands = lower confidence
	if commandComplexity > 5 {
		score -= 15
	} else if commandComplexity > 3 {
		score -= 10
	}

	// Ensure score is within bounds
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	// Determine confidence level
	var level ConfidenceLevel
	if score >= 90 {
		level = HighConfidence
	} else if score >= 70 {
		level = MediumConfidence
	} else {
		level = LowConfidence
	}

	return level, score
}

// CalculateCommandComplexity estimates how complex a command is.
func CalculateCommandComplexity(command string) int {
	complexity := 0

	// Count pipes
	complexity += strings.Count(command, "|")

	// Count redirects
	complexity += strings.Count(command, ">")
	complexity += strings.Count(command, "<")

	// Count flags
	complexity += strings.Count(command, " -")
	complexity += strings.Count(command, " --")

	// Count logical operators
	complexity += strings.Count(command, "&&")
	complexity += strings.Count(command, "||")

	// Count subshells
	complexity += strings.Count(command, "$(")
	complexity += strings.Count(command, "`")

	return complexity
}

// GetConfidenceEmoji returns an emoji representing the confidence level.
func GetConfidenceEmoji(level ConfidenceLevel) string {
	switch level {
	case HighConfidence:
		return "✅"
	case MediumConfidence:
		return "⚠️"
	case LowConfidence:
		return "❓"
	default:
		return "❓"
	}
}

// GetConfidenceColor returns an ANSI color code for the confidence level.
func GetConfidenceColor(level ConfidenceLevel) string {
	switch level {
	case HighConfidence:
		return "\033[0;32m" // Green
	case MediumConfidence:
		return "\033[0;33m" // Yellow
	case LowConfidence:
		return "\033[0;31m" // Red
	default:
		return "\033[0m" // Reset
	}
}

