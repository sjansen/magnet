package pages

import (
	"fmt"
	"html"
	"io"
)

var _ Response = &BrowserPage{}

// BrowserPage is the default application starting poing.
type BrowserPage struct {
	Page
	Prefix   string
	Prefixes []string
	Objects  map[string]string
}

// WriteContent writes an HTTP response body.
func (p *BrowserPage) WriteContent(w io.Writer) {
	p.BeginContent(w, p.Prefix)
	io.WriteString(w, `
<div class="bg-white shadow overflow-hidden sm:rounded-lg">
  <div class="px-4 py-5 sm:px-6">
  <h3 class="text-lg leading-6 font-medium text-gray-900">`)
	io.WriteString(w, html.EscapeString(p.Prefix))
	io.WriteString(w, `</h3>
  </div>`)
	io.WriteString(w, `
  <div class="bg-white shadow overflow-hidden sm:rounded-lg">`)
	io.WriteString(w, `
    <div class="border-t bg-gray-200">`)
	io.WriteString(w, `
    <dl>`)
	even := false
	for _, prefix := range p.Prefixes {
		if even = !even; even {
			io.WriteString(w, `
      <div class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">`)
		} else {
			io.WriteString(w, `
      <div class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">`)
		}
		fmt.Fprintf(w, `<a href="%[1]s">%[1]s</a>`, html.EscapeString(prefix))
		io.WriteString(w, `</div>`)
	}
	for object, icon := range p.Objects {
		if even = !even; even {
			io.WriteString(w, `
      <div class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">`)
		} else {
			io.WriteString(w, `
      <div class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">`)
		}
		fmt.Fprintf(w, `<img src="%s" height="20" width="20" />`,
			html.EscapeString(icon),
		)
		io.WriteString(w, html.EscapeString(object))
		io.WriteString(w, `</div>`)
	}
	io.WriteString(w, `
    </dl>`)
	io.WriteString(w, `
    </div>`)
	io.WriteString(w, `
  </div>`)
	io.WriteString(w, `
</div>
`)
	p.FinishContent(w)
}
