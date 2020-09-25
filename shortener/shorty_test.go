package shortener

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var shortServer = ShortServer{URLs: URLDictionary{
	"test": "http://test.local/",
}}

func TestGetURL(t *testing.T) {
	t.Run("Getting the given URLs", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/test", nil)
		response := httptest.NewRecorder()

		gotDest, err := shortServer.GetURL(response, request)
		wantDest := "http://test.local/"

		assert.Nil(t, err)

		got := response.Result().StatusCode
		want := 308 // redirect status code

		assert.Equal(t, got, want)
		assert.Equal(t, gotDest, wantDest)
	})

	t.Run("Requesting a non-existing shortcut", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/unknown", nil)
		response := httptest.NewRecorder()

		_, err := shortServer.GetURL(response, request)

		assert.NotNil(t, err)

		got := response.Result().StatusCode
		want := 404
		assert.Equal(t, got, want)
	})
}

func TestLookup(t *testing.T) {
	t.Run("Lookup and find", func(t *testing.T) {
		got, err := shortServer.URLs.Lookup("test")
		want := "http://test.local/"

		assert.Nil(t, err)
		assert.Equal(t, got, want)
	})

	t.Run("Lookup non existing key", func(t *testing.T) {
		url := "unknown"
		_, err := shortServer.URLs.Lookup(url)
		want := fmt.Sprintf("could not find the proper URL for %s (unknown)", url)

		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), want)
	})
}
