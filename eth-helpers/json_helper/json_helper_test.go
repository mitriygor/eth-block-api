package json_helper

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_PostRequest(t *testing.T) {
	tests := []struct {
		serverResp  interface{}
		wantErr     bool
		errContains string
	}{
		{map[string]interface{}{"data": "ok"}, false, ""},
		{map[string]interface{}{"not_json": "{"}, false, "invalid character"},
	}

	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(test.serverResp)
		}))

		var result map[string]interface{}
		err := PostRequest(ts.URL, map[string]interface{}{"key": "value"}, &result)
		if (err != nil) != test.wantErr {
			t.Errorf("unexpected error: got %v, wantErr %v", err, test.wantErr)
			continue
		}

		if err != nil && !bytes.Contains([]byte(err.Error()), []byte(test.errContains)) {
			t.Errorf("unexpected error: got %s, should contain %s", err.Error(), test.errContains)
		}
		ts.Close()
	}
}

func Test_HasError(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected bool
	}{
		{struct{ Error string }{"an error"}, true},
		{struct{ NoError string }{"no error field"}, false},
		{struct{ Error string }{""}, false},
	}

	for _, test := range tests {
		got := hasError(test.input)
		if got != test.expected {
			t.Errorf("hasError(%v) = %v; want %v", test.input, got, test.expected)
		}
	}
}
