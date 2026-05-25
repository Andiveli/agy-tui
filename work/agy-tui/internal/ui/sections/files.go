package sections

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/samael/agy-tui/internal/ui/kit"
)

// FileChanges displays recently changed files.
type FileChanges struct {
	files []fileEntry
}

type fileEntry struct {
	path   string
	action string
}

func NewFileChanges() *FileChanges {
	return &FileChanges{}
}

func (f *FileChanges) Name() string { return "Files" }

func (f *FileChanges) View(width int, styles kit.Styles) string {
	header := styles.SectionHeader.Render(" " + f.Name() + " ")
	var body string
	if len(f.files) == 0 {
		body = styles.Dimmed.Render(" no changes")
	} else {
		var lines []string
		maxFiles := 5
		if len(f.files) < maxFiles {
			maxFiles = len(f.files)
		}
		for _, fe := range f.files[:maxFiles] {
			color := lipgloss.Color("#f9e2af")
			icon := "~"
			switch fe.action {
			case "created":
				icon = "+"
				color = lipgloss.Color("#a6e3a1")
			case "deleted":
				icon = "-"
				color = lipgloss.Color("#f38ba8")
			}
			line := lipgloss.NewStyle().Foreground(color).Render(fmt.Sprintf(" %s %s", icon, fe.path))
			lines = append(lines, styles.ContextItem.Render(line))
		}
		body = lipgloss.JoinVertical(lipgloss.Top, lines...)
	}
	return lipgloss.NewStyle().Width(width).
		Render(lipgloss.JoinVertical(lipgloss.Top, header, body))
}

func (f *FileChanges) Height() int {
	h := 2
	if len(f.files) > 0 {
		if len(f.files) > 5 {
			h += 5
		} else {
			h += len(f.files)
		}
	} else {
		h++
	}
	return h
}

func (f *FileChanges) Update(msg tea.Msg) (kit.Section, tea.Cmd) {
	switch v := msg.(type) {
	case kit.FileChangedMsg:
		f.files = append([]fileEntry{{path: v.Path, action: v.Action}}, f.files...)
		if len(f.files) > 10 {
			f.files = f.files[:10]
		}
	}
	return f, nil
}
