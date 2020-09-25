package main

//import (
//	"github.com/valerius21/shorty/shortener"
//	"net/http"
//)
//
//var server = shortener.ShortServer{URLs: shortener.URLDictionary{
//	"lol": "htts://httpbin.org/",
//}}
//
//func handler(w http.ResponseWriter, r *http.Request) {
//	_, _ = server.GetURL(w, r)
//}

func main() {
	// Start HTTP ShortServer
	//handler := http.HandlerFunc(handler)
	//if err := http.ListenAndServe(":5000", handler); err != nil {
	//	log.Fatalf("could not listen on port 5000 %v", err)
	//}
	// Provide Start Page

	// If URL is requested
	// Lookup in DB
	// Redirect if entry is found [DONE]
	// Else 404 [DONE]

	// Provide POST add endpoint
	// Look in DB
	// Response if name is free with new url
}
