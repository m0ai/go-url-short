package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-url-short/internal/store"
	"log"
	"net/http"
)

type httpServer struct {
	Log   *log.Logger
	Store *store.InMemStore
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Log:   log.New(log.Writer(), "httpServer", log.LstdFlags),
		Store: store.NewInMemStore(),
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func NewHTTPServer(add, port string) *http.Server {
	s := newHTTPServer()

	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/health", s.handleHealthCheck).Methods("GET")
	r.HandleFunc("/shorten", s.handleShorten).Methods("POST")
	r.HandleFunc("/s/{shortURL}", s.handleRedirect)

	return &http.Server{
		Addr:    add + ":" + port,
		Handler: r,
	}
}

func (s *httpServer) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK, I'm alive!"))
}

func (s *httpServer) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request, missing url from form"))
		return
	}

	shortKey, err := s.Store.Set(originalURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	s.Log.Println(s.Store.Get(shortKey))
	w.WriteHeader(http.StatusCreated)

	result := fmt.Sprintf("http://localhost:8081/s/%s", shortKey)
	s.Log.Println(result)
	// TODO return json
	w.Write([]byte(result))
}

func (s *httpServer) handleRedirect(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	shortURL := params["shortURL"]

	if shortURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request, missing shortURL"))
		return
	}
	originalURL, err := s.Store.Get(shortURL)

	s.Log.Println(originalURL, err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error, Could not get original URL"))
		return
	}

	http.Redirect(w, r, originalURL, http.StatusPermanentRedirect)
}