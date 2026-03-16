package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/domain/entity"
)

// DefaultOllamaURL is the default Ollama API base.
const DefaultOllamaURL = "http://localhost:11434"

// OllamaClient implements gateway.LLMProvider for Ollama (OpenAI-compatible local API).
type OllamaClient struct {
	BaseURL   string
	ModelName string
	Client    *http.Client
}

// NewOllamaClient creates a client for Ollama. endpoint can be empty to use DefaultOllamaURL.
func NewOllamaClient(endpoint, model string, timeout time.Duration) *OllamaClient {
	if endpoint == "" {
		endpoint = DefaultOllamaURL
	}
	url := endpoint + "/v1/chat/completions"
	if timeout == 0 {
		timeout = 5 * time.Minute
	}
	return &OllamaClient{
		BaseURL:   url,
		ModelName: model,
		Client:    &http.Client{Timeout: timeout},
	}
}

// Call sends a chat completion request. Implements gateway.LLMProvider.
func (o *OllamaClient) Call(req entity.AgentRequest) (string, error) {
	payload := map[string]interface{}{
		"model": o.ModelName,
		"messages": []map[string]string{
			{"role": "system", "content": fmt.Sprintf("Act as: %s. %s", req.Role, req.Instruction)},
			{"role": "user", "content": req.InputData},
		},
		"temperature": req.Temperature,
		"stream":      false,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("json payload: %w", err)
	}
	httpReq, err := http.NewRequest("POST", o.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := o.Client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API %d: %s", resp.StatusCode, string(body))
	}
	return ParseChatResponse(resp.Body)
}
