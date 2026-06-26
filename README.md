# agent-executor

A local Docker-backed code execution service for AI agents.

`agent-executor` exposes an HTTP API that allows agents, developer tools, and local workflows to execute short code snippets inside Docker containers with practical execution controls.

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
- Artifact collection for generated files
- Inline content for small UTF-8 text artifacts
- Artifact size and count limits
- Tests for API handlers, runtimes, Docker execution, validation, workspace behavior, and artifact collection

## Project Status

Currently under active development.

Current focus areas:

- Binary/base64 artifact support
- Workspace cleanup hardening
- Additional Docker execution hardening
- Additional runtime support
- Improved observability and error reporting

## Why this exists

There is an increasing need for AI agents to execute code as part of a workflow.

This project explores what a small, local-first execution service for that use case might look like.

The goal is not to build a production-grade sandbox. The goal is to build a testable, extensible execution service with practical guardrails and documented tradeoffs.

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

~~~http
GET /health
~~~

Returns service health.

### List runtimes

~~~http
GET /runtimes
~~~

Returns the available execution runtimes.

### Execute code

~~~http
POST /executions
~~~

Executes code in one of the supported runtimes.

Example request:

~~~json
{
  "language": "python",
  "code": "print('hello from python')"
}
~~~

Example response:

~~~json
{
  "stdout": "hello from python\n",
  "stderr": "",
  "exitCode": 0
}
~~~

### Execute code with input files

Execution requests may include input files.

Input files are validated by the API layer, written to a temporary host workspace, and mounted into the Docker container at `/workspace`.

Code runs from `/workspace`, so relative paths such as `data/input.txt` work inside the execution.

Example request:

~~~json
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
~~~

Example response:

~~~json
{
  "stdout": "hello from workspace\n",
  "stderr": "",
  "exitCode": 0
}
~~~

Input file paths must be safe relative paths.

Supported paths:

~~~text
data/input.txt
src/main.py
fixtures/sample.json
~~~

Rejected paths:

~~~text
/etc/passwd
../secret.txt
safe/../../secret.txt
data\input.txt
~~~

Absolute paths, parent directory traversal, empty paths, and backslash-based paths are rejected.

### Execute code with generated artifacts

Executed code can create files in the workspace. Generated files are returned as artifacts in the execution response.

Original input files are not returned as generated artifacts.

Example request:

~~~json
{
  "language": "python",
  "code": "open(\"output.txt\", \"w\").write(\"hello artifact\")"
}
~~~

Example response:

~~~json
{
  "stdout": "",
  "stderr": "",
  "exitCode": 0,
  "artifacts": [
    {
      "path": "output.txt",
      "size": 14,
      "content": "hello artifact",
      "encoding": "utf-8",
      "contentType": "text/plain; charset=utf-8"
    }
  ]
}
~~~

Small UTF-8 text artifacts may include inline `content`.

Larger files or binary files are returned as metadata-only artifacts for now.

### Execute code with input files and artifacts

Input files and generated artifacts can be used together.

Example request:

~~~json
{
  "language": "python",
  "code": "message = open(\"input/message.txt\").read()\nopen(\"output/result.txt\", \"w\").write(message.upper())",
  "files": [
    {
      "path": "input/message.txt",
      "content": "hello from input file"
    }
  ]
}
~~~

Example response:

~~~json
{
  "stdout": "",
  "stderr": "",
  "exitCode": 0,
  "artifacts": [
    {
      "path": "output/result.txt",
      "size": 21,
      "content": "HELLO FROM INPUT FILE",
      "encoding": "utf-8",
      "contentType": "text/plain; charset=utf-8"
    }
  ]
}
~~~

## Supported Runtimes

| Language | Status |
|---|---|
| Python | Supported |
| JavaScript | Supported |

## Development

Run tests:

~~~bash
go test ./...
~~~

Run the service:

~~~bash
go run ./cmd/server
~~~
