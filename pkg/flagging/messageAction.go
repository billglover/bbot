package flagging

type MessageAction struct {
	Type        string  `json:"type"`
	CallbackID  string  `json:"callback_id"`
	Team        Team    `json:"team"`
	Channel     Channel `json:"channel"`
	User        User    `json:"user"`
	ActionTs    float64 `json:"action_ts"`
	MessageTs   float64 `json:"message_ts"`
	Message     Message `json:"message"`
	ResponseURL string  `json:"response_url"`
	TriggerID   string  `json:"trigger_id"`
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
