package slack

import (
	"encoding/json"
	"net/url"

	"github.com/pkg/errors"
)

// MessageAction is the message received from the Slack API in response
// to a user performing an action on a message.
type MessageAction struct {
	Type        string      `json:"type"`
	CallbackID  string      `json:"callback_id"`
	Team        Team        `json:"team"`
	Channel     Channel     `json:"channel"`
	User        User        `json:"user"`
	ActionTs    json.Number `json:"action_ts"`
	MessageTs   json.Number `json:"message_ts"`
	Message     Message     `json:"message"`
	ResponseURL string      `json:"response_url"`
	TriggerID   string      `json:"trigger_id"`
}

// Team is a Slack team.
type Team struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
}

// Channel is a Slack channel.
type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ParseAction parses the payload of a request, a string, and returns a MessageAction.
func ParseAction(b string) (MessageAction, error) {
	ma := MessageAction{}

	form, err := url.ParseQuery(b)
	if err != nil {
		return ma, errors.Wrap(err, "failed to parse request body")
	}

	err = json.Unmarshal([]byte(form.Get("payload")), &ma)
	if err != nil {
		return ma, errors.Wrap(err, "failed to parse request body")
	}
	return ma, err
}
