# Archive: Sidebar Data Flow

**Change**: `sidebar-data-flow`
**Status**: ✅ Archived

## Summary
Connected all 6 sidebar sections with real data. ChatModel emits domain messages after agy responses. App probes MCP/LSP/workspace at startup.

## Files
```
MOD internal/ui/kit/kit.go   — ChatResponseMsg → ChatCompletedMsg
MOD internal/ui/chat.go      — Emit ProgressMsg, SessionChangedMsg, FileChangedMsg
MOD internal/ui/app.go       — Startup probing + domain msg routing
```

## Next
Change #3: Streaming responses
