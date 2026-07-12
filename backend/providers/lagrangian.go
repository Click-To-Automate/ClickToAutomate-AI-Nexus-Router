package providers

import (
	"sort"
	"sync"
)

// ModelStats represents intrinsic constraints of a model for the primal problem
type ModelStats struct {
	Cost    float64 // normalized 0.0-1.0
	Latency float64 // normalized 0.0-1.0
	Quality float64 // normalized 0.0-1.0
}

// Baseline knowledge of our model families
var providerStats = map[string]ModelStats{
	"cerebras":  {Cost: 0.1, Latency: 0.1, Quality: 0.60},
	"groq":      {Cost: 0.3, Latency: 0.2, Quality: 0.85},
	"openai":    {Cost: 0.2, Latency: 0.2, Quality: 0.75}, // assuming 4o-mini
	"gemini":    {Cost: 0.5, Latency: 0.6, Quality: 0.88},
	"mistral":   {Cost: 0.7, Latency: 0.5, Quality: 0.90},
	"anthropic": {Cost: 0.6, Latency: 0.5, Quality: 0.95},
}

// multipliers tracks the Lagrange multipliers (λ) for each provider
// If a provider fails or rate limits, its multiplier goes up, aggressively penalizing it.
var multipliers = struct {
	sync.RWMutex
	m map[string]float64
}{m: make(map[string]float64)}

// IncrementPenalty increases the Lagrange multiplier for a provider (e.g., on 429 error)
func IncrementPenalty(providerID string) {
	multipliers.Lock()
	defer multipliers.Unlock()
	current := multipliers.m[providerID]
	multipliers.m[providerID] = current + 1.0 // Add heavy penalty
}

// DecayPenalties simulates adaptive convergence by gradually reducing multipliers over time
// This could be run in a background ticker.
func DecayPenalties() {
	multipliers.Lock()
	defer multipliers.Unlock()
	for k, v := range multipliers.m {
		if v > 0 {
			multipliers.m[k] = v * 0.9 // Decay by 10%
		}
	}
}

// RouteScore represents an evaluated model path
type RouteScore struct {
	ProviderID string
	Objective  float64
}

// CalculateOptimalRoutes performs the Dual Decomposition
// Goal: Minimize Cost + Latency + (Lambda * Penalty) subject to Quality > Complexity
// A lower Objective is better.
func CalculateOptimalRoutes(complexity float64, availableProviders []string) []RouteScore {
	var scores []RouteScore

	multipliers.RLock()
	defer multipliers.RUnlock()

	for _, pid := range availableProviders {
		stats, ok := providerStats[pid]
		if !ok {
			// For the 158 unknown generic API key providers, assume standard flagship quality
			stats = ModelStats{Cost: 0.5, Latency: 0.5, Quality: 0.85}
		}

		// Quality Constraint Check: If quality is significantly lower than complexity, heavily penalize.
		qualityPenalty := 0.0
		if stats.Quality < complexity {
			qualityPenalty = (complexity - stats.Quality) * 5.0 // Slack variable penalty
		}

		// Lambda multiplier for dynamic rate-limits
		lambda := multipliers.m[pid]

		// Objective Function: L(x, lambda)
		// We want to minimize cost and latency, but penalize missing the quality constraint and past failures.
		objective := stats.Cost + stats.Latency + qualityPenalty + lambda

		scores = append(scores, RouteScore{
			ProviderID: pid,
			Objective:  objective,
		})
	}

	// Sort mathematically: Lowest objective score first
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Objective < scores[j].Objective
	})

	return scores
}
