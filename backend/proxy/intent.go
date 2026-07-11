package proxy

import (
	"strings"
)

// AnalyzeTask examines the chat messages to determine the optimal routing tier.
// Returns "simple", "code", or "complex"
func AnalyzeTask(messages []interface{}) string {
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

	// 0. Massive Payload Detection (prevent 413 Payload Too Large on strict limits)
	if approxTokens > 10000 {
		return "massive"
	}

	// 1. Code Detection
	codeKeywords := []string{
		"```", "function", "def ", "class ", "interface ", 
		"refactor", "bug", "error:", "compile", "script",
	}
	for _, kw := range codeKeywords {
		if strings.Contains(fullText, kw) {
			return "code"
		}
	}

	// 2. Complex/Analysis Detection
	complexKeywords := []string{
		"analyze", "summarize", "explain the difference", 
		"compare", "architecture", "design", "essay",
	}
	for _, kw := range complexKeywords {
		if strings.Contains(fullText, kw) {
			return "complex"
		}
	}

	if length > 1000 {
		return "complex"
	}

	// 3. Default to Simple
	return "simple"
}
