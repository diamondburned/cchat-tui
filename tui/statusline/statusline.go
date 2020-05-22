package statusline

import (
	"github.com/diamondburned/cchat-tui/tui/log"
	"github.com/diamondburned/cchat-tui/tui/statusline/mode"
	"github.com/diamondburned/cchat-tui/tui/ti"
	"github.com/rivo/tview"
)

type Container struct {
	*tview.Flex
	Mode *mode.Indicator
	Line *tview.Pages // contains Cmd and Log

	Cmd *tview.InputField // hidden unless CommandMode
	Log *log.OneLiner
}

func NewContainer(modestate *mode.State, d ti.Drawer) *Container {
	indicator := mode.NewIndicator(modestate, d)

	// Logger line that shares a global buffer.
	log := log.NewOneLiner(d)
	// Command input, hidden by default.
	cmd := tview.NewInputField()
	cmd.SetFieldBackgroundColor(-1)

	line := tview.NewPages()
	line.SetBackgroundColor(-1)
	line.AddPage("cmd", cmd, true, false)
	line.AddPage("log", log, true, true)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexColumn)
	flex.SetBackgroundColor(-1)
	flex.AddItem(indicator, 14, 1, true)
	flex.AddItem(line, 0, 1, true)

	// Calculate the mode's width after changing mode.
	modestate.OnChange(func(m mode.Mode) {
		flex.ResizeItem(indicator, indicator.Width(), 1)

		// Set the input visible if we're in command mode.
		if m == mode.Command {
			line.SwitchToPage("cmd")
		} else {
			line.SwitchToPage("log")
		}
	})

	return &Container{
		flex,
		indicator,
		line,
		cmd,
		log,
	}
}

// Focus calls delegate() on the command input. Only focus if the mode is
// Command.
func (c *Container) Focus(delegate func(tview.Primitive)) {
	delegate(c.Cmd)
}

// Blur empties the command input.
func (c *Container) Blur() {
	c.Cmd.SetText("")
}
