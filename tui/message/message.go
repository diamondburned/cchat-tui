package message

import (
	"bytes"
	"fmt"
	"time"

	"github.com/diamondburned/cchat"
	"github.com/diamondburned/cchat-tui/tui/app"
	"github.com/diamondburned/cchat-tui/tui/humanize"
	"github.com/diamondburned/cchat-tui/tui/log"
	"github.com/diamondburned/cchat/text"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

type MessageView struct {
	*tview.Flex
	Container *Container
	Input     *tview.InputField

	current cchat.ServerMessage
	sender  cchat.ServerMessageSender
}

func NewMessageView() *MessageView {
	flex := tview.NewFlex()
	flex.SetBackgroundColor(-1)
	flex.SetDirection(tview.FlexRow)

	container := NewContainer()
	input := tview.NewInputField()
	input.SetFieldBackgroundColor(-1)
	input.SetPlaceholder("Message...")

	flex.AddItem(container, 0, 1, true)
	flex.AddItem(input, 1, 1, true)

	return &MessageView{
		Flex:      flex,
		Container: container,
		Input:     input,
	}
}

// JoinServer is not thread-safe.
func (v *MessageView) JoinServer(server cchat.ServerMessage) {
	if v.current != nil {
		if err := v.current.LeaveServer(); err != nil {
			log.Error(errors.Wrap(err, "Error leaving server"))
		}
	}

	var name string
	if s, ok := server.(cchat.Server); ok {
		name, _ = s.Name()
	}

	v.current = server

	if s, ok := server.(cchat.ServerMessageSender); ok {
		v.sender = s
		// Only allow inputs if the server allows sending.
		v.Input.SetAcceptanceFunc(acceptAll)
		v.Input.SetPlaceholder("Message #" + name)
	} else {
		v.sender = nil
		v.Input.SetAcceptanceFunc(acceptNone)
		v.Input.SetPlaceholder("You cannot send messages here.")
	}

	if err := v.current.JoinServer(v.Container); err != nil {
		log.Error(errors.Wrap(err, "Failed to join server"))
	}
}

// SendMessage is not thread-safe.
func (v *MessageView) SendMessage() {
	if v.sender == nil {
		return
	}

	var server = v.sender

	var send = SendMessage(v.Input.GetText())
	if send == "" {
		return
	}

	v.Input.SetText("")

	go func() {
		if err := server.SendMessage(send); err != nil {
			log.Error(errors.Wrap(err, "Failed to send message"))
		}
	}()
}

type SendMessage string

func (s SendMessage) Content() string { return string(s) }

func acceptNone(textToCheck string, lastChar rune) bool {
	return false
}
func acceptAll(textToCheck string, lastChar rune) bool {
	return true
}

type Container struct {
	*tview.TextView
	Messages []Message

	focused   int
	renderBuf bytes.Buffer
}

func NewContainer() *Container {
	text := tview.NewTextView()
	text.SetBackgroundColor(-1)
	text.SetToggleHighlights(true)
	text.SetRegions(true)
	text.SetDynamicColors(true)
	text.SetScrollable(true)

	return &Container{TextView: text, focused: -1}
}

// FocusMessage is not thread-safe.
func (c *Container) FocusMessage(i int) bool {
	if i >= len(c.Messages) || i < 0 {
		return false
	}

	c.TextView.Highlight(c.Messages[i].ID)
	c.focused = i
	return true
}

// UnfocusMessage is not thread-safe.
func (c *Container) UnfocusMessage() {
	c.focused = -1
	c.TextView.Highlight()
}

func (c *Container) findByID(id string) (index int) {
	for i, m := range c.Messages {
		if m.ID == id {
			return i
		}
	}
	return -1
}

func (c *Container) rerender() {
	for _, msg := range c.Messages {
		c.renderBuf.WriteString(msg.Render())
		c.renderBuf.WriteByte('\n')
	}
	c.TextView.SetText(c.renderBuf.String())
	c.renderBuf.Reset()

	// Resture selection.
	if c.focused > -1 {
		c.FocusMessage(c.focused)
	}
}

// Reset is not thread-safe.
func (c *Container) Reset() {
	c.Messages = nil
	c.TextView.SetText("")
}

// CreateMessage is thread-safe.
func (c *Container) CreateMessage(msg cchat.MessageCreate) {
	var msgc = NewMessage(msg)
	app.QueueUpdateDraw(func() {
		c.Messages = append(c.Messages, msgc)
		// lazy render
		c.TextView.Write([]byte(msgc.Render() + "\n"))
	})
}

// UpdateMessage is thread-safe.
func (c *Container) UpdateMessage(msg cchat.MessageUpdate) {
	app.QueueUpdateDraw(func() {
		// Find the message.
		i := c.findByID(msg.ID())
		if i < 0 {
			return
		}

		// TODO: edited timestamp

		if author := msg.Author(); !author.Empty() {
			c.Messages[i].UpdateAuthor(author)
		}
		if content := msg.Content(); !content.Empty() {
			c.Messages[i].UpdateContent(content)
		}

		// Re-render the entire buffer.
		c.rerender()
	})
}

// DeleteMessage is thread-safe.
func (c *Container) DeleteMessage(msg cchat.MessageDelete) {
	app.QueueUpdateDraw(func() {
		i := c.findByID(msg.ID())
		if i < 0 {
			return
		}

		c.Messages = append(c.Messages[:i], c.Messages[i+1:]...)
		c.rerender()
	})
}

type Message struct {
	ID    string
	Nonce string // TODO

	Timestamp string
	Author    string
	Content   string
}

func NewMessage(msg cchat.MessageCreate) Message {
	m := Message{ID: msg.ID()}
	m.UpdateTimestamp(msg.Time())
	m.UpdateAuthor(msg.Author())
	m.UpdateContent(msg.Content())

	return m
}

func (m *Message) UpdateTimestamp(t time.Time) {
	m.Timestamp = renderTimestamp(t)
}
func (m *Message) UpdateAuthor(author text.Rich) {
	m.Author = renderAuthor(author)
}
func (m *Message) UpdateContent(content text.Rich) {
	m.Content = renderContent(content)
}
func (m Message) Render() string {
	return fmt.Sprintf(`["%s"]%s %s: %s`, m.ID, m.Timestamp, m.Author, m.Content)
}

func renderTimestamp(t time.Time) string {
	return `[gray]` + humanize.TimeAgo(t) + "[-]"
}

func renderAuthor(txt text.Rich) string {
	// TODO: variable colors
	const color = "aqua"

	return "[" + color + "]" + tview.Escape(txt.Content) + "[-]"
}

func renderContent(txt text.Rich) string {
	return tview.Escape(txt.Content)
}
