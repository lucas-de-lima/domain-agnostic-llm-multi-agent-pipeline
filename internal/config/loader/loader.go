package loader

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/domain/entity"
)

//go:embed default/*.txt
var defaultPrompts embed.FS

// Load returns the specialization for the given id.
// id "default" loads embedded domain-agnostic prompts; any other id loads from specializations/<id>/ on disk.
// baseDir is the directory containing the "specializations" folder (e.g. working directory or executable dir). If empty, only "default" can be loaded.
func Load(id string, baseDir string) (*entity.Specialization, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		id = "default"
	}
	if id == "default" {
		return loadDefault()
	}
	return loadFromDir(baseDir, id)
}

func loadDefault() (*entity.Specialization, error) {
	ctx, _ := defaultPrompts.ReadFile("default/context.txt")
	ext, _ := defaultPrompts.ReadFile("default/extract.txt")
	syn, _ := defaultPrompts.ReadFile("default/synthesize.txt")
	aud, _ := defaultPrompts.ReadFile("default/audit.txt")
	return &entity.Specialization{
		ID:               "default",
		Name:             "Default",
		Description:      "Domain-agnostic pipeline for any technical or educational content.",
		ContextPrompt:    string(ctx),
		ExtractPrompt:    string(ext),
		SynthesizePrompt: string(syn),
		AuditPrompt:      string(aud),
	}, nil
}

func loadFromDir(baseDir, id string) (*entity.Specialization, error) {
	dir := filepath.Join(baseDir, "specializations", id)
	manifestPath := filepath.Join(dir, "manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("read manifest: %w", err)
	}
	var manifest struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}
	if manifest.ID == "" {
		manifest.ID = id
	}
	ctx, err := os.ReadFile(filepath.Join(dir, "context.txt"))
	if err != nil {
		return nil, fmt.Errorf("context.txt: %w", err)
	}
	ext, err := os.ReadFile(filepath.Join(dir, "extract.txt"))
	if err != nil {
		return nil, fmt.Errorf("extract.txt: %w", err)
	}
	syn, err := os.ReadFile(filepath.Join(dir, "synthesize.txt"))
	if err != nil {
		return nil, fmt.Errorf("synthesize.txt: %w", err)
	}
	aud, err := os.ReadFile(filepath.Join(dir, "audit.txt"))
	if err != nil {
		return nil, fmt.Errorf("audit.txt: %w", err)
	}
	return &entity.Specialization{
		ID:               manifest.ID,
		Name:             manifest.Name,
		Description:      manifest.Description,
		ContextPrompt:    string(ctx),
		ExtractPrompt:    string(ext),
		SynthesizePrompt: string(syn),
		AuditPrompt:      string(aud),
	}, nil
}
