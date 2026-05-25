# Tasks: auto-named-session

## Budget

Estimated: ~50 lines changed across 3 files. Well under 400-line threshold.
Delivery: single PR, feature-branch-chain not needed.

## Tasks

### Task 1: Add `ContinueLastStreaming` to Client

**File**: `internal/backend/agy.go`
**Change**: Add method that mirrors `StartStreaming` with `--continue` flag
**Lines**: ~25 new
**Test**: Manual — verify `ContinueLastStreaming` instantiates a command with `--continue` in args

### Task 2: Wire auto-naming + continuation in ChatModel

**File**: `internal/ui/chat.go`
**Changes**:
- Add `sessionName string` field
- Add `slugify()` function
- In `enter` handler: derive name on first prompt, emit `SessionChangedMsg`
- In `startStream`: call `ContinueLastStreaming` when `promptCount >= 1`
**Lines**: ~20 new/modified

### Task 3: Make SessionInfo an UpdatableSection

**File**: `internal/ui/sections/session.go`
**Change**: Add `Update(tea.Msg)` method handling `SessionChangedMsg`
**Lines**: ~10 new
**Test**: Verify section receives and displays session name after `Update`

## Verification

1. `go build ./...` — compiles
2. `go vet ./...` — no issues
3. `go test ./...` — all existing tests pass
4. Manual: run `go run ./cmd/agy-tui/` and observe session name in sidebar
