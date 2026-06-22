# Roadmap

This roadmap tracks the major implementation areas for `agent-executor`.

## Completed

- HTTP API
- `GET /health`
- `GET /runtimes`
- `POST /executions`
- Python runtime
- JavaScript runtime
- Runtime abstraction
- Docker-backed execution
- Timeout handling
- Output limiting
- Request body validation
- Code size validation
- Docker configuration defaults
- Tests for API handlers
- Tests for runtimes
- Tests for limited writer
- Tests for Docker execution

## Current Focus

### 1. Docker hardening

Improve the Docker execution configuration with practical local safety controls.

Potential work:

- Disable unnecessary privileges
- Add `no-new-privileges`
- Apply memory limits
- Apply CPU limits
- Consider disabling network access by default
- Consider read-only filesystem mode
- Define writable temp locations
- Avoid mounting sensitive host paths
- Document tradeoffs for each Docker option

### 2. File input support

Allow execution requests to include files that are made available to the running code.

Potential work:

- Define file request model
- Add file count limits
- Add per-file size limits
- Add total file size limits
- Validate file paths
- Reject absolute paths
- Reject path traversal
- Preserve nested relative directories safely

### 3. Temporary workspace lifecycle

Create an isolated temporary workspace for each execution.

Potential work:

- Create workspace before execution
- Write input files into workspace
- Execute code inside workspace
- Ensure workspace path is not user-controlled
- Clean up workspace after execution
- Test cleanup on success
- Test cleanup on failure
- Test cleanup on timeout

### 4. Artifact collection

Allow executed code to produce files that can be returned as execution artifacts.

Potential work:

- Define artifact response model
- Collect files from workspace after execution
- Limit artifact count
- Limit individual artifact size
- Limit total artifact size
- Return artifact metadata
- Optionally return artifact content for small text files
- Avoid exposing unsafe host paths

## Future Ideas

### Additional runtimes

Potential runtimes:

- Go
- Ruby
- Bash
- TypeScript

### Configuration improvements

Potential work:

- Environment-based config
- Config validation
- Per-runtime config
- Per-request execution limits with maximum caps

### Observability

Potential work:

- Structured logging
- Execution duration metrics
- Timeout metrics
- Error classification
- Request IDs

### API improvements

Potential work:

- Better error response model
- Runtime metadata endpoint
- Execution status model
- More detailed validation errors

### Developer experience

Potential work:

- Docker Compose setup
- Example requests
- Example agent integration
- Makefile
- GitHub Actions CI
- Contribution guide
