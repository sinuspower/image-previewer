package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	internal_cache "github.com/sinuspower/image-previewer/internal/cache"
	internal_settings "github.com/sinuspower/image-previewer/internal/settings"
	"github.com/stretchr/testify/require"
)

var imageServerHandleFunc = func(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Clone()
	// copy headers
	for key, values := range header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	f, err := os.Open("test/testdata/_gopher_original_1024x504.jpg")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func TestGetPreviews(t *testing.T) {
	initVariables(t)
	log.SetOutput(ioutil.Discard)

	imageServer := httptest.NewServer(http.HandlerFunc(imageServerHandleFunc))
	defer imageServer.Close()

	previewServer := httptest.NewServer(http.HandlerFunc(fillHandler))
	defer previewServer.Close()

	var testCases = []struct { //nolint:go-lint
		name        string
		urlTemplate string
		filepath    string
	}{
		{"50x50", "%s/fill/50/50/%s/images/source.jpg", "test/testdata/gopher_50x50.jpg"},
		{"200x700", "%s/fill/200/700/%s/images/source.jpg", "test/testdata/gopher_200x700.jpg"},
		{"256x126", "%s/fill/256/126/%s/images/source.jpg", "test/testdata/gopher_256x126.jpg"},
		{"333x666", "%s/fill/333/666/%s/images/source.jpg", "test/testdata/gopher_333x666.jpg"},
		{"500x500", "%s/fill/500/500/%s/images/source.jpg", "test/testdata/gopher_500x500.jpg"},
		{"1024x252", "%s/fill/1024/252/%s/images/source.jpg", "test/testdata/gopher_1024x252.jpg"},
		{"1024x504", "%s/fill/1024/504/%s/images/source.jpg", "test/testdata/gopher_1024x504.jpg"},
		{"2000x1000", "%s/fill/2000/1000/%s/images/source.jpg", "test/testdata/gopher_2000x1000.jpg"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf(tc.urlTemplate, previewServer.URL, imageServer.URL) //nolint:go-lint
			testGetPreview(t, url, tc.filepath)                                    //nolint:go-lint
		})
	}
}

func TestProxyHeaders(t *testing.T) {
	initVariables(t)
	log.SetOutput(ioutil.Discard)

	imageServer := httptest.NewServer(http.HandlerFunc(imageServerHandleFunc))
	defer imageServer.Close()

	previewServer := httptest.NewServer(http.HandlerFunc(fillHandler))
	defer previewServer.Close()

	url := fmt.Sprintf("%s/fill/50/50/%s/images/source.jpg", previewServer.URL, imageServer.URL)
	rq, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	require.NoError(t, err)

	rq.Header.Add("Header-One", "test-header-one")
	rq.Header.Add("Header-Two", "test-header-two")

	client := new(http.Client)
	rs, err := client.Do(rq)
	require.NoError(t, err)
	defer rs.Body.Close()
	require.Equal(t, 200, rs.StatusCode)

	require.Equal(t, "test-header-one", rs.Header["Header-One"][0])
	require.Equal(t, "test-header-two", rs.Header["Header-Two"][0])

	err = cache.Clear()
	require.NoError(t, err)
}

func initVariables(t *testing.T) {
	env := environment{"8080", "5", "50", "50", "2000", "2000"}
	setEnv(env)
	settings = new(internal_settings.Settings)
	err := settings.ParseEnv()
	require.NoError(t, err)

	cache, err = internal_cache.NewCache(settings.GetCacheSize(), "cache")
	require.NoError(t, err)
}

func testGetPreview(t *testing.T, url string, filepath string) {
	rs, err := http.Get(url) //nolint:go-lint
	require.NoError(t, err)
	defer rs.Body.Close()
	require.Equal(t, http.StatusOK, rs.StatusCode)

	f, err := os.Open(filepath)
	require.NoError(t, err)

	expBytes, err := ioutil.ReadAll(f)
	require.NoError(t, err)

	actBytes, err := ioutil.ReadAll(rs.Body)
	require.NoError(t, err)

	require.Equal(t, expBytes, actBytes)

	err = cache.Clear()
	require.NoError(t, err)
}
