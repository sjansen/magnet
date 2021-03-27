package pages

import (
	"fmt"
	"io"

	"github.com/sjansen/magnet/internal/util/s3path"
)

var _ Response = &PrefixPage{}

// PrefixPage enables browsing to specific objects.
type PrefixPage struct {
	BrowsePage
	Prefixes []*s3path.S3Path
	Objects  map[string]string
}

// WriteContent writes an HTTP response body.
func (p *PrefixPage) WriteContent(w io.Writer) {
	if err := tmpls.ExecuteTemplate(w, "prefix.html", p); err != nil {
		fmt.Println(err)
	}
}
