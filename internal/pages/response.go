package pages

import (
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
