package pages

import (
	"fmt"
	"html/template"
	"net/http"
)

var tmpls *template.Template

func init() {
	fmt.Println("Parsing templates...")
	tmpls = template.Must(template.ParseGlob("templates/*.html"))
}

type Href struct {
	Text string
	URL  string
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
