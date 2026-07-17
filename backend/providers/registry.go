package providers

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"ainexusrouter-core/config"
	"ainexusrouter-core/db"
	"ainexusrouter-core/discovery"
)

type ResolvedProvider struct {
	Name              string
	BaseURL           string
	RequiresCustomURL bool
	APIKey            string
	AuthType          string // "bearer" | "cookie" | "oauth" | "bearer_token"
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
							Name:              p.ID,
							BaseURL:           p.BaseURL,
							RequiresCustomURL: p.RequiresCustomURL,
							APIKey:            key,
							AuthType:          p.AuthType,
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
					Name:              p.ID,
					BaseURL:           p.BaseURL,
					RequiresCustomURL: p.RequiresCustomURL,
					APIKey:            key,
					AuthType:          p.AuthType,
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("no provider configured for model: %s", model)
}

// GetProviderByID returns the provider configuration explicitly by ID
func GetProviderByID(providerID string) (*ResolvedProvider, error) {
	for _, p := range config.GlobalConfig.Providers {
		if p.ID == providerID {
			key := db.GetKey(p.ID)
			if key == "" && p.EnvKey != "" {
				key = os.Getenv(p.EnvKey)
			}
			return &ResolvedProvider{
				Name:              p.ID,
				BaseURL:           p.BaseURL,
				RequiresCustomURL: p.RequiresCustomURL,
				APIKey:            key,
				AuthType:          p.AuthType,
			}, nil
		}
	}
	return nil, fmt.Errorf("provider %s not found", providerID)
}

// SmartRoute represents a resolved provider and the actual model payload we should rewrite to
type SmartRoute struct {
	Provider    ResolvedProvider
	ActualModel string
}

type IntentProfile struct {
	Complexity float64
	IsCoding   bool
	IsReasoning bool
	HasImage   bool
}

// GetSmartRoutes builds an ordered fallback chain of available providers based on mathematical complexity and intent
func GetSmartRoutes(profile IntentProfile) ([]SmartRoute, error) {
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

	var routes []SmartRoute

	// 2. Check if user has defined a custom preferred models priority
	preferredModelsJSON := db.GetSetting("preferred_models", "[]")
	var preferredProviders []string
	if err := json.Unmarshal([]byte(preferredModelsJSON), &preferredProviders); err == nil && len(preferredProviders) > 0 {
		// Strict priority routing
		for _, prefItem := range preferredProviders {
			parts := strings.SplitN(prefItem, "::", 2)
			if len(parts) != 2 {
				// Fallback logic for old configs that only had provider IDs
				prefID := parts[0]
				for _, p := range activeConfigs {
					if p.ID == prefID {
						bestModel := discovery.GetOptimalModel(p.ID, profile.Complexity, profile.IsCoding, profile.IsReasoning, profile.HasImage)
						key := db.GetKey(p.ID)
						if key == "" && p.EnvKey != "" {
							key = os.Getenv(p.EnvKey)
						}
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
				continue
			}

			prefID := parts[0]
			prefModel := parts[1]

			// Find the provider config
			for _, p := range activeConfigs {
				if p.ID == prefID {
					key := db.GetKey(p.ID)
					if key == "" && p.EnvKey != "" {
						key = os.Getenv(p.EnvKey)
					}
					routes = append(routes, SmartRoute{
						Provider: ResolvedProvider{
							Name:     p.ID,
							BaseURL:  p.BaseURL,
							APIKey:   key,
							AuthType: p.AuthType,
						},
						ActualModel: prefModel, // Use the explicit model specified by the user!
					})
					break
				}
			}
		}
		
		if len(routes) > 0 {
			return routes, nil
		}
	}

	// 3. Fallback: Run Lagrangian Dual Decomposition to find the mathematical optimal order
	optimalScores := CalculateOptimalRoutes(profile.Complexity, availableProviders)

	// 4. Map scores back to actual models
	for _, score := range optimalScores {
		for _, p := range activeConfigs {
			if p.ID == score.ProviderID {
				bestModel := discovery.GetOptimalModel(p.ID, profile.Complexity, profile.IsCoding, profile.IsReasoning, profile.HasImage)

				key := db.GetKey(p.ID)
				if key == "" && p.EnvKey != "" {
					key = os.Getenv(p.EnvKey)
				}

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
