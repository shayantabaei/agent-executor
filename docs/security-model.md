# Security Model

`agent-executor` executes arbitrary user-provided code inside Docker containers.

This project is designed for local and development workflows. It is not intended to be a hardened production sandbox or a multi-tenant code execution platform.

Docker provides useful isolation, but Docker should not be treated as a complete security boundary against malicious code.

`agent-executor` can be used through either the HTTP API or the MCP stdio server. Both entrypoints share the same execution service so execution-specific validation, timeout handling, workspace behavior, Docker execution, and artifact collection are centralized.

## Goals

The service aims to provide practical guardrails for local agent execution workflows.

Current controls include:

- Limit execution time
- Limit captured output size
- Limit HTTP request body size
- Limit code payload size
- Limit input file count
- Limit individual input file size
- Limit total input file size
- Limit artifact count
- Limit individual artifact size
- Limit total artifact size
- Limit inline artifact content size
- Apply Docker memory limits
- Apply Docker CPU limits
- Disable Docker networking by default
- Apply Docker PID limits
- Restrict privilege escalation with `no-new-privileges`
- Avoid exposing host files by default
- Validate file paths before creating workspaces
- Use temporary workspaces for file input and artifact output
- Clean up workspaces after execution
- Keep MCP stdio protocol output separate from logs by writing logs to stderr

## Non-goals

`agent-executor` does not currently aim to provide:

- Strong isolation against malicious users
- Production-grade arbitrary code execution security
- Multi-tenant sandboxing
- Kernel-level sandboxing
- Protection against container escape vulnerabilities
- Abuse prevention for public unauthenticated APIs
- Guaranteed network isolation beyond configured Docker behavior
- Protection from every possible denial-of-service vector
- Long-term storage of execution files or artifacts
- Secure remote MCP server hosting

## Entrypoints

### HTTP API

The HTTP API exposes code execution over local HTTP endpoints.

HTTP-specific protections include:

- Request body size limits
- Request decoding
- HTTP error mapping

Execution-specific protections, such as code size validation, file validation, timeout behavior, Docker execution limits, and artifact limits, are handled by the shared execution service.

### MCP stdio server

The MCP stdio server exposes `agent-executor` to MCP-compatible clients.

The MCP server exposes:

- Tool: `execute_code`
- Resource: `agent-executor://runtimes`
- Resource: `agent-executor://capabilities`

The MCP server does not call the HTTP API. The `execute_code` tool calls the shared execution service directly.

This means MCP execution requests use the same execution-specific guardrails as HTTP execution requests.

MCP stdio uses stdout for JSON-RPC protocol messages. The MCP server must not write logs, debug output, or human-readable status messages to stdout. Logs should be written to stderr.

## Threats considered

### Infinite loops

User code may run forever.

Mitigation:

- Execution timeouts

### Excessive output

User code may print an unbounded amount of output.

Mitigation:

- Limited stdout/stderr capture

### Large HTTP request bodies

Clients may send very large HTTP request bodies.

Mitigation:

- HTTP request body size limits

### Large MCP tool inputs

MCP clients may send large tool inputs.

Mitigation:

- Code size validation in the execution service
- Input file count validation in the execution service
- Individual input file size validation in the execution service
- Total input file size validation in the execution service

### Large code payloads

Clients may send code payloads larger than the service is intended to handle.

Mitigation:

- Code size validation

### Excessive input files

Clients may send too many input files or files that are too large.

Mitigation:

- Maximum input file count
- Maximum individual input file size
- Maximum total input file size

### Excessive memory usage

User code may allocate too much memory.

Mitigation:

- Docker memory limits

### Excessive CPU usage

User code may consume too much CPU.

Mitigation:

- Docker CPU limits

### Excessive process creation

User code may attempt to create many processes.

Mitigation:

- Docker PID limits

### Network abuse

User code may attempt to access external networks.

Mitigation:

- Docker networking is disabled by default with `--network none`

### Path traversal

File support introduces risks where users attempt to write outside the intended workspace.

Examples:

```text
../secret.txt
../../etc/passwd
/tmp/unsafe.txt
safe/../../secret.txt
```

Mitigation:

- Reject absolute paths
- Reject path traversal
- Reject empty paths
- Reject backslash-based paths
- Normalize and validate paths
- Ensure all files remain inside the temporary workspace

### Host file exposure

Mounting host directories into containers can expose sensitive files.

Mitigation:

- Avoid mounting sensitive host paths
- Use temporary execution workspaces
- Keep workspace scope narrow
- Mount only the temporary workspace into the container
- Clean up workspace after execution

### Generated artifact abuse

Executed code may create many files, large files, or files that are unsafe to inline in the API or MCP response.

Mitigation:

- Artifact count is limited
- Individual artifact size is limited
- Total artifact size is limited
- Inline artifact content size is limited
- Original input files are excluded from generated artifacts
- Binary files are returned as metadata-only artifacts for now

### MCP protocol corruption

MCP stdio servers use stdout for JSON-RPC protocol messages.

Writing logs, debug messages, or startup messages to stdout can corrupt the MCP protocol stream.

Mitigation:

- Do not write logs to stdout
- Write logs to stderr
- Avoid printing startup banners or human-readable status messages from the MCP stdio command

## Docker execution hardening

The Docker execution command applies a small set of practical hardening controls by default:

- Network access is disabled with `--network none`
- Memory usage is limited with `--memory`
- CPU usage is limited with `--cpus`
- Process creation is limited with `--pids-limit`
- Privilege escalation is restricted with `--security-opt no-new-privileges:true`
- Runtime images are not pulled implicitly when pull policy is set to `never`
- The temporary workspace is mounted into the container at `/workspace`
- The container runs with `/workspace` as its working directory

These controls reduce common local execution risks, but they do not make arbitrary code execution safe.

## Temporary workspace exposure

Input files are written to a temporary host workspace and mounted into the container at `/workspace`.

Mitigation:

- File paths are validated before workspace creation
- Absolute paths are rejected
- Path traversal is rejected
- Empty paths are rejected
- Backslash-based paths are rejected to keep API paths portable
- Each execution receives a separate temporary workspace
- Workspaces are cleaned up after execution

The workspace is intentionally short-lived and scoped to a single execution request.

## Artifact response limits

Generated files can increase response size or expose data produced during execution.

Mitigation:

- Artifact count is limited
- Individual artifact size is limited
- Total artifact size is limited
- Inline content size is limited
- Original input files are excluded from generated artifacts
- Binary files are returned as metadata-only artifacts for now

Small UTF-8 text artifacts may be returned inline in the API or MCP response.

Larger files and binary files are returned as metadata-only artifacts for now. Future binary artifact support may add base64-encoded content for small binary files.

## Risks that remain

Executing arbitrary code remains inherently risky.

Remaining risks include:

- Container escape vulnerabilities
- Bugs in Docker configuration
- Bugs in path validation
- Bugs in artifact collection
- Bugs in workspace cleanup
- Network abuse if networking is enabled
- Disk usage abuse if workspace limits are incomplete
- Malicious code targeting the Docker daemon or host environment
- Resource exhaustion outside configured limits
- Sensitive data exposure if unsafe host paths are mounted in the future
- Unsafe behavior if this service is exposed publicly without authentication
- Unsafe behavior if exposed as a remote MCP server without additional authentication, authorization, and sandboxing controls
- MCP client misuse if a client allows untrusted prompts to invoke the `execute_code` tool without review

## Read-only filesystem mode

Read-only container filesystems are not enabled yet.

A read-only filesystem is a useful hardening control, but it requires a clear writable workspace model. Many runtimes expect writable temporary directories, caches, or working directories. Enabling `--read-only` before defining those writable locations could break normal Python or JavaScript execution.

This will be revisited now that temporary workspace support exists.

Future work may define:

- Where input files are written
- Where code is executed
- Which temporary directories are writable
- Where generated artifacts are collected
- How workspaces are cleaned up
- Whether runtime-specific writable paths are needed
- Whether `--read-only` can be enabled safely by default

## Intended usage

Recommended usage:

- Local development
- AI agent experimentation
- Controlled demos
- Personal portfolio/testing environments
- Internal tools where inputs are trusted or semi-trusted
- Local MCP-compatible agent workflows

Not recommended usage:

- Public unauthenticated APIs
- Multi-user hosted execution
- Production execution of untrusted code
- Sensitive environments without additional sandboxing
- Running arbitrary code from unknown users
- Remote MCP server exposure without additional controls

## Security posture

The project takes a practical defense-in-depth approach for a local execution service.

It does not claim to make arbitrary code execution safe. Instead, it documents the risks clearly and implements incremental controls to reduce common failure modes.
