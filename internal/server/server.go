package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/timduly4/mcp-server/internal/resource"
)

// MCPServer wraps the MCP server functionality
type MCPServer struct {
	server  *server.MCPServer
	adapter *resource.Adapter
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer(adapter *resource.Adapter) *MCPServer {
	// Create MCP server with metadata
	s := server.NewMCPServer(
		"GitHub Starred Repos MCP Server",
		"1.0.0",
		server.WithResourceCapabilities(true, false), // subscribe = false
		server.WithToolCapabilities(false),
		server.WithPromptCapabilities(false),
	)

	mcpServer := &MCPServer{
		server:  s,
		adapter: adapter,
	}

	// Register resources
	mcpServer.registerResources()

	return mcpServer
}

// registerResources sets up all MCP resource endpoints
func (m *MCPServer) registerResources() {
	// Static resource: List all starred repositories
	starredListResource := mcp.NewResource(
		"github://starred",
		"All Starred Repositories",
		mcp.WithMIMEType("application/json"),
		mcp.WithResourceDescription("List of all GitHub repositories starred by the authenticated user"),
	)

	m.server.AddResource(starredListResource, m.handleListStarred)

	// Dynamic resource template: Individual starred repository
	starredRepoTemplate := mcp.NewResourceTemplate(
		"github://starred/{owner}/{repo}",
		"Starred Repository Details",
		mcp.WithTemplateMIMEType("application/json"),
		mcp.WithTemplateDescription("Details of a specific starred repository"),
	)

	m.server.AddResourceTemplate(starredRepoTemplate, m.handleGetStarredRepo)
}

// handleListStarred handles requests for all starred repositories
func (m *MCPServer) handleListStarred(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	log.Printf("Fetching all starred repositories")

	resources, err := m.adapter.ListStarredResources()
	if err != nil {
		return nil, fmt.Errorf("failed to list starred resources: %w", err)
	}

	// Convert to JSON
	jsonData, err := m.adapter.ToJSON(resources)
	if err != nil {
		return nil, fmt.Errorf("failed to convert resources to JSON: %w", err)
	}

	// Return as MCP resource contents
	contents := []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "github://starred",
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}

	log.Printf("Returning %d starred repositories", len(resources))
	return contents, nil
}

// handleGetStarredRepo handles requests for a specific starred repository
func (m *MCPServer) handleGetStarredRepo(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Extract owner and repo from URI
	// URI format: github://starred/{owner}/{repo}
	log.Printf("Fetching starred repository: %s", request.Params.URI)

	// Parse the URI to extract owner/repo
	// For simplicity, we'll search through all starred repos
	// In production, you might want to parse the URI more carefully
	fullName := extractFullNameFromURI(request.Params.URI)
	if fullName == "" {
		return nil, fmt.Errorf("invalid URI format: %s", request.Params.URI)
	}

	resource, err := m.adapter.GetStarredResource(fullName)
	if err != nil {
		return nil, fmt.Errorf("failed to get starred resource: %w", err)
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource to JSON: %w", err)
	}

	contents := []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}

	log.Printf("Returning starred repository: %s", fullName)
	return contents, nil
}

// extractFullNameFromURI extracts owner/repo from github://starred/{owner}/{repo}
func extractFullNameFromURI(uri string) string {
	// Simple URI parsing - in production, use a proper URI parser
	// Expected format: github://starred/{owner}/{repo}
	const prefix = "github://starred/"
	if len(uri) <= len(prefix) {
		return ""
	}

	// Extract the part after the prefix
	fullName := uri[len(prefix):]
	return fullName
}

// Start starts the MCP server using stdio transport
func (m *MCPServer) Start(ctx context.Context) error {
	log.Println("Starting MCP server on stdio...")
	return server.ServeStdio(m.server)
}
