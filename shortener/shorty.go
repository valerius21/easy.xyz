package shortener

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"net/http"
	"strings"
)

type ShortServer struct {
	DB  *bolt.DB
	URL string
}

type URLPair struct {
	Shorthand string
	Target    string
}

// GetURL redirects to the requested URL
func (s ShortServer) GetURL(w http.ResponseWriter, r *http.Request) {
	targetURL, err := s.Lookup(r.URL.String())

	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, targetURL, 308)
}

// AddURL adds an URL to the database.
func (s ShortServer) AddURL(w http.ResponseWriter, r *http.Request) {
	var urlPair URLPair

	err := json.NewDecoder(r.Body).Decode(&urlPair)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.Add(urlPair)

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
}

// Lookup an URL for a given Shorthand
func (s ShortServer) Lookup(shorthand string) (targetURL string, err error) {
	clearedURL := strings.ReplaceAll(shorthand, "/", "")

	err = s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("urls"))
		v := b.Get([]byte(clearedURL))

		if v == nil {
			errorMessage := fmt.Sprintf("could not find the proper URL for %s (unknown)", clearedURL)
			return errors.New(errorMessage)
		}

		targetURL = string(v)
		return nil
	})

	return targetURL, err
}

// Add an urlPair to the Database
func (s ShortServer) Add(urlPair URLPair) (err error) {
	return s.DB.Update(func(tx *bolt.Tx) error {
		targetURL, err := s.Lookup(urlPair.Shorthand)
		if err != nil && err.Error() != fmt.Sprintf("could not find the proper URL for %s (unknown)",
			urlPair.Shorthand) {
			return err
		}

		if targetURL != "" {
			return errors.New(fmt.Sprintf("%s does already exist!", urlPair.Shorthand))
		}

		bucket, err := tx.CreateBucketIfNotExists([]byte("urls"))
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(urlPair.Shorthand), []byte(urlPair.Target))
		if err != nil {
			return err
		}

		return nil
	})
}
