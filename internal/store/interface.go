package store

type Store interface {
	// Get returns the original URL for the given short key
	Get(shortKey string) (string, error)
	// Set saves the original URL and returns the short key
	Set(originalURL string) (string, error)
	// DbClose closes the database connection
	DbClose()
}
