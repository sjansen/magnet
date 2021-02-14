package pages

import (
	"io"
)

var _ Response = &PublicPage{}

// PublicPage is a page the doesn't require authentication.
type PublicPage struct {
	Page
}

// WriteContent writes an HTTP response body.
func (p *PublicPage) WriteContent(w io.Writer) {
	p.BeginContent(w, "Public")
	io.WriteString(w, `<h1 class="bg-gray-50 rounded-lg shadow text-5xl text-center">Spoon!</h1>`)
	p.FinishContent(w)
}
