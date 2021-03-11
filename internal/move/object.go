package move

import (
	"crypto/md5"
	"fmt"
	"io"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gabriel-vasile/mimetype"
)

type object struct {
	Bucket string
	Key    string

	ETag      *string
	MD5       string
	MimeType  string
	Extension string
	Metadata  map[string]*string
}

func (m *Mover) inspect(bucket, key string) (*object, error) {
	resp, err := m.client.GetObject(&s3.GetObjectInput{
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
	fmt.Println("  content-disposition:", aws.StringValue(resp.ContentDisposition))
	fmt.Println("  content-encoding:", aws.StringValue(resp.ContentEncoding))
	fmt.Println("  content-language:", aws.StringValue(resp.ContentLanguage))
	fmt.Println("  content-type:", aws.StringValue(resp.ContentType))
	fmt.Println("  Metadata:")
	for k, v := range resp.Metadata {
		fmt.Printf("    %s: %q\n", k, aws.StringValue(v))
	}
	fmt.Println("  Modified:", aws.TimeValue(resp.LastModified))

	return &object{
		Bucket: bucket,
		Key:    key,

		ETag:      resp.ETag,
		MD5:       md5,
		MimeType:  mime.String(),
		Extension: mime.Extension(),
		Metadata:  resp.Metadata,
	}, nil
}

func (m *Mover) move(bucket, key string) error {
	fmt.Println("Inspecting", bucket, key)
	src, err := m.inspect(bucket, key)
	if err != nil {
		return err
	}

	base := path.Base(src.Key)
	base = base[:len(base)-len(path.Ext(base))]
	dst := &object{
		Bucket: bucket,
		Key:    path.Join("review", base, src.MD5+src.Extension),

		ETag:     src.ETag,
		MimeType: src.MimeType,
		Metadata: src.Metadata,
	}

	fmt.Println("Creating", dst.Bucket, dst.Key)
	_, err = m.client.CopyObject(&s3.CopyObjectInput{
		Bucket: aws.String(dst.Bucket),
		Key:    aws.String(dst.Key),

		CopySource:        aws.String(path.Join(src.Bucket, src.Key)),
		CopySourceIfMatch: src.ETag,

		ContentType:       aws.String(dst.MimeType),
		MetadataDirective: aws.String("REPLACE"),
		Metadata:          dst.Metadata,
	})
	if err != nil {
		return err
	}
	fmt.Println("SUCCESS")

	return nil
}
