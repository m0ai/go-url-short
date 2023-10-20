package shorten

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const keyLength = 6

func GenerateShortKey() string {
	// rand.Seed() was deprecated
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[r.Intn(len(charset))]
	}

	return string(shortKey)
}
