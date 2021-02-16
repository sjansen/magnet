package pages

import (
	"fmt"
	"html"
	"io"
	"strings"
)

var _ Response = &BrowserPage{}

// BrowserPage is the default application starting poing.
type BrowserPage struct {
	Page
	Prefix   string
	Prefixes []string
	Objects  map[string]string
	Content  struct {
		Timestamp string
		Size      string
	}
}

// WriteContent writes an HTTP response body.
func (p *BrowserPage) WriteContent(w io.Writer) {
	p.BeginContent(w, p.Prefix)
	io.WriteString(w, `
<div class="bg-white shadow overflow-hidden sm:rounded-lg">
  <div class="px-4 py-5 sm:px-6">`)
	p.writeHeader(w)
	io.WriteString(w, `
  </div>
  <div class="bg-white shadow overflow-hidden sm:rounded-lg">
    <div class="border-t bg-gray-200">`)
	if p.Content.Timestamp != "" {
		p.writeObject(w)
	} else {
		p.writeListing(w)
	}
	io.WriteString(w, `
    </div>`)
	io.WriteString(w, `
  </div>`)
	io.WriteString(w, `
</div>
`)
	if p.Content.Timestamp != "" {
		io.WriteString(w, `
 <div class="text-2xl text-center px-4 py-5 sm:px-6">
   <a href="/`)
		io.WriteString(w, html.EscapeString(p.Prefix))
		io.WriteString(w, `">Download</a>
</div>`)
	}
	p.FinishContent(w)
}

func (p *BrowserPage) writeHeader(w io.Writer) {
	io.WriteString(w, `
	<h3 class="text-lg leading-6 font-medium text-gray-900">`)
	parts := strings.Split(p.Prefix, "/")
	dirs := len(parts) - 1
	stop := dirs - 1
	for i := 0; i < stop; i++ {
		io.WriteString(w, `<a class="hover:underline" href="`)
		for j := i; j < stop; j++ {
			io.WriteString(w, `../`)
		}
		io.WriteString(w, `">`)
		io.WriteString(w, html.EscapeString(parts[i]))
		io.WriteString(w, `</a>/`)
	}
	if strings.HasSuffix(p.Prefix, "/") {
		io.WriteString(w, html.EscapeString(parts[stop]))
		io.WriteString(w, `/`)
	} else {
		io.WriteString(w, `<a class="hover:underline" href=".">`)
		io.WriteString(w, html.EscapeString(parts[stop]))
		io.WriteString(w, `</a>/`)
		io.WriteString(w, html.EscapeString(parts[dirs]))
	}
	io.WriteString(w, `</h3>`)
}

func (p *BrowserPage) writeListing(w io.Writer) {
	even := false
	for _, prefix := range p.Prefixes {
		if even = !even; even {
			io.WriteString(w, `
      <div class="bg-gray-50 px-4 py-5 sm:px-6">`)
		} else {
			io.WriteString(w, `
      <div class="bg-white px-4 py-5 sm:px-6">`)
		}
		fmt.Fprintf(w, `<a href="%[1]s">%[1]s</a>`, html.EscapeString(prefix))
		io.WriteString(w, `</div>`)
	}
	for object, icon := range p.Objects {
		if even = !even; even {
			io.WriteString(w, `
      <div class="bg-gray-50 px-4 py-5 sm:px-6">`)
		} else {
			io.WriteString(w, `
      <div class="bg-white px-4 py-5 sm:px-6">`)
		}
		io.WriteString(w, `
		<a class="hover:underline" href="`)
		io.WriteString(w, html.EscapeString(object))
		io.WriteString(w, `">`)
		fmt.Fprintf(w, `
		<img src="%s" height="20" width="20" class="inline mr-2" />`,
			html.EscapeString(icon),
		)
		io.WriteString(w, html.EscapeString(object))
		io.WriteString(w, `</a>`)
		io.WriteString(w, `
	  </div>`)
	}
}

func (p *BrowserPage) writeObject(w io.Writer) {
	io.WriteString(w, `
    <dl>`)

	io.WriteString(w, `
      <div class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-4 sm:gap-4 sm:px-6">`)
	io.WriteString(w, `
		<dt class="text-sm font-medium text-gray-500">Size</dt>`)
	io.WriteString(w, `
		<dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-3">`)
	io.WriteString(w, html.EscapeString(p.Content.Size))
	io.WriteString(w, `</dd>`)
	io.WriteString(w, `
	  </div>`)

	io.WriteString(w, `
      <div class="bg-white px-4 py-5 sm:grid sm:grid-cols-4 sm:gap-4 sm:px-6">`)
	io.WriteString(w, `
		<dt class="text-sm font-medium text-gray-500">Timestamp</dt>`)
	io.WriteString(w, `
		<dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-3">`)
	io.WriteString(w, html.EscapeString(p.Content.Timestamp))
	io.WriteString(w, `</dd>`)
	io.WriteString(w, `
	  </div>`)

	io.WriteString(w, `
    </dl>`)
}
