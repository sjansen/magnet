package move

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	// TODO S3 svc
}

// HandleEvent handles Lambda events.
func HandleEvent(ctx context.Context, events events.S3Event) error {
	for i, event := range events.Records {
		fmt.Printf("%d: %#v\n", i, event)
		fmt.Printf("%#v\n", event.S3.Bucket)
		fmt.Printf("%#v\n", event.S3.Object)
	}
	return nil
}

// StartLambdaHandler waits for and processes events from AWS Lambda.
func StartLambdaHandler() {
	lambda.Start(HandleEvent)
}
