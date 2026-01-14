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

// LMStudioClient implements the gateway.LLMProvider interface
type LMStudioClient struct {
	BaseURL   string
	ModelName string
	Client    *http.Client
}

// NewLMStudioClient creates a new instance with the configured timeout
func NewLMStudioClient(url, model string, timeout time.Duration) *LMStudioClient {
	return &LMStudioClient{
		BaseURL:   url,
		ModelName: model,
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Call makes a model call using an OpenAI-compatible payload
func (l *LMStudioClient) Call(req entity.AgentRequest) (string, error) {
	// Build payload as expected by LM Studio/OpenAI
	payload := map[string]interface{}{
		"model": l.ModelName,
		"messages": []map[string]string{
			{"role": "system", "content": fmt.Sprintf("Atue como: %s. %s", req.Role, req.Instruction)},
			{"role": "user", "content": req.InputData},
		},
		"temperature": req.Temperature,
		"max_tokens":  -1,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error creating json payload: %w", err)
	}

	httpReq, err := http.NewRequest("POST", l.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := l.Client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to connect to LM Studio: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API Error %d: %s", resp.StatusCode, string(body))
	}

	return l.parseResponse(resp.Body)
}

func (l *LMStudioClient) parseResponse(body io.Reader) (string, error) {
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("empty model response")
}
