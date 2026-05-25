package ui

import (
	"testing"

	"github.com/samael/agy-tui/internal/backend"
	"github.com/samael/agy-tui/internal/ui/kit"
)

var testStyles = kit.DeriveStyles(kit.DefaultCatppuccinMocha())

func testClient() *backend.Client {
	return backend.NewClient()
}

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
		{
			name:     "JSON config file",
			input:    "Updated plugin/plugin.json",
			expected: 1,
		},
		{
			name:     "Go sum file",
			input:    "Updated go.sum",
			expected: 1,
		},
		{
			name:     "Go mod file",
			input:    "Updated go.mod",
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

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "How do I fix the authentication bug",
			expected: "how-do-i-fix",
		},
		{
			input:    "hello",
			expected: "hello",
		},
		{
			input:    "",
			expected: "session",
		},
		{
			input:    "  spaces   everywhere  ",
			expected: "spaces-everywhere",
		},
		{
			input:    "UPPERCASE Things",
			expected: "uppercase-things",
		},
		{
			input:    "special!@#characters???",
			expected: "special!@#characters",
		},
		{
			input:    "a b c d e f g",
			expected: "a-b-c-d",
		},
		{
			input:    "123 numbers 456",
			expected: "123-numbers-456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := slugify(tt.input)
			if result != tt.expected {
				t.Errorf("slugify(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNewChatModel_Defaults(t *testing.T) {
	styles := testStyles
	client := testClient()
	cm := NewChatModel(styles, client)

	if cm == nil {
		t.Fatal("NewChatModel returned nil")
	}
	if cm.sessionMgr == nil {
		t.Error("sessionMgr not initialized")
	}
	if cm.loading {
		t.Error("new model should not be loading")
	}
	if cm.promptCount != 0 {
		t.Errorf("promptCount = %d, want 0", cm.promptCount)
	}
	if cm.streamingText != "" {
		t.Errorf("streamingText = %q, want empty", cm.streamingText)
	}
	if cm.streamReader != nil {
		t.Error("streamReader should be nil initially")
	}
	if cm.streamScanner != nil {
		t.Error("streamScanner should be nil initially")
	}
}

func TestNewChatModel_SessionManager(t *testing.T) {
	cm := NewChatModel(testStyles, testClient())
	if cm.sessionMgr.CurrentSession() != nil {
		t.Error("session should not be started yet")
	}

	// StartSession should work
	cm.sessionMgr.StartSession()
	if cm.sessionMgr.CurrentSession() == nil {
		t.Error("session should be active after StartSession")
	}
}
