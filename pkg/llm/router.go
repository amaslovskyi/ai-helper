package llm

import (
	"strings"
)

// RouterRule defines a condition → model mapping
type RouterRule struct {
	Keywords []string
	Model    Model
	Priority int // Higher priority rules are checked first
}

// Router selects the best model based on command context
type Router struct {
	rules []RouterRule
}

// NewRouter creates a new model router with default rules
func NewRouter() *Router {
	return &Router{
		rules: []RouterRule{
			// Infrastructure & Orchestration (8B - needs deep reasoning)
			{
				Keywords: []string{"kubectl", "helm", "terraform", "terragrunt", "aws", "gcloud", "azure"},
				Model:    Qwen38B,
				Priority: 100,
			},
			// Container & CI/CD (8B - complex configs)
			{
				Keywords: []string{"docker", "podman", "buildah", "gitlab-ci", "jenkins", "circleci"},
				Model:    Qwen38B,
				Priority: 90,
			},
			// Monitoring & Observability (8B - log analysis)
			{
				Keywords: []string{"prometheus", "grafana", "datadog", "kubectl logs", "stern"},
				Model:    Qwen38B,
				Priority: 80,
			},
			// Config management (8B - needs deep reasoning)
			{
				Keywords: []string{"ansible", "salt", "puppet"},
				Model:    Qwen38B,
				Priority: 70,
			},
			// ML/Data Engineering (4B instruction-tuned - stack traces)
			{
				Keywords: []string{"python", "pip", "conda", "poetry", "jupyter", "mlflow", "kubeflow", "ray", "spark"},
				Model:    Gemma34B,
				Priority: 60,
			},
			// Ultra-fast for trivial shell errors (1.7B - instant response)
			{
				Keywords: []string{"^cp ", "^mv ", "^rm ", "^mkdir ", "^touch ", "^grep ", "^find ", "^awk ", "^sed "},
				Model:    Qwen317B,
				Priority: 50,
			},
		},
	}
}

// SelectModel chooses the best model for a given command
func (r *Router) SelectModel(command string, mode RequestMode) Model {
	// Proactive mode always uses 8B for better quality
	if mode == ModeProactive {
		return Qwen38B
	}

	cmdLower := strings.ToLower(command)

	// Find matching rule with highest priority
	var selectedModel Model = Qwen34B // Default fast model
	highestPriority := -1

	for _, rule := range r.rules {
		for _, keyword := range rule.Keywords {
			// Check for prefix match (for patterns like ^cp)
			if strings.HasPrefix(keyword, "^") {
				keyword = strings.TrimPrefix(keyword, "^")
				if strings.HasPrefix(cmdLower, keyword) {
					if rule.Priority > highestPriority {
						selectedModel = rule.Model
						highestPriority = rule.Priority
					}
					break
				}
			} else {
				// Check for substring match
				if strings.Contains(cmdLower, keyword) {
					if rule.Priority > highestPriority {
						selectedModel = rule.Model
						highestPriority = rule.Priority
					}
					break
				}
			}
		}
	}

	return selectedModel
}

// AddRule adds a custom routing rule
func (r *Router) AddRule(rule RouterRule) {
	r.rules = append(r.rules, rule)
}

