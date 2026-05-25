package sections

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/samael/agy-tui/internal/ui/kit"
)

// ContextPanel displays workspace context items.
type ContextPanel struct {
	items []string
}

func NewContextPanel() *ContextPanel {
	return &ContextPanel{}
}

func (c *ContextPanel) Name() string { return "Context" }

func (c *ContextPanel) View(width int, styles kit.Styles) string {
	header := styles.SectionHeader.Render(" " + c.Name() + " ")
	var body string
	if len(c.items) == 0 {
		body = styles.Dimmed.Render(" no context")
	} else {
		var lines []string
		maxItems := 5
		if len(c.items) < maxItems {
			maxItems = len(c.items)
		}
		for _, item := range c.items[:maxItems] {
			lines = append(lines, styles.ContextItem.Render(" " + item))
		}
		body = lipgloss.JoinVertical(lipgloss.Top, lines...)
	}
	return lipgloss.NewStyle().Width(width).
		Render(lipgloss.JoinVertical(lipgloss.Top, header, body))
}

func (c *ContextPanel) Height() int {
	h := 2
	if len(c.items) > 0 {
		if len(c.items) > 5 {
			h += 5
		} else {
			h += len(c.items)
		}
	} else {
		h++
	}
	return h
}

func (c *ContextPanel) Update(msg tea.Msg) (kit.Section, tea.Cmd) {
	switch v := msg.(type) {
	case kit.SessionChangedMsg:
		if v.Context != "" {
			c.items = append(c.items, v.Context)
		}
	}
	return c, nil
}
