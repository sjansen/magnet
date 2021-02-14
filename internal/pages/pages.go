package pages

import (
	"html"
	"io"
	"net/http"
)

// Response is an HTTP response.
type Response interface {
	ContentType() string
	StatusCode() int
	WriteContent(w io.Writer)
}

// WriteResponse writes an HTTP response.
func WriteResponse(w http.ResponseWriter, resp Response) {
	w.Header().Set("Content-Type", resp.ContentType())
	w.WriteHeader(resp.StatusCode())
	resp.WriteContent(w)
}

// Page is an abstract class providing a standard page structure.
type Page struct {
	Status int
	Title  string
}

// ContentType returns a MIME type.
func (p *Page) ContentType() string {
	return "text/html; charset=utf-8"
}

// StatusCode returns an HTTP status code.
func (p *Page) StatusCode() int {
	if p.Status == 0 {
		return http.StatusOK
	}
	return p.Status
}

// BeginContent writes the start of an HTMl page.
func (p *Page) BeginContent(w io.Writer, title string) {
	io.WriteString(w,
		`<!doctype html>
<html>
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <link href="https://unpkg.com/tailwindcss@^2/dist/tailwind.min.css" rel="stylesheet">
  <title>`)
	if title != "" {
		io.WriteString(w, html.EscapeString(title))
		io.WriteString(w, " - ")
	}
	io.WriteString(w, `Magnet</title>
</head>
<body>
  <main class="p-8">`)
}

// FinishContent writes the end of an HTMl page.
func (p *Page) FinishContent(w io.Writer) {
	io.WriteString(w,
		`  </main>
</body>
</html>`)
}
