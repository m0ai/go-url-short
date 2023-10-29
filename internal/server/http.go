package server

import (
	"encoding/json"
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

	if dbConfig.Host == "" {
		st = store.NewInMemStore()
	} else {
		st = store.NewPostgresStore(dbConfig)
	}
	return st
}

type HTTPServerArgs struct {
	Port     string `default:"8080" envconfig:"PORT" required:"true" desc:"Port to listen on"`
	Host     string `default:"localhost" envconfig:"HOST" required:"true" desc:"Address to listen on"`
	Prefix   string `default:"/" envconfig:"PREFIX" required:"true" desc:"Prefix for all routes"`
	DbConfig *store.DatabaseConfig
}

func NewHTTPServer(config *HTTPServerArgs) *http.Server {
	httpLog := log.New(log.Writer(), "HTTPSERVER:", log.LstdFlags)
	s := &httpServer{
		Log:   httpLog,
		Store: configureStore(config.DbConfig),
	}

	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpLog.Println(r.RequestURI)
			next.ServeHTTP(w, r)
		})
	})
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
	w.Write([]byte(`{"message": "ok, I'm healthy"}`))
}

func (s *httpServer) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "Method not allowed"}`))
		json.NewEncoder(w).Encode(&ErrorResponse{r.Method + " Method not allowed"})
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&ErrorResponse{"Missing original url parmas"})
		return
	}

	shortKey, err := s.Store.Set(originalURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&ErrorResponse{"Unhandled Error"})
		return
	}

	// TODO: Delete Hardcoded URL
	host := r.Host
	if r.TLS != nil {
		host = "https://" + host
	} else {
		host = "http://" + host
	}
	result := fmt.Sprintf("%s/%s", host, shortKey)
	s.Log.Printf("Generated short url %s form ", result, originalURL)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&ShortUrlResponse{
		ShortUrl: result,
		Url:      originalURL,
	})
}

func (s *httpServer) handleRedirect(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	shortURL := params["shortURL"]
	if shortURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&ErrorResponse{"Missing short url"})
		return
	}

	originalURL, err := s.Store.Get(shortURL)
	if err != nil && errors.Is(store.ErrKeyNotFound, err) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&ErrorResponse{"Not Found key(" + shortURL + ")"})
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&ErrorResponse{"Unhandled Error"})
		return
	}

	s.Log.Printf("Redirecting key(%s) to %s", shortURL, originalURL)
	http.Redirect(w, r, originalURL, http.StatusPermanentRedirect)
}
