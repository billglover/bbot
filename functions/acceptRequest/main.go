package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	ss "github.com/billglover/bbot/secrets"
	bot "github.com/billglover/bbot/slackbot"
)

var slackClientID string
var slackClientSecret string
var slackSigningSecret string

// Handler is our lambda handler invoked by the `lambda.Start` function call
func handler(ctx context.Context, req bot.Request) (bot.Response, error) {

	if bot.ValidateRequest(req, slackSigningSecret) == false {
		return bot.ErrorResponse("invalid request, check request signature", http.StatusBadRequest)
	}

	resp := bot.Response{
		StatusCode: http.StatusAccepted,
		Headers:    map[string]string{"Content-Type": "application/json"}}
	return resp, nil
}

func main() {
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

	lambda.Start(handler)
}
