# Architecture

`agent-executor` is organized around a small HTTP API, an MCP stdio server, a runtime abstraction, a temporary workspace lifecycle, artifact collection, and a Docker-backed execution service.

The service accepts code execution requests through HTTP or MCP, validates them, writes optional input files into a temporary workspace, executes code inside Docker, captures output, collects generated artifacts, cleans up the workspace, and returns a normalized result.

## Core Flow

1. Client sends an execution request through HTTP or MCP.
2. The entrypoint layer decodes the request and applies entrypoint-specific validation such as HTTP request body size limits.
3. Execution service validates execution limits and selects the requested runtime.
4. Runtime defines how the code should be executed.
5. Temporary workspace is created.
6. Input files are written into the workspace.
7. Docker execution layer starts a container with configured limits.
8. Workspace is mounted into the container at `/workspace`.
9. Code runs inside the container with `/workspace` as the working directory.
10. Output is captured through limited writers.
11. Generated artifacts are collected from the workspace.
12. Workspace is cleaned up.
13. Execution result is returned to the caller.

## Components

### API Layer

The API layer is responsible for:

- HTTP routing
- Request decoding
- Request body size validation
- Response formatting
- Mapping errors to HTTP status codes

Current endpoints:

- `GET /health`
- `GET /runtimes`
- `POST /executions`

### MCP Server

The MCP server exposes agent-executor as a local Model Context Protocol stdio server for MCP-compatible clients.

The MCP server lives under:

```text
cmd/agent-executor-mcp
internal/mcpserver
```

The MCP server exposes:

- Tool: `execute_code`
- Resource: `agent-executor://runtimes`
- Resource: `agent-executor://capabilities`

MCP requests do not call the HTTP API. The MCP tool calls the shared execution service directly:

```text
MCP client
  -> MCP stdio server
  -> Execution service
  -> Runtime selection
  -> Temporary workspace creation
  -> Docker execution with workspace mounted at /workspace
  -> Limited output capture
  -> Artifact collection
  -> Workspace cleanup
  -> MCP tool result
```

This keeps execution validation, timeout handling, runtime lookup, file input handling, workspace creation, Docker execution, and artifact collection centralized in the execution package.

Because MCP stdio uses stdout for JSON-RPC protocol messages, the MCP server must not write logs or debug output to stdout. Logs should be written to stderr.

### Request Validation

Execution requests are validated before they reach the Docker execution layer.

HTTP-specific request validation, such as request body size validation, happens in the API layer. Execution-specific validation happens in the execution service so both HTTP and MCP requests share the same execution limits.

Validation includes:

- Required language/runtime
- Required code
- Maximum code size
- Maximum file count
- Maximum individual file size
- Maximum total file size
- Safe relative file paths

Input file paths must use portable slash-based relative paths such as:

```text
data/input.txt
src/main.py
fixtures/sample.json
```

The API rejects unsafe paths such as:

```text
/etc/passwd
../secret.txt
safe/../../secret.txt
data\input.txt
```

### Execution Service

The execution service coordinates execution requests.

Responsibilities:

- Validate requested language/runtime
- Select the correct runtime implementation
- Create a temporary workspace
- Write input files into the workspace
- Apply execution timeout behavior
- Call the Docker execution layer
- Collect generated artifacts
- Clean up the workspace after execution
- Return a normalized execution result

### Runtime Abstraction

Runtimes define how a language should be executed.

Each runtime is responsible for describing the command and execution behavior needed for a specific language.

Current runtimes:

- Python
- JavaScript

The runtime abstraction makes it possible to add more languages without changing the API contract.

### Workspace Lifecycle

Each execution creates a temporary host workspace.

Input files are written into this workspace before Docker execution starts. Safe nested paths are preserved.

Example input file:

```json
{
  "path": "data/input.txt",
  "content": "hello from workspace"
}
```

The service writes that file into a temporary host directory:

```text
<temp-workspace>/data/input.txt
```

The temporary workspace is then mounted into the Docker container at:

```text
/workspace
```

The container runs with `/workspace` as its working directory, so user code can access files with relative paths:

```python
open("data/input.txt").read()
```

After execution completes, fails, or times out, generated artifacts are collected and the workspace is cleaned up.

### Docker Execution Layer

The Docker execution layer is responsible for running code inside containers.

Responsibilities:

- Build Docker run arguments
- Start a container using the selected runtime image
- Apply memory limits
- Apply CPU limits
- Disable networking by default
- Apply PID limits
- Apply `no-new-privileges`
- Mount the temporary workspace
- Set `/workspace` as the working directory
- Execute the runtime command
- Capture stdout and stderr
- Stop execution when timeout is reached
- Return the exit code and output

### Limited Writers

Limited writers prevent unbounded output from consuming too much memory.

They cap how much stdout and stderr can be captured from a running execution.

This protects the service from code that prints excessively or produces unexpectedly large output.

### Artifact Collection

After execution completes, the workspace is scanned for generated files.

Artifacts represent files created by the executed code. Original input files are excluded from artifact results.

Each artifact includes:

- `path`
- `size`
- optional inline `content`
- optional `encoding`
- optional `contentType`

Small UTF-8 text artifacts may be returned inline in the API response.

Example artifact:

```json
{
  "path": "output.txt",
  "size": 14,
  "content": "hello artifact",
  "encoding": "utf-8",
  "contentType": "text/plain; charset=utf-8"
}
```

Larger files and binary files are returned as metadata-only artifacts for now.

Artifact collection applies configurable limits:

- Maximum artifact count
- Maximum individual artifact size
- Maximum total artifact size
- Maximum inline artifact content size

Future binary artifact support may add base64-encoded content for small binary files.

## Current Execution Model

The HTTP execution model is:

```text
HTTP request
  -> API request decoding/body size validation
  -> Execution service
  -> Request validation
  -> Runtime selection
  -> Temporary workspace creation
  -> Input file writing
  -> Docker execution with workspace mounted at /workspace
  -> Limited output capture
  -> Artifact collection
  -> Workspace cleanup
  -> JSON response
```

The MCP execution model is:

```text
MCP client
  -> MCP stdio server
  -> Execution service
  -> Request validation
  -> Runtime selection
  -> Temporary workspace creation
  -> Input file writing
  -> Docker execution with workspace mounted at /workspace
  -> Limited output capture
  -> Artifact collection
  -> Workspace cleanup
  -> MCP tool result
```

## Design Principles

### Keep the API small

The API should stay focused on execution workflows.

### Keep entrypoints thin

HTTP and MCP entrypoints should decode requests, call the shared execution service, and format responses for their respective protocols.

Execution behavior should stay centralized in the execution layer rather than being duplicated across entrypoints.

### Keep runtimes isolated

Language-specific behavior should live in runtime implementations rather than spreading across the API, MCP server, or Docker execution code.

### Validate before execution

Invalid requests should be rejected before Docker execution starts.

### Make limits explicit

Timeouts, output limits, request limits, file limits, memory limits, CPU limits, PID limits, and artifact limits should be configurable and visible in the code.

### Use temporary workspaces

Each execution should receive an isolated temporary workspace that can be mounted into the container and cleaned up after execution.

### Collect artifacts safely

Generated files should be collected only after execution and before workspace cleanup.

Artifact collection should enforce count and size limits before reading file content into memory.

Original input files should not be returned as generated artifacts.

### Prefer simple security controls first

This project is not a hardened sandbox, but it should still apply practical controls that reduce obvious risks.

### Test behavior at boundaries

Important behavior should be covered by tests, especially:

- Request validation
- File validation
- Runtime selection
- Timeout handling
- Output limiting
- Docker execution behavior
- Workspace creation
- Input file writing
- Artifact collection
- Artifact limits
- Workspace cleanup
- MCP resource listing and reading
- MCP tool execution behavior
