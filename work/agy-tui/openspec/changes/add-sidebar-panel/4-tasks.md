# Tasks: Right-side Sidebar Panel

**Change**: `add-sidebar-panel`
**Status**: ✅ Tasked
**Date**: 2026-05-25

---

## Executive Summary

13 tasks across 5 slices for feature-branch-chain. Total ~380 lines. Fits within the 400-line budget, so a single PR is viable, but the chain structure is maintained for clean review boundaries.

---

## Task Breakdown

### Slice 1: Theme System Refactoring (tasks 1-2)
*PR target: feature/agy-tui/sidebar-panel (tracker branch)*

| ID | Task | Files | Deps | Est. LOC | Tests? |
|----|------|-------|------|----------|--------|
| T1 | Refactor theme.go: vars → Theme struct + DeriveStyles() | `internal/ui/theme.go` | none | +20/-20 (net 0) | Yes |
| T2 | Add DefaultCatppuccinMocha constant + backward compat | `internal/ui/theme.go` | T1 | +25 | Yes |

**Verification**: `DeriveStyles(DefaultCatppuccinMocha)` produces styles identical to old vars. Tests pass.

---

### Slice 2: Core Interfaces + Messages (tasks 3-4)
*PR targets: previous PR branch*

| ID | Task | Files | Deps | Est. LOC | Tests? |
|----|------|-------|------|----------|--------|
| T3 | Create section.go: Section + UpdatableSection interfaces | `internal/ui/section.go` | none | +35 | Yes |
| T4 | Create messages.go: all domain message types | `internal/ui/messages.go` | none | +20 | No (types only) |

**Verification**: Interfaces compile, message types are constructable.

---

### Slice 3: Placeholder Components (tasks 5-7)
*PR targets: previous PR branch*

| ID | Task | Files | Deps | Est. LOC | Tests? |
|----|------|-------|------|----------|--------|
| T5 | Create MainPanel placeholder model | `internal/ui/main_panel.go` | T1 (styles) | +30 | Yes |
| T6 | Create SessionInfo section | `internal/ui/sections/session.go` | T3, T4 | +25 | Yes |
| T7 | Create SubAgentProgress section (with bubbles/spinner) | `internal/ui/sections/progress.go` | T3, T4 | +40 | Yes |

**Verification**: MainPanel renders placeholder, sections render with headers.

---

### Slice 4: Remaining Sections + Sidebar Model (tasks 8-11)
*PR targets: previous PR branch*

| ID | Task | Files | Deps | Est. LOC | Tests? |
|----|------|-------|------|----------|--------|
| T8 | Create MCPStatus section | `internal/ui/sections/mcp.go` | T3, T4 | +20 | Yes |
| T9 | Create LSPStatus section | `internal/ui/sections/lsp.go` | T3, T4 | +20 | Yes |
| T10 | Create ContextPanel section | `internal/ui/sections/context.go` | T3, T4 | +30 | Yes |
| T11 | Create FileChanges section | `internal/ui/sections/files.go` | T3, T4 | +30 | Yes |

**Verification**: All sections render, respond to their respective messages.

---

### Slice 5: Sidebar + App Root + Theme Editor (tasks 12-14)
*PR targets: previous PR branch → final merge to tracker → main*

| ID | Task | Files | Deps | Est. LOC | Tests? |
|----|------|-------|------|----------|--------|
| T12 | Create Sidebar model (sections slice, height calc, toggle) | `internal/ui/sidebar.go` | T5-T11 | +45 | Yes |
| T13 | Create App root model (NewApp, layout, msg routing) | `internal/ui/app.go` | T1, T2, T5, T12 | +40 | Yes |
| T14 | Create ThemeEditor overlay model | `internal/ui/theme_editor.go` | T1, T2 | +50 | Yes |

**Verification**: `go build ./cmd/agy-tui` succeeds. `go test ./internal/ui/...` passes.

---

## Slice Dependency Graph

```
Slice 1 (Theme) ──→ Slice 2 (Interfaces) ──→ Slice 3 (Placeholders) ──→ Slice 4 (Sections) ──→ Slice 5 (Integration)
```

---

## Review Workload Forecast

| Metric | Value |
|--------|-------|
| Total estimated LOC | ~380 lines |
| 400-line budget risk | Low (~380 < 400) |
| Chained PRs recommended | No (within budget) |
| Decision needed before apply | No |

**Recommendation**: Single PR is viable. However, slices 1-4 can be reviewed independently for clarity. Recommend merging as a single PR to the feature branch given the low total.

---

## Rollback Plan

Each slice is independently revertible:
- Revert Slice 1: restore `theme.go` vars, command `git revert <theme-commit>`
- Revert Slice 5: `git revert <app-commit>`, sidebar disappears, main panel remains
- Full rollback: `git revert <theme-commit>..<last-commit>` — restores pre-change state
