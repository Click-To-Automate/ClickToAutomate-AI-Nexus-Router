# Lagrangian Routing & Subgradient Fallbacks

The **ClickToAutomate AI Nexus Router** removes the burden of manually selecting LLM models. It uses **Lagrangian Dual Decomposition** to calculate the absolute mathematically optimal route for every single request, dynamically balancing cost, latency, and quality.

## 1. Complexity Scoring (`intent.go`)

When a request arrives, `AnalyzeComplexity()` scans the concatenated prompt text and calculates a rigorous mathematical complexity score between `0.0` and `1.0`.

- **Simple Context**: A baseline task yields a score of `0.5`.
- **Complex Code**: Finding tokens like `function`, `class`, or `error:` increases the required complexity score to `0.8+`.
- **Massive Payload**: A prompt with >10,000 tokens instantly triggers a `0.95` complexity constraint.

This score acts as the minimum **Quality Constraint ($q_{min}$)** for the mathematical solver.

## 2. The Lagrangian Solver (`lagrangian.go`)

Based on the complexity constraint, the router builds a mathematically sorted fallback chain of available providers.

1. **Intrinsic Constraints**: Each model family (e.g., LLaMA 8B vs GPT-4o) has hardcoded baseline `Cost`, `Latency`, and `Quality` scores.
2. **Objective Function**: The solver calculates $L(x, \lambda)$ for every available model. It attempts to minimize Cost + Latency, while strictly penalizing models whose `Quality` falls below the task's `ComplexityScore`.
3. **Multiplier Adjustments ($\lambda$)**: The solver applies live Lagrange Multipliers (penalties) for each provider.

## 3. Subgradient Coordination Loop (`stream.go`)

The AI Nexus Router is fully self-healing.

If the solver chooses Anthropic, but Anthropic returns a `429 Too Many Requests` or `413 Payload Too Large`, the router catches the error instantly.

- **Zero Interruption**: The proxy silently cascades to the next best model mathematically sorted by the solver.
- **Feedback Loop**: It instantly fires `IncrementPenalty(provider)` to spike the Lagrange multiplier for Anthropic. The next request to hit the router will see Anthropic's objective score severely penalized, forcing the decomposition solver to route traffic away from the congested provider until the penalties eventually decay.

## Dynamic Model Discovery (`discovery.go`)

Providers constantly deprecate models. To ensure the solver never attempts to hit a dead endpoint:
- On `wails dev` or executable launch, the router queries the `GET /v1/models` endpoint for your providers.
- It caches the *live* list of currently supported models in memory.
- When the solver picks a target (e.g., "give me an 8B model on Groq"), it uses `discovery.GetBestModel("groq", "8b")` to dynamically retrieve the active model string rather than relying on brittle hardcoded strings.
