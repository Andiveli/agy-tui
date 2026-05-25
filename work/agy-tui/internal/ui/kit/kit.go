// Package kit provides shared types for the agy-tui UI layer.
// It exists to break the import cycle between ui and ui/sections.
package kit

import (
	"io"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ─── Theme ───────────────────────────────────────────────────────────

// Theme holds the complete color palette.
type Theme struct {
	Base     lipgloss.Color
	Surface0 lipgloss.Color
	Surface1 lipgloss.Color
	Surface2 lipgloss.Color
	Overlay  lipgloss.Color
	Text     lipgloss.Color
	Subtext  lipgloss.Color
	Lavender lipgloss.Color
	Blue     lipgloss.Color
	Sapphire lipgloss.Color
	Green    lipgloss.Color
	Yellow   lipgloss.Color
	Peach    lipgloss.Color
	Red      lipgloss.Color
	Mauve    lipgloss.Color
	Pink     lipgloss.Color
}

// DefaultCatppuccinMocha returns the default theme.
func DefaultCatppuccinMocha() Theme {
	return Theme{
		Base:     lipgloss.Color("#1e1e2e"),
		Surface0: lipgloss.Color("#313244"),
		Surface1: lipgloss.Color("#45475a"),
		Surface2: lipgloss.Color("#585b70"),
		Overlay:  lipgloss.Color("#6c7086"),
		Text:     lipgloss.Color("#cdd6f4"),
		Subtext:  lipgloss.Color("#a6adc8"),
		Lavender: lipgloss.Color("#b4befe"),
		Blue:     lipgloss.Color("#89b4fa"),
		Sapphire: lipgloss.Color("#74c7ec"),
		Green:    lipgloss.Color("#a6e3a1"),
		Yellow:   lipgloss.Color("#f9e2af"),
		Peach:    lipgloss.Color("#fab387"),
		Red:      lipgloss.Color("#f38ba8"),
		Mauve:    lipgloss.Color("#cba6f7"),
		Pink:     lipgloss.Color("#f5c2e7"),
	}
}

// ─── Styles ──────────────────────────────────────────────────────────

// Styles holds all derived lipgloss styles.
type Styles struct {
	Border         lipgloss.Style
	ActiveBorder   lipgloss.Style
	Title          lipgloss.Style
	StatusBar      lipgloss.Style
	StatusAccent   lipgloss.Style
	KeyHint        lipgloss.Style
	KeyBinding     lipgloss.Style
	UserMessage    lipgloss.Style
	AgentMessage   lipgloss.Style
	SystemMessage  lipgloss.Style
	Error          lipgloss.Style
	Spinner        lipgloss.Style
	Dimmed         lipgloss.Style
	ContextHeader  lipgloss.Style
	ContextItem    lipgloss.Style
	ToolCall       lipgloss.Style
	InputPrompt    lipgloss.Style
	SectionHeader  lipgloss.Style
	SectionContent lipgloss.Style
}

// DeriveStyles creates a complete Styles set from a Theme.
func DeriveStyles(t Theme) Styles {
	return Styles{
		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(t.Surface2),

		ActiveBorder: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(t.Lavender),

		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Lavender).
			PaddingLeft(1).
			PaddingRight(1),

		StatusBar: lipgloss.NewStyle().
			Background(t.Surface0).
			Foreground(t.Text).
			PaddingLeft(1).
			PaddingRight(1),

		StatusAccent: lipgloss.NewStyle().
			Background(t.Lavender).
			Foreground(t.Base).
			Bold(true).
			PaddingLeft(1).
			PaddingRight(1),

		KeyHint: lipgloss.NewStyle().
			Foreground(t.Overlay),

		KeyBinding: lipgloss.NewStyle().
			Foreground(t.Lavender).
			Bold(true),

		UserMessage: lipgloss.NewStyle().
			Foreground(t.Blue).
			Bold(true),

		AgentMessage: lipgloss.NewStyle().
			Foreground(t.Text),

		SystemMessage: lipgloss.NewStyle().
			Foreground(t.Yellow).
			Italic(true),

		Error: lipgloss.NewStyle().
			Foreground(t.Red).
			Bold(true),

		Spinner: lipgloss.NewStyle().
			Foreground(t.Mauve),

		Dimmed: lipgloss.NewStyle().
			Foreground(t.Overlay),

		ContextHeader: lipgloss.NewStyle().
			Foreground(t.Sapphire).
			Bold(true),

		ContextItem: lipgloss.NewStyle().
			Foreground(t.Subtext).
			PaddingLeft(1),

		ToolCall: lipgloss.NewStyle().
			Foreground(t.Peach),

		InputPrompt: lipgloss.NewStyle().
			Foreground(t.Green).
			Bold(true),

		SectionHeader: lipgloss.NewStyle().
			Foreground(t.Sapphire).
			Bold(true).
			PaddingLeft(1).
			PaddingRight(1),

		SectionContent: lipgloss.NewStyle().
			Foreground(t.Subtext).
			PaddingLeft(1).
			PaddingRight(1),
	}
}

// ─── Section ─────────────────────────────────────────────────────────

// Section is a renderable unit within the sidebar.
type Section interface {
	Name() string
	View(width int, styles Styles) string
	Height() int
}

// UpdatableSection extends Section with message handling.
type UpdatableSection interface {
	Section
	Update(msg tea.Msg) (Section, tea.Cmd)
}

// ─── Messages ────────────────────────────────────────────────────────

type MCPStatusMsg struct {
	Connected bool
	Status    string
}

type LSPStatusMsg struct {
	Connected bool
	Status    string
}

type ProgressMsg struct {
	SubAgent string
	Status   string
	Progress int
}

type FileChangedMsg struct {
	Path   string
	Action string
}

type SessionChangedMsg struct {
	Name      string
	Context   string
	ConvCount int    // number of saved agy conversations
	ConvID    string // current conversation ID (if resumed)
}

type ThemeChangedMsg struct{}

// ThemeChangedMsgCmd returns a command that broadcasts a theme change.
func ThemeChangedMsgCmd() tea.Cmd {
	return func() tea.Msg {
		return ThemeChangedMsg{}
	}
}

// ─── Chat ────────────────────────────────────────────────────────────

// ChatMessage represents a single message in the chat history.
type ChatMessage struct {
	Role    string // "user" or "agent"
	Content string
}

// StreamReadyMsg carries the agy pipe reader from the startup goroutine.
// Update stores the reader on the model to avoid writing model state from inside cmd goroutines.
type StreamReadyMsg struct {
	Reader io.ReadCloser
}

// ChatStreamChunkMsg is sent for each line of a streaming response.
type ChatStreamChunkMsg struct {
	Text string
}

// ChatCompletedMsg carries a completed agy response with sidebar-relevant metadata.
type ChatCompletedMsg struct {
	Content     string
	SessionName string
	FilePaths   []string
}

// ChatErrorMsg carries an agy error back from the goroutine.
type ChatErrorMsg struct {
	Err error
}

// ─── Theme Editor ────────────────────────────────────────────────────

// ThemeEditor is an overlay model for editing theme colors at runtime.
type ThemeEditor struct {
	Open    bool
	Inputs  []textinput.Model
	Preview Theme
	Styles  Styles
	Focused int
	Saved   bool
}

// NewThemeEditor creates a ThemeEditor overlay with inputs for each color.
func NewThemeEditor(theme Theme, styles Styles) *ThemeEditor {
	fields := ThemeFields(theme)
	inputs := make([]textinput.Model, len(fields))

	for i, f := range fields {
		ti := textinput.New()
		ti.Placeholder = f.Value
		ti.SetValue(f.Value)
		ti.CharLimit = 7
		ti.Width = 10
		if i == 0 {
			ti.Focus()
		}
		inputs[i] = ti
	}

	return &ThemeEditor{
		Inputs:  inputs,
		Preview: theme,
		Styles:  styles,
		Focused: 0,
	}
}

// Init implements a minimal tea.Model interface for compatibility.
func (e *ThemeEditor) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (e *ThemeEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !e.Open {
		return e, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			e.Open = false
			e.Saved = false
			return e, nil
		case "enter":
			e.Saved = true
			e.Open = false
			return e, ThemeChangedMsgCmd()
		case "tab", "down":
			e.Focused = (e.Focused + 1) % len(e.Inputs)
			e.FocusInputs()
		case "shift+tab", "up":
			e.Focused = (e.Focused - 1 + len(e.Inputs)) % len(e.Inputs)
			e.FocusInputs()
		}
	}
	var cmd tea.Cmd
	e.Inputs[e.Focused], cmd = e.Inputs[e.Focused].Update(msg)
	e.UpdatePreviewFromInputs()
	return e, cmd
}

// View implements tea.Model (minimal — actual rendering is done via RenderThemeEditorOverlay).
func (e *ThemeEditor) View() string { return "" }

func (e *ThemeEditor) FocusInputs() {
	for i := range e.Inputs {
		if i == e.Focused {
			e.Inputs[i].Focus()
		} else {
			e.Inputs[i].Blur()
		}
	}
}

func (e *ThemeEditor) UpdatePreviewFromInputs() {
	fields := ThemeFields(e.Preview)
	for i, f := range fields {
		val := e.Inputs[i].Value()
		if val == "" {
			val = f.Value
		}
		if val[0] != '#' {
			val = "#" + val
		}
		c := lipgloss.Color(val)
		switch f.Name {
		case "Base":
			e.Preview.Base = c
		case "Surface0":
			e.Preview.Surface0 = c
		case "Surface1":
			e.Preview.Surface1 = c
		case "Surface2":
			e.Preview.Surface2 = c
		case "Overlay":
			e.Preview.Overlay = c
		case "Text":
			e.Preview.Text = c
		case "Subtext":
			e.Preview.Subtext = c
		case "Lavender":
			e.Preview.Lavender = c
		case "Blue":
			e.Preview.Blue = c
		case "Sapphire":
			e.Preview.Sapphire = c
		case "Green":
			e.Preview.Green = c
		case "Yellow":
			e.Preview.Yellow = c
		case "Peach":
			e.Preview.Peach = c
		case "Red":
			e.Preview.Red = c
		case "Mauve":
			e.Preview.Mauve = c
		case "Pink":
			e.Preview.Pink = c
		}
	}
}

// FieldDef describes a single color field in the Theme struct.
type FieldDef struct {
	Name  string
	Value string
}

// ThemeFields returns the list of editable color fields.
func ThemeFields(t Theme) []FieldDef {
	return []FieldDef{
		{"Base", string(t.Base)},
		{"Surface0", string(t.Surface0)},
		{"Surface1", string(t.Surface1)},
		{"Surface2", string(t.Surface2)},
		{"Overlay", string(t.Overlay)},
		{"Text", string(t.Text)},
		{"Subtext", string(t.Subtext)},
		{"Lavender", string(t.Lavender)},
		{"Blue", string(t.Blue)},
		{"Sapphire", string(t.Sapphire)},
		{"Green", string(t.Green)},
		{"Yellow", string(t.Yellow)},
		{"Peach", string(t.Peach)},
		{"Red", string(t.Red)},
		{"Mauve", string(t.Mauve)},
		{"Pink", string(t.Pink)},
	}
}
