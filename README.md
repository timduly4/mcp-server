# GitHub Starred Repos MCP Server

An MCP (Model Context Protocol) server in Go that exposes GitHub repositories starred by a user as resources for LLM applications.

Built with **Bazel** for reproducible, hermetic builds across all platforms.

## Features

- Expose starred GitHub repositories as MCP resources
- List all starred repositories for authenticated user
- Query individual starred repository details
- Full MCP compliance with JSON-RPC over stdio
- OAuth-secured GitHub API integration
- Modular architecture with dependency injection using uber-go/fx
- **Bazel build system** for reproducible builds
- Hermetic Go SDK management with rules_go
- Auto-generated BUILD files with Gazelle

## Prerequisites

### Required
- **Bazel 8.4+** ([install instructions](https://bazel.build/install))
- GitHub Personal Access Token with appropriate scopes:
  - `public_repo`
  - `read:user`

### Optional (for traditional Go workflow)
- Go 1.22+ (tested with Go 1.25.5)
- Note: Bazel manages its own Go SDK (1.24.0) automatically

## Installation

### Clone the Repository

```bash
git clone https://github.com/timduly4/mcp-server.git
cd mcp-server
```

### Build with Bazel (Recommended)

Bazel provides reproducible, hermetic builds:

```bash
# Build the binary
bazel build //cmd/server:mcp-server

# Or use the convenience script
./bazel.sh build
```

The binary will be at: `bazel-bin/cmd/server/mcp-server_/mcp-server`

### Build with Go (Alternative)

```bash
# Install dependencies
go mod download

# Build the binary
go build -o bin/mcp-server ./cmd/server
```

### Configuration

1. Copy the example environment file:
```bash
cp .env.example .env
```

2. Edit `.env` and add your GitHub Personal Access Token:
```bash
GITHUB_TOKEN=your_github_token_here
```

To generate a GitHub token:
1. Go to https://github.com/settings/tokens
2. Click "Generate new token (classic)"
3. Select scopes: `public_repo`, `read:user`
4. Generate and copy the token

## Usage

### Running the Server

#### With Bazel (Recommended)

```bash
# Build and run using the convenience script
./bazel.sh run

# Or manually
source .env
bazel-bin/cmd/server/mcp-server_/mcp-server
```

#### With Go Binary

```bash
# Load environment variables and run
source .env
./bin/mcp-server
```

Or set the environment variable inline:

```bash
GITHUB_TOKEN=your_token_here ./bin/mcp-server
```

### MCP Resources

The server exposes the following resources:

#### 1. List All Starred Repositories

**URI:** `github://starred`

**Description:** Returns a list of all GitHub repositories starred by the authenticated user.

**Response Format:**
```json
[
  {
    "uri": "github://starred/owner/repo",
    "name": "owner/repo",
    "description": "Repository description",
    "mimeType": "application/json",
    "contents": {
      "name": "repo",
      "full_name": "owner/repo",
      "owner": "owner",
      "description": "Repository description",
      "url": "https://api.github.com/repos/owner/repo",
      "html_url": "https://github.com/owner/repo",
      "language": "Go",
      "stars": 42,
      "forks": 10,
      "updated_at": "2024-01-01"
    }
  }
]
```

#### 2. Get Specific Starred Repository

**URI Template:** `github://starred/{owner}/{repo}`

**Description:** Returns details of a specific starred repository.

**Example:** `github://starred/mark3labs/mcp-go`

## Architecture

The project follows a modular architecture with clear separation of concerns:

```
mcp-server/
├── cmd/server/              # Main application entry point
│   └── main.go             # fx dependency injection setup
├── internal/
│   ├── config/             # Configuration management
│   │   └── config.go       # Environment variable loading
│   ├── github/             # GitHub API client
│   │   ├── client.go       # GitHub REST API wrapper
│   │   └── client_test.go  # Unit tests
│   ├── resource/           # MCP resource adapter
│   │   ├── adapter.go      # Maps GitHub data to MCP format
│   │   └── adapter_test.go # Unit tests
│   └── server/             # MCP server implementation
│       └── server.go       # MCP protocol handling
├── bin/                    # Compiled binaries
├── tests/                  # Integration tests
├── .env.example            # Example environment configuration
├── .gitignore              # Git ignore rules
├── CLAUDE.md               # Project requirements and design
├── go.mod                  # Go module dependencies
└── README.md               # This file
```

### Key Components

1. **Server Module** (`internal/server`)
   - Implements MCP server using `mcp-go` framework
   - Registers resource endpoints
   - Handles JSON-RPC requests over stdio

2. **GitHub Client Module** (`internal/github`)
   - Wraps GitHub REST API v3
   - Handles authentication with OAuth tokens
   - Provides normalized data structures
   - Supports pagination

3. **Resource Adapter Module** (`internal/resource`)
   - Maps GitHub API responses to MCP resource schema
   - Converts data to JSON format
   - Handles error normalization

4. **Dependency Injection** (`cmd/server`)
   - Uses uber-go/fx for wiring components
   - Provides testable interfaces
   - Manages application lifecycle

## Development

### Running Tests

#### With Bazel (Recommended)

```bash
# Run all tests
bazel test //...

# Or use the convenience script
./bazel.sh test

# Run specific test
bazel test //internal/github:github_test

# Run with detailed output
bazel test --test_output=all //...
```

#### With Go

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests verbosely
go test -v ./internal/github/... ./internal/resource/...
```

### Building

#### With Bazel

```bash
# Build for current platform
bazel build //cmd/server:mcp-server

# Build with optimization
bazel build --config=opt //cmd/server:mcp-server

# Build for debugging
bazel build --config=debug //cmd/server:mcp-server

# Cross-compilation (requires platform configs)
bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/server:mcp-server
```

#### With Go

```bash
# Build for current platform
go build -o bin/mcp-server ./cmd/server

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/mcp-server-linux ./cmd/server

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o bin/mcp-server-darwin ./cmd/server

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o bin/mcp-server.exe ./cmd/server
```

### Bazel Commands Reference

```bash
# Update BUILD files automatically
bazel run //:gazelle

# Clean build artifacts
bazel clean

# Clean everything (including downloaded dependencies)
bazel clean --expunge

# Query all targets
bazel query //...

# View dependency graph
bazel query --output=graph //cmd/server:mcp-server

# Build everything
bazel build //...
```

### Convenience Script

The `bazel.sh` script provides shortcuts:

```bash
./bazel.sh build    # Build the binary
./bazel.sh test     # Run all tests
./bazel.sh run      # Build and run the server
./bazel.sh clean    # Clean artifacts
./bazel.sh gazelle  # Update BUILD files
./bazel.sh format   # Format BUILD files
./bazel.sh query    # Query all targets
./bazel.sh help     # Show help
```

### Code Organization

- Keep business logic in `internal/` packages
- Write unit tests alongside implementation files (`*_test.go`)
- Use interfaces for testability
- Follow Go standard project layout
- Each package has a `BUILD.bazel` file defining build targets
- Use Gazelle to auto-generate BUILD files: `bazel run //:gazelle`

### Adding New Dependencies

1. Update `go.mod`:
   ```bash
   go get github.com/example/new-package
   ```

2. Update BUILD files:
   ```bash
   bazel run //:gazelle
   ```

3. If needed, update `MODULE.bazel` to explicitly declare dependencies

4. Commit changes to `go.mod`, `go.sum`, and all `BUILD.bazel` files

## MCP Compliance

This server conforms to the [MCP specification](https://modelcontextprotocol.io/docs/develop/build-server):

- ✅ JSON-RPC 2.0 protocol
- ✅ stdio transport
- ✅ Resource capabilities
- ✅ Static resources
- ✅ Dynamic resource templates (URI templates)
- ✅ Proper error handling
- ✅ OAuth security

## Dependencies

### Runtime Dependencies
- **mcp-go** (github.com/mark3labs/mcp-go) v0.43.1 - MCP server framework
- **go-github** (github.com/google/go-github/v57) v57.0.0 - GitHub API client
- **oauth2** (golang.org/x/oauth2) v0.33.0 - OAuth 2.0 authentication
- **fx** (go.uber.org/fx) v1.24.0 - Dependency injection framework

### Build Dependencies
- **Bazel** 8.4+ - Build system
- **rules_go** v0.50.1 - Bazel Go rules
- **Gazelle** v0.39.1 - BUILD file generator
- **Go SDK** 1.24.0 - Managed by Bazel (hermetic)

## Future Enhancements

Potential extensions as outlined in CLAUDE.md:

- [ ] Additional GitHub resources (owned repos, issues, pull requests)
- [ ] Caching layer for improved performance
- [ ] Prompt templates for repo metadata summarization
- [ ] WebSocket/HTTP transport options
- [ ] Rate limiting and request throttling
- [ ] Metrics and observability

## Troubleshooting

### Common Issues

**Error: "GITHUB_TOKEN environment variable is required"**
- Solution: Make sure to set the `GITHUB_TOKEN` environment variable before running the server.

**Error: "failed to fetch starred repos: 401 Unauthorized"**
- Solution: Verify your GitHub token is valid and has the required scopes.

**Build fails with "package X is not in GOROOT"**
- Solution: Upgrade Go to version 1.22 or higher using `brew upgrade go` (macOS) or download from https://go.dev/dl/

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## References

- [Model Context Protocol Specification](https://modelcontextprotocol.io/)
- [mcp-go Documentation](https://github.com/mark3labs/mcp-go)
- [GitHub REST API](https://docs.github.com/en/rest)
- [uber-go/fx Documentation](https://uber-go.github.io/fx/)
