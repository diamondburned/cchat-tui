package log

import (
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/diamondburned/cchat-tui/tui/ti"
	"github.com/gdamore/tcell"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

var LogBuffer = newLogBuffer()

type logBuffer struct {
	mutex   sync.Mutex
	entries []LogEntry
	onEntry []func(LogEntry)
}

type LogEntry struct {
	Time time.Time
	Msg  string
}

func NewLogToWriter(w io.Writer) func(LogEntry) {
	return func(entry LogEntry) {
		w.Write([]byte(entry.Time.String() + ": " + entry.Msg))
	}
}

func NewLogToFile(file string) (func(LogEntry), error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0750)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open/create file")
	}
	return NewLogToWriter(f), nil
}

func newLogBuffer() *logBuffer {
	return &logBuffer{}
}

func Error(err error) {
	if err != nil {
		Write("Error: " + err.Error())
	}
}
func Write(msg string) {
	LogBuffer.Write(msg)
}
func OnEntry(fn func(LogEntry)) {
	LogBuffer.OnEntry(fn)
}

func (l *logBuffer) Write(msg string) {
	entry := LogEntry{
		Time: time.Now(),
		Msg:  msg,
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.entries = append(l.entries, entry)
	l.callEntry(entry)
}

func (l *logBuffer) OnEntry(fn func(LogEntry)) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.onEntry = append(l.onEntry, fn)
}

func (l *logBuffer) callEntry(entry LogEntry) {
	for _, fn := range l.onEntry {
		fn(entry)
	}
}

type OneLiner struct {
	*tview.TextView
	drawer ti.Drawer
}

func NewOneLiner(d ti.Drawer) (l *OneLiner) {
	tv := tview.NewTextView()
	tv.SetWrap(false)
	tv.SetBackgroundColor(-1)
	tv.SetDynamicColors(false)

	l = &OneLiner{tv, d}
	OnEntry(l.OnEntry)
	return
}

// OnEntry is thread-safe, as it will be called by the logger.
func (l *OneLiner) OnEntry(entry LogEntry) {
	msg := strings.Replace(entry.Msg, "\n", "↵", -1)

	l.drawer.QueueUpdateDraw(func() {
		var color = tcell.Color(-1)
		if strings.HasPrefix(msg, "Error:") {
			color = tcell.ColorRed
		}

		l.SetTextColor(color)
		l.SetText(msg)
	})
}

// Clear is not thread-safe. It should be wrapped with QueueUpdateDraw.
func (l *OneLiner) Clear() {
	l.SetText("")
}
