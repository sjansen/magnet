package move

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sjansen/magnet/internal/config"
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
func StartLambdaHandler(cfg *config.Move) {
	lambda.Start(HandleEvent)
}
