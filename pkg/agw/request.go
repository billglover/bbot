package agw

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Request is of type APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

// IsValid takes a signing key and returns true if the request signature is
// valid. It returns false in all other cases.
func (r *Request) IsValid(key string) bool {
	if r.HTTPMethod != http.MethodPost {
		return false
	}

	ts, ok := r.Headers["X-Slack-Request-Timestamp"]
	if ok == false {
		return false
	}

	sig, ok := r.Headers["X-Slack-Signature"]
	if ok == false {
		return false
	}

	if checkHMAC(r.Body, ts, sig, key) == false {
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
