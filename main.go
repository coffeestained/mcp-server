package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"mcp-server/config"
	"mcp-server/logger"

	"mcp-server/handlers"
	"mcp-server/providers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Initialize Logger
	logger.InitLogger()

	// Load Configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Fatal error loading configuration", "error", err)
		return
	}
	slog.Info("Configuration loaded successfully")

	// Initialize Domains
	gitProvider := providers.NewGitProvider(&cfg.Github)
	slog.Info("Git provider initialized.")

	var soProvider *providers.StackOverflowProvider
	if cfg.StackExchange.APIKey == "" {
		slog.Warn("StackExchange apiKey is missing. Requests will have a lower rate limit.")
	}
	soProvider = providers.NewStackOverflowProvider(&cfg.StackExchange)
	slog.Info("Stack Overflow provider initialized.")

	var openapiProvider *providers.OpenAPIProvider
	if len(cfg.OpenAPI.Schemas) > 0 {
		openapiProvider = providers.NewOpenAPIProvider(&cfg.OpenAPI)
		slog.Info("OpenAPI provider initialized with configured schemas.")
	} else {
		slog.Warn("OpenAPI provider not configured (no schemas listed). This feature will be disabled.")
	}

	// Setup Router and Routes
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		// --- Git Routes ---
		gitHandler := handlers.NewGitHandler(gitProvider)
		r.Get("/repos", gitHandler.ListRepos)
		r.Get("/repos/{repoName}/tree/*", gitHandler.GetTree)
		r.Get("/repos/{repoName}/blob/*", gitHandler.GetBlob)

		// --- OpenAPI Routes ---
		if openapiProvider != nil {
			openapiHandler := handlers.NewOpenAPIHandler(openapiProvider)
			r.Get("/openapi", openapiHandler.ListSchemas)
			r.Get("/openapi/{schemaName}", openapiHandler.GetSchema)
		} else {
			r.Get("/openapi/*", handlers.FeatureNotConfiguredHandler("OpenAPI"))
		}

		// --- Stack Overflow Route ---
		soHandler := handlers.NewStackOverflowHandler(soProvider)
		r.Get("/stackoverflow/search", soHandler.Search)
	})

	// Start Server
	port := fmt.Sprintf(":%s", cfg.Server.Port)
	slog.Info("Starting MCP server", "port", port)
	if err := http.ListenAndServe(port, r); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}