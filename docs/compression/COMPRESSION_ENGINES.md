# Compression Engines

The **CTA AI Nexus Router** features a dual-engine compression pipeline that runs on every request before forwarding to the upstream provider. Its goal is to minimize the token footprint of the payload, directly cutting API costs.

## Pipeline Overview

Compression runs in the proxy layer (`proxy/stream.go`) on the most recent user message after it is split into semantic chunks by `engine.go`:

```
User Message
     │
     ▼
ChunkContent()          — Split content into typed segments (code, JSON, prose)
     │
     ├─ JSON chunk > 2000 chars?
     │        │
     │        ▼
     │   TOON Extraction (Phase 1)  ← Groq llama-3.1-8b-instant
     │        │ (replaces raw JSON with compact TOON notation)
     │
     └─ Any chunk (post-TOON if applicable)
              │
              ▼
     Vision Cost < Text Cost?  AND  isVisionProvider?  AND  !fidelityRequired?
              │
              ▼
     Vision Compression (Phase 2)  ← gg canvas renderer → Base64 PNG
              │
              ▼
     Rewritten Payload → Upstream Provider
```

---

## Phase 1: TOON Extraction Pipeline (`proxy/toon.go`)

**When it activates**: Any chunk whose raw content is valid JSON and exceeds **2,000 characters**.

**Why it exists**: Many AI workflows — code agents, database tools, API debuggers — embed massive JSON blobs (payloads, schemas, API responses) directly in the prompt. Forwarding 30,000 tokens of JSON to a provider is expensive and often hits context limits.

**How it works**:
1. `engine.go` detects valid JSON blocks via `json.Valid()` during content chunking.
2. The proxy pauses and fires a synchronous HTTP call to Groq's `llama-3.1-8b-instant` model (chosen for its exceptional speed and low cost).
3. A strict system prompt instructs the model to flatten the JSON into **TOON (Token-Oriented Object Notation)** — a highly compressed tabular format that strips all brackets, quotes, nesting indentation, and redundant keys while preserving semantic fidelity.
4. The raw JSON in the payload is replaced with the resulting TOON string.

**Example savings**:
| Input | Output |
|-------|--------|
| 30,000 token JSON schema | ~5,000 tokens of TOON (~83% reduction) |
| 8,000 token API response | ~1,800 tokens of TOON (~77% reduction) |

**Failure mode**: If the Groq call fails (network error, rate limit), the engine logs the failure and continues with the original JSON — the request is never dropped.

---

## Phase 2: Vision Compression Engine (`proxy/context_image.go`)

**When it activates**: For vision-capable providers (Anthropic, Gemini, Mistral, OpenAI) when the cost math favors images over text.

**Why it exists**: Vision APIs bill based on image tile count (typically 170 tokens per 512×512 tile), which can be dramatically cheaper than billing for raw text tokens for very large, densely-packed context payloads.

**Cost Decision Formula**:
```
text_cost  = len(chunk) / 4              (chars ÷ 4 ≈ tokens)
vision_cost = ceil(len(chunk) / 2000) × 170   (tiles × 170 tokens each)

Activate if: vision_cost < text_cost  AND  len(chunk) > 2000  AND  !fidelityRequired
```

**How it works** (`context_image.go`):
1. The engine allocates an in-memory 2D canvas using Go's `gg` (2D graphics) library.
2. It renders the raw text in high-contrast black on white using a dense `7×13` monospace font.
3. **Pagination**: Massive payloads are split into multiple pages (capped at ~2,000px height per page) to remain within LLM vision input limits.
4. Each page is encoded as a Base64 PNG and injected into the message as an OpenAI-compatible `image_url` content block.
5. The model field is automatically rewritten to `pixtral-12b-2409` if the route targets Mistral (the only Mistral vision model).

**Example savings**:
| Input | Output |
|-------|--------|
| 50,000 token codebase context | ~1,500 vision tokens (~97% reduction) |
| 12,000 token log file | ~510 vision tokens (~96% reduction) |

**Failure mode**: If canvas rendering fails, the engine falls back to plain text for that chunk — the request is never dropped.

---

## Telemetry

After every successful request, `db.IncrementUsage(providerID, totalSavedTokens)` logs the cumulative tokens saved to the `usage_stats` SQLite table. This data is surfaced on the **Analytics** page of the dashboard, providing a real-time view of your compression savings over time.

---

## The Chunker (`proxy/engine.go`)

Before either compression phase, `ChunkContent(text)` splits the raw user message into typed segments:

| Chunk Type | `FidelityRequired` | Notes |
|------------|-------------------|-------|
| Code block (` ``` `) | `true` | Never vision-compressed (must remain exact) |
| JSON block | `false` | TOON candidate, then vision candidate |
| Prose / Natural language | `false` | Vision candidate for large blocks |
