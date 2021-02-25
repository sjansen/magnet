package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/oklog/ulid/v2"

	"github.com/sjansen/magnet/internal/config"
	"github.com/sjansen/magnet/internal/pages"
	"github.com/sjansen/magnet/internal/util/s3form"
)

// Uploader can be used to add objects to a bucket.
type Uploader struct {
	base string
	bkt  string
	svc  *s3.S3
}

// NewUploader creates a new object uploader.
func NewUploader(base string, cfg *config.Config, svc *s3.S3) *Uploader {
	return &Uploader{
		bkt: cfg.Bucket,
		svc: svc,
	}
}

// Handler can be used to add objects to a bucket.
func (up *Uploader) Handler(w http.ResponseWriter, r *http.Request) {
	creds, err := up.svc.Client.Config.Credentials.Get()
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

	region := up.svc.Client.SigningRegion
	bucket := up.bkt
	form, err := s3form.New(region, bucket).
		Prefix(prefix).
		UseAccelerateEndpoint().
		Sign(credentials, expiration)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pages.WriteResponse(w, &pages.UploadPage{
		Form: form,
	})
}
