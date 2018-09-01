package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	bot "github.com/billglover/bbot/slackbot"
	"github.com/pkg/errors"
)

var slackClientID string
var slackClientSecret string
var slackSigningSecret string

// Handler is our lambda handler invoked by the `lambda.Start` function call
func handler(ctx context.Context, req bot.Request) (bot.Response, error) {

	if bot.ValidateRequest(req, slackClientSecret) == false {
		resp := bot.Response{StatusCode: http.StatusBadRequest}
		return resp, nil
	}

	resp := bot.Response{StatusCode: http.StatusAccepted}
	return resp, nil
}

func main() {
	secrets, err := getSecrets([]string{
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

