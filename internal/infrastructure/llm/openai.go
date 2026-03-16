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

// OpenAIClient implements gateway.LLMProvider for OpenAI (or any OpenAI-compatible API with Bearer token).
type OpenAIClient struct {
	BaseURL   string
	ModelName string
	APIKey    string
	Client    *http.Client
}

// NewOpenAIClient creates a client for OpenAI or compatible APIs.
func NewOpenAIClient(baseURL, model, apiKey string, timeout time.Duration) *OpenAIClient {
	if timeout == 0 {
		timeout = 60 * time.Second
	}
	return &OpenAIClient{
		BaseURL:   baseURL,
		ModelName: model,
		APIKey:    apiKey,
		Client:    &http.Client{Timeout: timeout},
	}
}

// Call sends a chat completion request. Implements gateway.LLMProvider.
func (o *OpenAIClient) Call(req entity.AgentRequest) (string, error) {
	payload := map[string]interface{}{
		"model": o.ModelName,
		"messages": []map[string]string{
			{"role": "system", "content": fmt.Sprintf("Act as: %s. %s", req.Role, req.Instruction)},
			{"role": "user", "content": req.InputData},
		},
		"temperature": req.Temperature,
		"max_tokens":  4096,
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
	if o.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+o.APIKey)
	}
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
