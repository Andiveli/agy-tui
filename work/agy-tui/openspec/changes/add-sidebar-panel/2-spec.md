# Spec: Right-side Sidebar Panel

**Change**: `add-sidebar-panel`
**Status**: ✅ Specified
**Date**: 2026-05-25

---

## 1. app-root

### Requirements
- `NewApp()` returns a `tea.Model` (root `App` struct)
- Composes `Sidebar` (right side, ~30 chars) + `MainPanel` (remaining width)
- Handles `tea.WindowSizeMsg` to compute and propagate widths
- Routes domain messages to correct child models
- Keybindings: `Ctrl+B` toggle sidebar visibility, `Ctrl+T` open theme editor
- Derives `Styles` from `Theme` on init and on `ThemeChangedMsg`

### Scenarios
| Scenario | Input | Expected Behavior |
|----------|-------|-------------------|
| App starts | `tea.Program.Run()` | Sidebar right + main panel left, styled with Catppuccin Mocha |
| Window resizes | `tea.WindowSizeMsg{Width: 120, Height: 40}` | Sidebar=30, MainPanel=90 |
| Sidebar toggle | `tea.KeyMsg{Ctrl+B}` | Sidebar toggles visible/hidden, main panel adjusts |
| Open theme editor | `tea.KeyMsg{Ctrl+T}` | Theme editor overlay opens |
| Theme changed | `ThemeChangedMsg` | All components rerender with new `Styles` |

### Acceptance Criteria
- Layout renders correctly at minimum 80×24 terminal size
- Sidebar is exactly 30 chars wide when visible
- Main panel fills remaining width
- Ctrl+B toggles sidebar without errors
- Ctrl+T opens theme editor overlay
- Window resize recalculates widths without visual glitches

### Keybindings
| Key | Action |
|-----|--------|
| Ctrl+B | Toggle sidebar visibility |
| Ctrl+T | Open theme editor overlay |
| Ctrl+C | Quit |

---

## 2. sidebar-panel

### Requirements
- `Section` interface: `Name() string`, `View(width int, styles Styles) string`, `Height() int`
- Optional `UpdatableSection` interface: `Update(msg tea.Msg) (Section, tea.Cmd)`
- 6 concrete sections: `SessionInfo`, `SubAgentProgress`, `MCPStatus`, `LSPStatus`, `ContextPanel`, `FileChanges`
- Each section header styled with `StyleContextHeader`
- Sidebar model handles `visible` toggle
- `SubAgentProgress` uses `bubbles/spinner` for animation
- Sections exceeding available height: `SessionInfo` + `SubAgentProgress` are priority, rest scroll or collapse

### Scenarios
| Scenario | Input | Expected Behavior |
|----------|-------|-------------------|
| Sidebar visible | Init state | All 6 sections render with styled headers |
| Spinner ticks | `spinner.TickMsg` | SubAgentProgress updates frame |
| MCP status updated | `MCPStatusMsg{Connected: true}` | MCP section shows "Connected" in green |
| LSP status updated | `LSPStatusMsg{Connected: false}` | LSP section shows "Disconnected" in red |
| File changed | `FileChangedMsg{Path: "foo.go", Action: "modified"}` | FileChanges section shows entry |
| Session changed | `SessionChangedMsg{Name: "my-session"}` | SessionInfo updates |
| Sidebar hidden | Ctrl+B | Only main panel visible |
| Narrow terminal | Width < 100 chars | Sidebar auto-hides |

### Acceptance Criteria
- All 6 sections render with correct headers and content
- Spinner animates at correct frame rate
- Messages update correct sections within 1 frame
- Toggle works correctly (visible → hidden → visible)
- Auto-hide on narrow terminals

---

## 3. theme-system

### Requirements
- `Theme` struct with all 22 color fields (Catppuccin Mocha palette)
- `DeriveStyles(t Theme) Styles` function — derives all 15 lipgloss styles
- `DefaultCatppuccinMocha` constant matching current hardcoded colors exactly
- Backward compatibility: existing color var names become `DefaultCatppuccinMocha.FieldName`
- Theme editor: overlay with `textinput.Model` per color field + live preview + close on Escape
- On save: broadcast `ThemeChangedMsg{}`

### Scenarios
| Scenario | Input | Expected Behavior |
|----------|-------|-------------------|
| Theme loads | Init | Styles derived from DefaultCatppuccinMocha |
| Color edited | User changes ColorBase | Live preview updates immediately |
| Theme saved | User confirms | ThemeChangedMsg broadcast, all components rerender |
| Theme canceled | Escape | Editor closes, no changes applied |
| New theme applied | ThemeChangedMsg | All styles recalculated, View() uses new styles |

### Acceptance Criteria
- All 22 colors map to `lipgloss.Color()` correctly
- `DeriveStyles(DefaultCatppuccinMocha)` produces styles identical to current hardcoded ones
- Theme editor shows all 22 inputs, navigable by Tab
- Live preview shows sample text with current color edits
- Escape closes without saving
- ThemeChangedMsg reaches all components
