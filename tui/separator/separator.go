package separator

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Vertical struct {
	*tview.Box
	Style tcell.Style
}

const VSeparator = '▏'

func NewVertical() (v *Vertical) {
	v = &Vertical{
		Box:   tview.NewBox(),
		Style: tcell.StyleDefault.Dim(true),
	}
	v.Box.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (ix, iy, iw, ih int) {
		for i := y; i < y+h; i++ {
			screen.SetContent(x, i, VSeparator, nil, v.Style)
		}
		return v.Box.GetInnerRect()
	})
	return
}
