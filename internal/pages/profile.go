package pages

import (
	"io"
)

var _ Response = &ProfilePage{}

// ProfilePage shows information about a user.
type ProfilePage struct {
	Page
	Attrs map[string][]string
}

// WriteContent writes an HTTP response body.
func (p *ProfilePage) WriteContent(w io.Writer) {
	tmpls.ExecuteTemplate(w, "profile.html", p)
}
