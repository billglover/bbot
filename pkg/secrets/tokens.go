package secrets

import "github.com/billglover/bbot/pkg/storage"

// AuthRecord represents the access token we store in DynamoDB for
// every authenticated workspace.
type AuthRecord struct {
	UID            string `json:"uid"`
	AccessToken    string `json:"access_token"`
	Scope          string `json:"scope"`
	UserID         string `json:"user_id"`
	TeamName       string `json:"team_name"`
	TeamID         string `json:"team_id"`
	BotUserID      string `json:"bot_user_id"`
	BotAccessToken string `json:"bot_access_token"`
	CoCURL         string `json:"code_of_conduct_URL"`
	AdminChannel   string `json:"admin_channel"`
}

// GetTeamTokens takes a Team ID and returns an AuthRecord containing the
// access tokens for the team. It returns an error if unable to retrieve
// the tokens.
func GetTeamTokens(db *storage.DynamoDB, teamID string) (AuthRecord, error) {
	record := AuthRecord{}
	err := db.Retrieve("uid", teamID, &record)
	return record, err
}

// SaveTeamTokens takes an AuthRecord containing the access tokens for a team.
// It returns an error if unable to store the tokens in the database.
func SaveTeamTokens(db *storage.DynamoDB, teamTokens AuthRecord) error {
	err := db.Save(teamTokens)
	return err
}
