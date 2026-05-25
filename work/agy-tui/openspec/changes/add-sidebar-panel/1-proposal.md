# Proposal: Right-side Sidebar Panel

**Change**: `add-sidebar-panel`
**Status**: ✅ Proposed
**Date**: 2026-05-25

---

## Intent

The app has no UI model — `ui.NewApp()` is called but doesn't exist. We need a root model that composes a right-side info panel and a main panel. The sidebar provides at-a-glance session, agent, and system status without cluttering the chat area. Success: a working sidebar with styled sections, updating via messages.

## Scope

### In Scope

- Theme struct refactoring (prerequisite for runtime color editing)
- Theme editor (keybind to modify colors at runtime)
- Sidebar component with pluggable sections
- 6 sections: Session Info, Subagent Progress, MCP Status, LSP Status, Context Panel, File Changes
- Root `App` model composing Sidebar + MainPanel placeholder
- Message types for sidebar data flow
- Keybind toggles: show/hide sidebar, open theme editor

### Out of Scope

- Main chat panel implementation (placeholder only)
- Backend polling for MCP/LSP — pure UI slots
- agy conversation flow
- Plugin system
- Persistent theme config (future)

## Capabilities

No existing specs to modify — all capabilities are new.

### New Capabilities

- `app-root`: Root `tea.Model` composing sidebar + main panel, message dispatch
- `sidebar-panel`: Right-side panel with pluggable Section interface, show/hide toggle
- `theme-system`: Theme struct replacing package vars, runtime color editing via theme editor

## Approach

**Composed model** — `App` owns child models and delegates rendering/layout via lipgloss `JoinHorizontal`. Sidebar uses a `Section` interface (`View()`, `Update()`), each section is its own mini-model. Theme refactor: move package vars into a `Theme` struct, pass as dependency. The theme editor is a `bubbles.textinput`-based overlay.

```
App
├── Sidebar (hideable)
│   ├── SessionInfo
│   ├── SubagentProgress
│   ├── MCPStatus
│   ├── LSPStatus
│   ├── ContextPanel
│   └── FileChanges
├── MainPanel (placeholder)
└── ThemeEditor (overlay)
```

### Key Decisions

- **Composed model over flat model**: Each section is its own `tea.Model` under the `Section` interface — cleaner separation and easier to test in isolation.
- **Theme struct refactoring is in-scope**: Required by the runtime theme editor. Package vars can't be mutated safely at runtime.
- **No backend integration yet**: MCP/LSP sections are pure UI slots with typed messages. Polling comes later.

## Risks and Mitigations

| Risk | Mitigation |
|------|-----------|
| Sidebar too narrow for theme editor | Keep sidebar at ~30 chars for status. Theme editor gets a full-screen overlay. |
| Sections don't fit vertically | Make sections skippable/scrollable. Use `lipgloss.Height()` to measure. |
| Spinner sections consume CPU | Bubbles spinner uses `time.Tick` — batch under shared tick if needed. |
| Theme refactor breaks existing styles | Keep backward compat with `DefaultCatppuccinMocha` constant + derive function. |

## Delivery Plan

Given **auto-chain** strategy with **feature-branch-chain**:

### Potential slices (ordered)

1. **Theme struct refactoring** — convert `theme.go` vars to `Theme` struct + `DeriveStyles()`. Test: styles render correctly.
2. **Root App model** — `app.go` + `main_panel.go` placeholder. Wire `NewApp()`. Test: app initializes and renders layout.
3. **Sidebar model + Section interface** — `sidebar.go` + `section.go`. Empty sidebar renders. Test: section ordering, height.
4. **Individual sections** — one per sub-slice or batched: Session, Progress (spinner), MCP, LSP, Context, Files.
5. **Theme editor overlay** — keybind, textinput-based color editing, live preview.

**Estimated total**: ~250-400 lines of Go across all slices.

## Rollback

Trivial — only `internal/ui/` is affected, and `cmd/agy-tui/main.go` is already wired to call `ui.NewApp()`. Revert removes all new components; existing backend is untouched.
