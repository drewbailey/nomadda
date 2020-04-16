package ui

import (
	"github.com/derailed/tview"
	"github.com/gdamore/tcell"
)

type App struct {
	*tview.Application
	Main  *Pages
	views map[string]tview.Primitive
}

func NewApp() *App {
	app := &App{
		Application: tview.NewApplication(),
		Main:        NewPages(),
		views: map[string]tview.Primitive{
			"logo": NewLogo(),
		},
	}
	// clear out to allow for default color
	app.Application.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		screen.Clear()
		return false
	})
	return app
}

func (a *App) Views() map[string]tview.Primitive {
	return a.views
}

func (a *App) Init() {
	a.Main.SetBackgroundColor(tcell.ColorDefault)
	// a.Main.AddPage("header", a.header(), true, true)
	a.SetRoot(a.Main, true)
}

func (a *App) Logo() *Logo {
	return a.views["logo"].(*Logo)
}
