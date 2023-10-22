package store

import (
	generator "go-url-short/internal/shorten"
	"log"
)

type InMemStore struct {
	urls map[string]string
	Log  *log.Logger
}

func NewInMemStore() *InMemStore {
	return &InMemStore{
		urls: make(map[string]string),
		Log:  log.New(log.Writer(), "INMEMSTORE:", log.LstdFlags),
	}
}
func (s InMemStore) DbClose() {
	s.Log.Println("Closing database connection")
	s.urls = make(map[string]string)
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

	if shortKey, found := s.urls[shortKey]; found {
		return shortKey, ErrKeyAlreadyExists{}
	}

	if _, err := s.Get(shortKey); err == nil {
		return "", ErrKeyAlreadyExists{}
	}

	s.urls[shortKey] = originalURL
	return shortKey, nil
}
