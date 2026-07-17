package providers

import (
	"sort"
	"sync"
	"ainexusrouter-core/db"
)

// ModelStats represents intrinsic constraints of a model for the primal problem
type ModelStats struct {
	Cost    float64 // Tiered: 0.0 (Free), 10.0 (Low Cost), 20.0 (Paid)
	Latency float64 // normalized 0.0-1.0
	Quality float64 // normalized 0.0-1.0
}

// Baseline knowledge of our model families
var providerStats = map[string]ModelStats{
	// --- FREE TIER (Cost 0.0) ---
	"groq":      {Cost: 0.0, Latency: 0.2, Quality: 0.85},
	"cerebras":  {Cost: 0.0, Latency: 0.1, Quality: 0.85},
	"nvidia":    {Cost: 0.0, Latency: 0.3, Quality: 0.90}, // NIM free credits
	"xai":       {Cost: 0.0, Latency: 0.4, Quality: 0.85},
	"together":  {Cost: 0.0, Latency: 0.3, Quality: 0.85},

	// --- LOW COST TIER (Cost 10.0) ---
	"deepseek":  {Cost: 10.0, Latency: 0.5, Quality: 0.90},
	"mistral":   {Cost: 10.0, Latency: 0.5, Quality: 0.85},
	"gemini":    {Cost: 10.0, Latency: 0.6, Quality: 0.95},

	// --- PAID TIER (Cost 20.0) ---
	"openai":    {Cost: 20.0, Latency: 0.3, Quality: 0.98},
	"anthropic": {Cost: 20.0, Latency: 0.4, Quality: 0.98},
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
	
	realLatencies := db.GetLatencies()

	for _, pid := range availableProviders {
		stats, ok := providerStats[pid]
		if !ok {
			// For unknown providers, assume they are Low Cost by default so they don't unexpectedly drain Paid credits, but don't blindly trust them as Free
			stats = ModelStats{Cost: 10.0, Latency: 0.5, Quality: 0.85}
		}

		// Quality Constraint Check: If quality is significantly lower than complexity, heavily penalize.
		qualityPenalty := 0.0
		if stats.Quality < complexity {
			qualityPenalty = (complexity - stats.Quality) * 5.0 // Slack variable penalty
		}

		// Lambda multiplier for dynamic rate-limits
		lambda := multipliers.m[pid]
		
		// Use dynamic latency if available (normalized roughly by max expected 5000ms)
		actualLatency := stats.Latency
		if realLatM, ok := realLatencies[pid]; ok && realLatM > 0 {
			normalized := float64(realLatM) / 5000.0
			if normalized > 1.0 {
				normalized = 1.0
			}
			actualLatency = normalized
		}

		// Objective Function: L(x, lambda)
		// We want to minimize cost and latency, but penalize missing the quality constraint and past failures.
		objective := stats.Cost + actualLatency + qualityPenalty + lambda

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
