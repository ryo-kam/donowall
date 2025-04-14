package main

import (
	"errors"
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
	if time.Now().Sub(cache.lastFetched) > 24*time.Hour {
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

	for i := 0; i < len(emotes); i++ {
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
		return "", errors.New(fmt.Sprintf("No such emote: %s", code))
	}

	return url, nil
}
