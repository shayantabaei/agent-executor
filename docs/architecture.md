# Architecture

`agent-executor` is organized around a small HTTP API, a runtime abstraction, and a Docker-backed execution service.

The service accepts code execution requests, validates them, selects the requested runtime, executes the code inside Docker, captures the result, and returns a normalized JSON response.

## Core Flow

1. Client sends an execution request.
2. API layer validates request body size and code size.
3. Execution service selects the requested runtime.
4. Runtime defines how the code should be executed.
5. Docker execution layer starts a container with configured limits.
6. Code runs inside the container.
7. Output is captured through a limited writer.
8. Execution result is returned as JSON.

## Components

### API Layer

The API layer is responsible for:

- HTTP routing
- Request decoding
- Request body size validation
- Code size validation
- Response formatting
- Mapping errors to HTTP status codes

Current endpoints:

- `GET /health`
- `GET /runtimes`
- `POST /executions`

### Execution Service

The execution service coordinates execution requests.

Responsibilities:

- Validate requested language/runtime
- Select the correct runtime implementation
- Apply execution timeout behavior
- Call the Docker execution layer
- Return a normalized execution result

### Runtime Abstraction

Runtimes define how a language should be executed.

Each runtime is responsible for describing the command or container execution behavior needed for a specific language.

Current runtimes:

- Python
- JavaScript

The runtime abstraction makes it possible to add more languages without changing the API contract.

### Docker Execution Layer

The Docker execution layer is responsible for running code inside containers.

Responsibilities:

- Start a container using the selected runtime image
- Apply memory limits
- Apply CPU limits
- Execute the runtime command
- Capture stdout and stderr
- Stop execution when timeout is reached
- Return the exit code and output

### Limited Writer

The limited writer prevents unbounded output from consuming too much memory.

It caps how much stdout and stderr can be captured from a running execution.

This protects the service from code that prints excessively or produces unexpectedly large output.

## Current Execution Model

The current execution model is intentionally simple:

```text
HTTP request
  -> API validation
  -> Execution service
  -> Runtime selection
  -> Docker execution
  -> Limited output capture
  -> JSON response
```

## Planned Workspace Model

File support will introduce a temporary workspace lifecycle:

```text
HTTP request with files
  -> Validate file metadata and paths
  -> Create temporary workspace
  -> Write input files into workspace
  -> Execute code inside workspace
  -> Collect generated artifacts
  -> Clean up workspace
  -> Return execution result and artifact metadata
```

## Design Principles

### Keep the API small

The API should stay focused on execution workflows.

### Keep runtimes isolated

Language-specific behavior should live in runtime implementations rather than spreading across the API or Docker execution code.

### Make limits explicit

Timeouts, output limits, request limits, memory limits, and CPU limits should be configurable and visible in the code.

### Prefer simple security controls first

This project is not a hardened sandbox, but it should still apply practical controls that reduce obvious risks.

### Test behavior at boundaries

Important behavior should be covered by tests, especially:

- Request validation
- Runtime selection
- Timeout handling
- Output limiting
- Docker execution behavior
- File path validation
- Workspace cleanup
