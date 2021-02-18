package pages

import (
	"strings"
)

type BrowseHeader struct {
	Head []Href
	Tail string
}

type BrowsePage struct {
	Page
	Key string
}

func (p *BrowsePage) Header() *BrowseHeader {
	parts := strings.Split(p.Key, "/")
	dirs := len(parts) - 1
	hrefs := make([]Href, 0, len(parts))

	stop := dirs - 1
	sb := strings.Builder{}
	sb.Grow(3 * stop)

	for i := 0; i < stop; i++ {
		for j := i; j < stop; j++ {
			sb.WriteString("../")
		}
		hrefs = append(hrefs, Href{
			Text: parts[i],
			URL:  sb.String(),
		})
		sb.Reset()
	}

	var tail string
	if strings.HasSuffix(p.Key, "/") {
		sb.WriteString(parts[stop])
		sb.WriteRune('/')
		tail = sb.String()
	} else {
		hrefs = append(hrefs, Href{
			Text: parts[stop],
			URL:  ".",
		})
		tail = parts[dirs]
	}

	return &BrowseHeader{
		Head: hrefs,
		Tail: tail,
	}
}
