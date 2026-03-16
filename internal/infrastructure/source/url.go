package source

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// URLSource implements gateway.ContentSource by fetching HTTP(S) URLs.
// HTML responses are stripped of tags to produce plain text.
type URLSource struct {
	Client *http.Client
}

// NewURLSource creates a ContentSource for HTTP URLs.
func NewURLSource() *URLSource {
	return &URLSource{
		Client: &http.Client{Timeout: 60 * time.Second},
	}
}

// Fetch fetches the URL and returns body as text. Implements gateway.ContentSource.
// If Content-Type suggests HTML, tags are stripped.
func (u *URLSource) Fetch(ctx context.Context, source string) (string, error) {
	urlStr := strings.TrimPrefix(source, "url:")
	urlStr = strings.TrimPrefix(urlStr, "//")
	if urlStr == "" {
		return "", fmt.Errorf("invalid url source: empty url")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	resp, err := u.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}
	text := string(body)
	ct := resp.Header.Get("Content-Type")
	if strings.Contains(ct, "html") {
		text = stripHTML(text)
	}
	return text, nil
}

// stripHTML removes HTML tags and collapses whitespace.
func stripHTML(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			b.WriteRune(r)
		}
	}
	// Collapse whitespace and trim
	lines := strings.Fields(b.String())
	return strings.Join(lines, " ")
}
