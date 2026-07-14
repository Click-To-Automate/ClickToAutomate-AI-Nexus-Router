# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-07-14

### Added
- **Initial Release** of AI Nexus Router.
- **Dynamic Traffic Controller**: Intelligently routes incoming requests between local models and remote providers (OpenAI, Anthropic, local nodes) based on query complexity.
- **Global Sanitization Layer**: Pre-processes context, applies system instructions globally, and enforces context window limitations.
- **GUI and CLI interfaces**: Bundled with a desktop application built on Wails, alongside a headless standalone CLI runner.
- **Failover / Fallback Mechanisms**: Automatically retries prompts with alternative providers if the primary provider times out or fails.
- **Streaming Proxy**: Fully transparent HTTP proxying and streaming for chat completion workflows, designed to drop-in replace standard AI endpoints in popular IDEs.