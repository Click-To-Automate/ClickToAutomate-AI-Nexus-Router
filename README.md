# ClickToAutomate AI Nexus Router (Golang Edition)

ClickToAutomate AI Nexus Router is a free, open-source AI Gateway designed to route, manage, and optimize multi-provider LLM requests. Built with high-performance Golang, it is paired with a React frontend and distributed as both a standalone executable and a desktop application via Tauri.

## Overview

ClickToAutomate AI Nexus Router aims to provide maximum performance, minimal resource usage, and effortless portability for all your AI integration needs.

By leveraging Go's single-binary compilation, users can download just one executable and immediately access the full power of ClickToAutomate AI Nexus Router's auto-fallback, routing strategies, and token compression.

### Features

- **Single Binary Distribution**: The React dashboard is fully embedded into the Go executable. Download, run, and start building.
- **Tauri Desktop App**: A seamless native desktop experience built with Tauri, using the Go core as a powerful sidecar engine.
- **Smart Routing & Resilience**: Auto-fallback across multiple providers ensuring you never hit rate limits.
- **One Endpoint for Every Tool**: Connect Claude Code, Codex, Cursor, and any OpenAI-compatible tool to a single local endpoint.

## Getting Started

### Standalone Executable
Simply run the compiled executable. The Go server will start the AI Gateway and automatically serve the React dashboard on `http://localhost:20128`.

### Desktop Application
Launch the ClickToAutomate AI Nexus Router desktop app. The native GUI will automatically manage the Go backend sidecar in the background.

## Project Structure

- `/backend` - The Golang core handling API routes, provider management, and serving the embedded frontend.
- `/frontend` - The React SPA (Single Page Application) dashboard.
- `/desktop` - The Tauri Rust shell that wraps the web dashboard and bundles the Go binary.

## Documentation

- [System Architecture](ARCHITECTURE.md)
- [Contributing Guidelines](CONTRIBUTING.md)
