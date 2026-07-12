# Contributing to ClickToAutomate AI Nexus Router (Golang Edition)

We welcome contributions from the community! This document outlines how to get your development environment set up and how to submit your changes.

## Prerequisites

To work on all aspects of the project, you will need:
- [Go 1.22+](https://go.dev/doc/install)
- [Node.js 18+](https://nodejs.org/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) (Required for the desktop app)

## Repository Structure

The project is split into three main directories:
- `/backend`: The Golang AI router and HTTP server.
- `/frontend`: The React dashboard.
- `/`: The root directory contains the Wails desktop application configuration.

## Local Development Workflow

### 1. Working on the Core (Backend + Frontend)
If you are developing features for the AI router or the web dashboard, you typically do not need to run the Wails shell. You can develop in pure web mode.

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

1. Ensure you have the Wails CLI installed (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`).
2. Run the Wails dev server from the root directory:
```bash
wails dev
```

## Submitting Pull Requests

1. Fork the repository and create your branch from `main`.
2. Ensure your code passes all linting and formatting rules (`gofmt` for Go, `eslint/prettier` for frontend).
3. Add tests for new core routing features.
4. Write descriptive commit messages outlining *what* was changed and *why*.
