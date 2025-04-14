package main

import (
	"fmt"
)

type Client struct {
	cache *EmoteCache
}

type SlackResponse struct {
	ResponseType string              `json:"response_type"`
	Blocks       []map[string]string `json:"blocks"`
	Text         string              `json:"text"`
}

func NewClient() (*Client, error) {
	cache, err := NewCache()

	if err != nil {
		return nil, err
	}

	return &Client{cache}, nil
}

func (client *Client) Get(code string) SlackResponse {

	url, err := client.cache.getEmote(code)

	if err != nil {
		return SlackResponse{
			ResponseType: "in_channel",
			Text:         fmt.Sprintf("error: %s", err.Error()),
		}
	} else {
		imgBlock := make(map[string]string)

		imgBlock["type"] = "image"
		imgBlock["alt_text"] = code
		imgBlock["image_url"] = url

		return SlackResponse{
			ResponseType: "in_channel",
			Text:         code,
			Blocks:       []map[string]string{imgBlock},
		}
	}
}
