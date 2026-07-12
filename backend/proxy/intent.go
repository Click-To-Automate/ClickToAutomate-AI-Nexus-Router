package proxy

import (
	"strings"
)

// AnalyzeComplexity examines the chat messages to determine the mathematical complexity score (0.0 to 1.0)
func AnalyzeComplexity(messages []interface{}) float64 {
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
		return 0.95
	}

	score := 0.5 // Baseline complexity

	// 1. Code Detection -> Increases complexity requirement
	codeKeywords := []string{
		"```", "function", "def ", "class ", "interface ", 
		"refactor", "bug", "error:", "compile", "script",
	}
	for _, kw := range codeKeywords {
		if strings.Contains(fullText, kw) {
			score += 0.3
			break
		}
	}

	// 2. Complex/Analysis Detection -> Increases complexity requirement
	complexKeywords := []string{
		"analyze", "summarize", "explain the difference", 
		"compare", "architecture", "design", "essay",
	}
	for _, kw := range complexKeywords {
		if strings.Contains(fullText, kw) {
			score += 0.2
			break
		}
	}

	if length > 2000 {
		score += 0.1
	}

	if score > 0.95 {
		score = 0.95 // Cap maximum complexity requirement slightly below 1.0
	}

	return score
}
