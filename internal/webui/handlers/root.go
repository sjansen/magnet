package handlers

import (
	"net/http"

	"github.com/sjansen/magnet/internal/build"
	"github.com/sjansen/magnet/internal/config"
	"github.com/sjansen/magnet/internal/webui/pages"
)

// Root is the default app starting page.
type Root struct {
	staticRoot string
}

// NewRoot creates a new root page handler.
func NewRoot(cfg *config.WebUI) *Root {
	return &Root{
		staticRoot: cfg.StaticRoot,
	}
}

// ServeHTTP handles reqeusts for the root page.
func (p *Root) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	page := &pages.RootPage{
		GitSHA:    build.GitSHA,
		Timestamp: build.Timestamp,
	}
	page.StaticRoot = p.staticRoot
	pages.WriteResponse(w, page)
}
