package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/billglover/bbot/pkg/slack"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/billglover/bbot/pkg/messaging"
	"github.com/billglover/bbot/pkg/queue"
	"github.com/billglover/bbot/pkg/secrets"
	"github.com/billglover/bbot/pkg/storage"
)

var (
	clientID     string
	clientSecret string
	region       string
	authTable    string
)

func main() {
	// retrieve secrets from the AWS parameter store
	s, err := secrets.GetSecrets([]string{
		"/bbot/env/SLACK_CLIENT_ID",
		"/bbot/env/SLACK_CLIENT_SECRET",
	})
	if err != nil {
		fmt.Println("ERROR: unable to retrieve secrets from parameter store:", err)
		os.Exit(1)
	}
	clientID = s["/bbot/env/SLACK_CLIENT_ID"]
	clientSecret = s["/bbot/env/SLACK_CLIENT_SECRET"]

	region = os.Getenv("BUDDYBOT_REGION")
	if region == "" {
		region = os.Getenv("BUDDYBOT_REGION")
		fmt.Println("ERROR: BUDDYBOT_REGION environment variable not set")
		os.Exit(1)
	}

	authTable = os.Getenv("BUDDYBOT_AUTH_TABLE")
	if authTable == "" {
		fmt.Println("ERROR: BUDDYBOT_AUTH_TABLE environment variable not set")
		os.Exit(1)
	}

	lambda.Start(handler)
}

func handler(ctx context.Context, evt queue.SQSEvent) error {

	for _, msg := range evt.Records {
		e := messaging.Envelope{}
		err := json.Unmarshal([]byte(msg.Body), &e)
		if err != nil {
			// If we return an error from the handler, the message remains on the
			// queue for future processing. Without separate error handline (e.g.
			// a dead-letter queue) this can result in an infinite (expensive) loop.
			fmt.Println("ERROR: unable to parse envelope:", err)
			return nil
		}

		// TODO: send the message to Slack
		db := storage.DynamoDB{
			Region: region,
			Table:  authTable,
		}
		ar, err := secrets.GetTeamTokens(&db, e.TeamID)
		if err != nil {
			fmt.Println("ERROR: unable to fetch team tokens:", err)
			//return errors.Wrap(err, "unable to fetch team tokens")
			return nil
		}

		ws, err := slack.New(ar.BotAccessToken, ar.AccessToken, ar.BotUserID)
		if err != nil {
			fmt.Println("ERROR: unable to establish Slack workspace:", err)
			//return errors.Wrap(err, "unable to establish Slack workspace")
			return nil
		}

		err = ws.SendMessage(e)
		if err != nil {
			fmt.Println("ERROR: unable to send message to Slack:", err)
			//return errors.Wrap(err, "unable to send message to Slack")
			return nil
		}

		//fmt.Printf("%+v\n", e)
	}

	return nil
}
