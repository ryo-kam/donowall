package main

import "fmt"

type Client struct {
	cache *EmoteCache
}

type SlackResponse = map[string]string

func NewClient() (*Client, error) {
	cache, err := NewCache()

	if err != nil {
		return nil, err
	}

	return &Client{cache}, nil
}

func (client *Client) Get(code string) SlackResponse {
	jsonResponse := make(SlackResponse)
	jsonResponse["response_type"] = "ephemeral"

	url, err := client.cache.getEmote(code)

	if err != nil {
		jsonResponse["text"] = fmt.Sprintf("error: %s", err.Error())
	} else {
		jsonResponse["alt_text"] = code
		jsonResponse["image_url"] = url
	}

	return jsonResponse
}
