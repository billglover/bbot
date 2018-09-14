package main

import (
	"context"
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

		case "command":
			fmt.Println("command request")

		case "action":
			fmt.Println("action request")
			err := handleAction(ctx, req, q, qName)
			if err != nil {
				fmt.Println("WARN:", err)
				bot.ErrorResponse("unable to handle message action", http.StatusBadRequest)
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

func handleAction(ctx context.Context, req bot.Request, q *sqs.SQS, qName string) error {
	ma, err := bot.ParseAction(req.Body)
	if err != nil {
		return err
	}
	fmt.Println(ma)

	delay := aws.Int64(0)
	msgmap := make(map[string]*sqs.MessageAttributeValue)
	body := aws.String("Test Message Body")

	msgmap["Type"] = &sqs.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: aws.String(ma.CallbackID),
	}

	res, err := q.SendMessage(&sqs.SendMessageInput{DelaySeconds: delay, MessageAttributes: msgmap, MessageBody: body, QueueUrl: &qName})
	if err != nil {
		return err
	}
	fmt.Println("Success", *res.MessageId)

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
