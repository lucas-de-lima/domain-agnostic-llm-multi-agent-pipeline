// Package main is a backward-compatible entrypoint that runs the pipeline with default specialization.
// Prefer cmd/pipeline for full flags (--specialization, multiple sources).
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/config/loader"
	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/domain/gateway"
	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/infrastructure/downloader"
	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/infrastructure/llm"
	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/infrastructure/source"
	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/infrastructure/textproc"
	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/usecase/pipeline"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: ./distiller <URL_OR_SOURCE>")
	}
	sourceArg := os.Args[1]

	llmProvider, err := llm.NewFromEnv()
	if err != nil {
		log.Fatalf("LLM config: %v", err)
	}

	tempDir := "output/temp"
	if d := os.Getenv("TEMP_DIR"); d != "" {
		tempDir = d
	}
	subLangs := []string{"en", "pt"}
	if s := os.Getenv("SUBTITLE_LANGS"); s != "" {
		parts := strings.Split(s, ",")
		subLangs = make([]string, 0, len(parts))
		for _, p := range parts {
			if t := strings.TrimSpace(p); t != "" {
				subLangs = append(subLangs, t)
			}
		}
	}
	ytdlp := downloader.NewYtDlpClientWithLangs(tempDir, subLangs)

	var youtube gateway.ContentSource = source.NewYouTubeSource(ytdlp)
	var file gateway.ContentSource = source.NewFileSource()
	var url gateway.ContentSource = source.NewURLSource()
	var stdin gateway.ContentSource = source.NewStdinSource()
	resolver := source.NewResolver(youtube, file, url, stdin)
	sanitizer := textproc.NewRegexSanitizer()

	baseDir, _ := os.Getwd()
	spec, err := loader.Load("default", baseDir)
	if err != nil {
		log.Fatalf("Load default specialization: %v", err)
	}

	uc := pipeline.NewDistillerUseCase(llmProvider, resolver, sanitizer, spec)

	log.Println("Starting pipeline...")
	start := time.Now()
	result, err := uc.Run(context.Background(), sourceArg)
	if err != nil {
		log.Fatalf("Pipeline error: %v", err)
	}

	outDir := "output"
	if d := os.Getenv("OUTPUT_DIR"); d != "" {
		outDir = d
	}
	_ = os.MkdirAll(outDir, 0755)
	outputFilename := filepath.Join(outDir, fmt.Sprintf("knowledge_%d.md", time.Now().Unix()))
	if err := os.WriteFile(outputFilename, []byte(result), 0644); err != nil {
		log.Fatalf("Write output: %v", err)
	}
	log.Printf("Done in %s. Output: %s", time.Since(start), outputFilename)
}
