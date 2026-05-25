package sections

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/samael/agy-tui/internal/ui/kit"
)

// SessionInfo displays the current session name and context.
type SessionInfo struct {
	name      string
	context   string
	convCount int
	convID    string
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
	var extras string
	if s.convCount > 0 {
		extras = styles.Dimmed.Render(" conversations: " + fmt.Sprint(s.convCount))
	}
	if s.convID != "" {
		// Show short form of the conversation ID
		id := s.convID
		if len(id) > 8 {
			id = id[:8] + "…"
		}
		extras += styles.Dimmed.Render(" ID: " + id)
	}
	section := lipgloss.JoinVertical(lipgloss.Top, header, name, ctx, extras)
	return lipgloss.NewStyle().Width(width).Render(section)
}

func (s *SessionInfo) Height() int {
	h := 3
	if s.convCount > 0 {
		h++
	}
	if s.convID != "" {
		h++
	}
	return h
}

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
		s.convCount = v.ConvCount
		if v.ConvID != "" {
			s.convID = v.ConvID
		}
	}
	return s, nil
}
