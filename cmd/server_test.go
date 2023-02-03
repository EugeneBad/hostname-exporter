package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerHandler(t *testing.T) {
	os.Setenv("HOSTNAME", "test_hostname")
	test_srv := NewServer(context.Background(), "9393")
	s := httptest.NewServer(http.HandlerFunc(test_srv.ServerHandler))
	t.Cleanup(s.Close)

	invalidReq, _ := http.NewRequest("POST", s.URL+"/hostname", nil)
	resp, _ := http.DefaultClient.Do(invalidReq)
	assert.Equal(t, 400, resp.StatusCode)

	// Valid Requests
	validReq, _ := http.NewRequest("GET", s.URL+"/metric", nil)
	resp, _ = http.DefaultClient.Do(validReq)
	assert.Equal(t, 200, resp.StatusCode)

	validReq, _ = http.NewRequest("GET", s.URL+"/hostname", nil)
	resp, _ = http.DefaultClient.Do(validReq)
	assert.Equal(t, 200, resp.StatusCode)

	respObject := &struct {
		Timestamp string `json:"timestamp"`
		Hostname  string `json:"hostname"`
	}{}

	json.NewDecoder(resp.Body).Decode(respObject)
	assert.NotEmpty(t, respObject.Hostname)
	assert.Equal(t, "test_hostname", respObject.Hostname)
	assert.NotEmpty(t, respObject.Timestamp)

}
