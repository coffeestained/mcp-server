package handlers

import (
	"net/http"
	"strings"
	"mcp-server/providers"
	"github.com/go-chi/chi/v5"
)

type GitHandler struct {
	provider *providers.GitProvider
}

// Bootstraps the Git Handler
func NewGitHandler(p *providers.GitProvider) *GitHandler {
	return &GitHandler{provider: p}
}

// 
func (h *GitHandler) ListRepos(w http.ResponseWriter, r *http.Request) {
	repos := h.provider.ListConfiguredRepos()
	respondWithJSON(w, http.StatusOK, repos)
}

func (h *GitHandler) GetTree(w http.ResponseWriter, r *http.Request) {
	repoName := chi.URLParam(r, "repoName")
	path := chi.URLParam(r, "*")
	path = strings.TrimPrefix(path, "/")

	files, err := h.provider.ListFiles(r.Context(), repoName, path)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	type FileInfo struct {
		Name string `json:"name"`
		Path string `json:"path"`
		Type string `json:"type"`
	}

	var response []FileInfo
	for _, f := range files {
		response = append(response, FileInfo{
			Name: f.GetName(),
			Path: f.GetPath(),
			Type: f.GetType(),
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *GitHandler) GetBlob(w http.ResponseWriter, r *http.Request) {
	repoName := chi.URLParam(r, "repoName")
	path := chi.URLParam(r, "*")
	path = strings.TrimPrefix(path, "/")

	content, err := h.provider.GetFileContent(r.Context(), repoName, path)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"content": content})
}