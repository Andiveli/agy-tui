# Design: Right-side Sidebar Panel

**Change**: `add-sidebar-panel`
**Status**: ✅ Designed
**Date**: 2026-05-25

---

## 1. Theme System

### Theme Struct
```go
// internal/ui/theme.go

type Theme struct {
    Base     lipgloss.Color `yaml:"base"`
    Surface0 lipgloss.Color `yaml:"surface0"`
    Surface1 lipgloss.Color `yaml:"surface1"`
    Surface2 lipgloss.Color `yaml:"surface2"`
    Overlay  lipgloss.Color `yaml:"overlay"`
    Text     lipgloss.Color `yaml:"text"`
    Subtext  lipgloss.Color `yaml:"subtext"`
    Lavender lipgloss.Color `yaml:"lavender"`
    Blue     lipgloss.Color `yaml:"blue"`
    Sapphire lipgloss.Color `yaml:"sapphire"`
    Green    lipgloss.Color `yaml:"green"`
    Yellow   lipgloss.Color `yaml:"yellow"`
    Peach    lipgloss.Color `yaml:"peach"`
    Red      lipgloss.Color `yaml:"red"`
    Mauve    lipgloss.Color `yaml:"mauve"`
    Pink     lipgloss.Color `yaml:"pink"`
}
```

### Styles Struct
```go
type Styles struct {
    Border          lipgloss.Style
    ActiveBorder    lipgloss.Style
    Title           lipgloss.Style
    StatusBar       lipgloss.Style
    StatusAccent    lipgloss.Style
    KeyHint         lipgloss.Style
    KeyBinding      lipgloss.Style
    UserMessage     lipgloss.Style
    AgentMessage    lipgloss.Style
    SystemMessage   lipgloss.Style
    Error           lipgloss.Style
    Spinner         lipgloss.Style
    Dimmed          lipgloss.Style
    ContextHeader   lipgloss.Style
    ContextItem     lipgloss.Style
    ToolCall        lipgloss.Style
    InputPrompt     lipgloss.Style
    SectionHeader   lipgloss.Style
    SectionContent  lipgloss.Style
}
```

### Derive Function
```go
var DefaultCatppuccinMocha = Theme{
    Base:     lipgloss.Color("#1e1e2e"),
    Surface0: lipgloss.Color("#313244"),
    // ... all 22 fields
}

func DeriveStyles(t Theme) Styles {
    return Styles{
        Border: lipgloss.NewStyle().
            BorderStyle(lipgloss.RoundedBorder()).
            BorderForeground(t.Surface2),
        Title: lipgloss.NewStyle().
            Bold(true).
            Foreground(t.Lavender),
        // ... all 15 styles derived from t
    }
}
```

---

## 2. App Model

```go
// internal/ui/app.go

type App struct {
    theme       Theme
    styles      Styles
    sidebar     Sidebar
    mainPanel   MainPanel
    showSidebar bool
    width       int
    height      int
}

func NewApp() *App {
    t := DefaultCatppuccinMocha
    return &App{
        theme:       t,
        styles:      DeriveStyles(t),
        sidebar:     NewSidebar(DeriveStyles(t)),
        mainPanel:   NewMainPanel(DeriveStyles(t)),
        showSidebar: true,
    }
}

func (a *App) Init() tea.Cmd {
    return tea.Batch(a.sidebar.Init(), a.mainPanel.Init())
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        a.width = msg.Width
        a.height = msg.Height
        sidebarWidth := 0
        if a.showSidebar && msg.Width >= 100 {
            sidebarWidth = 30
        }
        a.sidebar.width = sidebarWidth
        a.sidebar.height = msg.Height
        a.mainPanel.width = msg.Width - sidebarWidth - 2
        a.mainPanel.height = msg.Height

    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+b":
            a.showSidebar = !a.showSidebar
            // trigger resize recalc
            return a, nil
        case "ctrl+t":
            // open theme editor
            return a, nil
        case "ctrl+c":
            return a, tea.Quit
        }

    case ThemeChangedMsg:
        a.styles = DeriveStyles(a.theme)
        a.sidebar.styles = a.styles
        a.mainPanel.styles = a.styles
        return a, nil

    default:
        var cmd tea.Cmd
        a.sidebar, cmd = a.sidebar.Update(msg)
        return a, cmd
    }
    return a, nil
}

func (a *App) View() string {
    if !a.showSidebar {
        return a.mainPanel.View()
    }
    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        a.mainPanel.View(),
        a.sidebar.View(),
    )
}
```

---

## 3. Sidebar Model

```go
// internal/ui/sidebar.go

type Sidebar struct {
    width    int
    height   int
    visible  bool
    sections []Section
    styles   Styles
}

type Section interface {
    Name() string
    View(width int, styles Styles) string
    Height() int
}

type UpdatableSection interface {
    Section
    Update(msg tea.Msg) (Section, tea.Cmd)
}

func NewSidebar(styles Styles) Sidebar {
    return Sidebar{
        sections: []Section{
            NewSessionSection(),
            NewProgressSection(),
            NewMCPSection(),
            NewLSPSection(),
            NewContextSection(),
            NewFileChangesSection(),
        },
        styles: styles,
    }
}

func (s *Sidebar) Init() tea.Cmd {
    var cmds []tea.Cmd
    for _, sec := range s.sections {
        if upd, ok := sec.(tea.Model); ok {
            cmds = append(cmds, upd.Init())
        }
    }
    return tea.Batch(cmds...)
}

func (s *Sidebar) Update(msg tea.Msg) (Sidebar, tea.Cmd) {
    switch msg := msg.(type) {
    case MCPStatusMsg, LSPStatusMsg, ProgressMsg,
         FileChangedMsg, SessionChangedMsg:
        for i, sec := range s.sections {
            if upd, ok := sec.(UpdatableSection); ok {
                updated, cmd := upd.Update(msg)
                s.sections[i] = updated.(Section)
                // collect cmd
                _ = cmd
            }
        }
    }
    return *s, nil
}

func (s *Sidebar) View() string {
    if s.width == 0 {
        return ""
    }
    var rendered []string
    remaining := s.height
    for _, sec := range s.sections {
        h := sec.Height()
        if h > remaining {
            // mark as collapsed
            rendered = append(rendered, s.renderCollapsed(sec.Name()))
            continue
        }
        rendered = append(rendered, sec.View(s.width-2, s.styles))
        remaining -= h
    }
    return lipgloss.NewStyle().
        Width(s.width).
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(s.styles.Border.GetBorderStyle()).
        Render(lipgloss.JoinVertical(lipgloss.Top, rendered...))
}

func (s *Sidebar) renderCollapsed(name string) string {
    return s.styles.Dimmed.Render("⋯ " + name)
}
```

---

## 4. Section Interface + Implementations

### Section Interface
```go
// internal/ui/section.go

type Section interface {
    Name() string
    View(width int, styles Styles) string
    Height() int
}

type UpdatableSection interface {
    Section
    Update(msg tea.Msg) (Section, tea.Cmd)
}
```

### Concrete Sections

Each section follows this pattern:
```
┌─ Section Name ──────────────┐
│ content line 1              │
│ content line 2              │
│ ...                         │
└─────────────────────────────┘
```

**SessionInfo** (`sections/session.go`):
- Height: 3 lines (header + name + context)
- No Update (static until SessionChangedMsg)
- Shows session name and brief context

**SubAgentProgress** (`sections/progress.go`):
- Height: 4 lines (header + spinner + agent name + status)
- Has `bubbles/spinner.Model` internally
- Updates on `ProgressMsg` and `spinner.TickMsg`

**MCPStatus** (`sections/mcp.go`):
- Height: 2 lines (header + status line)
- Updates on `MCPStatusMsg`

**LSPStatus** (`sections/lsp.go`):
- Height: 2 lines (header + status line)
- Updates on `LSPStatusMsg`

**ContextPanel** (`sections/context.go`):
- Height: variable (header + N context items, capped at 5)
- Updates on context changes

**FileChanges** (`sections/files.go`):
- Height: variable (header + N recent files, capped at 5)
- Updates on `FileChangedMsg`

---

## 5. MainPanel (placeholder)

```go
// internal/ui/main_panel.go

type MainPanel struct {
    width  int
    height int
    styles Styles
}

func NewMainPanel(styles Styles) MainPanel {
    return MainPanel{styles: styles}
}

func (m MainPanel) Init() tea.Cmd { return nil }

func (m MainPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return m, nil
}

func (m MainPanel) View() string {
    return lipgloss.NewStyle().
        Width(m.width).
        Height(m.height).
        Align(lipgloss.Center).
        AlignVertical(lipgloss.Center).
        Render("Chat area — coming soon")
}
```

---

## 6. Theme Editor Overlay

```go
// internal/ui/theme_editor.go

type ThemeEditor struct {
    open    bool
    inputs  []textinput.Model
    preview Theme
    styles  Styles
    focused int
    save    bool
}

func NewThemeEditor(theme Theme, styles Styles) ThemeEditor {
    inputs := make([]textinput.Model, 0)
    // Create textinput for each color field
    // ...
    return ThemeEditor{
        inputs:  inputs,
        preview: theme,
        styles:  styles,
    }
}
```

**Keybindings** (when open):
- `Tab` / `Shift+Tab`: cycle focus between inputs
- `Enter`: save theme and close
- `Escape`: close without saving
- `Up` / `Down`: navigate inputs

---

## 7. Message Types

```go
// internal/ui/messages.go

type MCPStatusMsg struct {
    Connected bool
    Status    string // e.g., "Connected to 3 servers"
}

type LSPStatusMsg struct {
    Connected bool
    Status    string
}

type ProgressMsg struct {
    SubAgent string // agent name
    Status   string // "running", "completed", "failed"
    Progress int    // 0-100 percentage
}

type FileChangedMsg struct {
    Path   string
    Action string // "created", "modified", "deleted"
}

type SessionChangedMsg struct {
    Name    string
    Context string // file/directory context
}

type ThemeChangedMsg struct{}
```

---

## 8. File Structure

```
internal/ui/
├── app.go               — Root App model (NewApp, Update routing, View layout)
├── app_test.go          — App tests: resize, toggle, message routing
├── sidebar.go           — Sidebar model (sections slice, Update dispatch, View)
├── sidebar_test.go      — Sidebar tests: height calc, toggle, auto-hide
├── section.go           — Section + UpdatableSection interfaces
├── main_panel.go        — MainPanel placeholder
├── main_panel_test.go   — MainPanel tests
├── theme.go             — Theme struct, Styles struct, DeriveStyles(), DefaultCatppuccinMocha
├── theme_test.go        — Theme tests: derive matches old styles, serialization
├── theme_editor.go      — ThemeEditor overlay model
├── theme_editor_test.go — ThemeEditor tests
├── messages.go          — All domain message types
└── sections/
    ├── session.go       — SessionInfo section
    ├── progress.go      — SubAgentProgress section (with spinner)
    ├── mcp.go           — MCPStatus section
    ├── lsp.go           — LSPStatus section
    ├── context.go       — ContextPanel section
    └── files.go         — FileChanges section
```

---

## 9. Error Handling

| Condition | Behavior |
|-----------|----------|
| Width < 80 chars | Minimum layout: sidebar auto-hides, main panel takes full width |
| Width < 100 chars | Sidebar auto-hides (too narrow for sidebar + main) |
| Height insufficient for all sections | Priority order: SessionInfo > SubAgentProgress > rest |
| Theme editor overflow | Editor opens as full-screen overlay, not constrained to sidebar width |
| Zero sections | Sidebar renders empty bordered panel |
| nil styles | Use DefaultCatppuccinMocha as fallback |
