package llm

import "context"

// Model represents an Ollama model
type Model string

const (
	// Primary models
	Qwen38B     Model = "qwen3:8b-q4_K_M"
	Qwen34B     Model = "qwen3:4b-q4_K_M"
	Gemma34B    Model = "gemma3:4b-it-q4_K_M"
	Qwen317B    Model = "qwen3:1.7b-q4_K_M"
	Gemma31B    Model = "gemma3:1b-it-q4_K_M"
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
	ModeReactive   RequestMode = "reactive"   // Fix failed command
	ModeProactive  RequestMode = "proactive"  // Generate command from natural language
)

// Response represents an AI response
type Response struct {
	Suggestion  string // The corrected command or generated command
	RootCause   string // Why it failed or what it does
	Tip         string // Best practice or safety note
	Cached      bool   // Whether this came from cache
	Model       Model  // Which model generated this
	Confidence  float64 // Confidence score (0-1)
}

// Client interface for LLM interactions
type Client interface {
	// Query sends a request to the LLM and returns a response
	Query(ctx context.Context, req Request) (*Response, error)
	
	// IsAvailable checks if the LLM service is available
	IsAvailable(ctx context.Context) error
	
	// ListModels returns available models
	ListModels(ctx context.Context) ([]Model, error)
}

