package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OllamaClient implements the Client interface for Ollama
type OllamaClient struct {
	baseURL    string
	httpClient *http.Client
	router     *Router
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(baseURL string) *OllamaClient {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	return &OllamaClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		router: NewRouter(ProviderOllama),
	}
}

// ollamaRequest represents an Ollama API request
type ollamaRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// ollamaResponse represents an Ollama API response
type ollamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

// Query sends a request to Ollama
func (c *OllamaClient) Query(ctx context.Context, req Request) (*Response, error) {
	// Select appropriate model
	model := c.router.SelectModel(req.Command, req.Mode)

	// Build prompt
	prompt := c.buildPrompt(req)

	// Create Ollama request
	ollamaReq := ollamaRequest{
		Model:  string(model),
		Prompt: prompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.7,
			"num_ctx":     4096,
		},
	}

	reqBody, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Send request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/generate", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Parse AI response into structured format
	return c.parseResponse(ollamaResp.Response, model), nil
}

// buildPrompt constructs the prompt based on request mode
func (c *OllamaClient) buildPrompt(req Request) string {
	if req.Mode == ModeProactive {
		return fmt.Sprintf(`You are a senior DevOps/SRE. Convert this natural language query to a command.

CRITICAL RULES:
1. DO NOT output "Thinking..." or any reasoning process
2. START IMMEDIATELY with ✓ followed by the command
3. NO verbose reasoning, NO process explanation

Query: %s
Context: %s
Dir: %s

REQUIRED OUTPUT FORMAT (start immediately):
✓ [command]
Root: [1 sentence what this does]
Tip: [optional safety note or best practice]

Your first line MUST be: ✓ [command]`, req.Command, req.Context, req.Directory)
	}

	// Reactive mode
	return fmt.Sprintf(`You are a senior DevOps/SRE. Fix this failed command.

CRITICAL RULES:
1. DO NOT output "Thinking..." or any reasoning process
2. DO NOT start with "Okay," "Let me," "Wait," or any explanation
3. START IMMEDIATELY with ✓ followed by the corrected command
4. NO thinking blocks, NO verbose reasoning, NO process explanation

Command: %s
Error: %s
Exit: %d
Dir: %s

REQUIRED OUTPUT FORMAT (start immediately, no preamble):
✓ [corrected command]
Root: [1 sentence why it failed]
Tip: [optional best practice]

Your first line MUST be: ✓ [command]`, req.Command, req.Error, req.ExitCode, req.Directory)
}

// parseResponse extracts structured data from AI response
func (c *OllamaClient) parseResponse(text string, model Model) *Response {
	response := &Response{
		Model:      model,
		Confidence: 0.8, // Default confidence
	}

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "✓"):
			response.Suggestion = strings.TrimSpace(strings.TrimPrefix(line, "✓"))
		case strings.HasPrefix(line, "Root:"):
			response.RootCause = strings.TrimSpace(strings.TrimPrefix(line, "Root:"))
		case strings.HasPrefix(line, "Tip:"):
			response.Tip = strings.TrimSpace(strings.TrimPrefix(line, "Tip:"))
		}
	}

	return response
}

// IsAvailable checks if Ollama is running
func (c *OllamaClient) IsAvailable(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/tags", nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ollama not available: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	return nil
}

// GetProvider returns the provider type
func (c *OllamaClient) GetProvider() Provider {
	return ProviderOllama
}

// ListModels returns available Ollama models
func (c *OllamaClient) ListModels(ctx context.Context) ([]Model, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/tags", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	models := make([]Model, 0, len(result.Models))
	for _, m := range result.Models {
		models = append(models, Model(m.Name))
	}

	return models, nil
}
