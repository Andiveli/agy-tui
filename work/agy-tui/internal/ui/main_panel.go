package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/samael/agy-tui/internal/ui/kit"
)

// MainPanel is a placeholder for the future chat/viewport area.
type MainPanel struct {
	width  int
	height int
	styles kit.Styles
}

func NewMainPanel(styles kit.Styles) MainPanel {
	return MainPanel{width: 80, height: 24, styles: styles}
}

func (m MainPanel) Init() tea.Cmd { return nil }

func (m MainPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m MainPanel) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}
	content := lipgloss.NewStyle().
		Width(m.width).Height(m.height).
		Align(lipgloss.Center).AlignVertical(lipgloss.Center).
		Render("Chat area — coming soon")
	return m.styles.Border.Width(m.width).Height(m.height).Render(content)
}
