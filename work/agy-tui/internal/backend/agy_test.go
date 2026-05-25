package backend

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient()
	if c.BinaryPath != "agy" {
		t.Errorf("BinaryPath = %q, want %q", c.BinaryPath, "agy")
	}
	if c.Timeout != DefaultTimeout {
		t.Errorf("Timeout = %v, want %v", c.Timeout, DefaultTimeout)
	}
}

func TestFormatTimeout(t *testing.T) {
	tests := []struct {
		input    time.Duration
		expected string
	}{
		{30 * time.Second, "30s"},
		{5 * time.Minute, "5m0s"},
		{90 * time.Second, "1m30s"},
		{1 * time.Hour, "1h0m0s"},
		{0, "0s"},
	}

	for _, tt := range tests {
		result := formatTimeout(tt.input)
		if result != tt.expected {
			t.Errorf("formatTimeout(%v) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestCheckBinary(t *testing.T) {
	c := NewClient()

	// agy might be installed or not — we just verify the method returns without panic
	err := c.CheckBinary()

	if err != nil {
		// If it fails, it must mention the binary name
		if !strings.Contains(err.Error(), "agy") {
			t.Errorf("CheckBinary error should mention agy, got: %v", err)
		}
	} else {
		// If it passes, we know agy is available
		t.Log("agy binary found in PATH")
	}
}

func TestCheckBinary_CustomPath(t *testing.T) {
	// Non-existent binary should fail
	c := &Client{BinaryPath: "nonexistent-cli-xyz-123"}
	err := c.CheckBinary()
	if err == nil {
		t.Error("expected error for non-existent binary, got nil")
	}
}

func TestSendPrompt_BinaryNotFound(t *testing.T) {
	c := &Client{BinaryPath: "nonexistent-cli-xyz-123", Timeout: time.Second}
	_, err := c.SendPrompt(context.Background(), "test")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed") {
		t.Errorf("error should mention failure, got: %v", err)
	}
}

func TestStartStreaming_BinaryNotFound(t *testing.T) {
	c := &Client{BinaryPath: "nonexistent-cli-xyz-123", Timeout: time.Second}
	_, err := c.StartStreaming(context.Background(), "test")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "agy start") {
		t.Errorf("error should mention 'agy start', got: %v", err)
	}
}

func TestContinueLastStreaming_Delegates(t *testing.T) {
	c := &Client{BinaryPath: "nonexistent-cli-xyz-123", Timeout: time.Second}
	_, err := c.ContinueLastStreaming(context.Background(), "test")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// Should fail with same error type as StartStreaming
	if !strings.Contains(err.Error(), "agy start") {
		t.Errorf("error should mention 'agy start', got: %v", err)
	}
}

func TestContinueConversationStreaming_Delegates(t *testing.T) {
	c := &Client{BinaryPath: "nonexistent-cli-xyz-123", Timeout: time.Second}
	_, err := c.ContinueConversationStreaming(context.Background(), "conv-id", "test")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "agy start") {
		t.Errorf("error should mention 'agy start', got: %v", err)
	}
}
