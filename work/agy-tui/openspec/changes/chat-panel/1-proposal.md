# Proposal: Real Chat Panel

**Change**: `chat-panel`
**Status**: ✅ Proposed
**Date**: 2026-05-25

---

## Intent

MainPanel currently shows "Chat area — coming soon". Users can't send prompts or see responses. We need a functional chat panel: scrollable message history + text input connected to the agy CLI backend.

## Scope

### In Scope
- Viewport (bubbles/viewport) for scrolling message history
- Text input (bubbles/textinput) at panel bottom for composing prompts
- User messages styled in blue, agent responses in normal text
- Enter to send prompt, PgUp/PgDown/↑/↓ for viewport scrolling
- Loading indicator while agy processes the prompt
- Connection to `backend.SendPrompt()` via goroutine + channel
- Session awareness: fresh session on startup, `ContinueLast` for subsequent turns
- Replace `MainPanel` struct type with `ChatModel` — update type in `App`

### Out of Scope
- Streaming responses (future)
- Markdown / glamour rendering (future)
- File attachments
- Sidebar data flow integration (separate change)
- Persistent conversation history (future)

## Capabilities

No existing specs — all capabilities are new.

### New Capabilities
- `chat-viewport`: Scrollable message history using bubbles/viewport with styled messages (user=blue, agent=normal, system=italic yellow)
- `prompt-input`: Text input at panel bottom with Enter-to-send, integrated with backend
- `chat-panel`: Composite tea.Model owning viewport + input + loading state + backend client + session tracking

### Modified Capabilities
- None

## Approach

```
ChatModel (replaces MainPanel)
├── viewport.Model        — scrollable message history
├── textinput.Model       — prompt composition
├── []Message             — in-memory message log
├── backend.Client        — agy prompt execution
├── SessionManager        — track/continue conversations
└── loading bool          — show spinner while awaiting response
```

On Enter: disable input, show spinner, launch goroutine → `backend.SendPrompt()`, send result back via `tea.Cmd` (channel message). Append both prompt + response to `[]Message`, re-render viewport. On WindowSizeMsg: resize viewport height minus input line.

`App.Update` currently routes sidebar msgs but ignores `mainPanel.Update` result — need to thread `mainPanel` update cmds.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/ui/main_panel.go` | Replaced | Rewrite as `ChatModel` composing viewport + textinput |
| `internal/ui/app.go` | Modified | `MainPanel` → `ChatModel` type; thread `mainPanel.Update` cmds into batch |
| `internal/backend/agy.go` | Unchanged | `SendPrompt` / `ContinueLast` used as-is |
| `internal/ui/kit/kit.go` | Unchanged | Styles already defined: `UserMessage`, `AgentMessage`, `InputPrompt` |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| agy CLI slow (up to 5 min) blocks UI | High | Goroutine + chan msg — UI stays responsive with loading indicator |
| Textinput + viewport key collision | Med | Route keys: viewport gets PgUp/PgDown/↑/↓ when input unfocused; Enter always goes to input |
| New goroutine not cancelled on quit | Low | Use `ctx` from `tea.Quit` signal, cancel context in `chatModel.Update` |

## Rollback

Revert `internal/ui/main_panel.go` to original placeholder (3 files: main_panel.go + app.go changes). Backend untouched. `openspec/changes/chat-panel/` deleted. No schema or data migration needed.

## Dependencies

- `github.com/charmbracelet/bubbles` — already imported (viewport + textinput)

## Success Criteria

- [ ] User types a prompt, presses Enter, prompt appears in viewport in blue
- [ ] Loading indicator visible while agy runs
- [ ] Response appears in normal text after agy finishes
- [ ] Scroll works (PgUp/PgDown/↑/↓) for history
- [ ] Subsequent prompts continue the same conversation
