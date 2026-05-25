# Spec: Sidebar Data Flow

**Change**: `sidebar-data-flow`
**Status**: ✅ Specified

## Requirements

### 1. Session Info
- After first agy response, emit SessionChangedMsg with session name
- SessionManager tracks conversation IDs from agy responses

### 2. Sub-Agent Progress
- When chat sends a prompt → emit ProgressMsg{running}
- When chat receives response → emit ProgressMsg{completed}
- Track at least one "agent" entry in the sidebar

### 3. File Changes
- Parse agy response text for file paths (lines containing `path/to/file.ext:` pattern)
- Emit FileChangedMsg for each detected file
- Ignore common non-file lines (error messages, status lines)

### 4. MCP Status
- At startup, read `~/.gemini/antigravity-cli/mcp_config.json`
- If file exists and has servers → MCPStatusMsg{Connected: true}
- Emit on app init

### 5. LSP Status
- At startup, check running processes for common LSPs (gopls, typescript-language-server, etc.)
- If any found → LSPStatusMsg{Connected: true}
- Linux/macOS only (uses `pgrep`)

### 6. Context Panel
- Show current workspace directory
- Emit SessionChangedMsg with context = CWD
