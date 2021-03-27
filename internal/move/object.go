package move

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gabriel-vasile/mimetype"
)

type object struct {
	Bucket string
	Key    string

	ETag      *string
	MD5       string
	MimeType  string
	Extension string
	Metadata  map[string]string
}

func (m *Mover) inspect(ctx context.Context, bucket, key string) (*object, error) {
	resp, err := m.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	f := resp.Body
	defer f.Close()

	h := md5.New()
	r := io.TeeReader(f, h)
	mime, err := mimetype.DetectReader(r)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	md5 := fmt.Sprintf("%x", h.Sum(nil))

	fmt.Println("  MD5:", md5)
	fmt.Println("  Mimetype:", mime.String(), mime.Extension())
	fmt.Println("  content-disposition:", aws.ToString(resp.ContentDisposition))
	fmt.Println("  content-encoding:", aws.ToString(resp.ContentEncoding))
	fmt.Println("  content-language:", aws.ToString(resp.ContentLanguage))
	fmt.Println("  content-type:", aws.ToString(resp.ContentType))
	fmt.Println("  Metadata:")
	for k, v := range resp.Metadata {
		fmt.Printf("    %s: %q\n", k, v)
	}
	fmt.Println("  Modified:", aws.ToTime(resp.LastModified))

	base := strings.TrimPrefix(key, "inbox/")
	base = base[:len(base)-len(path.Ext(base))]
	if idx := strings.IndexRune(base, '/'); idx > -1 {
		base = base[idx+1:]
	}
	if prefix, ok := resp.Metadata["batch-prefix"]; ok {
		delete(resp.Metadata, "batch-prefix")
		key = path.Join("review", prefix, base, md5+mime.Extension())
	} else {
		key = path.Join("review", base, md5+mime.Extension())
	}

	return &object{
		Bucket: bucket,
		Key:    key,

		ETag:     resp.ETag,
		MD5:      md5,
		MimeType: mime.String(),
		Metadata: resp.Metadata,
	}, nil
}

func (m *Mover) move(ctx context.Context, bucket, key string) error {
	fmt.Println("Inspecting", bucket, key)
	dst, err := m.inspect(ctx, bucket, key)
	if err != nil {
		return err
	}

	fmt.Println("Creating", dst.Bucket, dst.Key)
	_, err = m.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket: aws.String(dst.Bucket),
		Key:    aws.String(dst.Key),

		CopySource:        aws.String(path.Join(bucket, key)),
		CopySourceIfMatch: dst.ETag,

		ContentType:       aws.String(dst.MimeType),
		MetadataDirective: types.MetadataDirectiveReplace,
		Metadata:          dst.Metadata,
	})
	if err != nil {
		return err
	}
	fmt.Println("SUCCESS")

	return nil
}
