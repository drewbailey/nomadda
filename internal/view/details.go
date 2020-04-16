package view

import (
	"github.com/derailed/tview"
	"github.com/gdamore/tcell"
)

type Details struct {
	*tview.TextView

	app *App
}

func NewDetails(app *App) *Details {
	return &Details{
		TextView: tview.NewTextView(),
		app:      app,
	}
}

func (d *Details) Init() {
	d.SetBorder(true)
	d.SetBackgroundColor(tcell.ColorDefault)
	d.SetScrollable(true).SetWrap(true).SetRegions(true)
	d.SetDynamicColors(true)
	d.SetHighlightColor(tcell.ColorOrange)
	d.SetTitleColor(tcell.ColorGreenYellow)
	d.SetBorderPadding(0, 0, 1, 1)
	d.SetChangedFunc(func() {
		d.app.Draw()
	})

	// TODO handle keyboard capture
	// d.SetInputCapture()

}
