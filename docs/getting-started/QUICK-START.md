# CTA AI Router - Quick Start Guide

Welcome to the **CTA AI Router**. This robust, high-performance Golang backend serves as a seamless OpenAI-compatible proxy that drastically reduces LLM API costs through context compression and intelligent routing.

## Prerequisites
- Go 1.25.12+
- Node.js (for the frontend dashboard)
- Air (Live reloader for Go)

## 1. Start the Golang Backend
The entire backend runs on a lightweight SQLite database and proxies requests on port `20128`.

```bash
cd "AI Router Golang"
cd backend
air
```

The server will initialize the SQLite database at `C:\Users\<user>\Documents\.cta-ai-nexus\router.db` and start the **Dynamic Model Discovery Engine** to ping all providers for live models.

## 2. Start the Frontend Dashboard
The React frontend provides a beautiful UI for managing API keys and viewing live telemetry.

```bash
cd "AI Router Golang"
cd frontend
npm run dev
```
Visit `http://localhost:5173` in your browser.

## 3. Configure API Keys
1. Open the frontend dashboard.
2. Navigate to the **Keys** section.
3. Enter your API keys for Groq, Mistral, Anthropic, Gemini, or Cerebras.
4. The router instantly securely stores these in the local SQLite database.

## 4. Usage in Cursor / IDEs
To use the router in any IDE like Cursor:
- Set the Base URL to: `http://localhost:20128`
- Model Name: `cta-ai-nexus`
- API Key: `any-string` (Authentication is handled internally via your saved keys).

You are now ready to enjoy mathematically optimized LLM requests!
