package tui

import (
	"github.com/diamondburned/cchat-tui/tui/message"
	"github.com/diamondburned/cchat-tui/tui/service"
	"github.com/diamondburned/cchat-tui/tui/statusline"
	"github.com/diamondburned/cchat-tui/tui/statusline/mode"
	"github.com/diamondburned/cchat-tui/tui/ti"
	"github.com/rivo/tview"
)

type Application struct {
	*tview.Application
	Window    *Window
	ModeState *mode.State
}

func NewApplication() *Application {
	app := tview.NewApplication()
	modestate := mode.NewState()
	return &Application{
		Application: app,
		Window:      NewWindow(modestate, app),
		ModeState:   modestate,
	}
}

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
		*tview.Flex
		Services *service.Container
		Messages *message.MessageView
	}
	StatusLine *statusline.Container
}

func NewWindow(modestate *mode.State, d ti.Drawer) *Window {
	flex := tview.NewFlex()
	flex.SetBackgroundColor(-1)
	flex.SetDirection(tview.FlexRow)

	stline := statusline.NewContainer(modestate, d)

	win := &Window{
		Flex:       flex,
		StatusLine: stline,
	}

	// Make the top container
	win.Top.Flex = tview.NewFlex()
	win.Top.Flex.SetBackgroundColor(-1)
	win.Top.Flex.SetDirection(tview.FlexColumn)

	win.Top.Services = service.NewContainer()
	win.Top.Messages = nil // TODO

	win.Top.AddItem(win.Top.Services, 40, 1, false)
	win.Top.AddItem(win.Top.Messages, 0, 1, false)

	flex.AddItem(win.Top, 0, 1, false)
	flex.AddItem(win.StatusLine, 1, 1, false)

	return win
}
