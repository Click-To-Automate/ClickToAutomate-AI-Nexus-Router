# Core Technologies

The **CTA AI Nexus Router** is built on three core technical pillars:

1. **RTK + Caveman Stacked Compression** — minimizes token costs before sending to providers
2. **Lagrangian Dual Decomposition** — picks the mathematically optimal provider for every request
3. **Dynamic Model Discovery** — keeps the router current with live provider model catalogs

---

## Part 1: RTK + Caveman Stacked Compression (Token Reducer)

When running coding assistants or agentic pipelines, prompts contain two very different types of content: machine-generated outputs (logs, traces, JSON payloads) and human-written conversational context.

Instead of a single generic compressor, the AI Nexus Router feeds prompts through a **two-phase pipeline** running sequentially in `proxy/stream.go`:

```
Raw Prompt
    │
    ▼
[ Phase 1: TOON Extraction ]  ─── Groq llama-3.1-8b-instant
    │         Strips JSON → compact tabular TOON notation
    ▼
[ Phase 2: Vision Compression ]  ─── gg canvas library → Base64 PNG
              Renders large prose/context to image tiles
    │
    ▼
Compressed Payload → Upstream Provider
```

### Phase 1: TOON Extraction (Terminal/Structured Output Filter)

- **Target**: Machine-generated content — large JSON blobs, API responses, database payloads
- **Mechanism**: Detected via `json.Valid()` during content chunking. Chunks over 2,000 chars are sent to `llama-3.1-8b-instant` with a system prompt that aggressively flattens the structure into **Token-Oriented Object Notation (TOON)**
- **Result**: 77–83% reduction on structured data

### Phase 2: Vision Compression Engine (Natural Language / Context Reducer)

- **Target**: Large prose blocks and remaining context that exceeds the token cost threshold for vision
- **Mechanism**: Uses Go's `gg` 2D graphics library to render text as high-contrast PNG images in memory. Forwarded as Base64 `image_url` content blocks to vision-capable providers (Anthropic, Gemini, OpenAI, Mistral)
- **Result**: 93–97% reduction on large context payloads

### Combined Savings

By chaining both phases, the cumulative savings on a mixed prompt (JSON + prose) can exceed 95%:

$$\text{Total Savings} \approx 1 - (1 - 0.80_{\text{TOON}}) \times (1 - 0.95_{\text{Vision}}) \approx 99\%$$

> **Note**: Code blocks (```` ``` ````) are always preserved verbatim (`FidelityRequired = true`) and never subjected to lossy compression.

---

## Part 2: Lagrangian Dual Decomposition (Intelligent Router)

Once the payload is compressed, the router must choose the best provider. It models this as a **constrained optimization problem** rather than using brittle rule-based forwarding.

### The Problem

Every request has three competing objectives:
1. **Minimize Cost + Latency**: Use the cheapest/fastest provider available
2. **Satisfy Quality Constraint**: The chosen provider must be capable enough for the prompt's complexity
3. **Avoid Penalized Providers**: Providers that recently failed (rate limits, errors) should be deprioritized

### The Solution: Lagrangian Objective Function

For each available provider, the solver computes:

$$L(x, \lambda) = \text{Cost}(x) + \text{Latency}(x) + \lambda \cdot \max(0, q_{min} - \text{Quality}(x))$$

- **$q_{min}$**: The minimum quality threshold, derived from the prompt's complexity score (0.0–1.0)
- **$\lambda$** (Lagrange multiplier): A penalty value that spikes when a provider returns errors and decays over time, automatically routing traffic away from congested or failing endpoints
- Providers are sorted by ascending $L$ score, producing the optimal cascade order

### Self-Healing via Multiplier Updates

When a provider returns `429`, `413`, or `5xx`, `IncrementPenalty(providerID)` spikes the $\lambda$ for that provider. The effect is immediate: the next request will mathematically deprioritize it, creating a self-healing load distribution without any manual intervention.

---

## Part 3: Dynamic Model Discovery

Providers constantly add and deprecate models. Hardcoded model strings break silently.

At startup, `discovery.RunDiscovery()` pings `GET /v1/models` for every provider with a configured API key. Live models are stored in `discovery.ProviderModels` (a `sync.RWMutex`-guarded concurrent map).

When the router needs to select a model for a provider (e.g., "give me an 8B model on Groq"), it calls:
```go
discovery.GetBestModel("groq", "8b", "llama3-8b-8192")
```

This returns the best live match from the discovered catalog, falling back to the hardcoded default only if discovery failed for that provider.

---

## Part 4: Frontend — React + Vite Dashboard

The frontend is a React + TypeScript application built with Vite, embedded directly into the Go binary at compile time.

Key UI features:
- **Connections** page: Browse 158+ providers with dynamic filters (Integration Type, Pricing, Account requirement), search, and multi-key management
- **Playground**: Full ChatGPT-like interface with model dropdown, chat history (localStorage-backed), Memory toggle + sidebar, and `<think>` tag rendering (collapsible reasoning accordion with Markdown support inside)
- **Analytics**: Live telemetry charts for per-provider request counts and cumulative tokens saved
- **Semantic Cache**: View and purge all cached prompt/response pairs
- **Compression**: Dashboard for viewing compression activity and saved tokens
- **MCP**: Model Context Protocol tool bridge configuration

---

## Part 5: Database Schema (SQLite — `db/client.go`)

The router uses `modernc.org/sqlite` (pure Go, no CGO) for all persistent state. The database lives at `Documents/.cta-ai-nexus/router.db`.

| Table | Schema | Purpose |
|-------|--------|---------|
| `provider_keys_multi` | `(id, provider_id, api_key, UNIQUE(provider_id, api_key))` | Multi-key storage with automatic load balancing via `ORDER BY RANDOM()` |
| `usage_stats` | `(provider_id, request_count, tokens_saved)` | Telemetry data per provider |
| `app_settings` | `(key, value)` | Key-value config store |
| `semantic_cache` | `(hash, prompt, response, tokens_saved, created_at)` | SHA-256 prompt → response cache |
