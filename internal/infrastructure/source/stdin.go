package source

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// StdinSource implements gateway.ContentSource by reading from os.Stdin.
// Use source "-" or "stdin:" to read from stdin.
type StdinSource struct{}

// NewStdinSource creates a ContentSource that reads from stdin.
func NewStdinSource() *StdinSource {
	return &StdinSource{}
}

// Fetch reads from stdin. Implements gateway.ContentSource.
// The source value is ignored; it should be "-" or "stdin:" for clarity.
func (s *StdinSource) Fetch(ctx context.Context, source string) (string, error) {
	_ = source
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("read stdin: %w", err)
	}
	return strings.TrimSpace(string(data)), nil
}
