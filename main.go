package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/diamondburned/cchat-tui/tui"
	"github.com/diamondburned/cchat-tui/tui/app"
	"github.com/diamondburned/cchat-tui/tui/log"
	"github.com/diamondburned/cchat-tui/tui/statusline/mode"
	"github.com/diamondburned/cchat/services"

	_ "github.com/diamondburned/cchat-mock"
	_ "github.com/diamondburned/cchat/services/plugins"
)

func main() {
	// TODO: make a plugin repository

	// logging
	f, err := log.NewLogToFile(filepath.Join(os.TempDir(), "cchat-tui.log"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open log file:", err)
	} else {
		log.OnEntry(f)
	}

	modestate := mode.NewState()
	win := tui.NewWindow(modestate)

	// Load plugins.
	srvcs, errs := services.Get()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
	}

	for _, srvc := range srvcs {
		win.Top.Services.AddService(srvc)
	}

	// Disable stderr before starting the TUI.
	log.SetStderr(false)

	app := app.Init(win)
	app.SetInputCapture(modestate.InputCapturer)
	app.Run()
	app.Stop()
}
