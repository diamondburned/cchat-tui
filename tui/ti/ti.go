// Package ti provides interfaces around tview.
package ti

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Drawer interface {
	Draw() *tview.Application
	ForceDraw() *tview.Application
	QueueUpdate(func()) *tview.Application
	QueueUpdateDraw(func()) *tview.Application
}

var _ Drawer = (*tview.Application)(nil)

type InputHandlerFunc = func(event *tcell.EventKey, setFocus func(tview.Primitive))
