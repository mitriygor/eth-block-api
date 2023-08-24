package cron_job

import (
	"testing"
	"time"
)

func Test_GetInterval(t *testing.T) {
	tests := []struct {
		input    int
		expected time.Duration
	}{
		{0, 100 * time.Second},
		{5, 5 * time.Second},
		{10, 10 * time.Second},
		{-5, -5 * time.Second},
	}

	for _, test := range tests {
		got := GetInterval(test.input)
		if got != test.expected {
			t.Errorf("For input %d, expected %v, but got %v", test.input, test.expected, got)
		}
	}
}
