# System Architecture

The **CTA AI Nexus Router** is a high-performance, self-contained OpenAI-compatible proxy written in Go. It is designed for extreme throughput, low latency, and zero dependency on heavy external infrastructure (no PostgreSQL, no Redis). It runs entirely locally via a single Go binary backed by SQLite.

## Directory Structure

```
backend/
├── api/
│   ├── server.go              # HTTP mux + route registration
│   └── handlers/              # Individual REST API handlers
│       ├── chat.go            # /v1/chat/completions
│       ├── keys.go            # /v1/keys (CRUD)
│       ├── providers.go       # /v1/providers (list)
│       ├── usage.go           # /v1/usage (telemetry)
│       ├── cache.go           # /v1/cache (semantic cache)
│       ├── settings.go        # /v1/settings (app config)
│       └── mcp.go             # /v1/mcp (MCP tool bridge)
├── proxy/
│   ├── stream.go              # Core proxy engine & cascade logic
│   ├── intent.go              # AnalyzeComplexity() scorer
│   ├── engine.go              # Content chunker (ChunkContent)
│   ├── context_image.go       # Vision compression engine
│   └── toon.go                # TOON extraction pipeline
├── providers/
│   ├── registry.go            # GetProviderForModel + GetSmartRoutes
│   └── lagrangian.go          # Lagrangian Dual Decomposition solver
├── discovery/
│   └── discovery.go           # Dynamic model discovery engine
├── config/
│   └── loader.go              # Reads providers.json into GlobalConfig
├── db/
│   └── client.go              # SQLite DB init, CRUD helpers
└── main.go                    # Entry point
```

## Core Request Flow (`stream.go`)

1. **Request Interception** — An OpenAI-compatible `POST /v1/chat/completions` arrives at `http://localhost:20128`.
2. **Model Resolution** — The `model` field in the JSON body determines the routing path:
   - `"cta-ai-nexus"` → triggers the **Lagrangian Auto-Router**
   - Any other model string → **Direct Routing** via `GetProviderForModel()`
3. **Intent Analysis** (Auto-Router path only) — `AnalyzeComplexity(messages)` scores the prompt between `0.0` and `1.0`.
4. **Smart Route Construction** — `GetSmartRoutes(complexity)` runs the Lagrangian solver, building a sorted fallback chain of available providers.
5. **Content Compression Pipeline** — For each route attempt, the proxy runs the hybrid compression pipeline (`ChunkContent → TOON → Vision`) to minimize the token footprint before forwarding.
6. **Payload Sanitization** — Strict providers (Mistral, Cerebras) reject proprietary client fields. The proxy deletes the `"user"` key and any offending fields before forwarding.
7. **Upstream Forwarding** — The rewritten payload is sent to the target provider over a standard `net/http` client with `Authorization: Bearer <key>`.
8. **Error Cascade** — On `429`, `413`, `5xx`, or model-dead errors, `IncrementPenalty()` spikes the Lagrange multiplier for that provider and the next route in the chain is tried.
9. **Streaming Proxy** — On success, the response is streamed back to the client in 4KB chunks via `http.Flusher`, with telemetry logged to SQLite.

## Internal State Management (`db/client.go`)

The router stores all persistent state in a local SQLite database using `modernc.org/sqlite` (a pure-Go port — no CGO required).

**Database Location**: `C:\Users\<username>\Documents\.cta-ai-nexus\router.db`

| Table | Purpose |
|-------|---------|
| `provider_keys_multi` | Multi-key storage per provider. `GetKey()` uses `ORDER BY RANDOM()` to load-balance across configured keys. |
| `usage_stats` | Per-provider request counts and cumulative tokens saved, updated on every successful request. |
| `app_settings` | Key-value store for application configuration (e.g., compression thresholds). |
| `semantic_cache` | SHA-256-keyed response cache for deduplicating identical prompts. |

> Note: The legacy `provider_keys` table is auto-migrated to `provider_keys_multi` on startup to support multiple keys per provider with automatic load balancing.

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/chat/completions` | `POST` | Main proxy endpoint (OpenAI-compatible) |
| `/v1/models` | `GET` | Returns all live discovered models |
| `/v1/keys` | `GET/POST` | CRUD for provider API keys |
| `/v1/providers` | `GET` | Returns provider list from `providers.json` |
| `/v1/usage` | `GET` | Returns per-provider telemetry |
| `/v1/cache` | `GET/DELETE` | Manages the semantic response cache |
| `/v1/settings` | `GET/POST` | App-level configuration |
| `/v1/mcp` | `POST` | MCP (Model Context Protocol) tool bridge |

## Performance Profile

- Written in Go to handle hundreds of concurrent streams without blocking.
- Stateless per-request proxy — no global mutable state in the hot path.
- In-memory canvas rendering takes ~5–15ms for a 40,000-character context payload.
- SQLite with `ORDER BY RANDOM()` provides zero-config multi-key load balancing.
