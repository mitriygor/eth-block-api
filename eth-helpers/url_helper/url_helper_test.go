package url_helper

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func Test_GetUrl(t *testing.T) {
	tests := []struct {
		endpoint string
		version  string
		key      string
		expected string
	}{
		{"http://api.example.com", "v1", "12345", "http://api.example.com/v1/12345"},
		{"http://api.example.com", "v1", "", "http://api.example.com/v1"},
		{"http://api.example.com", "", "12345", "http://api.example.com/12345"},
		{"http://api.example.com", "", "", "http://api.example.com"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			output := GetUrl(test.endpoint, test.version, test.key)
			if output != test.expected {
				t.Fatalf("Expected %s but got %s", test.expected, output)
			}
		})
	}
}

func Test_GetRandId(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	id := GetRandId()
	if id == 0 {
		t.Fatalf("Expected a non-zero number but got %d", id)
	}

	min := 1
	max := math.MaxInt - 100
	if id < min || id > max {
		t.Fatalf("Expected a number between %d and %d but got %d", min, max, id)
	}
}
