package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/billglover/bbot/cmd/actionHandler/router"
	"github.com/billglover/bbot/pkg/secrets"
)

func main() {

	// Retrieve the Slack signing secret from the AWS parameter store. This is
	// used to ensure incoming requests orginated from Slack. If we can't
	// retrieve the certificate we terminate the program as there is nothing
	// we can do without it.
	s, err := secrets.GetSecrets([]string{"/bbot/env/SLACK_SIGNING_SECRET"})
	if err != nil {
		fmt.Println("ERROR: unable to retrieve signing secret from parameter store:", err)
		os.Exit(1)
	}
	signingSecret, ok := s["/bbot/env/SLACK_SIGNING_SECRET"]
	if ok == false {
		fmt.Println("ERROR: unable to retrieve signing secret from parameter store:", err)
		os.Exit(1)
	}

	// When we receive a message action from Slack we place it onto a queue for
	// processing. The location of these queues are stored in environment
	// variables.
	flagMessageQ := os.Getenv("SQS_QUEUE_FLAGMESSAGE")
	if flagMessageQ == "" {
		fmt.Println("ERROR: SQS_QUEUE_FLAGMESSAGE environment variable not set")
		os.Exit(1)
	}

	// The router is responsible for validating the signature on requests that
	// we receive and then identifying the action being requested before placing
	// the message action onto the appropriate queue for processing. We
	// configure the router by registering a mapping between actions and queues.
	// If we are unable to register any routes we terminate the program as
	// it offers no functionality without route mappings.
	r, err := router.New(router.SigningSecret(signingSecret))
	err = r.RegisterRoute("flagMessage", flagMessageQ)
	if err != nil {
		fmt.Println("ERROR: unable to register queue:", err)
		os.Exit(1)
	}

	// We tell AWS Lambda to start routing incoming message actions using our
	// router. The router is responsible for sending the appropriate responses
	// to all requests.
	lambda.Start(r.Route)
}
