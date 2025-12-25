package interactive

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/amaslovskyi/ai-helper/pkg/ui"
)

// MenuOption represents a single menu choice
type MenuOption struct {
	Key         string
	Label       string
	Description string
	Action      string // For internal tracking
}

// MenuResult contains the user's choice and any additional data
type MenuResult struct {
	Action   string
	Input    string // Additional input if needed
	Canceled bool
}

// Menu displays an interactive menu and returns the user's choice
type Menu struct {
	Title   string
	Options []MenuOption
	reader  *bufio.Reader
}

// NewMenu creates a new interactive menu
func NewMenu(title string) *Menu {
	return &Menu{
		Title:   title,
		Options: []MenuOption{},
		reader:  bufio.NewReader(os.Stdin),
	}
}

// AddOption adds an option to the menu
func (m *Menu) AddOption(key, label, description, action string) {
	m.Options = append(m.Options, MenuOption{
		Key:         key,
		Label:       label,
		Description: description,
		Action:      action,
	})
}

// Show displays the menu and waits for user input
func (m *Menu) Show() *MenuResult {
	// Display title
	fmt.Println()
	fmt.Println(ui.Colorize(ui.CyanBold, "🤖 "+m.Title))
	fmt.Println()

	// Display options
	for _, opt := range m.Options {
		fmt.Printf("  %s %s %s\n",
			ui.Colorize(ui.Yellow, "["+opt.Key+"]"),
			ui.Colorize(ui.White, opt.Label),
			ui.Colorize(ui.Dim, "- "+opt.Description))
	}

	fmt.Println()
	fmt.Print(ui.Colorize(ui.Green, "Your choice: "))

	// Read user input
	input, err := m.reader.ReadString('\n')
	if err != nil {
		return &MenuResult{Canceled: true}
	}

	// Clean input
	choice := strings.TrimSpace(strings.ToLower(input))

	// Handle empty input (treat as skip)
	if choice == "" {
		return &MenuResult{Action: "skip", Canceled: false}
	}

	// Find matching option
	for _, opt := range m.Options {
		if strings.ToLower(opt.Key) == choice {
			return &MenuResult{
				Action:   opt.Action,
				Input:    choice,
				Canceled: false,
			}
		}
	}

	// Invalid choice
	ui.PrintWarning(fmt.Sprintf("Invalid choice: %s", choice))
	return &MenuResult{Action: "skip", Canceled: false}
}

// ShowErrorMenu displays the standard error handling menu
func ShowErrorMenu(command string, errorMsg string) *MenuResult {
	menu := NewMenu("Command failed. What would you like to do?")
	
	// Standard options
	menu.AddOption("1", "Get AI suggestion", "Let AI analyze and suggest a fix", "ai")
	menu.AddOption("2", "Show manual", "Display manual page for this command", "manual")
	menu.AddOption("3", "Skip", "Continue without fixing", "skip")
	menu.AddOption("4", "Disable AI for session", "Turn off AI until terminal restart", "disable")
	
	return menu.Show()
}

// ShowProactiveMenu displays menu for proactive mode suggestions
func ShowProactiveMenu(suggestion string) *MenuResult {
	menu := NewMenu("AI generated this command:")
	
	fmt.Println()
	fmt.Println(ui.Colorize(ui.GreenBold, "  ✓ "+suggestion))
	fmt.Println()
	
	menu.AddOption("y", "Execute", "Run this command", "execute")
	menu.AddOption("c", "Copy", "Copy to clipboard (manual execution)", "copy")
	menu.AddOption("m", "Modify", "Edit before executing", "modify")
	menu.AddOption("n", "Skip", "Don't execute, continue", "skip")
	
	return menu.Show()
}

// ShowConfirmation displays a simple yes/no confirmation
func ShowConfirmation(message string) bool {
	fmt.Printf("%s %s ",
		ui.Colorize(ui.Yellow, message),
		ui.Colorize(ui.Dim, "(y/n):"))
	
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	
	response := strings.TrimSpace(strings.ToLower(input))
	return response == "y" || response == "yes"
}

// ShowDangerousCommandWarning displays warning for dangerous commands
func ShowDangerousCommandWarning(command string, warning string) bool {
	fmt.Println()
	fmt.Println(ui.Colorize(ui.RedBold, "⚠️  DANGEROUS COMMAND DETECTED"))
	fmt.Println()
	fmt.Println(ui.Colorize(ui.Yellow, "Command: ")+command)
	fmt.Println(ui.Colorize(ui.Red, "Warning: ")+warning)
	fmt.Println()
	
	return ShowConfirmation("Are you sure you want to proceed?")
}

// Prompt displays a prompt and returns user input
func Prompt(message string) string {
	fmt.Print(ui.Colorize(ui.Cyan, message+" "))
	
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	
	return strings.TrimSpace(input)
}

