package mode

import (
	"github.com/diamondburned/cchat-tui/tui/ti"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Mode uint8

const (
	_ Mode = iota
	NormalMode
	InsertMode
	VisualMode
	CommandMode
)

func ModeFromEventKey(evk *tcell.EventKey) Mode {
	// If modifiers are pressed, then we shouldn't handle it.
	if evk.Modifiers() != tcell.ModNone {
		return 0
	}

	switch evk.Key() {
	case tcell.KeyESC:
		return NormalMode
	case tcell.KeyRune:
		break
	default:
		return 0
	}

	switch evk.Rune() {
	case 'i':
		return InsertMode
	case 'v':
		return VisualMode
	case ';':
		return CommandMode
	}

	return 0
}

func (m Mode) Valid() bool {
	return m == 0
}

func (m Mode) String() string {
	switch m {
	case NormalMode:
		return "Normal"
	case InsertMode:
		return "Insert"
	case VisualMode:
		return "Visual"
	case CommandMode:
		return "Command"
	}
	return ""
}

func (m Mode) DisplayString() string {
	switch m {
	case NormalMode:
		return "-- NORMAL --"
	case InsertMode:
		return "-- INSERT --"
	case VisualMode:
		return "-- VISUAL --"
	case CommandMode:
		return ":"
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
	return &State{Mode: NormalMode}
}

// OnChange appends a mode handler. This function is not thread-safe. Calls on
// this function, however, will be inside the tview thread.
func (state *State) OnChange(fn func(Mode)) {
	state.onChange = append(state.onChange, fn)
	fn(state.Mode)
}

func (state *State) InputHandler() ti.InputHandlerFunc {
	return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if mode := ModeFromEventKey(event); mode > 0 {
			state.changeMode(mode)
		}
	}
}

func (state *State) changeMode(mode Mode) {
	state.Mode = mode
	for _, fn := range state.onChange {
		fn(mode)
	}
}

type Indicator tview.TextView

func NewIndicator(s *State, d ti.Drawer) *Indicator {
	tv := tview.NewTextView()
	tv.SetDynamicColors(false)
	tv.SetScrollable(false)
	tv.SetRegions(false)
	tv.SetTextColor(-1)
	tv.SetBackgroundColor(-1)
	tv.SetWrap(false)
	tv.SetToggleHighlights(false)

	s.OnChange(func(mode Mode) {
		tv.SetText(mode.DisplayString())
		d.ForceDraw()
	})

	return (*Indicator)(tv)
}
