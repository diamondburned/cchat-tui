package app

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var app *tview.Application
var pages *tview.Pages // global navigator

func init() {
	pages = tview.NewPages()
	pages.SetBackgroundColor(-1)

	app = tview.NewApplication()
	app.EnableMouse(false)
	app.SetBeforeDrawFunc(func(s tcell.Screen) bool {
		s.Clear()
		return false
	})
	app.SetRoot(pages, true)
}

func Init(main tview.Primitive) *tview.Application {
	pages.AddPage("main", main, true, true)
	return app
}

func SwitchToView(p tview.Primitive) {
	pages.AddAndSwitchToPage("view", p, true)
}

func SwitchToMain() {
	pages.SwitchToPage("main")
	pages.RemovePage("view") // unreference
}

func Draw() {
	app.Draw()
}

func ForceDraw() {
	app.ForceDraw()
}

func QueueUpdateDraw(fn func()) {
	app.QueueUpdateDraw(fn)
}

func SetFocus(p tview.Primitive) {
	app.SetFocus(p)
}

func GetFocus() tview.Primitive {
	return app.GetFocus()
}
