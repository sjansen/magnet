package pages

import (
	"fmt"
	"html"
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
	p.BeginContent(w, "Profile")
	io.WriteString(w, `
<div class="bg-white shadow overflow-hidden sm:rounded-lg">
  <div class="px-4 py-5 sm:px-6">
  <h3 class="text-lg leading-6 font-medium text-gray-900">
    Profile
  </h3>
  </div>`)
	io.WriteString(w, `
  <div class="bg-white shadow overflow-hidden sm:rounded-lg">`)
	io.WriteString(w, `
    <div class="border-t bg-gray-200">`)
	io.WriteString(w, `
    <dl>`)
	even := false
	for key, values := range p.Attrs {
		if even = !even; even {
			io.WriteString(w, `
      <div class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">`)
		} else {
			io.WriteString(w, `
      <div class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">`)
		}
		fmt.Fprintf(w, `
        <dt class=" row-span-%d `, len(values))
		io.WriteString(w, `text-sm font-medium text-gray-500">`)
		io.WriteString(w, html.EscapeString(key))
		io.WriteString(w, `</dt>`)
		for _, value := range values {
			io.WriteString(w, `
        <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">`)
			io.WriteString(w, html.EscapeString(value))
			io.WriteString(w, `</dd>`)
		}
		io.WriteString(w, `
      </div>`)
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
