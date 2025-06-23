# MCP (Multi-Context Provider) Server

A lightweight Go server designed to provide specific, read-only contexts (Git, OpenAPI, Stack Overflow) to assist an AI coding agent. The server is resilient, allowing features to be enabled or disabled via the configuration file.

## Note

This is still very early in development process. Many features are still being finalized.

## Features

*   Provides read-access to configured Git repositories via the GitHub API. Public repos work without an API key.
*   Serves multiple, named OpenAPI 3.0 schemas from URLs or local files.
*   Queries Stack Overflow and returns the top 10 results with comment trees.
*   Gracefully handles misconfigurations and disables features on a per-module basis with clear logging.

## Requirements

*   **Go 1.24+**
*   (Optional) A GitHub Personal Access Token for higher rate limits and access to private repos.
*   (Optional) A Stack Exchange API Key to increase request limits.

## Setup & Running

1.  **Clone the repository:**

    git clone <your-repo-url>
    cd mcp-server

2.  **Create your configuration file.** Create a file named `config.yaml`. You can start by copying the content below.

    **`config.yaml.sample`**

    # Server configuration is optional, defaults to port 8080
    server:
      port: "8080"
    
    # --- All provider sections below are optional. ---
    
    # GitHub Provider: apiKey is optional. 
    # Leave empty for public repos only (lower rate limit).
    github:
      apiKey: "YOUR_GITHUB_PERSONAL_ACCESS_TOKEN"
      repositories:
        chi: "go-chi/chi"
        mcp-server: "your-github-username/mcp-server"
    
    # Stack Exchange Provider: apiKey is optional.
    stackexchange:
      apiKey: "YOUR_STACKEXCHANGE_API_KEY"
    
    # OpenAPI Provider: A map of schema names to their path.
    openapi:
      schemas:
        petstore: "https://petstore3.swagger.io/api/v3/openapi.json"
        local_api: "docs/my_local_api.yaml"

3.  **Install dependencies.**

    go mod tidy

4.  **Run the server.**

    go run main.go
    
    The server will start on the configured port, with detailed startup logs in the console.

## API Endpoints

All endpoints are prefixed with `/api/v1`.

#### List Configured Repositories
`GET /repos`

    curl http://localhost:8080/api/v1/repos

#### List Files in a Repository Path
`GET /repos/{repoName}/tree/{path}`

    curl http://localhost:8080/api/v1/repos/chi/tree/

#### Get File Content
`GET /repos/{repoName}/blob/{path}`

    curl http://localhost:8080/api/v1/repos/chi/blob/README.md

#### List Available OpenAPI Schemas
`GET /openapi`

    curl http://localhost:8080/api/v1/openapi

#### Get a Specific OpenAPI Schema
`GET /openapi/{schemaName}`

    curl http://localhost:8080/api/v1/openapi/petstore

#### Search Stack Overflow
`GET /stackoverflow/search?q={query}`

    curl "http://localhost:8080/api/v1/stackoverflow/search?q=golang%20http%20post%20json"