package llm

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/domain/gateway"
)

// Env keys for LLM configuration.
const (
	EnvLLMProvider = "LLM_PROVIDER"   // lmstudio, openai, ollama
	EnvLLMEndpoint = "LLM_ENDPOINT"   // base URL or full chat completions URL
	EnvLLMModel   = "LLM_MODEL"       // model name
	EnvLLMTimeout = "LLM_TIMEOUT"     // seconds, e.g. 600
	EnvOpenAIKey  = "OPENAI_API_KEY"  // for provider=openai
)

// NewFromEnv returns an LLMProvider from environment variables.
// LLM_PROVIDER: lmstudio | openai | ollama (default: lmstudio)
// LLM_ENDPOINT: URL (e.g. http://localhost:1234/v1/chat/completions for LM Studio)
// LLM_MODEL: model name
// LLM_TIMEOUT: seconds (default 600)
// OPENAI_API_KEY: required when LLM_PROVIDER=openai
func NewFromEnv() (gateway.LLMProvider, error) {
	provider := strings.ToLower(strings.TrimSpace(os.Getenv(EnvLLMProvider)))
	if provider == "" {
		provider = "lmstudio"
	}
	endpoint := strings.TrimSpace(os.Getenv(EnvLLMEndpoint))
	model := strings.TrimSpace(os.Getenv(EnvLLMModel))
	if model == "" {
		model = "local-model"
	}
	timeoutSec := 600
	if s := os.Getenv(EnvLLMTimeout); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			timeoutSec = n
		}
	}
	timeout := time.Duration(timeoutSec) * time.Second

	switch provider {
	case "lmstudio":
		if endpoint == "" {
			endpoint = "http://host.docker.internal:1234/v1/chat/completions"
		}
		return NewLMStudioClient(endpoint, model, timeout), nil
	case "openai":
		if endpoint == "" {
			endpoint = "https://api.openai.com/v1/chat/completions"
		}
		key := strings.TrimSpace(os.Getenv(EnvOpenAIKey))
		return NewOpenAIClient(endpoint, model, key, timeout), nil
	case "ollama":
		if endpoint == "" {
			endpoint = DefaultOllamaURL
		}
		// Ollama base URL is without /v1/chat/completions; NewOllamaClient appends it
		return NewOllamaClient(endpoint, model, timeout), nil
	default:
		return nil, fmt.Errorf("unknown LLM_PROVIDER %q (use lmstudio, openai, or ollama)", provider)
	}
}
