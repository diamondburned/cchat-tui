package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/diamondburned/cchat-tui/tui/app"
	"github.com/gdamore/tcell"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

var LogBuffer = newLogBuffer()

type logBuffer struct {
	mutex   sync.Mutex
	entries []LogEntry
	onEntry []func(LogEntry)
	stderr  bool
}

type LogEntry struct {
	Time time.Time
	Msg  string
}

func (entry LogEntry) String() string {
	return entry.Time.Format(time.Stamp) + ": " + entry.Msg
}

func NewLogToWriter(w io.Writer) func(LogEntry) {
	return func(entry LogEntry) {
		w.Write([]byte(entry.String() + "\n"))
	}
}

func NewLogToFile(file string) (func(LogEntry), error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_SYNC, 0750)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open/create file")
	}
	return NewLogToWriter(f), nil
}

func newLogBuffer() *logBuffer {
	return &logBuffer{stderr: true}
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

func SetStderr(stderr bool) {
	LogBuffer.SetStderr(stderr)
}

func (l *logBuffer) Write(msg string) {
	entry := LogEntry{
		Time: time.Now(),
		Msg:  msg,
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.stderr {
		fmt.Fprintln(os.Stderr, entry.String())
	}

	l.entries = append(l.entries, entry)

	for _, fn := range l.onEntry {
		fn(entry)
	}
}

func (l *logBuffer) SetStderr(stderr bool) {
	l.mutex.Lock()
	l.stderr = stderr
	l.mutex.Unlock()
}

func (l *logBuffer) OnEntry(fn func(LogEntry)) {
	l.mutex.Lock()
	l.onEntry = append(l.onEntry, fn)
	l.mutex.Unlock()
}

type OneLiner struct {
	*tview.TextView
}

func NewOneLiner() (l *OneLiner) {
	tv := tview.NewTextView()
	tv.SetWrap(false)
	tv.SetBackgroundColor(-1)
	tv.SetDynamicColors(false)

	l = &OneLiner{tv}
	OnEntry(l.OnEntry)
	return
}

// OnEntry is thread-safe, as it will be called by the logger.
func (l *OneLiner) OnEntry(entry LogEntry) {
	msg := strings.Replace(entry.Msg, "\n", "↵", -1)

	// Prevent deadlocks.
	go app.QueueUpdateDraw(func() {
		var color = tcell.Color(-1)
		if strings.HasPrefix(msg, "Error:") {
			color = tcell.ColorRed
		}

		parts := strings.Split(msg, ": ")
		msg = parts[len(parts)-1]

		l.SetTextColor(color)
		l.SetText(msg)
	})
}

// Clear is not thread-safe. It should be wrapped with QueueUpdateDraw.
func (l *OneLiner) Clear() {
	l.SetText("")
}
