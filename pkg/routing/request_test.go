package routing

import (
	"net/http"
	"testing"
)

// TestValidateRequest runs through a number of tests to confirm the
// ValidateRequest function correctly validates an inbound request using a
// request signature.
func TestValidateRequest(t *testing.T) {
	t.Run("valid signature", validateReqSuccess)
	t.Run("invalid signature", validateReqInvalidSig)
	t.Run("invalid method", validateReqInvalidMethod)
	t.Run("invalid timestamp", validateReqInvalidTimestamp)
}

// ValidateReqSuccess tests that ValidateRequest correctly returns true for a
// valid inbound request.
func validateReqSuccess(t *testing.T) {
	secret := "8f742231b10e8888abcd99yyyzzz85a5"
	req := Request{
		Body:       "token=xyzz0WbapA4vBCDEFasx0q6G&team_id=T1DC2JH3J&team_domain=testteamnow&channel_id=G8PSS9T3V&channel_name=foobar&user_id=U2CERLKJA&user_name=roadrunner&command=%2Fwebhook-collect&text=&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT1DC2JH3J%2F397700885554%2F96rGlfmibIGlgcZRskXaIFfN&trigger_id=398738663015.47445629121.803a0bc887a14d10d2c447fce8b6703c",
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"X-Slack-Request-Timestamp": "1531420618",
			"X-Slack-Signature":         "v0=a2114d57b48eac39b9ad189dd8316235a7b4a8d21a10bd27519666489c69b503"},
	}

	valid := validateRequest(req, secret)
	if valid == false {
		t.Fail()
	}
}

// ValidateReqInvalidSig tests that ValidateRequest ensures that only requests
// with valid signatures are considered valid.
func validateReqInvalidSig(t *testing.T) {
	secret := "8f742231b10e8888abcd99yyyzzz85a5"
	req := Request{
		Body:       "token=xyzz0WbapA4vBCDEFasx0q6G&team_id=T1DC2JH3J&team_domain=testteamnow&channel_id=G8PSS9T3V&channel_name=foobar&user_id=U2CERLKJA&user_name=roadrunner&command=%2Fwebhook-collect&text=&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT1DC2JH3J%2F397700885554%2F96rGlfmibIGlgcZRskXaIFfN&trigger_id=398738663015.47445629121.803a0bc887a14d10d2c447fce8b6703c",
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"X-Slack-Request-Timestamp": "1531420618",
			"X-Slack-Signature":         "v0=a2114d57b48eac39b9ad189dd8316235a7b4a8d21a10bd27519666489c69b403"},
	}

	valid := validateRequest(req, secret)
	if valid == true {
		t.Fail()
	}
}

// validateReqInvalidMethod tests that ValidateRequest ensures that only
// requests with a POST method are considered valid.
func validateReqInvalidMethod(t *testing.T) {
	secret := "8f742231b10e8888abcd99yyyzzz85a5"
	req := Request{
		Body:       "token=xyzz0WbapA4vBCDEFasx0q6G&team_id=T1DC2JH3J&team_domain=testteamnow&channel_id=G8PSS9T3V&channel_name=foobar&user_id=U2CERLKJA&user_name=roadrunner&command=%2Fwebhook-collect&text=&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT1DC2JH3J%2F397700885554%2F96rGlfmibIGlgcZRskXaIFfN&trigger_id=398738663015.47445629121.803a0bc887a14d10d2c447fce8b6703c",
		HTTPMethod: http.MethodGet,
		Headers: map[string]string{
			"X-Slack-Request-Timestamp": "1531420618",
			"X-Slack-Signature":         "v0=a2114d57b48eac39b9ad189dd8316235a7b4a8d21a10bd27519666489c69b503"},
	}

	valid := validateRequest(req, secret)
	if valid == true {
		t.Fail()
	}
}

// validateReqInvalidTimestamp tests that ValidateRequest ensures that only
// requests with a valid timestamp are considered valid.
func validateReqInvalidTimestamp(t *testing.T) {
	secret := "8f742231b10e8888abcd99yyyzzz85a5"
	req := Request{
		Body:       "token=xyzz0WbapA4vBCDEFasx0q6G&team_id=T1DC2JH3J&team_domain=testteamnow&channel_id=G8PSS9T3V&channel_name=foobar&user_id=U2CERLKJA&user_name=roadrunner&command=%2Fwebhook-collect&text=&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT1DC2JH3J%2F397700885554%2F96rGlfmibIGlgcZRskXaIFfN&trigger_id=398738663015.47445629121.803a0bc887a14d10d2c447fce8b6703c",
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"X-Slack-Request-Timestamp": "2531420618",
			"X-Slack-Signature":         "v0=a2114d57b48eac39b9ad189dd8316235a7b4a8d21a10bd27519666489c69b503"},
	}

	valid := validateRequest(req, secret)
	if valid == true {
		t.Fail()
	}
}
