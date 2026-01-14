package pipeline

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/domain/entity"
	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/domain/gateway"
)

type DistillerUseCase struct {
	llm        gateway.LLMProvider
	downloader gateway.ContentDownloader
	sanitizer  gateway.TextSanitizer
}

func NewDistillerUseCase(llm gateway.LLMProvider, dl gateway.ContentDownloader, san gateway.TextSanitizer) *DistillerUseCase {
	return &DistillerUseCase{
		llm:        llm,
		downloader: dl,
		sanitizer:  san,
	}
}

func (uc *DistillerUseCase) Run(sourceURL string) (string, error) {
	// 1. Acquisition and cleaning
	log.Println("📥 Starting download and sanitization...")
	rawText, err := uc.downloader.Download(sourceURL)
	if err != nil {
		return "", fmt.Errorf("download error: %w", err)
	}

	cleanText, err := uc.sanitizer.Clean(rawText)
	if err != nil {
		return "", fmt.Errorf("sanitization error: %w", err)
	}

	// Safety check to avoid processing empty or trivial content
	if len(cleanText) < 50 {
		return "", fmt.Errorf("extracted text too short or empty")
	}

	// 2. Agent 0: Context Architect (Router)
	log.Println("🔮 Agent 0: Identifying context and specialists...")
	dynamicContext, err := uc.identifyContext(cleanText)
	if err != nil {
		return "", fmt.Errorf("Agent 0 failed: %w", err)
	}

	log.Printf("🎯 Context: %s | Level: %s\n", dynamicContext.MainSubject, dynamicContext.ComplexityLevel)
	log.Printf("👥 Team Called: [1]%s [2]%s [3]%s\n",
		dynamicContext.ExpertRole1, dynamicContext.ExpertRole2, dynamicContext.ExpertRole3)

	// 3. Agent 1: Structural Extractor (The Architect)
	log.Printf("🕵️  Agent 1 (%s): Extracting structure...", dynamicContext.ExpertRole1)
	extractionJSON, err := uc.runExtraction(cleanText, dynamicContext)
	if err != nil {
		return "", fmt.Errorf("Agent 1 failed: %w", err)
	}

	// 4. Agent 2: The Synthesizer (The Writer)
	log.Printf("✍️  Agent 2 (%s): Writing draft...", dynamicContext.ExpertRole2)
	draftText, err := uc.runSynthesis(cleanText, extractionJSON, dynamicContext)
	if err != nil {
		return "", fmt.Errorf("Agent 2 failed: %w", err)
	}

	// 5. Agent 3: The Auditor (The Critic)
	log.Printf("🧐 Agent 3 (%s): Validating and refining...", dynamicContext.ExpertRole3)
	finalContent, err := uc.runAudit(draftText, cleanText, dynamicContext)
	if err != nil {
		return "", fmt.Errorf("Agent 3 failed: %w", err)
	}

	return finalContent, nil
}

// --- Métodos Privados de Orquestração ---

func (uc *DistillerUseCase) identifyContext(input string) (*entity.DynamicContext, error) {
	instruction := `
	Analyze the provided text (input). Your goal is to classify the knowledge domain and determine the best expert roles to work on it.

	Return ONLY a JSON with the following structure:
	{
		"main_subject": "The main subject (e.g., Quantum Physics, French Cuisine, DevOps)",
		"complexity_level": "Technical level of the text (Beginner, Intermediate, Advanced)",
		"expert_role_1": "Technical role name for data extraction (e.g., Theoretical Physicist, Saucier, SRE Engineer)",
		"expert_role_2": "Role name for writing educational content (e.g., University Professor, Cookbook Editor, Tech Lead)",
		"expert_role_3": "Role name for auditing mistakes (e.g., Scientific Reviewer, Food Critic, Security Auditor)",
		"target_audience": "Ideal target audience for the summary"
	}
	`
	resp, err := uc.llm.Call(entity.AgentRequest{
		Role:        "Senior Content Classification Analyst",
		Instruction: instruction,
		InputData:   sampleText(input, 2000), // Send only the first 2k chars to save/accelerate classification
		Temperature: 0.1,
	})
	if err != nil {
		return nil, err
	}

	var ctx entity.DynamicContext
	if err := json.Unmarshal([]byte(uc.cleanJSON(resp)), &ctx); err != nil {
		return nil, fmt.Errorf("Agent 0 JSON parse error: %v | Raw: %s", err, resp)
	}
	return &ctx, nil
}

func (uc *DistillerUseCase) runExtraction(input string, ctx *entity.DynamicContext) (string, error) {
	instruction := fmt.Sprintf(`
	You are a %s.
	Your task is to analyze the raw text and extract the vital technical data about %s.
	Ignore irrelevant conversation. Focus on logical structure, facts, numbers, ingredients or commands.

	Return a generic JSON that represents the "truth" of this content.
	Example generic structure (adapt to domain):
	{
		"key_concepts": [],
		"procedural_steps": [],
		"required_tools": [],
		"critical_alerts": []
	}
	`, ctx.ExpertRole1, ctx.MainSubject)

	return uc.llm.Call(entity.AgentRequest{
		Role:        ctx.ExpertRole1,
		Instruction: instruction,
		InputData:   input,
		Temperature: 0.1,
	})
}

func (uc *DistillerUseCase) runSynthesis(originalInput, extractionJSON string, ctx *entity.DynamicContext) (string, error) {
	inputComposto := fmt.Sprintf("--- STRUCTURED DATA ---\n%s\n\n--- ORIGINAL TEXT ---\n%s", extractionJSON, originalInput)

	instruction := fmt.Sprintf(`
	You are a %s writing for %s.
	Use the STRUCTURED DATA as the source of truth and the ORIGINAL TEXT for nuance.

	Goal: Produce a final document in Markdown, professional and highly educational about %s.

	Guidelines:
	1. Correct incorrect or confusing jargon from the original text.
	2. Organize into Title, Summary, Sections and Conclusion.
	3. Use rich formatting (bold, lists).
	`, ctx.ExpertRole2, ctx.TargetAudience, ctx.MainSubject)

	return uc.llm.Call(entity.AgentRequest{
		Role:        ctx.ExpertRole2,
		Instruction: instruction,
		InputData:   inputComposto,
		Temperature: 0.4, // Um pouco mais criativo para escrever bem
	})
}

func (uc *DistillerUseCase) runAudit(draft, originalInput string, ctx *entity.DynamicContext) (string, error) {
	instruction := fmt.Sprintf(`
	You are a %s. Your role is to ensure technical and logical integrity.
	Review the Draft below.

	Check:
	1. If there are hallucinations (things not present in the original or technically impossible in %s).
	2. If the language is appropriate for %s.
	3. If the step-by-step makes logical/physical sense.

	If it's perfect, return the original Draft.
	If there are issues, rewrite the problematic section while keeping Markdown style.
	`, ctx.ExpertRole3, ctx.MainSubject, ctx.TargetAudience)

	return uc.llm.Call(entity.AgentRequest{
		Role:        ctx.ExpertRole3,
		Instruction: instruction,
		InputData:   draft,
		Temperature: 0.1,
	})
}

// Helpers

func (uc *DistillerUseCase) cleanJSON(raw string) string {
	raw = strings.ReplaceAll(raw, "```json", "")
	raw = strings.ReplaceAll(raw, "```", "")
	return strings.TrimSpace(raw)
}

func sampleText(text string, limit int) string {
	if len(text) > limit {
		return text[:limit] + "..."
	}
	return text
}
