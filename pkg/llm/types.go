package llm

import "context"

// Provider represents an LLM provider
type Provider string

const (
	// ProviderOllama - Local Ollama server
	ProviderOllama Provider = "ollama"

	// ProviderOpenCode - OpenCode AI coding agent
	ProviderOpenCode Provider = "opencode"
)

// Model represents an LLM model
type Model string

const (
	// Ollama models
	Qwen38B  Model = "qwen3:8b-q4_K_M"
	Qwen34B  Model = "qwen3:4b-q4_K_M"
	Gemma34B Model = "gemma3:4b-it-q4_K_M"
	Qwen317B Model = "qwen3:1.7b-q4_K_M"
	Gemma31B Model = "gemma3:1b-it-q4_K_M"

	// OpenCode models (provider/model format)
	OpenCodeClaudeSonnet Model = "anthropic/claude-sonnet-4-20250514"
	OpenCodeClaudeOpus   Model = "anthropic/claude-opus-4-20250514"
	OpenCodeGPT4o        Model = "openai/gpt-4o"
	OpenCodeGPT4oMini    Model = "openai/gpt-4o-mini"
	OpenCodeGPT35Turbo   Model = "openai/gpt-3.5-turbo"
	OpenCodeGeminiPro    Model = "google/gemini-pro"
	OpenCodeLlama3       Model = "ollama/llama3"
	OpenCodeQwen         Model = "ollama/qwen"
)

// Request represents an AI request
type Request struct {
	Command   string
	Error     string
	ExitCode  int
	Directory string
	Context   string
	Mode      RequestMode
}

// RequestMode defines the type of request
type RequestMode string

const (
	ModeReactive  RequestMode = "reactive"  // Fix failed command
	ModeProactive RequestMode = "proactive" // Generate command from natural language
)

// Response represents an AI response
type Response struct {
	Suggestion string   // The corrected command or generated command
	RootCause  string   // Why it failed or what it does
	Tip        string   // Best practice or safety note
	Cached     bool     // Whether this came from cache
	Model      Model    // Which model generated this
	Confidence float64  // Confidence score (0-1)
	Provider   Provider // Which provider generated this
}

// Client interface for LLM interactions
type Client interface {
	// Query sends a request to the LLM and returns a response
	Query(ctx context.Context, req Request) (*Response, error)

	// IsAvailable checks if the LLM service is available
	IsAvailable(ctx context.Context) error

	// ListModels returns available models
	ListModels(ctx context.Context) ([]Model, error)

	// GetProvider returns the provider type
	GetProvider() Provider
}
