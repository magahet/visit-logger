// Package p contains an HTTP Cloud Function.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	lru "github.com/hashicorp/golang-lru"
)

const (
	cacheSize = 100
)

func init() {
	// Disable log prefixes such as the default timestamp.
	// Prefix text prevents the message from being parsed as JSON.
	// A timestamp is added when shipping logs to Cloud Logging.
	log.SetFlags(0)
}

// Entry defines a log entry.
type Entry struct {
	Name     string    `json:"name"`
	Title    string    `json:"title"`
	Url      string    `json:"url"`
	LastSeen time.Time `json:"lastSeen"`
	Count    int       `json:"count"`
}

func (e Entry) String() string {
	j, _ := json.Marshal(e)
	return string(j)
}

type server struct {
	caches   map[string]*lru.Cache
	validate *validator.Validate
}

type Report struct {
	Name    string   `json:"name"`
	Entries []*Entry `json:"entries"`
}

func (s *server) generateReport(name string) []byte {
	cache, ok := s.caches[name]
	if !ok {
		j, _ := json.Marshal(map[string]string{"error": "No entries found for given name"})
		log.Println("No entries found for given name:", name)
		return j
	}
	entries := make([]*Entry, cache.Len())
	for i, k := range cache.Keys() {
		eInt, ok := cache.Peek(k)
		if !ok {
			log.Println("key not found in cache (should not happen)", k)
			continue
		}
		e := eInt.(*Entry)
		entries[i] = e
	}
	r := Report{name, entries}
	j, _ := json.Marshal(r)
	return j
}

func (s *server) rootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		name := r.URL.Query().Get("name")
		if err := s.validate.Var(name, "required,lowercase,alpha,gte=1,lte=20"); err != nil {
			j, _ := json.Marshal(map[string]string{"error": "name is not valid"})
			w.Write(j)
			return

		}

		w.Write(s.generateReport(name))
	case "POST":
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusNoContent)
			log.Println("Sending CORS headers")
			return
		}

		// Set CORS headers for the main request.
		w.Header().Set("Access-Control-Allow-Origin", "*")

		r.Body = http.MaxBytesReader(w, r.Body, 4096)

		var e Entry
		err := json.NewDecoder(r.Body).Decode(&e)
		if err != nil {
			log.Println("Could not parse POST body")
			return
		}

		if e.Url == "" || e.Title == "" || e.Name == "" {
			log.Println("Some or all required fields are missing")
			return
		}

		c, ok := s.caches[e.Name]
		if !ok {
			c, _ = lru.New(cacheSize)
			s.caches[e.Name] = c
			log.Println("New cache created for:", e.Name)
		}

		pInt, ok := c.Get(e.Title)

		// Entry exists in cache
		if ok {
			p := pInt.(*Entry)
			p.LastSeen = time.Now()
			p.Count += 1
			c.Add(p.Title, p)
			log.Println("Entry updated:", p)
		} else {
			e.Count = 1
			e.LastSeen = time.Now()
			c.Add(e.Title, &e)
			log.Println("Entry added:", e)
		}
	}

}

func main() {
	c := make(map[string]*lru.Cache)
	v := validator.New()
	s := server{c, v}

	http.HandleFunc("/", s.rootHandler)
	http.ListenAndServe(":5000", nil)
}
