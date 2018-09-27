package slack

import (
	"fmt"

	"github.com/billglover/bbot/pkg/messaging"
	api "github.com/nlopes/slack"
	"github.com/pkg/errors"
)

// Workspace represents a Slack workspace.
type Workspace struct {
	client       *api.Client
	botToken     string
	botUserToken string
	botUser      string
}

// New returns a Workspace. It requires a botToken and a botUserToken to allow
// it to perform operations in the Workspace. If either of these are empty, an
// error is returned.
func New(botToken, botUserToken, botUser string) (*Workspace, error) {
	w := new(Workspace)
	w.botToken = botToken
	w.botUserToken = botUserToken
	w.botUser = botUser
	if w.botToken == "" || w.botUserToken == "" || w.botUser == "" {
		return w, errors.New("botToken, botUserToken and botUser must be provided")
	}
	w.client = api.New(botToken)
	return w, nil
}

// SendMessage sends a message to Slack.
func (w *Workspace) SendMessage(e messaging.Envelope) error {

	// TODO:
	// - we need to know the channel for ephemeral messages

	// Searching for the admin channel requires the Bot User token rather than the Bot token.
	adminGroup := ""
	userClient := api.New(w.botUserToken)
	grps, err := userClient.GetGroups(true)
	if err != nil {
		fmt.Println("WARN: unable to retrieve list of groups:", err)
		return nil
	}
	fmt.Printf("INFO: groups: %+v\n", grps)
	for _, g := range grps {
		fmt.Println("\t", g.NameNormalized)
		if g.NameNormalized == "admins" {
			adminGroup = g.ID
		}
	}
	fmt.Println("INFO: adminGroup:", adminGroup)

	if e.Destination.Type == messaging.ChannelDestination {
		msgParams := api.PostMessageParameters{
			Username: w.botUser,
			AsUser:   true,
			Markdown: true,
		}

		fmt.Println("INFO: msgParams", msgParams)

		ch, ts, err := w.client.PostMessage(adminGroup, e.Message.Text, msgParams)
		if err != nil {
			fmt.Printf("INFO: %+v\n", e)
			return errors.Wrap(err, "unable to send message to channel")
		}
		fmt.Printf("INFO: message posted in channel %s at %s", ch, ts)
	}

	return nil
}
