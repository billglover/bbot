package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/billglover/bbot/pkg/agw"
	"github.com/billglover/bbot/pkg/secrets"
	"github.com/billglover/bbot/pkg/slack"
	"github.com/billglover/bbot/pkg/storage"
)

var (
	clientID     string
	clientSecret string
	region       string
	authTable    string
)

func main() {

	region = os.Getenv("BUDDYBOT_REGION")
	if region == "" {
		region = os.Getenv("BUDDYBOT_REGION")
		fmt.Println("ERROR: BUDDYBOT_REGION environment variable not set")
		os.Exit(1)
	}

	authTable = os.Getenv("BUDDYBOT_AUTH_TABLE")
	if authTable == "" {
		fmt.Println("ERROR: BUDDYBOT_AUTH_TABLE environment variable not set")
		os.Exit(1)
	}

	stage := os.Getenv("BUDDYBOT_STAGE")
	if stage == "" {
		fmt.Println("ERROR: BUDDYBOT_STAGE environment variable not set")
		os.Exit(1)
	}

	// retrieve secrets from the AWS parameter store
	s, err := secrets.GetSecrets([]string{
		"/bbot/" + stage + "/SLACK_CLIENT_ID",
		"/bbot/" + stage + "/SLACK_CLIENT_SECRET",
	})
	if err != nil {
		fmt.Println("ERROR: unable to retrieve secrets from parameter store:", err)
		os.Exit(1)
	}
	clientID = s["/bbot/"+stage+"/SLACK_CLIENT_ID"]
	clientSecret = s["/bbot/"+stage+"/SLACK_CLIENT_SECRET"]

	lambda.Start(handler)
}

func handler(ctx context.Context, req agw.Request) (agw.Response, error) {

	// change the temporary code for an API access token
	v := url.Values{}
	v.Set("code", req.QueryStringParameters["code"])
	v.Set("name", "https://ro9agrx7m2.execute-api.eu-west-1.amazonaws.com/dev/endpoint/auth")

	r, err := http.NewRequest(http.MethodPost, "https://slack.com/api/oauth.access", strings.NewReader(v.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.SetBasicAuth(clientID, clientSecret)
	client := http.DefaultClient
	resp, err := client.Do(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		fmt.Println("ERROR: unable to request auth token:", err)
		return agw.ErrorResponse("unable to request auth token", http.StatusInternalServerError)
	}

	ar := new(slack.AuthResponse)
	err = json.NewDecoder(resp.Body).Decode(ar)
	if err != nil {
		fmt.Println("ERROR: unable to decode auth token:", err)
		return agw.ErrorResponse("unable to decode auth token", http.StatusInternalServerError)
	}

	t := secrets.AuthRecord{
		UID:            ar.TeamID,
		UserID:         ar.UserID,
		AccessToken:    ar.AccessToken,
		Scope:          ar.Scope,
		TeamName:       ar.TeamName,
		TeamID:         ar.TeamID,
		BotUserID:      ar.Bot.BotUserID,
		BotAccessToken: ar.Bot.BotAccessToken,
	}

	db := storage.DynamoDB{
		Region: region,
		Table:  authTable,
	}

	err = secrets.SaveTeamTokens(&db, t)
	if err != nil {
		fmt.Println("ERROR: unable to save auth token:", err)
		return agw.ErrorResponse("unable to save auth token", http.StatusInternalServerError)
	}

	fmt.Printf("INFO: authorisation granted for team: %s (%s)\n", t.TeamName, t.TeamID)

	return agw.SuccessResponse()
}
