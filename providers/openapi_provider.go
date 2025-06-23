package providers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"mcp-server/config"
)

type OpenAPIProvider struct {
	schemas map[string]string
}

// Bootstraps the OpenAPI Provider
func NewOpenAPIProvider(cfg *config.OpenAPIConfig) *OpenAPIProvider {
	return &OpenAPIProvider{schemas: cfg.Schemas}
}

// GetSchema now takes a name to look up in the map.
func (p *OpenAPIProvider) GetSchema(name string) ([]byte, error) {
	path, ok := p.schemas[name]
	if !ok {
		return nil, fmt.Errorf("schema with name '%s' not configured", name)
	}

	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		resp, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return io.ReadAll(resp.Body)
	}

	return os.ReadFile(path)
}

// ListSchemas returns the names of all configured schemas.
func (p *OpenAPIProvider) ListSchemas() []string {
	names := make([]string, 0, len(p.schemas))
	for name := range p.schemas {
		names = append(names, name)
	}
	return names
}