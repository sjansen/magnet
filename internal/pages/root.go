package pages

import (
	"io"
)

var _ Response = &RootPage{}

// RootPage is the default application starting poing.
type RootPage struct {
	Page
	GitSHA    string
	Timestamp string
}

// WriteContent writes an HTTP response body.
func (p *RootPage) WriteContent(w io.Writer) {
	tmpls.ExecuteTemplate(w, "root.html", p)
}
