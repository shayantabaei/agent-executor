# agent-executor

A local docker backed code execution service for AI agents.

`agent-executor` exposes an API that allows agents, developer tools or local workflows to execute short code snippets inside Docker containers with execution controls.

## Features

- HTTP API with `/health`, `/executions` and `/runtimes` endpoints
- Docker backed execution service
- Python and JavaScript runtime support
- Timeout handling
- Output limiting
- Request body size and code payload size validation
- Configurable Docker memory and CPU limits
- Tests

## Project Status

Currently under active development.

Current focus areas:

- Docker execution hardening
- File input support
- Temporary execution workspace
- Artifact collection
- Workspace cleanup

## Why this exists

There is an increasing need for AI agents to execute code as part of a workflow in a safe manner. This project explores what a small, local first execution service might look like for that case.

The goal is not to build a production grade sandbox, but rather to build a testable, extensible execution service with practical guardrails.

## Security Model

`agent-executor` executes arbitrary code inside Docker containers.
Intended for local development workflows.

The service applies practical controls such as request validation, execution timeouts, output limits, and Docker resource limits, but arbitrary code execution remains inherently risky.

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

Example payload:

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
### Supported Runtimes

| Language | Status |
|---|---|
| Python | Supported |
| JavaScript | Supported |

### Development
Run tests: `go test ./...`

Run the service: `go run ./cmd/server`
