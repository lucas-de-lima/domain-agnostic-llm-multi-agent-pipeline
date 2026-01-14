package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/infrastructure/downloader"
	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/infrastructure/llm"
	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/infrastructure/textproc"
	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/usecase/pipeline"
)

func main() {
	// 1. Configurações (Poderiam vir de ENV)
	if len(os.Args) < 2 {
		log.Fatal("Uso: ./distiller <URL_DO_VIDEO>")
	}
	videoURL := os.Args[1]

	// Configurações do container vs local
	lmEndpoint := "http://host.docker.internal:1234/v1/chat/completions"
	modelName := "local-model" // Ou o nome específico carregado no LM Studio

	// 2. Inicialização da Infraestrutura (Dependências)
	llmClient := llm.NewLMStudioClient(lmEndpoint, modelName, 10*time.Minute) // Timeout longo para modelos locais
	dlClient := downloader.NewYtDlpClient("output/temp")
	sanitizer := textproc.NewRegexSanitizer()

	// 3. Inicialização do Caso de Uso (Injeção de Dependência)
	distiller := pipeline.NewDistillerUseCase(llmClient, dlClient, sanitizer)

	// 4. Execução
	fmt.Println("🚀 Iniciando Distiller v2.0 (Dynamic Agent Pipeline)...")
	start := time.Now()

	result, err := distiller.Run(videoURL)
	if err != nil {
		log.Fatalf("❌ Erro fatal na pipeline: %v", err)
	}

	// 5. Output
	outputFilename := fmt.Sprintf("output/knowledge_%d.md", time.Now().Unix())
	if err := os.WriteFile(outputFilename, []byte(result), 0644); err != nil {
		log.Fatalf("Erro ao salvar arquivo: %v", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("\n✅ Processo concluído em %s.\n📄 Resultado salvo em: %s\n", elapsed, outputFilename)
}
