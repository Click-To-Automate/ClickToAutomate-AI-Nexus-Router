package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"ainexusrouter-core/proxy"
)

// AnthropicMessagesRequest matches Anthropic Messages API input
type AnthropicMessagesRequest struct {
	Model     string        `json:"model"`
	Messages  []interface{} `json:"messages"`
	System    interface{}   `json:"system,omitempty"`
	MaxTokens int           `json:"max_tokens,omitempty"`
	Stream    bool          `json:"stream,omitempty"`
}

// HandleAnthropicMessages maps downstream Anthropic protocol calls to standard OpenAI proxy formats
func HandleAnthropicMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-api-key, anthropic-version")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var anthropicReq AnthropicMessagesRequest
	if err := json.Unmarshal(body, &anthropicReq); err != nil {
		http.Error(w, "Invalid Anthropic JSON payload", http.StatusBadRequest)
		return
	}

	if anthropicReq.Model == "" {
		http.Error(w, "Model is required", http.StatusBadRequest)
		return
	}

	// Translate Anthropic format to OpenAI format
	openAIMessages := []map[string]interface{}{}

	// Handle System Prompt if present
	if anthropicReq.System != nil {
		var systemText string
		switch sys := anthropicReq.System.(type) {
		case string:
			systemText = sys
		case []interface{}:
			// Simple fallback if system is blocks
			for _, block := range sys {
				if blockMap, ok := block.(map[string]interface{}); ok {
					if t, ok := blockMap["text"].(string); ok {
						systemText += t
					}
				}
			}
		}
		if systemText != "" {
			openAIMessages = append(openAIMessages, map[string]interface{}{
				"role":    "system",
				"content": systemText,
			})
		}
	}

	// Map Messages
	for _, msg := range anthropicReq.Messages {
		if msgMap, ok := msg.(map[string]interface{}); ok {
			role := msgMap["role"]
			content := msgMap["content"]
			openAIMessages = append(openAIMessages, map[string]interface{}{
				"role":    role,
				"content": content,
			})
		}
	}

	// Build OpenAI payload
	openAIPayload := map[string]interface{}{
		"model":    anthropicReq.Model,
		"messages": openAIMessages,
		"stream":   anthropicReq.Stream,
	}
	if anthropicReq.MaxTokens > 0 {
		openAIPayload["max_tokens"] = anthropicReq.MaxTokens
	}

	openAIBody, err := json.Marshal(openAIPayload)
	if err != nil {
		http.Error(w, "Internal serialization error", http.StatusInternalServerError)
		return
	}

	log.Printf("[Anthropic-Adapter] Intercepted /v1/messages call -> model: %s | stream: %v\n", anthropicReq.Model, anthropicReq.Stream)

	// Intercept response to convert from OpenAI back to Anthropic
	if anthropicReq.Stream {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		// Create streaming translation wrapper
		flusher, _ := w.(http.Flusher)
		adapter := &anthropicStreamAdapter{
			writer:  w,
			flusher: flusher,
			model:   anthropicReq.Model,
		}

		// Initial start events
		adapter.SendEvent("message_start", map[string]interface{}{
			"type": "message_start",
			"message": map[string]interface{}{
				"id":            "msg_adapted_" + fmt.Sprintf("%d", len(anthropicReq.Model)),
				"type":          "message",
				"role":          "assistant",
				"content":       []interface{}{},
				"model":         anthropicReq.Model,
				"stop_reason":   nil,
				"stop_sequence": nil,
				"usage": map[string]interface{}{
					"input_tokens":  0,
					"output_tokens": 0,
				},
			},
		})
		adapter.SendEvent("content_block_start", map[string]interface{}{
			"type":  "content_block_start",
			"index": 0,
			"content_block": map[string]interface{}{
				"type": "text",
				"text": "",
			},
		})

		err = proxy.HandleProxyRequest(adapter, r, openAIBody, anthropicReq.Model, true)
		if err != nil {
			log.Printf("[Anthropic-Adapter] Streaming Proxy Error: %v\n", err)
		}

		// Final stop events
		adapter.SendEvent("content_block_stop", map[string]interface{}{
			"type":  "content_block_stop",
			"index": 0,
		})
		adapter.SendEvent("message_delta", map[string]interface{}{
			"type": "message_delta",
			"delta": map[string]interface{}{
				"stop_reason":   "end_turn",
				"stop_sequence": nil,
			},
			"usage": map[string]interface{}{
				"output_tokens": 0,
			},
		})
		adapter.SendEvent("message_stop", map[string]interface{}{
			"type": "message_stop",
		})
	} else {
		// Non-streaming interceptor
		recorder := &responseRecorder{
			header: http.Header{},
			buf:    &bytes.Buffer{},
		}

		err = proxy.HandleProxyRequest(recorder, r, openAIBody, anthropicReq.Model, false)
		if err != nil {
			log.Printf("[Anthropic-Adapter] Non-Streaming Proxy Error: %v\n", err)
			w.WriteHeader(recorder.code)
			w.Write(recorder.buf.Bytes())
			return
		}

		// Convert OpenAI static JSON to Anthropic static JSON
		var openAIResp struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}

		if err := json.Unmarshal(recorder.buf.Bytes(), &openAIResp); err != nil {
			// Write unmodified if it's not OpenAI standard (e.g. error message)
			w.WriteHeader(recorder.code)
			w.Write(recorder.buf.Bytes())
			return
		}

		textResponse := ""
		if len(openAIResp.Choices) > 0 {
			textResponse = openAIResp.Choices[0].Message.Content
		}

		anthropicResp := map[string]interface{}{
			"id":   "msg_adapted_static",
			"type": "message",
			"role": "assistant",
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": textResponse,
				},
			},
			"model":         anthropicReq.Model,
			"stop_reason":   "end_turn",
			"stop_sequence": nil,
			"usage": map[string]interface{}{
				"input_tokens":  0,
				"output_tokens": 0,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(anthropicResp)
	}
}

// adapter for SSE writes from OpenAI proxy to Anthropic protocol
type anthropicStreamAdapter struct {
	writer  http.ResponseWriter
	flusher http.Flusher
	model   string
	buffer  string
}

func (a *anthropicStreamAdapter) Header() http.Header {
	return a.writer.Header()
}

func (a *anthropicStreamAdapter) WriteHeader(statusCode int) {
	// Already written or handled
}

func (a *anthropicStreamAdapter) Write(b []byte) (int, error) {
	a.buffer += string(b)

	for {
		lineIdx := strings.Index(a.buffer, "\n")
		if lineIdx == -1 {
			break
		}

		line := strings.TrimSpace(a.buffer[:lineIdx])
		a.buffer = a.buffer[lineIdx+1:]

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "data:") {
			data := strings.TrimPrefix(line, "data:")
			data = strings.TrimSpace(data)

			if data == "[DONE]" {
				continue
			}

			var chunk struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}

			if err := json.Unmarshal([]byte(data), &chunk); err == nil {
				if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
					text := chunk.Choices[0].Delta.Content
					a.SendEvent("content_block_delta", map[string]interface{}{
						"type":  "content_block_delta",
						"index": 0,
						"delta": map[string]interface{}{
							"type": "text_delta",
							"text": text,
						},
					})
				}
			}
		}
	}

	return len(b), nil
}

func (a *anthropicStreamAdapter) SendEvent(event string, data interface{}) {
	payload, _ := json.Marshal(data)
	fmt.Fprintf(a.writer, "event: %s\ndata: %s\n\n", event, string(payload))
	if a.flusher != nil {
		a.flusher.Flush()
	}
}

type responseRecorder struct {
	header http.Header
	buf    *bytes.Buffer
	code   int
}

func (r *responseRecorder) Header() http.Header {
	return r.header
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	return r.buf.Write(b)
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.code = statusCode
}
