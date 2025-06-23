package providers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http" 
	"strings"

	"mcp-server/config"

	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"
)

type GitProvider struct {
	client *github.Client
	repos  map[string]string
}

// Bootstraps the GitProvider
func NewGitProvider(cfg *config.GithubConfig) *GitProvider {
	var httpClient *http.Client
	
	if cfg.APIKey != "" {
		slog.Info("GitHub API key found, initializing authenticated client.")
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.APIKey})
		httpClient = oauth2.NewClient(ctx, ts)
	} else {
		slog.Warn("GitHub API key not found. Initializing unauthenticated client. Rate limits will be lower.")
		// A nil httpClient tells the GitHub client to use http.DefaultClient
		httpClient = nil 
	}
	
	client := github.NewClient(httpClient)

	return &GitProvider{
		client: client,
		repos:  cfg.Repositories,
	}
}

// Returns repo details by shortname
func (p *GitProvider) getRepoFullName(shortName string) (owner, repo string, err error) {
	fullName, ok := p.repos[shortName]
	if !ok {
		return "", "", fmt.Errorf("repository '%s' not configured", shortName)
	}
	parts := strings.SplitN(fullName, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid repository format for '%s': expected 'owner/repo', got '%s'", shortName, fullName)
	}
	return parts[0], parts[1], nil
}

// Lists files
func (p *GitProvider) ListFiles(ctx context.Context, repoShortName, path string) ([]*github.RepositoryContent, error) {
	owner, repo, err := p.getRepoFullName(repoShortName)
	if err != nil {
		return nil, err
	}
	_, dirContent, _, err := p.client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return nil, fmt.Errorf("could not get contents for %s/%s at path '%s': %w", owner, repo, path, err)
	}
	return dirContent, nil
}

// Retrieves File Contents
func (p *GitProvider) GetFileContent(ctx context.Context, repoShortName, path string) (string, error) {
	owner, repo, err := p.getRepoFullName(repoShortName)
	if err != nil {
		return "", err
	}
	fileContent, _, _, err := p.client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return "", fmt.Errorf("could not get file content for %s/%s at path '%s': %w", owner, repo, path, err)
	}
	if fileContent == nil {
		return "", errors.New("path is a directory, not a file, or does not exist")
	}
	if fileContent.GetType() == "dir" {
		return "", errors.New("path is a directory, not a file")
	}
	content, err := fileContent.GetContent()
	if err == nil {
		return content, nil
	}
	if fileContent.GetEncoding() == "base64" && fileContent.Content != nil {
		rawB64Content := *fileContent.Content
		decoded, b64err := base64.StdEncoding.DecodeString(rawB64Content)
		if b64err != nil {
			return "", fmt.Errorf("base64 decoding failed for path %s: %w", path, b64err)
		}
		return string(decoded), nil
	}
	return "", fmt.Errorf("could not get or decode file content for %s: %w", path, err)
}

// List all available Repos
func (p *GitProvider) ListConfiguredRepos() map[string]string {
	return p.repos
}