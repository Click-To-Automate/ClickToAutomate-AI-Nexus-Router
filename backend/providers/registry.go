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
	Name     string
	BaseURL  string
	APIKey   string
	AuthType string // "bearer" | "cookie" | "oauth" | "bearer_token"
}

// GetProviderForModel determines which provider to use based on the model name
// using the dynamic discovery engine and fallback prefixes.
func GetProviderForModel(model string) (*ResolvedProvider, error) {
	discovery.Mu.RLock()
	for pid, models := range discovery.ProviderModels {
		for _, m := range models {
			if m == model {
				// Found the exact model in the dynamic discovery map!
				for _, p := range config.GlobalConfig.Providers {
					if p.ID == pid {
						discovery.Mu.RUnlock()
						key := db.GetKey(p.ID)
						if key == "" && p.EnvKey != "" {
							key = os.Getenv(p.EnvKey)
						}
						return &ResolvedProvider{
							Name:     p.ID,
							BaseURL:  p.BaseURL,
							APIKey:   key,
							AuthType: p.AuthType,
						}, nil
					}
				}
			}
		}
	}
	discovery.Mu.RUnlock()

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
					Name:     p.ID,
					BaseURL:  p.BaseURL,
					APIKey:   key,
					AuthType: p.AuthType,
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

// GetSmartRoutes builds an ordered fallback chain of available providers based on mathematical complexity
func GetSmartRoutes(complexity float64) ([]SmartRoute, error) {
	var availableProviders []string
	var activeConfigs []config.ProviderDef

	// 1. Find all providers that have keys
	for _, p := range config.GlobalConfig.Providers {
		key := db.GetKey(p.ID)
		if key == "" && p.EnvKey != "" {
			key = os.Getenv(p.EnvKey)
		}
		if key != "" {
			availableProviders = append(availableProviders, p.ID)
			activeConfigs = append(activeConfigs, p)
		}
	}

	if len(availableProviders) == 0 {
		return nil, fmt.Errorf("no API keys configured to fulfill cta-ai-nexus auto-route")
	}

	// 2. Run Lagrangian Dual Decomposition to find the mathematical optimal order
	optimalScores := CalculateOptimalRoutes(complexity, availableProviders)

	var routes []SmartRoute

	// 3. Map scores back to actual models
	for _, score := range optimalScores {
		for _, p := range activeConfigs {
			if p.ID == score.ProviderID {
				// Based on complexity, pick the target tier
				keyword := ""
				fallback := ""
				if complexity >= 0.8 {
					if p.ID == "groq" { keyword = "70b"; fallback = "llama3-70b-8192" 
					} else if p.ID == "anthropic" { keyword = "sonnet"; fallback = "claude-3-5-sonnet-20240620" 
					} else if p.ID == "gemini" { keyword = "pro"; fallback = "gemini-1.5-pro" 
					} else if p.ID == "mistral" { keyword = "large"; fallback = "mistral-large-latest" 
					} else if p.ID == "openai" { keyword = "4o"; fallback = "gpt-4o" 
					} else { keyword = "pro" } // Generic matching for 158 providers
				} else {
					if p.ID == "groq" { keyword = "8b"; fallback = "llama3-8b-8192" 
					} else if p.ID == "cerebras" { keyword = "8b"; fallback = "llama3.1-8b" 
					} else if p.ID == "gemini" { keyword = "flash"; fallback = "gemini-1.5-flash" 
					} else if p.ID == "openai" { keyword = "mini"; fallback = "gpt-4o-mini" 
					} else { keyword = "mini" } // Generic matching for 158 providers
				}

				key := db.GetKey(p.ID)
				if key == "" && p.EnvKey != "" {
					key = os.Getenv(p.EnvKey)
				}

				bestModel := discovery.GetBestModel(p.ID, keyword, fallback)
				routes = append(routes, SmartRoute{
					Provider: ResolvedProvider{
						Name:     p.ID,
						BaseURL:  p.BaseURL,
						APIKey:   key,
						AuthType: p.AuthType,
					},
					ActualModel: bestModel,
				})
				break
			}
		}
	}

	return routes, nil
}
