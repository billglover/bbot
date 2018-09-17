package main

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/billglover/bbot/queue"
	bot "github.com/billglover/bbot/slackbot"
)

var tcs = []struct {
	Name     string
	Request  bot.Request
	Response bot.Response
}{
	{
		Name: "valid event request signature",
		Request: bot.Request{
			PathParameters: map[string]string{"type": "event"},
			Body:           "token=xyzz0WbapA4vBCDEFasx0q6G&team_id=T1DC2JH3J&team_domain=testteamnow&channel_id=G8PSS9T3V&channel_name=foobar&user_id=U2CERLKJA&user_name=roadrunner&command=%2Fwebhook-collect&text=&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT1DC2JH3J%2F397700885554%2F96rGlfmibIGlgcZRskXaIFfN&trigger_id=398738663015.47445629121.803a0bc887a14d10d2c447fce8b6703c",
			HTTPMethod:     http.MethodPost,
			Headers: map[string]string{
				"X-Slack-Request-Timestamp": "1531420618",
				"X-Slack-Signature":         "v0=a2114d57b48eac39b9ad189dd8316235a7b4a8d21a10bd27519666489c69b503"},
		},
		Response: bot.Response{
			StatusCode:      http.StatusNotImplemented,
			Body:            `{"status":"501","message":"events API not yet implemented"}`,
			Headers:         map[string]string{"Content-Type": "application/json"},
			IsBase64Encoded: false,
		},
	},
	{
		Name: "valid action request signature",
		Request: bot.Request{
			PathParameters: map[string]string{"type": "action"},
			Body:           "payload=%7B%22type%22%3A%22message_action%22%2C%22token%22%3A%22w06XkgVo2IlRaQRamizypPQl%22%2C%22action_ts%22%3A%221537222293.313519%22%2C%22team%22%3A%7B%22id%22%3A%22TBLG57ECT%22%2C%22domain%22%3A%22buddybotdev%22%7D%2C%22user%22%3A%7B%22id%22%3A%22UBLKAG9K4%22%2C%22name%22%3A%22bill%22%7D%2C%22channel%22%3A%7B%22id%22%3A%22CBLPRTX3P%22%2C%22name%22%3A%22general%22%7D%2C%22callback_id%22%3A%22flagMessage%22%2C%22trigger_id%22%3A%22437591536869.394549252435.702a7192c8705ab95f942133862968e4%22%2C%22message_ts%22%3A%221535813905.000100%22%2C%22message%22%3A%7B%22type%22%3A%22message%22%2C%22user%22%3A%22UBLKAG9K4%22%2C%22text%22%3A%22hello%22%2C%22client_msg_id%22%3A%2254360322-36c9-4cd8-a7ab-13b5322362a6%22%2C%22ts%22%3A%221535813905.000100%22%7D%2C%22response_url%22%3A%22https%3A%5C%2F%5C%2Fhooks.slack.com%5C%2Fapp%5C%2FTBLG57ECT%5C%2F437591536901%5C%2Fx38f3pfe0HPyw8QovCcAA0mU%22%7D",
			HTTPMethod:     http.MethodPost,
			Headers: map[string]string{
				"X-Slack-Request-Timestamp": "1531420618",
				"X-Slack-Signature":         "v0=8042e798f4806956fe1e91674837a4228dc7020f01f4bb358d19b025ad081d93"},
		},
		Response: bot.Response{
			StatusCode:      http.StatusAccepted,
			Body:            `{"status":"202","message":""}`,
			Headers:         map[string]string{"Content-Type": "application/json"},
			IsBase64Encoded: false,
		},
	},
	{
		Name: "invalid request signature",
		Request: bot.Request{
			Body:       "token=xyzz0WbapA4vBCDEFasx0q6G&team_id=T1DC2JH3J&team_domain=testteamnow&channel_id=G8PSS9T3V&channel_name=foobar&user_id=U2CERLKJA&user_name=roadrunner&command=%2Fwebhook-collect&text=&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT1DC2JH3J%2F397700885554%2F96rGlfmibIGlgcZRskXaIFfN&trigger_id=398738663015.47445629121.803a0bc887a14d10d2c447fce8b6703c",
			HTTPMethod: http.MethodPost,
			Headers: map[string]string{
				"X-Slack-Request-Timestamp": "1531420618",
				"X-Slack-Signature":         "v0=b2114d57b48eac39b9ad189dd8316235a7b4a8d21a10bd27519666489c69b503"},
		},
		Response: bot.Response{
			StatusCode:      http.StatusBadRequest,
			Body:            `{"status":"400","message":"invalid request, check request signature"}`,
			Headers:         map[string]string{"Content-Type": "application/json"},
			IsBase64Encoded: false,
		},
	},
	{
		Name: "invalid endpoint type",
		Request: bot.Request{
			PathParameters: map[string]string{"type": "unknown"},
			Body:           "token=xyzz0WbapA4vBCDEFasx0q6G&team_id=T1DC2JH3J&team_domain=testteamnow&channel_id=G8PSS9T3V&channel_name=foobar&user_id=U2CERLKJA&user_name=roadrunner&command=%2Fwebhook-collect&text=&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT1DC2JH3J%2F397700885554%2F96rGlfmibIGlgcZRskXaIFfN&trigger_id=398738663015.47445629121.803a0bc887a14d10d2c447fce8b6703c",
			HTTPMethod:     http.MethodPost,
			Headers: map[string]string{
				"X-Slack-Request-Timestamp": "1531420618",
				"X-Slack-Signature":         "v0=a2114d57b48eac39b9ad189dd8316235a7b4a8d21a10bd27519666489c69b503"},
		},
		Response: bot.Response{
			StatusCode:      http.StatusNotFound,
			Body:            `{"status":"404","message":"invalid request, check endpoint type"}`,
			Headers:         map[string]string{"Content-Type": "application/json"},
			IsBase64Encoded: false,
		},
	},
}

func TestHandler(t *testing.T) {

	slackSigningSecret = "8f742231b10e8888abcd99yyyzzz85a5"

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			qSuccess, _ := NewSuccessQueue("")
			handler := requestHandler(qSuccess)
			resp, err := handler(ctx, tc.Request)

			if err != nil {
				t.Error(err)
			}

			if got, want := resp.StatusCode, tc.Response.StatusCode; got != want {
				t.Errorf("unexpected StatusCode: got %d, want %d", got, want)
			}

			if got, want := resp.Body, tc.Response.Body; got != want {
				t.Errorf("unexpected Body:\n  got : %s\n  want: %s", got, want)
			}

			if got, want := resp.Headers, tc.Response.Headers; reflect.DeepEqual(got, want) == false {
				t.Errorf("unexpected headers:\n  got : %+v\n  want: %+v", got, want)
			}

		})
	}

}

// TestSuccessQueue implements the Queue interface
type TestSuccessQueue struct{}

// Queue takes message headers and a body and places a message on the queue.
// It always returns without error.
func (q *TestSuccessQueue) Queue(h queue.Headers, b queue.Body) error {
	return nil
}

// NewSuccessQueue returns a Queuer that successfully queues messages.
func NewSuccessQueue(name string) (*TestSuccessQueue, error) {
	q := new(TestSuccessQueue)
	return q, nil
}

// TestFailQueue implements the Queue interface
type TestFailQueue struct{}

// Queue takes message headers and a body and places a message on the queue.
// It always returns an error.
func (q *TestFailQueue) Queue(h queue.Headers, b queue.Body) error {
	return errors.New("failed to queue message")
}

// NewFailQueue returns a Queuer that fails to successfully queue messages.
func NewFailQueue(name string) (*TestFailQueue, error) {
	q := new(TestFailQueue)
	return q, nil
}
