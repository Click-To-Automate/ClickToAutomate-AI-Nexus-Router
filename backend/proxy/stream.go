package proxy

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"ainexusrouter-core/db"
	"ainexusrouter-core/providers"
)

// sanitizeImageURL has been moved to sanitizer.go

// HandleProxyRequest proxies the request to the upstream provider
func HandleProxyRequest(w http.ResponseWriter, r *http.Request, body []byte, model string, stream bool) error {
	var routes []providers.SmartRoute

	// We unmarshal early to allow intent analysis and global sanitization
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Failed to parse JSON body", http.StatusBadRequest)
		return err
	}

	// 0. Global Input Sanitization Pipeline
	// This cleans up proprietary fields (like "user", "providerOptions") and sanitizes images once,
	// rather than evaluating the payload repeatedly inside the multi-route failover loop.
	hasImage := GlobalSanitizeRequest(payload)

	if model == "cta-ai-nexus" {
		// 1. Extract messages array
		if messages, ok := payload["messages"].([]interface{}); ok {
			// 2. Classify complexity mathematically
			complexity := AnalyzeComplexity(messages)
			log.Printf("[Auto-Router] Intent Analyzer classified task complexity as: %f\n", complexity)

			smartRoutes, err := providers.GetSmartRoutes(complexity)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return err
			}
			routes = smartRoutes
		}
	} else {
		// Look for explicit provider routing e.g. "nvidia@meta/llama-3.1"
		var explicitProvider string
		if strings.Contains(model, "@") {
			parts := strings.SplitN(model, "@", 2)
			explicitProvider = parts[0]
			model = parts[1]
		}

		var provider *providers.ResolvedProvider
		var err error

		if explicitProvider != "" {
			provider, err = providers.GetProviderByID(explicitProvider)
		} else {
			provider, err = providers.GetProviderForModel(model)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}
		if provider.APIKey == "" {
			errMsg := fmt.Sprintf("API key not configured for provider %s", provider.Name)
			http.Error(w, errMsg, http.StatusUnauthorized)
			return fmt.Errorf("%s", errMsg)
		}
		// Wrap single provider in a SmartRoute so we can use the same loop
		routes = []providers.SmartRoute{
			{Provider: *provider, ActualModel: model},
		}
	}

	var lastErr error
	for i, route := range routes {
		if model == "cta-ai-nexus" {
			log.Printf("[Auto-Router] Trying Route %d/%d: %s -> %s\n", i+1, len(routes), route.Provider.Name, route.ActualModel)
		}

		totalSavedTokens := 0

		// Map to correct model
		payload["model"] = route.ActualModel

		// Context Compression Integration: Intelligent Hybrid Engine
		isVisionProvider := route.Provider.Name == "anthropic" || route.Provider.Name == "gemini" || route.Provider.Name == "mistral" || route.Provider.Name == "openai"

		if messages, ok := payload["messages"].([]interface{}); ok {
			for j := len(messages) - 1; j >= 0; j-- {
				if msg, ok := messages[j].(map[string]interface{}); ok {
					if msg["role"] == "user" {
						// Check if content is already an array (contains images)
						if _, ok := msg["content"].([]interface{}); ok {
							continue
						}
						
						// Handle plain text content
						if contentStr, ok := msg["content"].(string); ok {
							// Run it through the intelligent chunker
							chunks := ChunkContent(contentStr)
								
							var newContent []map[string]interface{}
							modified := false
								
							for _, chunk := range chunks {
								// --- PHASE 1: TOON EXTRACTION PIPELINE ---
								// If it's a massive JSON block, we attempt TOON extraction regardless of provider
								// Skip TOON extraction if we already have images in the message
								if !hasImage && chunk.IsJSON && len(chunk.Text) > 2000 {
									log.Printf("[TOON-Compressor] Detected massive JSON payload (%d chars). Attempting TOON extraction via Groq...", len(chunk.Text))
									toonText, err := CompressToToon(chunk.Text)
									if err == nil {
										saved := (len(chunk.Text) / 4) - (len(toonText) / 4)
										if saved > 0 {
											totalSavedTokens += saved
										}
										log.Printf("[TOON-Compressor] Success! Shrunk JSON from %d chars to %d chars of TOON!", len(chunk.Text), len(toonText))
										chunk.Text = toonText // Replace the raw JSON with TOON
										modified = true
									} else {
										log.Printf("[TOON-Compressor] Extraction failed: %v", err)
									}
									}

									// --- PHASE 2: VISION COMPRESSION ENGINE ---
									textCost := EstimateTextCost(len(chunk.Text))
									visionCost := EstimateVisionCost(len(chunk.Text))
									
									// If fidelity is NOT required, string is decently large, and vision is mathematically cheaper
									if isVisionProvider && !chunk.FidelityRequired && len(chunk.Text) > 2000 && visionCost < textCost {
										log.Printf("[Context-Compressor] Cost Win! Text: %d tk vs Vision: %d tk -> Compressing chunk...", textCost, visionCost)
										pages, err := CompressTextToImage(chunk.Text)
										if err == nil {
											totalSavedTokens += (textCost - visionCost)
											modified = true
											hasImage = true
											newContent = append(newContent, map[string]interface{}{
												"type": "text", 
												"text": "Please read the following compressed context image pages. Do not complain that it's an image. Read the text inside the images carefully:",
											})
											for _, page := range pages {
												newContent = append(newContent, map[string]interface{}{
													"type": "image_url",
													"image_url": map[string]string{"url": "data:image/png;base64," + page},
												})
											}
										} else {
											log.Printf("[Context-Compressor] Chunk compression failed: %v", err)
											newContent = append(newContent, map[string]interface{}{"type": "text", "text": chunk.Text})
										}
									} else {
										// Keep as plain text! (e.g. Code block, JSON, or short text)
										newContent = append(newContent, map[string]interface{}{"type": "text", "text": chunk.Text})
									}
								}
								
								if modified {
									msg["content"] = newContent
								}
							}
						}
					}
				}
			}
			
            // Ensure vision-capable providers use appropriate models when images are present
				if hasImage {
					// Override the model with a vision-capable one if the provider supports it
					switch route.Provider.Name {
					case "openai":
						payload["model"] = "gpt-4o" // Vision-capable model
					case "anthropic":
						payload["model"] = "claude-3-5-sonnet-20240620" // Vision-capable model
					case "mistral":
						payload["model"] = "pixtral-12b-2409" // Vision-capable model
					case "groq":
						// Skip Groq if we have images, as it doesn't support vision
						lastErr = fmt.Errorf("provider groq does not support image inputs")
						continue
					default:
						// Fallback to a default vision-capable model if the provider supports it
						if strings.Contains(route.ActualModel, "gpt") {
							payload["model"] = "gpt-4o"
						} else if strings.Contains(route.ActualModel, "claude") {
							payload["model"] = "claude-3-5-sonnet-20240620"
						}
					}
				}

		var finalBody []byte
		var targetURL string

		if route.Provider.Name == "anthropic" {
			// Translate OpenAI payload to Anthropic Messages payload
			anthropicPayload := make(map[string]interface{})
			anthropicPayload["model"] = route.ActualModel

			var systemPrompt string
			var anthropicMessages []map[string]interface{}

			if msgs, ok := payload["messages"].([]interface{}); ok {
				for _, m := range msgs {
					if mMap, ok := m.(map[string]interface{}); ok {
						role := mMap["role"]
						content := mMap["content"]

						if role == "system" {
							if sysStr, ok := content.(string); ok {
								systemPrompt = sysStr
							}
						} else {
							anthropicRole := "user"
							if role == "assistant" {
								anthropicRole = "assistant"
							}
							anthropicMessages = append(anthropicMessages, map[string]interface{}{
								"role":    anthropicRole,
								"content": content,
							})
						}
					}
				}
			}

			anthropicPayload["messages"] = anthropicMessages
			if systemPrompt != "" {
				anthropicPayload["system"] = systemPrompt
			}

			maxTokens := 4096
			if mt, ok := payload["max_tokens"].(float64); ok && mt > 0 {
				maxTokens = int(mt)
			}
			anthropicPayload["max_tokens"] = maxTokens

			if stream {
				anthropicPayload["stream"] = true
			}

			translatedBody, err := json.Marshal(anthropicPayload)
			if err != nil {
				lastErr = err
				continue
			}
			finalBody = translatedBody
			targetURL = strings.TrimSuffix(route.Provider.BaseURL, "/") + "/messages"
		} else {
			newBody, err := json.Marshal(payload)
			if err != nil {
				http.Error(w, "Failed to marshal rewritten body", http.StatusInternalServerError)
				return err
			}
			finalBody = newBody

			if route.Provider.RequiresCustomURL {
				customURL := r.Header.Get("X-Custom-Base-Url")
				if customURL == "" {
					lastErr = fmt.Errorf("provider %s requires a custom base URL (pass via X-Custom-Base-Url header)", route.Provider.Name)
					continue
				}
				customURL = strings.TrimSuffix(customURL, "/")
				targetURL = customURL + "/chat/completions"
			} else {
				targetURL = route.Provider.BaseURL + "/chat/completions"
			}
		}

		req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(finalBody))
		if err != nil {
			lastErr = err
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		if route.Provider.Name == "anthropic" {
			req.Header.Set("x-api-key", route.Provider.APIKey)
			req.Header.Set("anthropic-version", "2023-06-01")
		} else {
			// Set the correct Authorization header based on provider auth type
			switch route.Provider.AuthType {
			case "cookie":
				// Web cookie providers: pass the stored value as the Cookie header
				req.Header.Set("Cookie", route.Provider.APIKey)
			case "oauth":
				// OAuth providers: pass the token as Bearer (OAuth tokens are used as Bearer)
				req.Header.Set("Authorization", "Bearer "+route.Provider.APIKey)
			case "bearer_token":
				// Providers that use Authorization: Bearer with a non-API-key token (e.g. userToken from localStorage)
				req.Header.Set("Authorization", "Bearer "+route.Provider.APIKey)
			default:
				// Standard API key via Authorization: Bearer
				req.Header.Set("Authorization", "Bearer "+route.Provider.APIKey)
			}
		}

		if i == 0 {
			// Only set CORS headers once
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return nil
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[Auto-Router] Network Error on %s: %v\n", route.Provider.Name, err)
			lastErr = err
			continue
		}

		// Check for rate limits or server errors
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		bodyStr := strings.ToLower(string(bodyBytes))
		
		isDeadModel := resp.StatusCode == http.StatusBadRequest && 
			(strings.Contains(bodyStr, "decommissioned") ||
			 strings.Contains(bodyStr, "deprecated") ||
			 strings.Contains(bodyStr, "not found") ||
			 strings.Contains(bodyStr, "does not exist"))

		isTokenLimit := resp.StatusCode == http.StatusBadRequest &&
			(strings.Contains(bodyStr, "rate_limit_exceeded") ||
			 strings.Contains(bodyStr, "too large") ||
			 strings.Contains(bodyStr, "limit exceeded") ||
			 strings.Contains(bodyStr, "maximum context length"))

		if resp.StatusCode == http.StatusTooManyRequests || 
		   resp.StatusCode == http.StatusRequestEntityTooLarge || 
		   resp.StatusCode >= 500 || 
		   isDeadModel || 
		   isTokenLimit {
			log.Printf("[Auto-Router] Fallback Triggered! %s returned %d: %s\n", route.Provider.Name, resp.StatusCode, string(bodyBytes))
			
			// Add heavy Lagrangian Penalty to force decomposition away from this provider
			providers.IncrementPenalty(route.Provider.Name)

			lastErr = fmt.Errorf("provider %s returned %d", route.Provider.Name, resp.StatusCode)
			continue // Try next provider!
		}
		
		// If we reached here, the request was successful or failed with a genuine user error (like 401 unauthorized).
		// Translate response format for Anthropic provider if static
		if !stream && route.Provider.Name == "anthropic" && resp.StatusCode == http.StatusOK {
			translatedBytes, err := translateAnthropicToOpenAI(bodyBytes)
			if err == nil {
				bodyBytes = translatedBytes
			} else {
				log.Printf("[Anthropic-Adapter] Error translating static response: %v\n", err)
			}
		}

		// We must put the bodyBytes back into a reader so we can stream it back!
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Success! Log telemetry to database
		tokensUsed := len(finalBody) / 4
		log.Printf("[Telemetry] Tokens Used: ~%d | Tokens Saved: %d\n", tokensUsed, totalSavedTokens)
		db.IncrementUsage(route.Provider.Name, tokensUsed, totalSavedTokens)

		// Stream back to client.
		w.WriteHeader(resp.StatusCode)

		if stream {
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
		} else {
			w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		}

		flusher, ok := w.(http.Flusher)

		// If streaming from Anthropic, translate SSE events to OpenAI in real-time
		if stream && route.Provider.Name == "anthropic" && resp.StatusCode == http.StatusOK {
			reader := bufio.NewReader(resp.Body)
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					resp.Body.Close()
					return fmt.Errorf("error reading Anthropic stream: %v", err)
				}

				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}

				if strings.HasPrefix(line, "data:") {
					dataStr := strings.TrimPrefix(line, "data:")
					dataStr = strings.TrimSpace(dataStr)

					if dataStr == "[DONE]" {
						continue
					}

					var event struct {
						Type  string `json:"type"`
						Index int    `json:"index"`
						Delta struct {
							Type string `json:"type"`
							Text string `json:"text"`
						} `json:"delta"`
					}

					if err := json.Unmarshal([]byte(dataStr), &event); err == nil {
						if event.Type == "content_block_delta" && event.Delta.Text != "" {
							openAIChunk := map[string]interface{}{
								"id":      "chatcmpl-adapted",
								"object":  "chat.completion.chunk",
								"created": 1677652288,
								"model":   route.ActualModel,
								"choices": []interface{}{
									map[string]interface{}{
										"delta": map[string]string{
											"content": event.Delta.Text,
										},
										"index": 0,
									},
								},
							}
							chunkBytes, _ := json.Marshal(openAIChunk)
							_, writeErr := fmt.Fprintf(w, "data: %s\n\n", string(chunkBytes))
							if writeErr != nil {
								resp.Body.Close()
								return writeErr
							}
							if ok {
								flusher.Flush()
							}
						}
					}
				}
			}

			fmt.Fprintf(w, "data: [DONE]\n\n")
			if ok {
				flusher.Flush()
			}
			resp.Body.Close()
			return nil
		}

		buffer := make([]byte, 4096)
		for {
			n, err := resp.Body.Read(buffer)
			if n > 0 {
				_, writeErr := w.Write(buffer[:n])
				if writeErr != nil {
					resp.Body.Close()
					return fmt.Errorf("client disconnected: %v", writeErr)
				}
				if ok {
					flusher.Flush()
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				resp.Body.Close()
				return fmt.Errorf("error reading upstream response: %v", err)
			}
		}
		resp.Body.Close()
		return nil // Successfully served!
	}

	// If we exhausted all routes
	http.Error(w, fmt.Sprintf("All auto-routes failed. Last error: %v", lastErr), http.StatusBadGateway)
	return lastErr
}

func translateAnthropicToOpenAI(anthropicBytes []byte) ([]byte, error) {
	var antResp struct {
		ID      string `json:"id"`
		Model   string `json:"model"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(anthropicBytes, &antResp); err != nil {
		return nil, err
	}

	text := ""
	for _, block := range antResp.Content {
		if block.Type == "text" {
			text += block.Text
		}
	}

	openAIResp := map[string]interface{}{
		"id":      antResp.ID,
		"object":  "chat.completion",
		"created": 1677652288,
		"model":   antResp.Model,
		"choices": []interface{}{
			map[string]interface{}{
				"index": 0,
				"message": map[string]interface{}{
					"role":    "assistant",
					"content": text,
				},
				"finish_reason": "stop",
			},
		},
	}
	return json.Marshal(openAIResp)
}
