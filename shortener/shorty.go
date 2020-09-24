package shortener

import (
	"errors"
	"fmt"
)

//func main() {
// Start HTTP Server
// Provide Start Page

// If URL is requested
// Lookup in DB
// Redirect if entry is found
// Else 404

// Provide POST add endpoint
// Look in DB
// Response if name is free with new url

//}

// ShortServer listens and serves the shortener
//func ShortServer(w http.ResponseWriter, r *http.Request) {
//	//http.Redirect(w,r, , 302)
//}

type URLDictionary map[string]string

func (urls URLDictionary) LookupURL(shortURL string) (string, error) {
	destinationURL, ok := urls[shortURL]

	if !ok {
		errorMessage := fmt.Sprintf("could not find the proper URL for %s (unknown)", shortURL)
		return "", errors.New(errorMessage)
	}

	return destinationURL, nil
}
