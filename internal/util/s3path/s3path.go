package s3path

import (
	"net/url"
	"strings"
)

// S3Path simplifies working with S3 paths
type S3Path struct {
	Display string
	Escaped string
}

func FromS3(key string) (*S3Path, error) {
	unescaped, err := url.QueryUnescape(key)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(unescaped, "/")
	for idx, part := range parts {
		parts[idx] = url.QueryEscape(part)
	}
	escaped := strings.Join(parts, "/")
	return &S3Path{
		Display: unescaped,
		Escaped: escaped,
	}, nil
}
