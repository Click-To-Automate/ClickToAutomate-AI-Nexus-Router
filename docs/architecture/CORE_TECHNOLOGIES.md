# Core Technologies

**ClickToAutomate AI Nexus Router** uses two core components to execute an efficient LLM cascade: **RTK + Caveman Stacked Compression** (to reduce token size) and **Lagrangian Dual Decomposition** (to decide which model gets the request). [1, 2, 3] 

Here is a breakdown of how both technologies function within the router.

---

## Part 1: RTK + Caveman Stacked Compression (The Token Reducer)

When running coding assistants or cascading pipelines, your prompts contain two different types of text: noisy machine code (logs, errors, terminals) and human-written conversational context. [3, 4] 

Instead of using one generic compressor, the AI Nexus Router feeds the prompt through a stacked pipeline where two specialized algorithmic engines run sequentially to compress data by 15% to 95%. [2, 3] 

Raw Multi-Modal Prompt ──> [ 1. RTK Engine ] ──> [ 2. Caveman Engine ] ──> Compressed Prompt 
                        (Strips terminal logs)   (Condenses human prose)   (78% - 94% tokens saved)

### Step 1: The RTK Engine (Command & Terminal Filter) [3] 

* **What it targets**: Machine-generated outputs such as terminal logs, build scripts, package manager noise, stack traces, Git diffs, and Docker container outputs. [3] 
* **How it works**: It acts as an AST (Abstract Syntax Tree) and heuristic parser. It recognizes the structures of known terminal commands and strips out redundant strings, repetitive padding lines, and non-essential trace data while leaving the core error or structural payload intact. [3, 5] 
* **Savings**: It yields roughly 60% to 90% savings on raw terminal outputs. [3] 

### Step 2: The Caveman Engine (Natural Language Condenser)

* **What it targets**: The remaining human conversational prose, assistant histories, and instructions.
* **How it works**: It applies programmatic textual compression rules to strip conversational fluff. It eliminates stop-words, filler language, and redundant grammatical clauses, leaving the text looking like basic, blunt "caveman speech" that an LLM can still interpret with 100% semantic accuracy.
* **Savings**: It yields roughly 46% savings on standard natural language. [4, 6, 7, 8] 

### The Cumulative "Stacked" Effect

By chaining them together (RTK -> Caveman), machine output is cleaned first, followed by prose condensation. Mathematically, the AI Nexus Router calculates the combined pipeline efficiency as:

$$\text{Total Savings} = 1 - (1 - 0.80_{\text{RTK}}) \times (1 - 0.46_{\text{Caveman}}) \approx 89.2\% \text{ saved tokens}$$ 
[3, 4, 8] 

---

## Part 2: Lagrangian Dual Decomposition (The Intelligent Router)

Once the prompt is compressed, the AI Nexus Router must choose the best model cascade (e.g., Llama 8B or a frontier model). It models this routing dilemma as a constrained optimization problem rather than using simple rule-based forwarding. [1, 9] 

### The Problem (The Primal Problem)

You want to serve an AI query while balancing conflicting realities:
1. **Minimize Cost & Latency**: Use the absolute cheapest/fastest model possible.
2. **Maximize Quality (Performance Constraints)**: Ensure the model chosen is smart enough to handle the prompt accurately.
3. **Capacity Constraints**: Keep track of rate limits and token quotas across multiple providers. [1, 2, 9] 

Solving this globally for thousands of fast-moving API requests simultaneously is mathematically complex (NP-hard). [10] 

### The Solution: Lagrangian Dual Decomposition [11] 

The AI Nexus Router utilizes Lagrangian Duality to turn this hard problem inside out. [1] 

1. **Relaxing Constraints via Multipliers ($\lambda$)**: The framework takes your hard constraints (like "Quality must be higher than 85%" or "Do not exceed a 10-cent budget") and converts them into soft penalties using mathematical variables called Lagrange multipliers. If a model configuration violates your rules, its penalty price skyrockets. [11, 12, 13] 
2. **Decomposition (Splitting the Chunk)**: The "Decomposition" aspect means the master routing engine breaks one massive, complex network routing problem into independent, parallelized sub-problems. It evaluates each available model option individually and concurrently on separate processor cores. [10, 12] 
3. **Adaptive Convergence**: The system utilizes a subgradient coordination loop. It dynamically tweaks the multipliers step-by-step based on live network feedback (e.g., if OpenAI suddenly throttles your rate limit, the multiplier penalty for OpenAI increases instantly). [1, 10, 11] 

### Why This Matters for Your Cascade

Instead of randomly guessing if Llama 8B can handle a prompt, the AI Nexus Router predicts the query's complexity. If the Lagrangian optimization determines that a cheap model can safely meet your quality threshold, it routes it there. If the prompt is too dense, it elevates it to a higher tier. [1, 6, 9] 

This optimization technique achieves over a 10% reduction in API compute costs while boosting total response accuracy. [1] 
