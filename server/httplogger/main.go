// Package p contains an HTTP Cloud Function.
package main

import (
	"encoding/json"
	"io"
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
	Count	int
}

type server struct {
	caches map[string]*lru.Cache
}

type Report struct {
	Name string `json:"name"`
	Entries string `json:"entries"`
}

func (s *server) generateReport() string {
	entries := [s.cache.Len()]*Entry
	for i, k := range s.cache.Keys() {
		eInt, ok := s.cache.Peek(k)
		if !ok {
			continue
		}
		e := eInt.(*Entry)
		entries[i] = e
	}
	r := Report{name, entries}
	return json.Marshal()
}

func (s *server) rootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		io.WriteString(w, s.generateReport())
	case "POST":
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

		r.Body = http.MaxBytesReader(w, r.Body, 4096)

		var e Entry
		err := json.NewDecoder(r.Body).Decode(&e)
		if err != nil {
			return
		}

		if e.Url == "" || e.Title == "" || e.Name == "" {
			return
		}

		pInt, ok := s.cache.Get(e.Title)

		// Entry exists in cache
		if ok {
			p := pInt.(*Entry)
			p.LastSeen = time.Now()
			p.Count += 1
			s.cache.Add(p.Title, p)
		} else {
			e.Count = 1
			e.LastSeen = time.Now()
			s.cache.Add(e.Title, e)
		}
	}

}

func main() {
	cache, err := lru.New(100)
	if err != nil {
		log.Fatal(err)
	}

	s := server{cache: cache}

	http.HandleFunc("/", s.rootHandler)
	http.ListenAndServe(":5000", nil)
}
