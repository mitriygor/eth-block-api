package eth_block_helper

import (
	"testing"
)

func Test_IsInt(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123", true},
		{"-123", true},
		{"abc", false},
		{"12.3", false},
	}

	for _, test := range tests {
		got := IsInt(test.input)
		if got != test.expected {
			t.Errorf("IsInt(%s) = %v; want %v", test.input, got, test.expected)
		}
	}
}

func Test_IsHex(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"0x1f4", true},
		{"0x1F4", true},
		{"0xZZZ", false},
		{"1f4", false},
		{"0x1f4a5f9c", true},
	}

	for _, test := range tests {
		got := IsHex(test.input)
		if got != test.expected {
			t.Errorf("IsHex(%s) = %v; want %v", test.input, got, test.expected)
		}
	}
}

func Test_StringToInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"123", 123},
		{"-123", -123},
		{"abc", -1},
		{"", -1},
	}

	for _, test := range tests {
		got := StringToInt(test.input)
		if got != test.expected {
			t.Errorf("StringToInt(%s) = %d; want %d", test.input, got, test.expected)
		}
	}
}

func Test_IntToHex(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{123, "0x7b"},
		{0, "0x0"},
	}

	for _, test := range tests {
		got := IntToHex(test.input)
		if got != test.expected {
			t.Errorf("IntToHex(%d) = %s; want %s", test.input, got, test.expected)
		}
	}
}
