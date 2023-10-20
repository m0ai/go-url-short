package store

import (
	generator "go-url-short/internal/shorten"
)

type InMemStore struct {
	urls map[string]string
}

func NewInMemStore() *InMemStore {
	return &InMemStore{
		urls: make(map[string]string),
	}
}

func (s InMemStore) Get(shortKey string) (string, error) {
	var originalURL string
	var found bool

	if originalURL, found = s.urls[shortKey]; !found {
		return "", ErrKeyNotFound{}
	}

	return originalURL, nil
}

func (s InMemStore) Set(originalURL string) (string, error) {
	shortKey := generator.GenerateShortKey()
	if _, err := s.Get(shortKey); err == nil {
		return "", ErrKeyAlreadyExists{}
	}

	s.urls[shortKey] = originalURL
	return shortKey, nil
}
