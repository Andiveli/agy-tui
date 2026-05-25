# Design: auto-named-session

## Architecture

Minimal change — no new types, no new files. Three file modifications.

## File-by-File Design

### 1. `internal/backend/agy.go` — Add `ContinueLastStreaming`

```go
func (c *Client) ContinueLastStreaming(ctx context.Context, prompt string) (io.ReadCloser, error) {
    ctx, cancel := context.WithTimeout(ctx, c.Timeout)
    cmd := exec.CommandContext(ctx, c.BinaryPath,
        "--print", prompt,
        "--print-timeout", formatTimeout(c.Timeout),
        "--continue",
    )
    stdout, err := cmd.StdoutPipe()
    // ...same error handling as StartStreaming
    return &cmdReadCloser{ReadCloser: stdout, cancel: cancel, cmd: cmd}, nil
}
```

Same as `StartStreaming` but appends `"--continue"` to args.

### 2. `internal/ui/chat.go` — Auto-naming + continuation logic

**New fields on ChatModel:**
```go
sessionName  string   // set once from first prompt
```

**Flow changes:**

In `Update` → `enter` key handler:
- On first prompt (promptCount == 0): derive `sessionName` via `slugify(prompt)`, emit `SessionChangedMsg` with the name
- In `startStream`: if `promptCount > 0`, call `client.ContinueLastStreaming` instead of `client.StartStreaming`

**New function:**
```go
func slugify(prompt string) string {
    // Take first 4 words, lowercase, replace spaces with hyphens, strip non-alphanumeric
}
```

### 3. `internal/ui/sections/session.go` — UpdatableSection

```go
func (s *SessionInfo) Update(msg tea.Msg) (kit.Section, tea.Cmd) {
    switch v := msg.(type) {
    case kit.SessionChangedMsg:
        if v.Name != "" {
            s.name = v.Name
        }
        if v.Context != "" {
            s.context = v.Context
        }
    }
    return s, nil
}
```

This makes SessionInfo implement `kit.UpdatableSection` — the existing sidebar routing already handles it.

## Message Flow

```
User presses Enter (prompt 1)
  → ChatModel.Update:
    1. slugify(prompt) → sessionName = "how-do-i-fix"
    2. loading = true
    3. return tea.Batch(
         emitProgress("running", 0),
         emitSession(sessionName, cwd),
         startStream(prompt),  // calls StartStreaming
       )

User presses Enter (prompt 2+)
  → ChatModel.Update:
    1. loading = true
    2. return tea.Batch(
         emitProgress("running", 0),
         startStream(prompt),  // calls ContinueLastStreaming
       )

→ SessionChangedMsg → App.Update → sidebar.Update → SessionInfo.Update
```

## Risks

- `--continue` before agy has any conversation → agy may error. Unlikely in practice since prompt 2+ implies prompt 1 succeeded.
