package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

// Client wraps the GitHub API client
type Client struct {
	client *github.Client
	ctx    context.Context
}

// StarredRepo represents a starred repository with relevant metadata
type StarredRepo struct {
	Name        string
	FullName    string
	Description string
	URL         string
	HTMLURL     string
	Language    string
	Stars       int
	Forks       int
	UpdatedAt   string
	Owner       string
}

// NewClient creates a new GitHub API client with OAuth token
func NewClient(ctx context.Context, token string) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &Client{
		client: github.NewClient(tc),
		ctx:    ctx,
	}
}

// GetStarredRepos fetches all starred repositories for the authenticated user
func (c *Client) GetStarredRepos() ([]StarredRepo, error) {
	var allRepos []StarredRepo
	opts := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		repos, resp, err := c.client.Activity.ListStarred(c.ctx, "", opts)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch starred repos: %w", err)
		}

		for _, repo := range repos {
			if repo.Repository == nil {
				continue
			}

			r := repo.Repository
			starredRepo := StarredRepo{
				Name:        getStringValue(r.Name),
				FullName:    getStringValue(r.FullName),
				Description: getStringValue(r.Description),
				URL:         getStringValue(r.URL),
				HTMLURL:     getStringValue(r.HTMLURL),
				Language:    getStringValue(r.Language),
				Stars:       getIntValue(r.StargazersCount),
				Forks:       getIntValue(r.ForksCount),
				Owner:       getOwnerLogin(r.Owner),
			}

			if r.UpdatedAt != nil {
				starredRepo.UpdatedAt = r.UpdatedAt.String()
			}

			allRepos = append(allRepos, starredRepo)
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepos, nil
}

// GetStarredReposForUser fetches starred repositories for a specific user
func (c *Client) GetStarredReposForUser(username string) ([]StarredRepo, error) {
	var allRepos []StarredRepo
	opts := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		repos, resp, err := c.client.Activity.ListStarred(c.ctx, username, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch starred repos for user %s: %w", username, err)
		}

		for _, repo := range repos {
			if repo.Repository == nil {
				continue
			}

			r := repo.Repository
			starredRepo := StarredRepo{
				Name:        getStringValue(r.Name),
				FullName:    getStringValue(r.FullName),
				Description: getStringValue(r.Description),
				URL:         getStringValue(r.URL),
				HTMLURL:     getStringValue(r.HTMLURL),
				Language:    getStringValue(r.Language),
				Stars:       getIntValue(r.StargazersCount),
				Forks:       getIntValue(r.ForksCount),
				Owner:       getOwnerLogin(r.Owner),
			}

			if r.UpdatedAt != nil {
				starredRepo.UpdatedAt = r.UpdatedAt.String()
			}

			allRepos = append(allRepos, starredRepo)
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepos, nil
}

// Helper functions to safely extract values from GitHub API responses
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func getIntValue(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func getOwnerLogin(owner *github.User) string {
	if owner == nil || owner.Login == nil {
		return ""
	}
	return *owner.Login
}
