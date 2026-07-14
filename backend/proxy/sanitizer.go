package proxy

import (
	"regexp"
)
// GlobalSanitizeRequest performs a pre-flight structural clean up of the payload.
// It explicitly removes proprietary fields injected by clients like Cursor
// before hitting the upstream provider.
// It also scans the messages to sanitize image structures and returns a flag if images are present.
func GlobalSanitizeRequest(payload map[string]interface{}) (hasImage bool) {
	// 1. Remove proprietary top-level fields
	delete(payload, "user")
	delete(payload, "providerOptions")

	hasImage = false

	// 2. Iterate and sanitize messages
	if messages, ok := payload["messages"].([]interface{}); ok {
		for j := len(messages) - 1; j >= 0; j-- {
			if msg, ok := messages[j].(map[string]interface{}); ok {
				if msg["role"] == "user" {
					// Check if content is already an array (contains images)
					if contentArray, ok := msg["content"].([]interface{}); ok {
						// Sanitize image_url objects to remove invalid fields
						sanitizedContent := sanitizeImageURL(contentArray)
						msg["content"] = sanitizedContent
						
						// Check for images and redact secrets in text
						for i, contentItem := range sanitizedContent {
							if item, ok := contentItem.(map[string]interface{}); ok {
								if item["type"] == "image_url" {
									hasImage = true
								} else if item["type"] == "text" {
									if textContent, ok := item["text"].(string); ok {
										item["text"] = RedactSecrets(textContent)
										sanitizedContent[i] = item
									}
								}
							}
						}
					} else if contentStr, ok := msg["content"].(string); ok {
						// If content is just a string, redact secrets from it directly
						msg["content"] = RedactSecrets(contentStr)
					}
				}
			}
		}
	}
	return hasImage
}

// sanitizeImageURL removes invalid fields from image_url objects to ensure compatibility with providers.
// Preserves the "url" and "detail" fields, and removes all other fields (e.g., "dimensions", "providerOptions").
func sanitizeImageURL(content []interface{}) []interface{} {
	var sanitizedContent []interface{}
	for _, item := range content {
		if itemMap, ok := item.(map[string]interface{}); ok {
			if itemMap["type"] == "image_url" {
				if imageURL, ok := itemMap["image_url"].(map[string]interface{}); ok {
					// Create a new image_url object with only the "url" and "detail" fields
					sanitizedImageURL := make(map[string]interface{})
					
					// Preserve the "url" field (must be a string)
					if url, ok := imageURL["url"].(string); ok {
						sanitizedImageURL["url"] = url
					}
					
					// Preserve the "detail" field if it exists (e.g., "high" or "low")
					if detail, ok := imageURL["detail"].(string); ok {
						sanitizedImageURL["detail"] = detail
					}
					
					newItemMap := make(map[string]interface{})
					newItemMap["type"] = "image_url"
					newItemMap["image_url"] = sanitizedImageURL
					itemMap = newItemMap
				}
			}
			sanitizedContent = append(sanitizedContent, itemMap)
		} else {
			sanitizedContent = append(sanitizedContent, item)
		}
	}
	return sanitizedContent
}

var secretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`sk-[a-zA-Z0-9]{32,}`),         // OpenAI
	regexp.MustCompile(`sk-ant-[a-zA-Z0-9_-]{32,}`),   // Anthropic
	regexp.MustCompile(`Bearer\s+[a-zA-Z0-9_\-\.]{32,}`), // Generic Bearer Tokens
}

// RedactSecrets scans a string for common credentials and replaces them with [REDACTED].
func RedactSecrets(text string) string {
	for _, pattern := range secretPatterns {
		text = pattern.ReplaceAllString(text, "[REDACTED]")
	}
	return text
}
