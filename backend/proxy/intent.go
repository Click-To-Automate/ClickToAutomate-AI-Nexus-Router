package proxy

import (
	"strings"
	"ainexusrouter-core/providers"
)

// AnalyzeIntent examines the chat messages to determine the intent and complexity
func AnalyzeIntent(messages []interface{}) providers.IntentProfile {
	var fullTextBuilder strings.Builder

	// Concatenate all user messages to analyze intent
	for _, msg := range messages {
		msgMap, ok := msg.(map[string]interface{})
		if !ok {
			continue
		}
		
		role, _ := msgMap["role"].(string)
		if role == "user" {
			content, ok := msgMap["content"].(string)
			if ok {
				fullTextBuilder.WriteString(content)
				fullTextBuilder.WriteString(" ")
			}
		}
	}

	fullText := strings.ToLower(fullTextBuilder.String())
	length := len(fullText)
	approxTokens := length / 4 // 1 token is roughly 4 characters

	// 0. Massive Payload Detection -> Extremely high complexity
	if approxTokens > 10000 {
		return providers.IntentProfile{
			Complexity:  0.95,
			IsCoding:    false,
			IsReasoning: false,
		}
	}

	score := 0.5 // Baseline complexity

	// 1. Code Detection -> Increases complexity requirement
	isCoding := false
	codeKeywords := []string{
		"```", "function", "def ", "class ", "interface ", 
		"refactor", "bug", "error:", "compile", "script", "json", "python", "golang",
	}
	for _, kw := range codeKeywords {
		if strings.Contains(fullText, kw) {
			score += 0.3
			isCoding = true
			break
		}
	}

	// 2. Complex/Analysis/Reasoning Detection -> Increases complexity requirement
	isReasoning := false
	complexKeywords := []string{
		"analyze", "summarize", "explain the difference", 
		"compare", "architecture", "design", "essay", "math", "calculate", "solve", "logic", "puzzle",
	}
	for _, kw := range complexKeywords {
		if strings.Contains(fullText, kw) {
			score += 0.2
			isReasoning = true
			break
		}
	}

	if length > 2000 {
		score += 0.1
	}

	if score > 0.95 {
		score = 0.95 // Cap maximum complexity requirement slightly below 1.0
	}

	return providers.IntentProfile{
		Complexity:  score,
		IsCoding:    isCoding,
		IsReasoning: isReasoning,
	}
}
