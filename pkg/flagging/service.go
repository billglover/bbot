package flagging

import (
	"fmt"
)

// Svc is the flagMessage service.
type Svc struct{}

// NewSvc returns a Svc with default configuration.
func NewSvc() *Svc {
	s := new(Svc)
	return s
}

// Handle is responsible for taking action in response to a flaggedMessage
// action being received. It is responsible for handling all errors that occur
// as there is no response path to the user who originally flagged the message.
func (s *Svc) Handle(m MessageAction) {
	fmt.Println("Message Text:", m.Message.Text)

	// What do we want to do?
	// 1) Let the user know the message has been flagged
	// 2) Send a note to the administrators group
	// 3) Send a note to the author of the original message indicating their message has been flagged
}
