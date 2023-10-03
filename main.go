package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// inmemory store
type URLShortener struct {
	urls map[string]string
}

const urlPrefix = "/short/"

func (us *URLShortener) handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[len(urlPrefix):]
	if shortKey == "" {
		http.Error(w, "Short key is missing", http.StatusBadRequest)
		return
	}

	originalURL, found := us.urls[shortKey]
	if !found {
		http.Error(w, "Short key is not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
	return
}

func (us *URLShortener) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// gnerate short key (TODO: using a hash function as MOD)
	shortKey := generateShortKey()
	us.urls[shortKey] = originalURL
	shortenedURL := fmt.Sprintf("http://%s%s%s", r.Host, urlPrefix, shortKey)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	responseHTML := fmt.Sprintf(`
        <h2>URL Shortener</h2>
        <p>Original URL: %s</p>
        <p>Shortened URL: <a href="%s">%s</a></p>
        <form method="post" action="/shorten">
            <input type="text" name="url" placeholder="Enter a URL">
            <input type="submit" value="Shorten">
        </form>
	`, originalURL, shortenedURL, shortenedURL)
	fmt.Fprintf(w, responseHTML)
	w.WriteHeader(http.StatusCreated)
	return
}
func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	// rand.Seed() was deprecated
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[r.Intn(len(charset))]
	}
	return string(shortKey)
}

func main() {
	s := &URLShortener{
		urls: make(map[string]string),
	}
	http.HandleFunc(urlPrefix, s.handleRedirect)
	http.HandleFunc("/shorten", s.handleShorten)

	var port string

	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}

	fmt.Println("Server is listening on port: ", port)
	http.ListenAndServe(":"+port, nil)
}
