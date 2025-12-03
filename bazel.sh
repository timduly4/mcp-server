#!/bin/bash
# Convenience script for common Bazel commands

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

usage() {
    echo "Usage: $0 <command>"
    echo ""
    echo "Commands:"
    echo "  build       - Build the MCP server binary"
    echo "  test        - Run all tests"
    echo "  run         - Build and run the MCP server"
    echo "  clean       - Clean Bazel artifacts"
    echo "  gazelle     - Update BUILD files with Gazelle"
    echo "  format      - Format BUILD files"
    echo "  query       - Query Bazel targets"
    echo "  help        - Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 build"
    echo "  $0 test"
    echo "  $0 run"
    exit 1
}

check_env() {
    if [ ! -f .env ]; then
        echo -e "${RED}Error: .env file not found${NC}"
        echo "Copy .env.example to .env and configure your GITHUB_TOKEN"
        exit 1
    fi
    source .env
    if [ -z "$GITHUB_TOKEN" ]; then
        echo -e "${RED}Error: GITHUB_TOKEN not set in .env${NC}"
        exit 1
    fi
}

cmd_build() {
    echo -e "${GREEN}Building MCP server with Bazel...${NC}"
    bazel build //cmd/server:mcp-server
    echo -e "${GREEN}Build complete!${NC}"
    echo "Binary location: bazel-bin/cmd/server/mcp-server_/mcp-server"
}

cmd_test() {
    echo -e "${GREEN}Running tests with Bazel...${NC}"
    bazel test //...
}

cmd_run() {
    check_env
    echo -e "${GREEN}Building and running MCP server...${NC}"
    bazel build //cmd/server:mcp-server
    echo -e "${GREEN}Starting server...${NC}"
    bazel-bin/cmd/server/mcp-server_/mcp-server
}

cmd_clean() {
    echo -e "${YELLOW}Cleaning Bazel artifacts...${NC}"
    bazel clean
    echo -e "${GREEN}Clean complete${NC}"
}

cmd_gazelle() {
    echo -e "${GREEN}Running Gazelle to update BUILD files...${NC}"
    bazel run //:gazelle
    echo -e "${GREEN}BUILD files updated${NC}"
}

cmd_format() {
    echo -e "${GREEN}Formatting BUILD files...${NC}"
    if command -v buildifier &> /dev/null; then
        find . -name "BUILD.bazel" -o -name "WORKSPACE" -o -name "*.bzl" | xargs buildifier
        echo -e "${GREEN}BUILD files formatted${NC}"
    else
        echo -e "${YELLOW}Warning: buildifier not installed${NC}"
        echo "Install with: go install github.com/bazelbuild/buildtools/buildifier@latest"
    fi
}

cmd_query() {
    echo -e "${GREEN}Querying Bazel targets...${NC}"
    bazel query //...
}

# Main command dispatcher
case "${1:-help}" in
    build)
        cmd_build
        ;;
    test)
        cmd_test
        ;;
    run)
        cmd_run
        ;;
    clean)
        cmd_clean
        ;;
    gazelle)
        cmd_gazelle
        ;;
    format)
        cmd_format
        ;;
    query)
        cmd_query
        ;;
    help|*)
        usage
        ;;
esac
