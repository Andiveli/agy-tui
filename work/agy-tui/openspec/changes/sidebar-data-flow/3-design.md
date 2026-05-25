# Design: Sidebar Data Flow

**Change**: `sidebar-data-flow`

## 1. ChatModel Changes

After successful agy response in ChatModel.Update:
```go
case kit.ChatResponseMsg:
    // ... existing response handling ...
    // Emit sidebar messages
    sessionMsg := kit.SessionChangedMsg{Name: "session-1", Context: cwd}
    progressMsg := kit.ProgressMsg{SubAgent: "agy", Status: "completed", Progress: 100}
    // Parse files from response
    for _, line := range parseFilePaths(msg.Content) {
        fileMsg := kit.FileChangedMsg{Path: line, Action: "modified"}
        // ... emit ...
    }
```

**Problem**: ChatModel can't directly emit to Sidebar (they're separate models). Solution: return commands that produce domain messages. App.Update receives them and forwards to sidebar.

So ChatModel.sendPrompt returns a cmd that produces ChatResponseMsg. After that, we need chained commands to produce sidebar messages.

Better approach: ChatModel produces a composite response:

```go
type ChatCompletedMsg struct {
    Content   string
    SessionName string
    FilePaths []string
}
```

Replace ChatResponseMsg with ChatCompletedMsg that carries sidebar-relevant data.

## 2. App Startup Probing

In App.Init(), add commands that probe:
- MCP config file
- LSP processes  
- Workspace directory

These return domain messages that App routes to Sidebar.

## 3. File Changes

```go
func (c *ChatModel) sendPrompt(prompt string) tea.Cmd {
    return func() tea.Msg {
        ctx := context.Background()
        resp, err := c.client.SendPrompt(ctx, prompt)
        if err != nil {
            return kit.ChatErrorMsg{Err: err}
        }
        return kit.ChatCompletedMsg{
            Content:     resp,
            SessionName: "session-1", // simplified
            FilePaths:   parseFilePaths(resp),
        }
    }
}

func parseFilePaths(output string) []string {
    // Match lines like: path/to/file.go:123 or path/to/file.go
    // Simple heuristic: lines containing ".go:", ".ts:", etc.
    var paths []string
    for _, line := range strings.Split(output, "\n") {
        trimmed := strings.TrimSpace(line)
        if matched, _ := filepath.Match("*.go:", trimmed); matched {
            paths = append(paths, trimmed)
        }
    }
    return paths
}
```
