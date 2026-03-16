package entity

// Specialization is a use-case preset: prompts (and optional output hints) for the 4 pipeline stages.
// Community contributions add new specializations (e.g. recipes, programming-lessons) via PR.
type Specialization struct {
	ID          string // e.g. "default", "recipes"
	Name        string
	Description string
	// Prompt templates for each stage. May contain placeholders like {{.MainSubject}}, {{.TargetAudience}}.
	ContextPrompt   string // Agent 0: classify domain and expert roles
	ExtractPrompt   string // Agent 1: extract structured data
	SynthesizePrompt string // Agent 2: write document from structure + original
	AuditPrompt     string // Agent 3: validate draft against original
}
