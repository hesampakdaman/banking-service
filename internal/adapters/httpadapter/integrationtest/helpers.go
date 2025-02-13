package integrationtest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func postJSON(t *testing.T, url string, body interface{}) *http.Response {
	t.Helper()

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("failed to send POST request: %v", err)
	}

	return resp
}

func getJSON(t *testing.T, url string) *http.Response {
	t.Helper()

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("failed to send GET request: %v", err)
	}

	return resp
}

func parseJSON(t *testing.T, resp *http.Response, target interface{}) {
	t.Helper()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	defer resp.Body.Close()

	err = json.Unmarshal(body, target)
	if err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
}
