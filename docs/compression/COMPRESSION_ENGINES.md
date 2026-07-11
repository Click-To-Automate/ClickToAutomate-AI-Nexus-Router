# Compression Engines

The crown jewel of the **CTA AI Router** is its dual-engine compression pipeline. These engines act as a financial optimizer, actively rewriting your prompts to drastically cut API costs.

## 1. Context Compression Engine (Text-to-Image)

Providers like Mistral (Pixtral), OpenAI, and Anthropic charge heavily for text tokens, but offer cheap, tile-based billing for Vision requests.

**How it works (`context_image.go`)**:
1. When a massive payload is detected, the engine predicts the math: text tokens (`chars / 4`) vs vision tokens (`170 per 512x512 tile`).
2. If vision is cheaper, Go initializes an in-memory 2D canvas using the `gg` graphics library.
3. It renders the raw codebase context in high-contrast black text on a white background using a dense `7x13` monospace font.
4. **Pagination**: To ensure readability by the LLM, the renderer safely splits massive payloads into multiple smaller Base64 PNG "pages" (capped at ~2000px height).
5. The router rewrites the JSON payload, replacing your raw text with an OpenAI-compatible Vision array `[{"type": "image_url"...}]`.

**Result**: A 50,000 token payload drops to ~1,500 vision tokens.

## 2. TOON Extraction Pipeline (LLM-based)

If you are forced to use a text-only model (like Groq's LLaMA 70B), image compression is impossible. For massive structured data (like JSON), we use the TOON Pipeline.

**How it works (`toon.go`)**:
1. `engine.go` uses Go's strict `json.Valid()` to detect massive JSON blobs inside the user's prompt.
2. If detected, the proxy pauses the request and fires a synchronous HTTP call to Groq's insanely fast `llama-3.1-8b-instant` model.
3. A strict system prompt forces the 8B model to aggressively flatten the JSON into **Token-Oriented Object Notation (TOON)**—a highly condensed tabular format that strips all brackets, quotes, and redundant keys.
4. The massive JSON is seamlessly replaced with the tiny TOON string in your prompt.

**Result**: A 30,000 token JSON file is reduced to ~5,000 tokens of TOON data, natively preserving perfect fidelity for the target model while circumventing tight rate limits.
