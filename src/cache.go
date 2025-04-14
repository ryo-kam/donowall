package main

import (
	"fmt"
	"time"

	"github.com/pajlada/gobttv"
)

type EmoteMap = map[string]string

type EmoteCache struct {
	instance    *gobttv.BTTVAPI
	lastFetched time.Time
	emotesCache EmoteMap
}

func NewCache() (*EmoteCache, error) {
	instance := gobttv.New()
	emotes, err := instance.GetEmotes()

	if err != nil {
		return nil, err
	}

	return &EmoteCache{
		instance:    instance,
		lastFetched: time.Now(),
		emotesCache: transformToMap(emotes),
	}, nil
}

func (cache *EmoteCache) getEmotes() (EmoteMap, error) {
	if time.Since(cache.lastFetched) > 24*time.Hour {
		emotes, err := cache.instance.GetEmotes()

		if err != nil {
			return nil, err
		}

		cache.emotesCache = transformToMap(emotes)
		cache.lastFetched = time.Now()
	}

	return cache.emotesCache, nil
}

func transformToMap(emotes []gobttv.Emote) EmoteMap {
	emoteMap := make(EmoteMap)

	for i := range emotes {
		emoteMap[emotes[i].Code] = emotes[i].URLs.X4
	}

	return emoteMap
}

func (cache *EmoteCache) getEmote(code string) (string, error) {
	emoteMap, err := cache.getEmotes()

	if err != nil {
		return "", err
	}

	url, ok := emoteMap[code]

	if !ok {
		return "", fmt.Errorf("no such emote: %s", code)
	}

	return url, nil
}
