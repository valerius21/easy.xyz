package main

import (
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/valerius21/easy.xyz/shortener"
	"log"
	"net/http"
)

func main() {
	// start with flags
	dbpath := flag.String("db", "./urls.db", "Path to the DB (will be created, if it does not exist")
	hostURL := flag.String("url", "easy.xyz", "host domain name")
	port := flag.Int("port", 8000, "port on which the server listens")

	db, err := bolt.Open(*dbpath, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

	// Start HTTP ShortServer
	shortServer := shortener.ShortServer{
		DB:  db,
		URL: *hostURL,
	}

	defer shortServer.DB.Close()

	addHandler := func() http.Handler {
		return http.HandlerFunc(shortServer.AddURL)
	}

	http.Handle("/add", addHandler())

	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)

	// Provide Start Page

	// If URL is requested [DONE]
	// Lookup in DB [DONE]
	// Redirect if entry is found [DONE]
	// Else 404 [DONE]

	// Provide POST add endpoint [DONE]
	// Look in DB [DONE]
	// Response if name is free with new url else error [DONE]

	// Logging
}
