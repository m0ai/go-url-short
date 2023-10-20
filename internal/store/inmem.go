package store

import (
	generator "go-url-short/internal/shorten"
)

type ErrKeyAlreadyExists struct {
	Err error
}

type ErrKeyNotFound struct {
	Err error
}

func (e ErrKeyAlreadyExists) Error() string {
	return "key already exists"
}

func (e ErrKeyNotFound) Error() string {
	return "key Not Found"
}

type Store interface {
	// Get returns the original URL for the given short key
	Get(shortKey string) (string, error)
	// Set saves the original URL and returns the short key
	Set(originalURL string) (string, error)
}

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
