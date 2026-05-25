# Verify: Sidebar Data Flow

**Change**: `sidebar-data-flow`
**Status**: ✅ Verified

| Check | Result |
|-------|--------|
| go build | ✅ |
| go vet | ✅ |

## Spec Compliance
- [x] ChatModel emits SessionChangedMsg after response
- [x] ChatModel emits ProgressMsg (running → completed)
- [x] ChatModel emits FileChangedMsg for parsed file paths
- [x] App probes MCP config at startup → MCPStatusMsg
- [x] App probes LSP processes at startup → LSPStatusMsg
- [x] App probes workspace at startup → SessionChangedMsg
- [x] App routes domain msgs to sidebar sections
