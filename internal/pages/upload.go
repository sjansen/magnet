package pages

import (
	"fmt"
	"io"

	"github.com/sjansen/magnet/internal/util/s3form"
)

var _ Response = &ProfilePage{}

// UploadPage shows information about a user.
type UploadPage struct {
	Page
	Form *s3form.SignedForm
}

// WriteContent writes an HTTP response body.
func (p *UploadPage) WriteContent(w io.Writer) {
	if err := tmpls.ExecuteTemplate(w, "upload.html", p); err != nil {
		fmt.Println(err)
	}
}
