package main

type BttvClient struct {
	cache *EmoteCache
}

func NewBttvClient() (*BttvClient, error) {
	cache, err := NewCache()

	if err != nil {
		return nil, err
	}

	return &BttvClient{cache}, nil
}

func (client *BttvClient) Get(code string) (string, error) {
	return client.cache.getEmote(code)
}
