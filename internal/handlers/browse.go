package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dustin/go-humanize"
	"github.com/sjansen/magnet/internal/pages"
)

var icons map[string]string = map[string]string{
	"":     "/icons/generic.svg",
	".gif": "/icons/image.svg",
	".jpg": "/icons/image.svg",
	".mp3": "/icons/audio.svg",
	".mp4": "/icons/video.svg",
	".svg": "/icons/image.svg",
}

var validBrowsePrefixes = map[string]struct{}{
	"inbox": {},
	"media": {},
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
	tmp := strings.SplitN(path, "/", 2)
	_, ok := validBrowsePrefixes[tmp[0]]
	if len(tmp) < 2 || !ok {
		// TODO custom 404 page
		w.WriteHeader(http.StatusNotFound)
		return
	}
	hasFinalSlash := strings.HasSuffix(path, "/")

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

	if !hasFinalSlash && len(result.CommonPrefixes) == 1 && len(result.Contents) == 0 {
		redirect := path + "/"
		if redirect == *result.CommonPrefixes[0].Prefix {
			http.Redirect(w, r, b.base+redirect, 302)
			return
		}
	}

	if !hasFinalSlash && len(result.CommonPrefixes) == 0 && len(result.Contents) == 1 {
		object := result.Contents[0]
		if path == *object.Key {
			page := &pages.ObjectPage{
				Timestamp: object.LastModified.String(),
				Size:      humanize.Bytes(uint64(*object.Size)),
			}
			page.Key = path
			page.Title = path

			pages.WriteResponse(w, page)
			return
		}
	}

	page := &pages.PrefixPage{
		Prefixes: make([]string, 0, len(result.CommonPrefixes)),
		Objects:  make(map[string]string, len(result.Contents)),
	}
	page.Title = path
	page.Key = path
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
