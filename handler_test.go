package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFooHandler(t *testing.T) {
	mux := setupRouter()
	srv := httptest.NewServer(mux)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/h1")

	println("I am here!!!!!!!!1")

	require.NoErrorf(t, err, "Error occured! Err: %s", err.Error())
	require.Truef(t, (resp.StatusCode <= 399), "Status Code recv %d", resp.StatusCode)

	r, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	require.NoError(t, err, "Error occured")

	assert.Equal(t, "foo", string(r))
}

func TestBarHandler(t *testing.T) {
	mux := setupRouter()
	srv := httptest.NewServer(mux)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/h2")

	require.NoErrorf(t, err, "Error occured! Err: %s", err.Error())
	require.Truef(t, (resp.StatusCode <= 399), "Status Code recv %d", resp.StatusCode)

	r, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	require.NoError(t, err, "Error occured")

	assert.Equal(t, "bar", string(r))
}
