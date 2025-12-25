package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ActivationMode defines how AI assistance is triggered
type ActivationMode string

const (
	// ModeAuto - AI triggers automatically on command failures (default)
	ModeAuto ActivationMode = "auto"
	
	// ModeInteractive - Shows interactive menu on failures, user chooses action
	ModeInteractive ActivationMode = "interactive"
	
	// ModeManual - AI only activates with explicit commands (ask, analyze)
	ModeManual ActivationMode = "manual"
	
	// ModeDisabled - AI is completely disabled
	ModeDisabled ActivationMode = "disabled"
)

// Config represents the user's configuration preferences
type Config struct {
	// ActivationMode controls how AI assistance is triggered
	ActivationMode ActivationMode `json:"activation_mode"`
	
	// AutoExecuteSafe allows automatic execution of safe, read-only commands
	AutoExecuteSafe bool `json:"auto_execute_safe"`
	
	// ShowConfidence displays confidence scores with AI suggestions
	ShowConfidence bool `json:"show_confidence"`
	
	// PreferredModel is the default Ollama model to use
	PreferredModel string `json:"preferred_model"`
	
	// ToolSpecificModes allows per-tool activation overrides
	// Example: {"kubectl": "interactive", "docker": "auto"}
	ToolSpecificModes map[string]ActivationMode `json:"tool_specific_modes"`
	
	// SessionDisabled is used for temporary session-level disabling
	// This is not saved to disk, only in-memory
	SessionDisabled bool `json:"-"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		ActivationMode:    ModeAuto,
		AutoExecuteSafe:   false,
		ShowConfidence:    true,
		PreferredModel:    "", // Empty means auto-select
		ToolSpecificModes: make(map[string]ActivationMode),
		SessionDisabled:   false,
	}
}

// Load loads configuration from disk, returns default if not found
func Load(configFile string) (*Config, error) {
	// Ensure directory exists
	dir := filepath.Dir(configFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}
	
	// If file doesn't exist, return default config
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}
	
	// Read and parse config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		// If config is corrupted, return default
		return DefaultConfig(), nil
	}
	
	// Ensure map is initialized
	if cfg.ToolSpecificModes == nil {
		cfg.ToolSpecificModes = make(map[string]ActivationMode)
	}
	
	return &cfg, nil
}

// Save saves configuration to disk
func (c *Config) Save(configFile string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	return os.WriteFile(configFile, data, 0644)
}

// GetModeForTool returns the activation mode for a specific tool
// Returns tool-specific mode if set, otherwise returns global mode
func (c *Config) GetModeForTool(tool string) ActivationMode {
	// Check if session is disabled
	if c.SessionDisabled {
		return ModeDisabled
	}
	
	// Check tool-specific override
	if mode, exists := c.ToolSpecificModes[tool]; exists {
		return mode
	}
	
	// Return global mode
	return c.ActivationMode
}

// ShouldTriggerAI determines if AI should be triggered based on mode
func (c *Config) ShouldTriggerAI(tool string) bool {
	mode := c.GetModeForTool(tool)
	return mode == ModeAuto
}

// ShouldShowMenu determines if interactive menu should be shown
func (c *Config) ShouldShowMenu(tool string) bool {
	mode := c.GetModeForTool(tool)
	return mode == ModeInteractive
}

// IsEnabled checks if AI is enabled (not disabled)
func (c *Config) IsEnabled(tool string) bool {
	mode := c.GetModeForTool(tool)
	return mode != ModeDisabled
}

// ValidateMode checks if a mode string is valid
func ValidateMode(mode string) bool {
	switch ActivationMode(mode) {
	case ModeAuto, ModeInteractive, ModeManual, ModeDisabled:
		return true
	default:
		return false
	}
}

