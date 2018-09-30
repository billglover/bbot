package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/billglover/bbot/pkg/messaging"
	"github.com/billglover/bbot/pkg/storage"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/billglover/bbot/pkg/queue"
	"github.com/billglover/bbot/pkg/secrets"
	"github.com/billglover/bbot/pkg/slack"
	"github.com/pkg/errors"
)

var (
	sendMessageQ string
	region       string
	authTable    string
)

func main() {

	// When we receive a message action from Slack we send out a number of
	// messages on Slack. We send these by placing messages on a queue for
	// processing. The location of this queue is stored as an environment
	// variable.
	sendMessageQ = os.Getenv("SQS_QUEUE_SENDMESSAGE")
	if sendMessageQ == "" {
		fmt.Println("ERROR: SQS_QUEUE_SENDMESSAGE environment variable not set")
		os.Exit(1)
	}

	// In order to retrieve values from the data store we need to know where
	// the database is located. The AWS Region and DynamoDB table name are stored
	// in environment variables. If these are not set the application is unable
	// to function and so we terminate.
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

	// We tell AWS Lambda to start handling incoming message actions using our
	// handler function.
	lambda.Start(handler)
}

// Handler reads messages off the messageAction queue, unmarshals them and
// passes them to the FlagMessage function.
//
// If an error is returned the message remains on the queue for future
// processing. Without additional error handling configuration on the
// queues this can lead to infinite loops. For now, we don't return errors
// opting to log them instead.
func handler(ctx context.Context, evt queue.SQSEvent) error {
	for _, msg := range evt.Records {

		m := slack.MessageAction{}
		err := json.Unmarshal([]byte(msg.Body), &m)
		if err != nil {
			fmt.Println("ERROR: unable to parse message action:", err)
			return nil
		}

		if err := FlagMessage(m); err != nil {
			fmt.Println("ERROR: unable to flag message:", err)
			return nil
		}

	}
	return nil
}

// FlagMessage takes a message action and flags the associated message for a
// potential Code of conduct violation. It notifies the reporter, author of
// the original message and the admins channel.
//
// Each of these actions are triggered independently. However, if one or more
// actions generates an error, an error is returned to the caller.
func FlagMessage(m slack.MessageAction) error {

	// Get the outbound queue for Slack messages
	q, err := queue.NewSQSQueue(sendMessageQ)
	if err != nil {
		return errors.Wrap(err, "unable to determine sendMessage queue")
	}

	// Send a message to the reporter to let them know their request has
	// been received. Don't immediately return on error.
	msg := msgForReporter(m)
	h := queue.Headers{"Team": msg.Destination.TeamID}
	errReporter := q.Queue(h, msg)
	if errReporter != nil {
		fmt.Println("ERROR: unable to notify reporting user:", errReporter)
	}

	// Send a message to the author to let them know one of their messages has
	// been flagged. Don't immediately return on error.
	msg = msgForAuthor(m)
	h = queue.Headers{"Team": msg.Destination.TeamID}
	errAuthor := q.Queue(h, msg)
	if errAuthor != nil {
		fmt.Println("ERROR: unable to notify author:", errAuthor)
	}

	// Query slack to find the admins channel so that we can notify the admins
	// that a message has been flagged.
	adminChan, errAdmin := getAdminChannel(m.Team.ID)
	if errAdmin != nil {
		fmt.Println("ERROR: unable to notify admins:", errAdmin)
		return errors.New("there were issues notifying all parties")
	}

	msg = msgForAdmins(m, adminChan)
	h = queue.Headers{"Team": msg.Destination.TeamID}
	errAdmin = q.Queue(h, msg)
	if errAdmin != nil {
		fmt.Println("ERROR: unable to notify admins:", errAdmin)
		return errors.New("there were issues notifying all parties")
	}

	return nil
}

// msgForReporter takes a message action and constructs a message that will be
// sent to the user who reported the message.
func msgForReporter(report slack.MessageAction) messaging.Envelope {
	e := messaging.Envelope{
		Destination: messaging.Address{
			TeamID:    report.Team.ID,
			ChannelID: report.Channel.ID,
			UserID:    report.User.ID,
		},

		Message: messaging.Message{
			Text: "Thanks, we are looking into your flagged message. We may be in touch for more detail.",
		},
		Ephemeral: true,
	}
	return e
}

// msgForAuthor takes a message action and constructs a message that will be
// sent to the user who originally authored the message.
func msgForAuthor(report slack.MessageAction) messaging.Envelope {
	e := messaging.Envelope{
		Destination: messaging.Address{
			TeamID:    report.Team.ID,
			ChannelID: report.Channel.ID,
			UserID:    report.Message.UserID,
		},
		Message: messaging.Message{
			Text: `One of your recent messages was flagged for a potential Code of Conduct issue.

With this in mind, we encourage you to re-read your recent messages and consider re-phrasing anything that may be interpreted in a way that is in breach of our Code of Conduct.

It is acknowledged that occasionally we all post messages in good faith that are subsequently read in a way that was never intended. Language can be nuanced, tone difficult to guage and cultural or personal context challenging to comprehend.

One of our admins will be in touch to help highlight the message that was flagged and answer any questions you may have.`,
		},
		Ephemeral: true,
	}
	return e
}

// msgForAdmins takes a message action and constructs a message that will be
// sent to the admins channel to allow admins to investigate the report.
func msgForAdmins(report slack.MessageAction, channel string) messaging.Envelope {
author, err := getUserName(report.Team.ID, report.Message.UserID)
if err != nil {
	fmt.Println("ERROR: unable to get author name")
	author = "unknown"
}

	e := messaging.Envelope{
		Destination: messaging.Address{
			TeamID:    report.Team.ID,
			ChannelID: channel,
		},
		Message: messaging.Message{
			Attachments: []messaging.Attachment{
				{
					Title: "Message Flagged",
					Description: "The following message has been flagged for a potential Code of Conduct violation.",
					Fields: []messaging.Field{
						{Name: "message", Value:report.Message.Text, Short: false},
						{Name: "reporter", Value: report.User.Name, Short: true},
						{Name: "author", Value: author, Short: true},
						{Name: "channel", Value: report.Channel.Name, Short: true},
					},
				},
			},
		},
		Ephemeral: false,
	}
	return e
}

// getAdminChannel takes a Slack Team ID and returns the ID of the admins channel.
// It returns an error if not found. It uses the access tokens from the data store
// to query Slack for a list of channels.
func getAdminChannel(t string) (string, error) {
	var adminChan string
	db := storage.DynamoDB{
		Region: region,
		Table:  authTable,
	}
	ar, err := secrets.GetTeamTokens(&db, t)
	if err != nil {
		return adminChan, errors.Wrap(err, "unable to fetch team tokens")
	}

	ws, err := slack.New(ar.BotAccessToken, ar.AccessToken, ar.BotUserID)
	if err != nil {
		return adminChan, errors.Wrap(err, "unable to establish slack workspace")
	}

	adminChan, err = ws.AdminChannelID()
	if err != nil {
		return adminChan, errors.Wrap(err, "unable to locate admins channel")
	}

	if adminChan == "" {
		return adminChan, errors.New("unable to locate admins channel")
	}

	return adminChan, nil
}

func getUserName(t, id string) (string, error){
	var userName string
	db := storage.DynamoDB{
		Region: region,
		Table:  authTable,
	}
	ar, err := secrets.GetTeamTokens(&db, t)
	if err != nil {
		return userName, errors.Wrap(err, "unable to fetch team tokens")
	}

	ws, err := slack.New(ar.BotAccessToken, ar.AccessToken, ar.BotUserID)
	if err != nil {
		return userName, errors.Wrap(err, "unable to establish slack workspace")
	}

	userName, err = ws.UserName(id)
	if err != nil {
		return userName, errors.Wrap(err, "unable to get user name")
	}

	return userName, nil
}