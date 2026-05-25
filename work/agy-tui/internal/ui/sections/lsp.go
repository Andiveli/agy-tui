package sections

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/samael/agy-tui/internal/ui/kit"
)

// LSPStatus displays the LSP connection state.
type LSPStatus struct {
	connected bool
	status    string
}

func NewLSPStatus() *LSPStatus {
	return &LSPStatus{status: "Disconnected"}
}

func (l *LSPStatus) Name() string { return "LSP" }

func (l *LSPStatus) View(width int, styles kit.Styles) string {
	header := styles.SectionHeader.Render(" " + l.Name() + " ")
	color := lipgloss.Color("#f38ba8")
	if l.connected {
		color = lipgloss.Color("#a6e3a1")
	}
	statusLine := lipgloss.NewStyle().Foreground(color).Render("● " + l.status)
	return lipgloss.NewStyle().Width(width).
		Render(lipgloss.JoinVertical(lipgloss.Top, header, styles.SectionContent.Render(statusLine)))
}

func (l *LSPStatus) Height() int { return 2 }

func (l *LSPStatus) Update(msg tea.Msg) (kit.Section, tea.Cmd) {
	switch v := msg.(type) {
	case kit.LSPStatusMsg:
		l.connected = v.Connected
		l.status = v.Status
	}
	return l, nil
}
