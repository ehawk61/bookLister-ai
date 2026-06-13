# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

bookLister-ai is a Go project. The `.gitignore` is configured for Go binaries, test artifacts, and a `.env` file for local secrets.

## Commands

```bash
# Build
go build ./...

# Run tests
go test ./...

# Run a single test
go test ./path/to/package -run TestName

# Lint (install golangci-lint if not present)
golangci-lint run
```

## Architecture

This repository is in early setup — no application code exists yet. As the codebase grows, update this section with:
- Entry points (e.g. `cmd/` binaries)
- Core packages and their responsibilities
- External dependencies and integrations (APIs, databases)
- Configuration and environment variable conventions
