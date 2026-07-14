# CTA AI Nexus Router — Documentation

Welcome to the documentation for the **CTA AI Nexus Router** — a high-performance, OpenAI-compatible proxy that slashes LLM API costs through intelligent routing, context compression, and semantic caching.

## Release Notes

### v1.0.0
- **Dual Build System**: Now ships with both a lightweight **CLI executable** and a full **GUI desktop application** (Wails).
- **Advanced Logging**: Real-time memory log buffer in the UI, plus automatic background physical `logs.txt` writer inside the database folder.
- **Enhanced Playground**: Smooth auto-scroll behavior for long streams and quick "Delete Chat" session management.
- **CORS Compatibility**: Proper preflight `OPTIONS` support on proxy endpoints, enabling in-browser web UI fetching directly.
- Fixed **image processing** to support vision-capable models (e.g., `gpt-4o`, `claude-3.5-sonnet`).

## Documentation Index

### Getting Started
- [Quick Start Guide](getting-started/QUICK-START.md) — Install, run, configure, and start routing requests in minutes
- [CLI Guide](getting-started/CLI-GUIDE.md) — How to run the lightweight headless router via command line

### Architecture
- [System Architecture](architecture/ARCHITECTURE.md) — Directory structure, request flow, API endpoints, and data persistence
- [Core Technologies](architecture/CORE_TECHNOLOGIES.md) — Deep dive into the compression pipeline, Lagrangian solver, discovery engine, and database schema

### Routing
- [Smart Routing](routing/SMART-ROUTING.md) — How complexity scoring, Lagrangian optimization, and the cascade fallback work

### Providers
- [Providers Guide](providers/PROVIDERS-GUIDE.md) — How provider resolution works, multi-key load balancing, payload normalization, and adding new providers

### Compression
- [Compression Engines](compression/COMPRESSION_ENGINES.md) — TOON extraction, Vision compression, the chunker, and telemetry

---

## Quick Overview

```
Client (Cursor / IDE / API call)
    │
    │  POST http://localhost:20128/v1/chat/completions
    │  { "model": "cta-ai-nexus", "messages": [...] }
    ▼
CTA AI Nexus Router (Go, port 20128)
    │
    ├─ Intent Analysis → Complexity Score (0.0 – 1.0)
    ├─ Lagrangian Solver → Optimal Provider Cascade
    ├─ TOON Extraction → JSON compression via Groq
    ├─ Vision Compression → Image rendering via gg canvas
    │
    ▼
Upstream Provider (Groq / Mistral / Anthropic / Gemini / ...)
    │
    ▼
Streamed Response → Client
```

## Port Reference

| Service | URL |
|---------|-----|
| Main Proxy + Dashboard | `http://localhost:20128` |
| Frontend Dev Server | `http://localhost:5173` (dev only) |
