package main

import (
	"fmt"
)

type Client struct {
	cache *EmoteCache
}

type SlackResponse struct {
	response_type string
	blocks        []map[string]string
	text          string
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
			response_type: "in_channel",
			text:          fmt.Sprintf("error: %s", err.Error()),
		}
	} else {
		imgBlock := make(map[string]string)

		imgBlock["type"] = "image"
		imgBlock["alt_text"] = code
		imgBlock["image_url"] = url

		return SlackResponse{
			response_type: "in_channel",
			text:          code,
			blocks:        []map[string]string{imgBlock},
		}
	}
}
