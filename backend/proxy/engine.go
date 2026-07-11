package proxy

import (
	"encoding/json"
	"math"
	"strings"
)

// Chunk represents a segmented piece of the user's prompt
type Chunk struct {
	Text             string
	FidelityRequired bool
	IsJSON           bool
}

// ChunkContent splits a user message into chunks. Code blocks (```...```) 
// and inline code are marked as FidelityRequired=true.
func ChunkContent(text string) []Chunk {
	var chunks []Chunk

	// Simple heuristic: split by "```" to find markdown code blocks
	parts := strings.Split(text, "```")
	
	for i, part := range parts {
		if part == "" {
			continue
		}
		
		// Every odd index in a split by "```" represents content INSIDE a code block
		isCodeBlock := (i % 2) != 0
		
		if isCodeBlock {
			chunks = append(chunks, Chunk{
				Text:             "```" + part + "```", // Re-add the markdown ticks so LLM formats it nicely
				FidelityRequired: true,
				IsJSON:           false,
			})
		} else {
			requiresFidelity := false
			isJSON := false
			
			// Securely detect if this string block is valid JSON!
			cleanPart := strings.TrimSpace(part)
			if (strings.HasPrefix(cleanPart, "{") && strings.HasSuffix(cleanPart, "}")) || 
			   (strings.HasPrefix(cleanPart, "[") && strings.HasSuffix(cleanPart, "]")) {
				if json.Valid([]byte(cleanPart)) {
					isJSON = true
					requiresFidelity = true // JSON implies exact data is needed
				}
			}

			// Also flag things that just look like code objects but aren't purely valid JSON
			if !isJSON && strings.Contains(part, "{") && strings.Contains(part, "}") {
				requiresFidelity = true
			}
			
			chunks = append(chunks, Chunk{
				Text:             part,
				FidelityRequired: requiresFidelity,
				IsJSON:           isJSON,
			})
		}
	}
	
	return chunks
}

// EstimateTextCost returns the approximate number of text tokens
func EstimateTextCost(chars int) int {
	return chars / 4
}

// EstimateVisionCost returns the approximate number of vision tokens
// based on Anthropic/OpenAI pricing for 1568px width images.
func EstimateVisionCost(chars int) int {
	// 1568px width / 7px char width = 224 chars per line
	charsPerLine := 224
	
	lines := float64(chars) / float64(charsPerLine)
	height := int(math.Ceil(lines)) * 14
	
	// If height is tiny, it might just be 1 tile
	if height == 0 {
		return 85
	}
	
	// Vision tile logic: 512x512 tiles.
	// Width 1568 -> ceil(1568/512) = 4 columns of tiles
	// Height -> ceil(height/512) rows of tiles
	cols := 4
	rows := int(math.Ceil(float64(height) / 512.0))
	
	tiles := cols * rows
	
	// Base cost + (tiles * tileCost)
	// OpenAI/Anthropic typically charge 85 base + 170 per tile, or similar variations.
	// Using a generic 170 per tile for safety margin.
	return 85 + (170 * tiles)
}
