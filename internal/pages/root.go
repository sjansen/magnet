package pages

import (
	"html"
	"io"
)

var _ Response = &RootPage{}

// RootPage is the default application starting poing.
type RootPage struct {
	Page
	GitSHA    string
	Timestamp string
}

// WriteContent writes an HTTP response body.
func (p *RootPage) WriteContent(w io.Writer) {
	p.BeginContent(w, "")
	io.WriteString(w, `
<div class="bg-white shadow overflow-hidden sm:rounded-lg">
  <div class="px-4 py-5 sm:px-6">
  <h3 class="text-lg leading-6 font-medium text-gray-900">
    Magnet
  </h3>
  </div>`)
	io.WriteString(w, `
  <div class="bg-white shadow overflow-hidden sm:rounded-lg">`)
	io.WriteString(w, `
    <div class="border-t bg-gray-200">`)
	io.WriteString(w, `
    <dl>`)

	io.WriteString(w, `
	<div class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">`)
	io.WriteString(w, `
	<dt class=" row-span-%d text-sm font-medium text-gray-500">`)
	io.WriteString(w, "GitSHA")
	io.WriteString(w, `</dt>`)
	io.WriteString(w, `
	<dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">`)
	io.WriteString(w, html.EscapeString(p.GitSHA))
	io.WriteString(w, `</dd>`)
	io.WriteString(w, `
    </div>`)

	io.WriteString(w, `
	<div class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">`)
	io.WriteString(w, `
	<dt class=" row-span-%d text-sm font-medium text-gray-500">`)
	io.WriteString(w, "Timestamp")
	io.WriteString(w, `</dt>`)
	io.WriteString(w, `
	<dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">`)
	io.WriteString(w, html.EscapeString(p.Timestamp))
	io.WriteString(w, `</dd>`)
	io.WriteString(w, `
    </div>`)

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
