package slack

import "encoding/json"

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

// User is a Slack user.
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Message is a Slack message.
type Message struct {
	Type     string `json:"type,omitempty"`
	UserID   string `json:"user,omitempty"`
	Text     string `json:"text,omitempty"`
	Ts       string `json:"ts,omitempty"`
	ThreadTs string `json:"thread_ts,omitempty"`
	SubType  string `json:"subtype,omitempty"`
	BotID    string `json:"bot_id,omitempty"`
	BotName  string `json:"username,omitempty"`
}
