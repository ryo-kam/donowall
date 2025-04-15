package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/pajlada/gobttv"
)

type EmoteMap = map[string]string

type EmoteCache struct {
	instance    *gobttv.BTTVAPI
	lastFetched time.Time
	emotesCache EmoteMap
	sync.Mutex
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

func (cache *EmoteCache) getEmote(code string) (string, error) {
	cache.Lock()
	isStale := time.Since(cache.lastFetched) > 24*time.Hour
	cache.Unlock()

	if isStale {
		emotes, err := cache.instance.GetEmotes()

		if err != nil {
			return "", err
		}

		cacheMap := transformToMap(emotes)

		cache.Lock()
		cache.emotesCache = cacheMap
		cache.lastFetched = time.Now()
		cache.Unlock()
	}

	cache.Lock()
	url, ok := cache.emotesCache[code]
	cache.Unlock()

	if !ok {
		return "", fmt.Errorf("no such emote: %s", code)
	}

	return url, nil
}

func transformToMap(emotes []gobttv.Emote) EmoteMap {
	emoteMap := make(EmoteMap)

	for i := range emotes {
		emoteMap[emotes[i].Code] = emotes[i].URLs.X4
	}

	return emoteMap
}
