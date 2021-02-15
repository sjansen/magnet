package handlers

import (
	"net/http"

	"github.com/sjansen/magnet/internal/build"
	"github.com/sjansen/magnet/internal/pages"
)

// Root is the app's root page.
func Root(w http.ResponseWriter, r *http.Request) {
	pages.WriteResponse(w, &pages.RootPage{
		GitSHA:    build.GitSHA,
		Timestamp: build.Timestamp,
	})
}
