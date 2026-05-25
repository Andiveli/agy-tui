package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/samael/agy-tui/internal/ui/kit"
	"github.com/samael/agy-tui/internal/ui/sections"
)

// Sidebar is the right-side information panel.
type Sidebar struct {
	width    int
	height   int
	visible  bool
	sections []kit.Section
	styles   kit.Styles
}

func NewSidebar(styles kit.Styles) *Sidebar {
	return &Sidebar{
		visible: true,
		styles:  styles,
		sections: []kit.Section{
			sections.NewSessionInfo(),
			sections.NewSubAgentProgress(),
			sections.NewMCPStatus(),
			sections.NewLSPStatus(),
			sections.NewContextPanel(),
			sections.NewFileChanges(),
		},
	}
}

func (s *Sidebar) Init() tea.Cmd {
	// Start spinner tick for progress section
	return sections.NewSubAgentProgress().Init()
}

func (s *Sidebar) Update(msg tea.Msg) (*Sidebar, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = 30
		s.height = msg.Height
		if msg.Width < 100 {
			s.visible = false
		}
	}

	for i, sec := range s.sections {
		if upd, ok := sec.(kit.UpdatableSection); ok {
			updated, cmd := upd.Update(msg)
			s.sections[i] = updated
			if cmd != nil {
				return s, cmd
			}
		}
	}
	return s, nil
}

func (s *Sidebar) ToggleVisibility() { s.visible = !s.visible }

func (s *Sidebar) View() string {
	if !s.visible || s.width == 0 {
		return ""
	}

	var rendered []string
	remaining := s.height - 2

	for _, sec := range s.sections {
		h := sec.Height()
		if h > remaining {
			rendered = append(rendered, s.styles.Dimmed.Render(" ⋯ "+sec.Name()))
			continue
		}
		rendered = append(rendered, sec.View(s.width-2, s.styles))
		remaining -= h
	}

	body := lipgloss.JoinVertical(lipgloss.Top, rendered...)
	return lipgloss.NewStyle().
		Width(s.width).Height(s.height).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(s.styles.Dimmed.GetForeground()).
		Render(body)
}
