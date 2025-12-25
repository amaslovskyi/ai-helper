package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/amaslovskyi/ai-helper/pkg/cache"
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

const version = "2.1.0-go"

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

	// Create cache
	cacheStore, err := cache.NewCache(cacheFile)
	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to initialize cache: %v", err))
		os.Exit(1)
	}

	// Create LLM client
	client := llm.NewOllamaClient("")

	// Create security scanner
	scanner := security.NewScanner()

	// Create validators (order matters: more specific first)
	validatorsList := []validators.Validator{
		kubectl.NewValidator(),      // k, kubectl
		terraform.NewValidator(),    // tf, terraform
		terragrunt.NewValidator(),   // tg, terragrunt
		helm.NewValidator(),         // h, helm
		git.NewValidator(),          // git, gco, gcb, gp, etc.
		docker.NewValidator(),       // docker, d, dc
		ansible.NewValidator(),      // ansible, ansible-playbook
		argocd.NewValidator(),       // argocd
	}

	// Parse command
	cmd := os.Args[1]

	switch cmd {
	case "analyze":
		handleAnalyze(client, cacheStore, scanner, validatorsList)
	case "proactive", "ask":
		handleProactive(client, scanner, validatorsList)
	case "version":
		fmt.Printf("AI Terminal Helper v%s (Go)\n", version)
	case "cache-stats":
		handleCacheStats(cacheStore)
	case "cache-clear":
		handleCacheClear(cacheStore)
	default:
		ui.PrintError(fmt.Sprintf("Unknown command: %s", cmd))
		printUsage()
		os.Exit(1)
	}
}

func handleAnalyze(client llm.Client, cacheStore *cache.Cache, scanner *security.Scanner, validators []validators.Validator) {
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

	// Try cache first
	if cachedResp, ok := cacheStore.Get(command, errorOutput); ok {
		fmt.Println(ui.Colorize(ui.MagentaBold, "💾 [Cached]"))
		printResponse(cachedResp)
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

func handleProactive(client llm.Client, scanner *security.Scanner, validators []validators.Validator) {
	if len(os.Args) < 3 {
		ui.PrintError("Usage: ai-helper proactive <query>")
		os.Exit(1)
	}

	query := strings.Join(os.Args[2:], " ")

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

func printUsage() {
	fmt.Printf(`AI Terminal Helper v%s (Go)

Usage:
  ai-helper analyze <command> <exit_code> [error_output]
  ai-helper proactive <query>
  ai-helper cache-stats
  ai-helper cache-clear
  ai-helper version

Examples:
  ai-helper analyze "kubectl get pods" 127 "command not found"
  ai-helper proactive "how do I list all docker containers"
  ai-helper cache-stats
`, version)
}

