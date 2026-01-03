package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/amaslovskyi/ai-helper/pkg/cache"
	"github.com/amaslovskyi/ai-helper/pkg/config"
	"github.com/amaslovskyi/ai-helper/pkg/interactive"
	"github.com/amaslovskyi/ai-helper/pkg/llm"
	"github.com/amaslovskyi/ai-helper/pkg/security"
	"github.com/amaslovskyi/ai-helper/pkg/ui"
	"github.com/amaslovskyi/ai-helper/pkg/validators"
	"github.com/amaslovskyi/ai-helper/pkg/validators/ansible"
	"github.com/amaslovskyi/ai-helper/pkg/validators/argocd"
	"github.com/amaslovskyi/ai-helper/pkg/validators/docker"
	"github.com/amaslovskyi/ai-helper/pkg/validators/git"
	"github.com/amaslovskyi/ai-helper/pkg/validators/helm"
	"github.com/amaslovskyi/ai-helper/pkg/validators/kubectl"
	"github.com/amaslovskyi/ai-helper/pkg/validators/terraform"
	"github.com/amaslovskyi/ai-helper/pkg/validators/terragrunt"
)

const version = "2.3.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Initialize components
	homeDir, err := os.UserHomeDir()
	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to get home directory: %v", err))
		os.Exit(1)
	}

	aiDir := filepath.Join(homeDir, ".ai")
	cacheFile := filepath.Join(aiDir, "cache.json")
	configFile := filepath.Join(aiDir, "config.json")

	// Load configuration
	cfg, err := config.Load(configFile)
	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to load config: %v", err))
		os.Exit(1)
	}

	// Create cache
	cacheStore, err := cache.NewCache(cacheFile)
	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to initialize cache: %v", err))
		os.Exit(1)
	}

	// Create LLM client based on provider configuration
	var client llm.Client
	switch cfg.Provider {
	case config.ProviderOpenCode:
		client = llm.NewOpenCodeClient(cfg.PreferredModel)
	default:
		client = llm.NewOllamaClient("")
	}

	// Create security scanner
	scanner := security.NewScanner()

	// Create validators (order matters: more specific first)
	validatorsList := []validators.Validator{
		kubectl.NewValidator(),    // k, kubectl
		terraform.NewValidator(),  // tf, terraform
		terragrunt.NewValidator(), // tg, terragrunt
		helm.NewValidator(),       // h, helm
		git.NewValidator(),        // git, gco, gcb, gp, etc.
		docker.NewValidator(),     // docker, d, dc
		ansible.NewValidator(),    // ansible, ansible-playbook
		argocd.NewValidator(),     // argocd
	}

	// Parse command
	cmd := os.Args[1]

	switch cmd {
	case "analyze":
		handleAnalyze(client, cacheStore, scanner, validatorsList, cfg)
	case "proactive", "ask":
		handleProactive(client, scanner, validatorsList, cfg)
	case "version", "-v", "--version", "-V":
		// Support common version flag conventions
		fmt.Printf("AI Terminal Helper v%s (Go)\n", version)
	case "cache-stats":
		handleCacheStats(cacheStore)
	case "cache-clear":
		handleCacheClear(cacheStore)
	case "config-show":
		handleConfigShow(cfg)
	case "config-set":
		handleConfigSet(cfg, configFile)
	case "config-reset":
		handleConfigReset(configFile)
	case "-h", "--help", "help":
		// Support common help flag conventions
		printUsage()
		os.Exit(0)
	default:
		ui.PrintError(fmt.Sprintf("Unknown command: %s", cmd))
		printUsage()
		os.Exit(1)
	}
}

func handleAnalyze(client llm.Client, cacheStore *cache.Cache, scanner *security.Scanner, validators []validators.Validator, cfg *config.Config) {
	if len(os.Args) < 4 {
		ui.PrintError("Usage: ai-helper analyze <command> <exit_code> [error_output]")
		os.Exit(1)
	}

	command := os.Args[2]
	exitCode := 1 // default
	fmt.Sscanf(os.Args[3], "%d", &exitCode)
	errorOutput := ""
	if len(os.Args) > 4 {
		errorOutput = os.Args[4]
	}

	// Extract tool name for mode checking
	toolName := extractToolName(command)

	// Check if AI is disabled
	if !cfg.IsEnabled(toolName) {
		os.Exit(exitCode) // Just exit with the original error code
	}

	// Try cache first
	if cachedResp, ok := cacheStore.Get(command, errorOutput); ok {
		fmt.Println(ui.Colorize(ui.MagentaBold, "💾 [Cached]"))
		printResponse(cachedResp)
		return
	}

	// Check activation mode
	if cfg.ShouldShowMenu(toolName) {
		// Show interactive menu
		result := interactive.ShowErrorMenu(command, errorOutput)

		switch result.Action {
		case "ai":
			// Continue to AI suggestion
		case "manual":
			fmt.Println(ui.Colorize(ui.Cyan, "📖 Tip: Use 'man "+toolName+"' for documentation"))
			return
		case "skip":
			return
		case "disable":
			cfg.SessionDisabled = true
			ui.PrintInfo("AI disabled for this session. Restart terminal to re-enable.")
			return
		default:
			return
		}
	} else if !cfg.ShouldTriggerAI(toolName) {
		// Manual mode - don't trigger AI automatically
		return
	}

	// Query AI
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cwd, _ := os.Getwd()
	req := llm.Request{
		Command:   command,
		Error:     errorOutput,
		ExitCode:  exitCode,
		Directory: cwd,
		Mode:      llm.ModeReactive,
	}

	resp, err := client.Query(ctx, req)
	if err != nil {
		ui.PrintError(fmt.Sprintf("AI query failed: %v", err))
		os.Exit(1)
	}

	// Validate the suggested command
	validationErr := validateCommand(resp.Suggestion, validators)
	if validationErr != nil {
		ui.PrintWarning(fmt.Sprintf("AI suggestion validation failed: %v", validationErr))
		ui.PrintInfo("Querying AI again with validation context...")

		// Query again with validation error context
		req.Context = fmt.Sprintf("Previous suggestion '%s' was invalid: %v", resp.Suggestion, validationErr)
		resp, err = client.Query(ctx, req)
		if err != nil {
			ui.PrintError(fmt.Sprintf("AI re-query failed: %v", err))
			os.Exit(1)
		}

		// Validate again
		validationErr = validateCommand(resp.Suggestion, validators)
		if validationErr != nil {
			ui.PrintError(fmt.Sprintf("AI still suggesting invalid command: %v", validationErr))
			ui.PrintInfo("Suggestion: " + resp.Suggestion)
			os.Exit(1)
		}
	}

	// Calculate confidence
	complexity := llm.CalculateCommandComplexity(command)
	confLevel, confScore := llm.CalculateConfidence(resp, validationErr, complexity)

	// Security scan
	dangerResult, err := scanner.Scan(resp.Suggestion)
	if err != nil {
		ui.PrintError(fmt.Sprintf("Security scan failed: %v", err))
		os.Exit(1)
	}

	if dangerResult.IsDangerous {
		ui.PrintDanger(dangerResult.Warning())
		os.Exit(1)
	}

	// Cache the response
	if err := cacheStore.Set(command, errorOutput, resp); err != nil {
		// Non-fatal, just log
		ui.PrintWarning(fmt.Sprintf("Failed to cache response: %v", err))
	}

	// Print response with confidence
	printResponseWithConfidence(resp, confLevel, confScore)
}

func handleProactive(client llm.Client, scanner *security.Scanner, validators []validators.Validator, cfg *config.Config) {
	if len(os.Args) < 3 {
		ui.PrintError("Usage: ai-helper proactive <query>")
		os.Exit(1)
	}

	query := strings.Join(os.Args[2:], " ")

	// Check if AI is enabled (use empty string for tool since this is general query)
	if !cfg.IsEnabled("") {
		os.Exit(0)
	}

	fmt.Println(ui.Colorize(ui.CyanBold, "🤖 Generating command for: ") + ui.Colorize(ui.Yellow, query))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cwd, _ := os.Getwd()
	req := llm.Request{
		Command:   query,
		Directory: cwd,
		Mode:      llm.ModeProactive,
	}

	resp, err := client.Query(ctx, req)
	if err != nil {
		ui.PrintError(fmt.Sprintf("AI query failed: %v", err))
		os.Exit(1)
	}

	// Validate the suggested command
	validationErr := validateCommand(resp.Suggestion, validators)
	if validationErr != nil {
		ui.PrintWarning(fmt.Sprintf("Validation failed: %v", validationErr))
		// In proactive mode, still show the suggestion but with warning
	}

	// Calculate confidence
	complexity := llm.CalculateCommandComplexity(query)
	confLevel, confScore := llm.CalculateConfidence(resp, validationErr, complexity)

	// Security scan
	dangerResult, err := scanner.Scan(resp.Suggestion)
	if err != nil {
		ui.PrintError(fmt.Sprintf("Security scan failed: %v", err))
		os.Exit(1)
	}

	if dangerResult.IsDangerous {
		ui.PrintDanger(dangerResult.Warning())
		os.Exit(1)
	}

	printResponseWithConfidence(resp, confLevel, confScore)
}

func handleCacheStats(cacheStore *cache.Cache) {
	stats := cacheStore.Stats()
	fmt.Println(ui.Colorize(ui.CyanBold, "📊 Cache Statistics:"))
	fmt.Printf("  %s %s%d%s\n",
		ui.Colorize(ui.Yellow, "Total patterns:"),
		ui.Colorize(ui.Green, ""),
		stats["total_entries"],
		ui.Colorize(ui.Reset, ""))
	fmt.Printf("  %s %s%d%s\n",
		ui.Colorize(ui.Yellow, "Total hits:"),
		ui.Colorize(ui.Green, ""),
		stats["total_hits"],
		ui.Colorize(ui.Reset, ""))
	fmt.Printf("  %s %s%s%s\n",
		ui.Colorize(ui.Yellow, "Cache file:"),
		ui.Colorize(ui.Blue, ""),
		stats["cache_file"],
		ui.Colorize(ui.Reset, ""))
}

func handleCacheClear(cacheStore *cache.Cache) {
	if err := cacheStore.Clear(); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to clear cache: %v", err))
		os.Exit(1)
	}
	ui.PrintSuccess("Cache cleared")
}

func validateCommand(command string, validators []validators.Validator) error {
	for _, v := range validators {
		if v.CanValidate(command) {
			return v.Validate(command)
		}
	}
	return nil
}

func printResponse(resp *llm.Response) {
	if resp.Suggestion != "" {
		fmt.Println(ui.Colorize(ui.GreenBold, "✓ "+resp.Suggestion))
	}
	if resp.RootCause != "" {
		fmt.Println(ui.Colorize(ui.Cyan, "Root: "+resp.RootCause))
	}
	if resp.Tip != "" {
		fmt.Println(ui.Colorize(ui.Yellow, "Tip: "+resp.Tip))
	}
}

func printResponseWithConfidence(resp *llm.Response, level llm.ConfidenceLevel, score int) {
	if resp.Suggestion != "" {
		fmt.Println(ui.Colorize(ui.GreenBold, "✓ "+resp.Suggestion))
	}
	if resp.RootCause != "" {
		fmt.Println(ui.Colorize(ui.Cyan, "Root: "+resp.RootCause))
	}
	if resp.Tip != "" {
		fmt.Println(ui.Colorize(ui.Yellow, "Tip: "+resp.Tip))
	}
	// Print confidence
	emoji := llm.GetConfidenceEmoji(level)
	color := llm.GetConfidenceColor(level)
	fmt.Printf("%sConfidence: %s %s (%d%%)%s\n",
		color, emoji, level, score, ui.Reset)
}

// extractToolName extracts the base tool name from a command
func extractToolName(command string) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

// handleConfigShow displays current configuration
func handleConfigShow(cfg *config.Config) {
	fmt.Println(ui.Colorize(ui.CyanBold, "⚙️  Configuration:"))
	fmt.Printf("  %s %s\n",
		ui.Colorize(ui.Yellow, "Activation Mode:"),
		ui.Colorize(ui.Green, string(cfg.ActivationMode)))
	fmt.Printf("  %s %v\n",
		ui.Colorize(ui.Yellow, "Auto Execute Safe:"),
		cfg.AutoExecuteSafe)
	fmt.Printf("  %s %v\n",
		ui.Colorize(ui.Yellow, "Show Confidence:"),
		cfg.ShowConfidence)
	fmt.Printf("  %s %s\n",
		ui.Colorize(ui.Yellow, "Provider:"),
		ui.Colorize(ui.Green, string(cfg.Provider)))
	if cfg.PreferredModel != "" {
		fmt.Printf("  %s %s\n",
			ui.Colorize(ui.Yellow, "Preferred Model:"),
			cfg.PreferredModel)
	}
	if len(cfg.ToolSpecificModes) > 0 {
		fmt.Println(ui.Colorize(ui.Yellow, "  Tool-Specific Modes:"))
		for tool, mode := range cfg.ToolSpecificModes {
			fmt.Printf("    %s: %s\n", tool, mode)
		}
	}
}

// handleConfigSet updates configuration
func handleConfigSet(cfg *config.Config, configFile string) {
	if len(os.Args) < 4 {
		fmt.Println(ui.Colorize(ui.Red, "Usage: ai-helper config-set <key> <value>"))
		fmt.Println()
		fmt.Println("Available keys:")
		fmt.Println("  mode <auto|interactive|manual|disabled> - Set activation mode")
		fmt.Println("  tool-mode <tool> <mode> - Set tool-specific mode")
		fmt.Println("  confidence <true|false> - Show/hide confidence scores")
		fmt.Println("  provider <ollama|opencode> - Set LLM provider")
		fmt.Println("  model <model-name> - Set preferred model")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  ai-helper config-set mode interactive")
		fmt.Println("  ai-helper config-set tool-mode kubectl interactive")
		fmt.Println("  ai-helper config-set confidence false")
		fmt.Println("  ai-helper config-set provider opencode")
		fmt.Println("  ai-helper config-set model anthropic/claude-sonnet-4-20250514")
		os.Exit(1)
	}

	key := os.Args[2]
	value := os.Args[3]

	switch key {
	case "mode":
		if !config.ValidateMode(value) {
			ui.PrintError("Invalid mode. Use: auto, interactive, manual, or disabled")
			os.Exit(1)
		}
		cfg.ActivationMode = config.ActivationMode(value)
		ui.PrintSuccess(fmt.Sprintf("Activation mode set to: %s", value))

	case "tool-mode":
		if len(os.Args) < 5 {
			ui.PrintError("Usage: ai-helper config-set tool-mode <tool> <mode>")
			os.Exit(1)
		}
		tool := os.Args[3]
		mode := os.Args[4]
		if !config.ValidateMode(mode) {
			ui.PrintError("Invalid mode. Use: auto, interactive, manual, or disabled")
			os.Exit(1)
		}
		cfg.ToolSpecificModes[tool] = config.ActivationMode(mode)
		ui.PrintSuccess(fmt.Sprintf("Mode for %s set to: %s", tool, mode))

	case "confidence":
		if value == "true" {
			cfg.ShowConfidence = true
		} else if value == "false" {
			cfg.ShowConfidence = false
		} else {
			ui.PrintError("Invalid value. Use: true or false")
			os.Exit(1)
		}
		ui.PrintSuccess(fmt.Sprintf("Show confidence set to: %s", value))

	case "provider":
		provider := config.LLMProvider(value)
		if provider != config.ProviderOllama && provider != config.ProviderOpenCode {
			ui.PrintError("Invalid provider. Use: ollama or opencode")
			os.Exit(1)
		}
		cfg.Provider = provider
		ui.PrintSuccess(fmt.Sprintf("Provider set to: %s", value))

	case "model":
		cfg.PreferredModel = value
		ui.PrintSuccess(fmt.Sprintf("Preferred model set to: %s", value))

	default:
		ui.PrintError(fmt.Sprintf("Unknown key: %s", key))
		os.Exit(1)
	}

	// Save configuration
	if err := cfg.Save(configFile); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to save config: %v", err))
		os.Exit(1)
	}
}

// handleConfigReset resets configuration to defaults
func handleConfigReset(configFile string) {
	if !interactive.ShowConfirmation("Reset configuration to defaults?") {
		ui.PrintInfo("Canceled")
		return
	}

	cfg := config.DefaultConfig()
	if err := cfg.Save(configFile); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to save config: %v", err))
		os.Exit(1)
	}
	ui.PrintSuccess("Configuration reset to defaults")
}

func printUsage() {
	fmt.Printf(`AI Terminal Helper v%s (Go)

Usage:
  ai-helper analyze <command> <exit_code> [error_output]
  ai-helper proactive <query>
  ai-helper cache-stats
  ai-helper cache-clear
  ai-helper config-show
  ai-helper config-set <key> <value>
  ai-helper config-reset
  ai-helper version | -v | --version
  ai-helper help | -h | --help

Examples:
  ai-helper analyze "kubectl get pods" 127 "command not found"
  ai-helper proactive "how do I list all docker containers"
  ai-helper cache-stats
  ai-helper config-set mode interactive
`, version)
}
