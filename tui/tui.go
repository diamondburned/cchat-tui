package tui

import (
	"github.com/diamondburned/cchat"
	"github.com/diamondburned/cchat-tui/tui/message"
	"github.com/diamondburned/cchat-tui/tui/service"
	"github.com/diamondburned/cchat-tui/tui/statusline"
	"github.com/diamondburned/cchat-tui/tui/statusline/mode"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Application struct {
	*tview.Application
	Window    *Window
	ModeState *mode.State
}

func NewApplication() *Application {
	app := tview.NewApplication()
	app.EnableMouse(false)
	app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		screen.Clear()
		return false
	})

	modestate := mode.NewState()
	win := NewWindow(modestate, app)

	app.SetRoot(win, true)

	return &Application{
		app,
		NewWindow(modestate, app),
		modestate,
	}
}

func (app *Application) AddService(sv cchat.Service) {
	app.Window.Top.Services.AddService(sv, app.Application)
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
		*tview.Grid
		Services *service.Container
		Messages *message.MessageView
	}
	StatusLine *statusline.Container
}

func NewWindow(modestate *mode.State, app *tview.Application) *Window {
	flex := tview.NewFlex()
	flex.SetBackgroundColor(-1)
	flex.SetDirection(tview.FlexRow)

	stline := statusline.NewContainer(modestate, app)

	win := &Window{
		Flex:       flex,
		StatusLine: stline,
	}

	// Make the top container
	win.Top.Grid = tview.NewGrid()
	win.Top.Grid.SetBackgroundColor(-1)

	win.Top.Services = service.NewContainer()
	win.Top.Messages = message.NewMessageView()

	// Set the minimum widths. 25px for each column.
	win.Top.SetMinSize(0, 25)
	// Set the proper column sizes.
	// - 25px for the left column on medium layout.
	// - 35px for the left column on large layout.
	// - The rest is for messages
	win.Top.SetColumns(25, 35, 0)

	// Smallest layout.
	win.Top.AddItem(win.Top.Services, 0, 0, 0, 0, 0, 0, true)
	win.Top.AddItem(win.Top.Messages, 0, 0, 1, 1, 0, 0, true)

	// Medium layout, columns > 70.
	win.Top.AddItem(win.Top.Services, 0, 0, 1, 1, 0, 70, true)
	win.Top.AddItem(win.Top.Messages, 0, 1, 1, 1, 0, 70, true)

	// Large layout, columns > 150.
	win.Top.AddItem(win.Top.Services, 0, 1, 1, 1, 0, 150, true)
	win.Top.AddItem(win.Top.Messages, 0, 2, 1, 1, 0, 150, true)

	flex.AddItem(win.Top, 0, 1, true)
	flex.AddItem(win.StatusLine, 1, 1, true)

	modestate.BindApplication(app)
	modestate.OnChange(func(m mode.Mode) {
		switch m {
		case mode.Normal:
			app.SetFocus(win.Top.Services)
		case mode.Insert:
			app.SetFocus(win.Top.Messages.Input)
		case mode.Visual:
			app.SetFocus(win.Top.Messages.Container)
		}

		// TODO: command mode
	})

	return win
}
