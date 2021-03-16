package move

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/sjansen/magnet/internal/aws"
	"github.com/sjansen/magnet/internal/config"
)

// Mover copies objects the inbox to start the review process.
type Mover struct {
	client *s3.Client
	config *config.Move
}

// Event is used to route lambda events based on the top-level field.
type Event struct {
	Promote  struct{}               // TODO
	S3Events []events.S3EventRecord `json:"Records"`
}

// New creates a new Mover.
func New() (*Mover, error) {
	ctx := context.Background()

	fmt.Println("Loading config...")
	cfg, err := config.LoadMoveConfig(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println("Preparing AWS clients...")
	aws, err := aws.New(ctx)
	if err != nil {
		return nil, err
	}

	return &Mover{
		client: aws.NewS3Client(),
		config: cfg,
	}, nil
}

// HandleEvent handles Lambda events.
func (m *Mover) HandleEvent(ctx context.Context, event Event) error {
	for _, e := range event.S3Events {
		err := m.move(ctx, e.S3.Bucket.Name, e.S3.Object.Key)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	return nil
}

// StartLambdaHandler waits for and processes events from AWS Lambda.
func (m *Mover) StartLambdaHandler() {
	lambda.Start(m.HandleEvent)
}
