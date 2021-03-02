package handlers

import (
	"net/http"

	"github.com/crewjam/saml/samlsp"

	"github.com/sjansen/magnet/internal/pages"
)

// WhoAmI shows information about the current user.
func WhoAmI(w http.ResponseWriter, r *http.Request) {
	var attrs samlsp.Attributes
	s := samlsp.SessionFromContext(r.Context())
	if sa, ok := s.(samlsp.SessionWithAttributes); ok {
		attrs = sa.GetAttributes()
	}

	pages.WriteResponse(w, &pages.ProfilePage{
		Attrs: attrs,
	})
}
