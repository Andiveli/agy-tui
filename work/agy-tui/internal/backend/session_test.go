package backend

import (
	"testing"
)

func TestNewSessionManager(t *testing.T) {
	sm := NewSessionManager()
	if sm == nil {
		t.Fatal("NewSessionManager returned nil")
	}
	if sm.CurrentSession() != nil {
		t.Error("expected nil current session before StartSession")
	}
}

func TestStartSession(t *testing.T) {
	sm := NewSessionManager()
	s := sm.StartSession()
	if s == nil {
		t.Fatal("StartSession returned nil")
	}
	if s.ID != "" {
		t.Errorf("expected empty ID for new session, got %q", s.ID)
	}
	if s.StartedAt.IsZero() {
		t.Error("expected StartedAt to be set")
	}
	// CurrentSession should return the same session
	if sm.CurrentSession() != s {
		t.Error("CurrentSession should return the session created by StartSession")
	}
}

func TestSetConversationID_AutoStarts(t *testing.T) {
	sm := NewSessionManager()
	sm.SetConversationID("conv-123")
	if sm.CurrentSession() == nil {
		t.Fatal("expected session to be auto-started")
	}
	if sm.CurrentSession().ID != "conv-123" {
		t.Errorf("expected ID 'conv-123', got %q", sm.CurrentSession().ID)
	}
}

func TestSetConversationID_UpdatesExisting(t *testing.T) {
	sm := NewSessionManager()
	sm.StartSession()
	sm.SetConversationID("conv-456")
	if sm.CurrentSession().ID != "conv-456" {
		t.Errorf("expected ID 'conv-456', got %q", sm.CurrentSession().ID)
	}
}

func TestListConversations_NoError(t *testing.T) {
	sm := NewSessionManager()
	// Should never return an error — missing dir returns nil, nil
	ids, err := sm.ListConversations()
	if err != nil {
		t.Errorf("ListConversations returned unexpected error: %v", err)
	}
	// ids may be nil or non-nil depending on whether agy conversations exist
	_ = ids
}

func TestSession_StartedAtSet(t *testing.T) {
	sm := NewSessionManager()
	s := sm.StartSession()
	if s.StartedAt.IsZero() {
		t.Error("StartedAt should not be zero")
	}
}

func TestCurrentSession_NilBeforeStart(t *testing.T) {
	sm := NewSessionManager()
	if sm.CurrentSession() != nil {
		t.Error("CurrentSession should be nil before StartSession")
	}
}

func TestMultipleSessions(t *testing.T) {
	sm := NewSessionManager()
	s1 := sm.StartSession()
	s1.ID = "first"

	s2 := sm.StartSession()
	s2.ID = "second"

	// StartSession replaces the current session
	if sm.CurrentSession().ID != "second" {
		t.Errorf("expected current session ID 'second', got %q", sm.CurrentSession().ID)
	}
}
