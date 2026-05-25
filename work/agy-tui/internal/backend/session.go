package backend

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Session represents a single conversation session.
type Session struct {
	ID         string
	StartedAt  time.Time
	LastPrompt string
}

// SessionManager tracks the active session and discovers existing conversations.
type SessionManager struct {
	current *Session
}

// NewSessionManager returns an initialized SessionManager.
func NewSessionManager() *SessionManager {
	return &SessionManager{}
}

// CurrentSession returns the active session, or nil if none has been started.
func (sm *SessionManager) CurrentSession() *Session {
	return sm.current
}

// StartSession creates and stores a new session, replacing any previous one.
func (sm *SessionManager) StartSession() *Session {
	sm.current = &Session{
		StartedAt: time.Now(),
	}
	return sm.current
}

// SetConversationID assigns a conversation ID to the current session.
// If no session exists, it starts one first.
func (sm *SessionManager) SetConversationID(id string) {
	if sm.current == nil {
		sm.StartSession()
	}
	sm.current.ID = id
}

// ListConversations reads .pb files from the agy conversations directory
// and returns their IDs (filenames without the .pb extension).
func (sm *SessionManager) ListConversations() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(home, ".gemini", "antigravity-cli", "conversations")

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var ids []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".pb") {
			ids = append(ids, strings.TrimSuffix(name, ".pb"))
		}
	}
	return ids, nil
}
