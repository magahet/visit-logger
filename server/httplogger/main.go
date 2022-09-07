// Package p contains an HTTP Cloud Function.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

func init() {
	// Disable log prefixes such as the default timestamp.
	// Prefix text prevents the message from being parsed as JSON.
	// A timestamp is added when shipping logs to Cloud Logging.
	log.SetFlags(0)
}

// Entry defines a log entry.
type Entry struct {
	// Logs Explorer allows filtering and display of this as `jsonPayload.component`.
	Name     string
	Title    string
	Url      string
	LastSeen time.Time
}

type server struct {
	cache *lru.Cache
}

func (s *server) log(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	var e Entry
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		return
	}

	if e.Url == "" || e.Title == "" || e.Name == "" {
		return
	}

	s.cache.Add(e.Title, e)

}

func main() {
	cache, err := lru.New(100)
	if err != nil {
		log.Fatal(err)
	}

	s := server{cache: cache}

	http.HandleFunc("/", s.log)
	http.ListenAndServe(":5000", nil)
}
