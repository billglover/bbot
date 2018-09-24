package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/billglover/bbot/pkg/queue"
	"github.com/billglover/bbot/pkg/slack"
	"github.com/pkg/errors"
)

func main() {

	// retrieve secrets from the AWS parameter store
	// s, err := secrets.GetSecrets([]string{
	// 	"/bbot/env/SLACK_CLIENT_ID",
	// 	"/bbot/env/SLACK_CLIENT_SECRET",
	// })
	// if err != nil {
	// 	fmt.Println("ERROR: unable to retrieve secrets from parameter store:", err)
	// 	os.Exit(1)
	// }

	// get queue names from env vars
	sendMessageQ := os.Getenv("SQS_QUEUE_SENDMESSAGE")
	if sendMessageQ == "" {
		fmt.Println("ERROR: SQS_QUEUE_SENDMESSAGE environment variable not set")
		os.Exit(1)
	}

	lambda.Start(handler)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func handler(ctx context.Context, evt queue.SQSEvent) error {
	for _, msg := range evt.Records {
		m := slack.MessageAction{}
		err := json.Unmarshal([]byte(msg.Body), &m)
		if err != nil {
			return errors.Wrap(err, "unable to parse message action")
		}

		fmt.Printf("%+v\n", m)

	}
	return nil
}
