package statusline

import (
	"github.com/diamondburned/cchat-tui/tui/log"
	"github.com/diamondburned/cchat-tui/tui/statusline/mode"
	"github.com/diamondburned/cchat-tui/tui/ti"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Container struct {
	*tview.Flex
	Mode *mode.Indicator
	Log  *log.OneLiner
}

func NewContainer(modestate *mode.State, d ti.Drawer) *Container {
	mode := mode.NewIndicator(modestate, d)
	logline := log.NewOneLiner(d)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexColumn)
	flex.SetBackgroundColor(-1)
	flex.AddItem(mode, 14, 1, true)
	flex.AddItem(logline, 0, 1, true)

	// Calculate the mode's width before redrawing.
	flex.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (ix, iy, iw, ih int) {
		// Calculate the mode's width.
		flex.ResizeItem(mode, mode.Width(), 1)
		return x, y, w, h
	})

	return &Container{
		flex,
		mode,
		logline,
	}
}
