package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	ss "github.com/billglover/bbot/secrets"
	bot "github.com/billglover/bbot/slackbot"
)

var slackClientID string
var slackClientSecret string
var slackSigningSecret string

type handler func(ctx context.Context, req bot.Request) (bot.Response, error)

func requestHandler(q *sqs.SQS, qName string) handler {
	handlerFunc := func(ctx context.Context, req bot.Request) (bot.Response, error) {
		if bot.ValidateRequest(req, slackSigningSecret) == false {
			return bot.ErrorResponse("invalid request, check request signature", http.StatusBadRequest)
		}

		switch req.PathParameters["type"] {

		case "event":
			fmt.Println("event request")
			return bot.ErrorResponse("events API not yet implemented", http.StatusNotImplemented)

		case "command":
			fmt.Println("command request")
			return bot.ErrorResponse("command API not yet implemented", http.StatusNotImplemented)

		case "action":
			err := sendAction(ctx, req, q, qName)
			if err != nil {
				fmt.Println("WARN:", err)
				return bot.ErrorResponse("unable to enqueue action for procesing", http.StatusBadRequest)
			}

		default:
			return bot.ErrorResponse("invalid request, check endpoint type", http.StatusNotFound)
		}

		resp := bot.Response{
			StatusCode: http.StatusAccepted,
			Headers:    map[string]string{"Content-Type": "application/json"}}
		return resp, nil
	}
	return handlerFunc
}

func sendAction(ctx context.Context, req bot.Request, q *sqs.SQS, qName string) error {
	messageAction, err := bot.ParseAction(req.Body)
	if err != nil {
		return err
	}

	body, err := json.Marshal(messageAction)
	if err != nil {
		return err
	}

	delay := aws.Int64(0)
	attributes := make(map[string]*sqs.MessageAttributeValue)

	attributes["Action"] = &sqs.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: aws.String(messageAction.CallbackID),
	}

	attributes["Team"] = &sqs.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: aws.String(messageAction.Team.ID),
	}

	msgInput := &sqs.SendMessageInput{
		DelaySeconds:      delay,
		MessageAttributes: attributes,
		MessageBody:       aws.String(string(body)),
		QueueUrl:          &qName,
	}

	resp, err := q.SendMessage(msgInput)
	if err != nil {
		return err
	}

	fmt.Println("INFO: enqueued action with ID", resp.MessageId)

	return nil
}

func main() {

	// retrieve secrets from the AWS parameter store
	secrets, err := ss.GetSecrets([]string{
		"/bbot/env/SLACK_CLIENT_ID",
		"/bbot/env/SLACK_CLIENT_SECRET",
		"/bbot/env/SLACK_SIGNING_SECRET",
	})
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}

	slackClientID = secrets["/bbot/env/SLACK_CLIENT_ID"]
	slackClientSecret = secrets["/bbot/env/SLACK_CLIENT_SECRET"]
	slackSigningSecret = secrets["/bbot/env/SLACK_SIGNING_SECRET"]

	// read environment configuration
	qName := os.Getenv("SQS_QUEUE_FLAGMESSAGE")
	if qName == "" {
		fmt.Println("ERROR: no queue name specified, check environemnt variables for SQS_QUEUE_FLAGMESSAGE")
		os.Exit(1)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	qSvc := sqs.New(sess)

	lambda.Start(requestHandler(qSvc, qName))
}
