# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

bookLister-ai is a Go project. The `.gitignore` is configured for Go binaries, test artifacts, and a `.env` file for local secrets.

## Commands

```bash
# Build (from middleware-api/)
cd middleware-api && go build ./...

# Run tests (from middleware-api/)
cd middleware-api && go test ./...

# Run a single test
cd middleware-api && go test ./path/to/package -run TestName

# Lint (install golangci-lint if not present)
cd middleware-api && golangci-lint run
```

## Architecture

All Go source lives under `middleware-api/`.

### Entry Point

- `middleware-api/cmd/server/main.go` — HTTP server, wires middleware and routes, listens on `PORT` (default `8080`)

### Core Packages

- `middleware-api/internal/handler` — HTTP handler functions (e.g. `Healthcheck`)
- `middleware-api/internal/middleware` — HTTP middleware (`CorrelationID`, `APIVersion`)
- `middleware-api/internal/response` — Shared response helpers (`ErrorResponse`, `WriteError`)

### API Routes

All endpoints are prefixed with `/api`. Current routes:
- `GET /api/healthcheck` — returns `{"status":"ok"}`

### Middleware Chain

Requests pass through: `CorrelationID` → `APIVersion` → handler. The correlation ID middleware is outermost and rejects requests missing a valid UUIDv4 `X-Correlation-ID` header.

### Configuration

- `PORT` env var — server listen port (default `8080`)
