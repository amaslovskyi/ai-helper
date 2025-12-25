package ui

import (
	"fmt"
	"os"
	"strings"
)

// Color codes for terminal output
const (
	Reset = "\033[0m"
	Bold  = "\033[1m"
	Dim   = "\033[2m"

	// Foreground colors
	Red     = "\033[0;31m"
	Green   = "\033[0;32m"
	Yellow  = "\033[0;33m"
	Blue    = "\033[0;34m"
	Magenta = "\033[0;35m"
	Cyan    = "\033[0;36m"
	White   = "\033[0;37m"

	// Bold colors
	RedBold     = "\033[1;31m"
	GreenBold   = "\033[1;32m"
	YellowBold  = "\033[1;33m"
	BlueBold    = "\033[1;34m"
	MagentaBold = "\033[1;35m"
	CyanBold    = "\033[1;36m"

	// Background colors
	BgRed    = "\033[41m"
	BgYellow = "\033[43m"
	BgGreen  = "\033[42m"
)

// ColorsEnabled checks if colors should be used
var ColorsEnabled = os.Getenv("NO_COLOR") == "" && isTerminal()

func isTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// Colorize wraps text with color codes
func Colorize(color, text string) string {
	if !ColorsEnabled {
		return text
	}
	return color + text + Reset
}

// FormatAIResponse formats AI output with colors
func FormatAIResponse(response string) string {
	if !ColorsEnabled {
		return response
	}

	lines := strings.Split(response, "\n")
	var formatted []string

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "✓"):
			// Green for corrected commands
			formatted = append(formatted, Colorize(GreenBold, line))
		case strings.HasPrefix(line, "Root:"):
			// Cyan for root cause
			formatted = append(formatted, Colorize(Cyan, line))
		case strings.HasPrefix(line, "Tip:"), strings.HasPrefix(line, "Check:"),
			strings.HasPrefix(line, "Fix:"), strings.HasPrefix(line, "Note:"):
			// Yellow for tips
			formatted = append(formatted, Colorize(Yellow, line))
		case strings.HasPrefix(line, "Error:"):
			// Red for errors
			formatted = append(formatted, Colorize(Red, line))
		case strings.HasPrefix(line, "💾"):
			// Magenta for cached responses
			formatted = append(formatted, Colorize(MagentaBold, line))
		default:
			formatted = append(formatted, line)
		}
	}

	return strings.Join(formatted, "\n")
}

// PrintSuccess prints a success message
func PrintSuccess(message string) {
	fmt.Println(Colorize(Green, "✅ "+message))
}

// PrintError prints an error message
func PrintError(message string) {
	fmt.Println(Colorize(Red, "❌ "+message))
}

// PrintWarning prints a warning message
func PrintWarning(message string) {
	fmt.Println(Colorize(Yellow, "⚠️  "+message))
}

// PrintInfo prints an info message
func PrintInfo(message string) {
	fmt.Println(Colorize(Cyan, "ℹ️  "+message))
}

// PrintDanger prints a danger warning with red background
func PrintDanger(message string) {
	fmt.Println(Colorize(RedBold+BgRed, " DANGER ") + " " + Colorize(RedBold, message))
}

