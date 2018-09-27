package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/billglover/bbot/pkg/storage"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/billglover/bbot/pkg/queue"
	"github.com/billglover/bbot/pkg/secrets"
	"github.com/billglover/bbot/pkg/slack"
	"github.com/pkg/errors"
)

var (
	clientID     string
	clientSecret string
	sendMessageQ string
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

	// get queue names from env vars
	sendMessageQ = os.Getenv("SQS_QUEUE_SENDMESSAGE")
	if sendMessageQ == "" {
		fmt.Println("ERROR: SQS_QUEUE_SENDMESSAGE environment variable not set")
		os.Exit(1)
	}

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

// Handler is our lambda handler invoked by the `lambda.Start` function call
func handler(ctx context.Context, evt queue.SQSEvent) error {
	for _, msg := range evt.Records {

		m := slack.MessageAction{}
		err := json.Unmarshal([]byte(msg.Body), &m)
		if err != nil {
			return nil //errors.Wrap(err, "unable to parse message action")
		}

		fmt.Printf("%+v\n", m)

		if err := FlagMessage(m); err != nil {
			fmt.Println("ERROR: unable to flag message", err)
			return nil //errors.Wrap(err, "unable to flag message")
		}

	}
	return nil
}

// FlagMessage takes a message action and flags it for a potential Code of
// Conduct violation. It returns an error if unable to flag the message.
func FlagMessage(m slack.MessageAction) error {

	db := storage.DynamoDB{
		Region: region,
		Table:  authTable,
	}
	ar, err := secrets.GetTeamTokens(&db, m.Team.ID)
	if err != nil {
		return errors.Wrap(err, "unable to fetch team tokens")
	}

	// TODO: rather than send messages to Slack directly, this function should
	// place messages onto the sendMessageQueue to be handled by the msgSender
	// function.
	w, err := slack.New(ar.BotAccessToken, ar.AccessToken)
	if err != nil {
		return errors.Wrap(err, "unable to create workspace")
	}

	if err = w.NotifyAdmins(slack.Message{}); err != nil {
		return errors.Wrap(err, "unable to notify admins")
	}

	if err = w.NotifyUser(slack.User{}, slack.Message{}); err != nil {
		return errors.Wrap(err, "unable to notify reporting user")
	}

	if err = w.NotifyUser(slack.User{}, slack.Message{}); err != nil {
		return errors.Wrap(err, "unable to notify message author")
	}

	fmt.Println("INFO: FlagMessage complete")

	return nil
}
