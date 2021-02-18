package pages

import (
	"fmt"
	"io"
)

var _ Response = &PrefixPage{}

type PrefixPage struct {
	BrowsePage
	Prefixes []string
	Objects  map[string]string
}

// WriteContent writes an HTTP response body.
func (p *PrefixPage) WriteContent(w io.Writer) {
	err := tmpls.ExecuteTemplate(w, "prefix.html", p)
	if err != nil {
		fmt.Println(err)
	}
}
