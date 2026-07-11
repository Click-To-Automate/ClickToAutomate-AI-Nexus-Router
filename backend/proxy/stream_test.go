package proxy

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"ainexusrouter-core/config"
)

func TestHandleProxyRequest(w *testing.T) {
	// Mock upstream provider server
	mockUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure auth header is passed
		if r.Header.Get("Authorization") != "Bearer test-key" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		// Write two chunks simulating SSE
		fmt.Fprintf(w, "data: {\"choices\": [{\"delta\": {\"content\": \"hello\"}}]}\n\n")
		flusher.Flush()
		time.Sleep(10 * time.Millisecond) // Simulated latency
		fmt.Fprintf(w, "data: {\"choices\": [{\"delta\": {\"content\": \" world\"}}]}\n\n")
		flusher.Flush()
	}))
	defer mockUpstream.Close()

	// Mock the GlobalConfig for tests
	config.GlobalConfig = config.ProviderConfig{
		Providers: []config.ProviderDef{
			{
				ID:       "openai",
				Name:     "OpenAI",
				BaseURL:  "https://api.openai.com/v1",
				EnvKey:   "OPENAI_API_KEY",
				Prefixes: []string{"gpt-"},
			},
		},
	}

	// Override OPENAI_API_KEY for the test
	os.Setenv("OPENAI_API_KEY", "test-key")

	// We can't easily mock the registry.go's direct BaseURL mapping without changing the code structure,
	// but for testing, let's just make sure the proxy function doesn't crash on bad configurations.
	// In a real TDD setup, we would inject the provider URL.
	// Let's test the error case first.
	
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewBuffer([]byte(`{"model": "unknown-model"}`)))
	recorder := httptest.NewRecorder()

	err := HandleProxyRequest(recorder, req, []byte(`{"model": "unknown-model"}`), "unknown-model", true)
	if err == nil {
		w.Errorf("Expected error for unknown model, got nil")
	}

	if recorder.Code != http.StatusBadRequest {
		w.Errorf("Expected status 400, got %d", recorder.Code)
	}
}
