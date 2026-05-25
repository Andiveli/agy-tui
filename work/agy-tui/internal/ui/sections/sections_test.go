package sections

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/samael/agy-tui/internal/ui/kit"
)

var testStyles = kit.DeriveStyles(kit.DefaultCatppuccinMocha())

// mustUpdate is a helper that casts a Section to UpdatableSection and calls Update.
func mustUpdate(t *testing.T, sec kit.Section, msg tea.Msg) {
	t.Helper()
	upd, ok := sec.(kit.UpdatableSection)
	if !ok {
		t.Fatalf("%T does not implement UpdatableSection", sec)
	}
	upd.Update(msg)
}

type sectionRenderCase struct {
	name   string
	sec    kit.Section
	update func() tea.Msg
	check  func(t *testing.T, sec kit.Section)
}

func TestSections_Render(t *testing.T) {
	tests := []sectionRenderCase{
		{
			name: "SessionInfo renders with correct height",
			sec:  NewSessionInfo(),
			check: func(t *testing.T, sec kit.Section) {
				if v := sec.View(30, testStyles); v == "" {
					t.Error("empty view")
				}
				if h := sec.Height(); h != 3 {
					t.Errorf("Height = %d, want 3", h)
				}
			},
		},
		{
			name: "SubAgentProgress renders",
			sec:  NewSubAgentProgress(),
			check: func(t *testing.T, sec kit.Section) {
				if v := sec.View(30, testStyles); v == "" {
					t.Error("empty view")
				}
				if h := sec.Height(); h < 2 {
					t.Errorf("Height = %d, want >= 2", h)
				}
			},
		},
		{
			name: "MCPStatus renders",
			sec:  NewMCPStatus(),
			check: func(t *testing.T, sec kit.Section) {
				if v := sec.View(30, testStyles); v == "" {
					t.Error("empty view")
				}
			},
		},
		{
			name: "MCPStatus updates",
			sec:  NewMCPStatus(),
			update: func() tea.Msg {
				return kit.MCPStatusMsg{Connected: true, Status: "Connected"}
			},
			check: func(t *testing.T, sec kit.Section) {
				if v := sec.View(30, testStyles); v == "" {
					t.Error("empty after update")
				}
			},
		},
		{
			name: "LSPStatus renders",
			sec:  NewLSPStatus(),
			check: func(t *testing.T, sec kit.Section) {
				if v := sec.View(30, testStyles); v == "" {
					t.Error("empty view")
				}
			},
		},
		{
			name: "ContextPanel renders",
			sec:  NewContextPanel(),
			check: func(t *testing.T, sec kit.Section) {
				if v := sec.View(30, testStyles); v == "" {
					t.Error("empty view")
				}
			},
		},
		{
			name: "FileChanges renders",
			sec:  NewFileChanges(),
			check: func(t *testing.T, sec kit.Section) {
				if v := sec.View(30, testStyles); v == "" {
					t.Error("empty view")
				}
			},
		},
		{
			name: "FileChanges updates",
			sec:  NewFileChanges(),
			update: func() tea.Msg {
				return kit.FileChangedMsg{Path: "test.go", Action: "modified"}
			},
			check: func(t *testing.T, sec kit.Section) {
				v := sec.View(30, testStyles)
				if v == "" {
					t.Error("empty after update")
				}
				if h := sec.Height(); h < 3 {
					t.Errorf("Height after update = %d, want >= 3", h)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.update != nil {
				if upd, ok := tt.sec.(kit.UpdatableSection); ok {
					upd.Update(tt.update())
				}
			}
			tt.check(t, tt.sec)
		})
	}
}

func TestSessionInfo_Update(t *testing.T) {
	sec := NewSessionInfo()
	mustUpdate(t, sec, kit.SessionChangedMsg{Name: "fix-auth", Context: "agy-tui"})

	v := sec.View(30, testStyles)
	if v == "" {
		t.Error("empty view after update")
	}
	if h := sec.Height(); h != 3 {
		t.Errorf("Height = %d, want 3", h)
	}
}

func TestSessionInfo_UpdateEmpty(t *testing.T) {
	sec := NewSessionInfo()
	mustUpdate(t, sec, kit.SessionChangedMsg{Name: "", Context: ""})
	v := sec.View(30, testStyles)
	if v == "" {
		t.Error("empty view after empty update")
	}
}

func TestSubAgentProgress_MultipleAgents(t *testing.T) {
	sec := NewSubAgentProgress()

	mustUpdate(t, sec, kit.ProgressMsg{SubAgent: "agy", Status: "running", Progress: 0})
	mustUpdate(t, sec, kit.ProgressMsg{SubAgent: "linter", Status: "completed", Progress: 100})
	mustUpdate(t, sec, kit.ProgressMsg{SubAgent: "tester", Status: "failed", Progress: 0})

	v := sec.View(30, testStyles)
	if v == "" {
		t.Error("empty view with multiple agents")
	}
	expectedH := 2 + 3
	if h := sec.Height(); h != expectedH {
		t.Errorf("Height = %d, want %d", h, expectedH)
	}
}

func TestContextPanel_MaxItems(t *testing.T) {
	sec := NewContextPanel()

	for i := 0; i < 10; i++ {
		mustUpdate(t, sec, kit.SessionChangedMsg{Context: "file.go"})
	}

	v := sec.View(30, testStyles)
	if v == "" {
		t.Error("empty view with items")
	}
	expectedH := 2 + 5
	if h := sec.Height(); h != expectedH {
		t.Errorf("Height = %d, want %d", h, expectedH)
	}
}

func TestFileChanges_Actions(t *testing.T) {
	sec := NewFileChanges()

	mustUpdate(t, sec, kit.FileChangedMsg{Path: "new.go", Action: "created"})
	mustUpdate(t, sec, kit.FileChangedMsg{Path: "old.go", Action: "deleted"})
	mustUpdate(t, sec, kit.FileChangedMsg{Path: "mod.go", Action: "modified"})

	v := sec.View(30, testStyles)
	if v == "" {
		t.Error("empty view with file changes")
	}
	expectedH := 2 + 3
	if h := sec.Height(); h != expectedH {
		t.Errorf("Height = %d, want %d", h, expectedH)
	}
}

func TestFileChanges_MaxFiles(t *testing.T) {
	sec := NewFileChanges()

	for i := 0; i < 15; i++ {
		mustUpdate(t, sec, kit.FileChangedMsg{Path: "f.go", Action: "modified"})
	}

	v := sec.View(30, testStyles)
	if v == "" {
		t.Error("empty view")
	}
	h := sec.Height()
	if h < 2 || h > 8 {
		t.Errorf("Height = %d, want between 2 and 8", h)
	}
}

func TestMCPStatus_Toggle(t *testing.T) {
	sec := NewMCPStatus()

	mustUpdate(t, sec, kit.MCPStatusMsg{Connected: false, Status: "Disconnected"})
	v1 := sec.View(30, testStyles)
	if v1 == "" {
		t.Error("empty view when disconnected")
	}

	mustUpdate(t, sec, kit.MCPStatusMsg{Connected: true, Status: "MCP configured"})
	v2 := sec.View(30, testStyles)
	if v2 == "" {
		t.Error("empty view when connected")
	}
	if v1 == v2 {
		t.Error("disconnected and connected views should differ")
	}
}

func TestLSPStatus_Update(t *testing.T) {
	sec := NewLSPStatus()

	mustUpdate(t, sec, kit.LSPStatusMsg{Connected: true, Status: "gopls, rust-analyzer"})
	v := sec.View(30, testStyles)
	if v == "" {
		t.Error("empty view after update")
	}
	if h := sec.Height(); h != 2 {
		t.Errorf("Height = %d, want 2", h)
	}
}
