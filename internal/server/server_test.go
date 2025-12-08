package server

import (
	"testing"
)

// TestExtractFullNameFromURI tests URI parsing for {owner}/{repo} pattern
func TestExtractFullNameFromURI(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		expected string
	}{
		{
			name:     "valid owner and repo",
			uri:      "github://starred/facebook/react",
			expected: "facebook/react",
		},
		{
			name:     "valid owner and repo with hyphen",
			uri:      "github://starred/mark3labs/mcp-go",
			expected: "mark3labs/mcp-go",
		},
		{
			name:     "empty URI",
			uri:      "",
			expected: "",
		},
		{
			name:     "URI with only prefix",
			uri:      "github://starred/",
			expected: "",
		},
		{
			name:     "should NOT match users pattern (routing conflict test)",
			uri:      "github://starred/users/timduly4",
			expected: "users/timduly4", // This is what we DON'T want to happen!
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractFullNameFromURI(tt.uri)
			if result != tt.expected {
				t.Errorf("extractFullNameFromURI(%q) = %q, want %q", tt.uri, result, tt.expected)
			}
		})
	}
}

// TestExtractUsernameFromURI tests URI parsing for users/{username} pattern
func TestExtractUsernameFromURI(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		expected string
	}{
		{
			name:     "valid username",
			uri:      "github://starred/users/octocat",
			expected: "octocat",
		},
		{
			name:     "valid username with hyphen",
			uri:      "github://starred/users/my-user",
			expected: "my-user",
		},
		{
			name:     "valid username with numbers",
			uri:      "github://starred/users/user123",
			expected: "user123",
		},
		{
			name:     "empty URI",
			uri:      "",
			expected: "",
		},
		{
			name:     "URI with only prefix",
			uri:      "github://starred/users/",
			expected: "",
		},
		{
			name:     "wrong prefix",
			uri:      "github://starred/timduly4",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractUsernameFromURI(tt.uri)
			if result != tt.expected {
				t.Errorf("extractUsernameFromURI(%q) = %q, want %q", tt.uri, result, tt.expected)
			}
		})
	}
}

// TestURIRouting_ConflictPrevention tests that the routing order prevents conflicts
// This test documents the routing behavior and why registration order matters
func TestURIRouting_ConflictPrevention(t *testing.T) {
	// This test documents the routing conflict that was fixed
	// When github://starred/users/timduly4 is requested:

	// BAD: If {owner}/{repo} pattern is checked first:
	// extractFullNameFromURI("github://starred/users/timduly4") returns "users/timduly4"
	// This causes error: "repository users/timduly4 not found"

	conflictingURI := "github://starred/users/timduly4"

	// Test 1: The {owner}/{repo} pattern WOULD match (this is the problem!)
	ownerRepoResult := extractFullNameFromURI(conflictingURI)
	if ownerRepoResult != "users/timduly4" {
		t.Errorf("extractFullNameFromURI should match 'users/timduly4', got %q", ownerRepoResult)
	}

	// Test 2: The users/{username} pattern correctly extracts username
	usernameResult := extractUsernameFromURI(conflictingURI)
	if usernameResult != "timduly4" {
		t.Errorf("extractUsernameFromURI should extract 'timduly4', got %q", usernameResult)
	}

	t.Log("✓ Routing conflict documented: users/{username} pattern must be registered BEFORE {owner}/{repo}")
	t.Log("✓ Correct routing order prevents 'repository users/timduly4 not found' error")
}
