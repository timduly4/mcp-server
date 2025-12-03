package resource

import (
	"encoding/json"
	"fmt"

	"github.com/timduly4/mcp-server/internal/github"
)

// Adapter converts GitHub data to MCP resource format
type Adapter struct {
	githubClient *github.Client
}

// NewAdapter creates a new resource adapter
func NewAdapter(githubClient *github.Client) *Adapter {
	return &Adapter{
		githubClient: githubClient,
	}
}

// MCPResource represents a resource in MCP format
type MCPResource struct {
	URI         string                 `json:"uri"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	MimeType    string                 `json:"mimeType"`
	Contents    map[string]interface{} `json:"contents"`
}

// ListStarredResources returns starred repositories as MCP resources
func (a *Adapter) ListStarredResources() ([]MCPResource, error) {
	repos, err := a.githubClient.GetStarredRepos()
	if err != nil {
		return nil, fmt.Errorf("failed to get starred repos: %w", err)
	}

	resources := make([]MCPResource, 0, len(repos))
	for _, repo := range repos {
		resource := a.repoToMCPResource(repo)
		resources = append(resources, resource)
	}

	return resources, nil
}

// GetStarredResource returns a specific starred repository as an MCP resource
func (a *Adapter) GetStarredResource(fullName string) (*MCPResource, error) {
	repos, err := a.githubClient.GetStarredRepos()
	if err != nil {
		return nil, fmt.Errorf("failed to get starred repos: %w", err)
	}

	for _, repo := range repos {
		if repo.FullName == fullName {
			resource := a.repoToMCPResource(repo)
			return &resource, nil
		}
	}

	return nil, fmt.Errorf("repository %s not found in starred repos", fullName)
}

// repoToMCPResource converts a GitHub starred repo to MCP resource format
func (a *Adapter) repoToMCPResource(repo github.StarredRepo) MCPResource {
	uri := fmt.Sprintf("github://starred/%s", repo.FullName)

	contents := map[string]interface{}{
		"name":        repo.Name,
		"full_name":   repo.FullName,
		"owner":       repo.Owner,
		"description": repo.Description,
		"url":         repo.URL,
		"html_url":    repo.HTMLURL,
		"language":    repo.Language,
		"stars":       repo.Stars,
		"forks":       repo.Forks,
		"updated_at":  repo.UpdatedAt,
	}

	description := repo.Description
	if description == "" {
		description = fmt.Sprintf("Starred repository: %s", repo.FullName)
	}

	return MCPResource{
		URI:         uri,
		Name:        repo.FullName,
		Description: description,
		MimeType:    "application/json",
		Contents:    contents,
	}
}

// ToJSON converts resources to JSON format
func (a *Adapter) ToJSON(resources []MCPResource) ([]byte, error) {
	data, err := json.MarshalIndent(resources, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resources to JSON: %w", err)
	}
	return data, nil
}
