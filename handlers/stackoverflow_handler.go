package handlers

import (
	"net/http"
	"mcp-server/providers"
)

type StackOverflowHandler struct {
	provider *providers.StackOverflowProvider 
}

// Bootstraps the StackOverflow Handler
func NewStackOverflowHandler(p *providers.StackOverflowProvider) *StackOverflowHandler {
	return &StackOverflowHandler{provider: p}
}

func (h *StackOverflowHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondWithError(w, http.StatusBadRequest, "Query parameter 'q' is required")
		return
	}

	results, err := h.provider.Search(query)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to search Stack Overflow: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, results)
}