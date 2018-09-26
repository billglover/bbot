package messsaging

// Message is an internal representation of messages that will be sent to a
// Slack channel or user.
type Message struct {
	Text string `json:"text"`
}

// Envelope provides routing information for a message.
type Envelope struct {
	Destination Address `json:"destination"`
	TeamID      string  `json:"team_id"`
	Ephemeral   bool    `json:"ephemeral,omitempty"`
	Message     Message `json:"message"`
}

// Address indicates where the message should be sent.
type Address struct {
	Type AddressType `json:"type"`
	ID   string      `json:"id"`
}

// AddressType indicates whether the address is for a user or a channel.
type AddressType int

// The list of possible message destinations.
const (
	UnknownDestination AddressType = iota
	UserDestination
	ChannelDestination
)

func (t AddressType) String() string {
	names := []string{
		"Unknown",
		"User",
		"Channel",
	}
	return names[int(t)]
}
