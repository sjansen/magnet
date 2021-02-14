package handlers

import (
	"net/http"

	"github.com/sjansen/magnet/internal/pages"
)

// Secret is a protected page.
func Secret(w http.ResponseWriter, r *http.Request) {
	pages.WriteResponse(w, &pages.SecretPage{
		Secret: "hunter2",
	})
}
