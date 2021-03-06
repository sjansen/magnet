package pages

import (
	"fmt"
	"io"
	"net/url"
)

var _ Response = &ObjectPage{}

// ObjectPage displays metadata for a specific object.
type ObjectPage struct {
	BrowsePage
	URL       *url.URL
	MimeType  string
	Size      string
	Timestamp string
	Metadata  map[string]string
}

// WriteContent writes an HTTP response body.
func (p *ObjectPage) WriteContent(w io.Writer) {
	if err := tmpls.ExecuteTemplate(w, "object.html", p); err != nil {
		fmt.Println(err)
	}
}
