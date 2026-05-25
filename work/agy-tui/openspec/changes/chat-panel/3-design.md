# Design: Real Chat Panel

**Change**: `chat-panel`
**Status**: ✅ Designed
**Date**: 2026-05-25

---

## 1. ChatMessage Type

Add to `internal/ui/kit/kit.go`:

```go
type ChatMessage struct {
    Role    string // "user" or "agent"
    Content string
}
```

## 2. ChatModel

New file `internal/ui/chat.go`:

```go
type ChatModel struct {
    viewport  viewport.Model
    input     textinput.Model
    messages  []kit.ChatMessage
    loading   bool
    styles    kit.Styles
    client    *backend.Client   // agy backend
    sessionID string            // current agy conversation ID
}

func NewChatModel(styles kit.Styles, client *backend.Client) ChatModel { ... }
func (c *ChatModel) Init() tea.Cmd { return textinput.Blink }
func (c *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { ... }
func (c *ChatModel) View() string { ... }
```

### Update routing

| Message | Action |
|---------|--------|
| tea.WindowSizeMsg | Resize viewport + input |
| tea.KeyMsg{Enter} | Send prompt (if !loading) |
| chatResponseMsg | Append response to viewport, clear loading |
| chatErrorMsg | Show error message, clear loading |
| textinput.KeyMsg | Forward to input |

### Async response handling

```go
type chatResponseMsg struct {
    Content string
}

type chatErrorMsg struct {
    Err error
}

func (c *ChatModel) sendPrompt(prompt string) tea.Cmd {
    return func() tea.Msg {
        ctx := context.Background()
        resp, err := c.client.SendPrompt(ctx, prompt)
        if err != nil {
            return chatErrorMsg{Err: err}
        }
        return chatResponseMsg{Content: resp}
    }
}
```

### View layout

```
+----------------------------------+
|  [viewport: scrollable messages] |
|                                  |
|  > Type a prompt...              |
+----------------------------------+
```

Viewport and input separated by a lipgloss border or gap.

## 3. App Integration

Modify `app.go`:
- Replace `mainPanel MainPanel` with `chat *ChatModel`
- Add `backend.Client` as dependency (created in NewApp)
- Route key messages to chat when not handled by app
- Pass WindowSizeMsg to chat with correct width

## 4. File Changes

| File | Action |
|------|--------|
| `internal/ui/kit/kit.go` | Add ChatMessage type |
| `internal/ui/chat.go` | NEW — ChatModel |
| `internal/ui/app.go` | MODIFY — replace MainPanel with ChatModel |
| `internal/ui/main_panel.go` | REMOVE or keep as dead code |

## 5. Error Handling

| Condition | Behavior |
|-----------|----------|
| agy timeout | Show "Request timed out" in viewport |
| agy not found | Show "agy binary not found" |
| Empty prompt | Ignore Enter |
| Loading state | Disable input, ignore Enter |
