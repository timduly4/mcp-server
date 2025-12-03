package main

import (
	"context"
	"log"

	"go.uber.org/fx"

	"github.com/timduly4/mcp-server/internal/config"
	"github.com/timduly4/mcp-server/internal/github"
	"github.com/timduly4/mcp-server/internal/resource"
	"github.com/timduly4/mcp-server/internal/server"
)

func main() {
	app := fx.New(
		// Provide configuration
		fx.Provide(config.Load),

		// Provide GitHub client
		fx.Provide(newGitHubClient),

		// Provide resource adapter
		fx.Provide(resource.NewAdapter),

		// Provide MCP server
		fx.Provide(server.NewMCPServer),

		// Invoke server startup
		fx.Invoke(runServer),
	)

	app.Run()
}

// newGitHubClient creates a GitHub client from configuration
func newGitHubClient(cfg *config.Config) *github.Client {
	ctx := context.Background()
	return github.NewClient(ctx, cfg.GitHubToken)
}

// runServer starts the MCP server
func runServer(lifecycle fx.Lifecycle, srv *server.MCPServer) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("GitHub Starred Repos MCP Server starting...")

			// Start the server in a goroutine so it doesn't block fx startup
			go func() {
				if err := srv.Start(ctx); err != nil {
					log.Fatalf("Server failed: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Server shutting down...")
			return nil
		},
	})
}
