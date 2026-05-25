package ui

import (
	"testing"
)

func TestParseFilePaths(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "empty input",
			input:    "",
			expected: 0,
		},
		{
			name:     "no file paths",
			input:    "Just some text without any file references",
			expected: 0,
		},
		{
			name:     "single Go file",
			input:    "Modified: internal/ui/app.go:42",
			expected: 1,
		},
		{
			name:     "multiple files",
			input:    "Changed:\n- internal/ui/app.go\n- internal/ui/chat.go\n- internal/ui/sidebar.go",
			expected: 3,
		},
		{
			name:     "TypeScript file",
			input:    "Updated src/components/App.tsx",
			expected: 1,
		},
		{
			name:     "markdown file",
			input:    "See README.md for details",
			expected: 1,
		},
		{
			name:     "no duplicates",
			input:    "app.go\napp.go\napp.go",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseFilePaths(tt.input)
			if len(result) != tt.expected {
				t.Errorf("parseFilePaths(%q) = %v (len=%d), want %d paths",
					tt.input, result, len(result), tt.expected)
			}
		})
	}
}
