# agent-executor

A local Docker-backed code execution service for AI agents.

`agent-executor` exposes an API that allows agents, developer tools, or local workflows to execute short code snippets inside Docker containers with practical execution controls.

## Features

- HTTP API with `/health`, `/executions`, and `/runtimes` endpoints
- Docker-backed execution service
- Python and JavaScript runtime support
- Runtime abstraction for adding new languages
- Timeout handling
- Output limiting
- Request body size and code payload size validation
- Input file validation
- Temporary execution workspaces
- Docker workspace mounting
- Configurable Docker memory and CPU limits
- Docker hardening controls such as disabled networking, PID limits, and `no-new-privileges`
- Tests for API handlers, runtimes, Docker execution, validation, and workspace behavior

## Project Status

Currently under active development.

Current focus areas:

- Artifact collection
- Workspace cleanup hardening
- Additional Docker execution hardening
- Additional runtime support
- Improved observability and error reporting

## Why this exists

There is an increasing need for AI agents to execute code as part of a workflow in a safe manner. This project explores what a small, local-first execution service might look like for that use case.

The goal is not to build a production-grade sandbox, but rather to build a testable, extensible execution service with practical guardrails and documented tradeoffs.

## Security Model

`agent-executor` executes arbitrary user-provided code inside Docker containers.

It is intended for local development workflows. It is not a hardened production sandbox or a multi-tenant code execution platform.

The service applies practical controls such as request validation, execution timeouts, output limits, Docker resource limits, disabled networking, PID limits, and temporary workspaces, but arbitrary code execution remains inherently risky.

See [docs/security-model.md](docs/security-model.md).

## Architecture

See [docs/architecture.md](docs/architecture.md).

## Roadmap

See [docs/roadmap.md](docs/roadmap.md).

## API Overview

### Health check

```http
GET /health
```

Returns service health.

### List runtimes

```http
GET /runtimes
```

Returns the available execution runtimes.

### Execute code

```http
POST /executions
```

Executes code in one of the supported runtimes.

Example request:

```json
{
  "language": "python",
  "code": "print('hello from python')"
}
```

Example response:

```json
{
  "stdout": "hello from python\n",
  "stderr": "",
  "exitCode": 0
}
```

### Execute code with input files

Execution requests may include input files.

Input files are validated by the API layer, written to a temporary host workspace, and mounted into the Docker container at `/workspace`. Code runs from `/workspace`, so relative paths such as `data/input.txt` work inside the execution.

Example request:

```json
{
  "language": "python",
  "code": "print(open(\"data/input.txt\").read())",
  "files": [
    {
      "path": "data/input.txt",
      "content": "hello from workspace"
    }
  ]
}
```

Example response:

```json
{
  "stdout": "hello from workspace\n",
  "stderr": "",
  "exitCode": 0
}
```

Input file paths must be safe relative paths. Absolute paths, path traversal, empty paths, and backslash-based paths are rejected.

Supported:

```text
data/input.txt
src/main.py
fixtures/sample.json
```

Rejected:

```text
/etc/passwd
../secret.txt
safe/../../secret.txt
data\input.txt
```

## Supported Runtimes

| Language | Status |
|---|---|
| Python | Supported |
| JavaScript | Supported |

## Development

Run tests:

```bash
go test ./...
```

Run the service:

```bash
go run ./cmd/server
```
