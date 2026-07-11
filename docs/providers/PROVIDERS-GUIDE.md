# Providers Guide

The **CTA AI Router** is completely LLM-agnostic, but applies aggressive normalizations to ensure payloads succeed across strict providers.

## Supported Providers

| Provider | Supported | Notes |
|----------|-----------|-------|
| **Groq** | ✅ Yes | Prioritized for Code & Simple tasks. Uses LLaMA 3.3. Powers the internal TOON extractor. |
| **Mistral** | ✅ Yes | Supports the Vision Context Engine by dynamically rewriting the target model to `pixtral-12b-2409` when images are detected. |
| **Anthropic**| ✅ Yes | Interfaced via standard OpenAI compatibility layers. |
| **Gemini** | ✅ Yes | Used as a reliable massive-context fallback. |
| **Cerebras** | ✅ Yes | Insanely fast inference, restricted to simple/small-context queries. |
| **OpenAI** | ✅ Yes | Used strictly as a premium tier fallback for complex tasks. |

## Payload Normalization

Many tools (like Cursor or Antigravity) append proprietary fields (like `"user": "cursor"`) to the JSON body. 

Strict providers (like Mistral or Cerebras) immediately reject the payload with a `400 Bad Request` if they encounter undocumented fields. 

To resolve this, the **CTA AI Router** aggressively sanitizes the payload directly inside `stream.go` before routing:
```go
// Sanitize payload
delete(payload, "user")
```

This guarantees flawless compatibility across all supported inference endpoints without requiring the user to modify their client software.
