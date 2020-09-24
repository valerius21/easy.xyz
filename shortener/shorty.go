package shortener

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Server struct {
	//RedirectHandler func(w http.ResponseWriter, r *http.Request)
	URLs URLDictionary
}

func (s Server) GetURL(w http.ResponseWriter, r *http.Request) (targetURL string, error error) {
	clearedURL := strings.ReplaceAll(r.URL.String(), "/", "")
	targetURL, error = s.URLs.Lookup(clearedURL)

	if error != nil {
		return "", error
	}

	http.Redirect(w, r, targetURL, 301)
	return targetURL, nil
}

// URLDictionary which holds the shorthands and the destination URLs
type URLDictionary map[string]string

// Lookup a given url in the URLDictionary
func (urls URLDictionary) Lookup(shortURL string) (string, error) {
	destinationURL, ok := urls[shortURL]

	if !ok {
		errorMessage := fmt.Sprintf("could not find the proper URL for %s (unknown)", shortURL)
		return "", errors.New(errorMessage)
	}

	return destinationURL, nil
}
