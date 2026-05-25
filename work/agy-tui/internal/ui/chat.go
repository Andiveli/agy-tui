package ui

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/samael/agy-tui/internal/backend"
	"github.com/samael/agy-tui/internal/ui/kit"
)

// ChatModel with streaming support.
type ChatModel struct {
	viewport      viewport.Model
	input         textinput.Model
	messages      []kit.ChatMessage
	loading       bool
	styles        kit.Styles
	client        *backend.Client
	sessionMgr    *backend.SessionManager
	sessionName   string
	promptCount   int
	streamingText string         // text accumulated during active streaming (shown in real-time)
	streamReader  io.ReadCloser  // active stream pipe (set while loading, nil otherwise)
	streamScanner *bufio.Scanner // persistent scanner across chunks (avoids losing buffered data)
	width         int
	height        int
}

var filePathPattern = regexp.MustCompile(`([a-zA-Z0-9_/.-]+\.(go|ts|tsx|js|jsx|py|rs|css|md|json|yaml|yml|mod|sum))`)

func NewChatModel(styles kit.Styles, client *backend.Client) *ChatModel {
	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1)
	ti := textinput.New()
	ti.Placeholder = "Type a prompt..."
	ti.PromptStyle = styles.InputPrompt
	ti.TextStyle = styles.AgentMessage
	ti.Width = 60
	ti.Focus()
	return &ChatModel{
		viewport:   vp,
		input:      ti,
		styles:     styles,
		client:     client,
		sessionMgr: backend.NewSessionManager(),
		width:      80,
		height:     24,
	}
}

func (c *ChatModel) Init() tea.Cmd { return textinput.Blink }

func (c *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		c.viewport.Width = msg.Width - 2
		c.viewport.Height = msg.Height - 4
		c.input.Width = msg.Width - 6
		return c, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if c.loading {
				return c, nil
			}
			prompt := c.input.Value()
			if prompt == "" {
				return c, nil
			}
			// Slash commands
			if prompt == "/quit" || prompt == "/exit" || prompt == "/q" {
				return c, tea.Quit
			}
			c.input.SetValue("")
			c.messages = append(c.messages, kit.ChatMessage{Role: "user", Content: prompt})
			c.loading = true
			c.promptCount++
			c.streamingText = ""
			c.updateViewport()

			cmds := []tea.Cmd{c.emitProgress("running", 0), c.startStream(prompt)}

			// On first prompt, derive session name
			if c.promptCount == 1 {
				c.sessionName = slugify(prompt)
				c.sessionMgr.StartSession()
				cmds = append(cmds, c.emitSession(c.sessionName, ""))
			}

			return c, tea.Batch(cmds...)

		case "ctrl+c":
			// If streaming, cancel the current request
			if c.loading {
				c.loading = false
				c.closeStreamReader()
				c.streamingText = ""
				c.messages = append(c.messages, kit.ChatMessage{Role: "agent", Content: "⚠️ Cancelled"})
				c.updateViewport()
				return c, tea.Batch(c.emitProgress("failed", 0), c.emitDisconnected())
			}
			if c.input.Value() != "" {
				c.input.SetValue("")
				return c, nil
			}
			return c, tea.Quit
		}

	// StreamReadyMsg arrives from startStream goroutine — stores reader/scanner on model.
	case kit.StreamReadyMsg:
		c.streamReader = msg.Reader
		c.streamScanner = bufio.NewScanner(msg.Reader)
		if c.streamScanner.Scan() {
			c.streamingText = c.streamScanner.Text() + "\n"
			c.updateViewport()
			return c, c.readNextChunk()
		}
		// Empty response — clean up immediately
		c.closeStreamReader()
		c.loading = false
		return c, c.emitProgress("completed", 100)

	case kit.ChatStreamChunkMsg:
		// Append to streaming text and update viewport in real-time
		c.streamingText += msg.Text + "\n"
		c.updateViewport()
		return c, c.readNextChunk()

	case kit.ChatCompletedMsg:
		c.loading = false
		c.closeStreamReader()

		// Detect conversation ID for SessionManager
		if convs, err := c.sessionMgr.ListConversations(); err == nil && len(convs) > 0 {
			c.sessionMgr.SetConversationID(convs[0])
		}

		// Append a single clean message with the full content
		content := msg.Content
		if content == "" && c.streamingText != "" {
			content = c.streamingText
		}
		c.messages = append(c.messages, kit.ChatMessage{Role: "agent", Content: content})
		c.streamingText = ""
		c.updateViewport()

		return c, tea.Batch(
			c.emitProgress("completed", 100),
			c.emitSession(msg.SessionName, ""),
			c.emitFileChanges(msg.FilePaths),
		)

	case kit.ChatErrorMsg:
		c.loading = false
		c.closeStreamReader()
		c.streamingText = ""

		errMsg := fmt.Sprintf("⚠️ Error: %s", msg.Err.Error())
		c.messages = append(c.messages, kit.ChatMessage{Role: "agent", Content: errMsg})
		c.updateViewport()
		return c, tea.Batch(c.emitProgress("failed", 0), c.emitDisconnected())
	}

	var cmd tea.Cmd
	c.input, cmd = c.input.Update(msg)
	return c, cmd
}

func (c *ChatModel) View() string {
	header := c.styles.Title.Render(" Chat ")
	content := c.viewport.View()
	var inputLine string
	if c.loading {
		if c.streamingText != "" {
			inputLine = c.styles.InputPrompt.Render(" Streaming... (Ctrl+C to cancel)")
		} else {
			inputLine = c.styles.Dimmed.Render(" Sending...")
		}
	} else {
		inputLine = c.input.View()
	}
	inputBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#585b70")).
		PaddingLeft(1).
		Render(inputLine)
	return lipgloss.JoinVertical(lipgloss.Top, header, content, inputBox)
}

// ── Streaming ──────────────────────────────────────────────────────

// startStream checks binary, starts agy, and returns the pipe reader.
// It does NOT write to the model — that happens in the StreamReadyMsg handler (main goroutine).
func (c *ChatModel) startStream(prompt string) tea.Cmd {
	return func() tea.Msg {
		if err := c.client.CheckBinary(); err != nil {
			return kit.ChatErrorMsg{Err: fmt.Errorf(
				"agy CLI not found — install Antigravity CLI first: %w", err)}
		}

		ctx := context.Background()

		var reader io.ReadCloser
		var err error
		if c.promptCount > 1 {
			reader, err = c.client.ContinueLastStreaming(ctx, prompt)
		} else {
			reader, err = c.client.StartStreaming(ctx, prompt)
		}
		if err != nil {
			return kit.ChatErrorMsg{Err: err}
		}

		return kit.StreamReadyMsg{Reader: reader}
	}
}

// readNextChunk reads the next line from the active stream scanner.
func (c *ChatModel) readNextChunk() tea.Cmd {
	return func() tea.Msg {
		// streamScanner is set in the StreamReadyMsg handler (main goroutine), not in a cmd goroutine.
		if c.streamScanner == nil {
			return kit.ChatCompletedMsg{
				Content:     c.streamingText,
				SessionName: c.sessionName,
				FilePaths:   parseFilePaths(c.streamingText),
			}
		}

		if c.streamScanner.Scan() {
			return kit.ChatStreamChunkMsg{Text: c.streamScanner.Text()}
		}

		// Stream ended — return completion data without touching model state.
		// Cleanup happens in the ChatCompletedMsg handler (main goroutine).
		return kit.ChatCompletedMsg{
			Content:     c.streamingText,
			SessionName: c.sessionName,
			FilePaths:   parseFilePaths(c.streamingText),
		}
	}
}

// closeStreamReader closes the active stream pipe and clears scanner state.
func (c *ChatModel) closeStreamReader() {
	if c.streamReader != nil {
		c.streamReader.Close()
		c.streamReader = nil
	}
	c.streamScanner = nil
}

// ── Sidebar messages ──────────────────────────────────────────────

func (c *ChatModel) emitProgress(status string, progress int) tea.Cmd {
	return func() tea.Msg {
		return kit.ProgressMsg{SubAgent: "agy", Status: status, Progress: progress}
	}
}

func (c *ChatModel) emitSession(name, context string) tea.Cmd {
	return func() tea.Msg {
		msg := kit.SessionChangedMsg{Name: name, Context: context}
		if c.sessionMgr != nil {
			if convs, err := c.sessionMgr.ListConversations(); err == nil {
				msg.ConvCount = len(convs)
			}
			if cur := c.sessionMgr.CurrentSession(); cur != nil {
				msg.ConvID = cur.ID
			}
		}
		return msg
	}
}

func (c *ChatModel) emitFileChanges(paths []string) tea.Cmd {
	return func() tea.Msg {
		if len(paths) > 0 {
			return kit.FileChangedMsg{Path: paths[0], Action: "modified"}
		}
		return nil
	}
}

func (c *ChatModel) emitDisconnected() tea.Cmd {
	return func() tea.Msg {
		return kit.MCPStatusMsg{Connected: false, Status: "agy unreachable"}
	}
}

// ── Viewport rendering ────────────────────────────────────────────

func (c *ChatModel) updateViewport() {
	var rendered []string
	var agentBuf strings.Builder

	for _, m := range c.messages {
		switch m.Role {
		case "user":
			if agentBuf.Len() > 0 {
				rendered = append(rendered, c.styles.AgentMessage.Render(agentBuf.String()))
				agentBuf.Reset()
			}
			rendered = append(rendered, c.styles.UserMessage.Render("You: "+m.Content))
		case "agent":
			agentBuf.WriteString(m.Content)
		}
	}
	if agentBuf.Len() > 0 {
		rendered = append(rendered, c.styles.AgentMessage.Render(agentBuf.String()))
	}

	// Show streaming text in real-time without adding it to messages
	if c.loading && c.streamingText != "" {
		rendered = append(rendered, c.styles.AgentMessage.Render(c.streamingText))
	}

	content := lipgloss.JoinVertical(lipgloss.Top, rendered...)
	c.viewport.SetContent(content)
	c.viewport.GotoBottom()
}

// ── Helpers ───────────────────────────────────────────────────────

// slugify derives a session name from the first ~4 words of a prompt.
func slugify(prompt string) string {
	words := strings.Fields(prompt)
	if len(words) > 4 {
		words = words[:4]
	}
	var clean []string
	for _, w := range words {
		w = strings.ToLower(w)
		w = strings.TrimFunc(w, func(r rune) bool {
			return !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9')
		})
		if w != "" {
			clean = append(clean, w)
		}
	}
	if len(clean) == 0 {
		return "session"
	}
	return strings.Join(clean, "-")
}

func parseFilePaths(output string) []string {
	var paths []string
	seen := make(map[string]bool)
	for _, line := range strings.Split(output, "\n") {
		matches := filePathPattern.FindStringSubmatch(line)
		if len(matches) > 1 {
			path := matches[1]
			if !seen[path] {
				seen[path] = true
				paths = append(paths, path)
			}
		}
	}
	return paths
}
