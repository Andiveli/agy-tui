package ui

import (
	"testing"

	"github.com/samael/agy-tui/internal/ui/kit"
)

func TestStatusBar_New(t *testing.T) {
	sb := NewStatusBar(testStyles)
	if sb == nil {
		t.Fatal("NewStatusBar returned nil")
	}
}

func TestStatusBar_View(t *testing.T) {
	sb := NewStatusBar(testStyles)
	v := sb.View(80)
	if v == "" {
		t.Error("empty view with width 80")
	}
}

func TestStatusBar_Narrow(t *testing.T) {
	sb := NewStatusBar(testStyles)
	v := sb.View(10)
	if v == "" {
		t.Error("empty view with narrow width")
	}
}

func TestStatusBar_Updates(t *testing.T) {
	sb := NewStatusBar(testStyles)

	sb.Update(kit.SessionChangedMsg{Name: "test-session"})
	v := sb.View(80)
	if v == "" {
		t.Error("empty view after session update")
	}

	sb.Update(kit.ProgressMsg{SubAgent: "agy", Status: "completed", Progress: 100})
	v2 := sb.View(80)
	if v2 == "" {
		t.Error("empty view after progress update")
	}

	sb.Update(kit.MCPStatusMsg{Connected: true, Status: "MCP configured"})
	v3 := sb.View(80)
	if v3 == "" {
		t.Error("empty view after MCP update")
	}

	sb.Update(kit.LSPStatusMsg{Connected: true, Status: "gopls"})
	v4 := sb.View(80)
	if v4 == "" {
		t.Error("empty view after LSP update")
	}
}
