package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type OpenCodeClient struct {
	model   Model
	timeout time.Duration
}

func NewOpenCodeClient(model string) *OpenCodeClient {
	if model == "" {
		model = string(OpenCodeClaudeSonnet)
	}
	return &OpenCodeClient{
		model:   Model(model),
		timeout: 60 * time.Second,
	}
}

func (c *OpenCodeClient) Query(ctx context.Context, req Request) (*Response, error) {
	prompt := c.buildPrompt(req)

	var cmd *exec.Cmd
	if strings.Contains(string(c.model), "/") {
		cmd = exec.CommandContext(ctx, "opencode", "run", "-m", string(c.model), prompt)
	} else {
		cmd = exec.CommandContext(ctx, "opencode", "run", prompt)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("opencode command failed: %w, stderr: %s", err, stderr.String())
	}

	return c.parseResponse(stdout.String()), nil
}

func (c *OpenCodeClient) buildPrompt(req Request) string {
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

func (c *OpenCodeClient) parseResponse(text string) *Response {
	response := &Response{
		Model:      c.model,
		Provider:   ProviderOpenCode,
		Confidence: 0.8,
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

func (c *OpenCodeClient) IsAvailable(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "opencode", "version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("opencode not available: %w", err)
	}
	return nil
}

func (c *OpenCodeClient) ListModels(ctx context.Context) ([]Model, error) {
	cmd := exec.CommandContext(ctx, "opencode", "models")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}

	var models []Model
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && strings.Contains(line, "/") {
			models = append(models, Model(line))
		}
	}

	return models, nil
}

func (c *OpenCodeClient) GetProvider() Provider {
	return ProviderOpenCode
}

type opencodeConfig struct {
	Model    string `json:"model"`
	Provider string `json:"provider"`
}

func (c *OpenCodeClient) GetConfig() (*opencodeConfig, error) {
	cmd := exec.Command("opencode", "config", "get", "--format", "json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	var config opencodeConfig
	if err := json.Unmarshal(output, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}
