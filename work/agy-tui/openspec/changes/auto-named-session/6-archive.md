# Archive: auto-named-session

## Summary

Minimal change (~45 lines, 3 files) adding session continuity to the agy-tui TUI:

### Changes Delivered

**`internal/backend/agy.go`** — Added `ContinueLastStreaming()` method, mirrors `StartStreaming()` with `--continue` flag.

**`internal/ui/chat.go`** — Replaced `sessionID` with `sessionName`, derived from first prompt via `slugify()`. On prompt 1: sets name + emits `SessionChangedMsg`. On prompt 2+: calls `ContinueLastStreaming()` instead of `StartStreaming()`.

**`internal/ui/sections/session.go`** — Added `Update(tea.Msg)` handling `SessionChangedMsg`. SessionInfo now implements `kit.UpdatableSection`.

### Metrics
- Files changed: 3
- Lines added: ~45
- Tests: 27 passing (0 regressions)
- Build: clean

### Artifacts
- `openspec/changes/auto-named-session/1-proposal.md`
- `openspec/changes/auto-named-session/2-spec.md`
- `openspec/changes/auto-named-session/3-design.md`
- `openspec/changes/auto-named-session/4-tasks.md`
- `openspec/changes/auto-named-session/5-verify.md`
- `openspec/changes/auto-named-session/6-archive.md`
- `openspec/specs/auto-named-session/spec.md`

### Next Steps
1. MCP/LSP polling (refresh status periodically, not just at startup)
2. Status bar with agy info (tasks, artifacts) + key hints
3. Error handling polish (agy not installed, timeouts)
