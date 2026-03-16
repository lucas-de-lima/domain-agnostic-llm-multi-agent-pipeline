package pipeline

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/domain/entity"
	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/domain/gateway"
)

type DistillerUseCase struct {
	llm        gateway.LLMProvider
	resolver   gateway.SourceResolver
	sanitizer  gateway.TextSanitizer
	spec       *entity.Specialization
}

func NewDistillerUseCase(llm gateway.LLMProvider, resolver gateway.SourceResolver, san gateway.TextSanitizer, spec *entity.Specialization) *DistillerUseCase {
	return &DistillerUseCase{
		llm:       llm,
		resolver:  resolver,
		sanitizer: san,
		spec:      spec,
	}
}

func (uc *DistillerUseCase) Run(ctx context.Context, sourceArg string) (string, error) {
	log.Println("Starting fetch and sanitization...")
	rawText, err := uc.resolver.FetchWithResolve(ctx, sourceArg)
	if err != nil {
		return "", fmt.Errorf("fetch: %w", err)
	}

	cleanText, err := uc.sanitizer.Clean(rawText)
	if err != nil {
		return "", fmt.Errorf("sanitization: %w", err)
	}

	if len(cleanText) < 50 {
		return "", fmt.Errorf("extracted text too short or empty")
	}

	log.Println("Agent 0: Identifying context and specialists...")
	dynamicContext, err := uc.identifyContext(cleanText)
	if err != nil {
		return "", fmt.Errorf("Agent 0: %w", err)
	}

	log.Printf("Context: %s | Level: %s", dynamicContext.MainSubject, dynamicContext.ComplexityLevel)
	log.Printf("Team: [1]%s [2]%s [3]%s", dynamicContext.ExpertRole1, dynamicContext.ExpertRole2, dynamicContext.ExpertRole3)

	log.Printf("Agent 1 (%s): Extracting structure...", dynamicContext.ExpertRole1)
	extractionJSON, err := uc.runExtraction(cleanText, dynamicContext)
	if err != nil {
		return "", fmt.Errorf("Agent 1: %w", err)
	}

	log.Printf("Agent 2 (%s): Writing draft...", dynamicContext.ExpertRole2)
	draftText, err := uc.runSynthesis(cleanText, extractionJSON, dynamicContext)
	if err != nil {
		return "", fmt.Errorf("Agent 2: %w", err)
	}

	log.Printf("Agent 3 (%s): Validating and refining...", dynamicContext.ExpertRole3)
	finalContent, err := uc.runAudit(draftText, cleanText, dynamicContext)
	if err != nil {
		return "", fmt.Errorf("Agent 3: %w", err)
	}

	return finalContent, nil
}

func (uc *DistillerUseCase) identifyContext(input string) (*entity.DynamicContext, error) {
	instruction := uc.spec.ContextPrompt
	resp, err := uc.llm.Call(entity.AgentRequest{
		Role:        "Senior Content Classification Analyst",
		Instruction: instruction,
		InputData:   sampleText(input, 2000),
		Temperature: 0.1,
	})
	if err != nil {
		return nil, err
	}

	var ctx entity.DynamicContext
	if err := json.Unmarshal([]byte(uc.cleanJSON(resp)), &ctx); err != nil {
		return nil, fmt.Errorf("Agent 0 JSON parse: %v | Raw: %s", err, resp)
	}
	return &ctx, nil
}

func (uc *DistillerUseCase) runExtraction(input string, ctx *entity.DynamicContext) (string, error) {
	instruction, err := uc.renderPrompt(uc.spec.ExtractPrompt, ctx)
	if err != nil {
		return "", err
	}
	return uc.llm.Call(entity.AgentRequest{
		Role:        ctx.ExpertRole1,
		Instruction: instruction,
		InputData:   input,
		Temperature: 0.1,
	})
}

func (uc *DistillerUseCase) runSynthesis(originalInput, extractionJSON string, ctx *entity.DynamicContext) (string, error) {
	inputComposto := fmt.Sprintf("--- STRUCTURED DATA ---\n%s\n\n--- ORIGINAL TEXT ---\n%s", extractionJSON, originalInput)
	instruction, err := uc.renderPrompt(uc.spec.SynthesizePrompt, ctx)
	if err != nil {
		return "", err
	}
	return uc.llm.Call(entity.AgentRequest{
		Role:        ctx.ExpertRole2,
		Instruction: instruction,
		InputData:   inputComposto,
		Temperature: 0.4,
	})
}

func (uc *DistillerUseCase) runAudit(draft, originalInput string, ctx *entity.DynamicContext) (string, error) {
	instruction, err := uc.renderPrompt(uc.spec.AuditPrompt, ctx)
	if err != nil {
		return "", err
	}
	inputWithOriginal := fmt.Sprintf("--- DRAFT ---\n%s\n\n--- ORIGINAL SOURCE ---\n%s", draft, originalInput)
	return uc.llm.Call(entity.AgentRequest{
		Role:        ctx.ExpertRole3,
		Instruction: instruction,
		InputData:   inputWithOriginal,
		Temperature: 0.1,
	})
}

func (uc *DistillerUseCase) renderPrompt(tpl string, ctx *entity.DynamicContext) (string, error) {
	t, err := template.New("").Parse(tpl)
	if err != nil {
		return "", fmt.Errorf("template: %w", err)
	}
	var b bytes.Buffer
	if err := t.Execute(&b, ctx); err != nil {
		return "", fmt.Errorf("template execute: %w", err)
	}
	return strings.TrimSpace(b.String()), nil
}

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
