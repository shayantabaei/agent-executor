# Roadmap

This roadmap tracks the major implementation areas for `agent-executor`.

`agent-executor` is being built iteratively as a local Docker-backed execution service for AI agents. The project focuses on a small HTTP API, MCP-based agent integration, runtime extensibility, practical execution controls, temporary workspaces, artifact collection, structured execution results, and clear documentation around security tradeoffs.

## Completed

### Core API

- HTTP API
- `GET /health`
- `GET /runtimes`
- `POST /executions`

### MCP server

- MCP stdio server
- `execute_code` MCP tool
- `agent-executor://runtimes` MCP resource
- `agent-executor://capabilities` MCP resource
- MCP execution path using the shared execution service
- MCP resource listing and reading tests
- Stderr-only logging for stdio compatibility

### Runtime support

- Python runtime
- JavaScript runtime
- Runtime abstraction

### Execution controls

- Timeout handling
- Output limiting
- Request body validation
- Code size validation
- File count validation
- Individual file size validation
- Total file size validation
- Normalized execution result model
- Shared execution validation through the execution service

### Docker execution

- Docker-backed execution
- Docker configuration defaults
- Docker hardening configuration
- Disabled Docker networking by default
- Docker memory limits
- Docker CPU limits
- Docker PID limits
- Docker `no-new-privileges`
- Docker pull policy
- Docker run argument construction tests

### File input support

- File input request model
- File input validation
- Maximum file count
- Maximum individual file size
- Maximum total file size
- Safe relative file path validation
- Temporary execution workspace creation
- Input file writing
- Docker workspace mounting
- Running code from `/workspace`
- Workspace cleanup after execution

### Artifact collection

- Artifact response model
- Generated artifact collection
- Input file exclusion from artifacts
- Artifact count limits
- Individual artifact size limits
- Total artifact size limits
- Inline content size limits
- Inline content for small UTF-8 text artifacts
- Metadata-only handling for larger or binary artifacts

### Tests

- Tests for API handlers
- Tests for runtimes
- Tests for limited writer
- Tests for Docker execution
- Tests for validation
- Tests for execution service validation
- Tests for execution service timeout behavior
- Tests for workspace behavior
- Tests for artifact collection
- Tests for artifact limits
- Tests for MCP resource listing
- Tests for MCP resource reading

### Documentation

- README
- Architecture documentation
- Security model documentation
- Roadmap documentation
- MCP server documentation

## Current Focus

### 1. Binary/base64 artifact support

Support returning small binary artifacts inline using base64 encoding.

Potential work:

- Detect binary artifact content types
- Base64 encode small binary artifacts
- Apply inline binary size limits
- Return `encoding: "base64"` for binary content
- Return appropriate `contentType`
- Add tests for PNG, PDF, and generic binary artifact behavior
- Document binary artifact response examples

### 2. Workspace cleanup hardening

Improve confidence that temporary workspaces are cleaned up correctly.

Potential work:

- Ensure cleanup runs after successful execution
- Ensure cleanup runs after failed execution
- Ensure cleanup runs after timeout
- Add tests for cleanup behavior around executor errors
- Add logging for cleanup failures
- Consider making cleanup errors observable without hiding execution results

### 3. MCP improvements

Continue improving the MCP integration now that the initial stdio server is in place.

Potential work:

- Add MCP tool execution tests
- Add tests around MCP request/response mapping
- Document local MCP client configuration examples
- Document usage with MCP-compatible clients
- Add richer runtime/capability metadata
- Add clearer MCP error mapping
- Consider whether MCP responses should include structured content in addition to text content
- Verify artifact response behavior through MCP clients

### 4. Additional Docker hardening

Continue improving Docker execution controls without breaking normal language runtime behavior.

Potential work:

- Consider read-only filesystem mode
- Define writable temp locations if read-only mode is enabled
- Consider running containers as a non-root user
- Consider dropping Linux capabilities
- Consider limiting writable directories
- Document tradeoffs for each hardening option

### 5. Observability and error reporting

Improve debugging and operational visibility.

Potential work:

- Structured logging
- Execution duration reporting
- Timeout metrics
- Error classification
- Request IDs
- Workspace cleanup metrics
- Clearer Docker execution errors
- Clearer validation errors
- MCP-specific error formatting

## Future Ideas

### Additional runtimes

Potential runtimes:

- Go
- Rust
- Ruby
- Bash
- TypeScript
- Java

Each runtime should be added through the runtime abstraction without changing the public execution API.

### Rust runtime support

Add Rust as a supported runtime so agents can execute small Rust programs through the same execution API.

This is especially useful for future agent workflows that need fast local computation, simulations, ranking logic, or structured analysis.

Potential work:

- Add a Rust runtime implementation
- Define how Rust snippets are compiled and executed
- Decide whether Rust execution uses `rustc`, `cargo script`, or a temporary Cargo project
- Add timeout behavior for both compile time and run time
- Capture compiler errors in stderr
- Add tests for successful Rust execution
- Add tests for Rust compile failures
- Add tests for timeout behavior during Rust execution
- Document Rust runtime examples

### Structured output support

Improve support for agent workflows that expect machine-readable execution results.

Potential work:

- Document conventions for writing JSON results to stdout
- Document conventions for writing JSON artifacts
- Add examples for returning structured recommendations
- Consider optional JSON output validation
- Consider a reserved artifact path such as `result.json`
- Consider response helpers for common structured outputs

### Agent workflow examples

Add examples showing how agents can use `agent-executor` as a local compute tool.

Potential work:

- Example agent request that runs a small data analysis script
- Example agent request that passes JSON or CSV input files
- Example agent request that returns generated artifacts
- Example workflow that produces a ranked recommendation report
- Example workflow that validates structured output
- Documentation showing how an agent can call `/executions`
- Documentation showing how an agent can use the MCP `execute_code` tool
- Documentation showing how artifacts can be used as generated reports

### Configuration improvements

Potential work:

- Environment-based config
- Config validation
- Per-runtime config
- Per-request execution limits with maximum caps
- Workspace size limits
- Artifact size limits
- Runtime image configuration

### API improvements

Potential work:

- Better error response model
- Runtime metadata endpoint
- Execution status model
- More detailed validation errors
- More consistent timeout response fields
- Execution duration in responses
- Runtime metadata in execution responses
- Artifact response model refinements

### MCP API improvements

Potential work:

- Runtime metadata in `agent-executor://runtimes`
- Execution limit metadata in `agent-executor://capabilities`
- More detailed tool descriptions
- More precise input schema descriptions
- Example prompts or usage notes for MCP clients
- Optional MCP resource for security assumptions
- Optional MCP resource for execution limits

### Developer experience

Potential work:

- Docker Compose setup
- Example requests
- Example agent integration
- Example MCP client configuration
- Makefile
- GitHub Actions CI
- `govulncheck` in CI
- Contribution guide

### Streaming output

Support streaming stdout and stderr while code is running.

Potential work:

- Stream stdout and stderr incrementally
- Preserve output limits
- Return final execution metadata after completion
- Consider Server-Sent Events or another simple streaming mechanism
- Consider MCP-compatible streaming/progress behavior

### Persistent execution history

Optional storage for past execution metadata.

Potential work:

- Store execution status
- Store runtime/language
- Store timing metadata
- Store stdout/stderr metadata
- Store artifact metadata
- Avoid storing sensitive code or file contents by default

### Authentication and API keys

Add optional authentication for non-local usage.

Potential work:

- API key support
- Local-only default mode
- Documentation for safe deployment assumptions
- Rate limiting if exposed beyond localhost

## Long-Term Ideas

- Queue-based execution model
- Async execution API
- Web UI for submitting executions
- Execution templates
- Agent evaluation workflows
- Fantasy football agent integration using the executor as a compute backend
- Simulation-heavy agent workflows
- Deterministic execution options for repeatable simulations
- Sandboxed package installation strategy
- More advanced resource accounting
- Better isolation through stronger sandbox technologies

## Non-Goals For Now

The project is not currently trying to become:

- A production-grade sandbox
- A public code execution platform
- A multi-tenant hosted service
- A replacement for dedicated secure sandboxing systems
- A general-purpose CI runner
- A domain-specific fantasy football platform

The current goal is to build a clear, local-first execution service that demonstrates practical API design, MCP integration, Docker execution, runtime abstraction, workspace handling, artifact collection, structured agent workflows, and security tradeoff documentation.
