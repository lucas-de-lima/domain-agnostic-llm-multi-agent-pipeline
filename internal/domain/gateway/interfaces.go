package gateway

import (
	"context"

	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/domain/entity"
)

// LLMProvider define como conversamos com a IA (LM Studio, OpenAI, Ollama)
type LLMProvider interface {
	Call(req entity.AgentRequest) (string, error)
}

// ContentSource define como obtemos o texto cru (YouTube, arquivo local, URL, stdin).
// O source pode ser URL, path de arquivo, ou convenção (ex.: "-" ou "stdin:").
type ContentSource interface {
	Fetch(ctx context.Context, source string) (raw string, err error)
}

// SourceResolver resolves a user-supplied argument (URL, path, "-") and fetches raw content.
// Use case depends on this to avoid depending on concrete source implementations.
type SourceResolver interface {
	FetchWithResolve(ctx context.Context, arg string) (raw string, err error)
}

// ContentDownloader define como obtemos o texto cru (Youtube, Arquivo Local, URL).
// Mantido para compatibilidade; preferir ContentSource em código novo.
type ContentDownloader interface {
	Download(source string) (string, error)
}

// TextSanitizer define como limpamos o lixo inicial (VTT, Logs, HTML)
type TextSanitizer interface {
	Clean(raw string) (string, error)
}
