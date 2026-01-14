package entity

type AgentRequest struct {
	Role        string  // Dynamic persona (e.g., "Senior SRE", "Executive Chef")
	Instruction string  // System prompt
	InputData   string  // Content to be processed
	Temperature float64 // Creativity vs. precision
}

// AgentConfig defines global settings (model name, endpoint)
type AgentConfig struct {
	ModelName string
	Endpoint  string
}
