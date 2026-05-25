package sections

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/samael/agy-tui/internal/ui/kit"
)

// MCPStatus displays the MCP server connection state.
type MCPStatus struct {
	connected bool
	status    string
}

func NewMCPStatus() *MCPStatus {
	return &MCPStatus{status: "Disconnected"}
}

func (m *MCPStatus) Name() string { return "MCP" }

func (m *MCPStatus) View(width int, styles kit.Styles) string {
	header := styles.SectionHeader.Render(" " + m.Name() + " ")
	color := lipgloss.Color("#f38ba8")
	indicator := "●"
	if m.connected {
		color = lipgloss.Color("#a6e3a1")
	}
	statusLine := lipgloss.NewStyle().Foreground(color).Render(indicator + " " + m.status)
	return lipgloss.NewStyle().Width(width).
		Render(lipgloss.JoinVertical(lipgloss.Top, header, styles.SectionContent.Render(statusLine)))
}

func (m *MCPStatus) Height() int { return 2 }

func (m *MCPStatus) Update(msg tea.Msg) (kit.Section, tea.Cmd) {
	switch v := msg.(type) {
	case kit.MCPStatusMsg:
		m.connected = v.Connected
		m.status = v.Status
	}
	return m, nil
}
