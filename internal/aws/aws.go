package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// NewSession creates a new aws-sdk session.
func NewSession() (*session.Session, error) {
	aws, err := session.NewSession(
		aws.NewConfig().
			WithCredentialsChainVerboseErrors(true),
	)
	return aws, err
}

func NewS3Client(aws *session.Session) *s3.S3 {
	return s3.New(aws)
}
