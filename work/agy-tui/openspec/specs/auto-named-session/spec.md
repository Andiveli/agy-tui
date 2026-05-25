# Spec: auto-named-session

## Description

Session lifecycle with auto-naming from the first user prompt and continuation via agy's `--continue` flag, so the chat feels like a persistent conversation.

## Requirements

### R1: Session auto-naming
- First user prompt in a session derives a human-readable name from the first ~4 words, slugified (lowercased, spaces → hyphens, remove special chars)
- Example: `"How do I fix the N+1 query in Go?"` → `"how-do-i-fix"`
- Session name is set immediately on first prompt (before response arrives)
- A new SDD cycle resets the session (no persistence across restarts)

### R2: Conversation continuation
- Prompt 1: calls `Client.StartStreaming()` (fresh invocation)
- Prompt 2+: calls `Client.ContinueLastStreaming()` (reuses last conversation with `--continue`)
- Both methods stream output line-by-line with the same chunk protocol

### R3: SessionInfo display
- SessionInfo section in the sidebar shows the session name
- Section implements `UpdatableSection` to receive `SessionChangedMsg`
- When session name updates, the section refreshes immediately

## Acceptance Criteria

- [AC1] After sending first prompt, sidebar shows a slugified name derived from the prompt text
- [AC2] Second prompt reuses the same session via `--continue` (not a fresh agy invocation)
- [AC3] Session name in sidebar stays consistent across all turns in the same session
- [AC4] An error in the first invocation does not prevent continuation
- [AC5] No errors logged when `--continue` is called for the first time

## Edge Cases

| Case | Expected Behavior |
|------|-------------------|
| First prompt < 4 words | Use all available words |
| Prompt has special chars (`!@#$`) | Strip special chars from slug |
| First prompt is empty | No session starts |
| agy `--continue` fails (no prior session) | Fall back to StartStreaming? For now, let it fail — edge case |
| Session name collision | Cosmetic only, no functional impact |
