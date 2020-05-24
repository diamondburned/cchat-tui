package tui

import (
	"github.com/diamondburned/cchat-tui/tui/app"
	"github.com/diamondburned/cchat-tui/tui/message"
	"github.com/diamondburned/cchat-tui/tui/service"
	"github.com/diamondburned/cchat-tui/tui/statusline"
	"github.com/diamondburned/cchat-tui/tui/statusline/mode"
	"github.com/rivo/tview"
)

type Window struct {
	// Contains
	//    message view
	//    ---
	//    modeline
	*tview.Flex

	Top struct {
		// Contains
		//    services  |  messages
		//              |
		*tview.Grid
		Services *service.Container
		Messages *message.MessageView
	}
	StatusLine *statusline.Container
}

func NewWindow(modestate *mode.State) *Window {
	flex := tview.NewFlex()
	flex.SetBackgroundColor(-1)
	flex.SetDirection(tview.FlexRow)

	stline := statusline.NewContainer(modestate)

	win := &Window{
		Flex:       flex,
		StatusLine: stline,
	}

	// Make the top container
	win.Top.Grid = tview.NewGrid()
	win.Top.Grid.SetBorders(false)
	win.Top.Grid.SetBackgroundColor(-1)

	win.Top.Services = service.NewContainer()
	win.Top.Messages = message.NewMessageView()

	// Set the proper column sizes.
	// - 25ch for the left column on medium layout.
	// - 40ch for the left column on large layout.
	// - The rest is for messages
	win.Top.SetColumns(25, 15, -1)

	// Smallest layout.
	win.Top.AddItem(win.Top.Services, 0, 0, 0, 0, 0, 0, true)
	win.Top.AddItem(win.Top.Messages, 0, 0, 1, 3, 0, 0, true)

	// Medium layout, columns > 70.
	win.Top.AddItem(win.Top.Services, 0, 0, 1, 1, 0, 70, true)
	win.Top.AddItem(win.Top.Messages, 0, 1, 1, 2, 0, 70, true)

	// Large layout, columns > 150.
	win.Top.AddItem(win.Top.Services, 0, 0, 1, 2, 0, 150, true)
	win.Top.AddItem(win.Top.Messages, 0, 2, 1, 1, 0, 150, true)

	flex.AddItem(win.Top, 0, 1, true)
	flex.AddItem(win.StatusLine, 1, 1, true)

	modestate.OnChange(func(m mode.Mode) {
		switch m {
		case mode.Normal:
			app.SetFocus(win.Top.Services)
		case mode.Insert:
			app.SetFocus(win.Top.Messages.Input)
		case mode.Visual:
			app.SetFocus(win.Top.Messages.Container)
		case mode.Command:
			app.SetFocus(win.StatusLine)
		}
	})

	return win
}
