# ClickToAutomate AI Nexus Router (Golang Edition)

ClickToAutomate AI Nexus Router is a free, open-source AI Gateway designed to route, manage, and optimize multi-provider LLM requests. Built with high-performance Golang, it is paired with a React frontend and distributed as both a standalone executable and a desktop application via Wails (as Go integrates much better natively with Wails than Tauri).

## Overview

ClickToAutomate AI Nexus Router aims to provide maximum performance, minimal resource usage, and effortless portability for all your AI integration needs.

By leveraging Go's single-binary compilation, users can download just one executable and immediately access the full power of ClickToAutomate AI Nexus Router's auto-fallback, routing strategies, and token compression.

### Features

- **Single Binary Distribution**: The React dashboard is fully embedded into the Go executable. Download, run, and start building.
- **Wails Desktop App**: A seamless native desktop experience built with Wails, natively embedding the Go backend without the need for complex sidecar architectures.
- **Smart Routing & Resilience**: Auto-fallback across multiple providers ensuring you never hit rate limits.
- **One Endpoint for Every Tool**: Connect Claude Code, Codex, Cursor, and any OpenAI-compatible tool to a single local endpoint.
- **Logs Section**: View real-time backend logs directly in the dashboard for debugging and monitoring.
- **Image Processing**: Seamlessly analyze images alongside prompts using vision-capable models (e.g., `gpt-4o`, `claude-3.5-sonnet`).

## Getting Started

### CLI Application (Headless)
Run the compiled `AI Router CLI.exe` from your terminal or command prompt. It will start the proxy server on port `20128` and stream logs to your terminal. Even in headless mode, you can still manage your providers by opening `http://localhost:20128` in any web browser.

**Available Flags:**
- `--port <number>`: Override the default port (e.g. `--port 8080`)
- `--db <path>`: Specify a custom SQLite database path
- `--help`: View all available commands

### GUI Desktop Application
Launch the `AI Router GUI.exe` desktop app. The native window (powered by Wails) will automatically manage the Go backend in the background and present the dashboard directly on your desktop. Logs are accessible via the built-in Logs tab.

## Project Structure

- `/backend` - The Golang core handling API routes, provider management, and serving the embedded frontend.
- `/frontend` - The React SPA (Single Page Application) dashboard.
- `/` (Root) - The Wails desktop wrapper that directly compiles the Go backend and React frontend into a single native application.

## Changelog

### v1.0.0
- **Dual Build System**: Now ships with both a lightweight **CLI executable** and a full **GUI desktop application** (Wails).
- **Advanced Logging**: Real-time memory log buffer in the UI, plus automatic background physical `logs.txt` writer inside the database folder.
- **Enhanced Playground**: Smooth auto-scroll behavior for long streams and quick "Delete Chat" session management.
- **CORS Compatibility**: Proper preflight `OPTIONS` support on proxy endpoints, enabling in-browser web UI fetching directly.
- Fixed **image processing** to support vision-capable models (e.g., `gpt-4o`, `claude-3.5-sonnet`).

## Build Instructions

### Wails Build
To build the desktop application for Windows, run:
```bash
wails build -clean -platform windows -webview2 embed
```

The executable will be available in the `build/bin` directory.

## Documentation

- [System Architecture](ARCHITECTURE.md)
- [Contributing Guidelines](CONTRIBUTING.md)
