package handlers

import (
	"net/http"

	"github.com/sjansen/magnet/internal/pages"
)

// Public doesn't require authentication.
func Public(w http.ResponseWriter, r *http.Request) {
	pages.WriteResponse(w, &pages.PublicPage{})
}
