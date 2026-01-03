package llm

import (
	"strings"
)

type RouterRule struct {
	Keywords []string
	Model    Model
	Priority int // Higher priority rules are checked first
}

type Router struct {
	rules    []RouterRule
	provider Provider
}

func NewRouter(provider Provider) *Router {
	r := &Router{
		provider: provider,
	}

	if provider == ProviderOpenCode {
		r.rules = []RouterRule{
			{
				Keywords: []string{"kubectl", "helm", "terraform", "terragrunt", "aws", "gcloud", "azure"},
				Model:    OpenCodeClaudeSonnet,
				Priority: 100,
			},
			{
				Keywords: []string{"docker", "podman", "buildah", "gitlab-ci", "jenkins", "circleci"},
				Model:    OpenCodeClaudeSonnet,
				Priority: 90,
			},
			{
				Keywords: []string{"prometheus", "grafana", "datadog", "kubectl logs", "stern"},
				Model:    OpenCodeClaudeSonnet,
				Priority: 80,
			},
			{
				Keywords: []string{"ansible", "salt", "puppet"},
				Model:    OpenCodeClaudeSonnet,
				Priority: 70,
			},
			{
				Keywords: []string{"python", "pip", "conda", "poetry", "jupyter", "mlflow", "kubeflow", "ray", "spark"},
				Model:    OpenCodeGPT4o,
				Priority: 60,
			},
			{
				Keywords: []string{"^cp ", "^mv ", "^rm ", "^mkdir ", "^touch ", "^grep ", "^find ", "^awk ", "^sed "},
				Model:    OpenCodeGPT4oMini,
				Priority: 50,
			},
		}
	} else {
		r.rules = []RouterRule{
			{
				Keywords: []string{"kubectl", "helm", "terraform", "terragrunt", "aws", "gcloud", "azure"},
				Model:    Qwen38B,
				Priority: 100,
			},
			{
				Keywords: []string{"docker", "podman", "buildah", "gitlab-ci", "jenkins", "circleci"},
				Model:    Qwen38B,
				Priority: 90,
			},
			{
				Keywords: []string{"prometheus", "grafana", "datadog", "kubectl logs", "stern"},
				Model:    Qwen38B,
				Priority: 80,
			},
			{
				Keywords: []string{"ansible", "salt", "puppet"},
				Model:    Qwen38B,
				Priority: 70,
			},
			{
				Keywords: []string{"python", "pip", "conda", "poetry", "jupyter", "mlflow", "kubeflow", "ray", "spark"},
				Model:    Gemma34B,
				Priority: 60,
			},
			{
				Keywords: []string{"^cp ", "^mv ", "^rm ", "^mkdir ", "^touch ", "^grep ", "^find ", "^awk ", "^sed "},
				Model:    Qwen317B,
				Priority: 50,
			},
		}
	}

	return r
}

func (r *Router) SelectModel(command string, mode RequestMode) Model {
	if mode == ModeProactive {
		if r.provider == ProviderOpenCode {
			return OpenCodeClaudeSonnet
		}
		return Qwen38B
	}

	cmdLower := strings.ToLower(command)

	var selectedModel Model
	if r.provider == ProviderOpenCode {
		selectedModel = OpenCodeGPT4o
	} else {
		selectedModel = Qwen34B
	}
	highestPriority := -1

	for _, rule := range r.rules {
		for _, keyword := range rule.Keywords {
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

func (r *Router) AddRule(rule RouterRule) {
	r.rules = append(r.rules, rule)
}
