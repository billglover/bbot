package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
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

// Request is of type APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

// Handler is our lambda handler invoked by the `lambda.Start` function call
func handler(ctx context.Context, req Request) (Response, error) {

	if validateRequest(req, slackClientSecret) == false {
		resp := Response{StatusCode: http.StatusBadRequest}
		return resp, nil
	}

	resp := Response{StatusCode: http.StatusAccepted}
	return resp, nil
}

func validateRequest(req Request, secret string) bool {
	if req.HTTPMethod != http.MethodPost {
		return false
	}

	ts, ok := req.Headers["X-Slack-Request-Timestamp"]
	if ok == false {
		return false
	}

	sig, ok := req.Headers["X-Slack-Signature"]
	if ok == false {
		return false
	}

	if checkHMAC(req.Body, ts, sig, secret) == false {
		return false
	}

	return true
}

// CheckHMAC reports whether msgHMAC is a valid HMAC tag for msg.
func checkHMAC(body, timestamp, msgHMAC, key string) bool {
	msgHMAC = msgHMAC[3:]
	msg := "v0:" + timestamp + ":" + body
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write([]byte(msg))

	expectedKey := hash.Sum(nil)
	actualKey, _ := hex.DecodeString(msgHMAC)
	return hmac.Equal(expectedKey, actualKey)
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
