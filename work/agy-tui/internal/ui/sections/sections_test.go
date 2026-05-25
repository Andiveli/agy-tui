package sections

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/samael/agy-tui/internal/ui/kit"
)

var testStyles = kit.DeriveStyles(kit.DefaultCatppuccinMocha())

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
