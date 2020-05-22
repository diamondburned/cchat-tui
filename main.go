package main

import (
	"os"
	"path/filepath"

	"github.com/diamondburned/cchat-tui/tui"
	"github.com/diamondburned/cchat-tui/tui/log"
)

func main() {
	// TODO: make a plugin repository

	// logging
	f, err := log.NewLogToFile(filepath.Join(os.TempDir(), "cchat-tui.log"))
	if err == nil {
		log.OnEntry(f)
	}

	app := tui.NewApplication()
	app.Run()
	app.Stop()
}
