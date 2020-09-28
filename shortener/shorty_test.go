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
	shortServer := ShortServer{DB: db}

	defer db.Close()
	defer os.RemoveAll(path)

	t.Run("Getting the given URLs", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/test", nil)
		response := httptest.NewRecorder()

		actualDest, err := shortServer.GetURL(response, request)
		expectedDest := "https://test.local/"

		assert.NoError(t, err)

		actual := response.Result().StatusCode
		expected := 308 // redirect status code

		assert.Equal(t, expected, actual)
		assert.Equal(t, expectedDest, actualDest)
	})

	t.Run("Requesting a non-existing shortcut", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/unknown", nil)
		response := httptest.NewRecorder()

		_, err := shortServer.GetURL(response, request)

		assert.Error(t, err)

		actual := response.Result().StatusCode
		expected := 404
		assert.Equal(t, expected, actual)
	})
}

func TestLookup(t *testing.T) {
	db, path := setupDatabase(t)
	shortServer := ShortServer{DB: db}

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
	shortServer := ShortServer{DB: db}

	defer db.Close()
	defer os.RemoveAll(path)

	t.Run("Adding an URL", func(t *testing.T) {
		urlPair := URLPair{Shorthand: "foo", Target: "https://bar.local/"}
		jsonPair, err := json.Marshal(urlPair)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/add", bytes.NewBuffer(jsonPair))
		response := httptest.NewRecorder()
		assert.NoError(t, err)

		message, err := shortServer.AddURL(response, request)
		assert.NoError(t, err)

		expected := 200
		actual := response.Result().StatusCode
		assert.Equal(t, expected, actual)
		assert.Equal(t, "OK", message)

		verifyRequest, err := http.NewRequest(http.MethodGet,
			fmt.Sprintf("/%s", urlPair.Shorthand), nil)
		assert.NoError(t, err)

		response = httptest.NewRecorder()
		targetURL, err := shortServer.GetURL(response, verifyRequest)
		assert.NoError(t, err)

		expected = 308
		actual = response.Result().StatusCode
		assert.Equal(t, expected, actual)
		assert.Equal(t, urlPair.Target, targetURL)
	})

	t.Run("Adding an existing URL", func(t *testing.T) {

	})

	t.Run("Empty fields", func(t *testing.T) {

	})
}

func TestAdd(t *testing.T) {
	db, path := setupDatabase(t)
	shortServer := ShortServer{DB: db}

	defer db.Close()
	defer os.RemoveAll(path)

	urlPair := URLPair{Shorthand: "foo", Target: "https://bar.local/"}

	err := shortServer.Add(urlPair)
	assert.NoError(t, err)
	actual, err := shortServer.Lookup(urlPair.Shorthand)
	assert.NoError(t, err)
	assert.Equal(t, urlPair.Target, actual)
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
