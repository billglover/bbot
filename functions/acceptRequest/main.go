package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/billglover/bbot/queue"
	ss "github.com/billglover/bbot/secrets"
	bot "github.com/billglover/bbot/slackbot"
)

var slackClientID string
var slackClientSecret string
var slackSigningSecret string

type handler func(ctx context.Context, req bot.Request) (bot.Response, error)

func requestHandler(q queue.Queuer) handler {
	handlerFunc := func(ctx context.Context, req bot.Request) (bot.Response, error) {
		if bot.ValidateRequest(req, slackSigningSecret) == false {
			return bot.ErrorResponse("invalid request, check request signature", http.StatusBadRequest)
		}

		switch req.PathParameters["type"] {

		case "event":
			return handleEvent(ctx, req, q)

		case "command":
			return handleCommand(ctx, req, q)

		case "action":
			return handleAction(ctx, req, q)

		default:
			return bot.ErrorResponse("invalid request, check endpoint type", http.StatusNotFound)
		}
	}
	
	return handlerFunc
}

func handleEvent(ctx context.Context, req bot.Request, q queue.Queuer) (bot.Response, error) {
	return bot.ErrorResponse("events API not yet implemented", http.StatusNotImplemented)
}

func handleCommand(ctx context.Context, req bot.Request, q queue.Queuer) (bot.Response, error) {
	return bot.ErrorResponse("command API not yet implemented", http.StatusNotImplemented)
}

func handleAction(ctx context.Context, req bot.Request, q queue.Queuer) (bot.Response, error) {
	action, err := bot.ParseAction(req.Body)
	if err != nil {
		fmt.Println("WARN:", err)
		return bot.ErrorResponse("unable to parse action", http.StatusBadRequest)
	}

	h := queue.Headers{
		"Team":   action.Team.ID,
		"Action": action.CallbackID,
	}

	err = q.Queue(h, action)
	if err != nil {
		fmt.Println("ERROR:", err)
		return bot.ErrorResponse("unable to queue action for further processing", http.StatusInternalServerError)
	}

	return bot.SuccessResponse("")
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

	// establish the outbound queue
	q, err := queue.NewSQSQueue(qName)
	if err != nil {
		fmt.Println("ERROR: unable to establish queue")
		os.Exit(1)
	}

	lambda.Start(requestHandler(q))
}
