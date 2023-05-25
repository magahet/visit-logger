// Package p contains an HTTP Cloud Function.
package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"

	lru "github.com/hashicorp/golang-lru"
)

const (
	cacheSize  = 100
	apiKeyPath = "/secrets/apikey"
)

var (
	caches   map[string]*lru.Cache
	validate *validator.Validate
	apiKey   string
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

type Report struct {
	Name    string   `json:"name"`
	Entries []*Entry `json:"entries"`
}

func generateReport(name string) []byte {
	cache, ok := caches[name]
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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		name := r.URL.Query().Get("name")
		if err := validate.Var(name, "required,lowercase,alpha,gte=1,lte=20"); err != nil {
			j, _ := json.Marshal(map[string]string{"error": "name is not valid"})
			w.Write(j)
			return
		}

		// Auth check
		if r.Header.Get("X-Api-Key") != apiKey {
			w.WriteHeader(http.StatusForbidden)
			log.Println("Auth header invalid:", r.Header.Get("X-Api-Key"))
			return
		}

		w.Write(generateReport(name))

	case http.MethodPost:
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

		c, ok := caches[e.Name]
		if !ok {
			c, _ = lru.New(cacheSize)
			caches[e.Name] = c
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

func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Api-Key")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusNoContent)
			log.Println("Sending CORS headers")
			return
		}
		f(w, r)
	}
}

func namesHandler(w http.ResponseWriter, r *http.Request) {
	// Auth check
	if r.Header.Get("X-Api-Key") != apiKey {
		w.WriteHeader(http.StatusForbidden)
		log.Printf("Auth header invalid: [%s] != [%s]", r.Header.Get("X-Api-Key"), apiKey)
		return
	}

	names := make([]string, 0, len(caches))
	for n := range caches {
		names = append(names, n)
	}

	j, _ := json.Marshal(map[string][]string{"names": names})
	w.Write(j)
}

func main() {

	flag.StringVar(&apiKey, "key", "", "API Key")
	flag.Parse()

	caches = make(map[string]*lru.Cache)
	validate = validator.New()

	if apiKey == "" {
		apiKeyBytes, err := os.ReadFile(apiKeyPath)
		if err != nil {
			log.Fatal(err)
		}
		apiKey = strings.TrimSpace(string(apiKeyBytes))
	}

	http.HandleFunc("/logs", cors(rootHandler))
	http.HandleFunc("/names", cors(namesHandler))
	http.ListenAndServe(":5000", nil)
}
