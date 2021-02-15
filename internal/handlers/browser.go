package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sjansen/magnet/internal/pages"
)

var icons map[string]string = map[string]string{
	"":     "/icons/generic.svg",
	".gif": "/icons/image.svg",
	".jpg": "/icons/image.svg",
	".mp3": "/icons/audio.svg",
	".mp4": "/icons/video.svg",
}

// Browser can be used to browse the objects in a bucket.
type Browser struct {
	base string
	bkt  *string
	svc  *s3.S3
}

// NewBrowser creates a new bucket browser.
func NewBrowser(base string, bucket string, svc *s3.S3) *Browser {
	return &Browser{
		base: base,
		bkt:  aws.String(bucket),
		svc:  svc,
	}
}

// Handler can be used to browse the objects in a bucket.
func (b *Browser) Handler(w http.ResponseWriter, r *http.Request) {
	// TODO move trimming to router?
	path := strings.TrimPrefix(r.URL.Path, b.base)
	if !strings.HasPrefix(path, "media/") {
		// TODO custom 404 page
		w.WriteHeader(http.StatusNotFound)
		return
	}

	result, err := b.svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    b.bkt,
		Prefix:    aws.String(path),
		Delimiter: aws.String("/"),
		MaxKeys:   aws.Int64(100),
		// TODO ContinuationToken
	})
	if err != nil {
		fmt.Printf("Error browsing path=%q err=%s\n", path, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	page := &pages.BrowserPage{
		Prefix:   path,
		Prefixes: make([]string, 0, len(result.CommonPrefixes)),
		Objects:  make(map[string]string, len(result.Contents)),
	}
	for _, x := range result.CommonPrefixes {
		prefix := strings.TrimPrefix(*x.Prefix, path)
		page.Prefixes = append(page.Prefixes, prefix)
	}
	for _, object := range result.Contents {
		if key := strings.TrimPrefix(*object.Key, path); key != "" {
			if icon, ok := icons[strings.ToLower(filepath.Ext(key))]; ok {
				page.Objects[key] = icon
			} else {
				page.Objects[key] = icons[""]
			}
		}
	}

	pages.WriteResponse(w, page)
}
