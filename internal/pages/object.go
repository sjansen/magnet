package pages

import (
	"fmt"
	"io"
)

var _ Response = &ObjectPage{}

// ObjectPage displays metadata for a specific object.
type ObjectPage struct {
	BrowsePage
	Timestamp string
	Size      string
}

// WriteContent writes an HTTP response body.
func (p *ObjectPage) WriteContent(w io.Writer) {
	err := tmpls.ExecuteTemplate(w, "object.html", p)
	if err != nil {
		fmt.Println(err)
	}
}
