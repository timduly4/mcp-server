# Quick Start Guide

Get up and running with the GitHub Starred Repos MCP Server in minutes using **Bazel**.

## Prerequisites

- **Bazel 8.4+** installed ([download here](https://bazel.build/install))
- A GitHub account
- Git installed
- *(Optional)* Go 1.22+ for traditional workflow ([download here](https://go.dev/dl/))

## Step 1: Get the Code

```bash
git clone https://github.com/timduly4/mcp-server.git
cd mcp-server
```

## Step 2: Set Up GitHub Token

1. Go to https://github.com/settings/tokens
2. Click "Generate new token (classic)"
3. Give it a name (e.g., "MCP Server")
4. Select these scopes:
   - ✅ `public_repo`
   - ✅ `read:user`
5. Click "Generate token"
6. Copy the token (you won't see it again!)

## Step 3: Configure Environment

```bash
# Copy the example config
cp .env.example .env

# Edit .env and paste your token
nano .env
# or
vim .env
# or use your favorite editor
```

Update the line:
```
GITHUB_TOKEN=your_github_token_here
```

## Step 4: Build and Run

### Option A: Using Bazel with Convenience Script (Recommended)

```bash
# Build the binary
./bazel.sh build

# Run the server
./bazel.sh run
```

### Option B: Using Bazel Directly

```bash
# Build
bazel build //cmd/server:mcp-server

# Run
source .env
bazel-bin/cmd/server/mcp-server_/mcp-server
```

### Option C: Using Go Commands Directly

```bash
# Install dependencies
go mod download

# Build
go build -o bin/mcp-server ./cmd/server

# Run
source .env && ./bin/mcp-server
```

## Step 5: Run Tests

Verify everything works by running the test suite:

### Using Bazel (Recommended)

```bash
# Run all tests
./bazel.sh test

# Or directly
bazel test //...
```

### Using Go

```bash
go test ./...
```

## Step 6: Test It

The MCP server runs on stdio and communicates via JSON-RPC. You can test it by:

### Using an MCP Client

If you have an MCP client application (like Claude Desktop), configure it to use this server:

**For Bazel build:**
```json
{
  "mcpServers": {
    "github-starred": {
      "command": "/path/to/mcp-server/bazel-bin/cmd/server/mcp-server_/mcp-server",
      "env": {
        "GITHUB_TOKEN": "your_token_here"
      }
    }
  }
}
```

**For Go build:**
```json
{
  "mcpServers": {
    "github-starred": {
      "command": "/path/to/mcp-server/bin/mcp-server",
      "env": {
        "GITHUB_TOKEN": "your_token_here"
      }
    }
  }
}
```

### Manual Testing (JSON-RPC)

Send a JSON-RPC request to list resources:

**With Bazel:**
```bash
echo '{"jsonrpc":"2.0","id":1,"method":"resources/list","params":{}}' | bazel-bin/cmd/server/mcp-server_/mcp-server
```

**With Go:**
```bash
echo '{"jsonrpc":"2.0","id":1,"method":"resources/list","params":{}}' | ./bin/mcp-server
```

Or test reading the starred repos resource:

**With Bazel:**
```bash
echo '{"jsonrpc":"2.0","id":1,"method":"resources/read","params":{"uri":"github://starred"}}' | bazel-bin/cmd/server/mcp-server_/mcp-server
```

**With Go:**
```bash
echo '{"jsonrpc":"2.0","id":1,"method":"resources/read","params":{"uri":"github://starred"}}' | ./bin/mcp-server
```

## Available Resources

Once running, the server exposes:

1. **List all starred repos**: `github://starred`
2. **Get specific repo**: `github://starred/{owner}/{repo}`

## Troubleshooting

### "GITHUB_TOKEN environment variable is required"

Make sure you:
1. Created the `.env` file
2. Added your token to it
3. Used `source .env` before running

### "401 Unauthorized"

Your token might be:
- Expired
- Invalid
- Missing the required scopes

Generate a new token with the correct scopes.

### Build errors about missing packages

Update your Go version:

```bash
# macOS with Homebrew
brew upgrade go

# Or download from https://go.dev/dl/
```

## Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Explore the code in `internal/`
- Run tests: `./bazel.sh test` or `go test ./...`
- Check out [CLAUDE.md](CLAUDE.md) for architecture details

## Common Commands

### Bazel Commands

```bash
./bazel.sh build    # Build the binary
./bazel.sh test     # Run tests
./bazel.sh run      # Build and run
./bazel.sh clean    # Clean build artifacts
./bazel.sh gazelle  # Update BUILD files
./bazel.sh help     # Show all available commands
```

### Go Commands

```bash
go build -o bin/mcp-server ./cmd/server  # Build the binary
go test ./...                             # Run tests
go mod download                           # Download dependencies
go run ./cmd/server                       # Run directly
```

## Why Bazel?

Bazel provides several advantages:

- **Reproducible builds**: Same input always produces the same output
- **Hermetic**: Isolated from system dependencies
- **Fast**: Aggressive caching and incremental builds
- **Scalable**: Works for projects of any size
- **Cross-platform**: Consistent builds across macOS, Linux, and Windows

## Getting Help

- Check the [README.md](README.md) for detailed docs
- Open an issue on GitHub
- Review the [MCP specification](https://modelcontextprotocol.io/)

Happy coding! 🚀
