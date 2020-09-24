package shortener

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var server = Server{URLs: URLDictionary{
	"go":  "https://google.de/",
	"hmm": "https://wikipedia.org/",
	"bin": "https://httpbin.org/",
}}

func TestGetURL(t *testing.T) {
	t.Run("Getting the given URLs", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/go", nil)
		response := httptest.NewRecorder()

		gotDest, err := server.GetURL(response, request)
		wantDest := "https://google.de/"

		assertNoError(t, err)

		got := response.Result().StatusCode
		want := 301 // redirect status code

		assertInt(t, got, want)
		assertStrings(t, gotDest, wantDest)
	})

	t.Run("Requesting a non-existing shortcut", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/lol", nil)
		response := httptest.NewRecorder()

		_, err := server.GetURL(response, request)

		assertError(t, err)

		got := response.Result().StatusCode
		want := 404

		assertInt(t, got, want)
	})
}

func TestLookup(t *testing.T) {
	t.Run("Lookup and find", func(t *testing.T) {
		got, err := server.URLs.Lookup("bin")
		want := "https://httpbin.org/"

		assertNoError(t, err)
		assertSuccessfulLookup(t, got, want)
	})

	t.Run("Lookup non existing key", func(t *testing.T) {
		url := "unknown"
		_, err := server.URLs.Lookup(url)
		want := fmt.Sprintf("could not find the proper URL for %s (unknown)", url)

		assertError(t, err)
		assertStrings(t, err.Error(), want)
	})
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error but didn't get any.")
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal("have error but didn't expect any!")
	}
}

func assertStrings(t *testing.T, got string, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertSuccessfulLookup(t *testing.T, got string, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertInt(t *testing.T, got int, want int) {
	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}
