package shortener

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestGetURL(t *testing.T) {
	db, path := setupDatabase(t)
	shortServer := ShortServer{DB: db, URL: "easy.xyz"}

	defer db.Close()
	defer os.RemoveAll(path)

	t.Run("Getting the given URLs", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/test", nil)
		response := httptest.NewRecorder()

		shortServer.GetURL(response, request)

		actual := response.Result().StatusCode
		expected := http.StatusPermanentRedirect // redirect status code

		assert.Equal(t, expected, actual)
	})

	t.Run("Requesting a non-existing shortcut", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/unknown", nil)
		response := httptest.NewRecorder()

		shortServer.GetURL(response, request)

		actual := response.Result().StatusCode
		expected := http.StatusNotFound
		assert.Equal(t, expected, actual)
	})
}

func TestLookup(t *testing.T) {
	db, path := setupDatabase(t)
	shortServer := ShortServer{DB: db, URL: "easy.xyz"}

	defer db.Close()
	defer os.RemoveAll(path)

	t.Run("Lookup and find", func(t *testing.T) {
		expected := "https://test.local/"
		actual, err := shortServer.Lookup("/test")
		assert.Equal(t, expected, actual)
		assert.NoError(t, err)
	})

	t.Run("Lookup non existing key", func(t *testing.T) {
		url := "unknown"
		_, err := shortServer.Lookup(url)
		expected := fmt.Sprintf("could not find the proper URL for %s (unknown)", url)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), expected)
	})
}

func TestAddURL(t *testing.T) {
	db, path := setupDatabase(t)
	shortServer := ShortServer{DB: db, URL: "easy.xyz"}

	defer db.Close()
	defer os.RemoveAll(path)

	t.Run("Adding an URL", func(t *testing.T) {
		urlPair := URLPair{Shorthand: "foo", Target: "https://bar.local/"}
		jsonPair, err := json.Marshal(urlPair)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/add", bytes.NewBuffer(jsonPair))
		response := httptest.NewRecorder()
		assert.NoError(t, err)

		shortServer.AddURL(response, request)
		assert.NoError(t, err)

		expected := http.StatusOK
		actual := response.Result().StatusCode
		assert.Equal(t, expected, actual)
	})

	t.Run("Adding an existing URL", func(t *testing.T) {
		urlPair := URLPair{Shorthand: "test", Target: "https://test.local/"}
		jsonPair, err := json.Marshal(urlPair)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/add", bytes.NewBuffer(jsonPair))
		response := httptest.NewRecorder()
		assert.NoError(t, err)

		shortServer.AddURL(response, request)

		expected := http.StatusConflict
		actual := response.Result().StatusCode
		assert.Equal(t, expected, actual)
	})

	t.Run("Empty fields", func(t *testing.T) {

	})
}

func TestAdd(t *testing.T) {
	db, path := setupDatabase(t)
	shortServer := ShortServer{DB: db, URL: "easy.xyz"}

	defer db.Close()
	defer os.RemoveAll(path)
	t.Run("Adding a new Pair", func(t *testing.T) {
		urlPair := URLPair{Shorthand: "foo", Target: "https://bar.local/"}

		err := shortServer.Add(urlPair)
		assert.NoError(t, err)
		actual, err := shortServer.Lookup(urlPair.Shorthand)
		assert.NoError(t, err)
		assert.Equal(t, urlPair.Target, actual)
	})

	t.Run("Adding an existing Pair", func(t *testing.T) {
		urlPair := URLPair{Shorthand: "test", Target: "https://test.local/"}
		err := shortServer.Add(urlPair)
		assert.Error(t, err)
	})
}

func setupDatabase(t *testing.T) (*bolt.DB, string) {
	t.Helper()
	// using tmp files for integration testing
	dir, err := ioutil.TempDir("", "database")
	assert.NoError(t, err)

	tmpDB := filepath.Join(dir, "test.db")

	db, err := bolt.Open(tmpDB, 0600, nil)
	assert.NoError(t, err)

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("urls"))
		if err != nil {
			return err
		}
		err = bucket.Put([]byte("test"), []byte("https://test.local/"))
		if err != nil {
			return err
		}
		return nil
	})

	assert.NoError(t, err)

	return db, dir
}
