package messaging

// Message is an internal representation of messages that will be sent to a
// Slack channel or user.
type Message struct {
	Text string `json:"text"`
}

// Envelope provides routing information for a message.
type Envelope struct {
	Destination Address `json:"destination"`
	Ephemeral   bool    `json:"ephemeral,omitempty"`
	Message     Message `json:"message"`
}

// Address indicates where the message should be sent.
type Address struct {
	TeamID    string `json:"team_id"`
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
}
