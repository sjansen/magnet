package pages

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
)

var tmpls *template.Template

func init() {
	fmt.Println("Parsing templates...")
	tmpls = template.New("pages")
	tmpls.Funcs(template.FuncMap{
		"HasPrefix": strings.HasPrefix,
	})
	template.Must(tmpls.ParseGlob("templates/*.html"))
}

// Href is an HTML <a href="...">...</a>
type Href struct {
	Text string
	URL  string
}

// Page is an abstract class providing a standard page structure.
type Page struct {
	Status     int
	Title      string
	StaticRoot string
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
