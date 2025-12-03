package tests

import (
	"context"
	"os"
	"testing"

	"github.com/timduly4/mcp-server/internal/config"
	"github.com/timduly4/mcp-server/internal/github"
	"github.com/timduly4/mcp-server/internal/resource"
)

// TestIntegration_GitHubClient tests the GitHub client with real API calls
// This test requires a valid GITHUB_TOKEN environment variable
func TestIntegration_GitHubClient(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("Skipping integration test: GITHUB_TOKEN not set")
	}

	ctx := context.Background()
	client := github.NewClient(ctx, token)

	// Test fetching starred repos
	repos, err := client.GetStarredRepos()
	if err != nil {
		t.Fatalf("Failed to fetch starred repos: %v", err)
	}

	t.Logf("Fetched %d starred repositories", len(repos))

	if len(repos) == 0 {
		t.Log("No starred repositories found (this is okay if the user hasn't starred any repos)")
		return
	}

	// Verify first repo has expected fields
	repo := repos[0]
	if repo.FullName == "" {
		t.Error("Repository FullName should not be empty")
	}
	if repo.URL == "" {
		t.Error("Repository URL should not be empty")
	}

	t.Logf("First starred repo: %s", repo.FullName)
}

// TestIntegration_ResourceAdapter tests the resource adapter with GitHub client
func TestIntegration_ResourceAdapter(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("Skipping integration test: GITHUB_TOKEN not set")
	}

	ctx := context.Background()
	client := github.NewClient(ctx, token)
	adapter := resource.NewAdapter(client)

	// Test listing starred resources
	resources, err := adapter.ListStarredResources()
	if err != nil {
		t.Fatalf("Failed to list starred resources: %v", err)
	}

	t.Logf("Got %d MCP resources", len(resources))

	if len(resources) == 0 {
		t.Log("No starred resources found (this is okay if the user hasn't starred any repos)")
		return
	}

	// Verify first resource has MCP format
	res := resources[0]
	if res.URI == "" {
		t.Error("Resource URI should not be empty")
	}
	if res.MimeType != "application/json" {
		t.Errorf("Resource MimeType = %s, want application/json", res.MimeType)
	}
	if res.Contents == nil {
		t.Error("Resource Contents should not be nil")
	}

	t.Logf("First resource URI: %s", res.URI)

	// Test JSON conversion
	jsonData, err := adapter.ToJSON(resources)
	if err != nil {
		t.Fatalf("Failed to convert to JSON: %v", err)
	}

	if len(jsonData) == 0 {
		t.Error("JSON data should not be empty")
	}

	t.Logf("JSON data length: %d bytes", len(jsonData))
}

// TestIntegration_Config tests configuration loading
func TestIntegration_Config(t *testing.T) {
	// Save original env var
	originalToken := os.Getenv("GITHUB_TOKEN")
	defer os.Setenv("GITHUB_TOKEN", originalToken)

	// Test with token
	testToken := "test_token_123"
	os.Setenv("GITHUB_TOKEN", testToken)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.GitHubToken != testToken {
		t.Errorf("GitHubToken = %s, want %s", cfg.GitHubToken, testToken)
	}

	// Test defaults
	if cfg.ServerPort != "8080" {
		t.Errorf("ServerPort = %s, want 8080", cfg.ServerPort)
	}
	if cfg.ServerHost != "localhost" {
		t.Errorf("ServerHost = %s, want localhost", cfg.ServerHost)
	}
}

// TestIntegration_ConfigMissingToken tests config loading without token
func TestIntegration_ConfigMissingToken(t *testing.T) {
	// Save original env var
	originalToken := os.Getenv("GITHUB_TOKEN")
	defer os.Setenv("GITHUB_TOKEN", originalToken)

	// Unset token
	os.Unsetenv("GITHUB_TOKEN")

	_, err := config.Load()
	if err == nil {
		t.Error("Expected error when GITHUB_TOKEN is missing, got nil")
	}
}
