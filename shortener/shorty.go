package shortener

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"net/http"
	"strings"
)

type ShortServer struct {
	DB *bolt.DB
}

// GetURL redirects to the requested URL
func (s ShortServer) GetURL(w http.ResponseWriter, r *http.Request) (targetURL string, err error) {
	targetURL, err = s.Lookup(r.URL.String())

	if err != nil {
		http.NotFound(w, r)
		return "", err
	}

	http.Redirect(w, r, targetURL, 308)
	return targetURL, nil
}

func (s ShortServer) Lookup(url string) (targetURL string, err error) {
	clearedURL := strings.ReplaceAll(url, "/", "")

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
