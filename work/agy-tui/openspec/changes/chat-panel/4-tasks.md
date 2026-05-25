# Tasks: Real Chat Panel

**Change**: `chat-panel`
**Status**: ✅ Tasked
**Date**: 2026-05-25

---

## Task Breakdown

| ID | Task | Files | Est. LOC |
|----|------|-------|----------|
| T1 | Add ChatMessage type to kit.go | `internal/ui/kit/kit.go` | +5 |
| T2 | Create ChatModel (viewport + input + send) | `internal/ui/chat.go` | +120 |
| T3 | Update App to use ChatModel | `internal/ui/app.go` | -20/+30 |
| T4 | Remove MainPanel placeholder | `internal/ui/main_panel.go` | -30 |

**Total**: ~105 lines added, ~50 removed, ~155 net

## Dependencies
T3 depends on T1 + T2. T4 independent.

## Review Workload
~155 lines net — well within 400-line budget. Single PR.
