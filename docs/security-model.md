# Security Model

`agent-executor` executes arbitrary user-provided code inside Docker containers.

This project is designed for local and development workflows. It is not intended to be a hardened production sandbox or a multi-tenant code execution platform.

Docker provides useful isolation, but Docker should not be treated as a complete security boundary against malicious code.

## Goals

The service aims to provide practical guardrails for local agent execution workflows.

Current and planned controls include:

- Limit execution time
- Limit captured output size
- Limit request body size
- Limit code payload size
- Apply Docker memory limits
- Apply Docker CPU limits
- Avoid exposing host files by default
- Validate file paths before creating workspaces
- Use temporary workspaces for file input and artifact output
- Clean up workspaces after execution

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

## Threats considered

### Infinite loops

User code may run forever.

Mitigation:

- Execution timeouts

### Excessive output

User code may print an unbounded amount of output.

Mitigation:

- Limited stdout/stderr capture

### Large request bodies

Clients may send very large HTTP request bodies.

Mitigation:

- Request body size limits

### Large code payloads

Clients may send code payloads larger than the service is intended to handle.

Mitigation:

- Code size validation

### Excessive memory usage

User code may allocate too much memory.

Mitigation:

- Docker memory limits

### Excessive CPU usage

User code may consume too much CPU.

Mitigation:

- Docker CPU limits

### Path traversal

File support may introduce risks where users attempt to write outside the intended workspace.

Examples:

```text
../secret.txt
../../etc/passwd
/tmp/unsafe.txt
```

Planned mitigation:

- Reject absolute paths
- Reject path traversal
- Normalize and validate paths
- Ensure all files remain inside the temporary workspace

### Host file exposure

Mounting host directories into containers can expose sensitive files.

Mitigation:

- Avoid mounting sensitive host paths
- Use temporary workspaces
- Keep workspace scope narrow
- Clean up workspace after execution

## Docker execution hardening

The Docker execution command applies a small set of practical hardening controls by default:

- Network access is disabled with `--network none`
- Memory usage is limited with `--memory`
- CPU usage is limited with `--cpus`
- Process creation is limited with `--pids-limit`
- Privilege escalation is restricted with `--security-opt no-new-privileges:true`
- Runtime images are not pulled implicitly when pull policy is set to `never`

These controls reduce common local execution risks, but they do not make arbitrary code execution safe.

## Risks that remain

Executing arbitrary code remains inherently risky.

Remaining risks include:

- Container escape vulnerabilities
- Bugs in Docker configuration
- Bugs in path validation
- Bugs in workspace cleanup
- Network abuse if networking is enabled
- Disk usage abuse if workspace limits are incomplete
- Malicious code targeting the Docker daemon or host environment
- Resource exhaustion outside configured limits

### Read-only filesystem mode

Read-only container filesystems are not enabled yet.

A read-only filesystem is a useful hardening control, but it requires a clear writable workspace model. Many runtimes expect writable temporary directories, caches, or working directories. Enabling `--read-only` before defining those writable locations could break normal Python or JavaScript execution.

This will be revisited with temporary workspace support, where the project can define:

- where input files are written
- where code is executed
- which temporary directories are writable
- where generated artifacts are collected
- how workspaces are cleaned up

## Intended usage

Recommended usage:

- Local development
- AI agent experimentation
- Controlled demos
- Personal portfolio/testing environments
- Internal tools where inputs are trusted or semi-trusted

Not recommended usage:

- Public unauthenticated APIs
- Multi-user hosted execution
- Production execution of untrusted code
- Sensitive environments without additional sandboxing
- Running arbitrary code from unknown users

## Security posture

The project takes a practical defense-in-depth approach for a local execution service.

It does not claim to make arbitrary code execution safe. Instead, it documents the risks clearly and implements incremental controls to reduce common failure modes.
