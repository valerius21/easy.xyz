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

		actualDest, err := shortServer.GetURL(response, request)
		expectedDest := "http://test.local/"

		assert.Nil(t, err)

		actual := response.Result().StatusCode
		expected := 308 // redirect status code

		assert.Equal(t, expected, actual)
		assert.Equal(t, expectedDest, actualDest)
	})

	t.Run("Requesting a non-existing shortcut", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/unknown", nil)
		response := httptest.NewRecorder()

		_, err := shortServer.GetURL(response, request)

		assert.NotNil(t, err)

		actual := response.Result().StatusCode
		expected := 404
		assert.Equal(t, expected, actual)
	})
}

func TestLookup(t *testing.T) {
	t.Run("Lookup and find", func(t *testing.T) {
		actual, err := shortServer.URLs.Lookup("test")
		expected := "http://test.local/"

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Lookup non existing key", func(t *testing.T) {
		url := "unknown"
		_, err := shortServer.URLs.Lookup(url)
		expected := fmt.Sprintf("could not find the proper URL for %s (unknown)", url)

		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), expected)
	})
}
