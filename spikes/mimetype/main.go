package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gabriel-vasile/mimetype"
)

type Object struct {
	Bucket string
	Key    string

	ETag      *string
	MD5       string
	MimeType  string
	Extension string
	Metadata  map[string]string
}

func main() {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	svc := s3.NewFromConfig(cfg)

	url, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	src, err := getObject(svc, url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		url, err := url.Parse(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		base := path.Base(src.Key)
		base = base[:len(base)-len(path.Ext(base))]
		key := path.Join(url.Path[1:], base, src.MD5+src.Extension)
		dst := &Object{
			Bucket: url.Host,
			Key:    key,

			ETag:     src.ETag,
			MimeType: src.MimeType,
			Metadata: src.Metadata,
		}

		err = copyObject(svc, src, dst)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func copyObject(svc *s3.Client, src, dst *Object) error {
	fmt.Println(dst.Bucket, dst.Key)
	ctx := context.Background()
	resp, err := svc.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket: aws.String(dst.Bucket),
		Key:    aws.String(dst.Key),

		CopySource:        aws.String(path.Join(src.Bucket, src.Key)),
		CopySourceIfMatch: src.ETag,

		ContentType:       aws.String(dst.MimeType),
		MetadataDirective: types.MetadataDirectiveReplace,
		Metadata:          dst.Metadata,
	})
	fmt.Println("  ETag:", aws.ToString(resp.CopyObjectResult.ETag))
	fmt.Println("  Modified:", aws.ToTime(resp.CopyObjectResult.LastModified))
	return err
}

func getObject(svc *s3.Client, url *url.URL) (*Object, error) {
	bucket := url.Host
	key := url.Path[1:]
	fmt.Println(bucket, key)

	ctx := context.Background()
	resp, err := svc.GetObject(ctx, &s3.GetObjectInput{
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

	fmt.Println("  content-disposition:", aws.ToString(resp.ContentDisposition))
	fmt.Println("  content-encoding:", aws.ToString(resp.ContentEncoding))
	fmt.Println("  content-language:", aws.ToString(resp.ContentLanguage))
	fmt.Println("  content-type:", aws.ToString(resp.ContentType))
	fmt.Println("  MD5:", md5)
	fmt.Println("  Mimetype:", mime.String(), mime.Extension())
	fmt.Println("  Metadata:")
	for k, v := range resp.Metadata {
		fmt.Printf("    %s: %q\n", k, v)
	}
	fmt.Println("  Modified:", aws.ToTime(resp.LastModified))

	return &Object{
		Bucket: bucket,
		Key:    key,

		ETag:      resp.ETag,
		MD5:       md5,
		MimeType:  mime.String(),
		Extension: mime.Extension(),
		Metadata:  resp.Metadata,
	}, nil
}
