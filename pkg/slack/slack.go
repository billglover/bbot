package slack

import "github.com/pkg/errors"

// Workspace represents a Slack workspace.
type Workspace struct {
	botToken     string
	botUserToken string
}

// New returns a Workspace. It requires a botToken and a botUserToken to allow
// it to perform operations in the Workspace. If either of these are empty, an
// error is returned.
func New(botToken, botUserToken string) (*Workspace, error) {
	w := new(Workspace)
	w.botToken = botToken
	w.botUserToken = botUserToken
	if w.botToken == "" || w.botUserToken == "" {
		return w, errors.New("both botToken and botUserToken must be provided")
	}
	return w, nil
}

// NotifyAdmins will send the message to the Admins channel. An error is returned
// if unable to send the message.
func (w *Workspace) NotifyAdmins(m Message) error {
	return errors.New("NotifyAdmins not implemented")
}

// NotifyUser will send an ephemeral message to the user. An error is returned
// if unable to send the mssage.
func (w *Workspace) NotifyUser(u User, m Message) error {
	return errors.New("NotifyUser not implemented")
}
