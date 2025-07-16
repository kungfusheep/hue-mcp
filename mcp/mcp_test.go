package mcp

import (
	"testing"
)

func TestColorConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"red color name", "red", "#FF0000"},
		{"green color name", "green", "#00FF00"},
		{"blue color name", "blue", "#0000FF"},
		{"mixed case", "RED", "#FF0000"},
		{"hex passthrough", "#FF00FF", ""},
		{"invalid name", "notacolor", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := namedColorToHex(tt.input)
			if result != tt.expected {
				t.Errorf("namedColorToHex(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHexColorValidation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid hex", "#FF0000", true},
		{"valid hex lowercase", "#ff0000", true},
		{"missing hash", "FF0000", false},
		{"too short", "#FF00", false},
		{"too long", "#FF00000", false},
		{"invalid characters", "#GGGGGG", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidHexColor(tt.input)
			if result != tt.expected {
				t.Errorf("isValidHexColor(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}