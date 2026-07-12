# Smart Routing: Lagrangian Dual Decomposition

The **CTA AI Nexus Router** removes the burden of manually selecting LLM providers. When you set `model: "cta-ai-nexus"`, the router uses **Lagrangian Dual Decomposition** to build a mathematically optimal, ordered fallback chain for every request — dynamically balancing cost, latency, and quality against your prompt's complexity.

## 1. Complexity Scoring (`proxy/intent.go`)

When a request arrives, `AnalyzeComplexity(messages)` scans the concatenated prompt text and computes a rigorous complexity score between `0.0` and `1.0`.

| Signal | Score Adjustment |
|--------|-----------------|
| Baseline (every request) | `+0.5` |
| Code keywords: `function`, `class`, `def`, `error:`, `exception`, `=>` | `+0.05` each |
| Massive payload (>10,000 tokens) | Forces `≥ 0.95` |
| Math/reasoning keywords: `prove`, `theorem`, `integral` | `+0.05` each |

This score becomes the minimum **Quality Constraint** ($q_{min}$) fed into the Lagrangian solver.

## 2. Lagrangian Dual Decomposition Solver (`providers/lagrangian.go`)

The router models provider selection as a constrained optimization problem. For each available provider:

$$L(x, \lambda) = \text{Cost}(x) + \text{Latency}(x) + \lambda \cdot \max(0, q_{min} - \text{Quality}(x))$$

Where:
- **Cost** and **Latency** are hardcoded baseline scores per provider (e.g., Cerebras = very low, GPT-4o = very high)
- **Quality** is a baseline capability score per provider (e.g., Cerebras = `0.7`, GPT-4o = `1.0`)
- **$\lambda$** is the current Lagrange multiplier (penalty) for that provider — spiked on errors, decayed over time

The solver evaluates all available providers and sorts them by ascending Lagrangian score, producing the optimal fallback chain.

## 3. Route Construction (`providers/registry.go`)

`GetSmartRoutes(complexity)` performs three steps:

**Step 1** — Enumerate all providers with configured API keys (from `provider_keys_multi` table or env vars).

**Step 2** — Run `CalculateOptimalRoutes(complexity, availableProviders)` to get the Lagrangian-sorted order.

**Step 3** — Map each provider to the best available model using `discovery.GetBestModel()`:

| Complexity | Tier Selected | Example |
|------------|---------------|---------|
| `>= 0.8` (Complex) | High-tier models | Groq `70b`, Anthropic `sonnet`, Gemini `pro`, Mistral `large`, OpenAI `4o` |
| `< 0.8` (Simple) | Efficient models | Groq `8b`, Cerebras `8b`, Gemini `flash`, OpenAI `mini` |

The `GetBestModel(providerID, keyword, fallback)` function searches the live discovery map for a model containing `keyword`, falling back to a hardcoded string if the provider wasn't successfully discovered at startup.

## 4. Dynamic Model Discovery (`discovery/discovery.go`)

Providers constantly deprecate models and add new ones. The router never relies on hardcoded model strings when it doesn't have to.

At startup, `RunDiscovery()` concurrently pings `GET /v1/models` for every configured provider:
- Results are parsed from OpenAI format (`{"data": [{"id": "..."}]}`) or Google format (`{"models": [{"name": "models/..."}]}`)
- Non-chat models (embeddings, TTS, Whisper, OCR) are filtered out via `isChatModel()`
- Live models are stored in `discovery.ProviderModels` (a concurrent-safe `sync.RWMutex`-guarded map)
- When routing, `GetBestModel()` queries this live map rather than using brittle hardcoded strings

## 5. Self-Healing Cascade (`proxy/stream.go`)

The router is fully fault-tolerant:

1. The Lagrangian solver produces an ordered route chain (e.g., `Groq → Mistral → Anthropic → Gemini`).
2. The proxy tries each route sequentially.
3. On any cascade-triggering error (see Providers Guide), it calls `IncrementPenalty(providerID)`.
4. **Penalty Effect**: The spiked $\lambda$ multiplier for that provider means the next request will mathematically deprioritize it — routing traffic away from the congested provider until penalties decay.
5. The cascade continues until a route succeeds or all routes are exhausted (returns `502 Bad Gateway`).

```
Request → Groq [429 → Penalty↑] → Mistral [Success ✅] → Response
                                    ↓ telemetry logged to SQLite
```

## 6. Semantic Cache (Compression Layer)

Before hitting any provider, the `chat.go` handler checks the `semantic_cache` SQLite table for a matching SHA-256 hash of the prompt. On a cache hit, the response is returned instantly with zero upstream API cost.

Cache entries are viewable and purgeable from the **Semantic Cache** page in the dashboard.
