# Proposal: Sidebar Data Flow

**Change**: `sidebar-data-flow`
**Status**: ✅ Proposed
**Date**: 2026-05-25

---

## Intent

Sidebar sections show static defaults ("Disconnected", "default", "no workspace") after chat. Users see no session info, no file changes, no agent progress. Wire real data into each section by emitting domain messages from ChatModel after agy responses and at startup.

## Scope

### In Scope
- ChatModel emits `SessionChangedMsg` after first agy response (sets session ID from SessionManager)
- ChatModel emits `ProgressMsg` ("sending..." during agy call, "completed" on response)
- ChatModel emits `FileChangedMsg` by parsing agy response text for file paths
- App checks MCP config file (`~/.gemini/antigravity-cli/mcp_config.json`) at startup, emits `MCPStatusMsg`
- App checks common language servers (gopls, typescript-language-server) via `pgrep` at startup, emits `LSPStatusMsg`
- ContextPanel shows current workspace directory at startup

### Out of Scope
- Real-time filesystem watching (fsnotify)
- Full MCP/LSP protocol integration
- agy sub-agent status polling (API doesn't expose it)
- Streaming sidebar updates

## Capabilities

No new capabilities — this change modifies existing specs.

### New Capabilities
None

### Modified Capabilities
- `chat-panel`: Must emit domain messages (`SessionChangedMsg`, `ProgressMsg`, `FileChangedMsg`) after agy responses to wire sidebar sections
- `sidebar-panel`: Must receive live data at startup for MCP, LSP, and workspace context; add startup probes

## Approach

```
ChatModel.sendPrompt() → ChatResponseMsg{Content}
  ├── parse file paths with regex → FileChangedMsg
  ├── read session ID from client → SessionChangedMsg
  └── emit ProgressMsg("completed") → sidebar.Update()

App startup:
  ├── read ~/.gemini/antigravity-cli/mcp_config.json → MCPStatusMsg
  ├── pgrep gopls / typescript-language-server → LSPStatusMsg
  └── os.Getwd() → SessionChangedMsg(context=workspace)
```

All domain messages flow through `App.Update()` → `sidebar.Update()` → `UpdatableSection.Update()`. No new message types needed — existing `kit.MCPStatusMsg`, `kit.LSPStatusMsg`, `kit.ProgressMsg`, `kit.FileChangedMsg`, `kit.SessionChangedMsg` are already handled by sections.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/ui/chat.go` | Modified | Add domain message emission in `sendPrompt` response handler |
| `internal/ui/app.go` | Modified | Add startup probes after `NewApp`; wire probe results into sidebar |
| `internal/backend/agy.go` | Unchanged | Used as-is |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| File path parsing heuristic misses paths | Medium | Use regex `[\w/.-]+\.[a-z]+`; iteratively improve |
| mcp_config.json format varies | Low | Graceful degradation — log parse error, default to Disconnected |
| LSP process check Linux/macOS only | Low | Wrap with `runtime.GOOS` check |

## Rollback

Revert `chat.go` and `app.go`. No schema or data migration. Sidebar returns to static defaults.

## Dependencies

None — uses existing `os`, `os/exec`, `regexp`, `runtime` from stdlib.

## Success Criteria

- [ ] SessionInfo shows real session name after first agy response
- [ ] SubAgentProgress shows spinner during agy call, checkmark on completion
- [ ] FileChanges shows files referenced in agy responses
- [ ] MCPStatus reflects mcp_config.json presence
- [ ] ContextPanel shows real workspace directory
- [ ] LSPStatus shows connected when known servers are running
