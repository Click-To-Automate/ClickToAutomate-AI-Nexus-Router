package proxy

import (
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

// HandleProxyRequest proxies the request to the upstream provider
func HandleProxyRequest(w http.ResponseWriter, r *http.Request, body []byte, model string, stream bool) error {
	var routes []providers.SmartRoute

	// We unmarshal early to allow intent analysis
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Failed to parse JSON body", http.StatusBadRequest)
		return err
	}

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
		provider, err := providers.GetProviderForModel(model)
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
		
		// Sanitize payload: Many strict providers (like Mistral/Cerebras) reject the 'user' field that Cursor/Antigravity sends
		delete(payload, "user")

		// Context Compression Integration: Intelligent Hybrid Engine
		isVisionProvider := route.Provider.Name == "anthropic" || route.Provider.Name == "gemini" || route.Provider.Name == "mistral" || route.Provider.Name == "openai"

		hasImage := false
			
			if messages, ok := payload["messages"].([]interface{}); ok {
				for j := len(messages) - 1; j >= 0; j-- {
					if msg, ok := messages[j].(map[string]interface{}); ok {
						if msg["role"] == "user" {
							if contentStr, ok := msg["content"].(string); ok {
								// Run it through the intelligent chunker
								chunks := ChunkContent(contentStr)
								
								var newContent []map[string]interface{}
								modified := false
								
								for _, chunk := range chunks {
									// --- PHASE 1: TOON EXTRACTION PIPELINE ---
									// If it's a massive JSON block, we attempt TOON extraction regardless of provider
									if chunk.IsJSON && len(chunk.Text) > 2000 {
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
			
			// Ensure Mistral uses pixtral if we ended up rendering an image
			if hasImage && route.Provider.Name == "mistral" {
				payload["model"] = "pixtral-12b-2409"
			}

		newBody, err := json.Marshal(payload)
		if err != nil {
			http.Error(w, "Failed to marshal rewritten body", http.StatusInternalServerError)
			return err
		}

		var targetURL string
		if route.Provider.Name == "openai" {
			targetURL = route.Provider.BaseURL + "/chat/completions"
		} else {
			targetURL = route.Provider.BaseURL + "/chat/completions"
		}

		req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(newBody))
		if err != nil {
			lastErr = err
			continue
		}

		req.Header.Set("Content-Type", "application/json")

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
		// We must put the bodyBytes back into a reader so we can stream it back!
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Success! Log telemetry to database
		db.IncrementUsage(route.Provider.Name, totalSavedTokens)

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
