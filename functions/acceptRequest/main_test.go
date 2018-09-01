package main

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	bot "github.com/billglover/bbot/slackbot"
)

var tcs = []struct {
	Name     string
	Request  bot.Request
	Response bot.Response
}{
	{
		Name: "valid request signature",
		Request: bot.Request{
			Body:       "token=xyzz0WbapA4vBCDEFasx0q6G&team_id=T1DC2JH3J&team_domain=testteamnow&channel_id=G8PSS9T3V&channel_name=foobar&user_id=U2CERLKJA&user_name=roadrunner&command=%2Fwebhook-collect&text=&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT1DC2JH3J%2F397700885554%2F96rGlfmibIGlgcZRskXaIFfN&trigger_id=398738663015.47445629121.803a0bc887a14d10d2c447fce8b6703c",
			HTTPMethod: http.MethodPost,
			Headers: map[string]string{
				"X-Slack-Request-Timestamp": "1531420618",
				"X-Slack-Signature":         "v0=a2114d57b48eac39b9ad189dd8316235a7b4a8d21a10bd27519666489c69b503"},
		},
		Response: bot.Response{
			StatusCode:      http.StatusAccepted,
			Body:            "",
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
}

func TestHandler(t *testing.T) {

	slackSigningSecret = "8f742231b10e8888abcd99yyyzzz85a5"

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
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
