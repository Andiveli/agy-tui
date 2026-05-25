package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/samael/agy-tui/internal/ui/kit"
)

// StatusBar is a minimal bottom bar showing session info and key hints.
type StatusBar struct {
	sessionName string
	agyOnline   bool
	mcpOnline   bool
	lspOnline   bool
	styles      kit.Styles
}

func NewStatusBar(styles kit.Styles) *StatusBar {
	return &StatusBar{styles: styles}
}

func (s *StatusBar) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case kit.SessionChangedMsg:
		if msg.Name != "" {
			s.sessionName = msg.Name
		}
	case kit.ProgressMsg:
		// Consider agy "online" if we get any progress update
		s.agyOnline = true
	case kit.MCPStatusMsg:
		s.mcpOnline = msg.Connected
	case kit.LSPStatusMsg:
		s.lspOnline = msg.Connected
	}
}

func (s *StatusBar) View(width int) string {
	if width <= 0 {
		return ""
	}

	// Left: session name
	left := s.styles.StatusAccent.Render(" " + s.sessionName + " ")
	if s.sessionName == "" {
		left = s.styles.StatusAccent.Render(" session ")
	}

	// Center: key hints
	hints := []string{
		s.styles.KeyHint.Render("Ctrl+B ") + s.styles.KeyBinding.Render("sidebar"),
		s.styles.KeyHint.Render("Ctrl+T ") + s.styles.KeyBinding.Render("theme"),
		s.styles.KeyHint.Render("Ctrl+C ") + s.styles.KeyBinding.Render("quit"),
	}
	center := lipgloss.JoinHorizontal(lipgloss.Center, hints...)

	// Right: connection indicators
	agyIcon := s.dot(false)
	if s.agyOnline {
		agyIcon = s.dot(true)
	}
	right := fmt.Sprintf("agy:%s mcp:%s lsp:%s", agyIcon, s.dot(s.mcpOnline), s.dot(s.lspOnline))
	right = s.styles.KeyHint.Render(right)

	// Layout: left | center | right with flex spacing
	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		left,
		lipgloss.NewStyle().Width(width-lipgloss.Width(left)-lipgloss.Width(right)).Render(center),
		right,
	)

	return s.styles.StatusBar.Width(width).Render(bar)
}

func (s *StatusBar) dot(on bool) string {
	if on {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1")).Render("●")
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#585b70")).Render("●")
}
