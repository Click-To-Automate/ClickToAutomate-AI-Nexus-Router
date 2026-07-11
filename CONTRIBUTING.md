# Contributing to ClickToAutomate AI Nexus Router (Golang Edition)

We welcome contributions from the community! This document outlines how to get your development environment set up and how to submit your changes.

## Prerequisites

To work on all aspects of the project, you will need:
- [Go 1.22+](https://go.dev/doc/install)
- [Node.js 18+](https://nodejs.org/)
- [Rust & Cargo](https://rustup.rs/) (Required for the Tauri desktop app)
- [Tauri CLI](https://tauri.app/v1/guides/getting-started/setup/)

## Repository Structure

The project is split into three main directories:
- `/backend`: The Golang AI router and HTTP server.
- `/frontend`: The React dashboard.
- `/desktop`: The Tauri desktop application shell.

## Local Development Workflow

### 1. Working on the Core (Backend + Frontend)
If you are developing features for the AI router or the web dashboard, you typically do not need to run the Tauri shell.

**Start the Frontend (Dev Server):**
```bash
cd frontend
npm install
npm run dev
```

**Start the Backend:**
```bash
cd backend
go run main.go
```
Ensure that the frontend is configured to point its API calls to the Go backend's port.

### 2. Working on the Desktop App
When testing the native desktop experience:

1. You must first build the Go backend executable for your platform. The Tauri `tauri.conf.json` is configured to look for this sidecar binary.
2. Build the Go binary and place it in the `desktop/src-tauri/binaries/` folder (naming conventions apply based on your OS).
3. Run the Tauri dev server:
```bash
cd desktop
npm install
npm run tauri dev
```

## Submitting Pull Requests

1. Fork the repository and create your branch from `main`.
2. Ensure your code passes all linting and formatting rules (`gofmt` for Go, `eslint/prettier` for frontend).
3. Add tests for new core routing features.
4. Write descriptive commit messages outlining *what* was changed and *why*.
