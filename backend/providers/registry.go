package providers

import (
	"fmt"
	"os"
	"strings"

	"ainexusrouter-core/config"
	"ainexusrouter-core/db"
	"ainexusrouter-core/discovery"
)

type ResolvedProvider struct {
	Name    string
	BaseURL string
	APIKey  string
}

// GetProviderForModel determines which provider to use based on the model name
// using the dynamic providers configuration.
func GetProviderForModel(model string) (*ResolvedProvider, error) {
	for _, p := range config.GlobalConfig.Providers {
		for _, prefix := range p.Prefixes {
			if strings.HasPrefix(model, prefix) {
				// 1. Check local SQLite DB first
				key := db.GetKey(p.ID)
				
				// 2. Fallback to Environment Variable
				if key == "" && p.EnvKey != "" {
					key = os.Getenv(p.EnvKey)
				}
				
				// Strip aggregator prefixes if needed (e.g. "openrouter/llama" -> "llama")
				// We will handle actual model translation in the proxy layer if needed.
				return &ResolvedProvider{
					Name:    p.ID,
					BaseURL: p.BaseURL,
					APIKey:  key,
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("no provider configured for model: %s", model)
}

// SmartRoute represents a resolved provider and the actual model payload we should rewrite to
type SmartRoute struct {
	Provider    ResolvedProvider
	ActualModel string
}

// GetSmartRoutes builds an ordered fallback chain of available providers based on task intent
func GetSmartRoutes(taskType string) ([]SmartRoute, error) {
	// The desired fallback order based on intent
	var desiredOrder []struct {
		providerID string
		keyword    string
		fallback   string
	}

	switch taskType {
	case "massive":
		desiredOrder = []struct{ providerID, keyword, fallback string }{
			{"mistral", "large", "mistral-large-latest"},
			{"gemini", "pro", "gemini-1.5-pro"},
			{"anthropic", "sonnet", "claude-3-5-sonnet-20240620"},
		}
	case "code":
		desiredOrder = []struct{ providerID, keyword, fallback string }{
			{"groq", "70b", "llama3-70b-8192"},
			{"anthropic", "sonnet", "claude-3-5-sonnet-20240620"},
			{"openai", "4o-mini", "gpt-4o-mini"},
			{"mistral", "large", "mistral-large-latest"},
			{"gemini", "pro", "gemini-1.5-pro"},
		}
	case "simple":
		desiredOrder = []struct{ providerID, keyword, fallback string }{
			{"cerebras", "8b", "llama3.1-8b"},
			{"mistral", "large", "mistral-large-latest"},
			{"gemini", "flash", "gemini-1.5-flash"},
		}
	case "complex":
		desiredOrder = []struct{ providerID, keyword, fallback string }{
			{"mistral", "large", "mistral-large-latest"},
			{"gemini", "pro", "gemini-1.5-pro"},
			{"groq", "70b", "llama3-70b-8192"},
			{"anthropic", "sonnet", "claude-3-5-sonnet-20240620"},
		}
	default:
		desiredOrder = []struct{ providerID, keyword, fallback string }{
			{"mistral", "large", "mistral-large-latest"},
			{"cerebras", "8b", "llama3.1-8b"},
			{"groq", "70b", "llama3-70b-8192"},
			{"gemini", "flash", "gemini-1.5-flash"},
			{"openai", "4o-mini", "gpt-4o-mini"},
			{"anthropic", "sonnet", "claude-3-5-sonnet-20240620"},
		}
	}

	var routes []SmartRoute

	for _, req := range desiredOrder {
		// Find provider config
		for _, p := range config.GlobalConfig.Providers {
			if p.ID == req.providerID {
				key := db.GetKey(p.ID)
				if key == "" && p.EnvKey != "" {
					key = os.Getenv(p.EnvKey)
				}
				
				// Only add to chain if we actually have a key!
				if key != "" {
					// Dynamically pick the best live model discovered at startup!
					bestModel := discovery.GetBestModel(p.ID, req.keyword, req.fallback)
					
					routes = append(routes, SmartRoute{
						Provider: ResolvedProvider{
							Name:    p.ID,
							BaseURL: p.BaseURL,
							APIKey:  key,
						},
						ActualModel: bestModel,
					})
				}
				break
			}
		}
	}

	if len(routes) == 0 {
		return nil, fmt.Errorf("no API keys configured to fulfill cta-ai-nexus auto-route")
	}

	return routes, nil
}
