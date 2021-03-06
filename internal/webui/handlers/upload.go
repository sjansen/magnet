package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/oklog/ulid/v2"

	"github.com/sjansen/magnet/internal/config"
	"github.com/sjansen/magnet/internal/util/s3form"
	"github.com/sjansen/magnet/internal/webui/pages"
)

// Uploader can be used to add objects to a bucket.
type Uploader struct {
	aws        aws.Config
	base       string
	bucket     string
	staticRoot string
}

// NewUploader creates a new object uploader.
func NewUploader(base string, cfg *config.WebUI, bucket string) *Uploader {
	return &Uploader{
		aws:        cfg.AWS.Config,
		bucket:     bucket,
		staticRoot: cfg.StaticRoot,
	}
}

// ServeHTTP can be used to add objects to a bucket.
func (u *Uploader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	creds, err := u.aws.Credentials.Retrieve(r.Context())
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	credentials := &s3form.Credentials{
		AccessKeyID:     creds.AccessKeyID,
		SecretAccessKey: creds.SecretAccessKey,
		SessionToken:    creds.SessionToken,
	}

	now := time.Now()
	expiration := now.Add(30 * time.Minute)

	// TODO use sync.Pool to reuse entropy source
	entropy := ulid.Monotonic(rand.New(rand.NewSource(now.UnixNano())), 0)
	ulid, err := ulid.New(ulid.Timestamp(now), entropy)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	prefix := "inbox/" + ulid.String() + "/"

	region := u.aws.Region
	bucket := u.bucket
	form, err := s3form.New(region, bucket).
		Prefix(prefix).
		UseAccelerateEndpoint().
		AddField("x-amz-meta-batch", "", &s3form.StartsWith{}).
		AddField("x-amz-meta-creator", "", &s3form.StartsWith{}).
		AddField("x-amz-meta-license", "", &s3form.StartsWith{}).
		AddField("x-amz-meta-source", "", &s3form.StartsWith{}).
		Sign(credentials, expiration)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p := &pages.UploadPage{
		Form: form,
	}
	p.StaticRoot = u.staticRoot
	pages.WriteResponse(w, p)
}
