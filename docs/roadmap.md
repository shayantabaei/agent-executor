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
- Docker hardening configuration
- Disabled Docker networking by default
- Docker memory limits
- Docker CPU limits
- Docker PID limits
- Docker `no-new-privileges`
- Docker pull policy
- File input request model
- File input validation
- Safe relative file path validation
- Temporary execution workspace creation
- Input file writing
- Docker workspace mounting
- Running code from `/workspace`
- Workspace cleanup after execution
- Tests for API handlers
- Tests for runtimes
- Tests for limited writer
- Tests for Docker execution
- Tests for validation
- Tests for workspace behavior
- Artifact response model
- Generated artifact collection
- Input file exclusion from artifacts
- Artifact count limits
- Individual artifact size limits
- Total artifact size limits
- Inline content for small UTF-8 text artifacts
- Metadata-only handling for larger or binary artifacts

## Current Focus

### 1. Binary/base64 artifact support

Support returning small binary artifacts inline using base64 encoding.

Potential work:

- Detect binary artifact content types
- Base64 encode small binary artifacts
- Apply inline binary size limits
- Return `encoding: "base64"` for binary content
- Return appropriate `contentType`
- Add tests for PNG/PDF/binary artifact behavior

### 2. Workspace cleanup hardening

Improve confidence that temporary workspaces are cleaned up correctly.

Potential work:

- Ensure cleanup runs after successful execution
- Ensure cleanup runs after failed execution
- Ensure cleanup runs after timeout
- Add tests for cleanup behavior around executor errors
- Add logging for cleanup failures
- Consider making cleanup errors observable without hiding execution results

### 3. Additional Docker hardening

Continue improving Docker execution controls without breaking normal language runtime behavior.

Potential work:

- Consider read-only filesystem mode
- Define writable temp locations if read-only mode is enabled
- Consider running containers as a non-root user
- Consider dropping Linux capabilities
- Consider limiting writable directories
- Document tradeoffs for each hardening option

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
- Workspace size limits
- Artifact size limits

### Observability

Potential work:

- Structured logging
- Execution duration metrics
- Timeout metrics
- Error classification
- Request IDs
- Workspace cleanup metrics

### API improvements

Potential work:

- Better error response model
- Runtime metadata endpoint
- Execution status model
- More detailed validation errors
- Artifact response model

### Developer experience

Potential work:

- Docker Compose setup
- Example requests
- Example agent integration
- Makefile
- GitHub Actions CI
- Contribution guide
