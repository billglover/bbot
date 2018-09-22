package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/billglover/bbot/pkg/flagging"
	"github.com/pkg/errors"
)

func main() {
	lambda.Start(Handler)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, evt events.SQSEvent) error {
	for _, msg := range evt.Records {
		m := flagging.MessageAction{}
		err := json.Unmarshal([]byte(msg.Body), &m)
		if err != nil {
			return errors.Wrap(err, "unable to parse message action")
		}

		switch m.CallbackID {

		case "flagMessage":
			svc := flagging.NewSvc()
			svc.Handle(m)

		default:
			fmt.Println("WARN: unknown messageAction:", m.CallbackID)
		}

	}
	return nil
}
