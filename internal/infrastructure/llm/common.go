package llm

import (
	"encoding/json"
	"fmt"
	"io"
)

// ParseChatResponse decodes OpenAI-compatible chat completion JSON and returns the first message content.
func ParseChatResponse(body io.Reader) (string, error) {
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("empty response")
}
