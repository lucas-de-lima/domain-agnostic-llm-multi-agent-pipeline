package entity

// DynamicContext represents Agent 0's initial understanding of the content.
// It defines which expert roles are required for the next steps.
type DynamicContext struct {
	MainSubject     string `json:"main_subject"`     // e.g., "Software Engineering", "Italian Cuisine"
	ComplexityLevel string `json:"complexity_level"` // e.g., "Senior", "Beginner", "Technical"
	ExpertRole1     string `json:"expert_role_1"`    // The information "Architect"
	ExpertRole2     string `json:"expert_role_2"`    // The "Writer"
	ExpertRole3     string `json:"expert_role_3"`    // The "Auditor"
	TargetAudience  string `json:"target_audience"`  // Who the content is for
}

// ExtractionResult is the generic output from Agent 1.
// We use map[string]interface{} because the structure depends on the domain.
type ExtractionResult map[string]interface{}
