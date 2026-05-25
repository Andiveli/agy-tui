# Proposal: Auto-Named Session

## Intent

Users start each conversation cold — agy runs as a fresh invocation every time, with no continuity and no way to identify sessions. We need auto-naming from the first prompt and invisible session continuation so the chat feels like a persistent conversation.

## Scope

### In Scope
- `ContinueLastStreaming()` on `backend.Client` — mirrors `StartStreaming` but passes `--continue`
- Auto-derived session name from first prompt (slugify first ~4 words)
- `ChatModel` uses `ContinueLastStreaming` for prompt 2+
- `SessionInfo` implements `UpdatableSection` receiving `SessionChangedMsg` with name updates

### Out of Scope
- Session listing / picker UI
- History browser / conversation search
- `SessionManager` integration (future work)
- Disk persistence beyond agy's `--continue`
- New UI sections or components

## Capabilities

### New Capabilities
- `auto-named-session`: Session lifecycle + auto-naming. First prompt derives a slug from its first ~4 words; subsequent prompts reuse the same session via `--continue`. Session name propagates to the sidebar `SessionInfo` through `UpdatableSection`.

### Modified Capabilities
- None

## Approach

1. **`ContinueLastStreaming(ctx, onRead)`** in `internal/backend/agy.go` — same signature as `StartStreaming` but runs `agy --continue` instead of `agy <prompt>`
2. **`ChatModel` auto-name** in `internal/ui/chat.go` — on first prompt, slugify first ~4 words → session name; on prompt 2+, call `ContinueLastStreaming`; publish `SessionChangedMsg`
3. **`SessionInfo` → `UpdatableSection`** in `internal/ui/sections/session.go` — implement `Update(SessionChangedMsg)` to update the displayed session name

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/backend/agy.go` | New method | `ContinueLastStreaming()` |
| `internal/ui/chat.go` | Modified | Session name derivation + dispatch logic |
| `internal/ui/sections/session.go` | Modified | Implement `UpdatableSection` |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| `ContinueLastStreaming` before any session | Low | Only called on prompt 2+; first prompt always uses `StartStreaming` |
| Slug collision on identical prompts | Low | Cosmetic only — not a primary key |

## Rollback Plan

Revert 3 files: `agy.go` (remove method), `chat.go` (always-`StartStreaming`), `session.go` (revert `UpdatableSection`). No data migration.

## Dependencies

- agy CLI `--continue` flag at the installed version

## Success Criteria

- [ ] First prompt generates a slugged session name visible in sidebar
- [ ] Prompt 2+ uses `--continue`, not a fresh invocation
- [ ] Session name persists across turns within the same session
