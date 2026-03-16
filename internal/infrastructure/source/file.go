package source

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileSource implements gateway.ContentSource by reading local files (txt, md, vtt).
type FileSource struct {
	// AllowedExt restricts which extensions are read (e.g. [".txt", ".md", ".vtt"]). Empty = all.
	AllowedExt []string
}

// NewFileSource creates a ContentSource for local files.
func NewFileSource() *FileSource {
	return &FileSource{AllowedExt: []string{".txt", ".md", ".vtt", ".srt"}}
}

// Fetch reads the file at the given path. Implements gateway.ContentSource.
// source must be a filesystem path. Prefix "file:" is stripped if present.
func (f *FileSource) Fetch(ctx context.Context, source string) (string, error) {
	_ = ctx
	path := strings.TrimPrefix(source, "file:")
	path = strings.TrimPrefix(path, "//")
	path = filepath.Clean(path)
	if path == "" || path == "." {
		return "", fmt.Errorf("invalid file source: empty path")
	}
	ext := strings.ToLower(filepath.Ext(path))
	if len(f.AllowedExt) > 0 {
		ok := false
		for _, e := range f.AllowedExt {
			if ext == e {
				ok = true
				break
			}
		}
		if !ok {
			return "", fmt.Errorf("file extension %q not allowed (allowed: %v)", ext, f.AllowedExt)
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}
	return string(data), nil
}
