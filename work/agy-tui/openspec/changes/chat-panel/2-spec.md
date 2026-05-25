# Spec: Real Chat Panel

**Change**: `chat-panel`
**Status**: ✅ Specified
**Date**: 2026-05-25

---

## 1. chat-viewport

### Requirements
- Uses `bubbles/viewport.Model` for scrollable message history
- Messages rendered top-to-bottom, newest at bottom
- Auto-scrolls to bottom on new messages
- Supports PgUp/PgDown/Home/End for navigation
- Width and height adapt to terminal size via `tea.WindowSizeMsg`

### Scenarios
| Scenario | Input | Expected |
|----------|-------|----------|
| New message received | Append to history | Viewport scrolls to bottom |
| User scrolls up | PgUp | Viewport scrolls up, new message doesn't auto-scroll |
| Terminal resized | WindowSizeMsg | Viewport dimensions update |
| Empty history | Init | Shows placeholder text |

### Acceptance Criteria
- Viewport renders all messages in order
- Auto-scroll enabled by default, disabled after manual scroll-up
- Resize recalculates content height

---

## 2. prompt-input

### Requirements
- Uses `bubbles/textinput.Model` at the bottom of the chat panel
- Placeholder text: "Type a prompt..."
- Enter sends the message and clears input
- Disabled while waiting for response
- Ctrl+C in empty input quits, Ctrl+C with text clears input

### Scenarios
| Scenario | Input | Expected |
|----------|-------|----------|
| Type prompt | typing | Text appears in input |
| Send | Enter | Input clears, message appears in viewport, loading shows |
| Send while loading | Enter | Ignored (input disabled) |
| Clear text | Ctrl+C with text | Input clears |
| Quit | Ctrl+C empty | App quits |

### Acceptance Criteria
- Input renders below viewport
- Enter triggers send
- Input clears after send
- Input disabled during loading

---

## 3. chat-panel

### Requirements
- Replaces MainPanel placeholder
- Composes viewport (top) + textinput (bottom) with lipgloss layout
- On send: creates `backend.SendPrompt(ctx, text)` in a goroutine, sends result back via `tea.Cmd`
- Shows loading indicator ("..." or spinner) while waiting
- Messages stored as `[]kit.ChatMessage{Role, Content}`
- User messages rendered in `UserMessage` style, agent in `AgentMessage` style

### Scenarios
| Scenario | Input | Expected |
|----------|-------|----------|
| Send prompt | Enter with text | "Sending..." shown, then response appears |
| agy returns error | Backend error | Error styled message in viewport |
| Session changed | New session | Chat history clears |

### Acceptance Criteria
- Messages render with correct styles per role
- agy backend is called asynchronously (non-blocking)
- Error responses don't crash the app
- Loading indicator visible during backend call
