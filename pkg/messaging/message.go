package messaging

// Message is an internal representation of messages that will be sent to a
// Slack channel or user.
type Message struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachment"`
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

// Attachment is an attachment to a message
type Attachment struct {
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	TitleLink   string  `json:"title_link,omitempty"`
	Fields      []Field `json:"fields,omitempty"`
}

// Field is a field in an Attachment
type Field struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}
