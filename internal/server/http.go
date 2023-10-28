package server

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"go-url-short/internal/store"
	"log"
	"net/http"
	"strings"
)

type httpServer struct {
	Log   *log.Logger
	Store store.Store
}

func configureStore(dbConfig *store.DatabaseConfig) store.Store {
	var st store.Store

	if dbConfig.Name == "" {
		st = store.NewInMemStore()
	} else {
		st = store.NewPostgresStore(dbConfig)
	}
	return st
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

type HTTPServerArgs struct {
	Port     string                `default:"8080" envconfig:"PORT" required:"true" desc:"Port to listen on"`
	Host     string                `default:"localhost" envconfig:"HOST" required:"true" desc:"Address to listen on"`
	Prefix   string                `default:"/" envconfig:"PREFIX" required:"true" desc:"Prefix for all routes"`
	DbConfig *store.DatabaseConfig `default:nil`
}

func NewHTTPServer(config *HTTPServerArgs) *http.Server {
	httpLog := log.New(log.Writer(), "HTTPSERVER:", log.LstdFlags)
	s := &httpServer{
		Log:   httpLog,
		Store: configureStore(config.DbConfig),
	}
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Log.Println("Method not allowed", r.RequestURI)
		http.Error(w, fmt.Sprintf("Method not allowed: %s", r.Method), http.StatusMethodNotAllowed)
	})
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Log.Println("Not found", r.RequestURI)
		http.Error(w, fmt.Sprintf("Not found: %s", r.RequestURI), http.StatusNotFound)
	})

	r.HandleFunc("/health", s.handleHealthCheck).Methods("GET")
	r.HandleFunc("/shorten", s.handleShorten).Methods("POST")
	r.HandleFunc("/{shortURL}", s.handleRedirect)
	return &http.Server{
		Addr:    strings.Join([]string{config.Host, ":", config.Port}, ""),
		Handler: r,
	}
}

func (s *httpServer) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"msg": "ok, I'm healthy"}`))

}

func (s *httpServer) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "Method not allowed"}`))
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "Missing url"}`))
		return
	}

	shortKey, err := s.Store.Set(originalURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "Internal server error"}`))
		return
	}
	w.WriteHeader(http.StatusCreated)

	// TODO: Delete Hardcoded URL
	host := r.Host
	if r.TLS != nil {
		host = "https://" + host
	} else {
		host = "http://" + host
	}
	result := fmt.Sprintf("%s/%s", host, shortKey)
	s.Log.Printf("Shortened URL %s from %s", result, originalURL)
	// TODO return json
	w.Write([]byte(result))
}

func (s *httpServer) handleRedirect(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	shortURL := params["shortURL"]
	if shortURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "Missing shortURL"}`))
		return
	}

	originalURL, err := s.Store.Get(shortURL)
	if err != nil && errors.Is(store.ErrKeyNotFound, err) {
		w.WriteHeader(http.StatusNotFound)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "Internal server error"}`))
		return
	}

	s.Log.Printf("Redirecting key(%s) to %s", shortURL, originalURL)
	http.Redirect(w, r, originalURL, http.StatusPermanentRedirect)
}
