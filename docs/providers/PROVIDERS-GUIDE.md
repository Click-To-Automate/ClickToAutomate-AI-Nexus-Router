# Providers Guide

The **CTA AI Nexus Router** is fully LLM-agnostic and supports 158+ providers. It uses a two-tier discovery system to map any model string to the correct upstream provider, and aggressively normalizes payloads to ensure compatibility with strict API endpoints.

## Provider Resolution: How It Works

When a request arrives with a model field like `qwen/qwen3-27b` or `gemma-3-27b-it`, the router resolves the correct provider in two stages (`registry.go`):

### Stage 1: Dynamic Discovery Lookup (Primary)

On startup, `discovery.go` pings the `GET /v1/models` endpoint of every provider that has a configured API key. The live model list is cached in `discovery.ProviderModels` (a `map[string][]string` guarded by a `sync.RWMutex`).

When a request arrives, `GetProviderForModel()` first scans this map for an **exact match**. This allows the router to correctly handle any model discovered at runtime — even brand-new models that were added to a provider after the router was compiled.

### Stage 2: Prefix Matching (Fallback)

If no exact match is found in the discovery map (e.g., for providers without discoverable model lists), the router falls back to **prefix matching** against the `prefixes` array defined for each provider in the embedded `backend/config/providers.json`:

```json
{ "id": "groq", "prefixes": ["groq/", "llama", "gemma", "mixtral", "whisper"] }
```

## Supported Provider Categories

| Category | Examples | Auth Method |
|----------|---------|-------------|
| **API Key** | OpenAI, Anthropic, Groq, Mistral, Cerebras, DeepSeek, SambaNova | Dashboard → Connections |
| **Free Tier** | Groq, Gemini, Mistral, Cerebras, SambaNova, LLM7, GitHub Models | Dashboard → Connections |
| **OAuth** | Claude Code, GitHub Copilot, Cursor IDE, Windsurf, Grok Build | OAuth flow (future) |
| **Web Cookie** | ChatGPT Web, Claude Web, DeepSeek Web, Gemini Web | Cookie paste (future) |
| **Local** | Ollama, LM Studio, vLLM | Local endpoint URL |
| **Aggregators** | OpenRouter, AgentRouter, AI/ML API | API Key |

## Multi-Key Load Balancing

Multiple API keys can be configured per provider. The `db.GetKey(providerID)` function always uses `ORDER BY RANDOM() LIMIT 1` to fetch a key, distributing load evenly across all configured keys automatically — no additional configuration required.

## Payload Normalization

Many AI clients (Cursor, Antigravity, Cline, etc.) append proprietary fields to the JSON body. Strict providers immediately reject payloads containing undocumented fields with `400 Bad Request`.

The proxy sanitizes the payload **before** forwarding (`stream.go`):
```go
delete(payload, "user")  // Strips client-injected user field
```

Additionally, if the Vision Compression Engine activates and generates image content for a Mistral route, the proxy automatically rewrites the model to `pixtral-12b-2409` (the only Mistral model that accepts vision inputs):
```go
if hasImage && route.Provider.Name == "mistral" {
    payload["model"] = "pixtral-12b-2409"
}
```

## Adding a New Provider

1. Add an entry to `backend/config/providers.json`:
```json
{
  "id": "my-provider",
  "name": "My Provider",
  "base_url": "https://api.myprovider.com/v1",
  "env_key": "MY_PROVIDER_API_KEY",
  "prefixes": ["my-provider/", "mymodel-"]
}
```
2. Configure an API key via the dashboard (Connections → search → click provider → Add Key), or set the `env_key` environment variable.
3. Restart the backend — the Discovery Engine will automatically enumerate its live models on next startup.

## Error Cascade & Fallback Behavior

The router treats the following upstream responses as cascade triggers (tries next provider):
- `429 Too Many Requests`
- `413 Request Entity Too Large`
- `5xx` Server Errors
- `400` with body containing: `decommissioned`, `deprecated`, `not found`, `does not exist`, `rate_limit_exceeded`, `too large`, `limit exceeded`, `maximum context length`

On each cascade, `IncrementPenalty(providerID)` is called to spike the Lagrange multiplier for that provider, deprioritizing it for subsequent requests until the penalty naturally decays.
