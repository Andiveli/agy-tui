package ui

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/samael/agy-tui/internal/ui/kit"
)

// ThemeEditor handles are re-exported from kit for convenience.
// The full implementation lives in internal/ui/kit/kit.go.
//
// Example usage:
//
//	editor := kit.NewThemeEditor(theme, styles)
//	editor.Open()
//	// ... in Update: editor, cmd = editor.Update(msg)
//	if editor.Saved {
//	    newTheme := editor.Theme()
//	    // apply theme ...
//	}

// RenderThemeEditorOverlay renders a centered theme editor overlay.
func RenderThemeEditorOverlay(editor *kit.ThemeEditor, width, height int) string {
	if !editor.Open {
		return ""
	}
	editorView := renderEditorView(editor)
	return lipgloss.Place(width, height,
		lipgloss.Center, lipgloss.Center,
		editorView,
	)
}

func renderEditorView(editor *kit.ThemeEditor) string {
	var lines []string
	lines = append(lines, editor.Styles.Title.Render(" Theme Editor "))
	lines = append(lines, "")

	fields := kit.ThemeFields(editor.Preview)
	for i, f := range fields {
		label := editor.Styles.Dimmed.Render(" " + f.Name + ":")
		inputView := editor.Inputs[i].View()
		lines = append(lines, label+" "+inputView)
	}

	lines = append(lines, "")
	lines = append(lines, renderPreview(editor))
	lines = append(lines, "")
	lines = append(lines, editor.Styles.KeyHint.Render(
		" Tab:next  Shift+Tab:prev  Enter:save  Esc:cancel "))

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(editor.Preview.Lavender).
		Padding(1).
		Render(lipgloss.JoinVertical(lipgloss.Top, lines...))
}

func renderPreview(editor *kit.ThemeEditor) string {
	var swatches []string
	colors := []lipgloss.Color{
		editor.Preview.Base, editor.Preview.Surface0, editor.Preview.Surface1,
		editor.Preview.Text, editor.Preview.Lavender, editor.Preview.Blue,
		editor.Preview.Green, editor.Preview.Yellow, editor.Preview.Peach,
		editor.Preview.Red, editor.Preview.Mauve,
	}
	for _, c := range colors {
		swatch := lipgloss.NewStyle().Background(c).Render("  ")
		swatches = append(swatches, swatch)
	}
	return editor.Styles.Dimmed.Render("Preview: ") +
		lipgloss.JoinHorizontal(lipgloss.Top, swatches...)
}
