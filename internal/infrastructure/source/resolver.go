package source

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/domain/gateway"
)

// Resolver picks the appropriate ContentSource and source string from a user-supplied argument.
type Resolver struct {
	YouTube gateway.ContentSource
	File    gateway.ContentSource
	URL     gateway.ContentSource
	Stdin   gateway.ContentSource
}

// NewResolver builds a resolver with the given sources. Any nil source is skipped for that kind.
func NewResolver(youtube, file, url, stdin gateway.ContentSource) *Resolver {
	return &Resolver{YouTube: youtube, File: file, URL: url, Stdin: stdin}
}

// Resolve returns the ContentSource and the source string to pass to Fetch.
// arg can be: a URL (http(s) or youtube), a path (file: or plain path), or "-" / "stdin:" for stdin.
func (r *Resolver) Resolve(arg string) (gateway.ContentSource, string, error) {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		return nil, "", fmt.Errorf("empty source argument")
	}
	lower := strings.ToLower(arg)

	// Stdin
	if arg == "-" || lower == "stdin" || lower == "stdin:" {
		if r.Stdin == nil {
			return nil, "", fmt.Errorf("stdin source not configured")
		}
		return r.Stdin, arg, nil
	}

	// Explicit file:
	if strings.HasPrefix(lower, "file:") {
		if r.File == nil {
			return nil, "", fmt.Errorf("file source not configured")
		}
		return r.File, arg, nil
	}

	// HTTP(S) URL
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		if isYouTubeURL(arg) && r.YouTube != nil {
			return r.YouTube, arg, nil
		}
		if r.URL != nil {
			return r.URL, arg, nil
		}
		return nil, "", fmt.Errorf("no URL or YouTube source configured")
	}

	// Treat as file path if it exists
	if _, err := os.Stat(arg); err == nil {
		if r.File == nil {
			return nil, "", fmt.Errorf("file source not configured")
		}
		return r.File, arg, nil
	}

	return nil, "", fmt.Errorf("cannot resolve source %q (use URL, file path, or - for stdin)", arg)
}

// FetchWithResolve resolves arg and calls Fetch on the chosen source. Convenience for one-shot use.
func (r *Resolver) FetchWithResolve(ctx context.Context, arg string) (string, error) {
	src, source, err := r.Resolve(arg)
	if err != nil {
		return "", err
	}
	return src.Fetch(ctx, source)
}

func isYouTubeURL(s string) bool {
	return strings.Contains(s, "youtube.com") || strings.Contains(s, "youtu.be")
}
