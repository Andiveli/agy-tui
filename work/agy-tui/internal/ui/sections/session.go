package sections

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/samael/agy-tui/internal/ui/kit"
)

// SessionInfo displays the current session name and context.
type SessionInfo struct {
	name    string
	context string
}

// NewSessionInfo creates a SessionInfo with default placeholder values.
func NewSessionInfo() *SessionInfo {
	return &SessionInfo{
		name:    "default",
		context: "no workspace",
	}
}

func (s *SessionInfo) Name() string { return "Session" }

func (s *SessionInfo) View(width int, styles kit.Styles) string {
	header := styles.SectionHeader.Render(" " + s.Name() + " ")
	name := styles.SectionContent.Render("  " + s.name)
	ctx := styles.Dimmed.Render(" " + s.context)
	section := lipgloss.JoinVertical(lipgloss.Top, header, name, ctx)
	return lipgloss.NewStyle().Width(width).Render(section)
}

func (s *SessionInfo) Height() int { return 3 }

// Update implements kit.UpdatableSection.
func (s *SessionInfo) Update(msg tea.Msg) (kit.Section, tea.Cmd) {
	switch v := msg.(type) {
	case kit.SessionChangedMsg:
		if v.Name != "" {
			s.name = v.Name
		}
		if v.Context != "" {
			s.context = v.Context
		}
	}
	return s, nil
}
