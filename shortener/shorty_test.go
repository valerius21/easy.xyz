package shortener

import (
	"fmt"
	"testing"
)

var urls = URLDictionary{
	"go":  "https://google.de/",
	"hmm": "https://wikipedia.org/",
	"bin": "https://httpbin.org/",
}

func TestGETUrl(t *testing.T) {

	//t.Run("Getting the given URLs", func(t *testing.T) {
	//	request, _ := http.NewRequest(http.MethodGet, "/go", nil)
	//	response := httptest.NewRecorder()
	//
	//	ShortServer(response, request)
	//
	//	got := response.Result().StatusCode
	//	want := 302 // redirect status code
	//
	//	if got != want {
	//		t.Errorf("got %d, wanted %d", got, want)
	//	}
	//})
}

func TestLookup(t *testing.T) {
	t.Run("Lookup and find", func(t *testing.T) {
		got, err := urls.LookupURL("bin")
		want := "https://httpbin.org/"

		assertNoError(t, err)
		assertSuccessfulLookup(t, got, want)
	})

	t.Run("Lookup non existing key", func(t *testing.T) {
		url := "unknown"
		_, err := urls.LookupURL(url)
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
