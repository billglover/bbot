package slack

import (
	"fmt"

	"github.com/billglover/bbot/pkg/messaging"
	api "github.com/nlopes/slack"
	"github.com/pkg/errors"
)

// Workspace represents a Slack workspace.
type Workspace struct {
	botClient    *api.Client
	userClient   *api.Client
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
	w.botClient = api.New(botToken)
	w.userClient = api.New(botUserToken)
	return w, nil
}

// SendMessage sends a message to Slack.
func (w *Workspace) SendMessage(e messaging.Envelope) error {

	switch {

	// Ephemeral messages are sent to an individual user
	case e.Ephemeral == true:
		if e.Destination.UserID == "" {
			return errors.New("ephemeral messages require a UserID")
		}

		msgOptsEphemeral := api.MsgOptionPostEphemeral2(e.Destination.UserID)
		msgOpts := api.MsgOptionText(e.Message.Text, true)
		ts, err := w.botClient.PostEphemeral(e.Destination.ChannelID, e.Destination.UserID, msgOptsEphemeral, msgOpts)
		if err != nil {
			return errors.Wrap(err, "failed to send ephemeral message")
		}
		fmt.Println("INFO: ephemeral emssage sent:", ts)

	// Standard messages without a UserID specified are sent to a channel
	case e.Ephemeral == false && e.Destination.UserID == "":

		msgParams := api.PostMessageParameters{
			Username: w.botUser,
			AsUser:   true,
			Markdown: true,
		}

		if e.Message.Attachments != nil {
			attachments := make([]api.Attachment, len(e.Message.Attachments))
			for i, a := range e.Message.Attachments {
				attachments[i] = api.Attachment{
					Title:     a.Title,
					TitleLink: a.TitleLink,
					Pretext:   a.Description,
				}

				// include all fields
				if a.Fields != nil {
					fields := make([]api.AttachmentField, len(a.Fields))
					for j, f := range a.Fields {
						fields[j] = api.AttachmentField{Title: f.Name, Value: f.Value, Short: f.Short}
					}
					attachments[i].Fields = fields
				}
			}
			msgParams.Attachments = attachments
		}

		ch, ts, err := w.botClient.PostMessage(e.Destination.ChannelID, e.Message.Text, msgParams)
		if err != nil {
			return errors.Wrap(err, "unable to send message to channel")
		}
		fmt.Printf("INFO: message posted in channel %s at %s", ch, ts)

		// Standard messages without a ChannelID specified ar esent to a user
	case e.Ephemeral == false && e.Destination.ChannelID == "":
		return errors.New("sending messages to an individual user is not currently supported")

	default:
		return errors.Errorf("unable to determine intended message destination: %+v", e.Destination)
	}

	return nil
}

// AdminChannelID returns the ChannelID for the admins channel in a workspace.
// It returns an error if it is unable to identify the admins channel.
func (w *Workspace) AdminChannelID() (string, error) {
	var id string

	// Note: private channels are known as Groups in Slack.
	//grps, err := w.botClient.GetGroups(true)

	// TODO: this will fail if it can't find the admins channel in the first 1000
	// results or if Slack decides to return fewer results. To fix this you need
	// to use the cursor to paginate through results.
	params := api.GetConversationsParameters{
		ExcludeArchived: "true",
		Types:           []string{"private_channel"},
		Limit:           1000,
	}
	grps, _, err := w.botClient.GetConversations(&params)
	if err != nil {
		return id, errors.Wrap(err, "unable to retrieve list of channels")
	}

	for _, g := range grps {
		if g.NameNormalized == "admins" {
			id = g.ID
		}
	}

	if id == "" {
		return id, errors.New("unable to locate 'admins' group")
	}

	return id, nil
}

// UserName takes a UserID and returns the corresponding UserName. It reutrns
// an error if it is unable to look up the user name.
func (w *Workspace) UserName(id string) (string, error) {

	user, err := w.botClient.GetUserInfo(id)
	return user.Name, err
}

// Permalink takes a message timestamp and returns the corresponding permalink.
// It returns an error if it is unable to look up the permalink.
func (w *Workspace) Permalink(ch, ts string) (string, error) {
	params := api.GetPermalinkParameters{
		Channel: ch,
		Ts:      ts,
	}

	permalink, err := w.botClient.GetPermalink(&params)
	return permalink, err
}
