package slackbot

import (
	"encoding/json"
	"net/url"

	"github.com/pkg/errors"
)

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

type Team struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
}

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type MessageAction struct {
	Type             string      `json:"type"`
	CallbackID       string      `json:"callback_id"`
	Team             Team        `json:"team"`
	Channel          Channel     `json:"channel"`
	User             User        `json:"user"`
	ActionTimestamp  json.Number `json:"action_ts"`
	MessageTimestamp json.Number `json:"message_ts"`
	Message          Message     `json:"message"`
	ResponseURL      string      `json:"response_url"`
	TriggerID        string      `json:"trigger_id"`
}

type Message struct {
	Msg
	SubMessage *Msg `json:"message,omitempty"`
}

type Msg struct {
	Type             string `json:"type,omitempty"`
	UserID           string `json:"user,omitempty"`
	Text             string `json:"text,omitempty"`
	Timestamp        string `json:"ts,omitempty"`
	ThreadTimestamp  string `json:"thread_ts,omitempty"`
	SubType          string `json:"subtype,omitempty"`
	Hidden           bool   `json:"hidden,omitempty"`
	DeletedTimestamp string `json:"deleted_ts,omitempty"`
	EventTimestamp   string `json:"event_ts,omitempty"`
	BotID            string `json:"bot_id,omitempty"`
	BotName          string `json:"username,omitempty"`
}
