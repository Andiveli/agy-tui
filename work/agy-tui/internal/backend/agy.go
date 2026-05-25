package backend

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

// DefaultTimeout is the maximum time allowed for a single agy invocation.
const DefaultTimeout = 5 * time.Minute

// Client wraps interactions with the agy CLI binary.
type Client struct {
	BinaryPath string
	Timeout    time.Duration
}

// NewClient returns a Client with sensible defaults.
func NewClient() *Client {
	return &Client{
		BinaryPath: "agy",
		Timeout:    DefaultTimeout,
	}
}

// SendPrompt runs `agy --print "<prompt>"` and returns the captured stdout.
func (c *Client) SendPrompt(ctx context.Context, prompt string) (string, error) {
	args := []string{
		"--print", prompt,
		"--print-timeout", formatTimeout(c.Timeout),
	}
	return c.run(ctx, args)
}

// ContinueConversation resumes an existing conversation.
func (c *Client) ContinueConversation(ctx context.Context, conversationID, prompt string) (string, error) {
	args := []string{
		"--print", prompt,
		"--print-timeout", formatTimeout(c.Timeout),
		"--conversation", conversationID,
	}
	return c.run(ctx, args)
}

// ContinueLast resumes the most recent conversation.
func (c *Client) ContinueLast(ctx context.Context, prompt string) (string, error) {
	args := []string{
		"--print", prompt,
		"--print-timeout", formatTimeout(c.Timeout),
		"--continue",
	}
	return c.run(ctx, args)
}

// StartStreaming runs agy and returns a ReadCloser for reading output incrementally.
func (c *Client) StartStreaming(ctx context.Context, prompt string) (io.ReadCloser, error) {
	return c.startStreaming(ctx, prompt)
}

// ContinueLastStreaming runs agy with --continue to continue the last conversation.
func (c *Client) ContinueLastStreaming(ctx context.Context, prompt string) (io.ReadCloser, error) {
	return c.startStreaming(ctx, prompt, "--continue")
}

// ContinueConversationStreaming runs agy with --conversation <id>.
func (c *Client) ContinueConversationStreaming(ctx context.Context, conversationID, prompt string) (io.ReadCloser, error) {
	return c.startStreaming(ctx, prompt, "--conversation", conversationID)
}

// startStreaming is the shared implementation for streaming agy invocations.
func (c *Client) startStreaming(ctx context.Context, prompt string, extraArgs ...string) (io.ReadCloser, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	args := []string{
		"--print", prompt,
		"--print-timeout", formatTimeout(c.Timeout),
	}
	args = append(args, extraArgs...)

	cmd := exec.CommandContext(ctx, c.BinaryPath, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	cmd.Stderr = cmd.Stdout // merge stderr into stdout for error visibility

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, fmt.Errorf("agy start: %w", err)
	}

	return &cmdReadCloser{
		ReadCloser: stdout,
		cancel:     cancel,
		cmd:        cmd,
	}, nil
}

// cmdReadCloser wraps stdout with context cancellation and process cleanup.
type cmdReadCloser struct {
	io.ReadCloser
	cancel context.CancelFunc
	cmd    *exec.Cmd
}

func (r *cmdReadCloser) Close() error {
	r.cancel()
	err := r.ReadCloser.Close()
	_ = r.cmd.Wait() // reap child process, discard error (already cancelled)
	return err
}

// run executes agy and returns captured stdout.
func (c *Client) run(ctx context.Context, args []string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, c.BinaryPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("agy timed out after %s", c.Timeout)
		}
		if ctx.Err() == context.Canceled {
			return "", fmt.Errorf("agy canceled: %w", ctx.Err())
		}
		stderrText := strings.TrimSpace(stderr.String())
		if stderrText != "" {
			return "", fmt.Errorf("agy failed: %s", stderrText)
		}
		return "", fmt.Errorf("agy failed: %w", err)
	}

	return strings.TrimSpace(stdout.String()), nil
}

// CheckBinary verifies that the agy binary is available in PATH.
func (c *Client) CheckBinary() error {
	_, err := exec.LookPath(c.BinaryPath)
	if err != nil {
		return fmt.Errorf("%s not found in $PATH: %w", c.BinaryPath, err)
	}
	return nil
}

func formatTimeout(d time.Duration) string {
	return d.Round(time.Second).String()
}
