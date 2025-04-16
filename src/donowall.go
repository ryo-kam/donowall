package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type SlackResponse struct {
	ResponseType    string `json:"response_type"`
	ReplaceOriginal bool   `json:"replace_original"`
	Text            string `json:"text"`
	Blocks          []any  `json:"blocks"`
}

func donowallHandler(bttvClient *BttvClient) socketmode.SocketmodeHandlerFunc {
	return func(event *socketmode.Event, slackClient *socketmode.Client) {
		cmd, ok := event.Data.(slack.SlashCommand)

		if !ok {
			slackClient.Debugf("Ignored %+v\n", event)
			return
		}

		slackClient.Debugf("[/donowall] invoked :%+v", cmd)

		ackPayload := SlackResponse{
			ResponseType:    slack.ResponseTypeInChannel,
			ReplaceOriginal: true,
			Text:            cmd.Text,
		}

		slackClient.Ack(*event.Request, ackPayload)

		code := strings.Split(cmd.Text, " ")[0]

		emoteUrl, err := bttvClient.Get(code)

		var response SlackResponse

		if err != nil {
			response = SlackResponse{
				ResponseType:    "in_channel",
				ReplaceOriginal: true,
				Text:            fmt.Sprintf("error: %s", err.Error()),
			}
		} else {
			imgBlock := make(map[string]string)

			imgBlock["type"] = "image"
			imgBlock["alt_text"] = code
			imgBlock["image_url"] = emoteUrl

			response = SlackResponse{
				ResponseType:    "in_channel",
				ReplaceOriginal: true,
				Text:            code,
				Blocks: []any{
					slack.ImageBlock{
						Type:     "image",
						ImageURL: emoteUrl,
						AltText:  code,
					},
				},
			}
		}

		jsonData, err := json.Marshal(response)

		if err != nil {
			slackClient.Debugf("[/donowall] error: %+v", err)
			return
		}

		resp, err := http.Post(cmd.ResponseURL, "application/json", bytes.NewBuffer(jsonData))

		if err != nil {
			slackClient.Debugf("[/donowall] error: %+v", err)
			return
		}

		defer resp.Body.Close()

		slackClient.Debugf("[/donowall] response: %+v", resp)
	}
}
