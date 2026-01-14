package gateway

import "github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/domain/entity"

// LLMProvider define como conversamos com a IA (LM Studio, OpenAI, Ollama)
type LLMProvider interface {
	Call(req entity.AgentRequest) (string, error)
}

// ContentDownloader define como obtemos o texto cru (Youtube, Arquivo Local, URL)
type ContentDownloader interface {
	Download(source string) (string, error)
}

// TextSanitizer define como limpamos o lixo inicial (VTT, Logs, HTML)
type TextSanitizer interface {
	Clean(raw string) (string, error)
}
