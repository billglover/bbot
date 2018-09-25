package slack

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
