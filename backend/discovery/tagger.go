package discovery

import (
	"log"
	"ainexusrouter-core/db"
)

const (
	TagVision    = "vision"
	TagCoding    = "coding"
	TagReasoning = "reasoning"
	TagComplex   = "complex"
	TagFast      = "fast"
	TagFree      = "free"
	TagLowCost   = "low-cost"
	TagPaid      = "paid"
	TagBest      = "best"
)

// SeedKnownTags populates the model_tags database with explicit capabilities for known models
func SeedKnownTags() {
	log.Println("[Discovery] Seeding explicit tags for known models...")

	// 1. Define explicit mappings based on latest 2026 data
	seedData := []struct {
		Provider string
		Model    string
		Tags     []string
	}{
		// OpenAI
		{"openai", "gpt-4o", []string{TagVision, TagCoding, TagComplex, TagFast, TagPaid, TagBest}},
		{"openai", "gpt-4o-mini", []string{TagVision, TagFast, TagLowCost}},
		{"openai", "o1-preview", []string{TagReasoning, TagComplex, TagCoding, TagPaid, TagBest}},
		{"openai", "o1-mini", []string{TagReasoning, TagCoding, TagLowCost}},
		{"openai", "gpt-4-turbo", []string{TagVision, TagCoding, TagComplex, TagPaid}},

		// Anthropic
		{"anthropic", "claude-3-5-sonnet-20240620", []string{TagVision, TagCoding, TagComplex, TagPaid, TagBest}},
		{"anthropic", "claude-3-haiku-20240307", []string{TagVision, TagFast, TagLowCost}},
		{"anthropic", "claude-3-opus-20240229", []string{TagVision, TagComplex, TagPaid}},

		// Google Gemini
		{"gemini", "gemini-1.5-pro", []string{TagVision, TagCoding, TagComplex, TagPaid, TagBest}},
		{"gemini", "gemini-1.5-flash", []string{TagVision, TagFast, TagFree}},

		// Groq (Llama and others) - All currently free tier
		{"groq", "llama-3.1-70b-versatile", []string{TagCoding, TagComplex, TagFree}},
		{"groq", "llama-3.1-8b-instant", []string{TagFast, TagFree}},
		{"groq", "llama-3.2-11b-vision-preview", []string{TagVision, TagFast, TagFree}},
		{"groq", "llama-3.2-90b-vision-preview", []string{TagVision, TagComplex, TagFree}},
		{"groq", "mixtral-8x7b-32768", []string{TagFast, TagFree}},
		{"groq", "gemma2-9b-it", []string{TagFast, TagFree}},

		// Mistral
		{"mistral", "mistral-large-latest", []string{TagComplex, TagPaid, TagBest}},
		{"mistral", "open-mistral-nemo", []string{TagFast, TagLowCost}},
		{"mistral", "codestral-latest", []string{TagCoding, TagFast, TagLowCost}},
		{"mistral", "pixtral-12b-2409", []string{TagVision, TagFast, TagLowCost}},
		{"mistral", "mistral-small-latest", []string{TagCoding, TagFast, TagLowCost}},
		{"mistral", "mistral-medium-latest", []string{TagCoding, TagComplex, TagPaid}},

		// Cerebras
		{"cerebras", "llama3.1-70b", []string{TagCoding, TagComplex, TagFree}},
		{"cerebras", "llama3.1-8b", []string{TagFast, TagFree}},

		// NVIDIA (NIM Free API Credits)
		{"nvidia", "meta/llama-3.1-70b-instruct", []string{TagCoding, TagComplex, TagFree}},
		{"nvidia", "meta/llama-3.1-8b-instruct", []string{TagFast, TagFree}},
		{"nvidia", "meta/llama-3.1-405b-instruct", []string{TagCoding, TagComplex, TagFree, TagBest}},
		{"nvidia", "meta/llama-3.2-90b-vision-instruct", []string{TagVision, TagComplex, TagFree}},
		{"nvidia", "meta/llama-3.2-11b-vision-instruct", []string{TagVision, TagFast, TagFree}},
		{"nvidia", "microsoft/phi-3-mini-128k-instruct", []string{TagFast, TagFree}},
		{"nvidia", "microsoft/phi-3-vision-128k-instruct", []string{TagVision, TagFast, TagFree}},
		{"nvidia", "google/gemma-2-9b-it", []string{TagFast, TagFree}},
		{"nvidia", "mistralai/mixtral-8x22b-instruct-v0.1", []string{TagComplex, TagFree}},
		{"nvidia", "nvidia/nemotron-4-340b-instruct", []string{TagComplex, TagFree}},

		// DeepSeek
		{"deepseek", "deepseek-v4-pro", []string{TagCoding, TagComplex, TagLowCost, TagBest}},
		{"deepseek", "deepseek-v4-flash", []string{TagFast, TagLowCost}},
		{"deepseek", "deepseek-r1", []string{TagReasoning, TagComplex, TagLowCost}},
		{"deepseek", "deepseek-coder", []string{TagCoding, TagLowCost}},

		// X.AI
		{"xai", "grok-4.5", []string{TagComplex, TagCoding, TagPaid, TagBest}},
		{"xai", "grok-4.3", []string{TagFast, TagLowCost}},
		{"xai", "grok-4.20-reasoning", []string{TagReasoning, TagComplex, TagPaid}},
		{"xai", "grok-imagine", []string{TagVision, TagPaid}},

		// Together AI
		{"together", "meta-llama/Meta-Llama-3.1-70B-Instruct-Turbo", []string{TagCoding, TagComplex, TagLowCost}},
		{"together", "meta-llama/Meta-Llama-3.1-8B-Instruct-Turbo", []string{TagFast, TagFree}},
		{"together", "meta-llama/Meta-Llama-3.1-405B-Instruct-Turbo", []string{TagCoding, TagComplex, TagPaid, TagBest}},
		{"together", "Qwen/Qwen2.5-Coder-32B-Instruct", []string{TagCoding, TagFast, TagLowCost}},
		{"together", "Qwen/Qwen2.5-72B-Instruct", []string{TagComplex, TagLowCost}},
		{"together", "deepseek-ai/DeepSeek-V3", []string{TagCoding, TagComplex, TagLowCost}},

		// Groq (Additional Models)
		{"groq", "allam-2-7b", []string{TagFast, TagFree}},
		{"groq", "openai/gpt-oss-20b", []string{TagFast, TagFree}},
		{"groq", "groq/compound", []string{TagFast, TagFree}},
		{"groq", "qwen/qwen3.6-27b", []string{TagComplex, TagFree}},
		{"groq", "meta-llama/llama-4-scout-17b-16e-instruct", []string{TagCoding, TagFast, TagFree}},
		{"groq", "llama-3.3-70b-versatile", []string{TagCoding, TagComplex, TagFree, TagBest}},
		{"groq", "groq/compound-mini", []string{TagFast, TagFree}},
		{"groq", "canopylabs/orpheus-v1-english", []string{TagFast, TagFree}},
		{"groq", "qwen/qwen3-32b", []string{TagComplex, TagFree}},
		{"groq", "openai/gpt-oss-safeguard-20b", []string{TagFast, TagFree}},
		{"groq", "openai/gpt-oss-120b", []string{TagComplex, TagFree, TagBest}},
		{"groq", "canopylabs/orpheus-arabic-saudi", []string{TagFast, TagFree}},

		// Llm7
		{"llm7", "llm7-default", []string{TagFast, TagLowCost}},
		{"llm7", "llm7-advanced", []string{TagComplex, TagCoding, TagPaid, TagBest}},

		// Mistral (Additional Models)
		{"mistral", "mistral-medium-2505", []string{TagComplex, TagPaid}},
		{"mistral", "mistral-medium-2508", []string{TagComplex, TagPaid}},
		{"mistral", "open-mistral-nemo-2407", []string{TagFast, TagLowCost}},
		{"mistral", "mistral-tiny-2407", []string{TagFast, TagLowCost}},
		{"mistral", "mistral-tiny-latest", []string{TagFast, TagLowCost}},
		{"mistral", "codestral-2508", []string{TagCoding, TagFast, TagLowCost}},
		{"mistral", "mistral-code-latest", []string{TagCoding, TagFast, TagLowCost}},
		{"mistral", "mistral-code-fim-latest", []string{TagCoding, TagFast, TagLowCost}},
		{"mistral", "devstral-2512", []string{TagCoding, TagComplex, TagPaid}},
		{"mistral", "devstral-medium-latest", []string{TagCoding, TagComplex, TagPaid}},
		{"mistral", "devstral-latest", []string{TagCoding, TagComplex, TagPaid}},
		{"mistral", "mistral-code-agent-latest", []string{TagCoding, TagPaid}},
		{"mistral", "mistral-small-2603", []string{TagFast, TagLowCost}},
		{"mistral", "mistral-vibe-cli-fast", []string{TagFast, TagLowCost}},
		{"mistral", "magistral-small-latest", []string{TagFast, TagLowCost}},
		{"mistral", "magistral-medium-2509", []string{TagComplex, TagPaid}},
		{"mistral", "magistral-medium-latest", []string{TagComplex, TagPaid}},
		{"mistral", "voxtral-small-2507", []string{TagFast, TagLowCost}},
		{"mistral", "voxtral-small-latest", []string{TagFast, TagLowCost}},
		{"mistral", "labs-leanstral-1-5-1", []string{TagFast, TagLowCost}},
		{"mistral", "labs-leanstral-1-5", []string{TagFast, TagLowCost}},
		{"mistral", "mistral-large-2512", []string{TagComplex, TagPaid, TagBest}},
		{"mistral", "ministral-3b-2512", []string{TagFast, TagLowCost}},
		{"mistral", "ministral-3b-latest", []string{TagFast, TagLowCost}},
		{"mistral", "ministral-8b-2512", []string{TagFast, TagLowCost}},
		{"mistral", "ministral-8b-latest", []string{TagFast, TagLowCost}},
		{"mistral", "ministral-14b-2512", []string{TagFast, TagLowCost}},
		{"mistral", "ministral-14b-latest", []string{TagFast, TagLowCost}},
		{"mistral", "mistral-medium", []string{TagComplex, TagPaid}},
		{"mistral", "mistral-medium-3-5", []string{TagComplex, TagPaid}},
		{"mistral", "mistral-medium-3.5", []string{TagComplex, TagPaid}},
		{"mistral", "mistral-medium-3", []string{TagComplex, TagPaid}},
		{"mistral", "mistral-medium-2604", []string{TagComplex, TagPaid}},
		{"mistral", "mistral-vibe-cli-latest", []string{TagComplex, TagPaid}},
		{"mistral", "mistral-vibe-cli-with-tools", []string{TagCoding, TagComplex, TagPaid}},
		{"mistral", "magistral-small-2509", []string{TagFast, TagLowCost}},
		{"mistral", "mistral-small-2506", []string{TagFast, TagLowCost}},
		{"mistral", "voxtral-mini-2602", []string{TagFast, TagLowCost}},
		{"mistral", "voxtral-mini-latest", []string{TagFast, TagLowCost}},

		// NVIDIA (Additional Models)
		{"nvidia", "01-ai/yi-large", []string{TagComplex, TagFree}},
		{"nvidia", "abacusai/dracarys-llama-3.1-70b-instruct", []string{TagCoding, TagComplex, TagFree}},
		{"nvidia", "adept/fuyu-8b", []string{TagVision, TagFast, TagFree}},
		{"nvidia", "ai21labs/jamba-1.5-large-instruct", []string{TagComplex, TagFree}},
		{"nvidia", "aisingapore/sea-lion-7b-instruct", []string{TagFast, TagFree}},
		{"nvidia", "baai/bge-m3", []string{TagFast, TagFree}},
		{"nvidia", "bigcode/starcoder2-15b", []string{TagCoding, TagFast, TagFree}},
		{"nvidia", "bytedance/seed-oss-36b-instruct", []string{TagComplex, TagFree}},
		{"nvidia", "databricks/dbrx-instruct", []string{TagCoding, TagComplex, TagFree}},
		{"nvidia", "deepseek-ai/deepseek-coder-6.7b-instruct", []string{TagCoding, TagFast, TagFree}},
		{"nvidia", "deepseek-ai/deepseek-v4-flash", []string{TagFast, TagFree}},
		{"nvidia", "deepseek-ai/deepseek-v4-pro", []string{TagCoding, TagComplex, TagFree, TagBest}},
		{"nvidia", "google/codegemma-1.1-7b", []string{TagCoding, TagFast, TagFree}},
		{"nvidia", "google/codegemma-7b", []string{TagCoding, TagFast, TagFree}},
		{"nvidia", "google/deplot", []string{TagVision, TagFast, TagFree}},
		{"nvidia", "google/diffusiongemma-26b-a4b-it", []string{TagVision, TagComplex, TagFree}},
		{"nvidia", "google/gemma-2-2b-it", []string{TagFast, TagFree}},
		{"nvidia", "google/gemma-2b", []string{TagFast, TagFree}},
		{"nvidia", "google/gemma-3-12b-it", []string{TagFast, TagFree}},
		{"nvidia", "google/gemma-3-4b-it", []string{TagFast, TagFree}},
		{"nvidia", "google/gemma-3n-e2b-it", []string{TagFast, TagFree}},
		{"nvidia", "google/gemma-3n-e4b-it", []string{TagFast, TagFree}},
		{"nvidia", "google/gemma-4-31b-it", []string{TagComplex, TagFree}},
		{"nvidia", "google/recurrentgemma-2b", []string{TagFast, TagFree}},
		{"nvidia", "ibm/granite-3.0-3b-a800m-instruct", []string{TagFast, TagFree}},
		{"nvidia", "ibm/granite-3.0-8b-instruct", []string{TagFast, TagFree}},
		{"nvidia", "ibm/granite-34b-code-instruct", []string{TagCoding, TagComplex, TagFree}},
		{"nvidia", "ibm/granite-8b-code-instruct", []string{TagCoding, TagFast, TagFree}},
		{"nvidia", "meta/codellama-70b", []string{TagCoding, TagComplex, TagFree}},
		{"nvidia", "meta/llama-3.2-1b-instruct", []string{TagFast, TagFree}},
		{"nvidia", "meta/llama-3.2-3b-instruct", []string{TagFast, TagFree}},
		{"nvidia", "meta/llama-3.3-70b-instruct", []string{TagCoding, TagComplex, TagFree}},
		{"nvidia", "meta/llama-4-maverick-17b-128e-instruct", []string{TagCoding, TagFast, TagFree}},
		{"nvidia", "meta/llama-guard-4-12b", []string{TagFast, TagFree}},
		{"nvidia", "meta/llama2-70b", []string{TagComplex, TagFree}},
		{"nvidia", "microsoft/kosmos-2", []string{TagVision, TagFast, TagFree}},
		{"nvidia", "microsoft/phi-3.5-moe-instruct", []string{TagFast, TagFree}},
		{"nvidia", "minimaxai/minimax-m2.7", []string{TagFast, TagFree}},
		{"nvidia", "minimaxai/minimax-m3", []string{TagComplex, TagFree}},
		{"nvidia", "mistralai/codestral-22b-instruct-v0.1", []string{TagCoding, TagComplex, TagFree}},
		{"nvidia", "mistralai/ministral-14b-instruct-2512", []string{TagFast, TagFree}},
		{"nvidia", "mistralai/mistral-7b-instruct-v0.3", []string{TagFast, TagFree}},
		{"nvidia", "mistralai/mistral-large", []string{TagComplex, TagFree}},
		{"nvidia", "mistralai/mistral-large-2-instruct", []string{TagComplex, TagFree}},
		{"nvidia", "mistralai/mistral-large-3-675b-instruct-2512", []string{TagComplex, TagFree, TagBest}},
		{"nvidia", "mistralai/mistral-medium-3.5-128b", []string{TagComplex, TagFree}},
		{"nvidia", "mistralai/mistral-nemotron", []string{TagFast, TagFree}},
		{"nvidia", "mistralai/mistral-small-4-119b-2603", []string{TagComplex, TagFree}},
		{"nvidia", "mistralai/mixtral-8x22b-v0.1", []string{TagComplex, TagFree}},
		{"nvidia", "mistralai/mixtral-8x7b-instruct-v0.1", []string{TagFast, TagFree}},
		{"nvidia", "moonshotai/kimi-k2.6", []string{TagComplex, TagFree}},
		{"nvidia", "nv-mistralai/mistral-nemo-12b-instruct", []string{TagFast, TagFree}},
		{"nvidia", "nvidia/ai-synthetic-video-detector", []string{TagVision, TagFast, TagFree}},
		{"nvidia", "nvidia/cosmos-reason2-8b", []string{TagReasoning, TagFast, TagFree}},
		{"nvidia", "nvidia/gliner-pii", []string{TagFast, TagFree}},
		{"nvidia", "nvidia/ising-calibration-1-35b-a3b", []string{TagComplex, TagFree}},
		{"nvidia", "nvidia/llama-3.1-nemoguard-8b-content-safety", []string{TagFast, TagFree}},
		{"nvidia", "nvidia/llama-3.1-nemoguard-8b-topic-control", []string{TagFast, TagFree}},
		{"nvidia", "nvidia/llama-3.1-nemotron-51b-instruct", []string{TagComplex, TagFree}},
		{"nvidia", "nvidia/llama-3.1-nemotron-70b-instruct", []string{TagComplex, TagFree}},
		{"nvidia", "nvidia/llama-3.1-nemotron-nano-8b-v1", []string{TagFast, TagFree}},
		{"nvidia", "nvidia/llama-3.1-nemotron-nano-vl-8b-v1", []string{TagVision, TagFast, TagFree}},
		{"nvidia", "nvidia/llama-3.1-nemotron-safety-guard-8b-v3", []string{TagFast, TagFree}},
		{"nvidia", "nvidia/llama-3.1-nemotron-ultra-253b-v1", []string{TagComplex, TagFree}},
		{"nvidia", "nvidia/llama-3.3-nemotron-super-49b-v1", []string{TagComplex, TagFree}},
		{"nvidia", "nvidia/llama-3.3-nemotron-super-49b-v1.5", []string{TagComplex, TagFree}},
		{"nvidia", "nvidia/llama3-chatqa-1.5-70b", []string{TagComplex, TagFree}},
		{"nvidia", "nvidia/mistral-nemo-minitron-8b-8k-instruct", []string{TagFast, TagFree}},
		{"nvidia", "nvidia/nemoretriever-parse", []string{TagFast, TagFree}},
		{"nvidia", "nvidia/nemotron-3-nano-30b-a3b", []string{TagComplex, TagFree}},
		{"nvidia", "nvidia/nemotron-3-nano-omni-30b-a3b-reasoning", []string{TagReasoning, TagComplex, TagFree}},
		{"nvidia", "nvidia/nemotron-3-super-120b-a12b", []string{TagComplex, TagFree}},
		{"nvidia", "nvidia/nemotron-3-ultra-550b-a55b", []string{TagComplex, TagFree, TagBest}},
		{"nvidia", "nvidia/nemotron-3.5-content-safety", []string{TagFast, TagFree}},
		{"nvidia", "nvidia/nemotron-4-340b-reward", []string{TagComplex, TagFree}},
		{"nvidia", "nvidia/nemotron-mini-4b-instruct", []string{TagFast, TagFree}},
		{"nvidia", "nvidia/nemotron-nano-12b-v2-vl", []string{TagVision, TagFast, TagFree}},
		{"nvidia", "nvidia/nemotron-nano-3-30b-a3b", []string{TagComplex, TagFree}},
		{"nvidia", "nvidia/nemotron-parse", []string{TagFast, TagFree}},
		{"nvidia", "nvidia/neva-22b", []string{TagVision, TagComplex, TagFree}},
		{"nvidia", "nvidia/nvclip", []string{TagVision, TagFast, TagFree}},
		{"nvidia", "nvidia/nvidia-nemotron-nano-9b-v2", []string{TagFast, TagFree}},
		{"nvidia", "nvidia/riva-translate-4b-instruct", []string{TagFast, TagFree}},
		{"nvidia", "nvidia/riva-translate-4b-instruct-v1.1", []string{TagFast, TagFree}},
		{"nvidia", "nvidia/vila", []string{TagVision, TagFast, TagFree}},
		{"nvidia", "openai/gpt-oss-120b", []string{TagComplex, TagFree, TagBest}},
		{"nvidia", "openai/gpt-oss-20b", []string{TagFast, TagFree}},
		{"nvidia", "qwen/qwen3-next-80b-a3b-instruct", []string{TagComplex, TagFree}},
		{"nvidia", "qwen/qwen3.5-122b-a10b", []string{TagComplex, TagFree}},
		{"nvidia", "qwen/qwen3.5-397b-a17b", []string{TagComplex, TagFree}},
		{"nvidia", "sarvamai/sarvam-m", []string{TagFast, TagFree}},
		{"nvidia", "stepfun-ai/step-3.5-flash", []string{TagFast, TagFree}},
		{"nvidia", "stepfun-ai/step-3.7-flash", []string{TagFast, TagFree}},
		{"nvidia", "upstage/solar-10.7b-instruct", []string{TagFast, TagFree}},
		{"nvidia", "writer/palmyra-creative-122b", []string{TagComplex, TagFree}},
		{"nvidia", "writer/palmyra-fin-70b-32k", []string{TagComplex, TagFree}},
		{"nvidia", "writer/palmyra-med-70b", []string{TagComplex, TagFree}},
		{"nvidia", "writer/palmyra-med-70b-32k", []string{TagComplex, TagFree}},
		{"nvidia", "z-ai/glm-5.2", []string{TagComplex, TagFree}},
		{"nvidia", "zyphra/zamba2-7b-instruct", []string{TagFast, TagFree}},
	}

	for _, entry := range seedData {
		err := db.SetTagsForModel(entry.Provider, entry.Model, entry.Tags)
		if err != nil {
			log.Printf("[Discovery] Failed to seed tags for %s/%s: %v", entry.Provider, entry.Model, err)
		}
	}
}
