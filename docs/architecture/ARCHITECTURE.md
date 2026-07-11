# System Architecture

The **CTA AI Router** is architected for extreme performance, low latency, and zero dependencies on heavy external infrastructure (like PostgreSQL or Redis). It runs entirely locally via Go and SQLite.

## Core Flow (`stream.go`)

1. **Request Interception**: An OpenAI-compatible request arrives at `http://localhost:20128/chat/completions`.
2. **Intent Analysis**: The proxy unmarshals the JSON body early to pass the messages to `AnalyzeTask()` (located in `intent.go`).
3. **Chunking & Pipeline Processing**: The user's prompt is split into chunks by the markdown parser (`engine.go`). 
    - The **TOON Extractor** runs synchronously for massive JSON blocks.
    - The **Context Compression Engine** calculates exact visual rendering costs versus text token costs.
4. **Smart Routing**: Based on the analyzed intent (`simple`, `code`, `complex`, `massive`), `registry.go` constructs a deterministic fallback chain (e.g., `Groq -> Anthropic -> OpenAI -> Mistral -> Gemini`).
5. **Streaming Proxy**: The backend rewrites the JSON payload dynamically (stripping unauthorized fields or substituting image bases), forwards it to the target provider, and streams the chunked response back to the client using `http.Flusher`.

## Internal State Management (`db/client.go`)

To maximize portability, the router stores its state locally using `modernc.org/sqlite` (a pure Go port of SQLite, removing CGO dependencies).

The database is located in your system's Documents folder (`.cta-ai-nexus/router.db`) and stores two primary tables:
1. `provider_keys`: Secure, local storage of API keys.
2. `usage_stats`: Incrementing request counts for the dashboard telemetry.

## Performance Profile
- Written in Go to handle hundreds of concurrent streams without blocking.
- In-memory canvas rendering takes ~5-15ms for a 40,000 character context payload.
- Zero-allocation string building for Intent Analysis.
