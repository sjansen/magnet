package handlers

import (
	"net/http"

	"github.com/crewjam/saml/samlsp"

	"github.com/sjansen/magnet/internal/config"
	"github.com/sjansen/magnet/internal/webui/pages"
)

// WhoAmI shows information about the current user.
type WhoAmI struct {
	staticRoot string
}

// NewRoot creates a new root page handler.
func NewWhoAmI(cfg *config.WebUI) *WhoAmI {
	return &WhoAmI{
		staticRoot: cfg.StaticRoot,
	}
}

// WhoAmI shows information about the current user.
func (p *WhoAmI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var attrs samlsp.Attributes
	s := samlsp.SessionFromContext(r.Context())
	if sa, ok := s.(samlsp.SessionWithAttributes); ok {
		attrs = sa.GetAttributes()
	}

	page := &pages.ProfilePage{
		Attrs: attrs,
	}
	page.StaticRoot = p.staticRoot
	pages.WriteResponse(w, page)
}
