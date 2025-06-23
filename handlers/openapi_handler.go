package handlers

import (
	"encoding/json"
	"net/http"
	"mcp-server/providers"
	"github.com/go-chi/chi/v5"
)

type OpenAPIHandler struct {
	provider *providers.OpenAPIProvider
}

// Bootstraps the OpenApi Handler
func NewOpenAPIHandler(p *providers.OpenAPIProvider) *OpenAPIHandler {
	return &OpenAPIHandler{provider: p}
}

// ListSchemas lists the names of available schemas.
func (h *OpenAPIHandler) ListSchemas(w http.ResponseWriter, r *http.Request) {
	names := h.provider.ListSchemas()
	respondWithJSON(w, http.StatusOK, map[string][]string{"available_schemas": names})
}

// GetSchema now gets the schema name from the URL path.
func (h *OpenAPIHandler) GetSchema(w http.ResponseWriter, r *http.Request) {
	schemaName := chi.URLParam(r, "schemaName")
	if schemaName == "" {
		respondWithError(w, http.StatusBadRequest, "Schema name is required in the URL path.")
		return
	}

	schemaBytes, err := h.provider.GetSchema(schemaName)
	if err != nil {
		// Could be not found or a fetch error.
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	var schemaData interface{}
	if err := json.Unmarshal(schemaBytes, &schemaData); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Fetched schema is not valid JSON: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, schemaData)
}