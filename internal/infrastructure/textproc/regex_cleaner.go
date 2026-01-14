package textproc

import (
	"regexp"
	"strings"
)

// RegexSanitizer implementa gateway.TextSanitizer
type RegexSanitizer struct {
	// Podemos adicionar opções de filtro configuráveis aqui no futuro
}

func NewRegexSanitizer() *RegexSanitizer {
	return &RegexSanitizer{}
}

func (r *RegexSanitizer) Clean(raw string) (string, error) {
	// Remove cabeçalho VTT
	raw = strings.ReplaceAll(raw, "WEBVTT", "")

	// Remove timestamps (00:00:00.000 --> ...)
	reTime := regexp.MustCompile(`(?m)^.*-->.*$`)
	raw = reTime.ReplaceAllString(raw, "")

	// Remove tags internas <c.v1>, <00:00>, etc
	reTags := regexp.MustCompile(`<[^>]*>`)
	raw = reTags.ReplaceAllString(raw, "")

	lines := strings.Split(raw, "\n")
	var cleanLines []string
	seen := make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Remove linhas vazias ou lixo curto
		if len(line) < 3 {
			continue
		}

		// Desduplicação simples
		if !seen[line] {
			cleanLines = append(cleanLines, line)
			seen[line] = true
		}
	}

	return strings.Join(cleanLines, " "), nil
}
