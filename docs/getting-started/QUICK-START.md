# Quick Start Guide

Welcome to the **CTA AI Nexus Router** — a high-performance, OpenAI-compatible proxy that dramatically cuts LLM API costs through context compression, semantic caching, and intelligent multi-provider routing.

## Prerequisites

| Dependency | Version | Purpose |
|------------|---------|---------|
| Go | 1.22+ | Backend runtime |
| Node.js | 18+ | Frontend dashboard |
| `air` (optional) | latest | Live-reloading in development |

## 1. Start the Backend

The entire backend is a single self-contained Go binary that:
- Initializes a SQLite database at startup
- Runs the Dynamic Model Discovery engine to ping all configured providers
- Serves both the API and the embedded React frontend on port `20128`

```bash
# Production mode — serves the pre-built React frontend
cd "AI Router Golang/backend"
go run .

# Development mode — auto-rebuilds on file changes
cd "AI Router Golang/backend"
air
```

The server will start on **`http://localhost:20128`** and log something like:
```
[Discovery] Initiating Dynamic Model Discovery...
[Discovery] ✅ GROQ -> Found 13 live models
[Discovery] ✅ MISTRAL -> Found 54 live models
[Discovery] Finished model discovery.
ClickToAutomate AI Nexus Router (Go API) starting on port 20128
```

## 2. Start the Frontend (Development Mode Only)

In production the React bundle is compiled into the Go binary's `public/` embed directory. For active UI development, run the Vite dev server separately:

```bash
cd "AI Router Golang/frontend"
npm install
npm run dev
```

Visit `http://localhost:5173` in your browser for hot-reloaded UI development.

To rebuild the frontend into the backend bundle:
```bash
npm run build   # outputs to ../backend/public/
```

## 3. Configure API Keys

1. Open the dashboard at `http://localhost:20128`.
2. Navigate to **Connections** in the sidebar.
3. Search or browse the 158+ listed providers, then click one to open the key manager.
4. Paste your API key and click **Add Key**.

The router securely stores all keys in the local SQLite database (`Documents/.cta-ai-nexus/router.db`). Multiple keys per provider are supported — the router automatically load-balances requests across all configured keys using `ORDER BY RANDOM()`.

## 4. Using the Auto-Router

To use the built-in Lagrangian Auto-Router, set your client's model to `cta-ai-nexus`:

- **Base URL**: `http://localhost:20128`
- **Model**: `cta-ai-nexus`
- **API Key**: `any-string` (authentication is handled by the locally stored keys)

The router will automatically analyze your prompt's complexity, select the mathematically optimal provider, compress the context, and cascade to a fallback if the first provider fails.

## 5. Direct Model Routing

You can also target any specific model from any configured provider. The router uses **Dynamic Model Discovery** to recognize which provider owns a model ID:

- **Base URL**: `http://localhost:20128`
- **Model**: `qwen/qwen3-27b`, `gemma-3-27b-it`, `mistral-large-latest`, etc.
- **API Key**: `any-string`

The router looks up the model in the live discovery map first, then falls back to prefix matching from `config/providers.json`.

## 6. Using in Cursor / Other IDEs

In any IDE that supports a custom OpenAI-compatible endpoint:

| Setting | Value |
|---------|-------|
| Base URL | `http://localhost:20128` |
| Model | `cta-ai-nexus` (auto-route) or any specific model ID |
| API Key | `any-string` |

You are now ready to enjoy mathematically optimized, cost-minimized LLM requests!
