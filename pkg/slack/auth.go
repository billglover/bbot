package slack

// AuthResponse represents the response we receive from Slack when
// a user adds our app to their workspace.
type AuthResponse struct {
	Ok          bool   `json:"ok"`
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	UserID      string `json:"user_id"`
	TeamName    string `json:"team_name"`
	TeamID      string `json:"team_id"`
	Bot         struct {
		BotUserID      string `json:"bot_user_id"`
		BotAccessToken string `json:"bot_access_token"`
	} `json:"bot"`
}
