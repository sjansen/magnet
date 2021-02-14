package pages

import (
	"html"
	"io"
)

var _ Response = &SecretPage{}

// SecretPage reveals a protected secret.
type SecretPage struct {
	Page
	Secret string
}

// WriteContent writes an HTTP response body.
func (p *SecretPage) WriteContent(w io.Writer) {
	p.BeginContent(w, "Secret")
	io.WriteString(w, `
<div class="bg-white shadow overflow-hidden sm:rounded-lg">
  <div class="px-4 py-5 sm:px-6">
  <h3 class="text-lg leading-6 font-medium text-gray-900">
    Secret
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
	io.WriteString(w, "Secret")
	io.WriteString(w, `</dt>`)
	io.WriteString(w, `
	<dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">`)
	io.WriteString(w, html.EscapeString(p.Secret))
	io.WriteString(w, `</dd>`)

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
