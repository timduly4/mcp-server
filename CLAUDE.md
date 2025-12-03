# CLAUDE.md — MCP Resource Server in Golang

## Goals
- Build an MCP resource server in Golang that exposes GitHub repositories starred by a user.
- Conform to the [MCP server spec](https://modelcontextprotocol.io/docs/develop/build-server).
- Provide modular, testable design using dependency injection.
- Use `mcp-go` (https://github.com/mark3labs/mcp-go) as the MCP server framework.
- Use Bazel as the build system for reproducible, hermetic builds.

## Constraints
- Language: Go 1.24+ (compatible with Bazel rules_go 0.50.1)
- Build System: Bazel 8.4+ with bzlmod (MODULE.bazel)
- Dependency injection: uber-go/fx
- MCP compliance: JSON-RPC over stdio, OAuth-secured
- External API: GitHub REST API v3
- Resources must remain simple: only starred repositories, no complex transformations.

## Architecture

### Build System
- **Bazel** - Modern, hermetic build system
  - bzlmod (MODULE.bazel) for dependency management
  - rules_go 0.50.1 for Go compilation
  - Gazelle 0.39.1 for auto-generating BUILD files
  - Go 1.24.0 SDK for compatibility
  - Custom build modes: fast, opt, debug, ci

### Agents / Modules
- **Server Module** (`internal/server`)
  - Initializes MCP server using `mcp-go`
  - Registers resource endpoints (static and templated)
  - Handles JSON-RPC requests over stdio transport
  - Returns resources in MCP-compliant format

- **GitHub Client Module** (`internal/github`)
  - Wraps GitHub REST API v3 calls
  - OAuth authentication with token-based auth
  - Query starred repos with automatic pagination
  - Returns normalized JSON objects with safe null handling

- **Resource Adapter Module** (`internal/resource`)
  - Maps GitHub data into MCP resource schema
  - Converts repo metadata to MCP resource format
  - Handles JSON serialization
  - Provides URI-based resource identification

- **Configuration Module** (`internal/config`)
  - Environment-based configuration
  - Token validation and defaults
  - Centralized settings management

- **DI Container** (uber-go/fx)
  - Wires server, client, and adapter using fx.Provide
  - Manages application lifecycle with fx.Lifecycle
  - Provides testable interfaces and dependency injection
  - Enables clean startup/shutdown hooks

## Hooks / Interfaces

### MCP Resource Interface
- **Static Resource**: `github://starred`
  - Returns list of all starred repos for authenticated user
  - Response: JSON array of MCP resources

- **Dynamic Resource Template**: `github://starred/{owner}/{repo}`
  - Returns details of a specific starred repository
  - URI variables: owner, repo
  - Response: Single MCP resource with full metadata

### GitHub API Interface
- `GET /user/starred` → authenticated user's starred repos
- `GET /users/{username}/starred` → specific user's starred repos
- Automatic pagination handling
- OAuth token-based authentication

### OAuth Flow
- Personal Access Token (PAT) authentication
- Required scopes: `public_repo`, `read:user`
- Token passed via GITHUB_TOKEN environment variable
- Secure token storage (not committed to repo)

### Dependency Wiring (uber-go/fx)
```
Config → GitHub Client → Resource Adapter → MCP Server
```

## Tests / Validation

### Unit Tests (Bazel)
```bash
bazel test //internal/github:github_test     # GitHub client helpers
bazel test //internal/resource:resource_test # Resource adapter logic
```

### Integration Tests
```bash
bazel test //tests:integration_test  # Full stack integration
```

### Test Coverage
- Unit tests for GitHub client (safe value extraction)
- Unit tests for resource adapter (MCP format conversion)
- Integration tests for config loading
- Integration tests for GitHub API calls (requires GITHUB_TOKEN)

### Security Validation
- OAuth tokens with minimum required scopes
- No token leakage in logs or responses
- Environment-based credential management
- Safe error handling without exposing sensitive data

## Build System Details

### Bazel Structure
```
MODULE.bazel              # bzlmod dependency management
WORKSPACE                 # Bazel workspace (minimal, uses bzlmod)
BUILD.bazel               # Root build file with Gazelle config
.bazelrc                  # Build configuration and modes
cmd/server/BUILD.bazel    # Binary target
internal/*/BUILD.bazel    # Library targets
tests/BUILD.bazel         # Test targets
```

### Common Bazel Commands
```bash
# Build the server
bazel build //cmd/server:mcp-server

# Run all tests
bazel test //...

# Run specific test
bazel test //internal/github:github_test

# Clean build artifacts
bazel clean

# Update BUILD files with Gazelle
bazel run //:gazelle

# Build with optimization
bazel build --config=opt //cmd/server:mcp-server

# Build for CI
bazel build --config=ci //...
```

### Convenience Script
Use `./bazel.sh` for common operations:
```bash
./bazel.sh build   # Build the binary
./bazel.sh test    # Run all tests
./bazel.sh run     # Build and run the server
./bazel.sh clean   # Clean artifacts
```

## Project Structure
```
mcp-server/
├── MODULE.bazel              # Bazel module configuration (bzlmod)
├── WORKSPACE                 # Bazel workspace
├── BUILD.bazel               # Root build file
├── .bazelrc                  # Bazel configuration
├── bazel.sh                  # Convenience script
├── cmd/
│   └── server/
│       ├── BUILD.bazel       # Binary build rules
│       └── main.go           # Application entry point (fx setup)
├── internal/
│   ├── config/
│   │   ├── BUILD.bazel       # Config library rules
│   │   └── config.go         # Environment configuration
│   ├── github/
│   │   ├── BUILD.bazel       # GitHub client rules
│   │   ├── client.go         # GitHub API wrapper
│   │   └── client_test.go    # Unit tests
│   ├── resource/
│   │   ├── BUILD.bazel       # Resource adapter rules
│   │   ├── adapter.go        # MCP resource mapping
│   │   └── adapter_test.go   # Unit tests
│   └── server/
│       ├── BUILD.bazel       # MCP server rules
│       └── server.go         # MCP protocol implementation
├── tests/
│   ├── BUILD.bazel           # Integration test rules
│   └── integration_test.go   # Integration tests
├── .env.example              # Environment template
├── .gitignore                # Git ignore (includes bazel-*)
├── go.mod                    # Go dependencies (source of truth)
├── go.sum                    # Go dependency checksums
├── CLAUDE.md                 # This file
├── README.md                 # User documentation
└── QUICKSTART.md             # Quick start guide
```

## Development Workflow

### Initial Setup
1. Install Bazel 8.4+
2. Set up environment: `cp .env.example .env`
3. Add GITHUB_TOKEN to .env
4. Build: `./bazel.sh build`

### Making Changes
1. Edit Go source files
2. Update BUILD files if needed (or run `bazel run //:gazelle`)
3. Run tests: `./bazel.sh test`
4. Build: `./bazel.sh build`

### Adding Dependencies
1. Update `go.mod`: `go get github.com/example/package`
2. Update MODULE.bazel if needed
3. Run Gazelle: `bazel run //:gazelle`
4. Commit changes to go.mod, go.sum, and BUILD files

## Future Extensions
- Add more GitHub resources (owned repos, issues, pull requests)
- Support caching layer for performance (Redis, in-memory)
- Expose prompt templates for summarizing repo metadata
- Add WebSocket/HTTP transports (in addition to stdio)
- Implement rate limiting and request throttling
- Add metrics and observability (Prometheus, OpenTelemetry)
- Support multiple authentication methods (OAuth app flow, GitHub App)
- Add Bazel remote caching for faster builds
- Container image builds with rules_docker/rules_oci
- Kubernetes deployment manifests with rules_k8s

## Dependencies (Managed by bzlmod + go.mod)

### Primary Dependencies
- **mcp-go** (github.com/mark3labs/mcp-go) v0.43.1 - MCP server framework
- **go-github** (github.com/google/go-github/v57) v57.0.0 - GitHub API client
- **oauth2** (golang.org/x/oauth2) v0.33.0 - OAuth 2.0 support
- **fx** (go.uber.org/fx) v1.24.0 - Dependency injection framework

### Build Dependencies
- **rules_go** (bazel_dep) v0.50.1 - Bazel Go rules
- **gazelle** (bazel_dep) v0.39.1 - BUILD file generator

### Transitive Dependencies
- go.uber.org/dig - Dependency injection (used by fx)
- go.uber.org/zap - Structured logging (used by fx)
- go.uber.org/multierr - Error handling (used by fx)
- golang.org/x/sys - System calls (used by various packages)

## Notes
- This project uses bzlmod (MODULE.bazel) instead of legacy WORKSPACE macros
- Go version 1.24.0 is used by Bazel for compatibility with rules_go 0.50.1
- The system may have Go 1.25.x installed, but Bazel uses its own hermetic Go SDK
- Bazel provides reproducible builds across different machines and OS
- All BUILD.bazel files can be auto-generated with Gazelle from Go source code
- The project supports both Bazel and traditional Go tooling (go build, go test)
