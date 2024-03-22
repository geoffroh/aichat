package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchRespFromOpenAI(t *testing.T) {

	// Mock OpenAI API response
	resp := `{"choices":[{"message":{"content":"Hello human"}}]}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(resp))
	}))
	defer ts.Close()

	// Set mock API url
	origUrl := os.Getenv("OPENAI_URL")
	os.Setenv("OPENAI_URL", ts.URL)
	defer os.Setenv("OPENAI_URL", origUrl)

	// Call function
	msg := fetchRespFromOpenAI(nil, "Hi")

	// Assert
	assert.Equal(t, "Hello human", msg)
}
