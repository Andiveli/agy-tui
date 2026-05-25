package sections

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/samael/agy-tui/internal/ui/kit"
)

// SubAgentProgress displays sub-agent execution progress with a spinner.
type SubAgentProgress struct {
	spinner spinner.Model
	agents  []agentEntry
}

type agentEntry struct {
	name     string
	status   string
	progress int
}

// NewSubAgentProgress creates a progress section with a spinner.
func NewSubAgentProgress() *SubAgentProgress {
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#cba6f7"))
	s.Spinner = spinner.Dot
	return &SubAgentProgress{spinner: s}
}

func (p *SubAgentProgress) Name() string { return "Sub-agents" }

func (p *SubAgentProgress) View(width int, styles kit.Styles) string {
	header := styles.SectionHeader.Render(" " + p.Name() + " ")

	var body string
	if len(p.agents) == 0 {
		body = styles.Dimmed.Render(" no active agents")
	} else {
		var lines []string
		for _, a := range p.agents {
			icon := p.spinner.View()
			statusColor := lipgloss.Color("#cba6f7")
			switch a.status {
			case "completed":
				icon = "✓"
				statusColor = lipgloss.Color("#a6e3a1")
			case "failed":
				icon = "✗"
				statusColor = lipgloss.Color("#f38ba8")
			}
			line := lipgloss.NewStyle().Foreground(statusColor).Render(fmt.Sprintf(" %s %s", icon, a.name))
			lines = append(lines, line)
		}
		body = lipgloss.JoinVertical(lipgloss.Top, lines...)
	}

	return lipgloss.NewStyle().Width(width).
		Render(lipgloss.JoinVertical(lipgloss.Top, header, body))
}

func (p *SubAgentProgress) Height() int {
	h := 2
	if len(p.agents) > 0 {
		h += len(p.agents)
	} else {
		h++
	}
	return h
}

func (p *SubAgentProgress) Init() tea.Cmd {
	return p.spinner.Tick
}

func (p *SubAgentProgress) Update(msg tea.Msg) (kit.Section, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		p.spinner, cmd = p.spinner.Update(msg)
		return p, cmd

	case kit.ProgressMsg:
		for i, a := range p.agents {
			if a.name == msg.SubAgent {
				p.agents[i].status = msg.Status
				p.agents[i].progress = msg.Progress
				return p, nil
			}
		}
		p.agents = append(p.agents, agentEntry{
			name: msg.SubAgent, status: msg.Status, progress: msg.Progress,
		})
	}
	return p, nil
}
