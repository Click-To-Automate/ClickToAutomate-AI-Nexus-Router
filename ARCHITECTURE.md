# System Architecture

The Golang version of ClickToAutomate AI Nexus Router utilizes a decoupled architecture where the core routing engine and the web dashboard are separated but seamlessly combined for distribution.

## High-Level Architecture

The system consists of three main components:

1. **Golang Backend Core**: The heart of the application. It provides the REST APIs (like `/v1/chat/completions`), manages state, handles AI provider connections, executes routing strategies, and enforces token compression.
2. **React SPA Frontend**: A client-side React application that provides the visual dashboard for managing providers, monitoring analytics, and configuring routing combos.
3. **Tauri Desktop Shell**: A Rust-based desktop application wrapper that bundles the Go backend and presents the React frontend in a native OS window.

## Deployment Modes

ClickToAutomate AI Nexus Router (Go Edition) is designed to be distributed in two primary ways to maximize accessibility.

### 1. Standalone Executable (Zero-Install)
In this mode, the React frontend is compiled into static assets (HTML, CSS, JS) and embedded directly into the Go binary using the `go:embed` directive. 

When a user executes the resulting binary:
- The Go HTTP server starts on a local port (e.g., `20128`).
- API requests to `/v1/*` are handled by the core routing logic.
- Browser requests to `/` or `/dashboard` are served the embedded static React files.

This provides a true "zero installation" experience. The user downloads a single `.exe`, runs it, and manages the system via their web browser, whilst routing their AI IDEs to `localhost:20128/v1`.

### 2. Desktop Application (Tauri Sidecar)
For users preferring a native application experience with system tray support and window management, the Tauri app is the solution.

- The React frontend is bundled directly into the Tauri application.
- The compiled Go binary is bundled as a Tauri **sidecar**.
- When the Tauri app launches, it automatically spawns the Go sidecar process in the background.
- The native UI communicates with the local Go backend over HTTP.
- Closing the Tauri app automatically cleans up and shuts down the Go sidecar process.

## Technology Stack

- **Backend**: Golang 1.22+
- **Frontend**: React, Vite, TypeScript
- **Desktop**: Tauri v2, Rust
