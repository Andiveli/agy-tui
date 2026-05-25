package ui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/samael/agy-tui/internal/backend"
	"github.com/samael/agy-tui/internal/ui/kit"
)

// pollInterval is how often MCP/LSP status is refreshed.
const pollInterval = 15 * time.Second

// App is the root Bubbletea model.
type App struct {
	theme       kit.Theme
	styles      kit.Styles
	sidebar     *Sidebar
	chat        *ChatModel
	statusBar   *StatusBar
	themeEditor *kit.ThemeEditor
	showSidebar bool
	width       int
	height      int
}

func NewApp() *App {
	t := kit.DefaultCatppuccinMocha()
	s := kit.DeriveStyles(t)
	return &App{
		theme:       t,
		styles:      s,
		sidebar:     NewSidebar(s),
		chat:        NewChatModel(s, backend.NewClient()),
		statusBar:   NewStatusBar(s),
		showSidebar: true,
		width:       80,
		height:      24,
	}
}

func (a *App) Init() tea.Cmd {
	return tea.Batch(
		a.sidebar.Init(),
		a.chat.Init(),
		a.probeMCP(),
		a.probeLSP(),
		a.probeWorkspace(),
		a.pollTick(),
	)
}

// pollTick returns a command that triggers periodic MCP/LSP re-probing.
func (a *App) pollTick() tea.Cmd {
	return tea.Tick(pollInterval, func(t time.Time) tea.Msg {
		return pollTickMsg{}
	})
}

// pollTickMsg is sent when it's time to re-probe MCP and LSP.
type pollTickMsg struct{}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If theme editor is open, route to it first
	if a.themeEditor != nil && a.themeEditor.Open {
		return a.handleThemeEditorUpdate(msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		return a, a.handleResize(msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+b":
			a.showSidebar = !a.showSidebar
			a.sidebar.ToggleVisibility()
			return a, nil
		case "ctrl+t":
			a.openThemeEditor()
			return a, nil
		}

	case kit.ThemeChangedMsg:
		a.styles = kit.DeriveStyles(a.theme)
		return a, nil

	// ── Tick ──

	case pollTickMsg:
		return a, tea.Batch(a.probeMCP(), a.probeLSP(), a.pollTick())

	// ── Domain messages routed to sidebar + status bar ──

	case kit.MCPStatusMsg:
		a.sidebar, _ = a.sidebar.Update(msg)
		a.statusBar.Update(msg)
		return a, nil
	case kit.LSPStatusMsg:
		a.sidebar, _ = a.sidebar.Update(msg)
		a.statusBar.Update(msg)
		return a, nil
	case kit.SessionChangedMsg:
		a.sidebar, _ = a.sidebar.Update(msg)
		a.statusBar.Update(msg)
		return a, nil
	case kit.ProgressMsg:
		a.sidebar, _ = a.sidebar.Update(msg)
		a.statusBar.Update(msg)
		return a, nil
	case kit.FileChangedMsg:
		a.sidebar, _ = a.sidebar.Update(msg)
		return a, nil
	}

	// Forward to chat
	chatModel, cmd := a.chat.Update(msg)
	a.chat = chatModel.(*ChatModel)

	// Forward to sidebar
	sidebar, sideCmd := a.sidebar.Update(msg)
	a.sidebar = sidebar

	return a, tea.Batch(cmd, sideCmd)
}

func (a *App) View() string {
	// Full-screen modal when theme editor is open
	if a.themeEditor != nil && a.themeEditor.Open {
		return RenderThemeEditorOverlay(a.themeEditor, a.width, a.height)
	}

	mainContent := a.chat.View()
	if a.showSidebar {
		sidebarContent := a.sidebar.View()
		if sidebarContent != "" {
			mainContent = lipgloss.JoinHorizontal(lipgloss.Top, mainContent, sidebarContent)
		}
	}

	bar := a.statusBar.View(a.width)
	return lipgloss.JoinVertical(lipgloss.Top, mainContent, bar)
}

func (a *App) handleResize(msg tea.WindowSizeMsg) tea.Cmd {
	// Reserve 1 line for the status bar
	contentHeight := msg.Height - 1
	if contentHeight < 5 {
		contentHeight = 5
	}

	sidebarWidth := 0
	if a.showSidebar && msg.Width >= 100 {
		sidebarWidth = 30
	}
	a.sidebar.width = sidebarWidth
	a.sidebar.height = contentHeight

	mainWidth := msg.Width - sidebarWidth
	if mainWidth < 10 {
		mainWidth = 10
	}
	_, cmd := a.chat.Update(tea.WindowSizeMsg{
		Width: mainWidth, Height: contentHeight,
	})
	return cmd
}

// ── Theme Editor ────────────────────────────────────────────────────

// openThemeEditor creates and opens the theme editor.
func (a *App) openThemeEditor() {
	a.themeEditor = kit.NewThemeEditor(a.theme, a.styles)
	a.themeEditor.Open = true
	a.themeEditor.Focused = 0
	if len(a.themeEditor.Inputs) > 0 {
		a.themeEditor.Inputs[0].Focus()
	}
}

// handleThemeEditorUpdate routes messages to the theme editor when it's open.
func (a *App) handleThemeEditorUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	if a.themeEditor == nil || !a.themeEditor.Open {
		return a, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			a.themeEditor.Open = false
			a.themeEditor = nil
			return a, nil
		case "enter":
			// Save theme
			a.theme = a.themeEditor.Preview
			a.themeEditor.Open = false
			a.themeEditor = nil
			return a, tea.Batch(
				kit.ThemeChangedMsgCmd(),
			)
		case "tab", "down":
			if a.themeEditor != nil {
				a.themeEditor.Focused = (a.themeEditor.Focused + 1) % len(a.themeEditor.Inputs)
				a.themeEditor.FocusInputs()
			}
			return a, nil
		case "shift+tab", "up":
			if a.themeEditor != nil {
				a.themeEditor.Focused = (a.themeEditor.Focused - 1 + len(a.themeEditor.Inputs)) % len(a.themeEditor.Inputs)
				a.themeEditor.FocusInputs()
			}
			return a, nil
		}
	}

	// Forward to editor for text input changes
	if a.themeEditor != nil && a.themeEditor.Focused < len(a.themeEditor.Inputs) {
		var cmd tea.Cmd
		a.themeEditor.Inputs[a.themeEditor.Focused], cmd = a.themeEditor.Inputs[a.themeEditor.Focused].Update(msg)
		a.themeEditor.UpdatePreviewFromInputs()
		return a, cmd
	}
	return a, nil
}

// ── Startup probing ────────────────────────────────────────────────

func (a *App) probeMCP() tea.Cmd {
	return func() tea.Msg {
		home, err := os.UserHomeDir()
		if err != nil {
			return kit.MCPStatusMsg{Connected: false, Status: "Cannot determine home dir"}
		}
		mcpPath := filepath.Join(home, ".gemini", "antigravity-cli", "mcp_config.json")
		data, err := os.ReadFile(mcpPath)
		if err != nil {
			return kit.MCPStatusMsg{Connected: false, Status: "No MCP config"}
		}
		if strings.Contains(string(data), "servers") || strings.Contains(string(data), "serverUrl") {
			return kit.MCPStatusMsg{Connected: true, Status: "MCP configured"}
		}
		return kit.MCPStatusMsg{Connected: false, Status: "No MCP servers"}
	}
}

func (a *App) probeLSP() tea.Cmd {
	return func() tea.Msg {
		servers := []string{"gopls", "typescript-language-server", "rust-analyzer", "pyright", "clangd"}
		var found []string
		for _, s := range servers {
			if err := exec.Command("pgrep", "-x", s).Run(); err == nil {
				found = append(found, s)
			}
		}
		if len(found) > 0 {
			return kit.LSPStatusMsg{Connected: true, Status: strings.Join(found, ", ")}
		}
		return kit.LSPStatusMsg{Connected: false, Status: "No LSP detected"}
	}
}

func (a *App) probeWorkspace() tea.Cmd {
	return func() tea.Msg {
		cwd, err := os.Getwd()
		base := "unknown"
		if err == nil {
			base = filepath.Base(cwd)
		}
		return kit.SessionChangedMsg{
			Name:    "default",
			Context: base,
		}
	}
}
