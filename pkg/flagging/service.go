package flagging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
)

// Event represents and AWS SQS Event
type Event events.SQSEvent

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, evt Event) error {
	for _, msg := range evt.Records {
		mA := MessageAction{}
		err := json.Unmarshal([]byte(msg.Body), &mA)
		if err != nil {
			return errors.Wrap(err, "unable to parse message action")
		}
		fmt.Printf("messageAction: %+v\n", mA)
	}
	return nil
}
