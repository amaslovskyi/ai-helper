package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yourusername/ai-helper/pkg/cache"
	"github.com/yourusername/ai-helper/pkg/llm"
	"github.com/yourusername/ai-helper/pkg/security"
	"github.com/yourusername/ai-helper/pkg/ui"
	"github.com/yourusername/ai-helper/pkg/validators"
	"github.com/yourusername/ai-helper/pkg/validators/docker"
)

const version = "2.0.0-go"

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

	// Create validators
	validators := []validators.Validator{
		docker.NewValidator(),
		// TODO: Add kubectl, terraform, git validators
	}

	// Parse command
	cmd := os.Args[1]

	switch cmd {
	case "analyze":
		handleAnalyze(client, cacheStore, scanner, validators)
	case "proactive", "ask":
		handleProactive(client, scanner, validators)
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
	if err := validateCommand(resp.Suggestion, validators); err != nil {
		ui.PrintWarning(fmt.Sprintf("AI suggestion validation failed: %v", err))
		ui.PrintInfo("Querying AI again with validation context...")
		
		// Query again with validation error context
		req.Context = fmt.Sprintf("Previous suggestion '%s' was invalid: %v", resp.Suggestion, err)
		resp, err = client.Query(ctx, req)
		if err != nil {
			ui.PrintError(fmt.Sprintf("AI re-query failed: %v", err))
			os.Exit(1)
		}
		
		// Validate again
		if err := validateCommand(resp.Suggestion, validators); err != nil {
			ui.PrintError(fmt.Sprintf("AI still suggesting invalid command: %v", err))
			ui.PrintInfo("Suggestion: " + resp.Suggestion)
			os.Exit(1)
		}
	}

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

	// Print response
	printResponse(resp)
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
	if err := validateCommand(resp.Suggestion, validators); err != nil {
		ui.PrintWarning(fmt.Sprintf("Validation failed: %v", err))
		// In proactive mode, still show the suggestion but with warning
	}

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

	printResponse(resp)
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

