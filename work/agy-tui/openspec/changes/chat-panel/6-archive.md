# Archive: Real Chat Panel

**Change**: `chat-panel`
**Status**: ✅ Archived
**Date**: 2026-05-25

## Summary

Replaced MainPanel placeholder with a functional chat panel. Users can type prompts, send to agy backend, and see responses in a scrollable viewport. Async via goroutine with loading indicator.

## Artifacts

| Phase | Openspec | Engram |
|-------|----------|--------|
| Proposal | `1-proposal.md` | ✅ |
| Spec | `2-spec.md` | ✅ |
| Design | `3-design.md` | ✅ |
| Tasks | `4-tasks.md` | ✅ |
| Apply | — | ✅ |
| Verify | `5-verify.md` | ✅ |
| Archive | `6-archive.md` | ✅ |

## Files

```
MOD internal/ui/kit/kit.go        — Added ChatMessage, ChatResponseMsg, ChatErrorMsg
NEW internal/ui/chat.go           — ChatModel (viewport + input + async send)
MOD internal/ui/app.go            — Replaced MainPanel with ChatModel
```

## Remaining

1. Connect sidebar data flow (session info, MCP/LSP status)
2. Tests (teatest + golden files)
3. Streaming responses
4. Theme editor wiring
