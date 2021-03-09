package convert

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	// TODO S3 svc
}

// HandleEvent handles Lambda events.
func HandleEvent(ctx context.Context, events json.RawMessage) error {
	fmt.Println(string(events))
	return nil
}

// StartLambdaHandler waits for and processes events from AWS Lambda.
func StartLambdaHandler() {
	lambda.Start(HandleEvent)
}
