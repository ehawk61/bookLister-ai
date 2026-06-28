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

### Entry Point

- `cmd/server/main.go` — HTTP server, wires middleware and routes, listens on `PORT` (default `8080`)

### Core Packages

- `internal/handler` — HTTP handler functions (e.g. `Healthcheck`)
- `internal/middleware` — HTTP middleware (`CorrelationID`, `APIVersion`)
- `internal/response` — Shared response helpers (`ErrorResponse`, `WriteError`)

### Middleware Chain

Requests pass through: `CorrelationID` → `APIVersion` → handler. The correlation ID middleware is outermost and rejects requests missing a valid UUIDv4 `X-Correlation-ID` header.

### Configuration

- `PORT` env var — server listen port (default `8080`)
