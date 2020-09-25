package shortener

import (
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

	defer assert.NoError(t, db.Close())
	defer assert.NoError(t, os.RemoveAll(path))

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

	defer assert.NoError(t, db.Close())
	defer assert.NoError(t, os.RemoveAll(path))

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
