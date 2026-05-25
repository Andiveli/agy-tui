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
	viewport    viewport.Model
	input       textinput.Model
	messages    []kit.ChatMessage
	loading     bool
	styles      kit.Styles
	client      *backend.Client
	sessionName string
	promptCount int
	streamBuf   strings.Builder
	width       int
	height      int
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
		viewport: vp,
		input:    ti,
		styles:   styles,
		client:   client,
		width:    80,
		height:   24,
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
			c.input.SetValue("")
			c.messages = append(c.messages, kit.ChatMessage{Role: "user", Content: prompt})
			c.loading = true
			c.promptCount++
			c.streamBuf.Reset()
			c.updateViewport()

			cmds := []tea.Cmd{c.emitProgress("running", 0), c.startStream(prompt)}

			// On first prompt, derive session name
			if c.promptCount == 1 {
				c.sessionName = slugify(prompt)
				cmds = append(cmds, c.emitSession(c.sessionName, ""))
			}

			return c, tea.Batch(cmds...)

		case "ctrl+c":
			if c.input.Value() != "" {
				c.input.SetValue("")
				return c, nil
			}
			return c, tea.Quit
		}

	case *backend.StreamChunk:
		// Append line to buffer and show in viewport
		c.streamBuf.WriteString(msg.Text + "\n")
		c.messages = append(c.messages, kit.ChatMessage{Role: "agent", Content: msg.Text + "\n"})
		c.updateViewport()
		// Schedule reading the next chunk
		return c, c.readNextChunk(msg.Scanner, msg.Reader)

	case kit.ChatCompletedMsg:
		c.loading = false
		// Replace the individual stream chunks with the final content
		c.removeLastAgentMessages()
		c.messages = append(c.messages, kit.ChatMessage{Role: "agent", Content: msg.Content})
		c.updateViewport()

		return c, tea.Batch(
			c.emitProgress("completed", 100),
			c.emitSession(msg.SessionName, ""),
			c.emitFileChanges(msg.FilePaths),
		)

	case kit.ChatErrorMsg:
		c.loading = false
		errMsg := fmt.Sprintf("Error: %s", msg.Err.Error())
		c.messages = append(c.messages, kit.ChatMessage{Role: "agent", Content: errMsg})
		c.updateViewport()
		return c, c.emitProgress("failed", 0)
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
		inputLine = c.styles.Dimmed.Render(" Waiting for response...")
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

// startStream launches agy and reads the first chunk.
func (c *ChatModel) startStream(prompt string) tea.Cmd {
	return func() tea.Msg {
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
		scanner := bufio.NewScanner(reader)
		if scanner.Scan() {
			return &backend.StreamChunk{
				Text:    scanner.Text(),
				Scanner: scanner,
				Reader:  reader,
			}
		}
		// Empty response
		reader.Close()
		return kit.ChatCompletedMsg{
			Content:     "",
			SessionName: c.sessionName,
		}
	}
}

// readNextChunk reads the next line from the scanner and returns a message.
func (c *ChatModel) readNextChunk(scanner *bufio.Scanner, reader io.ReadCloser) tea.Cmd {
	return func() tea.Msg {
		if scanner.Scan() {
			return &backend.StreamChunk{
				Text:    scanner.Text(),
				Scanner: scanner,
				Reader:  reader,
			}
		}
		reader.Close()
		content := c.streamBuf.String()
		return kit.ChatCompletedMsg{
			Content:     content,
			SessionName: c.sessionName,
			FilePaths:   parseFilePaths(content),
		}
	}
}

func (c *ChatModel) emitProgress(status string, progress int) tea.Cmd {
	return func() tea.Msg {
		return kit.ProgressMsg{SubAgent: "agy", Status: status, Progress: progress}
	}
}

func (c *ChatModel) emitSession(name, context string) tea.Cmd {
	return func() tea.Msg {
		return kit.SessionChangedMsg{Name: name, Context: context}
	}
}

func (c *ChatModel) emitFileChanges(paths []string) tea.Cmd {
	return func() tea.Msg {
		// Only send the first file to avoid flooding
		if len(paths) > 0 {
			return kit.FileChangedMsg{Path: paths[0], Action: "modified"}
		}
		return nil
	}
}

func (c *ChatModel) removeLastAgentMessages() {
	// Remove streaming chunk messages from the end
	for i := len(c.messages) - 1; i >= 0; i-- {
		if c.messages[i].Role == "user" {
			break
		}
		c.messages = c.messages[:i]
	}
}

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
	if c.loading && c.streamBuf.Len() == 0 {
		rendered = append(rendered, c.styles.Dimmed.Render(" Sending..."))
	}
	content := lipgloss.JoinVertical(lipgloss.Top, rendered...)
	c.viewport.SetContent(content)
	c.viewport.GotoBottom()
}

// slugify derives a session name from the first ~4 words of a prompt.
func slugify(prompt string) string {
	// Split into words, take first 4
	words := strings.Fields(prompt)
	if len(words) > 4 {
		words = words[:4]
	}
	// Lowercase and join with hyphens
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
