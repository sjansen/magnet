package pages

import (
	"fmt"
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
	if err := tmpls.ExecuteTemplate(w, "profile.html", p); err != nil {
		fmt.Println(err)
	}
}
