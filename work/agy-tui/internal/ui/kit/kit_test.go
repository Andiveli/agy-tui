package kit

import (
	"testing"
)

func TestDeriveStyles_DefaultTheme(t *testing.T) {
	styles := DeriveStyles(DefaultCatppuccinMocha())

	tests := []struct {
		name   string
		style  interface{ Render(a ...string) string }
		expect string // non-empty means we just verify it renders
	}{
		{"Title", styles.Title, ""},
		{"Border", styles.Border, ""},
		{"Error", styles.Error, ""},
		{"Dimmed", styles.Dimmed, ""},
		{"UserMessage", styles.UserMessage, ""},
		{"AgentMessage", styles.AgentMessage, ""},
		{"SectionHeader", styles.SectionHeader, ""},
		{"SectionContent", styles.SectionContent, ""},
		{"Spinner", styles.Spinner, ""},
		{"ContextItem", styles.ContextItem, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.style.Render("hello")
			if result == "" {
				t.Errorf("%s rendered empty string", tt.name)
			}
		})
	}
}

func TestThemeFields_Count(t *testing.T) {
	fields := ThemeFields(DefaultCatppuccinMocha())
	if len(fields) != 16 {
		t.Errorf("expected 16 theme fields, got %d", len(fields))
	}
}

func TestThemeChangedMsgCmd(t *testing.T) {
	cmd := ThemeChangedMsgCmd()
	msg := cmd()
	if _, ok := msg.(ThemeChangedMsg); !ok {
		t.Errorf("expected ThemeChangedMsg, got %T", msg)
	}
}
