# Tasks: Sidebar Data Flow

**Change**: `sidebar-data-flow`

| ID | Task | Files | LOC |
|----|------|-------|-----|
| T1 | Add ChatCompletedMsg to kit.go | `kit/kit.go` | +10 |
| T2 | Update ChatModel to emit ChatCompletedMsg with file paths | `chat.go` | +40 |
| T3 | Add startup probing in App.Init (MCP, LSP, CWD) | `app.go` | +50 |
| T4 | App routes domain msgs to Sidebar sections | `app.go` | +20 |

**Total**: ~120 lines. Single PR.
