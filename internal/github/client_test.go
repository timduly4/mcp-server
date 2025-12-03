package github

import (
	"testing"
)

func TestGetStringValue(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		expected string
	}{
		{
			name:     "nil pointer",
			input:    nil,
			expected: "",
		},
		{
			name:     "non-nil pointer",
			input:    stringPtr("test"),
			expected: "test",
		},
		{
			name:     "empty string",
			input:    stringPtr(""),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStringValue(tt.input)
			if result != tt.expected {
				t.Errorf("getStringValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetIntValue(t *testing.T) {
	tests := []struct {
		name     string
		input    *int
		expected int
	}{
		{
			name:     "nil pointer",
			input:    nil,
			expected: 0,
		},
		{
			name:     "non-nil pointer",
			input:    intPtr(42),
			expected: 42,
		},
		{
			name:     "zero value",
			input:    intPtr(0),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getIntValue(tt.input)
			if result != tt.expected {
				t.Errorf("getIntValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Helper functions for tests
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
