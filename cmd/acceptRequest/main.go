package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/billglover/bbot/pkg/routing"
	ss "github.com/billglover/bbot/pkg/secrets"
)

var slackClientID string
var slackClientSecret string
var slackSigningSecret string

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
	actionQ := os.Getenv("SQS_QUEUE_ACTION")
	if actionQ == "" {
		fmt.Println("ERROR: no queue name specified, check environemnt variables for SQS_QUEUE_FLAGMESSAGE")
		os.Exit(1)
	}

	// start the router
	r, err := routing.NewRouter(actionQ, "", "")
	if err != nil {
		fmt.Println("ERROR: unable to instantiate router:", err)
		os.Exit(1)
	}

	r.ReqSecret = slackSigningSecret
	lambda.Start(r.Route)
}
