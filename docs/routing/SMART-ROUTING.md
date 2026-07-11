# Smart Routing & Fallbacks

The **CTA AI Router** removes the burden of manually selecting LLM models from the developer. It intelligently dictates the routing path based on the structure of the prompt.

## The Intent Analyzer

When a request arrives, `AnalyzeTask()` scans the concatenated prompt text.
- **`CODE`**: Triggers if it finds ` ``` `, `function`, `class`, or `error:`.
- **`COMPLEX`**: Triggers if it finds words like `analyze`, `summarize`, `architecture`, or if the text is long.
- **`MASSIVE`**: Triggers for payloads exceeding 10,000 text tokens.
- **`SIMPLE`**: Default fallback for standard chat.

## Dynamic Fallback Chains (`registry.go`)

Based on the intent, the router builds a prioritized list of providers. 
For example, a `CODE` task relies heavily on Groq (for speed) and Anthropic (for capability):
1. Groq (`llama-3.3-70b-versatile`)
2. Anthropic (`claude-3-5-sonnet-latest`)
3. OpenAI (`gpt-4o`)
4. Mistral (`mistral-large-2512`)

If Groq returns a `413 Payload Too Large` or `429 Too Many Requests`, the proxy instantly catches the error and silently cascades to Anthropic in the exact same stream, with zero interruption to the user.

## Dynamic Model Discovery (`discovery.go`)

Providers constantly deprecate models (e.g., Groq deprecating `llama3-70b-8192`). 
To ensure the router never attempts to hit a dead endpoint:
- On `start.bat` launch, the router queries the `GET /v1/models` endpoint for Groq, Mistral, Anthropic, Cerebras, etc.
- It caches the *live* list of currently supported models in memory.
- When `registry.go` needs a model (e.g., "give me the best LLaMA 3.3 model on Groq"), it uses `discovery.GetBestModel("groq", "llama-3.3")` to dynamically retrieve the exact, active model string rather than relying on brittle hardcoded strings.
