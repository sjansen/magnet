package config

import "net/url"

// URL represents a parsed URL
type URL struct {
	url.URL
}

// Unmarshal converts an environment variable string to a URL
func (u *URL) Unmarshal(s string) error {
	return u.URL.UnmarshalBinary([]byte(s))
}
