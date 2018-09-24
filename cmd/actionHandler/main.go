package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/billglover/bbot/cmd/actionHandler/router"
	"github.com/billglover/bbot/pkg/secrets"
)

func main() {

	// retrieve secrets from the AWS parameter store
	s, err := secrets.GetSecrets([]string{
		"/bbot/env/SLACK_SIGNING_SECRET",
	})
	if err != nil {
		fmt.Println("ERROR: unable to retrieve secrets from parameter store:", err)
		os.Exit(1)
	}

	// get queue names from env vars
	flagMessageQ := os.Getenv("SQS_QUEUE_FLAGMESSAGE")
	if flagMessageQ == "" {
		fmt.Println("ERROR: SQS_QUEUE_FLAGMESSAGE environment variable not set")
		os.Exit(1)
	}

	// configure the router to handle flagMessage requests
	r, err := router.New(router.SigningSecret(s["/bbot/env/SLACK_SIGNING_SECRET"]))
	err = r.RegisterRoute("flagMessage", flagMessageQ)
	if err != nil {
		fmt.Println("ERROR: unable to register queue:", err)
		os.Exit(1)
	}

	lambda.Start(r.Route)
}
