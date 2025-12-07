package resource

import (
	"testing"

	"github.com/timduly4/mcp-server/internal/github"
)

func TestRepoToMCPResource(t *testing.T) {
	adapter := &Adapter{}

	repo := github.StarredRepo{
		Name:        "test-repo",
		FullName:    "owner/test-repo",
		Description: "A test repository",
		URL:         "https://api.github.com/repos/owner/test-repo",
		HTMLURL:     "https://github.com/owner/test-repo",
		Language:    "Go",
		Stars:       42,
		Forks:       10,
		UpdatedAt:   "2024-01-01",
		Owner:       "owner",
	}

	resource := adapter.repoToMCPResource(repo)

	// Test URI format
	expectedURI := "github://starred/owner/test-repo"
	if resource.URI != expectedURI {
		t.Errorf("URI = %v, want %v", resource.URI, expectedURI)
	}

	// Test name
	if resource.Name != repo.FullName {
		t.Errorf("Name = %v, want %v", resource.Name, repo.FullName)
	}

	// Test MIME type
	if resource.MimeType != "application/json" {
		t.Errorf("MimeType = %v, want application/json", resource.MimeType)
	}

	// Test contents
	if resource.Contents["name"] != repo.Name {
		t.Errorf("Contents[name] = %v, want %v", resource.Contents["name"], repo.Name)
	}

	if resource.Contents["stars"] != repo.Stars {
		t.Errorf("Contents[stars] = %v, want %v", resource.Contents["stars"], repo.Stars)
	}
}

func TestRepoToMCPResource_EmptyDescription(t *testing.T) {
	adapter := &Adapter{}

	repo := github.StarredRepo{
		Name:     "test-repo",
		FullName: "owner/test-repo",
		// Description is empty
	}

	resource := adapter.repoToMCPResource(repo)

	expectedDescription := "Starred repository: owner/test-repo"
	if resource.Description != expectedDescription {
		t.Errorf("Description = %v, want %v", resource.Description, expectedDescription)
	}
}

func TestToJSON(t *testing.T) {
	adapter := &Adapter{}

	resources := []MCPResource{
		{
			URI:         "github://starred/owner/repo1",
			Name:        "owner/repo1",
			Description: "Test repo 1",
			MimeType:    "application/json",
			Contents: map[string]interface{}{
				"name":  "repo1",
				"stars": 10,
			},
		},
	}

	jsonData, err := adapter.ToJSON(resources)
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}

	if len(jsonData) == 0 {
		t.Error("ToJSON() returned empty data")
	}

	// Verify it's valid JSON by checking for expected patterns
	jsonStr := string(jsonData)
	if !contains(jsonStr, "github://starred/owner/repo1") {
		t.Error("JSON does not contain expected URI")
	}
}

func TestRepoToMCPResourceForUser(t *testing.T) {
	adapter := &Adapter{}

	repo := github.StarredRepo{
		Name:        "test-repo",
		FullName:    "owner/test-repo",
		Description: "A test repository",
		URL:         "https://api.github.com/repos/owner/test-repo",
		HTMLURL:     "https://github.com/owner/test-repo",
		Language:    "Go",
		Stars:       42,
		Forks:       10,
		UpdatedAt:   "2024-01-01",
		Owner:       "owner",
	}

	username := "testuser"
	resource := adapter.repoToMCPResourceForUser(repo, username)

	// Test URI format includes username
	expectedURI := "github://starred/users/testuser/owner/test-repo"
	if resource.URI != expectedURI {
		t.Errorf("URI = %v, want %v", resource.URI, expectedURI)
	}

	// Test name
	if resource.Name != repo.FullName {
		t.Errorf("Name = %v, want %v", resource.Name, repo.FullName)
	}

	// Test MIME type
	if resource.MimeType != "application/json" {
		t.Errorf("MimeType = %v, want application/json", resource.MimeType)
	}

	// Test contents includes starred_by field
	if resource.Contents["starred_by"] != username {
		t.Errorf("Contents[starred_by] = %v, want %v", resource.Contents["starred_by"], username)
	}

	if resource.Contents["name"] != repo.Name {
		t.Errorf("Contents[name] = %v, want %v", resource.Contents["name"], repo.Name)
	}

	if resource.Contents["stars"] != repo.Stars {
		t.Errorf("Contents[stars] = %v, want %v", resource.Contents["stars"], repo.Stars)
	}
}

func TestRepoToMCPResourceForUser_EmptyDescription(t *testing.T) {
	adapter := &Adapter{}

	repo := github.StarredRepo{
		Name:     "test-repo",
		FullName: "owner/test-repo",
		// Description is empty
	}

	username := "testuser"
	resource := adapter.repoToMCPResourceForUser(repo, username)

	expectedDescription := "Repository owner/test-repo starred by testuser"
	if resource.Description != expectedDescription {
		t.Errorf("Description = %v, want %v", resource.Description, expectedDescription)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
