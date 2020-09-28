package main

import "flag"

func main() {
	// start with flags
	flag.String("db", "./urls.db", "Path to the DB (will be created, if it does not exist")
	flag.String("url", "easy.xyz", "host domain name")

	// Start HTTP ShortServer

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
