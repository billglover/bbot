package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
)

var slackClientID string
var slackClientSecret string
var slackSigningSecret string

// Response is of type APIGatewayProxyResponse
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {

	fmt.Println("slackClientID:", slackClientID[:4])
	fmt.Println("slackClientSecret:", slackClientSecret[:4])
	fmt.Println("slackSigningSecret:", slackSigningSecret[:4])

	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"message": "Go Serverless v1.0! Your function executed successfully!",
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "acceptRequest-handler",
		},
	}

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

	lambda.Start(Handler)
}

func getSecrets(keys []string) (map[string]string, error) {
	svc := ssm.New(session.New())

	paramsIn := ssm.GetParametersInput{
		Names:          aws.StringSlice(keys),
		WithDecryption: aws.Bool(true),
	}

	paramsOut, err := svc.GetParameters(&paramsIn)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get parameters from AWS parameter store")
	}

	secrets := make(map[string]string, len(paramsOut.Parameters))
	for _, p := range paramsOut.Parameters {
		secrets[*p.Name] = *p.Value
	}

	return secrets, nil
}
