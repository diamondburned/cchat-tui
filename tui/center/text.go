package center

import (
	"strings"

	"github.com/rivo/tview"
)

type Text struct {
	*Center
	TextView *tview.TextView
}

func NewText(text string) *Text {
	tv := tview.NewTextView()
	tv.SetBackgroundColor(-1)
	tv.SetDynamicColors(true)
	tv.SetTextAlign(tview.AlignCenter)

	c := New(tv)

	t := &Text{TextView: tv, Center: c}
	t.SetText(text)

	return t
}

func (t *Text) SetText(text string) {
	lines := strings.Split(text, "\n")
	// get the height
	t.Center.MaxHeight = len(lines)
	// find the max length
	t.Center.MaxWidth = 0
	for _, line := range lines {
		if len(line) > t.Center.MaxWidth {
			t.Center.MaxWidth = len(line)
		}
	}
	t.TextView.SetText(text)
}
