package mode

import (
	"github.com/diamondburned/cchat-tui/tui/log"
	"github.com/diamondburned/cchat-tui/tui/ti"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Mode uint8

const (
	Normal Mode = iota
	Insert
	Visual
	Command
)

func ModeFromEventKey(evk *tcell.EventKey) (Mode, bool) {
	// If modifiers are pressed, then we shouldn't handle it.
	if evk.Modifiers() != tcell.ModNone {
		return 0, false
	}

	switch evk.Key() {
	case tcell.KeyESC:
		return Normal, true
	case tcell.KeyRune:
		// no-op
	default:
		return 0, false
	}

	switch evk.Rune() {
	case 'i':
		return Insert, true
	case 'v':
		return Visual, true
	case ';':
		return Command, true
	}

	return 0, false
}

func (m Mode) Valid() bool {
	return m == 0
}

func (m Mode) String() string {
	switch m {
	case Normal:
		return "Normal"
	case Insert:
		return "Insert"
	case Visual:
		return "Visual"
	case Command:
		return "Command"
	}
	return ""
}

func (m Mode) DisplayString() string {
	switch m {
	case Normal:
		return "-- NORMAL -- "
	case Insert:
		return "-- INSERT -- "
	case Visual:
		return "-- VISUAL -- "
	case Command:
		return ";"
	}
	return ""
}

func (m Mode) DisplayWidth() int {
	return len(m.DisplayString())
}

type State struct {
	Mode     Mode
	onChange []func(Mode)
}

func NewState() *State {
	return &State{Mode: Normal}
}

// OnChange appends a mode handler. This function is not thread-safe. Calls on
// this function, however, will be inside the tview thread.
func (state *State) OnChange(fn func(Mode)) {
	state.onChange = append(state.onChange, fn)
	fn(state.Mode)
}

// BindApplication binds to the application a mode shortcut handler.
func (state *State) BindApplication(app *tview.Application) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Run this first since I don't trust type assertions to be fast.
		mo, ok := ModeFromEventKey(event)
		if !ok {
			return event
		}

		// If we're focused on an input handler, then we shouldn't handle
		// anything but an Escape.
		if _, ok := app.GetFocus().(*tview.InputField); ok {
			if event.Key() != tcell.KeyESC {
				return event
			}
		}

		if mo == state.Mode {
			return event
		}

		state.changeMode(mo)
		return nil
	})
}

func (state *State) changeMode(mode Mode) {
	state.Mode = mode
	for _, fn := range state.onChange {
		fn(mode)
	}
}

type Indicator struct {
	*tview.TextView
	state *State
}

func NewIndicator(s *State, d ti.Drawer) *Indicator {
	tv := tview.NewTextView()
	tv.SetDynamicColors(false)
	tv.SetScrollable(false)
	tv.SetRegions(false)
	tv.SetBackgroundColor(-1)
	tv.SetWrap(false)
	tv.SetToggleHighlights(false)

	s.OnChange(func(mode Mode) {
		log.Write("Indicator changing")
		tv.SetText(mode.DisplayString())
	})

	return &Indicator{tv, s}
}

func (i *Indicator) Width() int {
	return i.state.Mode.DisplayWidth()
}
