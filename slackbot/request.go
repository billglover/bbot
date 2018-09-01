package slackbot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Request is of type APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

// ValidateRequest confirms that a request is signed with the expected
// signature for a given secret. This allows you to ensure that the request
// originated from Slack. It takes a Request and a secret and returns a bool
// indicating whether the request signature was valid. All errors are handled
// and result in the signature being returned as invalid.
//
// https://api.slack.com/docs/verifying-requests-from-slack
func ValidateRequest(req Request, secret string) bool {
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
