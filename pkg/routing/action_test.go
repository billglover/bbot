package routing

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

var tcs = []struct {
	Name    string
	Request Request
	Action  MessageAction
}{
	{
		Name: "flagged user message",
		Request: Request{
			PathParameters: map[string]string{"type": "action"},
			Body:           "payload=%7B%22type%22%3A%22message_action%22%2C%22token%22%3A%22w06XkgVo2IlRaQRamizypPQl%22%2C%22action_ts%22%3A%221535885531.310842%22%2C%22team%22%3A%7B%22id%22%3A%22TBLG57ECT%22%2C%22domain%22%3A%22buddybotdev%22%7D%2C%22user%22%3A%7B%22id%22%3A%22UBLKAG9K4%22%2C%22name%22%3A%22bill%22%7D%2C%22channel%22%3A%7B%22id%22%3A%22CBLPRTX3P%22%2C%22name%22%3A%22general%22%7D%2C%22callback_id%22%3A%22flagMessage%22%2C%22trigger_id%22%3A%22427767237634.394549252435.6b8540a80ec073eab2ed85b4550a236a%22%2C%22message_ts%22%3A%221535813905.000100%22%2C%22message%22%3A%7B%22type%22%3A%22message%22%2C%22user%22%3A%22UBLKAG9K4%22%2C%22text%22%3A%22hello%22%2C%22client_msg_id%22%3A%2254360322-36c9-4cd8-a7ab-13b5322362a6%22%2C%22ts%22%3A%221535813905.000100%22%7D%2C%22response_url%22%3A%22https%3A%5C%2F%5C%2Fhooks.slack.com%5C%2Fapp%5C%2FTBLG57ECT%5C%2F429109817558%5C%2FgnYAErMdhKzXWvT1CmJpVPGG%22%7D",
			HTTPMethod:     http.MethodPost,
			Headers: map[string]string{
				"X-Slack-Request-Timestamp": "1531420618",
				"X-Slack-Signature":         "v0=a2114d57b48eac39b9ad189dd8316235a7b4a8d21a10bd27519666489c69b503"},
		},
		Action: MessageAction{
			Type:             "message_action",
			CallbackID:       "flagMessage",
			Team:             Team{ID: "TBLG57ECT", Domain: "buddybotdev"},
			Channel:          Channel{ID: "CBLPRTX3P", Name: "general"},
			User:             User{ID: "UBLKAG9K4", Name: "bill"},
			ActionTimestamp:  json.Number("1535885531.310842"),
			MessageTimestamp: json.Number("1535813905.000100"),
			Message: Message{
				Msg: Msg{
					UserID:    "UBLKAG9K4",
					Type:      "message",
					Text:      "hello",
					Timestamp: "1535813905.000100",
				},
			},
			ResponseURL: "https://hooks.slack.com/app/TBLG57ECT/429109817558/gnYAErMdhKzXWvT1CmJpVPGG",
			TriggerID:   "427767237634.394549252435.6b8540a80ec073eab2ed85b4550a236a",
		},
	},
	{
		Name: "flagged attachment message",
		Request: Request{
			PathParameters: map[string]string{"type": "action"},
			Body:           "payload=%7B%22type%22%3A%22message_action%22%2C%22token%22%3A%22w06XkgVo2IlRaQRamizypPQl%22%2C%22action_ts%22%3A%221536058881.467448%22%2C%22team%22%3A%7B%22id%22%3A%22TBLG57ECT%22%2C%22domain%22%3A%22buddybotdev%22%7D%2C%22user%22%3A%7B%22id%22%3A%22UBLKAG9K4%22%2C%22name%22%3A%22bill%22%7D%2C%22channel%22%3A%7B%22id%22%3A%22CBLPRTX3P%22%2C%22name%22%3A%22general%22%7D%2C%22callback_id%22%3A%22flagMessage%22%2C%22trigger_id%22%3A%22428810990741.394549252435.c931785a55575952f2c67edb4950858f%22%2C%22message_ts%22%3A%221536058875.000100%22%2C%22message%22%3A%7B%22type%22%3A%22message%22%2C%22text%22%3A%22%22%2C%22files%22%3A%5B%7B%22id%22%3A%22FCMLEM1FH%22%2C%22created%22%3A1536058872%2C%22timestamp%22%3A1536058872%2C%22name%22%3A%22BuddyBot+V2.png%22%2C%22title%22%3A%22BuddyBot+V2.png%22%2C%22mimetype%22%3A%22image%5C%2Fpng%22%2C%22filetype%22%3A%22png%22%2C%22pretty_type%22%3A%22PNG%22%2C%22user%22%3A%22UBLKAG9K4%22%2C%22editable%22%3Afalse%2C%22size%22%3A74905%2C%22mode%22%3A%22hosted%22%2C%22is_external%22%3Afalse%2C%22external_type%22%3A%22%22%2C%22is_public%22%3Atrue%2C%22public_url_shared%22%3Afalse%2C%22display_as_bot%22%3Afalse%2C%22username%22%3A%22%22%2C%22url_private%22%3A%22https%3A%5C%2F%5C%2Ffiles.slack.com%5C%2Ffiles-pri%5C%2FTBLG57ECT-FCMLEM1FH%5C%2Fbuddybot_v2.png%22%2C%22url_private_download%22%3A%22https%3A%5C%2F%5C%2Ffiles.slack.com%5C%2Ffiles-pri%5C%2FTBLG57ECT-FCMLEM1FH%5C%2Fdownload%5C%2Fbuddybot_v2.png%22%2C%22thumb_64%22%3A%22https%3A%5C%2F%5C%2Ffiles.slack.com%5C%2Ffiles-tmb%5C%2FTBLG57ECT-FCMLEM1FH-b95d8a6fea%5C%2Fbuddybot_v2_64.png%22%2C%22thumb_80%22%3A%22https%3A%5C%2F%5C%2Ffiles.slack.com%5C%2Ffiles-tmb%5C%2FTBLG57ECT-FCMLEM1FH-b95d8a6fea%5C%2Fbuddybot_v2_80.png%22%2C%22thumb_360%22%3A%22https%3A%5C%2F%5C%2Ffiles.slack.com%5C%2Ffiles-tmb%5C%2FTBLG57ECT-FCMLEM1FH-b95d8a6fea%5C%2Fbuddybot_v2_360.png%22%2C%22thumb_360_w%22%3A360%2C%22thumb_360_h%22%3A203%2C%22thumb_480%22%3A%22https%3A%5C%2F%5C%2Ffiles.slack.com%5C%2Ffiles-tmb%5C%2FTBLG57ECT-FCMLEM1FH-b95d8a6fea%5C%2Fbuddybot_v2_480.png%22%2C%22thumb_480_w%22%3A480%2C%22thumb_480_h%22%3A270%2C%22thumb_160%22%3A%22https%3A%5C%2F%5C%2Ffiles.slack.com%5C%2Ffiles-tmb%5C%2FTBLG57ECT-FCMLEM1FH-b95d8a6fea%5C%2Fbuddybot_v2_160.png%22%2C%22thumb_720%22%3A%22https%3A%5C%2F%5C%2Ffiles.slack.com%5C%2Ffiles-tmb%5C%2FTBLG57ECT-FCMLEM1FH-b95d8a6fea%5C%2Fbuddybot_v2_720.png%22%2C%22thumb_720_w%22%3A720%2C%22thumb_720_h%22%3A405%2C%22thumb_800%22%3A%22https%3A%5C%2F%5C%2Ffiles.slack.com%5C%2Ffiles-tmb%5C%2FTBLG57ECT-FCMLEM1FH-b95d8a6fea%5C%2Fbuddybot_v2_800.png%22%2C%22thumb_800_w%22%3A800%2C%22thumb_800_h%22%3A450%2C%22thumb_960%22%3A%22https%3A%5C%2F%5C%2Ffiles.slack.com%5C%2Ffiles-tmb%5C%2FTBLG57ECT-FCMLEM1FH-b95d8a6fea%5C%2Fbuddybot_v2_960.png%22%2C%22thumb_960_w%22%3A960%2C%22thumb_960_h%22%3A540%2C%22thumb_1024%22%3A%22https%3A%5C%2F%5C%2Ffiles.slack.com%5C%2Ffiles-tmb%5C%2FTBLG57ECT-FCMLEM1FH-b95d8a6fea%5C%2Fbuddybot_v2_1024.png%22%2C%22thumb_1024_w%22%3A1024%2C%22thumb_1024_h%22%3A576%2C%22image_exif_rotation%22%3A1%2C%22original_w%22%3A1280%2C%22original_h%22%3A720%2C%22permalink%22%3A%22https%3A%5C%2F%5C%2Fbuddybotdev.slack.com%5C%2Ffiles%5C%2FUBLKAG9K4%5C%2FFCMLEM1FH%5C%2Fbuddybot_v2.png%22%2C%22permalink_public%22%3A%22https%3A%5C%2F%5C%2Fslack-files.com%5C%2FTBLG57ECT-FCMLEM1FH-2609594d27%22%2C%22is_starred%22%3Afalse%7D%5D%2C%22upload%22%3Atrue%2C%22user%22%3A%22UBLKAG9K4%22%2C%22display_as_bot%22%3Afalse%2C%22ts%22%3A%221536058875.000100%22%7D%2C%22response_url%22%3A%22https%3A%5C%2F%5C%2Fhooks.slack.com%5C%2Fapp%5C%2FTBLG57ECT%5C%2F428532210739%5C%2FVSw8tqPU8SNKYcQ0ZIoYrh7G%22%7D",
			HTTPMethod:     http.MethodPost,
			Headers: map[string]string{
				"X-Slack-Request-Timestamp": "1531420618",
				"X-Slack-Signature":         "v0=a2114d57b48eac39b9ad189dd8316235a7b4a8d21a10bd27519666489c69b503"},
		},
		Action: MessageAction{
			Type:             "message_action",
			CallbackID:       "flagMessage",
			Team:             Team{ID: "TBLG57ECT", Domain: "buddybotdev"},
			Channel:          Channel{ID: "CBLPRTX3P", Name: "general"},
			User:             User{ID: "UBLKAG9K4", Name: "bill"},
			ActionTimestamp:  json.Number("1536058881.467448"),
			MessageTimestamp: json.Number("1536058875.000100"),
			Message: Message{
				Msg: Msg{
					UserID:    "UBLKAG9K4",
					Type:      "message",
					Text:      "",
					Timestamp: "1536058875.000100",
				},
			},
			ResponseURL: "https://hooks.slack.com/app/TBLG57ECT/428532210739/VSw8tqPU8SNKYcQ0ZIoYrh7G",
			TriggerID:   "428810990741.394549252435.c931785a55575952f2c67edb4950858f",
		},
	},
	{
		Name: "flagged bot message",
		Request: Request{
			PathParameters: map[string]string{"type": "action"},
			Body:           "payload=%7B%22type%22%3A%22message_action%22%2C%22token%22%3A%22w06XkgVo2IlRaQRamizypPQl%22%2C%22action_ts%22%3A%221536060699.383687%22%2C%22team%22%3A%7B%22id%22%3A%22TBLG57ECT%22%2C%22domain%22%3A%22buddybotdev%22%7D%2C%22user%22%3A%7B%22id%22%3A%22UBLKAG9K4%22%2C%22name%22%3A%22bill%22%7D%2C%22channel%22%3A%7B%22id%22%3A%22CBLPRTX3P%22%2C%22name%22%3A%22general%22%7D%2C%22callback_id%22%3A%22flagMessage%22%2C%22trigger_id%22%3A%22428827157461.394549252435.9a81b8671849fb4f378f33988f09d1b5%22%2C%22message_ts%22%3A%221533595230.000090%22%2C%22message%22%3A%7B%22text%22%3A%22Congrats+%3C%40UBLPTK0JH%3E%21+Score+now+at+8+%3Asmile%3A%22%2C%22username%22%3A%22buddybot%22%2C%22bot_id%22%3A%22BBL3GSL7K%22%2C%22mrkdwn%22%3Afalse%2C%22type%22%3A%22message%22%2C%22subtype%22%3A%22bot_message%22%2C%22ts%22%3A%221533595230.000090%22%7D%2C%22response_url%22%3A%22https%3A%5C%2F%5C%2Fhooks.slack.com%5C%2Fapp%5C%2FTBLG57ECT%5C%2F428670626322%5C%2FnQqdNXEWWZJpVgop3FTTM4BH%22%7D",
			HTTPMethod:     http.MethodPost,
			Headers: map[string]string{
				"X-Slack-Request-Timestamp": "1531420618",
				"X-Slack-Signature":         "v0=a2114d57b48eac39b9ad189dd8316235a7b4a8d21a10bd27519666489c69b503"},
		},
		Action: MessageAction{
			Type:             "message_action",
			CallbackID:       "flagMessage",
			Team:             Team{ID: "TBLG57ECT", Domain: "buddybotdev"},
			Channel:          Channel{ID: "CBLPRTX3P", Name: "general"},
			User:             User{ID: "UBLKAG9K4", Name: "bill"},
			ActionTimestamp:  json.Number("1536060699.383687"),
			MessageTimestamp: json.Number("1533595230.000090"),
			Message: Message{
				Msg: Msg{
					Type:      "message",
					Text:      "Congrats <@UBLPTK0JH>! Score now at 8 :smile:",
					Timestamp: "1533595230.000090",
					SubType:   "bot_message",
					BotID:     "BBL3GSL7K",
					BotName:   "buddybot",
				},
			},
			ResponseURL: "https://hooks.slack.com/app/TBLG57ECT/428670626322/nQqdNXEWWZJpVgop3FTTM4BH",
			TriggerID:   "428827157461.394549252435.9a81b8671849fb4f378f33988f09d1b5",
		},
	},
}

func TestParseAction(t *testing.T) {
	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			ma, err := ParseAction(tc.Request.Body)
			if err != nil {
				t.Errorf("unable to parse message body: %v", err)
			}

			if got, want := ma.Team.ID, "TBLG57ECT"; got != want {
				t.Errorf("unable to parse TeamID: got %s, want %s", got, want)
			}

			if got, want := ma, tc.Action; reflect.DeepEqual(got, want) == false {
				t.Errorf("\n got: %+v\nwant: %+v", got, want)
			}
		})
	}
}
