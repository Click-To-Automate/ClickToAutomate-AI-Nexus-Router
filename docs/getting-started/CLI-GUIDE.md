# CLI Guide

The **CTA AI Nexus Router** provides a lightweight, headless Command Line Interface (CLI) executable (`AI Router CLI.exe`) for users who want to run the proxy server without a GUI.

## When to use the CLI
- **Server Environments**: Running the router on a headless Linux/Windows server, Raspberry Pi, or Docker container.
- **Background Service**: Keeping the proxy running silently in the background while you code in your IDE (e.g., Cursor, VSCode).
- **Minimal Resource Usage**: The CLI uses significantly less RAM than the GUI desktop app since it doesn't spin up a WebView or native window.

## Usage

Simply run the executable from your terminal or command prompt:

```bash
# Windows
.\release\v1.0.0\"AI Router CLI.exe"
```

### Command-Line Arguments

You can customize the router's behavior using the following flags:

| Flag | Description | Default |
|------|-------------|---------|
| `--port` | Port to run the proxy server and dashboard on. | `20128` (or `30128` in production builds) |
| `--db` | Custom path to the sqlite database file. | `~/.cta-ai-router/providers.db` |
| `--help` | Displays the help message with all available commands. | |

**Example Usage with Arguments:**

```bash
.\release\v1.0.0\"AI Router CLI.exe" --port 8080 --db "C:\MyData\custom-providers.db"
```

Once running, the CLI will output logs directly to your terminal. It will initialize the database, load providers, discover models, and start listening for proxy requests.

## Endpoints

Even though you are using the CLI, the integrated React frontend dashboard is still served!

1. **Proxy Endpoint**: Point your tools (Cursor, Claude Code, scripts) to `http://localhost:20128/v1`
2. **Dashboard**: Open `http://localhost:20128` in your web browser to access the management UI and Playground.

## Background Logging

While running, the CLI automatically saves all application logs to a `logs.txt` file located in the same directory as your `sqlite3` database (usually `~/.cta-ai-router/` or equivalent AppData directory). This ensures you can always review the history of requests and errors, even if the terminal window is closed.
